package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	BaseRepository[models.Payment]
}

type paymentRepository struct {
	*baseRepository[models.Payment]
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	var payment models.Payment
	return &paymentRepository{
		baseRepository: &baseRepository[models.Payment]{
			model: &payment,
			db:    db,
		},
	}
}

// DiscountCodeRepository defines the interface for discount code data access
type DiscountCodeRepository interface {
	BaseRepository[models.DiscountCode]
}

type discountCodeRepository struct {
	*baseRepository[models.DiscountCode]
}

// NewDiscountCodeRepository creates a new discount code repository
func NewDiscountCodeRepository(db *gorm.DB) *discountCodeRepository {
	var discountCode models.DiscountCode
	return &discountCodeRepository{
		baseRepository: &baseRepository[models.DiscountCode]{
			model: &discountCode,
			db:    db,
		},
	}
}

// DiscountUsageRepository defines the interface for discount usage data access
type DiscountUsageRepository interface {
	BaseRepository[models.DiscountUsage]
}

type discountUsageRepository struct {
	*baseRepository[models.DiscountUsage]
}

// NewDiscountUsageRepository creates a new discount usage repository
func NewDiscountUsageRepository(db *gorm.DB) *discountUsageRepository {
	var discountUsage models.DiscountUsage
	return &discountUsageRepository{
		baseRepository: &baseRepository[models.DiscountUsage]{
			model: &discountUsage,
			db:    db,
		},
	}
}
