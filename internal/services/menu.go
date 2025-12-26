package services

import (
	"app-noti/common"
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"fmt"
	"strconv"
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

	categoryIDs := make([]int, 0, len(categories))
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID)
	}

	itemCountMap, err := s.getCategoryItemCounts(ctx, categoryIDs)
	if err != nil {
		itemCountMap = make(map[int]int)
	}

	items := make([]*models.MenuCategoryResponse, 0, len(categories))
	for _, category := range categories {
		itemCount := itemCountMap[category.ID]

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

func (s *Service) getCategoryItemCounts(ctx context.Context, categoryIDs []int) (map[int]int, error) {
	if len(categoryIDs) == 0 {
		return make(map[int]int), nil
	}

	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("category_id IN ? AND is_deleted = FALSE", categoryIDs)
		},
	}

	countMap, err := s.menuItemRepo.CountGroupByInt(ctx, "category_id", filters...)
	if err != nil {
		return nil, err
	}

	return countMap, nil
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

func (s *Service) GetMenuItems(ctx context.Context, request *models.ListMenuItemRequest) (*models.BaseListResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(request.Page, request.PageSize)

	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("menu_items.is_deleted = FALSE")
		},
	}

	if request.Status != nil && *request.Status != "" && *request.Status != "all" {
		status := *request.Status
		filters = append(filters, func(tx *gorm.DB) {
			statusMap := map[string]string{
				"available":   "available",
				"unavailable": "unavailable",
				"sold_out":    "sold_out",
			}
			if dbStatus, ok := statusMap[strings.ToLower(status)]; ok {
				tx.Where("menu_items.status = ?", dbStatus)
			}
		})
	}

	if request.Category != nil && *request.Category != "" && *request.Category != "all" {
		categoryName := *request.Category
		filters = append(filters, func(tx *gorm.DB) {
			tx.Joins("LEFT JOIN menu_categories ON menu_items.category_id = menu_categories.id").
				Where("LOWER(menu_categories.name) = ?", strings.ToLower(categoryName))
		})
	}

	if request.Search != nil && *request.Search != "" {
		search := "%" + strings.ToLower(*request.Search) + "%"
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("LOWER(menu_items.name) LIKE ?", search)
		})
	}

	totalCount, err := s.menuItemRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return &models.BaseListResponse{
			Total:    0,
			Page:     page,
			PageSize: pageSize,
			Items:    []*models.MenuItemResponse{},
		}, nil
	}

	queryParams := models.QueryParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	sortMap := map[string]string{
		"default":     "id.asc",
		"name":        "name.asc",
		"price_asc":   "price.asc",
		"price_desc":  "price.desc",
		"last_update": "updated_at.desc",
	}

	if sortOrder, ok := sortMap[request.Sort]; ok {
		queryParams.QuerySort.Origin = sortOrder
	} else {
		queryParams.QuerySort.Origin = sortMap["default"]
	}

	menuItems, err := s.menuItemRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	categoryIDs := make([]int, 0, len(menuItems))
	itemIDs := make([]int, 0, len(menuItems))
	for _, item := range menuItems {
		categoryIDs = append(categoryIDs, item.CategoryID)
		itemIDs = append(itemIDs, item.ID)
	}

	categoryMap, err := s.getCategoryMapByIDs(ctx, categoryIDs)
	if err != nil {
		categoryMap = make(map[int]string)
	}

	primaryImageMap, err := s.getPrimaryImageMap(ctx, itemIDs)
	if err != nil {
		primaryImageMap = make(map[int]string)
	}

	items := make([]*models.MenuItemResponse, 0, len(menuItems))
	for _, menuItem := range menuItems {
		categoryName := categoryMap[menuItem.CategoryID]

		statusMap := map[string]string{
			"available":   "Available",
			"unavailable": "Unavailable",
			"sold_out":    "Sold Out",
		}
		displayStatus := statusMap[menuItem.Status]
		if displayStatus == "" {
			displayStatus = menuItem.Status
		}

		lastUpdate := ""
		if menuItem.UpdatedAt != nil {
			lastUpdate = menuItem.UpdatedAt.Format("2006-01-02")
		}

		item := &models.MenuItemResponse{
			ID:              menuItem.ID,
			Name:            menuItem.Name,
			Category:        categoryName,
			Price:           menuItem.Price,
			Status:          displayStatus,
			LastUpdate:      lastUpdate,
			ChefRecommended: menuItem.IsChefRecommended,
			ImageURL:        primaryImageMap[menuItem.ID],
			Description:     menuItem.Description,
			PreparationTime: menuItem.PrepTimeMinutes,
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

func (s *Service) getCategoryMapByIDs(ctx context.Context, categoryIDs []int) (map[int]string, error) {
	if len(categoryIDs) == 0 {
		return make(map[int]string), nil
	}

	uniqueIDs := make(map[int]bool)
	for _, id := range categoryIDs {
		uniqueIDs[id] = true
	}

	ids := make([]int, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}

	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id IN ?", ids)
		},
	}

	categories, err := s.menuCategoryRepo.List(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	return categoryMap, nil
}

func (s *Service) getPrimaryImageMap(ctx context.Context, itemIDs []int) (map[int]string, error) {
	if len(itemIDs) == 0 {
		return make(map[int]string), nil
	}

	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("menu_item_id IN ? AND is_primary = TRUE", itemIDs)
		},
	}

	photos, err := s.menuItemPhotoRepo.List(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	imageMap := make(map[int]string)
	for _, photo := range photos {
		imageMap[photo.MenuItemID] = photo.Url
	}

	return imageMap, nil
}

