package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	storage "app-noti/services/digital_ocean_storage"
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

		restaurantId := table.RestaurantId
		menuItemsResponse, err := h.service.GetMenuItemsByRestaurant(c, restaurantId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, menuItemsResponse)
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

func (h *Handler) GetMenuItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.ListMenuItemRequest
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetMenuItems(c, &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) GetMenuItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuItemIDParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetMenuItemByID(c, params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) CreateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateMenuItemRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.CreateMenuItem(c, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(http.StatusCreated, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuItemIDParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateMenuItemRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateMenuItem(c, params.ID, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) DeleteMenuItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.MenuItemIDParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		err := h.service.DeleteMenuItem(c, params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{"message": "Menu item deleted successfully"}))
	}
}

func (h *Handler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		contentType := file.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG and JPG images are allowed"})
			return
		}

		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 10MB"})
			return
		}

		err, doStorage := storage.NewDOStorage("menu-items")
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		if err := doStorage.Run(); err != nil {
			common.AbortWithError(c, err)
			return
		}

		fileURL, err := doStorage.UploadFile(file, "smart-restaurant")
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{
			"url": fileURL,
		}))
	}
}

func (h *Handler) AssignMenuItemModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.AssignModifierToMenuItemRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.AssignMenuItemModifierGroup(c, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) DeleteMenuItemModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

		var uri models.DeleteMenuItemModifierGroupUri
		if err := c.ShouldBindUri(&uri); err != nil {
			common.AbortWithError(c, err)
			return
		}

		if err := h.service.DeleteMenuItemModifierGroup(c.Request.Context(), uri.MenuItemID, uri.GroupID); err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(nil))
	}
}
