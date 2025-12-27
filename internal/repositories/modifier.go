package repositories

import (
	"app-noti/internal/models"
	"context"
	"errors"

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

type MenuItemModifierGroupRepo struct {
	db *gorm.DB
	BaseRepository[models.MenuItemModifierGroup]
}

func (r *MenuItemModifierGroupRepo) FindByMenuItemIDAndGroupID(
	ctx context.Context,
	menuItemID int,
	groupID int,
) (*models.MenuItemModifierGroup, error) {
	var result models.MenuItemModifierGroup

	err := r.db.WithContext(ctx).
		Where("menu_item_id = ? AND group_id = ?", menuItemID, groupID).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *MenuItemModifierGroupRepo) DeleteByMenuItemIDAndGroupID(
	ctx context.Context,
	menuItemID int,
	groupID int,
) error {

	return r.db.WithContext(ctx).
		Where("menu_item_id = ? AND group_id = ?", menuItemID, groupID).
		Delete(&models.MenuItemModifierGroup{}).
		Error
}

func (r *ModifierOptionRepo) GetDB() *gorm.DB {
	return r.db
}

func NewMenuItemModifierGroupRepository(db *gorm.DB) *MenuItemModifierGroupRepo {
	baseRepo := NewBaseRepository[models.MenuItemModifierGroup](db)
	return &MenuItemModifierGroupRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