func (s *Service) getMenuItemModifiers(ctx context.Context, menuItemID int) ([]models.MenuItemModifier, error) {
	associationFilters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("menu_item_id = ?", menuItemID)
		},
	}

	associations, err := s.menuItemModifierGroupRepo.List(ctx, models.QueryParams{}, associationFilters...)
	if err != nil || len(associations) == 0 {
		return []models.MenuItemModifier{}, nil
	}

	groupIDs := make([]int, 0, len(associations))
	for _, assoc := range associations {
		groupIDs = append(groupIDs, assoc.GroupID)
	}

	groupFilters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id IN ?", groupIDs)
		},
	}

	groups, err := s.modifierGroupRepo.List(ctx, models.QueryParams{}, groupFilters...)
	if err != nil {
		return []models.MenuItemModifier{}, nil
	}

	optionFilters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("group_id IN ? AND status = ?", groupIDs, "active")
		},
	}

	options, err := s.modifierOptionRepo.List(ctx, models.QueryParams{}, optionFilters...)
	if err != nil {
		options = []*models.ModifierOption{}
	}

	optionsPreviewMap := make(map[int]string)
	optionsByGroup := make(map[int][]*models.ModifierOption)
	for _, option := range options {
		optionsByGroup[option.GroupID] = append(optionsByGroup[option.GroupID], option)
	}

	for groupID, opts := range optionsByGroup {
		previewParts := make([]string, 0)
		for i, opt := range opts {
			if i >= 3 {
				previewParts = append(previewParts, "...")
				break
			}
			priceText := ""
			if opt.PriceAdjustment > 0 {
				priceText = fmt.Sprintf(" (+$%.2f)", opt.PriceAdjustment)
			}
			previewParts = append(previewParts, opt.Name+priceText)
		}
		optionsPreviewMap[groupID] = common.JoinStrings(previewParts, ", ")
	}

	modifiers := make([]models.MenuItemModifier, 0, len(groups))
	for _, group := range groups {
		selectionTypeDisplay := "Single"
		if group.SelectionType == "multiple" {
			selectionTypeDisplay = "Multi"
		}

		modifier := models.MenuItemModifier{
			ID:              fmt.Sprintf("%d", group.ID),
			ModifierGroupID: fmt.Sprintf("%d", group.ID),
			Name:            group.Name,
			Required:        group.IsRequired,
			SelectionType:   selectionTypeDisplay,
			OptionsPreview:  optionsPreviewMap[group.ID],
		}
		modifiers = append(modifiers, modifier)
	}

	return modifiers, nil
}

