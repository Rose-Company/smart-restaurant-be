package repositories

import (
	"app-noti/internal/models"
	"context"

	"gorm.io/gorm"
)

type MenuCategoryRepo struct {
	db *gorm.DB
	BaseRepository[models.MenuCategory]
}

func NewMenuCategoryRepository(db *gorm.DB) *MenuCategoryRepo {
	baseRepo := NewBaseRepository[models.MenuCategory](db)
	return &MenuCategoryRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type MenuItemRepo struct {
	db *gorm.DB
	BaseRepository[models.MenuItem]
}

func (r *MenuItemRepo) FindByRestaurantID(
	ctx context.Context,
	restaurantID int,
) ([]*models.MenuItem, error) {

	var items []*models.MenuItem

	err := r.db.WithContext(ctx).
		Where("restaurant_id = ?", restaurantID).
		Find(&items).Error

	return items, err
}

func NewMenuItemRepository(db *gorm.DB) *MenuItemRepo {
	baseRepo := NewBaseRepository[models.MenuItem](db)
	return &MenuItemRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type MenuItemPhotoRepo struct {
	db *gorm.DB
	BaseRepository[models.MenuItemPhoto]
}

func NewMenuItemPhotoRepository(db *gorm.DB) *MenuItemPhotoRepo {
	baseRepo := NewBaseRepository[models.MenuItemPhoto](db)
	return &MenuItemPhotoRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func (r *MenuItemPhotoRepo) GetDB() *gorm.DB {
	return r.db
}
