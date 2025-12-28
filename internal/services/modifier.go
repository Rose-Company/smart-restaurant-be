package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func (s *Service) GetModifierGroup(ctx context.Context, request *models.ListModifierGroupRequest) (*models.BaseListResponse, error) {
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

	totalCount, err := s.modifierGroupRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []*models.ModifierGroup{},
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
		default:
			queryParams.QuerySort.Origin = "display_order.asc"
		}
	} else {
		queryParams.QuerySort.Origin = "display_order.asc"
	}

	modifierGroups, err := s.modifierGroupRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	modifierGroupIDs := make([]int, 0, len(modifierGroups))
	for _, modifierGroup := range modifierGroups {
		modifierGroupIDs = append(modifierGroupIDs, modifierGroup.ID)
	}

	modifierOptionFilters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("group_id IN ?", modifierGroupIDs)
		},
	}

	modifierOptions, err := s.modifierOptionRepo.ListByConditions(
		ctx,
		modifierOptionFilters...,
	)
	if err != nil {
		return nil, err
	}

	optionsMap := make(map[int][]*models.ModifierOption)
	for _, opt := range modifierOptions {
		optionsMap[opt.GroupID] = append(optionsMap[opt.GroupID], opt)
	}

	items := make([]*models.ModifierGroupResponse, 0, len(modifierGroups))

	for _, g := range modifierGroups {
		opts := optionsMap[g.ID]
		if opts == nil {
			opts = []*models.ModifierOption{}
		}

		items = append(items, &models.ModifierGroupResponse{
			ID:            g.ID,
			RestaurantID:  g.RestaurantID,
			Name:          g.Name,
			SelectionType: g.SelectionType,
			IsRequired:    g.IsRequired,
			MinSelections: g.MinSelections,
			MaxSelections: g.MaxSelections,
			DisplayOrder:  g.DisplayOrder,
			Status:        g.Status,
			Options:       opts,
			CreatedAt:     g.CreatedAt,
			UpdatedAt:     g.UpdatedAt,
		})
	}

	return &models.BaseListResponse{
		Total:    int(totalCount),
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

func (s *Service) CreateModifierGroup(ctx context.Context, request *models.CreateModifierGroupRequest) (*models.ModifierGroup, error) {
	modifierGroup := &models.ModifierGroup{
		Name:          request.Name,
		SelectionType: request.SelectionType,
		IsRequired:    request.IsRequired,
		MinSelections: request.MinSelections,
		MaxSelections: request.MaxSelections,
		DisplayOrder:  request.DisplayOrder,
		Status:        request.Status,
	}

	created, err := s.modifierGroupRepo.Create(ctx, modifierGroup)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) UpdateModifierGroup(ctx context.Context, id int, request *models.UpdateModifierGroupRequest) (*models.ModifierGroup, error) {
	existing, err := s.modifierGroupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	columns := make(map[string]interface{})
	if request.Name != nil {
		columns["name"] = *request.Name
	}
	if request.SelectionType != nil {
		columns["selection_type"] = *request.SelectionType
	}
	if request.IsRequired != nil {
		columns["is_required"] = *request.IsRequired
	}
	if request.MinSelections != nil {
		columns["min_selections"] = *request.MinSelections
	}
	if request.MaxSelections != nil {
		columns["max_selections"] = *request.MaxSelections
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

	updated, err := s.modifierGroupRepo.UpdateColumns(ctx, id, columns)

	return updated, nil
}

func (s *Service) DeleteModifierGroup(ctx context.Context, id int) error {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id = ?", id)
		},
	}

	_, err := s.modifierGroupRepo.GetDetailByConditions(ctx, filters...)
	if err != nil {
		return err
	}

	columns := map[string]interface{}{
		"status": "inactive",
	}

	_, err = s.modifierGroupRepo.UpdateColumns(ctx, id, columns)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateModifierOptions(ctx context.Context, groupId int, request *models.CreateModifierOptionRequest) (*models.ModifierOption, error) {
	modifierOption := &models.ModifierOption{
		GroupID:         groupId,
		Name:            request.Name,
		PriceAdjustment: request.PriceAdjustment,
		Status:          request.Status,
	}

	created, err := s.modifierOptionRepo.Create(ctx, modifierOption)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) UpdateModifierOptions(ctx context.Context, id int, request *models.UpdateModifierOptionRequest) (*models.ModifierOption, error) {
	existing, err := s.modifierOptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	columns := make(map[string]interface{})
	if request.Name != nil {
		columns["name"] = *request.Name
	}
	if request.PriceAdjustment != nil {
		columns["price_adjustment"] = *request.PriceAdjustment
	}
	if request.Status != nil {
		columns["status"] = *request.Status
	}

	if len(columns) == 0 {
		return existing, nil
	}

	updated, err := s.modifierOptionRepo.UpdateColumns(ctx, id, columns)

	return updated, nil
}

func (s *Service) DeleteModifierOptions(ctx context.Context, id int) error {
	rows, err := s.modifierOptionRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("modifier group not found")
	}

	return nil
}