func (s *Service) GetMenuItemByID(ctx context.Context, id int) (*models.MenuItemDetailResponse, error) {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id = ? AND is_deleted = FALSE", id)
		},
	}

	menuItem, err := s.menuItemRepo.GetDetailByConditions(ctx, filters...)
	if err != nil {
		return nil, err
	}

	category, err := s.menuCategoryRepo.GetByID(ctx, menuItem.CategoryID)
	categoryName := ""
	if err == nil {
		categoryName = category.Name
	}

	statusMap := map[string]string{
		"available":   "Available",
		"unavailable": "Unavailable",
		"sold_out":    "Sold Out",
	}
	displayStatus := statusMap[menuItem.Status]

	lastUpdate := ""
	if menuItem.UpdatedAt != nil {
		lastUpdate = menuItem.UpdatedAt.Format("2006-01-02")
	}

	photoFilters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("menu_item_id = ?", id)
		},
	}

	photos, err := s.menuItemPhotoRepo.List(ctx, models.QueryParams{}, photoFilters...)
	if err != nil {
		photos = []*models.MenuItemPhoto{}
	}

	imageRequests := make([]models.MenuItemPhotoRequest, 0, len(photos))
	primaryImageURL := ""
	for _, photo := range photos {
		imageRequests = append(imageRequests, models.MenuItemPhotoRequest{
			ID:        string(rune(photo.ID)),
			URL:       photo.Url,
			IsPrimary: photo.IsPrimary,
		})
		if photo.IsPrimary {
			primaryImageURL = photo.Url
		}
	}

	// Fetch modifiers with group details
	modifiers, err := s.getMenuItemModifiers(ctx, id)
	if err != nil {
		modifiers = []models.MenuItemModifier{}
	}

	response := &models.MenuItemDetailResponse{
		ID:              menuItem.ID,
		Name:            menuItem.Name,
		Category:        categoryName,
		Price:           menuItem.Price,
		Status:          displayStatus,
		LastUpdate:      lastUpdate,
		ChefRecommended: menuItem.IsChefRecommended,
		ImageURL:        primaryImageURL,
		Description:     menuItem.Description,
		PreparationTime: menuItem.PrepTimeMinutes,
		Images:          imageRequests,
		Modifiers:       modifiers,
	}

	return response, nil
}

func (s *Service) CreateMenuItem(ctx context.Context, request *models.CreateMenuItemRequest) (*models.MenuItem, error) {
	category, err := s.menuCategoryRepo.GetByID(ctx, request.CategoryID)
	if err != nil {
		return nil, err
	}

	menuItem := &models.MenuItem{
		RestaurantID:      category.RestaurantID,
		CategoryID:        request.CategoryID,
		Name:              request.Name,
		Description:       request.Description,
		Price:             request.Price,
		PrepTimeMinutes:   request.PrepTimeMinutes,
		Status:            request.Status,
		IsChefRecommended: request.IsChefRecommended,
		IsDeleted:         false,
	}

	created, err := s.menuItemRepo.Create(ctx, menuItem)
	if err != nil {
		return nil, err
	}

	// Create menu item photos
	if len(request.Images) > 0 {
		photos := make([]*models.MenuItemPhoto, 0, len(request.Images))
		for _, img := range request.Images {
			photo := &models.MenuItemPhoto{
				MenuItemID: created.ID,
				Url:        img.URL,
				IsPrimary:  img.IsPrimary,
			}
			photos = append(photos, photo)
		}
		if err := s.menuItemPhotoRepo.CreatesMultiple(ctx, photos); err != nil {
			// Log error but don't fail the whole operation
		}
	}

	// Create menu item modifier groups
	if len(request.Modifiers) > 0 {
		modifiers := make([]*models.MenuItemModifierGroup, 0, len(request.Modifiers))
		for _, mod := range request.Modifiers {
			modifierGroupID, _ := strconv.Atoi(mod.ModifierGroupID)
			modifier := &models.MenuItemModifierGroup{
				MenuItemID: created.ID,
				GroupID:    modifierGroupID,
			}
			modifiers = append(modifiers, modifier)
		}
		if err := s.menuItemModifierGroupRepo.CreatesMultiple(ctx, modifiers); err != nil {
			// Log error but don't fail the whole operation
		}
	}

	return created, nil
}

