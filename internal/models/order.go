package models

import (
	"app-noti/common"
	"time"
)

// ============================================
// ORDER MODELS
// ============================================

type Order struct {
	ID                  int         `json:"id" gorm:"primaryKey;autoIncrement"`
	RestaurantID        *int        `json:"restaurant_id" gorm:"column:restaurant_id"`
	TableID             int         `json:"table_id" gorm:"column:table_id;not null"`
	Table               *Table      `json:"table,omitempty" gorm:"foreignKey:TableID"`
	OrderNumber         string      `json:"order_number" gorm:"column:order_number;uniqueIndex;not null"`
	CustomerUserID      *string     `json:"customer_user_id" gorm:"column:customer_user_id;type:uuid"`
	CustomerUser        *User       `json:"customer,omitempty" gorm:"foreignKey:CustomerUserID"`
	CustomerName        *string     `json:"customer_name" gorm:"column:customer_name"`
	CustomerPhone       *string     `json:"customer_phone" gorm:"column:customer_phone"`
	CustomerEmail       *string     `json:"customer_email" gorm:"column:customer_email"`
	WaiterID            *string     `json:"waiter_id" gorm:"column:waiter_id;type:uuid"`
	Waiter              *User       `json:"waiter,omitempty" gorm:"foreignKey:WaiterID"`
	KitchenStaffID      *string     `json:"kitchen_staff_id" gorm:"column:kitchen_staff_id;type:uuid"`
	KitchenStaff        *User       `json:"kitchen_staff,omitempty" gorm:"foreignKey:KitchenStaffID"`
	Status              string      `json:"status" gorm:"column:status;default:'pending'"`
	Subtotal            float64     `json:"subtotal" gorm:"column:subtotal;type:decimal(10,2);default:0"`
	Tax                 float64     `json:"tax" gorm:"column:tax;type:decimal(10,2);default:0"`
	Discount            float64     `json:"discount" gorm:"column:discount;type:decimal(10,2);default:0"`
	Total               float64     `json:"total" gorm:"column:total;type:decimal(10,2);default:0;not null"`
	Notes               *string     `json:"notes" gorm:"column:notes;type:text"`
	SpecialInstructions *string     `json:"special_instructions" gorm:"column:special_instructions;type:text"`
	Priority            string      `json:"priority" gorm:"column:priority;default:'normal'"`
	Source              string      `json:"source" gorm:"column:source;default:'qr'"`
	EstimatedReadyTime  *time.Time  `json:"estimated_ready_time" gorm:"column:estimated_ready_time"`
	Meta                *string     `json:"meta" gorm:"column:meta;type:jsonb"`
	CreatedAt           time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	AcceptedAt          *time.Time  `json:"accepted_at" gorm:"column:accepted_at"`
	PreparingAt         *time.Time  `json:"preparing_at" gorm:"column:preparing_at"`
	ReadyAt             *time.Time  `json:"ready_at" gorm:"column:ready_at"`
	ServedAt            *time.Time  `json:"served_at" gorm:"column:served_at"`
	CompletedAt         *time.Time  `json:"completed_at" gorm:"column:completed_at"`
	CancelledAt         *time.Time  `json:"cancelled_at" gorm:"column:cancelled_at"`
	CancelledBy         *string     `json:"cancelled_by" gorm:"column:cancelled_by"`
	CancelReason        *string     `json:"cancel_reason" gorm:"column:cancel_reason;type:text"`
	Items               []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return common.POSTGRES_TABLE_NAME_ORDERS
}

