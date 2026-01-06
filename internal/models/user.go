package models

import (
	"app-noti/common"
	"time"
)

type User struct {
	ID          string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	FirstName   *string   `json:"first_name" gorm:"column:first_name"`
	LastName    *string   `json:"last_name" gorm:"column:last_name"`
	Password    string    `json:"-" gorm:"password"`
	RoleID      string    `json:"role_id" gorm:"column:role"`
	Role        *Role     `json:"role" gorm:"foreignKey:RoleID"`
	Status      string    `json:"status" gorm:"status"`
	IsActive    bool      `json:"is_active" gorm:"is_active"`
	PhoneNumber string    `json:"phone_number" gorm:"phone_number"`
	Provider    string    `gorm:"nullable"`
	DateCreated time.Time `json:"date_created" gorm:"column:date_created;autoCreateTime"`
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
