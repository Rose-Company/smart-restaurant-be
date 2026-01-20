package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (s *Service) GetTables(ctx context.Context, request *models.ListTablesRequest) (*models.BaseListResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(request.Page, request.PageSize)

	// Build filters
	filters := []repositories.Clause{}

	// Filter by status
	if request.Status != nil && *request.Status != "" && *request.Status != "all" {
		status := *request.Status
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("status = ?", status)
		})
	}

	// Filter by zone/location
	if request.Zone != nil && *request.Zone != "" && *request.Zone != "all" {
		zone := *request.Zone
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("location = ?", zone)
		})
	}

	// Search by table number
	if request.Search != nil && *request.Search != "" {
		search := "%" + strings.ToLower(*request.Search) + "%"
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("LOWER(table_number) LIKE ?", search)
		})
	}

	// Count total
	totalCount, err := s.tableRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []*models.TableWithOrderData{},
		}, nil
	}

	// Apply sorting
	queryParams := models.QueryParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	// Handle sorting
	if request.Sort != "" {
		switch request.Sort {
		case "tableNumber":
			queryParams.QuerySort.Origin = "table_number.asc"
		case "capacity":
			queryParams.QuerySort.Origin = "capacity.desc"
		case "recentlyCreated":
			queryParams.QuerySort.Origin = "created_at.desc"
		default:
			queryParams.QuerySort.Origin = "id.asc"
		}
	} else {
		queryParams.QuerySort.Origin = "id.asc"
	}

	// Get tables
	tables, err := s.tableRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	// Build response items
	items := make([]*models.TableWithOrderData, 0, len(tables))
	for _, table := range tables {
		item := &models.TableWithOrderData{
			ID:          table.ID,
			TableNumber: table.TableNumber,
			Capacity:    table.Capacity,
			Location:    table.Location,
			Status:      table.Status,
		}

		// If table is occupied, get order data
		if table.Status == "occupied" {
			orderData, err := s.getTableOrderData(ctx, table.ID)
			if err == nil && orderData != nil {
				item.OrderData = orderData
			}
		}

		items = append(items, item)
	}

	return &models.BaseListResponse{
		Total:    int(totalCount),
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// GetTablesForStaff - Get tables with orders for staff view (waiter/kitchen)
// Filter by is_help_needed and is_ready_to_bill, and only tables assigned to staff
func (s *Service) GetTablesForStaff(ctx context.Context, request *models.ListTablesForStaffRequest) (*models.BaseListResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(request.Page, request.PageSize)

	// First, get tables assigned to this staff member
	var assignedTableIDs []int
	if request.StaffID != "" {
		assignments, err := s.waiterTableAssignmentRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("waiter_id = ? AND is_active = true", request.StaffID).
				Select("table_id")
		})
		if err == nil && len(assignments) > 0 {
			for _, assignment := range assignments {
				assignedTableIDs = append(assignedTableIDs, assignment.TableID)
			}
		}
	}

	// If staff has no assigned tables, return empty
	if len(assignedTableIDs) == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []models.TableForStaffResponse{},
		}, nil
	}

	// Get occupied tables that are assigned to this staff
	allTables, err := s.tableRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("status = ? AND id IN ?", "occupied", assignedTableIDs)
	})
	if err != nil {
		return nil, err
	}

	if len(allTables) == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []models.TableForStaffResponse{},
		}, nil
	}

	// Extract table IDs for batch query
	tableIDs := make([]int, 0, len(allTables))
	for _, t := range allTables {
		tableIDs = append(tableIDs, t.ID)
	}

	// Get ALL orders for occupied tables with active status
	allOrders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("table_id IN ? ", tableIDs)
	})
	if err != nil {
		return nil, err
	}

	// Group orders by table_id
	ordersByTableID := make(map[int][]*models.Order)
	orderIDs := make([]int, 0)
	for _, order := range allOrders {
		ordersByTableID[order.TableID] = append(ordersByTableID[order.TableID], order)
		orderIDs = append(orderIDs, order.ID)
	}

	// Get ALL order items in ONE query
	var allOrderItems []*models.OrderItem
	if len(orderIDs) > 0 {
		allOrderItems, err = s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id IN ? AND status NOT IN ?", orderIDs, []string{"cancelled"})
		})
		if err != nil {
			return nil, err
		}
	}

	// Group order items by order_id
	itemsByOrderID := make(map[int][]*models.OrderItem)
	for _, item := range allOrderItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	// Filter tables based on order flags
	var filteredTableIDs []int
	for _, table := range allTables {
		orders := ordersByTableID[table.ID]
		if len(orders) == 0 {
			continue
		}

		// Check if table matches filter criteria
		shouldInclude := false

		if request.IsReadyToBill == nil && request.IsHelpNeeded == nil {
			// No filter, include all occupied tables with active orders
			shouldInclude = true
		} else if request.IsReadyToBill != nil && *request.IsReadyToBill {
			// Filter by is_ready_to_bill = true
			for _, order := range orders {
				if order.IsReadyToBill {
					shouldInclude = true
					break
				}
			}
		} else if request.IsHelpNeeded != nil && *request.IsHelpNeeded {
			// Filter by is_help_needed = true
			for _, order := range orders {
				if order.IsHelpNeeded {
					shouldInclude = true
					break
				}
			}
		} else if (request.IsReadyToBill != nil && !*request.IsReadyToBill) || (request.IsHelpNeeded != nil && !*request.IsHelpNeeded) {
			// Filter by both false (kitchen ready)
			allFalse := true
			for _, order := range orders {
				if order.IsReadyToBill || order.IsHelpNeeded {
					allFalse = false
					break
				}
			}
			if allFalse {
				shouldInclude = true
			}
		}

		if shouldInclude {
			filteredTableIDs = append(filteredTableIDs, table.ID)
		}
	}

	// Apply pagination
	totalCount := int64(len(filteredTableIDs))
	if totalCount == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []models.TableForStaffResponse{},
		}, nil
	}

	// Get paginated table IDs
	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(filteredTableIDs) {
		endIdx = len(filteredTableIDs)
	}

	paginatedTableIDs := filteredTableIDs[startIdx:endIdx]

	// Build response with detailed order information using pre-loaded data
	items := make([]models.TableForStaffResponse, 0, len(paginatedTableIDs))
	for _, tableID := range paginatedTableIDs {
		var table *models.Table
		for _, t := range allTables {
			if t.ID == tableID {
				table = t
				break
			}
		}
		if table == nil {
			continue
		}

		// Get orders for this table from pre-loaded data
		orders := ordersByTableID[tableID]
		if len(orders) == 0 {
			continue
		}

		// Build order summaries with items using pre-loaded data
		orderSummaries := make([]models.TableOrderSummary, 0)
		totalBill := 0.0
		var customerName string
		isHelpNeeded := false
		isReadyToBill := false

		for _, order := range orders {
			// Get order items from pre-loaded data
			orderItems := itemsByOrderID[order.ID]

			if order.CustomerName != nil {
				customerName = *order.CustomerName
			}

			// Track help needed and ready to bill flags at table level
			if order.IsHelpNeeded {
				isHelpNeeded = true
			}
			if order.IsReadyToBill {
				isReadyToBill = true
			}

			// Build item summaries
			itemSummaries := make([]models.OrderItemSummary, 0)
			for _, item := range orderItems {
				itemSummaries = append(itemSummaries, models.OrderItemSummary{
					ID:        item.ID,
					ItemName:  item.ItemName,
					Quantity:  item.Quantity,
					UnitPrice: item.UnitPrice,
					Status:    item.Status,
				})
			}

			orderSummaries = append(orderSummaries, models.TableOrderSummary{
				ID:            order.ID,
				OrderNumber:   order.OrderNumber,
				Status:        order.Status,
				TotalAmount:   order.Total,
				IsReadyToBill: order.IsReadyToBill,
				IsHelpNeeded:  order.IsHelpNeeded,
				ItemsCount:    len(orderItems),
				Items:         itemSummaries,
				CreatedAt:     order.CreatedAt,
				CustomerName:  customerName,
			})

			totalBill += order.Total
		}

		tableResponse := models.TableForStaffResponse{
			ID:                int(table.ID),
			TableNumber:       table.TableNumber,
			Capacity:          table.Capacity,
			Location:          table.Location,
			Status:            table.Status,
			Orders:            orderSummaries,
			ActiveOrdersCount: len(orders),
			TotalBill:         totalBill,
			IsHelpNeeded:      isHelpNeeded,
			IsReadyToBill:     isReadyToBill,
			CreatedAt:         table.CreatedAt,
			UpdatedAt:         table.UpdatedAt,
		}

		items = append(items, tableResponse)
	}

	return &models.BaseListResponse{
		Total:    int(totalCount),
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

func (s *Service) getTableOrderData(ctx context.Context, tableID int) (*models.TableOrderData, error) {
	var result struct {
		ActiveOrders int     `gorm:"column:active_orders"`
		TotalBill    float64 `gorm:"column:total_bill"`
	}

	err := s.tableRepo.GetDB().Raw(`
		SELECT 
			COUNT(*) as active_orders,
			COALESCE(SUM(total), 0) as total_bill
		FROM orders
		WHERE table_id = ? AND status IN ('pending', 'processing')
	`, tableID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	if result.ActiveOrders == 0 {
		return nil, nil
	}

	return &models.TableOrderData{
		ActiveOrders: result.ActiveOrders,
		TotalBill:    result.TotalBill,
	}, nil
}

// GetTableByID retrieves a single table by ID
func (s *Service) GetTableByID(ctx context.Context, id int) (*models.TableWithOrderData, error) {
	table, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &models.TableWithOrderData{
		ID:               table.ID,
		RestaurantId:     table.RestaurantId,
		TableNumber:      table.TableNumber,
		Capacity:         table.Capacity,
		Location:         table.Location,
		Status:           table.Status,
		QrToken:          table.QrToken,
		QrTokenCreatedAt: table.QrTokenCreatedAt,
		QrTokenExpiresAt: table.QrTokenExpiresAt,
	}

	// If table is occupied, get order data
	if table.Status == "occupied" {
		orderData, err := s.getTableOrderData(ctx, table.ID)
		if err == nil && orderData != nil {
			response.OrderData = orderData
		}
	}

	return response, nil
}

// CreateTable creates a new table
func (s *Service) CreateTable(ctx context.Context, request *models.CreateTableRequest) (*models.Table, error) {
	table := &models.Table{
		TableNumber: request.TableNumber,
		Capacity:    request.Capacity,
		Location:    request.Location,
		Status:      request.Status,
	}

	created, err := s.tableRepo.Create(ctx, table)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) UpdateTable(ctx context.Context, id int, request *models.UpdateTableRequest) (*models.Table, error) {
	existing, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build update columns
	columns := make(map[string]interface{})
	if request.TableNumber != nil {
		columns["table_number"] = *request.TableNumber
	}
	if request.Capacity != nil {
		columns["capacity"] = *request.Capacity
	}
	if request.Location != nil {
		columns["location"] = *request.Location
	}
	if request.Status != nil {
		columns["status"] = *request.Status
	}

	if len(columns) == 0 {
		return existing, nil
	}

	updated, err := s.tableRepo.UpdateColumns(ctx, id, columns)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("table number already exists")
		}
		return nil, err
	}

	return updated, nil
}

