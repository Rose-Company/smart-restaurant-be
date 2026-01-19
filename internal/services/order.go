package services

import (
	"app-noti/common"
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"app-noti/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateOrder - TASK-001: Create/Submit order from cart
func (s *Service) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.OrderResponse, error) {
	// Start transaction
	return s.createOrderWithTransaction(ctx, req)
}

func (s *Service) createOrderWithTransaction(ctx context.Context, req models.CreateOrderRequest) (*models.OrderResponse, error) {
	var order *models.Order
	var orderItems []*models.OrderItem
	var totalAmount float64
	var taxAmount float64

	// Execute in transaction
	err := s.tableRepo.ExecRaw(ctx, "BEGIN")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.tableRepo.ExecRaw(ctx, "ROLLBACK")
		}
	}()

	// 1. Validate table exists
	table, err := s.tableRepo.GetByID(ctx, req.TableID)
	if err != nil {
		return nil, common.ErrTableNotFound
	}

	// 2. Generate order number
	orderNumber := common.GenerateOrderNumber()

	// 3. Get customer info if provided
	var customerName string
	var customerPhone *string
	var customerEmail *string

	if req.CustomerID != nil {
		user, err := s.userRepo.GetByID(ctx, *req.CustomerID)
		if err == nil && user != nil {
			customerName = fmt.Sprintf("%s %s", *user.FirstName, *user.LastName)
			customerPhone = &user.PhoneNumber
			customerEmail = &user.Email
		}
	}

	if customerName == "" {
		customerName = "Guest"
	}

	// 4. Calculate totals and create order items
	for _, itemReq := range req.Items {
		// Get menu item
		menuItem, err := s.menuItemRepo.GetByID(ctx, itemReq.MenuItemID)
		if err != nil {
			return nil, common.ErrMenuItemNotFound
		}

		// Calculate item subtotal
		itemSubtotal := menuItem.Price * float64(itemReq.Quantity)
		modifiersTotal := 0.0

		// Create order item
		orderItem := &models.OrderItem{
			MenuItemID:          &menuItem.ID,
			ItemName:            menuItem.Name,
			ItemDescription:     menuItem.Description,
			Quantity:            itemReq.Quantity,
			UnitPrice:           menuItem.Price,
			Subtotal:            itemSubtotal,
			Status:              "pending",
			SpecialInstructions: itemReq.SpecialInstructions,
		}

		// Handle modifiers
		if len(itemReq.Modifiers) > 0 {
			for _, modReq := range itemReq.Modifiers {
				modOption, err := s.modifierOptionRepo.GetByID(ctx, modReq.ModifierOptionID)
				if err != nil {
					continue
				}

				// Get modifier group for validation (not used in response but kept for future use)
				_, err = s.modifierGroupRepo.GetByID(ctx, modReq.ModifierGroupID)
				if err != nil {
					continue
				}

				modifiersTotal += modOption.PriceAdjustment
			}
		}

		orderItem.ModifiersTotal = modifiersTotal
		totalAmount += itemSubtotal + modifiersTotal
		orderItems = append(orderItems, orderItem)
	}

	// 5. Calculate tax (10%)
	taxAmount = totalAmount * 0.10

	// 6. Create order
	order = &models.Order{
		TableID:        req.TableID,
		OrderNumber:    orderNumber,
		CustomerUserID: req.CustomerID,
		CustomerName:   &customerName,
		CustomerPhone:  customerPhone,
		CustomerEmail:  customerEmail,
		Status:         "pending",
		Subtotal:       totalAmount,
		Tax:            taxAmount,
		Discount:       0,
		Total:          totalAmount + taxAmount,
		Notes:          req.Notes,
		Priority:       "normal",
		Source:         "qr",
	}

	createdOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// 7. Create order items and modifiers
	var responseItems []models.OrderItemResponse
	for i, itemReq := range req.Items {
		orderItems[i].OrderID = createdOrder.ID

		createdItem, err := s.orderItemRepo.Create(ctx, orderItems[i])
		if err != nil {
			return nil, err
		}

		// Create modifiers
		var modifiers []models.OrderModifierResponse
		for _, modReq := range itemReq.Modifiers {
			modOption, _ := s.modifierOptionRepo.GetByID(ctx, modReq.ModifierOptionID)
			modGroup, _ := s.modifierGroupRepo.GetByID(ctx, modReq.ModifierGroupID)

			if modOption != nil && modGroup != nil {
				orderMod := &models.OrderModifier{
					OrderItemID:        createdItem.ID,
					ModifierGroupID:    modGroup.ID,
					ModifierGroupName:  modGroup.Name,
					ModifierOptionID:   modOption.ID,
					ModifierOptionName: modOption.Name,
					PriceAdjustment:    modOption.PriceAdjustment,
				}

				_, err := s.orderModifierRepo.Create(ctx, orderMod)
				if err == nil {
					modifiers = append(modifiers, models.OrderModifierResponse{
						ModifierGroupID:    modGroup.ID,
						ModifierGroupName:  modGroup.Name,
						ModifierOptionID:   modOption.ID,
						ModifierOptionName: modOption.Name,
						Price:              modOption.PriceAdjustment,
					})
				}
			}
		}

		responseItems = append(responseItems, models.OrderItemResponse{
			ID:                  createdItem.ID,
			MenuItemID:          createdItem.MenuItemID,
			MenuItemName:        createdItem.ItemName,
			Quantity:            createdItem.Quantity,
			UnitPrice:           createdItem.UnitPrice,
			Subtotal:            createdItem.Subtotal,
			SpecialInstructions: createdItem.SpecialInstructions,
			Status:              createdItem.Status,
			Modifiers:           modifiers,
		})
	}

	// 8. Create timeline entry
	timeline := &models.OrderTimeline{
		OrderID:       createdOrder.ID,
		Status:        "pending",
		UpdatedBy:     "customer",
		UpdatedByName: &customerName,
		Note:          common.StrPtr("Order placed"),
	}
	s.orderTimelineRepo.Create(ctx, timeline)

	// Commit transaction
	err = s.tableRepo.ExecRaw(ctx, "COMMIT")
	if err != nil {
		return nil, err
	}

	// 9. Build response
	response := &models.OrderResponse{
		ID:             createdOrder.ID,
		OrderNumber:    createdOrder.OrderNumber,
		TableID:        createdOrder.TableID,
		TableName:      table.TableNumber,
		CustomerID:     createdOrder.CustomerUserID,
		CustomerName:   *createdOrder.CustomerName,
		Status:         createdOrder.Status,
		TotalAmount:    createdOrder.Subtotal,
		TaxAmount:      createdOrder.Tax,
		DiscountAmount: createdOrder.Discount,
		FinalAmount:    createdOrder.Total,
		Notes:          createdOrder.Notes,
		CreatedAt:      createdOrder.CreatedAt,
		Items:          responseItems,
	}

	return response, nil
}

