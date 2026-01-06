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

	userRoutes := c.Group("/api/user")
	{
		userRoutes.POST("/signup", h.SignUp)
		userRoutes.POST("/login", h.LogIn)
	}

	authRoutes := c.Group("/api/auth")
	{
		authRoutes.POST("/request-signup-otp", h.RequestSignupOTP)
		authRoutes.POST("/validate-signup-otp", h.ValidateSignupOTP)
		authRoutes.POST("/request-reset-password", h.RequestResetPassword)
		authRoutes.POST("/validate-otp", h.ValidateOTP)
		authRoutes.POST("/reset-password", h.ResetPassword)
	}

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

			modifiersGroupAdmin := menuAdmin.Group("/modifier-groups")
			{
				modifiersGroupAdmin.GET("", h.GetModifierGroup())
				modifiersGroupAdmin.POST("", h.CreatModifierGroup())
				modifiersGroupAdmin.PUT("/:id", h.UpdateModifierGroup())
				modifiersGroupAdmin.DELETE("/:id", h.DeleteModifierGroup())
				modifiersGroupAdmin.POST("/:id/options", h.CreateModifierOptions())
			}

			modifiersOptionsAdmin := menuAdmin.Group("/modifier-options")
			{
				modifiersOptionsAdmin.PUT("/:id", h.UpdateModifierOptions())
				modifiersOptionsAdmin.DELETE("/:id", h.DeleteModifierOptions())
			}
		}
	}

	menu := c.Group("/api/menu")
	{
		menu.GET("", h.LoadMenu())

		menuItem := menu.Group("/items")
		{
			menuItem.GET("/:id", h.GetMenuItemByID())
			menuItem.POST("/:id/modifier-groups", h.AssignMenuItemModifierGroup())
			menuItem.DELETE("/:id/modifier-groups/:groupId", h.DeleteMenuItemModifierGroup())
		}
	}

}
