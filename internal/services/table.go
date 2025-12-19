package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"strings"

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

	// Build response items with order data
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
		if table.Status == "Occupied" {
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

// Helper function to get order data for a table
func (s *Service) getTableOrderData(ctx context.Context, tableID int) (*models.TableOrderData, error) {
	var result struct {
		ActiveOrders int     `gorm:"column:active_orders"`
		TotalBill    float64 `gorm:"column:total_bill"`
	}

	err := s.tableRepo.GetDB().Raw(`
		SELECT 
			COUNT(*) as active_orders,
			COALESCE(SUM(total_amount), 0) as total_bill
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
