package handlers

import (
	"app-noti/common"
	"net/http"

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
		// Extract VN-PAY callback parameters from query string
		params := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		// Validate required parameters
		if params["vnp_ResponseCode"] == "" || params["vnp_OrderInfo"] == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing required VN-PAY callback parameters",
			})
			return
		}

		// Handle callback using service
		response, err := h.service.HandleVNPayCallback(c.Request.Context(), params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		// Return success response
		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, response.Message, response))
	}
}
