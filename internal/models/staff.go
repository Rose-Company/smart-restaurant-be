package models

import (
	"time"
)

// StaffProfile - Extended staff information
type StaffProfile struct {
	ID                    int        `json:"id" gorm:"primaryKey"`
	UserID                string     `json:"user_id" gorm:"column:user_id;uniqueIndex"`
	EmployeeID            *string    `json:"employee_id" gorm:"column:employee_id;uniqueIndex"`
	Department            *string    `json:"department" gorm:"column:department"`
	Position              *string    `json:"position" gorm:"column:position"`
	HireDate              *time.Time `json:"hire_date" gorm:"column:hire_date"`
	EmploymentStatus      string     `json:"employment_status" gorm:"column:employment_status;default:'active'"`
	ShiftType             *string    `json:"shift_type" gorm:"column:shift_type"`
	WeeklyHours           *int       `json:"weekly_hours" gorm:"column:weekly_hours"`
	EmergencyContactName  *string    `json:"emergency_contact_name" gorm:"column:emergency_contact_name"`
	EmergencyContactPhone *string    `json:"emergency_contact_phone" gorm:"column:emergency_contact_phone"`
	EmergencyContactRel   *string    `json:"emergency_contact_relationship" gorm:"column:emergency_contact_relationship"`
	TotalOrdersServed     int        `json:"total_orders_served" gorm:"column:total_orders_served;default:0"`
	TotalOrdersPrepared   int        `json:"total_orders_prepared" gorm:"column:total_orders_prepared;default:0"`
	AverageOrderTime      *int       `json:"average_order_time" gorm:"column:average_order_time"`
	AverageRating         *float64   `json:"average_rating" gorm:"column:average_rating"`
	Notes                 *string    `json:"notes" gorm:"column:notes"`
	CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (StaffProfile) TableName() string {
	return "staff_profiles"
}

// WaiterTableAssignment - Table assignments for waiters
type WaiterTableAssignment struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	WaiterID   string    `json:"waiter_id" gorm:"column:waiter_id"`
	TableID    int       `json:"table_id" gorm:"column:table_id"`
	AssignedAt time.Time `json:"assigned_at" gorm:"column:assigned_at;autoCreateTime"`
	AssignedBy *string   `json:"assigned_by" gorm:"column:assigned_by"`
	IsActive   bool      `json:"is_active" gorm:"column:is_active;default:true"`
	Table      *Table    `json:"table" gorm:"foreignKey:TableID"`
}

func (WaiterTableAssignment) TableName() string {
	return "waiter_table_assignments"
}

// StaffInvitation - Email invitations for new staff
type StaffInvitation struct {
	ID         int        `json:"id" gorm:"primaryKey"`
	Email      string     `json:"email" gorm:"column:email"`
	RoleID     string     `json:"role_id" gorm:"column:role_id"`
	InvitedBy  string     `json:"invited_by" gorm:"column:invited_by"`
	Token      string     `json:"token" gorm:"column:token;uniqueIndex"`
	Status     string     `json:"status" gorm:"column:status;default:'pending'"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"column:expires_at"`
	AcceptedAt *time.Time `json:"accepted_at" gorm:"column:accepted_at"`
	AcceptedBy *string    `json:"accepted_by" gorm:"column:accepted_by"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (StaffInvitation) TableName() string {
	return "staff_invitations"
}

// Request/Response models for staff management

// ListStaffRequest - Query parameters for listing staff
type ListStaffRequest struct {
	Role     *string `form:"role"`      // admin|waiter|kitchen_staff
	Page     int     `form:"page"`      // default: 1
	PageSize int     `form:"page_size"` // default: 10
	Search   *string `form:"search"`    // search term (name, email)
	Status   *string `form:"status"`    // active|inactive
}

// CreateStaffRequest - Create new staff account
type CreateStaffRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	Phone          string `json:"phone" binding:"required"`
	Role           string `json:"role" binding:"required"` // admin|waiter|kitchen_staff
	AssignedTables []int  `json:"assigned_tables,omitempty"`
}

// UpdateStaffRequest - Update staff account
type UpdateStaffRequest struct {
	Name           *string `json:"name,omitempty"`
	Email          *string `json:"email,omitempty"`
	Role           *string `json:"role,omitempty"`
	AssignedTables []int   `json:"assigned_tables,omitempty"`
	Status         *string `json:"status,omitempty"` // active|inactive
	Phone          *string `json:"phone,omitempty"`
}

// SendInviteRequest - Send staff invitation
type SendInviteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"` // admin|waiter|kitchen_staff
}

// AssignTablesRequest - Assign tables to waiter
type AssignTablesRequest struct {
	TableIDs []int `json:"tableIds" binding:"required"`
}

// StaffResponse - Common staff response
type StaffResponse struct {
	ID                    int                   `json:"id"`
	Name                  string                `json:"name"`
	Email                 string                `json:"email"`
	Phone                 *string               `json:"phone,omitempty"`
	Role                  string                `json:"role"`
	Status                string                `json:"status"`
	AvatarURL             *string               `json:"avatar_url,omitempty"`
	AssignedTables        []int                 `json:"assigned_tables,omitempty"`
	AssignedTablesNames   []string              `json:"assigned_tables_names,omitempty"`
	CreatedAt             time.Time             `json:"created_at"`
	LastLogin             *time.Time            `json:"last_login,omitempty"`
	Statistics            *StaffStatistics      `json:"statistics,omitempty"`
	InviteSent            bool                  `json:"invite_sent,omitempty"`
	AssignedTablesDetails []TableDetailResponse `json:"assigned_tables_details,omitempty"`
	Permissions           []string              `json:"permissions,omitempty"`
}

// StaffStatistics - Role-specific statistics
type StaffStatistics struct {
	TotalOrdersServed   *int     `json:"total_orders_served,omitempty"`
	AverageRating       *float64 `json:"average_rating,omitempty"`
	OrdersToday         *int     `json:"orders_today,omitempty"`
	ActiveOrders        *int     `json:"active_orders,omitempty"`
	TotalOrdersPrepared *int     `json:"total_orders_prepared,omitempty"`
	AveragePrepTime     *int     `json:"average_prep_time,omitempty"`
	PendingOrders       *int     `json:"pending_orders,omitempty"`
}

// TableDetailResponse - Detailed table information for staff
type TableDetailResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}

// StaffListResponse - Paginated staff list
type StaffListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Items    []StaffResponse `json:"items"`
}

// UpdateStaffProfileRequest - Update own staff profile
type UpdateStaffProfileRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Phone  *string `json:"phone,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

// ChangeStaffPasswordRequest - Change staff password
type ChangeStaffPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