func (s *Service) UpdateTableStatus(ctx context.Context, id int, request *models.UpdateTableStatusRequest) (*models.Table, error) {
	_, err := s.tableRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	columns := map[string]interface{}{
		"status": request.Status,
	}

	updated, err := s.tableRepo.UpdateColumns(ctx, id, columns)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) GenerateQrCodeByTableId(ctx context.Context, tableId int) (*models.QrCodeData, error) {
	table, err := s.tableRepo.GetByID(ctx, tableId)
	if err != nil {
		return nil, err
	}

	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	table.QrToken = token
	table.QrTokenCreatedAt = &now
	table.QrTokenExpiresAt = &expiredAt

	_, err = s.tableRepo.Update(ctx, table.ID, table)
	if err != nil {
		return nil, err
	}

	return &models.QrCodeData{
		TableID:   table.ID,
		Token:     token,
		CreatedAt: &now,
		ExpiresAt: &expiredAt,
	}, nil
}

func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateSecureToken - Public function to generate secure tokens (used by payment service)
func GenerateSecureToken(n int) (string, error) {
	return generateSecureToken(n)
}

func (s *Service) GetAllTables(ctx context.Context) ([]*models.Table, error) {
	return s.tableRepo.GetAll(ctx, models.QueryParams{})
}

