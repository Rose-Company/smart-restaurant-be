package repositories

import (
	"app-noti/internal/models"
	"context"

	"gorm.io/gorm"
)

type TableRepo struct {
	db *gorm.DB
	BaseRepository[models.Table]
}

func (r *TableRepo) GetAll(ctx context.Context, params models.QueryParams) ([]*models.Table, error) {
	var tables []*models.Table

	if err := r.db.WithContext(ctx).Find(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
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
