package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"strings"

	"gorm.io/gorm"
)

func (s *Service) GetMenuCategories(ctx context.Context, request *models.ListMenuCategoryRequest) (*models.BaseListResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(request.Page, request.PageSize)

	filters := []repositories.Clause{}

	if request.Status != nil && *request.Status != "" && *request.Status != "all" {
		status := *request.Status
		filters = append(filters, func(tx *gorm.DB) {
			if status == "active" {
				tx.Where("status = ?", "active")
			} else if status == "inactive" {
				tx.Where("status = ?", "inactive")
			}
		})
	}

	if request.Search != nil && *request.Search != "" {
		search := "%" + strings.ToLower(*request.Search) + "%"
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("LOWER(name) LIKE ?", search)
		})
	}

	totalCount, err := s.menuCategoryRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []*models.MenuCategoryResponse{},
		}, nil
	}

	queryParams := models.QueryParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	if request.Sort != "" {
		switch request.Sort {
		case "display_order":
			queryParams.QuerySort.Origin = "display_order.asc"
		case "name":
			queryParams.QuerySort.Origin = "name.asc"
		case "item_count":
			queryParams.QuerySort.Origin = "display_order.asc"
		default:
			queryParams.QuerySort.Origin = "display_order.asc"
		}
	} else {
		queryParams.QuerySort.Origin = "display_order.asc"
	}

	categories, err := s.menuCategoryRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	items := make([]*models.MenuCategoryResponse, 0, len(categories))
	for _, category := range categories {
		itemCount, err := s.getCategoryItemCount(ctx, category.ID)
		if err != nil {
			itemCount = 0
		}

		item := &models.MenuCategoryResponse{
			ID:           category.ID,
			Name:         category.Name,
			Description:  category.Description,
			ItemCount:    itemCount,
			IsActive:     category.Status == "active",
			DisplayOrder: category.DisplayOrder,
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

func (s *Service) getCategoryItemCount(ctx context.Context, categoryID int) (int, error) {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("category_id = ? AND is_deleted = FALSE", categoryID)
		},
	}

	count, err := s.menuItemRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *Service) GetMenuCategoryByID(ctx context.Context, id int) (*models.MenuCategoryDetailResponse, error) {
	category, err := s.menuCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	itemCount, err := s.getCategoryItemCount(ctx, category.ID)
	if err != nil {
		itemCount = 0
	}

	response := &models.MenuCategoryDetailResponse{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		ItemCount:    itemCount,
		IsActive:     category.Status == "active",
		DisplayOrder: category.DisplayOrder,
	}

	return response, nil
}

func (s *Service) CreateMenuCategory(ctx context.Context, request *models.CreateMenuCategoryRequest) (*models.MenuCategory, error) {
	restaurantID := 1
	if request.RestaurantID != nil {
		restaurantID = *request.RestaurantID
	}

	displayOrder := 0
	if request.DisplayOrder != nil {
		displayOrder = *request.DisplayOrder
	} else {
		filters := []repositories.Clause{
			func(tx *gorm.DB) {
				tx.Where("restaurant_id = ?", restaurantID)
			},
		}
		categories, err := s.menuCategoryRepo.List(ctx, models.QueryParams{
			Limit: 1,
			QuerySort: models.QuerySort{
				Origin: "display_order.desc",
			},
		}, filters...)

		if err == nil && len(categories) > 0 {
			displayOrder = categories[0].DisplayOrder + 1
		}
	}

	category := &models.MenuCategory{
		RestaurantID: restaurantID,
		Name:         request.Name,
		Description:  request.Description,
		DisplayOrder: displayOrder,
		Status:       request.Status,
	}

	created, err := s.menuCategoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) UpdateMenuCategory(ctx context.Context, id int, request *models.UpdateMenuCategoryRequest) (*models.MenuCategory, error) {
	existing, err := s.menuCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	columns := make(map[string]interface{})
	if request.Name != nil {
		columns["name"] = *request.Name
	}
	if request.Description != nil {
		columns["description"] = *request.Description
	}
	if request.DisplayOrder != nil {
		columns["display_order"] = *request.DisplayOrder
	}
	if request.Status != nil {
		columns["status"] = *request.Status
	}

	if len(columns) == 0 {
		return existing, nil
	}

	updated, err := s.menuCategoryRepo.UpdateColumns(ctx, id, columns)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) UpdateMenuCategoryStatus(ctx context.Context, id int, request *models.UpdateMenuCategoryStatusRequest) (*models.MenuCategory, error) {
	_, err := s.menuCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	status := "inactive"
	if request.IsActive {
		status = "active"
	}

	columns := map[string]interface{}{
		"status": status,
	}

	updated, err := s.menuCategoryRepo.UpdateColumns(ctx, id, columns)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
