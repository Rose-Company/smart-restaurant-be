package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type ReviewItemRepo struct {
	db *gorm.DB
	BaseRepository[models.ReviewItem]
}

func NewReviewItemRepository(db *gorm.DB) *ReviewItemRepo {
	baseRepo := NewBaseRepository[models.ReviewItem](db)
	return &ReviewItemRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
