package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

// UserPreferencesRepository interface
type UserPreferencesRepository interface {
	BaseRepository[models.UserPreferences]
}

type userPreferencesRepository struct {
	*baseRepository[models.UserPreferences]
}

// NewUserPreferencesRepository creates a new user preferences repository
func NewUserPreferencesRepository(db *gorm.DB) UserPreferencesRepository {
	var prefs models.UserPreferences
	return &userPreferencesRepository{
		baseRepository: &baseRepository[models.UserPreferences]{
			model: &prefs,
			db:    db,
		},
	}
}
