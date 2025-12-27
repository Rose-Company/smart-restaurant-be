package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type RestaurantRepo struct {
	db *gorm.DB
	BaseRepository[models.Restaurant]
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepo {
	baseRepo := NewBaseRepository[models.Restaurant](db)
	return &RestaurantRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