type OrderItem struct {
	ID                  int             `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID             int             `json:"order_id" gorm:"column:order_id;not null"`
	Order               *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	MenuItemID          *int            `json:"menu_item_id" gorm:"column:menu_item_id"`
	MenuItem            *MenuItem       `json:"menu_item,omitempty" gorm:"foreignKey:MenuItemID"`
	ItemName            string          `json:"item_name" gorm:"column:item_name;not null"`
	ItemDescription     *string         `json:"item_description" gorm:"column:item_description;type:text"`
	Quantity            int             `json:"quantity" gorm:"column:quantity;default:1"`
	UnitPrice           float64         `json:"unit_price" gorm:"column:unit_price;type:decimal(10,2);not null"`
	Subtotal            float64         `json:"subtotal" gorm:"column:subtotal;type:decimal(10,2);not null"`
	ModifiersTotal      float64         `json:"modifiers_total" gorm:"column:modifiers_total;type:decimal(10,2);default:0"`
	TaxAmount           float64         `json:"tax_amount" gorm:"column:tax_amount;type:decimal(10,2);default:0"`
	DiscountAmount      float64         `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2);default:0"`
	Status              string          `json:"status" gorm:"column:status;default:'pending'"`
	SpecialInstructions *string         `json:"special_instructions" gorm:"column:special_instructions;type:text"`
	Meta                *string         `json:"meta" gorm:"column:meta;type:jsonb"`
	CreatedAt           time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Modifiers           []OrderModifier `json:"modifiers,omitempty" gorm:"foreignKey:OrderItemID"`
}

func (OrderItem) TableName() string {
	return common.POSTGRES_TABLE_NAME_ORDER_ITEMS
}

type OrderModifier struct {
	ID                 int       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderItemID        int       `json:"order_item_id" gorm:"column:order_item_id;not null"`
	ModifierGroupID    int       `json:"modifier_group_id" gorm:"column:modifier_group_id;not null"`
	ModifierGroupName  string    `json:"modifier_group_name" gorm:"column:modifier_group_name;not null"`
	ModifierOptionID   int       `json:"modifier_option_id" gorm:"column:modifier_option_id;not null"`
	ModifierOptionName string    `json:"modifier_option_name" gorm:"column:modifier_option_name;not null"`
	PriceAdjustment    float64   `json:"price_adjustment" gorm:"column:price_adjustment;type:decimal(10,2);default:0"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (OrderModifier) TableName() string {
	return "order_item_modifiers"
}

type OrderTimeline struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID       int       `json:"order_id" gorm:"column:order_id;not null"`
	Status        string    `json:"status" gorm:"column:status;not null"`
	Timestamp     time.Time `json:"timestamp" gorm:"column:timestamp;autoCreateTime"`
	UpdatedBy     string    `json:"updated_by" gorm:"column:updated_by;not null"`
	UpdatedByID   *string   `json:"updated_by_id" gorm:"column:updated_by_id;type:uuid"`
	UpdatedByName *string   `json:"updated_by_name" gorm:"column:updated_by_name"`
	Note          *string   `json:"note" gorm:"column:note;type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (OrderTimeline) TableName() string {
	return "order_timeline"
}

type KitchenAlert struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID        int        `json:"order_id" gorm:"column:order_id;not null"`
	Order          *Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	AlertType      string     `json:"alert_type" gorm:"column:alert_type;not null"`
	Message        string     `json:"message" gorm:"column:message;not null;type:text"`
	Priority       string     `json:"priority" gorm:"column:priority;default:'normal'"`
	SentTo         string     `json:"sent_to" gorm:"column:sent_to;not null"`
	WaiterID       *string    `json:"waiter_id" gorm:"column:waiter_id;type:uuid"`
	Waiter         *User      `json:"waiter,omitempty" gorm:"foreignKey:WaiterID"`
	Status         string     `json:"status" gorm:"column:status;default:'sent'"`
	SentAt         time.Time  `json:"sent_at" gorm:"column:sent_at;autoCreateTime"`
	AcknowledgedAt *time.Time `json:"acknowledged_at" gorm:"column:acknowledged_at"`
	ResolvedAt     *time.Time `json:"resolved_at" gorm:"column:resolved_at"`
}

func (KitchenAlert) TableName() string {
	return "kitchen_alerts"
}

// ============================================
// REQUEST/RESPONSE MODELS
// ============================================

type CreateOrderRequest struct {
	TableID    int                      `json:"table_id" binding:"required"`
	CustomerID *string                  `json:"customer_id"`
	Items      []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	Notes      *string                  `json:"notes"`
}

