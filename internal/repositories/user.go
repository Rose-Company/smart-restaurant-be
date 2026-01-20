package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
	BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	baseRepo := NewBaseRepository[models.User](db)
	return &UserRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type OTPRepo struct {
	db *gorm.DB
	BaseRepository[models.OTP]
}

func NewOTPRepository(db *gorm.DB) *OTPRepo {
	baseRepo := NewBaseRepository[models.OTP](db)
	return &OTPRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type OTPAttemptRepo struct {
	db *gorm.DB
	BaseRepository[models.OTPAttempt]
}

func NewOTPAttemptRepository(db *gorm.DB) *OTPAttemptRepo {
	baseRepo := NewBaseRepository[models.OTPAttempt](db)
	return &OTPAttemptRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type RoleRepo struct {
	db *gorm.DB
	BaseRepository[models.Role]
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	baseRepo := NewBaseRepository[models.Role](db)
	return &RoleRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
