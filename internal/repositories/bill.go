package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

// BillRepository defines the interface for bill data access
type BillRepository interface {
	BaseRepository[models.Bill]
}

type billRepository struct {
	*baseRepository[models.Bill]
}

// NewBillRepository creates a new bill repository
func NewBillRepository(db *gorm.DB) *billRepository {
	var bill models.Bill
	return &billRepository{
		baseRepository: &baseRepository[models.Bill]{
			model: &bill,
			db:    db,
		},
	}
}

// BillItemRepository defines the interface for bill item data access
type BillItemRepository interface {
	BaseRepository[models.BillItem]
}

type billItemRepository struct {
	*baseRepository[models.BillItem]
}

// NewBillItemRepository creates a new bill item repository
func NewBillItemRepository(db *gorm.DB) *billItemRepository {
	var billItem models.BillItem
	return &billItemRepository{
		baseRepository: &baseRepository[models.BillItem]{
			model: &billItem,
			db:    db,
		},
	}
}
