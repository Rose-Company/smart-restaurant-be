package server

import (
	"app-noti/pkg"
	"app-noti/pkg/rest_service"
	"app-noti/services/logger"
	"context"

	"github.com/go-redsync/redsync/v4"
)

type ServerContext interface {
	GetService(prefix string) interface{}
	GetLogger() logger.Loggers
	GetContext() context.Context
	SetUser(value interface{})
	GetUser() interface{}
	GetLoggerWithPrefix(prefix string) logger.Loggers
	GetRedisRedsync(prefix string) redsync.Redsync
	InitAuthorizationData()
	GetAuthConfig() *AuthorizationConfig
	SetTelegramService(service rest_service.RestInterface)
	GetTelegramService() rest_service.RestInterface
	GetAwsSes() *pkg.AWSSesClient
	SetAwsSes(service *pkg.AWSSesClient)
}
