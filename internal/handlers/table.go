package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
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
