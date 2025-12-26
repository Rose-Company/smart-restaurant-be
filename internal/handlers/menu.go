package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LoadMenu() gin.HandlerFunc {
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

func (h *Handler) GetMenuCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params = models.ListMenuCategoryRequest{}
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetMenuCategories(c, &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) GetMenuCategoryByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuCategoryParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetMenuCategoryByID(c, params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) CreateMenuCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateMenuCategoryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.CreateMenuCategory(c, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusCreated, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateMenuCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuCategoryParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateMenuCategoryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateMenuCategory(c, params.ID, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateMenuCategoryStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuCategoryParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateMenuCategoryStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateMenuCategoryStatus(c, params.ID, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}
