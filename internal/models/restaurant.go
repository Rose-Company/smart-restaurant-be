package models

import (
	"app-noti/common"
	"time"
)

type Restaurant struct {
	ID          int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"column:name"`
	Description *string    `json:"description,omitempty" gorm:"column:description"`
	Address     *string    `json:"address,omitempty" gorm:"column:address"`
	Phone       *string    `json:"phone,omitempty" gorm:"column:phone"`
	Email       *string    `json:"email,omitempty" gorm:"column:email"`
	LogoUrl     *string    `json:"logo_url,omitempty" gorm:"column:logo_url"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (Restaurant) TableName() string {
	return common.POSTGRES_TABLE_NAME_RESTAURANTS
}

type CreateRestaurantRequest struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone" binding:"omitempty,max=20"`
	Email       *string `json:"email" binding:"omitempty,max=255,email"`
	LogoUrl     *string `json:"logo_url"`
	Status      string  `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateRestaurantRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone" binding:"omitempty,max=20"`
	Email       *string `json:"email" binding:"omitempty,max=255,email"`
	LogoUrl     *string `json:"logo_url"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type ListRestaurantRequest struct {
	BaseRequestParamsUri
	Search *string `form:"search"`
	Status *string `form:"status"`
}
