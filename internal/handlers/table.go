package handlers

import (
	"app-noti/common"
	"app-noti/internal/models"

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