// GetOrders - TASK-002: Get orders list with dynamic filtering
func (s *Service) GetOrders(ctx context.Context, request *models.ListOrdersRequest) (*models.PaginatedOrdersResponse, error) {
	page, pageSize := utils.GetPageAndPageSize(request.Page, request.PageSize)

	// Build filters array
	filters := []repositories.Clause{}

	// Filter by status
	if request.Status != nil && *request.Status != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("status = ?", *request.Status)
		})
	}

	// Filter by table_id
	if request.TableID != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("table_id = ?", *request.TableID)
		})
	}

	// Filter by date range
	if request.DateFrom != nil && *request.DateFrom != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("created_at >= ?", *request.DateFrom)
		})
	}
	if request.DateTo != nil && *request.DateTo != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("created_at <= ?", *request.DateTo)
		})
	}

	// Search by order number or customer name
	if request.Search != nil && *request.Search != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("order_number ILIKE ? OR customer_name ILIKE ?", "%"+*request.Search+"%", "%"+*request.Search+"%")
		})
	}

	// Filter by menu category (if category filter is provided)
	// This is used when user wants to see orders with specific menu categories
	var categoryFilterIDs []int
	if request.Category != nil && *request.Category != "" {
		// Get menu items by category using repository
		menuItems, err := s.menuItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Joins("JOIN menu_categories mc ON menu_items.category_id = mc.id").
				Where("mc.name ILIKE ?", "%"+*request.Category+"%")
		})
		if err == nil && len(menuItems) > 0 {
			for _, item := range menuItems {
				categoryFilterIDs = append(categoryFilterIDs, item.ID)
			}
		}
	}

	if request.IsReadyToBill != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("is_ready_to_bill = ?", *request.IsReadyToBill)
		})
	}

	// Filter by is_help_needed
	if request.IsHelpNeeded != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("is_help_needed = ?", *request.IsHelpNeeded)
		})
	}

	// Sorting
	sortOrder := "created_at DESC"
	if request.Sort != "" {
		switch request.Sort {
		case "created_at_asc":
			sortOrder = "created_at ASC"
		case "total_amount_desc":
			sortOrder = "total DESC"
		case "total_amount_asc":
			sortOrder = "total ASC"
		}
	}
	filters = append(filters, func(tx *gorm.DB) {
		tx.Order(sortOrder)
	})

	// Preload relationships
	filters = append(filters, func(tx *gorm.DB) {
		tx.Preload("Table").Preload("CustomerUser").Preload("Waiter")
	})

	// Get total count
	total, err := s.orderRepo.Count(ctx, models.QueryParams{}, filters[:len(filters)-1]...) // Exclude preload from count
	if err != nil {
		return nil, err
	}

	// Get orders with pagination
	queryParams := models.QueryParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
	orders, err := s.orderRepo.List(ctx, queryParams, filters...)
	if err != nil {
		return nil, err
	}

	// Build response
	var items []models.OrderListItemResponse
	for _, order := range orders {
		// Get order items with modifiers preloaded
		orderItems, _ := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id = ?", order.ID)
			tx.Preload("MenuItem")
			tx.Preload("Modifiers")
		})

		// Filter items by category if provided
		var filteredItems []*models.OrderItem
		if len(categoryFilterIDs) > 0 {
			for _, item := range orderItems {
				if item.MenuItemID != nil {
					for _, catID := range categoryFilterIDs {
						if *item.MenuItemID == catID {
							filteredItems = append(filteredItems, item)
							break
						}
					}
				}
			}
		} else {
			filteredItems = orderItems
		}

		// Build order item responses
		var responseItems []models.OrderItemResponse
		for _, item := range filteredItems {
			// Build modifiers from preloaded data
			var modResponses []models.OrderModifierResponse
			if item.Modifiers != nil && len(item.Modifiers) > 0 {
				for _, mod := range item.Modifiers {
					modResponses = append(modResponses, models.OrderModifierResponse{
						ModifierGroupID:    mod.ModifierGroupID,
						ModifierGroupName:  mod.ModifierGroupName,
						ModifierOptionID:   mod.ModifierOptionID,
						ModifierOptionName: mod.ModifierOptionName,
						Price:              mod.PriceAdjustment,
					})
				}
			}

			// Get menu item image (MenuItem doesn't have Photos field)
			var menuItemImage *string

			responseItems = append(responseItems, models.OrderItemResponse{
				ID:                  item.ID,
				MenuItemID:          item.MenuItemID,
				MenuItemName:        item.ItemName,
				MenuItemImage:       menuItemImage,
				Quantity:            item.Quantity,
				UnitPrice:           item.UnitPrice,
				Subtotal:            item.Subtotal,
				SpecialInstructions: item.SpecialInstructions,
				Status:              item.Status,
				Modifiers:           modResponses,
			})
		}

		tableName := "Unknown"
		if order.Table != nil {
			tableName = order.Table.TableNumber
		}

		var waiterName *string
		if order.Waiter != nil {
			if order.Waiter.FirstName != nil && order.Waiter.LastName != nil {
				name := fmt.Sprintf("%s %s", *order.Waiter.FirstName, *order.Waiter.LastName)
				waiterName = &name
			}
		}

		// Get customer name safely
		customerName := ""
		if order.CustomerUser != nil {
			if order.CustomerUser.FirstName != nil && order.CustomerUser.LastName != nil {
				customerName = fmt.Sprintf("%s %s", *order.CustomerUser.FirstName, *order.CustomerUser.LastName)
			}
		}

		items = append(items, models.OrderListItemResponse{
			ID:                 order.ID,
			OrderNumber:        order.OrderNumber,
			TableID:            order.TableID,
			TableName:          tableName,
			CustomerID:         order.CustomerUserID,
			CustomerName:       customerName,
			Status:             order.Status,
			TotalAmount:        order.Total,
			ItemsCount:         len(filteredItems),
			Items:              responseItems,
			CreatedAt:          order.CreatedAt,
			UpdatedAt:          order.UpdatedAt,
			EstimatedReadyTime: order.EstimatedReadyTime,
			WaiterID:           order.WaiterID,
			WaiterName:         waiterName,
			IsReadyToBill:      order.IsReadyToBill,
			IsHelpNeeded:       order.IsHelpNeeded,
		})
	}

	return &models.PaginatedOrdersResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}, nil
}

