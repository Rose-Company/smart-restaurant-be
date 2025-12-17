package services

import (
	l "app-noti/pkg/logger"
	"app-noti/server"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(sc server.ServerContext) *Service {

	return &Service{
		logger: l.New(),
	}
}
