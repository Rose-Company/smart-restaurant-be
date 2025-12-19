package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) loadMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableIdStr := c.Query("table")
		token := c.Query("token")

		if tableIdStr == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"menu": false, "error": "table and token required"})
			return
		}

		tableId, err := strconv.Atoi(tableIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"menu": false, "error": "invalid table id"})
			return
		}

		table, err := h.service.GetTableByID(c, tableId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"menu": false})
			return
		}

		if table.QrToken != token || table.QrTokenExpiresAt == nil || time.Now().After(*table.QrTokenExpiresAt) {
			c.JSON(http.StatusForbidden, gin.H{"menu": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"menu": true})
	}
}
