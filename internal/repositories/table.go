package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type TableRepo struct {
	db *gorm.DB
	BaseRepository[models.Table]
}

func NewTableRepository(db *gorm.DB) *TableRepo {
	baseRepo := NewBaseRepository[models.Table](db)
	return &TableRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func (r *TableRepo) GetDB() *gorm.DB {
	return r.db
}
