package models

import (
	"app-noti/common"
	"time"
)

type ListTablesRequest struct {
	BaseRequestParamsUri
	Search *string `form:"search"`
	Status *string `form:"status"`
	Zone   *string `form:"zone"`
}

type Table struct {
	ID          int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TableNumber string     `json:"table_number" gorm:"column:table_number"`
	Capacity    int        `json:"capacity" gorm:"column:capacity"`
	Location    string     `json:"location" gorm:"column:location"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type TableOrderData struct {
	ActiveOrders int     `json:"active_orders"`
	TotalBill    float64 `json:"total_bill"`
}

type TableWithOrderData struct {
	ID          int             `json:"id"`
	TableNumber string          `json:"table_number"`
	Capacity    int             `json:"capacity"`
	Location    string          `json:"location"`
	Status      string          `json:"status"`
	OrderData   *TableOrderData `json:"order_data,omitempty"`
}

type CreateTableRequest struct {
	TableNumber string `json:"table_number" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	Location    string `json:"location" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=active occupied inactive"`
}

type UpdateTableRequest struct {
	TableNumber *string `json:"table_number,omitempty"`
	Capacity    *int    `json:"capacity,omitempty" binding:"omitempty,min=1"`
	Location    *string `json:"location,omitempty"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=active occupied inactive"`
}

type UpdateTableStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active occupied inactive"`
}

func (Table) TableName() string {
	return common.POSTGRES_TABLE_NAME_TABLES
}
