package models

import (
	"app-noti/common"
	"time"
)

type MenuCategory struct {
	ID           int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RestaurantID int        `json:"restaurant_id" gorm:"column:restaurant_id"`
	Name         string     `json:"name" gorm:"column:name"`
	Description  *string    `json:"description,omitempty" gorm:"column:description"`
	DisplayOrder int        `json:"display_order" gorm:"column:display_order"`
	Status       string     `json:"status" gorm:"column:status"`
	CreatedAt    *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (MenuCategory) TableName() string {
	return common.POSTGRES_TABLE_NAME_MENU_CATEGORIES
}

type CreateMenuCategoryRequest struct {
	Name         string  `json:"name" binding:"required,max=50"`
	Description  *string `json:"description"`
	DisplayOrder int     `json:"display_order" binding:"min=0"`
	Status       string  `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateMenuCategoryRequest struct {
	Name         *string `json:"name" binding:"max=50"`
	Description  *string `json:"description"`
	DisplayOrder *int    `json:"display_order" binding:"min=0"`
	Status       *string `json:"status" binding:"oneof=active inactive"`
}

type ListMenuCategoryRequest struct {
	BaseRequestParamsUri
	Search *string `form:"search"`
	Status *string `form:"status"`
}

type MenuItem struct {
	ID                int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RestaurantID      int        `json:"restaurant_id" gorm:"column:restaurant_id"`
	CategoryID        int        `json:"category_id" gorm:"column:category_id"`
	Name              string     `json:"name" gorm:"column:name"`
	Description       *string    `json:"description,omitempty" gorm:"column:description"`
	Price             float64    `json:"price" gorm:"column:price"`
	PrepTimeMinutes   int        `json:"prep_time_minutes" gorm:"column:prep_time_minutes"`
	Status            string     `json:"status" gorm:"column:status"`
	IsChefRecommended bool       `json:"is_chef_recommended" gorm:"column:is_chef_recommended"`
	IsDeleted         bool       `json:"is_deleted" gorm:"column:is_deleted"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (MenuItem) TableName() string {
	return common.POSTGRES_TABLE_NAME_MENU_ITEMS
}

type CreateMenuItemRequest struct {
	CategoryID        int     `json:"category_id" binding:"required"`
	Name              string  `json:"name" binding:"required,max=80"`
	Description       *string `json:"description"`
	Price             float64 `json:"price" binding:"required,gt=0"`
	PrepTimeMinutes   int     `json:"prep_time_minutes" binding:"min=0,max=240"`
	Status            string  `json:"status" binding:"required,oneof=available unavailable sold_out"`
	IsChefRecommended bool    `json:"is_chef_recommended"`
}

type UpdateMenuItemRequest struct {
	CategoryID        *int     `json:"category_id"`
	Name              *string  `json:"name" binding:"max=80"`
	Description       *string  `json:"description"`
	Price             *float64 `json:"price" binding:"gt=0"`
	PrepTimeMinutes   *int     `json:"prep_time_minutes" binding:"min=0,max=240"`
	Status            *string  `json:"status" binding:"oneof=available unavailable sold_out"`
	IsChefRecommended *bool    `json:"is_chef_recommended"`
}

type ListMenuItemRequest struct {
	BaseRequestParamsUri
	Search     *string `form:"search"`
	Status     *string `form:"status"`
	CategoryID *int    `form:"category_id"`
}

type MenuItemPhoto struct {
	ID         int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MenuItemID int        `json:"menu_item_id" gorm:"column:menu_item_id"`
	Url        string     `json:"url" gorm:"column:url"`
	IsPrimary  bool       `json:"is_primary" gorm:"column:is_primary"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func (MenuItemPhoto) TableName() string {
	return common.POSTGRES_TABLE_NAME_MENU_ITEM_PHOTOS
}

type CreateMenuItemPhotoRequest struct {
	Url       string `json:"url" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}
