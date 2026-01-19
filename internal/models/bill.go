package models

import (
	"time"
)

// ============================================
// BILL MODELS
// ============================================

// Bill represents a bill generated from an order
type Bill struct {
	ID              int         `json:"id" gorm:"primaryKey;autoIncrement"`
	BillNumber      string      `json:"bill_number" gorm:"column:bill_number;uniqueIndex;not null"`
	OrderID         int         `json:"order_id" gorm:"column:order_id;not null;index"`
	Order           *Order      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	RestaurantID    *int        `json:"restaurant_id" gorm:"column:restaurant_id;index"`
	TableID         *int        `json:"table_id" gorm:"column:table_id"`
	Table           *Table      `json:"table,omitempty" gorm:"foreignKey:TableID"`
	CustomerID      *string     `json:"customer_id" gorm:"column:customer_id;type:uuid;index"`
	Customer        *User       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Subtotal        float64     `json:"subtotal" gorm:"column:subtotal;type:decimal(10,2);default:0"`
	TaxAmount       float64     `json:"tax_amount" gorm:"column:tax_amount;type:decimal(10,2);default:0"`
	TaxRate         float64     `json:"tax_rate" gorm:"column:tax_rate;type:decimal(5,2);default:10"`
	DiscountAmount  float64     `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2);default:0"`
	DiscountCode    *string     `json:"discount_code" gorm:"column:discount_code"`
	ServiceCharge   float64     `json:"service_charge" gorm:"column:service_charge;type:decimal(10,2);default:0"`
	TotalAmount     float64     `json:"total_amount" gorm:"column:total_amount;type:decimal(10,2);not null"`
	Status          string      `json:"status" gorm:"column:status;default:'pending'"`  // pending, paid, cancelled, refunded
	BillType        string      `json:"type" gorm:"column:bill_type;default:'request'"` // request, generated, manual
	PaymentMethod   *string     `json:"payment_method" gorm:"column:payment_method"`    // cash, card, ewallet, stripe, other
	RequestedBy     *string     `json:"requested_by" gorm:"column:requested_by"`        // customer, waiter, system
	RequestedByID   *string     `json:"requested_by_id" gorm:"column:requested_by_id;type:uuid"`
	RequestedByUser *User       `json:"requested_by_user,omitempty" gorm:"foreignKey:RequestedByID"`
	CreatedAt       time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	PaidAt          *time.Time  `json:"paid_at" gorm:"column:paid_at"`
	BillItems       []*BillItem `json:"items,omitempty" gorm:"foreignKey:BillID"`
}

// TableName returns the table name for Bill
func (Bill) TableName() string {
	return "bills"
}

// BillItem represents a line item in a bill
type BillItem struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	BillID         int       `json:"bill_id" gorm:"column:bill_id;not null;index"`
	MenuItemID     *int      `json:"menu_item_id" gorm:"column:menu_item_id"`
	MenuItem       *MenuItem `json:"menu_item,omitempty" gorm:"foreignKey:MenuItemID"`
	MenuItemName   string    `json:"menu_item_name" gorm:"column:menu_item_name;not null"`
	Quantity       int       `json:"quantity" gorm:"column:quantity;not null"`
	UnitPrice      float64   `json:"unit_price" gorm:"column:unit_price;type:decimal(10,2);not null"`
	ModifiersTotal float64   `json:"modifiers_total" gorm:"column:modifiers_total;type:decimal(10,2);default:0"`
	Subtotal       float64   `json:"subtotal" gorm:"column:subtotal;type:decimal(10,2);not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

// TableName returns the table name for BillItem
func (BillItem) TableName() string {
	return "bill_items"
}

// ============================================
// REQUEST/RESPONSE MODELS
// ============================================

// CreateBillRequest represents the request to create a bill from a table
// Bill includes all orders and items on the table
type CreateBillRequest struct {
	TableID      int     `json:"table_id" binding:"required,min=1"`
	Type         string  `json:"type" binding:"required,oneof=generated request manual"`
	RequestedBy  string  `json:"requested_by" binding:"required,oneof=customer waiter system"`
	DiscountCode *string `json:"discount_code"`
}

// OrderItemForBill - Order item with order info for bill display
type OrderItemForBill struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ItemName  string  `json:"item_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Status    string  `json:"status"`
}

// OrderForBill - Order with items for bill display
type OrderForBill struct {
	ID          int                `json:"id"`
	OrderNumber string             `json:"order_number"`
	Status      string             `json:"status"`
	TotalAmount float64            `json:"total_amount"`
	Items       []OrderItemForBill `json:"items"`
}