func (s *Service) UpdateMenuItem(ctx context.Context, id int, request *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id = ? AND is_deleted = FALSE", id)
		},
	}

	existing, err := s.menuItemRepo.GetDetailByConditions(ctx, filters...)
	if err != nil {
		return nil, err
	}

	columns := make(map[string]interface{})

	if request.CategoryID != nil {
		_, err := s.menuCategoryRepo.GetByID(ctx, *request.CategoryID)
		if err != nil {
			return nil, err
		}
		columns["category_id"] = *request.CategoryID
	}

	if request.Name != nil {
		columns["name"] = *request.Name
	}
	if request.Description != nil {
		columns["description"] = *request.Description
	}
	if request.Price != nil {
		columns["price"] = *request.Price
	}
	if request.PrepTimeMinutes != nil {
		columns["prep_time_minutes"] = *request.PrepTimeMinutes
	}
	if request.Status != nil {
		columns["status"] = *request.Status
	}
	if request.IsChefRecommended != nil {
		columns["is_chef_recommended"] = *request.IsChefRecommended
	}

	if len(columns) == 0 && len(request.Images) == 0 && len(request.Modifiers) == 0 {
		return existing, nil
	}

	var updated *models.MenuItem
	if len(columns) > 0 {
		updated, err = s.menuItemRepo.UpdateColumns(ctx, id, columns)
		if err != nil {
			return nil, err
		}
	} else {
		updated = existing
	}

	// Update images if provided
	if len(request.Images) > 0 {
		// Delete existing photos
		deleteFilters := []repositories.Clause{
			func(tx *gorm.DB) {
				tx.Where("menu_item_id = ?", id)
			},
		}
		s.menuItemPhotoRepo.Delete(ctx, deleteFilters...)

		// Create new photos
		photos := make([]*models.MenuItemPhoto, 0, len(request.Images))
		for _, img := range request.Images {
			photo := &models.MenuItemPhoto{
				MenuItemID: id,
				Url:        img.URL,
				IsPrimary:  img.IsPrimary,
			}
			photos = append(photos, photo)
		}
		s.menuItemPhotoRepo.CreatesMultiple(ctx, photos)
	}

	// Update modifiers if provided
	if len(request.Modifiers) > 0 {
		// Delete existing modifiers
		deleteFilters := []repositories.Clause{
			func(tx *gorm.DB) {
				tx.Where("menu_item_id = ?", id)
			},
		}
		s.menuItemModifierGroupRepo.Delete(ctx, deleteFilters...)

		// Create new modifiers
		modifiers := make([]*models.MenuItemModifierGroup, 0, len(request.Modifiers))
		for _, mod := range request.Modifiers {
			modifierGroupID, _ := strconv.Atoi(mod.ModifierGroupID)
			modifier := &models.MenuItemModifierGroup{
				MenuItemID: id,
				GroupID:    modifierGroupID,
			}
			modifiers = append(modifiers, modifier)
		}
		s.menuItemModifierGroupRepo.CreatesMultiple(ctx, modifiers)
	}

	return updated, nil
}

func (s *Service) DeleteMenuItem(ctx context.Context, id int) error {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id = ? AND is_deleted = FALSE", id)
		},
	}

	_, err := s.menuItemRepo.GetDetailByConditions(ctx, filters...)
	if err != nil {
		return err
	}

	columns := map[string]interface{}{
		"is_deleted": true,
	}

	_, err = s.menuItemRepo.UpdateColumns(ctx, id, columns)
	return err
}