// GetOrderByID - TASK-003: Get order details (role-based response)
func (s *Service) GetOrderByID(ctx context.Context, orderID int, role string) (*models.OrderResponse, error) {
	// Get order with relationships
	order, err := s.orderRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("id = ?", orderID)
		tx.Preload("Table")
		tx.Preload("CustomerUser")
		tx.Preload("Waiter")
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrOrderNotFound
		}
		return nil, err
	}

	// Get order items
	items, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id = ?", orderID)
		tx.Preload("MenuItem")
	})
	if err != nil {
		return nil, err
	}

	// Build items response
	var responseItems []models.OrderItemResponse
	for _, item := range items {
		// Get modifiers
		modifiers, _ := s.orderModifierRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_item_id = ?", item.ID)
		})

		var modResponses []models.OrderModifierResponse
		for _, mod := range modifiers {
			modResponses = append(modResponses, models.OrderModifierResponse{
				ModifierGroupID:    mod.ModifierGroupID,
				ModifierGroupName:  mod.ModifierGroupName,
				ModifierOptionID:   mod.ModifierOptionID,
				ModifierOptionName: mod.ModifierOptionName,
				Price:              mod.PriceAdjustment,
			})
		}

		// Menu item image - MenuItem doesn't have Photos field in this version
		var menuItemImage *string

		responseItems = append(responseItems, models.OrderItemResponse{
			ID:                  item.ID,
			MenuItemID:          item.MenuItemID,
			MenuItemName:        item.ItemName,
			MenuItemImage:       menuItemImage,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			Subtotal:            item.Subtotal,
			SpecialInstructions: item.SpecialInstructions,
			Status:              item.Status,
			Modifiers:           modResponses,
		})
	}

	// Build response based on role
	customerName := "Guest"
	if order.CustomerName != nil {
		customerName = *order.CustomerName
	}

	response := &models.OrderResponse{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		TableID:        order.TableID,
		CustomerName:   customerName,
		Status:         order.Status,
		TotalAmount:    order.Subtotal,
		TaxAmount:      order.Tax,
		DiscountAmount: order.Discount,
		FinalAmount:    order.Total,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
		Items:          responseItems,
	}

	if order.Table != nil {
		response.TableName = order.Table.TableNumber
	}

	// Kitchen role - simplified response
	if role == "kitchen" {
		response.CustomerID = nil
		response.CustomerPhone = nil
		response.Notes = order.Notes
		response.EstimatedReadyTime = order.EstimatedReadyTime
		return response, nil
	}

	// Full response for customer/waiter/admin
	response.CustomerID = order.CustomerUserID
	response.CustomerPhone = order.CustomerPhone
	response.Notes = order.Notes
	response.SpecialInstructions = order.SpecialInstructions
	response.EstimatedReadyTime = order.EstimatedReadyTime
	response.WaiterID = order.WaiterID

	if order.Waiter != nil {
		if order.Waiter.FirstName != nil && order.Waiter.LastName != nil {
			waiterName := fmt.Sprintf("%s %s", *order.Waiter.FirstName, *order.Waiter.LastName)
			response.WaiterName = &waiterName
		}
	}

	// Get timeline
	timeline, _ := s.orderTimelineRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id = ?", orderID)
		tx.Order("timestamp ASC")
	})

	var timelineResponses []models.OrderTimelineResponse
	for _, t := range timeline {
		timelineResponses = append(timelineResponses, models.OrderTimelineResponse{
			ID:            t.ID,
			Status:        t.Status,
			Timestamp:     t.Timestamp,
			UpdatedBy:     t.UpdatedBy,
			UpdatedByName: t.UpdatedByName,
			Note:          t.Note,
		})
	}
	response.Timeline = timelineResponses

	return response, nil
}

