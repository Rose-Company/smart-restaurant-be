package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	BaseRepository[models.Order]
}

type orderRepository struct {
	baseRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		baseRepository: baseRepository[models.Order]{
			db:    db,
			model: &models.Order{},
		},
	}
}

type OrderItemRepository interface {
	BaseRepository[models.OrderItem]
}

type orderItemRepository struct {
	baseRepository[models.OrderItem]
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{
		baseRepository: baseRepository[models.OrderItem]{
			db:    db,
			model: &models.OrderItem{},
		},
	}
}

type OrderModifierRepository interface {
	BaseRepository[models.OrderModifier]
}

type orderModifierRepository struct {
	baseRepository[models.OrderModifier]
}

func NewOrderModifierRepository(db *gorm.DB) OrderModifierRepository {
	return &orderModifierRepository{
		baseRepository: baseRepository[models.OrderModifier]{
			db:    db,
			model: &models.OrderModifier{},
		},
	}
}

type OrderTimelineRepository interface {
	BaseRepository[models.OrderTimeline]
}

type orderTimelineRepository struct {
	baseRepository[models.OrderTimeline]
}

func NewOrderTimelineRepository(db *gorm.DB) OrderTimelineRepository {
	return &orderTimelineRepository{
		baseRepository: baseRepository[models.OrderTimeline]{
			db:    db,
			model: &models.OrderTimeline{},
		},
	}
}

type KitchenAlertRepository interface {
	BaseRepository[models.KitchenAlert]
}

type kitchenAlertRepository struct {
	baseRepository[models.KitchenAlert]
}

func NewKitchenAlertRepository(db *gorm.DB) KitchenAlertRepository {
	return &kitchenAlertRepository{
		baseRepository: baseRepository[models.KitchenAlert]{
			db:    db,
			model: &models.KitchenAlert{},
		},
	}
}
