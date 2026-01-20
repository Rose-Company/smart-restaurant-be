package models

import (
	"app-noti/common"
	"time"
)

type User struct {
	ID            string           `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email         string           `json:"email" gorm:"uniqueIndex"`
	FirstName     *string          `json:"first_name,omitempty" gorm:"column:first_name"`
	LastName      *string          `json:"last_name,omitempty" gorm:"column:last_name"`
	Password      string           `json:"-" gorm:"password"`
	RoleID        string           `json:"role_id" gorm:"column:role"`
	Role          *Role            `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Status        string           `json:"status,omitempty" gorm:"status"`
	IsActive      bool             `json:"is_active" gorm:"is_active"`
	PhoneNumber   string           `json:"phone_number,omitempty" gorm:"phone_number;uniqueIndex"`
	Provider      string           `json:"provider,omitempty" gorm:"column:provider"`
	DateCreated   time.Time        `json:"date_created" gorm:"column:date_created;autoCreateTime"`
	AvatarURL     *string          `json:"avatar_url,omitempty" gorm:"-"`
	DateOfBirth   *string          `json:"date_of_birth,omitempty" gorm:"-"`
	Gender        *string          `json:"gender,omitempty" gorm:"-"`
	StreetAddress *string          `json:"street_address,omitempty" gorm:"-"`
	City          *string          `json:"city,omitempty" gorm:"-"`
	State         *string          `json:"state,omitempty" gorm:"-"`
	PostalCode    *string          `json:"postal_code,omitempty" gorm:"-"`
	Country       *string          `json:"country,omitempty" gorm:"-"`
	Preferences   *UserPreferences `json:"preferences,omitempty" gorm:"foreignKey:CustomerID;references:ID"`
	LoyaltyPoints int              `json:"loyalty_points,omitempty" gorm:"-"`
	LoyaltyTier   string           `json:"loyalty_tier,omitempty" gorm:"-"`
	TotalOrders   int              `json:"total_orders,omitempty" gorm:"-"`
	TotalSpent    float64          `json:"total_spent,omitempty" gorm:"-"`
	EmailVerified bool             `json:"email_verified,omitempty" gorm:"-"`
	PhoneVerified bool             `json:"phone_verified,omitempty" gorm:"-"`
	LastLoginAt   *time.Time       `json:"last_login_at,omitempty" gorm:"-"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return common.POSTGRES_TABLE_NAME_USERS
}

type Role struct {
	ID   string `json:"id" gorm:"type:uuid;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}

func (Role) TableName() string {
	return common.POSTGRES_TABLE_NAME_ROLES
}

type OTP struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Target      string    `gorm:"size:255;not null" json:"target"`
	Type        string    `gorm:"size:50;not null" json:"type"`
	OTPCode     string    `gorm:"size:6;not null" json:"otp_code"`
	ExpiredAt   time.Time `gorm:"not null" json:"expired_at"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	VerifyToken string    `gorm:"size:255" json:"verify_token"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OTP) TableName() string {
	return common.POSTGRES_TABLE_NAME_OTPS
}

type OTPAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OTPID     uint      `gorm:"not null" json:"otp_id"`
	Value     string    `gorm:"size:6;not null" json:"value"`
	IsSuccess bool      `gorm:"default:false" json:"is_success"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (OTPAttempt) TableName() string {
	return common.POSTGRES_TABLE_NAME_OTP_ATTEMPTS
}

// Request models
type SignupRequest struct {
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required"`
	VerifyToken string `json:"verify_token" binding:"required"`
}

type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	IdToken  *string `json:"id_token,omitempty"`
	Role     *string `json:"role,omitempty"` // Optional: specify role for OAuth users (default: end_user)
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	VerifyToken string `json:"verify_token" binding:"required"`
}

type OTPValidateRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type RequestResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Google OAuth user
type GoogleUser struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Exp           string `json:"exp"`
	AtHash        string `json:"at_hash"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Typ           string `json:"typ"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// UserPreferences model
type UserPreferences struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	CustomerID          string    `json:"customer_id" gorm:"column:customer_id;uniqueIndex"`
	DietaryRestrictions string    `json:"dietary_restrictions" gorm:"column:dietary_restrictions;type:text"` // JSON array as string
	FavoriteCuisines    string    `json:"favorite_cuisine" gorm:"column:favorite_cuisines;type:text"`        // JSON array as string
	NotificationEnabled bool      `json:"notification_enabled" gorm:"column:notification_enabled;default:true"`
	EmailNotifications  bool      `json:"email_notifications" gorm:"column:email_notifications;default:true"`
	SMSNotifications    bool      `json:"sms_notifications" gorm:"column:sms_notifications;default:false"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserPreferences) TableName() string {
	return "customer_preferences"
}

// Request/Response models for customer profile

type UpdateProfileRequest struct {
	Name        *string             `json:"name,omitempty"`
	Phone       *string             `json:"phone,omitempty"`
	DateOfBirth *string             `json:"date_of_birth,omitempty"`
	Gender      *string             `json:"gender,omitempty"`
	Address     *AddressRequest     `json:"address,omitempty"`
	Preferences *PreferencesRequest `json:"preferences,omitempty"`
}

type AddressRequest struct {
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    *string `json:"country,omitempty"`
}

type PreferencesRequest struct {
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	FavoriteCuisines    []string `json:"favorite_cuisine,omitempty"`
	NotificationEnabled *bool    `json:"notification_enabled,omitempty"`
	EmailNotifications  *bool    `json:"email_notifications,omitempty"`
	SMSNotifications    *bool    `json:"sms_notifications,omitempty"`
}

type ProfileResponse struct {
	ID            int                  `json:"id"`
	Email         string               `json:"email"`
	Name          *string              `json:"name,omitempty"`
	Phone         *string              `json:"phone,omitempty"`
	AvatarURL     *string              `json:"avatar_url,omitempty"`
	DateOfBirth   *string              `json:"date_of_birth,omitempty"`
	Gender        *string              `json:"gender,omitempty"`
	Address       *AddressResponse     `json:"address,omitempty"`
	Preferences   *PreferencesResponse `json:"preferences,omitempty"`
	Loyalty       *LoyaltyResponse     `json:"loyalty,omitempty"`
	Statistics    *StatisticsResponse  `json:"statistics,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	EmailVerified bool                 `json:"email_verified"`
	PhoneVerified bool                 `json:"phone_verified"`
}

type AddressResponse struct {
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    *string `json:"country,omitempty"`
}

type PreferencesResponse struct {
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	FavoriteCuisines    []string `json:"favorite_cuisine,omitempty"`
	NotificationEnabled bool     `json:"notification_enabled"`
	EmailNotifications  bool     `json:"email_notifications"`
	SMSNotifications    bool     `json:"sms_notifications"`
}

type LoyaltyResponse struct {
	Points       int    `json:"points"`
	Tier         string `json:"tier"`
	NextTier     string `json:"next_tier"`
	PointsToNext int    `json:"points_to_next_tier"`
}

type StatisticsResponse struct {
	TotalOrders   int      `json:"total_orders"`
	TotalSpent    float64  `json:"total_spent"`
	FavoriteItems []string `json:"favorite_items"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordResponse struct {
	ChangedAt time.Time `json:"changed_at"`
}

type UploadAvatarResponse struct {
	AvatarURL    string    `json:"avatar_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type GetMeResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
}
