package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func (h *Handler) GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params = models.ListTablesRequest{}
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetTables(c, &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) GetTableByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetTableByID(c, id)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateTableRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.CreateTable(c, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateTableRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateTable(c, id, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateTableStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateTableStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateTableStatus(c, id, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) generateQrCodeByTableId() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		url, err := h.service.GenerateQrCodeByTableId(c, id)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{
			"url": url,
		}))
	}
}

func (h *Handler) downloadQrCodeByTableId() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		token := c.Query("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		table, err := h.service.GetTableByID(c, id)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		if table.QrToken != token || table.QrTokenExpiresAt == nil || time.Now().After(*table.QrTokenExpiresAt) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid or expired token"})
			return
		}

		url := fmt.Sprintf("https://smart-restaurant-fe.vercel.app//menu?table=%d&token=%s", table.ID, token)

		png, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=table_%d_qr.png", table.ID))
		c.Writer.Write(png)
	}
}

func (h *Handler) downloadAllQrCode() gin.HandlerFunc {
	return func(c *gin.Context) {

		tables, err := h.service.GetAllTables(c)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		for _, table := range tables {
			if table.QrToken == "" || table.QrTokenExpiresAt == nil || time.Now().After(*table.QrTokenExpiresAt) {
				continue
			}

			url := fmt.Sprintf("https://restaurant-domain.com/menu?table=%d&token=%s", table.ID, table.QrToken)
			png, err := qrcode.Encode(url, qrcode.Medium, 256)
			if err != nil {
				zipWriter.Close()
				common.AbortWithError(c, err)
				return
			}

			f, err := zipWriter.Create(fmt.Sprintf("table_%d.png", table.ID))
			if err != nil {
				zipWriter.Close()
				common.AbortWithError(c, err)
				return
			}
			f.Write(png)
		}

		zipWriter.Close()

		c.Header("Content-Type", "application/zip")
		c.Header("Content-Disposition", "attachment; filename=all_tables_qr.zip")
		c.Writer.Write(buf.Bytes())
	}
}
