package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateDiscount handles POST /api/discounts/validate
// TASK-015: Validate discount code before payment
func (h *Handler) ValidateDiscount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ValidateDiscountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.AbortWithError(c, common.ErrInvalidInput)
			return
		}

		// For now, validate discount with 0 amount - the service will handle order lookup
		// In production, you'd fetch the order to get the actual total
		discount, err := h.service.ValidateDiscountCode(c.Request.Context(), req.Code, req.OrderID, 0)
		if err != nil {
			// Return specific error responses for different scenarios
			if err == common.ErrExpiredDiscount {
				errorResponse := &models.ValidateDiscountErrorResponse{
					Code:        req.Code,
					IsValid:     false,
					ErrorReason: "expired",
				}
				c.JSON(http.StatusBadRequest, common.BaseResponseMess(http.StatusBadRequest, "Discount code has expired", errorResponse))
				return
			}

			if err == common.ErrMinOrderNotMet {
				errorResponse := &models.ValidateDiscountErrorResponse{
					Code:        req.Code,
					IsValid:     false,
					ErrorReason: "min_order_not_met",
				}
				c.JSON(http.StatusBadRequest, common.BaseResponseMess(http.StatusBadRequest, "Order amount does not meet minimum requirement", errorResponse))
				return
			}

			if err == common.ErrInvalidDiscount {
				errorResponse := &models.ValidateDiscountErrorResponse{
					Code:        req.Code,
					IsValid:     false,
					ErrorReason: "code_not_found",
				}
				c.JSON(http.StatusBadRequest, common.BaseResponseMess(http.StatusBadRequest, "Discount code is invalid", errorResponse))
				return
			}

			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Discount code is valid", discount))
	}
}
