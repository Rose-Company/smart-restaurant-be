package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// ============================================
// PAYMENT MODELS
// ============================================

// Payment represents a payment transaction
type Payment struct {
	ID             int      `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID      string   `json:"payment_id" gorm:"column:payment_id;uniqueIndex;not null"`
	BillID         int      `json:"bill_id" gorm:"column:bill_id;not null;index"`
	Bill           *Bill    `json:"bill,omitempty" gorm:"foreignKey:BillID"`
	OrderID        int      `json:"order_id" gorm:"column:order_id;not null;index"`
	Order          *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Amount         float64  `json:"amount" gorm:"column:amount;type:decimal(10,2);not null"`
	Method         string   `json:"method" gorm:"column:method;not null"`          // cash, card, ewallet, vnpay, stripe, other
	Status         string   `json:"status" gorm:"column:status;default:'pending'"` // pending, processing, succeeded, failed, cancelled, refunded
	ReceivedAmount *float64 `json:"received_amount" gorm:"column:received_amount;type:decimal(10,2)"`
	ChangeAmount   *float64 `json:"change_amount" gorm:"column:change_amount;type:decimal(10,2)"`
	// VN-PAY specific fields
	VNPayTransactionNo *string `json:"vnpay_transaction_no" gorm:"column:vnpay_transaction_no;index"`
	VNPayBankCode      *string `json:"vnpay_bank_code" gorm:"column:vnpay_bank_code"`
	VNPayCardType      *string `json:"vnpay_card_type" gorm:"column:vnpay_card_type"`
	VNPayResponseCode  *string `json:"vnpay_response_code" gorm:"column:vnpay_response_code"`
	// Stripe specific fields (kept for backward compatibility)
	StripePaymentIntentID *string           `json:"stripe_payment_intent_id" gorm:"column:stripe_payment_intent_id;index"`
	StripeChargeID        *string           `json:"stripe_charge_id" gorm:"column:stripe_charge_id"`
	StripePaymentMethodID *string           `json:"stripe_payment_method_id" gorm:"column:stripe_payment_method_id"`
	StripeCustomerID      *string           `json:"stripe_customer_id" gorm:"column:stripe_customer_id"`
	ReceiptURL            *string           `json:"receipt_url" gorm:"column:receipt_url"`
	ErrorCode             *string           `json:"error_code" gorm:"column:error_code"`
	ErrorMessage          *string           `json:"error_message" gorm:"column:error_message"`
	DeclineReason         *string           `json:"decline_reason" gorm:"column:decline_reason"`
	Metadata              datatypes.JSONMap `json:"metadata" gorm:"column:metadata;type:jsonb"`
	CreatedAt             time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	ProcessedAt           *time.Time        `json:"processed_at" gorm:"column:processed_at"`
	FailedAt              *time.Time        `json:"failed_at" gorm:"column:failed_at"`
	RefundedAt            *time.Time        `json:"refunded_at" gorm:"column:refunded_at"`
}

// TableName returns the table name for Payment
func (Payment) TableName() string {
	return "payments"
}

// ============================================
// DISCOUNT MODELS
// ============================================

// DiscountCode represents a promotional discount code
type DiscountCode struct {
	ID                   int           `json:"id" gorm:"primaryKey;autoIncrement"`
	Code                 string        `json:"code" gorm:"column:code;uniqueIndex;not null"`
	Description          *string       `json:"description" gorm:"column:description;type:text"`
	DiscountType         string        `json:"discount_type" gorm:"column:discount_type;not null"` // percentage, fixed_amount
	DiscountValue        float64       `json:"discount_value" gorm:"column:discount_value;type:decimal(10,2);not null"`
	MaxDiscountAmount    *float64      `json:"max_discount_amount" gorm:"column:max_discount_amount;type:decimal(10,2)"`
	MinOrderAmount       float64       `json:"min_order_amount" gorm:"column:min_order_amount;type:decimal(10,2);default:0"`
	UsageLimit           *int          `json:"usage_limit" gorm:"column:usage_limit"`
	UsageCount           int           `json:"usage_count" gorm:"column:usage_count;default:0"`
	PerCustomerLimit     *int          `json:"per_customer_limit" gorm:"column:per_customer_limit;default:1"`
	ValidFrom            time.Time     `json:"valid_from" gorm:"column:valid_from;not null"`
	ValidUntil           time.Time     `json:"valid_until" gorm:"column:valid_until;not null"`
	IsActive             bool          `json:"is_active" gorm:"column:is_active;default:true;index"`
	ApplicableCategories pq.Int64Array `json:"applicable_categories" gorm:"column:applicable_categories;type:integer[]"`
	ApplicableItems      pq.Int64Array `json:"applicable_items" gorm:"column:applicable_items;type:integer[]"`
	CreatedAt            time.Time     `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time     `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the table name for DiscountCode
func (DiscountCode) TableName() string {
	return "discount_codes"
}

// DiscountUsage represents a discount usage record
type DiscountUsage struct {
	ID             int           `json:"id" gorm:"primaryKey;autoIncrement"`
	DiscountID     int           `json:"discount_id" gorm:"column:discount_id;not null;index"`
	Discount       *DiscountCode `json:"discount,omitempty" gorm:"foreignKey:DiscountID"`
	CustomerID     *string       `json:"customer_id" gorm:"column:customer_id;type:uuid;index"`
	Customer       *User         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	OrderID        *int          `json:"order_id" gorm:"column:order_id;index"`
	Order          *Order        `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	BillID         *int          `json:"bill_id" gorm:"column:bill_id;index"`
	Bill           *Bill         `json:"bill,omitempty" gorm:"foreignKey:BillID"`
	DiscountAmount float64       `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2);not null"`
	UsedAt         time.Time     `json:"used_at" gorm:"column:used_at;autoCreateTime"`
}

