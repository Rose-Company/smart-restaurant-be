package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListStaff - TASK-021: GET /api/admin/staff
func (h *Handler) ListStaff(c *gin.Context) {
	var req models.ListStaffRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Call service
	result, err := h.service.ListStaff(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff list retrieved successfully", result))
}

// CreateStaff - TASK-022: POST /api/admin/staff
func (h *Handler) CreateStaff(c *gin.Context) {
	var req models.CreateStaffRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get admin ID from context
	adminID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Call service
	result, err := h.service.CreateStaff(c.Request.Context(), req, adminID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff account created successfully", result))
}

// GetStaffByID - TASK-023: GET /api/admin/staff/:id
func (h *Handler) GetStaffByID(c *gin.Context) {
	staffID := c.Param("id")

	// Call service
	result, err := h.service.GetStaffByID(c.Request.Context(), staffID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff details retrieved successfully", result))
}

// UpdateStaff - TASK-024: PUT /api/admin/staff/:id
func (h *Handler) UpdateStaff(c *gin.Context) {
	staffID := c.Param("id")

	var req models.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get admin ID from context
	adminID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Call service
	result, err := h.service.UpdateStaff(c.Request.Context(), staffID, req, adminID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff account updated successfully", result))
}

// DeleteStaff - TASK-025: DELETE /api/admin/staff/:id
func (h *Handler) DeleteStaff(c *gin.Context) {
	staffID := c.Param("id")

	// Call service
	err := h.service.DeleteStaff(c.Request.Context(), staffID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff account deleted successfully", nil))
}

// SendStaffInvite - TASK-026: POST /api/admin/staff/:id/send-invite
func (h *Handler) SendStaffInvite(c *gin.Context) {
	staffID := c.Param("id")

	var req models.SendInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get admin ID from context
	adminID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Call service
	err := h.service.SendStaffInvite(c.Request.Context(), staffID, req, adminID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Invitation sent successfully", map[string]interface{}{
		"email": req.Email,
		"role":  req.Role,
		"sent":  true,
	}))
}

// AssignTablesToWaiter - TASK-027: PATCH /api/admin/staff/:id/assign-tables
func (h *Handler) AssignTablesToWaiter(c *gin.Context) {
	staffID := c.Param("id")

	var req models.AssignTablesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Get admin ID from context
	adminID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Call service
	err := h.service.AssignTablesToWaiter(c.Request.Context(), staffID, req, adminID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Tables assigned successfully", map[string]interface{}{
		"waiter_id": staffID,
		"tables":    req.TableIDs,
	}))
}

// GetStaffProfile - TASK-028: GET /api/staff/profile
func (h *Handler) GetStaffProfile(c *gin.Context) {
	// Get staff ID from context (from JWT)
	staffID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	// Call service
	result, err := h.service.GetStaffProfile(c.Request.Context(), staffID.(string))
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff profile retrieved successfully", result))
}

// UpdateStaffProfile - TASK-029: PUT /api/staff/profile
func (h *Handler) UpdateStaffProfile(c *gin.Context) {
	// Get staff ID from context (from JWT)
	staffID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	var req models.UpdateStaffProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Call service
	result, err := h.service.UpdateStaffProfile(c.Request.Context(), staffID.(string), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Staff profile updated successfully", result))
}

// ChangeStaffPassword - TASK-030: PATCH /api/staff/password
func (h *Handler) ChangeStaffPassword(c *gin.Context) {
	// Get staff ID from context (from JWT)
	staffID, exists := c.Get("user_id")
	if !exists {
		common.AbortWithError(c, common.ErrUnauthorized)
		return
	}

	var req models.ChangeStaffPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	// Call service
	err := h.service.ChangeStaffPassword(c.Request.Context(), staffID.(string), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Password changed successfully", map[string]interface{}{
		"changed_at": time.Now().Format(time.RFC3339),
	}))
}
