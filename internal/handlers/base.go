package handlers

import (
	services "app-noti/internal/services"
	l "app-noti/pkg/logger"
	"app-noti/server"

	"github.com/gin-gonic/gin"
)

var ll = l.New()

type Handler struct {
	sc      server.ServerContext
	service *services.Service
}

func NewHandler(sc server.ServerContext) *Handler {
	return &Handler{
		sc:      sc,
		service: services.NewService(sc),
	}
}

func (h *Handler) RegisterRouter(c *gin.Engine) {
	// authConfig := h.sc.GetAuthConfig()
	// authenticator := middleware.NewAuthenticator(authConfig)

	admin := c.Group("/api/admin")
	{
		admin.POST("/upload", h.UploadImage())
		admin.GET("/tables", h.GetTables())
		admin.GET("/tables/:id", h.GetTableByID())
		admin.POST("/tables", h.CreateTable())
		admin.PUT("/tables/:id", h.UpdateTable())
		admin.PATCH("/tables/:id/status", h.UpdateTableStatus())
		admin.POST("/tables/:id/qr/generate", h.GenerateQrCodeByTableId())
		admin.GET("tables/:id/qr/download", h.DownloadQrCodeByTableId())
		admin.GET("tables/qr/download-all", h.DownloadAllQrCode())
		admin.GET("tables/:id/qr", h.GetQrCodeByTableId())

		menuAdmin := admin.Group("/menu")
		{
			menuAdmin.GET("/categories", h.GetMenuCategories())
			menuAdmin.GET("/categories/:id", h.GetMenuCategoryByID())
			menuAdmin.POST("/categories", h.CreateMenuCategory())
			menuAdmin.PUT("/categories/:id", h.UpdateMenuCategory())
			menuAdmin.PATCH("/categories/:id/status", h.UpdateMenuCategoryStatus())

			itemsAdmin := menuAdmin.Group("/items")
			{
				itemsAdmin.GET("", h.GetMenuItems())
				itemsAdmin.GET("/:id", h.GetMenuItemByID())
				itemsAdmin.POST("", h.CreateMenuItem())
				itemsAdmin.PUT("/:id", h.UpdateMenuItem())
				itemsAdmin.DELETE("/:id", h.DeleteMenuItem())
			}
		}
	}

	menu := c.Group("/api/menu")
	{
		menu.GET("", h.LoadMenu())
	}

}
