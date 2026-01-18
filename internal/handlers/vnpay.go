package handlers

import (
	"app-noti/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateVNPayPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData map[string]interface{}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		resp, err := h.service.CreateVNPayPaymentRes(requestData, c.Request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(200, resp)
	}
}

func (h *Handler) VnpayCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		status := c.Query("vnp_ResponseCode")
		orderID := c.Query("vnp_OrderInfo")
		amountStr := c.Query("vnp_Amount")

		if status == "" || orderID == "" || amountStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required params",
			})
			return
		}

		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid amount",
			})
			return
		}
		amount = amount / 100

		if status == "00" {
			// TODO: handle update bill status

			c.JSON(http.StatusOK, gin.H{
				"code":    status,
				"amount":  amount,
				"message": "success",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    status,
			"message": "payment failed",
		})
	}
}
