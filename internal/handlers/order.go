package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder handles POST /api/orders
// TASK-001: Create/Submit order from cart
func (h *Handler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get customer ID from JWT if user is authenticated (optional)
	if ok, profile := common.ProfileFromJwt(c); ok {
		req.CustomerID = &profile.Id
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.BaseResponseMess(http.StatusCreated, "Order created successfully", order))
}

// GetOrders handles GET /api/orders
// TASK-002: Get orders list with dynamic filtering
func (h *Handler) GetOrders(c *gin.Context) {
	var params models.ListOrdersRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		common.AbortWithError(c, err)
		return
	}

	result, err := h.service.GetOrders(c.Request.Context(), &params)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Orders retrieved successfully", result))
}

// GetOrderByID handles GET /api/orders/:id
// TASK-003: Get order details (role-based response)
func (h *Handler) GetOrderByID(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get user role from context (if authenticated)
	role := c.DefaultQuery("role", "customer")
	userRole, exists := c.Get("user_role")
	if exists && userRole != nil {
		role = userRole.(string)
	}

	order, err := h.service.GetOrderByID(c.Request.Context(), orderID, role)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order retrieved successfully", order))
}

// UpdateOrderStatus handles PATCH /api/orders/:id/status
// TASK-004: Update order status (role-based workflow)
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	result, err := h.service.UpdateOrderStatus(c.Request.Context(), orderID, req, userID, userName)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order status updated successfully", result))
}

// UpdateOrderItemStatus handles PATCH /api/orders/items/:itemId/status
// TASK-005: Update individual item status (batch update across multiple orders)
func (h *Handler) UpdateOrderItemStatus(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.UpdateOrderItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	result, err := h.service.UpdateOrderItemStatus(c.Request.Context(), itemID, req, userID, userName)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order item status updated successfully", result))

}

// UpdateOrderMultipleItemsStatus handles PATCH /api/orders/:id/items/status
// TASK-005b: Update multiple items status in specific order
func (h *Handler) UpdateOrderMultipleItemsStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.UpdateOrderMultipleItemsStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	result, err := h.service.UpdateOrderMultipleItemsStatus(c.Request.Context(), orderID, req, userID, userName)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order items status updated successfully", result))
}

// UpdateOrder handles PATCH /api/orders/:id
// TASK-006: Add notes/metadata to order
func (h *Handler) UpdateOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	result, err := h.service.UpdateOrder(c.Request.Context(), orderID, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order updated successfully", result))
}

// CancelOrder handles POST /api/orders/:id/cancel
// TASK-007: Cancel order
func (h *Handler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	result, err := h.service.CancelOrder(c.Request.Context(), orderID, req, userID, userRole)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Order cancelled successfully", result))
}

// CallStaff handles POST /api/tables/:id/call-staff
// Allows customer to request staff assistance for a table
func (h *Handler) CallStaff(c *gin.Context) {
	tableID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	result, err := h.service.CallStaff(c.Request.Context(), tableID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff called successfully", result))
}

// SendKitchenAlert handles POST /api/orders/:id/alert
// TASK-008: Send alert to waiter (item ready)
func (h *Handler) SendKitchenAlert(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	result, err := h.service.SendKitchenAlert(c.Request.Context(), orderID, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Alert sent successfully", result))
}

// CreateOrderReview handles POST /api/orders/:id/review
// TASK-009: Submit review for completed order
func (h *Handler) CreateOrderReview(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get customer ID from JWT - authentication required for reviews
	ok, profile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	result, err := h.service.CreateOrderReview(c.Request.Context(), orderID, profile.Id, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.BaseResponseMess(http.StatusCreated, "Review submitted successfully", result))
}

// GetOrderItemsSummaryByCategory handles GET /api/orders/items-summary
// TASK-010: Get items summary grouped by category for kitchen station
func (h *Handler) GetOrderItemsSummaryByCategory(c *gin.Context) {
	result, err := h.service.GetOrderItemsSummaryByCategory(c.Request.Context())
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Items summary retrieved successfully", result))
}
