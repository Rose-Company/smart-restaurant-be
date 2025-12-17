package rest_api_service

import (
	"app-noti/common"
	"app-noti/internal/handlers"
	"app-noti/middleware"
	"app-noti/server"
	"os"

	"github.com/gin-contrib/requestid"

	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RestHandler(sc server.ServerContext) func() *gin.Engine {
	return func() *gin.Engine {
		mode, ok := os.LookupEnv(common.ENV_GIN_DEBUG)
		if !ok {
			mode = "debug"
		}
		router := gin.New()
		gin.SetMode(mode)
		router.Use(requestid.New())
		router.Use(middleware.Logger(sc), middleware.Recovery(sc))
		router.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		}))

		sc.InitAuthorizationData()

		health := router.Group("/health")
		{
			health.GET("/status", handlers.Check(sc))
		}

		// Handler
		handler := handlers.NewHandler(sc)
		handler.RegisterRouter(router)

		return router
	}
}
