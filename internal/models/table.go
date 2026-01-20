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

// ListTablesForStaffRequest - Filter tables for staff view with order details
type ListTablesForStaffRequest struct {
	BaseRequestParamsUri
	Search        *string `form:"search"`
	Status        *string `form:"status"`
	IsReadyToBill *bool   `form:"is_ready_to_bill"`
	IsHelpNeeded  *bool   `form:"is_help_needed"`
	StaffID       string  `json:"staff_id"`
}

type TableParamsUri struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type Table struct {
	ID               int        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TableNumber      string     `json:"table_number" gorm:"column:table_number"`
	RestaurantId     int        `json:"restaurant_id"`
	Capacity         int        `json:"capacity" gorm:"column:capacity"`
	Location         string     `json:"location" gorm:"column:location"`
	Status           string     `json:"status" gorm:"column:status"`
	IsHelpNeeded     bool       `json:"is_help_needed" gorm:"column:is_help_needed"`
	IsReadyToBill    bool       `json:"is_ready_to_bill" gorm:"column:is_ready_to_bill"`
	QrToken          string     `json:"qr_token" gorm:"column:qr_token"`
	QrTokenCreatedAt *time.Time `json:"qr_token_created_at" gorm:"column:qr_token_created_at"`
	QrTokenExpiresAt *time.Time `json:"qr_token_expires_at" gorm:"column:qr_token_expires_at"`
	CreatedAt        *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type TableOrderData struct {
	ActiveOrders int     `json:"active_orders"`
	TotalBill    float64 `json:"total_bill"`
}

type TableWithOrderData struct {
	ID               int             `json:"id"`
	RestaurantId     int             `json:"restaurant_id"`
	TableNumber      string          `json:"table_number"`
	Capacity         int             `json:"capacity"`
	Location         string          `json:"location"`
	Status           string          `json:"status"`
	QrToken          string          `json:"qr_token" gorm:"column:qr_token"`
	QrTokenCreatedAt *time.Time      `json:"qr_token_created_at" gorm:"column:qr_token_created_at"`
	QrTokenExpiresAt *time.Time      `json:"qr_token_expires_at" gorm:"column:qr_token_expires_at"`
	OrderData        *TableOrderData `json:"order_data,omitempty"`
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

type UpdateTableFlagsRequest struct {
	IsReadyToBill *bool `json:"is_ready_to_bill,omitempty"`
	IsHelpNeeded  *bool `json:"is_help_needed,omitempty"`
}

type GenerateQrCodeRequest struct {
	TableNumber string `json:"table_number" binding:"required"`
}

type QrCodeData struct {
	TableID   int        `json:"table_id"`
	Token     string     `json:"token"`
	CreatedAt *time.Time `json:"create_at"`
	ExpiresAt *time.Time `json:"expire_at"`
}

type GenerateQrCodeResponse struct {
	Url string `json:"url"`
}

type QrCodeInfo struct {
	Token     string
	CreatedAt *time.Time
	ExpiresAt *time.Time
}

// TableOrderSummary - Summary of order for staff view
type TableOrderSummary struct {
	ID            int                `json:"id"`
	OrderNumber   string             `json:"order_number"`
	Status        string             `json:"status"`
	TotalAmount   float64            `json:"total_amount"`
	IsReadyToBill bool               `json:"is_ready_to_bill"`
	IsHelpNeeded  bool               `json:"is_help_needed"`
	ItemsCount    int                `json:"items_count"`
	Items         []OrderItemSummary `json:"items"`
	CreatedAt     time.Time          `json:"created_at"`
	CustomerName  string             `json:"customer_name,omitempty"`
}

// OrderItemSummary - Summary of order item
type OrderItemSummary struct {
	ID        int     `json:"id"`
	ItemName  string  `json:"item_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Status    string  `json:"status"`
}

// OrderItemDetailForStaff - Order item with order ID for staff view
type OrderItemDetailForStaff struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ItemName  string  `json:"item_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Status    string  `json:"status"`
}

// TableForStaffResponse - Table with orders for staff view
type TableForStaffResponse struct {
	ID                int                 `json:"id"`
	TableNumber       string              `json:"table_number"`
	Capacity          int                 `json:"capacity"`
	Location          string              `json:"location"`
	Status            string              `json:"status"`
	Orders            []TableOrderSummary `json:"orders"`
	ActiveOrdersCount int                 `json:"active_orders_count"`
	TotalBill         float64             `json:"total_bill"`
	IsHelpNeeded      bool                `json:"is_help_needed"`
	IsReadyToBill     bool                `json:"is_ready_to_bill"`
	CreatedAt         *time.Time          `json:"created_at,omitempty"`
	UpdatedAt         *time.Time          `json:"updated_at,omitempty"`
}

// TableDetailForStaffResponse - Table detail with all order items for staff
type TableDetailForStaffResponse struct {
	ID             int                       `json:"id"`
	TableNumber    string                    `json:"table_number"`
	Capacity       int                       `json:"capacity"`
	Location       string                    `json:"location"`
	Status         string                    `json:"status"`
	GuestCount     int                       `json:"guest_count"`
	TotalBill      float64                   `json:"total_bill"`
	OrderItems     []OrderItemDetailForStaff `json:"order_items"`
	AllOrdersCount int                       `json:"all_orders_count"`
	Bill           *Bill                     `json:"bill,omitempty"`
	CreatedAt      *time.Time                `json:"created_at,omitempty"`
}

func (Table) TableName() string {
	return common.POSTGRES_TABLE_NAME_TABLES
}
