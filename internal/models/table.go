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
	TableNumber string     `json:"tableNumber" gorm:"column:table_number"`
	Capacity    int        `json:"capacity" gorm:"column:capacity"`
	Location    string     `json:"location" gorm:"column:location"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
}

type TableOrderData struct {
	ActiveOrders int     `json:"activeOrders"`
	TotalBill    float64 `json:"totalBill"`
}

type TableWithOrderData struct {
	ID          int             `json:"id"`
	TableNumber string          `json:"tableNumber"`
	Capacity    int             `json:"capacity"`
	Location    string          `json:"location"`
	Status      string          `json:"status"`
	OrderData   *TableOrderData `json:"orderData,omitempty"`
}

func (Table) TableName() string {
	return common.POSTGRES_TABLE_NAME_TABLES
}