type CreateOrderItemRequest struct {
	MenuItemID          int                    `json:"menu_item_id" binding:"required"`
	Quantity            int                    `json:"quantity" binding:"required,min=1"`
	SpecialInstructions *string                `json:"special_instructions"`
	Modifiers           []OrderModifierRequest `json:"modifiers"`
}

type OrderModifierRequest struct {
	ModifierGroupID  int `json:"modifier_group_id" binding:"required"`
	ModifierOptionID int `json:"modifier_option_id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status    string  `json:"status" binding:"required"`
	UpdatedBy string  `json:"updated_by"`
	Reason    *string `json:"reason"`
}

type UpdateOrderItemStatusRequest struct {
	OrderIDs []int   `json:"order_ids" binding:"required,min=1"`
	Status   string  `json:"status" binding:"required"`
	Reason   *string `json:"reason"`
}

type UpdateOrderMultipleItemsStatusRequest struct {
	Items []struct {
		MenuItemID int    `json:"menu_item_id" binding:"required"`
		Status     string `json:"status" binding:"required"`
	} `json:"items" binding:"required,min=1"`
	Reason *string `json:"reason"`
}

type UpdateOrderMultipleItemsStatusResponse struct {
	TotalUpdated int                             `json:"total_updated"`
	UpdatedItems []OrderItemStatusUpdateResponse `json:"updated_items"`
}

type UpdateOrderRequest struct {
	Note     *string                `json:"note"`
	Metadata map[string]interface{} `json:"metadata"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type CreateAlertRequest struct {
	Message   string `json:"message" binding:"required"`
	AlertType string `json:"alert_type" binding:"required"`
	Priority  string `json:"priority"`
}

type ListOrdersRequest struct {
	BaseRequestParamsUri
	Role     *string `form:"role"`
	Status   *string `form:"status"`
	DateFrom *string `form:"date_from"`
	DateTo   *string `form:"date_to"`
	TableID  *int    `form:"table_id"`
	Search   *string `form:"search"`
	Category *string `form:"category"`
}

// Response models
type OrderResponse struct {
	ID                  int                     `json:"id"`
	OrderNumber         string                  `json:"order_number"`
	TableID             int                     `json:"table_id"`
	TableName           string                  `json:"table_name"`
	CustomerID          *string                 `json:"customer_id,omitempty"`
	CustomerName        string                  `json:"customer_name"`
	CustomerPhone       *string                 `json:"customer_phone,omitempty"`
	Status              string                  `json:"status"`
	TotalAmount         float64                 `json:"total_amount"`
	TaxAmount           float64                 `json:"tax_amount"`
	DiscountAmount      float64                 `json:"discount_amount"`
	FinalAmount         float64                 `json:"final_amount"`
	Notes               *string                 `json:"notes,omitempty"`
	SpecialInstructions *string                 `json:"special_instructions,omitempty"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
	EstimatedReadyTime  *time.Time              `json:"estimated_ready_time,omitempty"`
	WaiterID            *string                 `json:"waiter_id,omitempty"`
	WaiterName          *string                 `json:"waiter_name,omitempty"`
	Items               []OrderItemResponse     `json:"items"`
	Timeline            []OrderTimelineResponse `json:"timeline,omitempty"`
	Bill                *BillSummary            `json:"bill,omitempty"`
}

type OrderListItemResponse struct {
	ID                 int                 `json:"id"`
	OrderNumber        string              `json:"order_number"`
	TableID            int                 `json:"table_id"`
	TableName          string              `json:"table_name"`
	CustomerID         *string             `json:"customer_id,omitempty"`
	CustomerName       string              `json:"customer_name"`
	Status             string              `json:"status"`
	TotalAmount        float64             `json:"total_amount"`
	ItemsCount         int                 `json:"items_count"`
	Items              []OrderItemResponse `json:"items"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	EstimatedReadyTime *time.Time          `json:"estimated_ready_time,omitempty"`
	WaiterID           *string             `json:"waiter_id,omitempty"`
	WaiterName         *string             `json:"waiter_name,omitempty"`
}

type OrderItemResponse struct {
	ID                  int                     `json:"id"`
	MenuItemID          *int                    `json:"menu_item_id"`
	MenuItemName        string                  `json:"menu_item_name"`
	MenuItemImage       *string                 `json:"menu_item_image,omitempty"`
	Quantity            int                     `json:"quantity"`
	UnitPrice           float64                 `json:"unit_price"`
	Subtotal            float64                 `json:"subtotal"`
	SpecialInstructions *string                 `json:"special_instructions,omitempty"`
	Status              string                  `json:"status"`
	Modifiers           []OrderModifierResponse `json:"modifiers"`
}

type OrderModifierResponse struct {
	ModifierGroupID    int     `json:"modifier_group_id"`
	ModifierGroupName  string  `json:"modifier_group_name"`
	ModifierOptionID   int     `json:"modifier_option_id"`
	ModifierOptionName string  `json:"modifier_option_name"`
	Price              float64 `json:"price"`
}

type OrderTimelineResponse struct {
	ID            int       `json:"id"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	UpdatedBy     string    `json:"updated_by"`
	UpdatedByName *string   `json:"updated_by_name,omitempty"`
	Note          *string   `json:"note,omitempty"`
}

type BillSummary struct {
	ID            int        `json:"id"`
	Status        string     `json:"status"`
	PaymentMethod *string    `json:"payment_method"`
	PaidAt        *time.Time `json:"paid_at"`
}

type OrderStatusUpdateResponse struct {
	ID                 int        `json:"id"`
	OrderNumber        string     `json:"order_number"`
	Status             string     `json:"status"`
	PreviousStatus     string     `json:"previous_status"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          string     `json:"updated_by"`
	UpdatedByName      *string    `json:"updated_by_name,omitempty"`
	EstimatedReadyTime *time.Time `json:"estimated_ready_time,omitempty"`
}

type OrderItemStatusUpdateResponse struct {
	ItemID         int       `json:"item_id"`
	OrderID        int       `json:"order_id"`
	MenuItemName   string    `json:"menu_item_name"`
	Status         string    `json:"status"`
	PreviousStatus string    `json:"previous_status"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	UpdatedByName  *string   `json:"updated_by_name,omitempty"`
}

type OrderItemStatusUpdateBatchResponse struct {
	TotalUpdated int                             `json:"total_updated"`
	UpdatedItems []OrderItemStatusUpdateResponse `json:"updated_items"`
}

type AlertResponse struct {
	AlertID     int       `json:"alert_id"`
	OrderID     int       `json:"order_id"`
	OrderNumber string    `json:"order_number"`
	Message     string    `json:"message"`
	SentTo      string    `json:"sent_to"`
	WaiterID    *string   `json:"waiter_id,omitempty"`
	WaiterName  *string   `json:"waiter_name,omitempty"`
	SentAt      time.Time `json:"sent_at"`
	Status      string    `json:"status"`
}

type CancelOrderResponse struct {
	ID           int       `json:"id"`
	OrderNumber  string    `json:"order_number"`
	Status       string    `json:"status"`
	CancelledAt  time.Time `json:"cancelled_at"`
	CancelledBy  string    `json:"cancelled_by"`
	CancelReason string    `json:"cancel_reason"`
	RefundStatus string    `json:"refund_status"`
	RefundAmount float64   `json:"refund_amount"`
}

type PaginatedOrdersResponse struct {
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Items    []OrderListItemResponse `json:"items"`
}

// TASK-010: Kitchen Items Summary grouped by category
type OrderItemsSummaryResponse struct {
	Category string                      `json:"category"`
	Items    []OrderItemSummaryWithTable `json:"items"`
}

type OrderItemSummaryWithTable struct {
	ItemID      int                  `json:"item_id"`
	MenuItemID  int                  `json:"menu_item_id"`
	ItemName    string               `json:"item_name"`
	TotalQty    int                  `json:"total_qty"`
	OrderTables []OrderItemTableInfo `json:"order_tables"`
}

type OrderItemTableInfo struct {
	OrderID     int    `json:"order_id"`
	TableID     int    `json:"table_id"`
	TableNumber string `json:"table_number"`
	Quantity    int    `json:"quantity"`
	Status      string `json:"status"`
}