// BillDetailResponse - Detailed bill response with table, all orders, and all items
type BillDetailResponse struct {
	ID             int            `json:"id"`
	BillNumber     string         `json:"bill_number"`
	TableID        *int           `json:"table_id"`
	TableNumber    *string        `json:"table_number"`
	RestaurantID   *int           `json:"restaurant_id"`
	CustomerID     *string        `json:"customer_id"`
	CustomerName   *string        `json:"customer_name"`
	CustomerPhone  *string        `json:"customer_phone"`
	Orders         []OrderForBill `json:"orders"`
	OrdersCount    int            `json:"orders_count"`
	ItemsCount     int            `json:"items_count"`
	Subtotal       float64        `json:"subtotal"`
	TaxAmount      float64        `json:"tax_amount"`
	TaxRate        *float64       `json:"tax_rate"`
	DiscountAmount float64        `json:"discount_amount"`
	DiscountCode   *string        `json:"discount_code"`
	ServiceCharge  float64        `json:"service_charge"`
	TotalAmount    float64        `json:"total_amount"`
	Status         string         `json:"status"`
	Type           string         `json:"type"`
	PaymentMethod  *string        `json:"payment_method"`
	RequestedBy    *string        `json:"requested_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
	PaidAt         *time.Time     `json:"paid_at"`
	Breakdown      *BillBreakdown `json:"breakdown,omitempty"`
}

// BillResponse represents the response with bill details
type BillResponse struct {
	ID             int                 `json:"id"`
	BillNumber     string              `json:"bill_number"`
	OrderID        int                 `json:"order_id"`
	OrderNumber    *string             `json:"order_number"`
	TableID        *int                `json:"table_id"`
	TableName      *string             `json:"table_name"`
	CustomerID     *string             `json:"customer_id"`
	CustomerName   *string             `json:"customer_name"`
	CustomerPhone  *string             `json:"customer_phone"`
	Subtotal       float64             `json:"subtotal"`
	TaxAmount      float64             `json:"tax_amount"`
	TaxRate        *float64            `json:"tax_rate"`
	DiscountAmount float64             `json:"discount_amount"`
	DiscountCode   *string             `json:"discount_code"`
	ServiceCharge  float64             `json:"service_charge"`
	TotalAmount    float64             `json:"total_amount"`
	Status         string              `json:"status"`
	Type           string              `json:"type"`
	PaymentMethod  *string             `json:"payment_method"`
	RequestedBy    *string             `json:"requested_by"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at"`
	PaidAt         *time.Time          `json:"paid_at"`
	Items          []*BillItemResponse `json:"items"`
	Breakdown      *BillBreakdown      `json:"breakdown,omitempty"`
}

// BillItemResponse represents a bill item in response
type BillItemResponse struct {
	ID             int     `json:"id"`
	MenuItemName   string  `json:"menu_item_name"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	ModifiersTotal float64 `json:"modifiers_total"`
	Subtotal       float64 `json:"subtotal"`
}

// BillBreakdown represents the detailed breakdown of bill amounts
type BillBreakdown struct {
	ItemsTotal             float64 `json:"items_total"`
	ModifiersTotal         float64 `json:"modifiers_total"`
	SubtotalBeforeDiscount float64 `json:"subtotal_before_discount"`
	Discount               float64 `json:"discount"`
	SubtotalAfterDiscount  float64 `json:"subtotal_after_discount"`
	Tax                    float64 `json:"tax"`
	ServiceCharge          float64 `json:"service_charge"`
	GrandTotal             float64 `json:"grand_total"`
}

// UpdateBillRequest represents the request to update a bill
type UpdateBillRequest struct {
	DiscountCode   *string  `json:"discount_code"`
	DiscountAmount *float64 `json:"discount_amount"`
	Status         *string  `json:"status" binding:"omitempty,oneof=pending paid cancelled refunded"`
	PaymentMethod  *string  `json:"payment_method" binding:"omitempty,oneof=cash card ewallet vnpay stripe other"`
}

// GetBillRequest represents query params for getting a bill
type GetBillRequest struct {
	Format  string `form:"format" binding:"omitempty,oneof=json pdf"`
	Include string `form:"include"`
}

// BillPDFResponse represents the response for PDF bill generation
type BillPDFResponse struct {
	PdfURL    string    `json:"pdf_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