// UpdateOrderStatus - TASK-004: Update order status
func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int, req models.UpdateOrderStatusRequest, userID, userName interface{}) (*models.OrderStatusUpdateResponse, error) {
	// Get current order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	previousStatus := order.Status

	// Validate status transition
	if !common.IsValidOrderStatusTransition(previousStatus, req.Status) {
		return nil, fmt.Errorf("invalid status transition. Cannot change from '%s' to '%s'", previousStatus, req.Status)
	}

	// Update order status
	now := time.Now()
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": now,
	}

	// Update specific timestamps based on status
	switch req.Status {
	case "confirmed":
		updates["accepted_at"] = now
	case "preparing":
		updates["preparing_at"] = now
		updates["estimated_ready_time"] = now.Add(30 * time.Minute)
	case "ready":
		updates["ready_at"] = now
	case "served":
		updates["served_at"] = now
	case "completed":
		updates["completed_at"] = now
	}

	_, err = s.orderRepo.UpdateColumns(ctx, orderID, updates)
	if err != nil {
		return nil, err
	}

	// Create timeline entry
	var updatedByName *string
	if userName != nil {
		name := userName.(string)
		updatedByName = &name
	}

	timeline := &models.OrderTimeline{
		OrderID:       orderID,
		Status:        req.Status,
		UpdatedBy:     req.UpdatedBy,
		UpdatedByName: updatedByName,
		Note:          req.Reason,
	}
	s.orderTimelineRepo.Create(ctx, timeline)

	// Get updated order
	updatedOrder, _ := s.orderRepo.GetByID(ctx, orderID)

	return &models.OrderStatusUpdateResponse{
		ID:                 updatedOrder.ID,
		OrderNumber:        updatedOrder.OrderNumber,
		Status:             updatedOrder.Status,
		PreviousStatus:     previousStatus,
		UpdatedAt:          updatedOrder.UpdatedAt,
		UpdatedBy:          req.UpdatedBy,
		UpdatedByName:      updatedByName,
		EstimatedReadyTime: updatedOrder.EstimatedReadyTime,
	}, nil
}

