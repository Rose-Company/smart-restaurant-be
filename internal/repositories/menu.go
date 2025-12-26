package repositories

import (
	"app-noti/internal/models"

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