func (s *Service) GetQrCodeByTableID(ctx context.Context, tableID int) (*models.QrCodeInfo, error) {
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, err
	}

	return &models.QrCodeInfo{
		Token:     table.QrToken,
		CreatedAt: table.QrTokenCreatedAt,
		ExpiresAt: table.QrTokenExpiresAt,
	}, nil
}

// GetTableDetailForStaff - Get table details with all order items for staff view
func (s *Service) GetTableDetailForStaff(ctx context.Context, tableID int) (*models.TableDetailForStaffResponse, error) {
	// Get table
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, err
	}

	// Get all orders for this table (including completed/cancelled)
	allOrders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("table_id = ?", tableID)
	})
	if err != nil {
		return nil, err
	}

	// Extract order IDs for batch query
	orderIDs := make([]int, 0, len(allOrders))
	for _, order := range allOrders {
		orderIDs = append(orderIDs, order.ID)
	}

	// Get all order items for all orders in ONE query
	var orderItems []*models.OrderItem
	if len(orderIDs) > 0 {
		orderItems, err = s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id IN ?", orderIDs)
		})
		if err != nil {
			return nil, err
		}
	}

	// Build response with order items
	itemResponses := make([]models.OrderItemDetailForStaff, 0)
	totalBill := 0.0

	for _, order := range allOrders {
		// Get items for this order from the loaded data
		for _, item := range orderItems {
			if item.OrderID == order.ID {
				itemResponses = append(itemResponses, models.OrderItemDetailForStaff{
					ID:        item.ID,
					OrderID:   item.OrderID,
					ItemName:  item.ItemName,
					Quantity:  item.Quantity,
					UnitPrice: item.UnitPrice,
					Status:    item.Status,
				})
			}
		}
		totalBill += order.Total
	}

	return &models.TableDetailForStaffResponse{
		ID:             table.ID,
		TableNumber:    table.TableNumber,
		Capacity:       table.Capacity,
		Location:       table.Location,
		Status:         table.Status,
		GuestCount:     table.Capacity,
		TotalBill:      totalBill,
		OrderItems:     itemResponses,
		AllOrdersCount: len(allOrders),
		CreatedAt:      table.CreatedAt,
	}, nil
}
