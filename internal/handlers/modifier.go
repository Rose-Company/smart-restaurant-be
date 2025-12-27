package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params = models.ListModifierGroupRequest{}
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.GetModifierGroup(c, &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) CreatModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateModifierGroupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.CreateModifierGroup(c, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateModifierGroupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateModifierGroup(c, id, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) DeleteModifierGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.ModifierGroupIDParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		err := h.service.DeleteModifierGroup(c, params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{"message": "Modifier Group deleted successfully"}))
	}
}

func (h *Handler) CreateModifierOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		var params = models.CreateModifierOptionRequest{}
		if err := c.ShouldBindQuery(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.CreateModifierOptions(c, id, &params)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) UpdateModifierOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		var request models.UpdateModifierOptionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			common.AbortWithError(c, err)
			return
		}

		data, err := h.service.UpdateModifierOptions(c, id, &request)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) DeleteModifierOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.ModifierOptionIDParamsUri
		if err := c.ShouldBindUri(&params); err != nil {
			common.AbortWithError(c, err)
			return
		}

		err := h.service.DeleteModifierOptions(c, params.ID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(gin.H{"message": "Modifier Option deleted successfully"}))
	}
}
