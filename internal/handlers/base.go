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

	c.GET("/api/me", middleware.UserAuthentication(), h.GetMe) // Get current user: id, role_id, role_name

	admin := c.Group("/api/admin", middleware.OptionalUserAuthentication())
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
			menuItem.GET("/:id", h.GetMenuItemByIDDetailed())
			menuItem.POST("/:id/modifier-groups", h.AssignMenuItemModifierGroup())
			menuItem.DELETE("/:id/modifier-groups/:groupId", h.DeleteMenuItemModifierGroup())
		}
	}

	// Order Management APIs (TASK-001 to TASK-010)
	// Using optional authentication - guests can order via QR code, logged-in users get personalized experience
	orders := c.Group("/api/orders")
	{
		orders.POST("", middleware.OptionalUserAuthentication(), h.CreateOrder)                                      // TASK-001: Create order
		orders.GET("", middleware.OptionalUserAuthentication(), h.GetOrders)                                         // TASK-002: Get orders list
		orders.GET("/:id", middleware.OptionalUserAuthentication(), h.GetOrderByID)                                  // TASK-003: Get order details
		orders.PATCH("/:id/status", middleware.OptionalUserAuthentication(), h.UpdateOrderStatus)                    // TASK-004: Update order status
		orders.PATCH("/items/:itemId/status", middleware.OptionalUserAuthentication(), h.UpdateOrderItemStatus)      // TASK-005: Update item status (batch)
		orders.PATCH("/:id/items/status", middleware.OptionalUserAuthentication(), h.UpdateOrderMultipleItemsStatus) // TASK-005b: Update multiple items status (single order)
		orders.PATCH("/:id", middleware.OptionalUserAuthentication(), h.UpdateOrder)                                 // TASK-006: Add notes/metadata
		orders.POST("/:id/cancel", middleware.OptionalUserAuthentication(), h.CancelOrder)                           // TASK-007: Cancel order
		orders.POST("/:id/alert", middleware.OptionalUserAuthentication(), h.SendKitchenAlert)                       // TASK-008: Send kitchen alert
		orders.POST("/:id/review", middleware.OptionalUserAuthentication(), h.CreateOrderReview)                     // TASK-009: Submit review
		orders.GET("/summary/category", middleware.OptionalUserAuthentication(), h.GetOrderItemsSummaryByCategory)   // TASK-010: Get items summary by category
	}

	// Bill Management APIs (TASK-010 to TASK-012)
	// Bill is created on TABLE level, not order level
	// One bill per table can include multiple orders and items
	bills := c.Group("/api/bills")
	{
		bills.POST("", middleware.OptionalUserAuthentication(), h.CreateBill())      // TASK-010: Create bill from table (includes all orders & items on table)
		bills.GET("/:id", middleware.OptionalUserAuthentication(), h.GetBill())      // TASK-011: Get bill details (table info + all orders + all items)
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

	// VN-PAY Callback
	c.GET("/api/vnpay/callback", h.HandleVNPayCallback())

	customer := c.Group("/api/customer")
	{
		customer.GET("/profile", middleware.UserAuthentication(), h.GetProfile)         // TASK-016: Get customer profile
		customer.PUT("/profile", middleware.UserAuthentication(), h.UpdateProfile)      // TASK-017: Update customer profile
		customer.POST("/avatar", middleware.UserAuthentication(), h.UploadAvatar)       // TASK-018: Upload avatar
		customer.PATCH("/password", middleware.UserAuthentication(), h.ChangePassword)  // TASK-019: Change password
		customer.GET("/reviews", middleware.UserAuthentication(), h.GetCustomerReviews) // TASK-020: Get customer reviews
	}

	// Staff Management Endpoints (TASK-021 to TASK-027)
	adminStaff := c.Group("/api/admin/staff")
	{
		adminStaff.GET("", middleware.OptionalUserAuthentication(), h.ListStaff)                                // TASK-021: List staff with filters
		adminStaff.POST("", middleware.OptionalUserAuthentication(), h.CreateStaff)                             // TASK-022: Create staff account
		adminStaff.GET("/:id", middleware.OptionalUserAuthentication(), h.GetStaffByID)                         // TASK-023: Get staff details
		adminStaff.PUT("/:id", middleware.OptionalUserAuthentication(), h.UpdateStaff)                          // TASK-024: Update staff account
		adminStaff.DELETE("/:id", middleware.OptionalUserAuthentication(), h.DeleteStaff)                       // TASK-025: Delete staff account
		adminStaff.POST("/:id/send-invite", middleware.OptionalUserAuthentication(), h.SendStaffInvite)         // TASK-026: Send staff invitation
		adminStaff.PATCH("/:id/assign-tables", middleware.OptionalUserAuthentication(), h.AssignTablesToWaiter) // TASK-027: Assign tables to waiter
	}

	// Staff Profile Endpoints (TASK-028 to TASK-030)
	staff := c.Group("/api/staff")
	{
		staff.GET("/profile", middleware.OptionalUserAuthentication(), h.GetStaffProfile)             // TASK-028: Get staff profile
		staff.PUT("/profile", middleware.OptionalUserAuthentication(), h.UpdateStaffProfile)          // TASK-029: Update own profile
		staff.PATCH("/password", middleware.OptionalUserAuthentication(), h.ChangeStaffPassword)      // TASK-030: Change password
		staff.GET("/tables", middleware.OptionalUserAuthentication(), h.GetTablesForStaff())          // Staff view: Get tables with orders
		staff.GET("/tables/:id", middleware.OptionalUserAuthentication(), h.GetTableDetailForStaff()) // Staff view: Get table detail with items
	}

}
