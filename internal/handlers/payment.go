package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProcessPayment handles POST /api/payments
// TASK-013: Process payment (VN-PAY, Cash, Card, or E-wallet)
func (h *Handler) ProcessPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ProcessPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		// Get client IP
		req.ClientIP = c.ClientIP()

		payment, err := h.service.ProcessPayment(c.Request.Context(), &req)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusCreated, common.BaseResponseMess(http.StatusCreated, "Payment processed successfully", payment))
	}
}

// GetPaymentStatus handles GET /api/payments/:id/status
// TASK-014: Check payment status
func (h *Handler) GetPaymentStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("id")
		if paymentID == "" {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		status, err := h.service.GetPaymentStatus(c.Request.Context(), paymentID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Payment status retrieved successfully", status))
	}
}
