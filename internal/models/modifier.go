package models

import (
	"app-noti/common"
	"time"
)

type ModifierGroup struct {
	ID            int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RestaurantID  int        `json:"restaurant_id" gorm:"column:restaurant_id"`
	Name          string     `json:"name" gorm:"column:name"`
	SelectionType string     `json:"selection_type" gorm:"column:selection_type"`
	IsRequired    bool       `json:"is_required" gorm:"column:is_required"`
	MinSelections int        `json:"min_selections" gorm:"column:min_selections"`
	MaxSelections int        `json:"max_selections" gorm:"column:max_selections"`
	DisplayOrder  int        `json:"display_order" gorm:"column:display_order"`
	Status        string     `json:"status" gorm:"column:status"`
	CreatedAt     *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (ModifierGroup) TableName() string {
	return common.POSTGRES_TABLE_NAME_MODIFIER_GROUPS
}

type CreateModifierGroupRequest struct {
	Name          string `json:"name" binding:"required,max=80"`
	SelectionType string `json:"selection_type" binding:"required,oneof=single multiple"`
	IsRequired    bool   `json:"is_required"`
	MinSelections int    `json:"min_selections" binding:"min=0"`
	MaxSelections int    `json:"max_selections" binding:"min=0"`
	DisplayOrder  int    `json:"display_order" binding:"min=0"`
	Status        string `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateModifierGroupRequest struct {
	Name          *string `json:"name" binding:"max=80"`
	SelectionType *string `json:"selection_type" binding:"oneof=single multiple"`
	IsRequired    *bool   `json:"is_required"`
	MinSelections *int    `json:"min_selections" binding:"min=0"`
	MaxSelections *int    `json:"max_selections" binding:"min=0"`
	DisplayOrder  *int    `json:"display_order" binding:"min=0"`
	Status        *string `json:"status" binding:"oneof=active inactive"`
}

type ListModifierGroupRequest struct {
	BaseRequestParamsUri
	Search        *string `form:"search"`
	Status        *string `form:"status"`
	SelectionType *string `form:"selection_type"`
}

type ModifierGroupIDParamsUri struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type ModifierOption struct {
	ID              int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	GroupID         int        `json:"group_id" gorm:"column:group_id"`
	Name            string     `json:"name" gorm:"column:name"`
	PriceAdjustment float64    `json:"price_adjustment" gorm:"column:price_adjustment"`
	Status          string     `json:"status" gorm:"column:status"`
	CreatedAt       *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func (ModifierOption) TableName() string {
	return common.POSTGRES_TABLE_NAME_MODIFIER_OPTIONS
}

type CreateModifierOptionRequest struct {
	Name            string  `json:"name" binding:"required,max=80"`
	PriceAdjustment float64 `json:"price_adjustment" binding:"min=0"`
	Status          string  `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateModifierOptionRequest struct {
	Name            *string  `json:"name" binding:"max=80"`
	PriceAdjustment *float64 `json:"price_adjustment" binding:"min=0"`
	Status          *string  `json:"status" binding:"oneof=active inactive"`
}

type ListModifierOptionRequest struct {
	BaseRequestParamsUri
	Search *string `form:"search"`
	Status *string `form:"status"`
}

type ModifierOptionIDParamsUri struct {
	GroupID int `uri:"group_id" binding:"required,min=1"`
	ID      int `uri:"id" binding:"required,min=1"`
}

type MenuItemModifierGroup struct {
	ID         int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MenuItemID int        `json:"menu_item_id" gorm:"column:menu_item_id"`
	GroupID    int        `json:"group_id" gorm:"column:group_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func (MenuItemModifierGroup) TableName() string {
	return common.POSTGRES_TABLE_NAME_MENU_ITEM_MODIFIER_GROUPS
}

type AssignModifierToMenuItemRequest struct {
	MenuItemID int `json:"menu_item_id" binding:"required"`
	GroupID    int `json:"group_id" binding:"required"`
}
