package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type ModifierGroupRepo struct {
	db *gorm.DB
	BaseRepository[models.ModifierGroup]
}

func NewModifierGroupRepository(db *gorm.DB) *ModifierGroupRepo {
	baseRepo := NewBaseRepository[models.ModifierGroup](db)
	return &ModifierGroupRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type ModifierOptionRepo struct {
	db *gorm.DB
	BaseRepository[models.ModifierOption]
}

func NewModifierOptionRepository(db *gorm.DB) *ModifierOptionRepo {
	baseRepo := NewBaseRepository[models.ModifierOption](db)
	return &ModifierOptionRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func (r *ModifierOptionRepo) GetDB() *gorm.DB {
	return r.db
}

type MenuItemModifierGroupRepo struct {
	db *gorm.DB
	BaseRepository[models.MenuItemModifierGroup]
}

func NewMenuItemModifierGroupRepository(db *gorm.DB) *MenuItemModifierGroupRepo {
	baseRepo := NewBaseRepository[models.MenuItemModifierGroup](db)
	return &MenuItemModifierGroupRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
