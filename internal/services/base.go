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
	logger                    *zap.Logger
	tableRepo                 *repositories.TableRepo
	restaurantRepo            *repositories.RestaurantRepo
	menuCategoryRepo          *repositories.MenuCategoryRepo
	menuItemRepo              *repositories.MenuItemRepo
	menuItemPhotoRepo         *repositories.MenuItemPhotoRepo
	modifierGroupRepo         *repositories.ModifierGroupRepo
	modifierOptionRepo        *repositories.ModifierOptionRepo
	menuItemModifierGroupRepo *repositories.MenuItemModifierGroupRepo
}

func NewService(sc server.ServerContext) *Service {
	db := sc.GetService(common.PREFIX_MAIN_POSTGRES).(*gorm.DB)

	return &Service{
		logger:                    l.New(),
		tableRepo:                 repositories.NewTableRepository(db),
		restaurantRepo:            repositories.NewRestaurantRepository(db),
		menuCategoryRepo:          repositories.NewMenuCategoryRepository(db),
		menuItemRepo:              repositories.NewMenuItemRepository(db),
		menuItemPhotoRepo:         repositories.NewMenuItemPhotoRepository(db),
		modifierGroupRepo:         repositories.NewModifierGroupRepository(db),
		modifierOptionRepo:        repositories.NewModifierOptionRepository(db),
		menuItemModifierGroupRepo: repositories.NewMenuItemModifierGroupRepository(db),
	}
}
