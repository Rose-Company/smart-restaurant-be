package handlers

import (
	services "app-noti/internal/services"
	"app-noti/middleware"
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
			menuAdmin.DELETE("/categories/:id", h.DeleteMenuCategory())

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
		menu.GET("/search", h.LoadMenu())

		menuItem := menu.Group("/items")
		{
			menuItem.GET("/:id", h.GetMenuItemByID())
			menuItem.POST("/:id/modifier-groups", h.AssignMenuItemModifierGroup())
			menuItem.DELETE("/:id/modifier-groups/:groupId", h.DeleteMenuItemModifierGroup())
		}
	}

	// Order Management APIs (TASK-001 to TASK-009)
	// Using optional authentication - guests can order via QR code, logged-in users get personalized experience
	orders := c.Group("/api/orders")
	{
		orders.POST("", middleware.OptionalUserAuthentication(), h.CreateOrder)                                     // TASK-001: Create order
		orders.GET("", middleware.OptionalUserAuthentication(), h.GetOrders)                                        // TASK-002: Get orders list
		orders.GET("/:id", middleware.OptionalUserAuthentication(), h.GetOrderByID)                                 // TASK-003: Get order details
		orders.PATCH("/:id/status", middleware.OptionalUserAuthentication(), h.UpdateOrderStatus)                   // TASK-004: Update order status
		orders.PATCH("/:id/items/:itemId/status", middleware.OptionalUserAuthentication(), h.UpdateOrderItemStatus) // TASK-005: Update item status
		orders.PATCH("/:id", middleware.OptionalUserAuthentication(), h.UpdateOrder)                                // TASK-006: Add notes/metadata
		orders.POST("/:id/cancel", middleware.OptionalUserAuthentication(), h.CancelOrder)                          // TASK-007: Cancel order
		orders.POST("/:id/alert", middleware.OptionalUserAuthentication(), h.SendKitchenAlert)                      // TASK-008: Send kitchen alert
		orders.POST("/:id/review", middleware.OptionalUserAuthentication(), h.CreateOrderReview)                    // TASK-009: Submit review
	}

	// Bill Management APIs (TASK-010 to TASK-012)
	bills := c.Group("/api/bills")
	{
		bills.POST("", middleware.OptionalUserAuthentication(), h.CreateBill())      // TASK-010: Create bill from order
		bills.GET("/:id", middleware.OptionalUserAuthentication(), h.GetBill())      // TASK-011: Get bill details
		bills.PATCH("/:id", middleware.OptionalUserAuthentication(), h.UpdateBill()) // TASK-012: Update bill (add discount, mark paid)
	}

	// Payment Management APIs (TASK-013 to TASK-014)
	payments := c.Group("/api/payments")
	{
		payments.POST("", middleware.OptionalUserAuthentication(), h.ProcessPayment())             // TASK-013: Process payment
		payments.GET("/:id/status", middleware.OptionalUserAuthentication(), h.GetPaymentStatus()) // TASK-014: Check payment status
	}

	// Discount Management APIs (TASK-015)
	discounts := c.Group("/api/discounts")
	{
		discounts.POST("/validate", middleware.OptionalUserAuthentication(), h.ValidateDiscount()) // TASK-015: Validate discount code
	}

	payment := c.Group("/api/payment")
	{
		vnpay := payment.Group("/vnpay")
		{
			vnpay.POST("", h.CreateVNPayPayment())
			vnpay.GET("/vnpay-callback", h.VnpayCallbackHandler())
		}
	}

}