// checkAndUpdateOrderCompletion - Helper to check if all items are completed and update order status
// If items param is nil, it will fetch from DB; otherwise use provided items to avoid N+1 queries
func (s *Service) checkAndUpdateOrderCompletion(ctx context.Context, orderID int, items []*models.OrderItem) error {
	// Use provided items if available, otherwise fetch from DB
	var allItems []*models.OrderItem
	if items != nil {
		allItems = items
	} else {
		var err error
		allItems, err = s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("order_id = ?", orderID)
		})
		if err != nil || len(allItems) == 0 {
			return err
		}
	}

	// Check if all items are completed
	allCompleted := true
	for _, item := range allItems {
		if item.Status != "completed" {
			allCompleted = false
			break
		}
	}

	// If all items completed, update order status to completed
	if allCompleted {
		now := time.Now()
		_, err := s.orderRepo.UpdateColumns(ctx, orderID, map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
			"updated_at":   now,
		})
		if err != nil {
			return err
		}

		// Create timeline entry for order completion
		timeline := &models.OrderTimeline{
			OrderID:   orderID,
			Status:    "completed",
			UpdatedBy: "system",
			Note:      common.StrPtr("Order completed - all items finished"),
		}
		s.orderTimelineRepo.Create(ctx, timeline)
	}

	return nil
}

// UpdateOrderItemStatus - TASK-005: Update menu item status across multiple orders
func (s *Service) UpdateOrderItemStatus(ctx context.Context, menuItemID int, req models.UpdateOrderItemStatusRequest, userID, userName interface{}) (*models.OrderItemStatusUpdateBatchResponse, error) {
	// Validate order_ids are provided
	if len(req.OrderIDs) == 0 {
		return nil, fmt.Errorf("order_ids cannot be empty")
	}

	var updatedByName *string
	if userName != nil {
		name := userName.(string)
		updatedByName = &name
	}

	response := &models.OrderItemStatusUpdateBatchResponse{
		UpdatedItems: []models.OrderItemStatusUpdateResponse{},
	}

	// OPTIMIZATION: Fetch items for all orders at once
	// Use IN clause instead of individual queries
	allItems, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("menu_item_id = ? AND order_id IN ?", menuItemID, req.OrderIDs)
	})
	if err != nil || len(allItems) == 0 {
		return response, nil // No items found, return empty response
	}

	// Group items by order_id for easy lookup
	itemsByOrderID := make(map[int][]*models.OrderItem)
	for _, item := range allItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	// Update all matching items
	for _, item := range allItems {
		previousStatus := item.Status

		// Update item status
		_, err = s.orderItemRepo.UpdateColumns(ctx, item.ID, map[string]interface{}{
			"status":     req.Status,
			"updated_at": time.Now(),
		})
		if err != nil {
			continue
		}

		if req.Status == "completed" {
			// Update item in-memory for later check
			item.Status = req.Status
		}

		response.UpdatedItems = append(response.UpdatedItems, models.OrderItemStatusUpdateResponse{
			ItemID:         item.ID,
			OrderID:        item.OrderID,
			MenuItemName:   item.ItemName,
			Status:         req.Status,
			PreviousStatus: previousStatus,
			UpdatedAt:      time.Now(),
			UpdatedBy:      "kitchen",
			UpdatedByName:  updatedByName,
		})
	}

	// Check and update order completion for each affected order
	if req.Status == "completed" {
		for orderID, items := range itemsByOrderID {
			s.checkAndUpdateOrderCompletion(ctx, orderID, items)
		}
	}

	response.TotalUpdated = len(response.UpdatedItems)
	return response, nil
}

