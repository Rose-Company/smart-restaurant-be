package common

import "time"

const (
	SUCCESS_STATUS           = 200
	BAD_REQUEST_STATUS       = 400
	SERVER_ERROR_STATUS      = 500
	UNAUTHORIZED_STATUS      = 401
	PERMISSION_DENIED_STATUS = 403
	NOT_FOUND_STATUS         = 404
	INTERNAL_SERVER_ERR      = 500
)

const (
	PREFIX_MAIN_POSTGRES       = "MAIN_POSTGRES"
	PREFIX_YOUPASS_DO_STORAGE  = "YOUPASS_DO_STORAGE"
	PREFIX_CRONJOB_AUTO_CANCEL = "CRONJOB_AUTO_CANCEL"
)

const ( //must NOT edit this
	ENV_GIN_DEBUG  = "GIN_DEBUG"
	ENV_RABBIT_URI = "RABBIT"
)

const (
	ENVJWTSecretKey = "JWT__SECRET_KEY"
)

var (
	DATETIME_WITH_TIMEZONE = time.RFC3339
)

const (
	USER_JWT_KEY = "USER_JWT_PROFILE"
	UserId       = "user_id"
)

const (
	POSTGRES_TABLE_NAME_TABLES                    = "public.tables"
	POSTGRES_TABLE_NAME_ORDERS                    = "public.orders"
	POSTGRES_TABLE_NAME_ORDER_ITEMS               = "public.order_items"
	POSTGRES_TABLE_NAME_RESTAURANTS               = "public.restaurants"
	POSTGRES_TABLE_NAME_MENU_CATEGORIES           = "public.menu_categories"
	POSTGRES_TABLE_NAME_MENU_ITEMS                = "public.menu_items"
	POSTGRES_TABLE_NAME_MENU_ITEM_PHOTOS          = "public.menu_item_photos"
	POSTGRES_TABLE_NAME_MODIFIER_GROUPS           = "public.modifier_groups"
	POSTGRES_TABLE_NAME_MODIFIER_OPTIONS          = "public.modifier_options"
	POSTGRES_TABLE_NAME_MENU_ITEM_MODIFIER_GROUPS = "public.menu_item_modifier_groups"
	POSTGRES_TABLE_NAME_USERS                     = "public.users"
	POSTGRES_TABLE_NAME_ROLES                     = "public.roles"
	POSTGRES_TABLE_NAME_OTPS                      = "public.otps"
	POSTGRES_TABLE_NAME_OTP_ATTEMPTS              = "public.otp_attempts"
)

// User roles
const (
	ROLE_END_USER      = "end_user"
	ROLE_ADMIN         = "admin"
	ROLE_END_USER_UUID = "end_user_uuid" // Replace with actual UUID from database
	ROLE_ADMIN_UUID    = "admin_uuid"    // Replace with actual UUID from database
)

// User providers
const (
	USER_PROVIDER_LOCAL  = "local"
	USER_PROVIDER_GOOGLE = "google"
)

// OTP types
const (
	TypeResetPassword = "reset_password"
	TypeVerifyEmail   = "verify_email"
)

// User status
const (
	USER_STATUS_ACTIVE   = "active"
	USER_STATUS_INACTIVE = "inactive"
)
