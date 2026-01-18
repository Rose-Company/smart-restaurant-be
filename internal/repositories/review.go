package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

// ReviewRepository interface
type ReviewRepository interface {
	BaseRepository[models.Review]
}

type reviewRepository struct {
	*baseRepository[models.Review]
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	var review models.Review
	return &reviewRepository{
		baseRepository: &baseRepository[models.Review]{
			model: &review,
			db:    db,
		},
	}
}

// ReviewItemRepository interface
type ReviewItemRepository interface {
	BaseRepository[models.ReviewItem]
}

type reviewItemRepository struct {
	*baseRepository[models.ReviewItem]
}

// NewReviewItemRepository creates a new review item repository
func NewReviewItemRepository(db *gorm.DB) ReviewItemRepository {
	var item models.ReviewItem
	return &reviewItemRepository{
		baseRepository: &baseRepository[models.ReviewItem]{
			model: &item,
			db:    db,
		},
	}
}

// ReviewPhotoRepository interface
type ReviewPhotoRepository interface {
	BaseRepository[models.ReviewPhoto]
}

type reviewPhotoRepository struct {
	*baseRepository[models.ReviewPhoto]
}

// NewReviewPhotoRepository creates a new review photo repository
func NewReviewPhotoRepository(db *gorm.DB) ReviewPhotoRepository {
	var photo models.ReviewPhoto
	return &reviewPhotoRepository{
		baseRepository: &baseRepository[models.ReviewPhoto]{
			model: &photo,
			db:    db,
		},
	}
}

// RestaurantResponseRepository interface
type RestaurantResponseRepository interface {
	BaseRepository[models.RestaurantResponse]
}

type restaurantResponseRepository struct {
	*baseRepository[models.RestaurantResponse]
}

// NewRestaurantResponseRepository creates a new restaurant response repository
func NewRestaurantResponseRepository(db *gorm.DB) RestaurantResponseRepository {
	var response models.RestaurantResponse
	return &restaurantResponseRepository{
		baseRepository: &baseRepository[models.RestaurantResponse]{
			model: &response,
			db:    db,
		},
	}
}
