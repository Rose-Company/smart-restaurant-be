package services

import (
	"app-noti/common"
	"app-noti/internal/repositories"
	l "app-noti/pkg/logger"
	"app-noti/server"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	logger    *zap.Logger
	tableRepo *repositories.TableRepo
}

func NewService(sc server.ServerContext) *Service {
	db := sc.GetService(common.PREFIX_MAIN_POSTGRES).(*gorm.DB)

	return &Service{
		logger:    l.New(),
		tableRepo: repositories.NewTableRepository(db),
	}
}
