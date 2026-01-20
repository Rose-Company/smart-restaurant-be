package services

import (
	"app-noti/internal/models"
	"app-noti/internal/repositories"
	"context"

	"gorm.io/gorm"
)

func (s *Service) getReviewByItemId(ctx context.Context, itemID int) ([]models.ReviewItem, error) {
	filters := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("menu_item_id = ?", itemID)
		},
	}

	reviews, err := s.reviewItemRepo.List(ctx, models.QueryParams{}, filters...)
	if err != nil || len(reviews) == 0 {
		return []models.ReviewItem{}, nil
	}

	result := make([]models.ReviewItem, 0, len(reviews))
	for _, r := range reviews {
		if r == nil {
			continue
		}
		result = append(result, *r)
	}

	return result, nil
}