// UpdateOrderMultipleItemsStatus - TASK-005b: Update multiple items status in specific order
func (s *Service) UpdateOrderMultipleItemsStatus(ctx context.Context, orderID int, req models.UpdateOrderMultipleItemsStatusRequest, userID, userName interface{}) (*models.UpdateOrderMultipleItemsStatusResponse, error) {
	// Verify order exists
	_, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items cannot be empty")
	}

	response := &models.UpdateOrderMultipleItemsStatusResponse{
		UpdatedItems: []models.OrderItemStatusUpdateResponse{},
	}

	// OPTIMIZATION: Fetch ALL items for this order at once (avoid N+1 queries)
	allOrderItems, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id = ?", orderID)
	})
	if err != nil || len(allOrderItems) == 0 {
		return nil, fmt.Errorf("no items found for this order")
	}

	// Build map: menu_item_id -> []*OrderItem for fast lookup
	itemsByMenuID := make(map[int][]*models.OrderItem)
	for _, item := range allOrderItems {
		if item.MenuItemID != nil {
			itemsByMenuID[*item.MenuItemID] = append(itemsByMenuID[*item.MenuItemID], item)
		}
	}

	// Track if any item is being updated to completed status
	hasCompletedItems := false

	// Update each item from request
	for _, reqItem := range req.Items {
		items := itemsByMenuID[reqItem.MenuItemID]
		if len(items) == 0 {
			continue // Skip if item not found
		}

		for _, item := range items {
			previousStatus := item.Status

			// Update item status
			_, err = s.orderItemRepo.UpdateColumns(ctx, item.ID, map[string]interface{}{
				"status":     reqItem.Status,
				"updated_at": time.Now(),
			})
			if err != nil {
				continue
			}

			if reqItem.Status == "completed" {
				hasCompletedItems = true
				// Update item in-memory for later check
				item.Status = reqItem.Status
			}

			response.UpdatedItems = append(response.UpdatedItems, models.OrderItemStatusUpdateResponse{
				ItemID:         item.ID,
				OrderID:        orderID,
				MenuItemName:   item.ItemName,
				Status:         reqItem.Status,
				PreviousStatus: previousStatus,
				UpdatedAt:      time.Now(),
				UpdatedBy:      "kitchen",
				UpdatedByName:  nil,
			})
		}
	}

	// Check if all items are completed and auto-update order status
	// Pass allOrderItems to avoid additional query
	if hasCompletedItems {
		s.checkAndUpdateOrderCompletion(ctx, orderID, allOrderItems)
	}

	response.TotalUpdated = len(response.UpdatedItems)
	return response, nil
}

// UpdateOrder - TASK-006: Add notes/metadata to order
func (s *Service) UpdateOrder(ctx context.Context, orderID int, req models.UpdateOrderRequest) (*models.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Note != nil {
		currentNote := ""
		if order.Notes != nil {
			currentNote = *order.Notes + ". "
		}
		newNote := currentNote + *req.Note
		updates["notes"] = newNote
	}

	if req.Metadata != nil {
		metaJSON, _ := json.Marshal(req.Metadata)
		metaStr := string(metaJSON)
		updates["meta"] = metaStr
	}

	updatedOrder, err := s.orderRepo.UpdateColumns(ctx, orderID, updates)
	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, updatedOrder.ID, "admin")
}

