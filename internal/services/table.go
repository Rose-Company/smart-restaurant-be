package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

func (s *Service) GenerateQrCodeByTableId(ctx context.Context, tableId int) (string, error) {
	table, err := s.tableRepo.GetByID(ctx, tableId)
	if err != nil {
		return "", err
	}

	token, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	table.QrToken = token
	table.QrTokenCreatedAt = &now
	table.QrTokenExpiresAt = &expiredAt

	_, err = s.tableRepo.Update(ctx, table.ID, table)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://smart-restaurant-fe.vercel.app/menu?table=%d&token=%s",
		table.ID,
		token,
	)

	return url, nil
}

func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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
