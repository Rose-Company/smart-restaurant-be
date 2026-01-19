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

// HandleVNPayCallback handles GET /api/vnpay/callback
// Callback from VN-PAY after payment
func (h *Handler) HandleVNPayCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract VN-PAY response parameters
		params := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		// Process VN-PAY callback
		result, err := h.service.HandleVNPayCallback(c.Request.Context(), params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "VN-PAY callback processed", result))
	}
}