// CancelOrder - TASK-007: Cancel order
func (s *Service) CancelOrder(ctx context.Context, orderID int, req models.CancelOrderRequest, userID, userRole interface{}) (*models.CancelOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	// Check if order can be cancelled
	if order.Status == "completed" || order.Status == "cancelled" {
		return nil, fmt.Errorf("cannot cancel order. Order is already %s", order.Status)
	}

	if order.Status == "preparing" || order.Status == "ready" {
		return nil, errors.New("cannot cancel order. Order is already being prepared")
	}

	// Update order
	cancelledBy := "customer"
	if userRole != nil {
		cancelledBy = userRole.(string)
	}

	now := time.Now()
	_, err = s.orderRepo.UpdateColumns(ctx, orderID, map[string]interface{}{
		"status":        "cancelled",
		"cancelled_at":  now,
		"cancelled_by":  cancelledBy,
		"cancel_reason": req.Reason,
		"updated_at":    now,
	})
	if err != nil {
		return nil, err
	}

	// Create timeline entry
	timeline := &models.OrderTimeline{
		OrderID:   orderID,
		Status:    "cancelled",
		UpdatedBy: cancelledBy,
		Note:      &req.Reason,
	}
	s.orderTimelineRepo.Create(ctx, timeline)

	return &models.CancelOrderResponse{
		ID:           order.ID,
		OrderNumber:  order.OrderNumber,
		Status:       "cancelled",
		CancelledAt:  now,
		CancelledBy:  cancelledBy,
		CancelReason: req.Reason,
		RefundStatus: "pending",
		RefundAmount: order.Total,
	}, nil
}

// SendKitchenAlert - TASK-008: Send alert to waiter
func (s *Service) SendKitchenAlert(ctx context.Context, orderID int, req models.CreateAlertRequest) (*models.AlertResponse, error) {
	order, err := s.orderRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("id = ?", orderID)
		tx.Preload("Waiter")
	})
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	alert := &models.KitchenAlert{
		OrderID:   orderID,
		AlertType: req.AlertType,
		Message:   req.Message,
		Priority:  priority,
		SentTo:    "waiter",
		WaiterID:  order.WaiterID,
		Status:    "sent",
	}

	createdAlert, err := s.kitchenAlertRepo.Create(ctx, alert)
	if err != nil {
		return nil, err
	}

	var waiterName *string
	if order.Waiter != nil {
		name := fmt.Sprintf("%s %s", *order.Waiter.FirstName, *order.Waiter.LastName)
		waiterName = &name
	}

	return &models.AlertResponse{
		AlertID:     createdAlert.ID,
		OrderID:     orderID,
		OrderNumber: order.OrderNumber,
		Message:     req.Message,
		SentTo:      "waiter",
		WaiterID:    order.WaiterID,
		WaiterName:  waiterName,
		SentAt:      createdAlert.SentAt,
		Status:      createdAlert.Status,
	}, nil
}

// CreateOrderReview - TASK-009: Submit review for completed order
func (s *Service) CreateOrderReview(ctx context.Context, orderID int, customerID string, req models.CreateReviewRequest) (interface{}, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, common.ErrOrderNotFound
	}

	// Verify order is completed
	if order.Status != "completed" {
		return nil, errors.New("can only review completed orders")
	}

	// Verify customer owns this order
	if order.CustomerUserID == nil || *order.CustomerUserID != customerID {
		return nil, errors.New("unauthorized to review this order")
	}

	// TODO: Create review in customer_reviews table
	// This would need the CustomerReview model and repository

	return map[string]interface{}{
		"review_id":     4001,
		"order_id":      orderID,
		"order_number":  order.OrderNumber,
		"customer_id":   customerID,
		"customer_name": *order.CustomerName,
		"rating":        req.Rating,
		"comment":       req.Comment,
		"photos":        req.Photos,
		"created_at":    time.Now(),
		"status":        "published",
	}, nil
}

