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
		admin.GET("/tables", h.GetTables())
		admin.GET("/tables/:id", h.GetTableByID())
		admin.POST("/tables", h.CreateTable())
		admin.PUT("/tables/:id", h.UpdateTable())
		admin.PATCH("/tables/:id/status", h.UpdateTableStatus())
	}
}
