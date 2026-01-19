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

// GetTablesForStaff - Get tables with orders for staff view (waiter/kitchen)
func (h *Handler) GetTablesForStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params = models.ListTablesForStaffRequest{}
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetTablesForStaff(c.Request.Context(), &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Tables retrieved successfully", data))
	}
}

// GetTableDetailForStaff - Get table detail with all order items for staff
func (h *Handler) GetTableDetailForStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.TableParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetTableDetailForStaff(c.Request.Context(), params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusOK, common.BaseResponseMess(http.StatusOK, "Table detail retrieved successfully", data))
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

func (h *Handler) GetQrCodeByTableId() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		qrCodeInfo, err := h.service.GetQrCodeByTableID(c, id)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{
			"table_id":  id,
			"token":     qrCodeInfo.Token,
			"create_at": qrCodeInfo.CreatedAt,
			"expire_at": qrCodeInfo.ExpiresAt,
		}))
	}
}

func (h *Handler) GenerateQrCodeByTableId() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GenerateQrCodeByTableId(c, id)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) DownloadQrCodeByTableId() gin.HandlerFunc {
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

		url := fmt.Sprintf("https://smart-restaurant-fe.vercel.app/menu?table=%d&token=%s", table.ID, token)

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

func (h *Handler) DownloadAllQrCode() gin.HandlerFunc {
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

			url := fmt.Sprintf("https://smart-restaurant-fe.vercel.app/menu?table=%d&token=%s", table.ID, table.QrToken)
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