// GetOrderItemsSummaryByCategory - TASK-010: Get items summary grouped by category for kitchen station
func (s *Service) GetOrderItemsSummaryByCategory(ctx context.Context) ([]*models.OrderItemsSummaryResponse, error) {
	// Get active orders only (pending, accepted, preparing)
	orders, err := s.orderRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("status IN ?", []string{"pending", "accepted", "preparing"}).
			Preload("Table")
	})
	if err != nil || len(orders) == 0 {
		return []*models.OrderItemsSummaryResponse{}, err
	}

	// Get order items for these orders (excluding cancelled and completed)
	var orderIDs []int
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	orderItems, err := s.orderItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("order_id IN ? AND status NOT IN ?", orderIDs, []string{"cancelled", "completed"}).
			Preload("Modifiers")
	})
	if err != nil || len(orderItems) == 0 {
		return []*models.OrderItemsSummaryResponse{}, err
	}

	// Map items to orders for easier processing
	orderItemsMap := make(map[int][]*models.OrderItem)
	for _, item := range orderItems {
		orderItemsMap[item.OrderID] = append(orderItemsMap[item.OrderID], item)
	}

	menuItemIDSet := make(map[int]bool)
	for _, items := range orderItemsMap {
		for _, item := range items {
			menuItemIDSet[*item.MenuItemID] = true
		}
	}

	menuItemIDs := make([]int, 0, len(menuItemIDSet))
	for id := range menuItemIDSet {
		menuItemIDs = append(menuItemIDs, id)
	}

	// Fetch menu items with categories
	menuItems, err := s.menuItemRepo.ListByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("id IN ?", menuItemIDs).
			Preload("Category")
	})
	if err != nil {
		menuItems = []*models.MenuItem{}
		return nil, err
	}

	// Build map of menu item to category
	itemCategoryMap := make(map[int]string)
	for _, item := range menuItems {
		if item.Category != nil {
			itemCategoryMap[item.ID] = item.Category.Name
		}
	}

	// Group items by (menu_item_id, category)
	// Structure: map[category]map[menu_item_id]ItemData
	categoryItemsMap := make(map[string]map[int]*models.OrderItemSummaryWithTable)

	for _, order := range orders {
		items := orderItemsMap[order.ID]
		for _, item := range items {
			category := itemCategoryMap[*item.MenuItemID]
			if category == "" {
				category = "Other"
			}

			// Initialize category map if needed
			if categoryItemsMap[category] == nil {
				categoryItemsMap[category] = make(map[int]*models.OrderItemSummaryWithTable)
			}

			// Initialize item if needed
			if categoryItemsMap[category][*item.MenuItemID] == nil {
				categoryItemsMap[category][*item.MenuItemID] = &models.OrderItemSummaryWithTable{
					ItemID:      item.ID,
					MenuItemID:  *item.MenuItemID,
					ItemName:    item.ItemName,
					TotalQty:    0,
					OrderTables: []models.OrderItemTableInfo{},
				}
			}

			// Add to total quantity
			categoryItemsMap[category][*item.MenuItemID].TotalQty += item.Quantity

			// Add table info
			tableNumber := "Unknown"
			if order.Table != nil {
				tableNumber = fmt.Sprintf("Table %d", order.Table.TableNumber)
			}

			categoryItemsMap[category][*item.MenuItemID].OrderTables = append(
				categoryItemsMap[category][*item.MenuItemID].OrderTables,
				models.OrderItemTableInfo{
					OrderID:     order.ID,
					TableID:     order.TableID,
					TableNumber: tableNumber,
					Quantity:    item.Quantity,
					Status:      item.Status,
				},
			)
		}
	}

	// Build response
	result := make([]*models.OrderItemsSummaryResponse, 0)

	// Define category order
	categoryOrder := []string{"Appetizer", "Soup", "Salad", "Main Course", "Side", "Dessert", "Beverage", "Other"}
	categoryMap := make(map[string]bool)
	for _, cat := range categoryOrder {
		if len(categoryItemsMap[cat]) > 0 {
			categoryMap[cat] = true
		}
	}

	for _, cat := range categoryOrder {
		if !categoryMap[cat] {
			continue
		}

		items := make([]models.OrderItemSummaryWithTable, 0)
		for _, itemSummary := range categoryItemsMap[cat] {
			items = append(items, *itemSummary)
		}

		result = append(result, &models.OrderItemsSummaryResponse{
			Category: cat,
			Items:    items,
		})
	}

	// Add any remaining categories not in order
	for cat, items := range categoryItemsMap {
		found := false
		for _, orderedCat := range categoryOrder {
			if cat == orderedCat {
				found = true
				break
			}
		}
		if !found && len(items) > 0 {
			itemsList := make([]models.OrderItemSummaryWithTable, 0)
			for _, itemSummary := range items {
				itemsList = append(itemsList, *itemSummary)
			}
			result = append(result, &models.OrderItemsSummaryResponse{
				Category: cat,
				Items:    itemsList,
			})
		}
	}

	return result, nil
}