// TableName returns the table name for DiscountUsage
func (DiscountUsage) TableName() string {
	return "discount_usage"
}

// ============================================
// REQUEST/RESPONSE MODELS
// ============================================

// ProcessPaymentRequest represents the request to process a payment
type ProcessPaymentRequest struct {
	BillID                int                    `json:"bill_id" binding:"required,min=1"`
	Amount                float64                `json:"amount" binding:"required,gt=0"`
	Method                string                 `json:"method" binding:"required,oneof=cash card ewallet vnpay stripe"`
	StripeToken           *string                `json:"stripe_token"`
	StripePaymentMethodID *string                `json:"stripe_payment_method_id"`
	ReceivedAmount        *float64               `json:"received_amount"`
	ChangeAmount          *float64               `json:"change_amount"`
	ClientIP              string                 `json:"client_ip"`
	Metadata              map[string]interface{} `json:"metadata"`
}

// PaymentResponse represents the response with payment details
type PaymentResponse struct {
	PaymentID             string     `json:"payment_id"`
	BillID                int        `json:"bill_id"`
	BillNumber            *string    `json:"bill_number"`
	Amount                float64    `json:"amount"`
	Method                string     `json:"method"`
	Status                string     `json:"status"`
	VNPayURL              *string    `json:"vnpay_url,omitempty"`
	StripePaymentIntentID *string    `json:"stripe_payment_intent_id,omitempty"`
	StripeChargeID        *string    `json:"stripe_charge_id,omitempty"`
	ReceiptURL            *string    `json:"receipt_url,omitempty"`
	ReceivedAmount        *float64   `json:"received_amount,omitempty"`
	ChangeAmount          *float64   `json:"change_amount,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	ProcessedAt           *time.Time `json:"processed_at"`
}

// PaymentStatusResponse represents the response for payment status check
type PaymentStatusResponse struct {
	PaymentID                 string     `json:"payment_id"`
	BillID                    int        `json:"bill_id"`
	BillNumber                *string    `json:"bill_number"`
	Status                    string     `json:"status"`
	Amount                    float64    `json:"amount"`
	Method                    string     `json:"method"`
	CreatedAt                 time.Time  `json:"created_at"`
	ProcessedAt               *time.Time `json:"processed_at"`
	FailedAt                  *time.Time `json:"failed_at,omitempty"`
	ReceiptURL                *string    `json:"receipt_url,omitempty"`
	ErrorMessage              *string    `json:"error_message,omitempty"`
	ErrorCode                 *string    `json:"error_code,omitempty"`
	StripePaymentIntentStatus *string    `json:"stripe_payment_intent_status,omitempty"`
	VNPayTransactionNo        *string    `json:"vnpay_transaction_no,omitempty"`
	VNPayBankCode             *string    `json:"vnpay_bank_code,omitempty"`
}

// PaymentCallbackResponse represents the response after VN-PAY callback
type PaymentCallbackResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// ValidateDiscountRequest represents the request to validate a discount code
type ValidateDiscountRequest struct {
	Code    string `json:"code" binding:"required"`
	OrderID int    `json:"order_id" binding:"required,min=1"`
}

// ValidateDiscountResponse represents the response for discount validation
type ValidateDiscountResponse struct {
	DiscountID        int       `json:"discount_id"`
	Code              string    `json:"code"`
	IsValid           bool      `json:"is_valid"`
	DiscountType      string    `json:"discount_type"`
	DiscountValue     float64   `json:"discount_value"`
	DiscountAmount    float64   `json:"discount_amount"`
	MinOrderAmount    float64   `json:"min_order_amount"`
	MaxDiscountAmount *float64  `json:"max_discount_amount"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidUntil        time.Time `json:"valid_until"`
	UsageLimit        *int      `json:"usage_limit"`
	UsageCount        int       `json:"usage_count"`
	Description       *string   `json:"description"`
}

// ValidateDiscountErrorResponse represents an error response for discount validation
type ValidateDiscountErrorResponse struct {
	Code               string     `json:"code"`
	IsValid            bool       `json:"is_valid"`
	ErrorReason        string     `json:"error_reason"`
	ValidUntil         *time.Time `json:"valid_until,omitempty"`
	MinOrderAmount     *float64   `json:"min_order_amount,omitempty"`
	CurrentOrderAmount *float64   `json:"current_order_amount,omitempty"`
	RequiredAdditional *float64   `json:"required_additional,omitempty"`
}
