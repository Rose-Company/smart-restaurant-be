package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBill handles POST /api/bills
// TASK-010: Create bill from order
func (h *Handler) CreateBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateBillRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		// Get user info if authenticated
		if ok, _ := common.ProfileFromJwt(c); ok {
			req.RequestedBy = "customer"
		}

		bill, err := h.service.CreateBill(c.Request.Context(), &req)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusCreated, common.BaseResponseMess(http.StatusCreated, "Bill created successfully", bill))
	}
}

// GetBill handles GET /api/bills/:id
// TASK-011: Get bill details with format option
func (h *Handler) GetBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		billID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		var req models.GetBillRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		// Default format is json
		if req.Format == "" {
			req.Format = "json"
		}

		bill, err := h.service.GetBill(c.Request.Context(), billID, req.Format)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Bill retrieved successfully", bill))
	}
}

// UpdateBill handles PATCH /api/bills/:id
// TASK-012: Update bill (add discount, mark paid)
func (h *Handler) UpdateBill() gin.HandlerFunc {
	return func(c *gin.Context) {
		billID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		var req models.UpdateBillRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		bill, err := h.service.UpdateBill(c.Request.Context(), billID, &req)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Bill updated successfully", bill))
	}
}
