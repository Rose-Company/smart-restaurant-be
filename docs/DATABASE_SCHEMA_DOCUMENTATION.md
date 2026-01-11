# 📊 SMART RESTAURANT DATABASE SCHEMA DOCUMENTATION

## 🗂️ OVERVIEW

This document provides a comprehensive overview of the complete database schema for the Smart Restaurant System, including all existing and newly designed tables to support the full feature set.

---

## 📋 TABLE OF CONTENTS

1. [Core Tables (Existing)](#1-core-tables-existing)
2. [Order Management System](#2-order-management-system)
3. [Bills & Payments System](#3-bills--payments-system)
4. [Customer Profile & Reviews](#4-customer-profile--reviews)
5. [Staff Management & Permissions](#5-staff-management--permissions)
6. [Analytics & Reports](#6-analytics--reports)
7. [Relationships Diagram](#7-relationships-diagram)
8. [Indexes Summary](#8-indexes-summary)
9. [Views Summary](#9-views-summary)
10. [Triggers & Functions](#10-triggers--functions)

---

## 1. CORE TABLES (Existing)

### 1.1 `restaurants`
**Purpose:** Restaurant information  
**Migration:** `002_create_restaurant.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Restaurant ID |
| name | VARCHAR(255) | NOT NULL | Restaurant name |
| description | TEXT | | Restaurant description |
| address | TEXT | | Physical address |
| phone | VARCHAR(20) | | Contact phone |
| email | VARCHAR(255) | | Contact email |
| logo_url | TEXT | | Logo image URL |
| status | VARCHAR(20) | DEFAULT 'active', CHECK IN ('active', 'inactive') | Restaurant status |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.2 `tables`
**Purpose:** Restaurant table management  
**Migration:** `001_create_tables.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Table ID |
| restaurant_id | INT | FK → restaurants.id | Restaurant reference |
| table_number | VARCHAR(50) | NOT NULL, UNIQUE | Table identifier |
| capacity | INT | CHECK (1-20) | Maximum seats |
| location | VARCHAR(100) | | Table location/area |
| description | TEXT | | Additional info |
| status | TEXT | DEFAULT 'active' | active/occupied/inactive |
| qr_token | TEXT | | QR code token |
| qr_token_created_at | TIMESTAMP | | QR token creation |
| qr_token_expires_at | TIMESTAMP | | QR token expiry |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.3 `roles`
**Purpose:** User roles for RBAC  
**Migration:** `005_create_authentication.sql` + `008_create_staff_management.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Role ID |
| name | VARCHAR(50) | NOT NULL, UNIQUE | Role name (end_user, admin, waiter, kitchen_staff) |
| display_name | VARCHAR(100) | | Human-readable name |
| description | TEXT | | Role description |
| permissions | JSONB | | Legacy permissions field |
| is_staff | BOOLEAN | DEFAULT FALSE | Is staff role |
| is_active | BOOLEAN | DEFAULT TRUE | Is active |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.4 `users`
**Purpose:** User accounts (customers + staff)  
**Migration:** `005_create_authentication.sql` + `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | User ID |
| email | VARCHAR(255) | NOT NULL, UNIQUE | Email address |
| password | VARCHAR(255) | | Hashed password |
| first_name | VARCHAR(100) | | First name |
| last_name | VARCHAR(100) | | Last name |
| full_name | VARCHAR(255) | | Full name |
| phone_number | VARCHAR(20) | | Phone number |
| role | UUID | NOT NULL, FK → roles.id | User role |
| status | VARCHAR(20) | DEFAULT 'active' | User status |
| is_active | BOOLEAN | DEFAULT TRUE | Active flag |
| provider | VARCHAR(50) | DEFAULT 'local' | OAuth provider |
| avatar_url | TEXT | | Profile picture URL |
| date_of_birth | DATE | | Birth date |
| gender | VARCHAR(20) | CHECK IN ('male', 'female', 'other', 'prefer_not_to_say') | Gender |
| street_address | VARCHAR(255) | | Street address |
| city | VARCHAR(100) | | City |
| state | VARCHAR(100) | | State/Province |
| postal_code | VARCHAR(20) | | Postal code |
| country | VARCHAR(100) | DEFAULT 'USA' | Country |
| email_verified | BOOLEAN | DEFAULT FALSE | Email verified |
| phone_verified | BOOLEAN | DEFAULT FALSE | Phone verified |
| last_login_at | TIMESTAMP | | Last login time |
| profile_completed | BOOLEAN | DEFAULT FALSE | Profile complete |
| date_created | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Legacy creation field |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.5 `otps`
**Purpose:** OTP codes for email/phone verification  
**Migration:** `005_create_authentication.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | OTP ID |
| target | VARCHAR(255) | NOT NULL | Email/Phone |
| type | VARCHAR(50) | NOT NULL | signup/reset_password |
| otp_code | VARCHAR(6) | NOT NULL | 6-digit code |
| expired_at | TIMESTAMP | NOT NULL | Expiry time |
| is_verified | BOOLEAN | DEFAULT FALSE | Verified flag |
| verify_token | VARCHAR(255) | | Verification token |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.6 `otp_attempts`
**Purpose:** Track OTP verification attempts  
**Migration:** `005_create_authentication.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Attempt ID |
| otp_id | INT | NOT NULL, FK → otps.id | OTP reference |
| value | VARCHAR(6) | NOT NULL | Attempted code |
| is_success | BOOLEAN | DEFAULT FALSE | Success flag |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Attempt timestamp |

---

### 1.7 `menu_categories`
**Purpose:** Menu categories  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Category ID |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| name | VARCHAR(50) | NOT NULL | Category name |
| description | TEXT | | Description |
| display_order | INT | DEFAULT 0 | Sort order |
| status | VARCHAR(20) | DEFAULT 'active', CHECK IN ('active', 'inactive') | Status |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Unique:** `(restaurant_id, name)`

---

### 1.8 `menu_items`
**Purpose:** Menu items  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Item ID |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| category_id | INT | NOT NULL, FK → menu_categories.id | Category reference |
| name | VARCHAR(80) | NOT NULL | Item name |
| description | TEXT | | Item description |
| price | DECIMAL(12,2) | NOT NULL, CHECK > 0 | Base price |
| prep_time_minutes | INT | DEFAULT 0, CHECK 0-240 | Preparation time |
| status | VARCHAR(20) | NOT NULL, CHECK IN ('available', 'unavailable', 'sold_out') | Availability |
| is_chef_recommended | BOOLEAN | DEFAULT FALSE | Chef's recommendation |
| is_deleted | BOOLEAN | DEFAULT FALSE | Soft delete flag |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.9 `menu_item_photos`
**Purpose:** Menu item images  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Photo ID |
| menu_item_id | INT | NOT NULL, FK → menu_items.id | Item reference |
| url | TEXT | NOT NULL | Image URL |
| is_primary | BOOLEAN | DEFAULT FALSE | Primary photo flag |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Upload timestamp |

---

### 1.10 `modifier_groups`
**Purpose:** Modifier groups (e.g., "Steak Temperature")  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Group ID |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| name | VARCHAR(80) | NOT NULL | Group name |
| selection_type | VARCHAR(20) | NOT NULL, CHECK IN ('single', 'multiple') | Selection mode |
| is_required | BOOLEAN | DEFAULT FALSE | Required flag |
| min_selections | INT | DEFAULT 0 | Minimum selections |
| max_selections | INT | DEFAULT 0 | Maximum selections |
| display_order | INT | DEFAULT 0 | Sort order |
| status | VARCHAR(20) | DEFAULT 'active', CHECK IN ('active', 'inactive') | Status |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

---

### 1.11 `modifier_options`
**Purpose:** Modifier options (e.g., "Medium Rare")  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Option ID |
| group_id | INT | NOT NULL, FK → modifier_groups.id | Group reference |
| name | VARCHAR(80) | NOT NULL | Option name |
| price_adjustment | DECIMAL(12,2) | DEFAULT 0, CHECK >= 0 | Price modifier |
| status | VARCHAR(20) | DEFAULT 'active', CHECK IN ('active', 'inactive') | Status |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

---

### 1.12 `menu_item_modifier_groups`
**Purpose:** Link menu items to modifier groups  
**Migration:** `003_create_menu_modifiers.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Link ID |
| menu_item_id | INT | NOT NULL, FK → menu_items.id | Item reference |
| group_id | INT | NOT NULL, FK → modifier_groups.id | Group reference |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Link timestamp |

**Unique:** `(menu_item_id, group_id)`

---

## 2. ORDER MANAGEMENT SYSTEM

**Migration:** `006_create_orders_system.sql`

### 2.1 `orders` (Enhanced)
**Purpose:** Customer orders  
**Migration:** `001_create_tables.sql` + `006_create_orders_system.sql` (enhancements)

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Order ID |
| restaurant_id | INT | FK → restaurants.id | Restaurant reference |
| table_id | INT | NOT NULL, FK → tables.id | Table reference |
| order_number | VARCHAR(50) | NOT NULL, UNIQUE | Order identifier (e.g., ORD-2026-001001) |
| customer_user_id | UUID | FK → users.id | Customer reference |
| customer_name | VARCHAR(255) | | Customer name (for guests) |
| customer_phone | VARCHAR(20) | | Customer phone |
| customer_email | VARCHAR(255) | | Customer email |
| waiter_id | UUID | FK → users.id | Assigned waiter |
| kitchen_staff_id | UUID | FK → users.id | Assigned kitchen staff |
| status | VARCHAR(20) | DEFAULT 'pending', CHECK IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'cancelled') | Order status |
| subtotal | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Items subtotal |
| tax | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Tax amount |
| discount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Discount amount |
| total | DECIMAL(10,2) | NOT NULL, DEFAULT 0.00, CHECK >= 0 | Final total |
| notes | TEXT | | Order notes |
| special_instructions | TEXT | | Special instructions |
| priority | VARCHAR(20) | DEFAULT 'normal', CHECK IN ('low', 'normal', 'high', 'urgent') | Order priority |
| source | VARCHAR(20) | DEFAULT 'qr', CHECK IN ('qr', 'waiter', 'admin', 'customer_app') | Order source |
| estimated_ready_time | TIMESTAMP | | Estimated completion |
| meta | JSONB | | Additional metadata |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Order placed time |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update time |
| accepted_at | TIMESTAMP | | Accepted timestamp |
| preparing_at | TIMESTAMP | | Preparing started |
| ready_at | TIMESTAMP | | Ready for pickup |
| served_at | TIMESTAMP | | Served to customer |
| completed_at | TIMESTAMP | | Order completed |
| cancelled_at | TIMESTAMP | | Cancellation time |
| cancelled_by | VARCHAR(50) | | Who cancelled (customer/waiter/admin) |
| cancel_reason | TEXT | | Cancellation reason |

**Indexes:**
- `idx_orders_status` (status)
- `idx_orders_table_id` (table_id)
- `idx_orders_customer_user_id` (customer_user_id)
- `idx_orders_waiter_id` (waiter_id)
- `idx_orders_created_at` (created_at DESC)
- `idx_orders_order_number` (order_number)

---

### 2.2 `order_items` (Enhanced)
**Purpose:** Order line items  
**Migration:** `001_create_tables.sql` + `006_create_orders_system.sql` (enhancements)

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Item ID |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| menu_item_id | INT | FK → menu_items.id | Menu item reference |
| item_name | VARCHAR(255) | NOT NULL | Item name (snapshot) |
| item_description | TEXT | | Item description |
| quantity | INT | DEFAULT 1, CHECK > 0 | Quantity ordered |
| unit_price | DECIMAL(10,2) | NOT NULL, CHECK >= 0 | Unit price |
| subtotal | DECIMAL(10,2) | NOT NULL, CHECK >= 0 | Line subtotal |
| modifiers_total | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Total modifier adjustments |
| tax_amount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Item tax |
| discount_amount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Item discount |
| final_price | DECIMAL(10,2) | GENERATED ALWAYS AS (subtotal + modifiers_total - discount_amount) STORED | Final price |
| status | VARCHAR(20) | DEFAULT 'pending', CHECK IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'rejected', 'cancelled') | Item status |
| special_instructions | TEXT | | Item-specific notes |
| meta | JSONB | | Additional metadata |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_order_items_order_id` (order_id)
- `idx_order_items_menu_item_id` (menu_item_id)
- `idx_order_items_status` (status)

---

### 2.3 `order_item_modifiers` ✨ NEW
**Purpose:** Track selected modifiers for order items  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Modifier ID |
| order_item_id | INT | NOT NULL, FK → order_items.id | Order item reference |
| modifier_group_id | INT | NOT NULL, FK → modifier_groups.id | Modifier group |
| modifier_group_name | VARCHAR(80) | NOT NULL | Group name (snapshot) |
| modifier_option_id | INT | NOT NULL, FK → modifier_options.id | Selected option |
| modifier_option_name | VARCHAR(80) | NOT NULL | Option name (snapshot) |
| price_adjustment | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Price change |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Selection timestamp |

**Indexes:**
- `idx_order_item_modifiers_order_item_id` (order_item_id)

**Purpose:** When a customer orders "Ribeye Steak" with "Medium" temperature and "Grilled Vegetables" side, this table stores those modifier selections.

---

### 2.4 `order_timeline` ✨ NEW
**Purpose:** Audit trail of order status changes  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Timeline ID |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| status | VARCHAR(20) | NOT NULL, CHECK IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'cancelled') | Status |
| timestamp | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Status change time |
| updated_by | VARCHAR(50) | NOT NULL | Role who updated (customer/waiter/kitchen/admin/system) |
| updated_by_id | UUID | FK → users.id | User ID |
| updated_by_name | VARCHAR(255) | | User name (snapshot) |
| note | TEXT | | Status change note |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation |

**Indexes:**
- `idx_order_timeline_order_id` (order_id)
- `idx_order_timeline_timestamp` (timestamp DESC)

**Purpose:** Provides customer order tracking timeline (e.g., "Order placed → Confirmed by waiter → Kitchen started preparing → Ready for pickup → Served").

---

### 2.5 `order_notes` ✨ NEW
**Purpose:** Additional notes/alerts for orders  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Note ID |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| note | TEXT | NOT NULL | Note content |
| note_type | VARCHAR(20) | DEFAULT 'general', CHECK IN ('general', 'allergy', 'special_request', 'kitchen_note', 'alert') | Note category |
| created_by | VARCHAR(50) | NOT NULL | Role (customer/waiter/kitchen/admin) |
| created_by_id | UUID | FK → users.id | User ID |
| created_by_name | VARCHAR(255) | | User name |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Note timestamp |

**Indexes:**
- `idx_order_notes_order_id` (order_id)
- `idx_order_notes_created_at` (created_at DESC)

---

### 2.6 `kitchen_alerts` ✨ NEW
**Purpose:** Alerts sent from kitchen to waiters  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Alert ID |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| alert_type | VARCHAR(50) | NOT NULL, CHECK IN ('order_ready', 'item_ready', 'delay_warning', 'special_request', 'urgent') | Alert type |
| message | TEXT | NOT NULL | Alert message |
| priority | VARCHAR(20) | DEFAULT 'normal', CHECK IN ('low', 'normal', 'high', 'urgent') | Alert priority |
| sent_to | VARCHAR(50) | NOT NULL | Target (waiter/kitchen/all_waiters/specific_waiter) |
| waiter_id | UUID | FK → users.id | Specific waiter |
| status | VARCHAR(20) | DEFAULT 'sent', CHECK IN ('sent', 'acknowledged', 'resolved', 'dismissed') | Alert status |
| sent_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Send time |
| acknowledged_at | TIMESTAMP | | Acknowledgment time |
| resolved_at | TIMESTAMP | | Resolution time |

**Indexes:**
- `idx_kitchen_alerts_order_id` (order_id)
- `idx_kitchen_alerts_status` (status)
- `idx_kitchen_alerts_waiter_id` (waiter_id)
- `idx_kitchen_alerts_sent_at` (sent_at DESC)

**Purpose:** Kitchen sends "Order for Table 5 is ready for pickup" alert to assigned waiter.

---

## 3. BILLS & PAYMENTS SYSTEM

**Migration:** `006_create_orders_system.sql`

### 3.1 `bills` ✨ NEW
**Purpose:** Bills generated from orders  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Bill ID |
| bill_number | VARCHAR(50) | UNIQUE NOT NULL | Bill identifier (e.g., BILL-2026-002001) |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| restaurant_id | INT | FK → restaurants.id | Restaurant reference |
| table_id | INT | FK → tables.id | Table reference |
| customer_id | UUID | FK → users.id | Customer reference |
| subtotal | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Items subtotal |
| tax_amount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Tax amount |
| tax_rate | DECIMAL(5,2) | DEFAULT 10.00, CHECK 0-100 | Tax percentage |
| discount_amount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Discount amount |
| discount_code | VARCHAR(50) | | Applied discount code |
| service_charge | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Service charge |
| total_amount | DECIMAL(10,2) | NOT NULL, CHECK >= 0 | Final total |
| status | VARCHAR(20) | DEFAULT 'pending', CHECK IN ('pending', 'paid', 'cancelled', 'refunded') | Bill status |
| bill_type | VARCHAR(20) | DEFAULT 'request', CHECK IN ('request', 'generated', 'manual') | Bill type |
| payment_method | VARCHAR(20) | CHECK IN ('cash', 'card', 'ewallet', 'stripe', 'other') | Payment method |
| requested_by | VARCHAR(50) | | Who requested (customer/waiter/system) |
| requested_by_id | UUID | FK → users.id | Requester user ID |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Bill creation |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update |
| paid_at | TIMESTAMP | | Payment timestamp |

**Indexes:**
- `idx_bills_bill_number` (bill_number)
- `idx_bills_order_id` (order_id)
- `idx_bills_status` (status)
- `idx_bills_created_at` (created_at DESC)
- `idx_bills_customer_id` (customer_id)

---

### 3.2 `payments` ✨ NEW
**Purpose:** Payment transactions  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Payment ID |
| payment_id | VARCHAR(100) | UNIQUE NOT NULL | Payment identifier |
| bill_id | INT | NOT NULL, FK → bills.id | Bill reference |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| amount | DECIMAL(10,2) | NOT NULL, CHECK > 0 | Payment amount |
| method | VARCHAR(20) | NOT NULL, CHECK IN ('cash', 'card', 'ewallet', 'stripe', 'other') | Payment method |
| status | VARCHAR(20) | DEFAULT 'pending', CHECK IN ('pending', 'processing', 'succeeded', 'failed', 'cancelled', 'refunded') | Payment status |
| received_amount | DECIMAL(10,2) | CHECK >= 0 | Cash received (cash payments) |
| change_amount | DECIMAL(10,2) | CHECK >= 0 | Change returned |
| stripe_payment_intent_id | VARCHAR(255) | | Stripe payment intent |
| stripe_charge_id | VARCHAR(255) | | Stripe charge ID |
| stripe_payment_method_id | VARCHAR(255) | | Stripe payment method |
| stripe_customer_id | VARCHAR(255) | | Stripe customer ID |
| receipt_url | TEXT | | Receipt URL (Stripe) |
| error_code | VARCHAR(50) | | Error code |
| error_message | TEXT | | Error message |
| decline_reason | VARCHAR(255) | | Card decline reason |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Payment creation |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update |
| processed_at | TIMESTAMP | | Processing timestamp |
| failed_at | TIMESTAMP | | Failure timestamp |
| refunded_at | TIMESTAMP | | Refund timestamp |
| metadata | JSONB | | Additional metadata |

**Indexes:**
- `idx_payments_payment_id` (payment_id)
- `idx_payments_bill_id` (bill_id)
- `idx_payments_order_id` (order_id)
- `idx_payments_status` (status)
- `idx_payments_stripe_payment_intent_id` (stripe_payment_intent_id)
- `idx_payments_created_at` (created_at DESC)

---

### 3.3 `discount_codes` ✨ NEW
**Purpose:** Promotional discount codes  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Discount ID |
| code | VARCHAR(50) | UNIQUE NOT NULL | Discount code (e.g., SAVE10) |
| description | TEXT | | Code description |
| discount_type | VARCHAR(20) | NOT NULL, CHECK IN ('percentage', 'fixed_amount') | Discount type |
| discount_value | DECIMAL(10,2) | NOT NULL, CHECK > 0 | Discount value |
| max_discount_amount | DECIMAL(10,2) | CHECK >= 0 | Maximum discount cap |
| min_order_amount | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Minimum order required |
| usage_limit | INT | CHECK > 0 | Total usage limit |
| usage_count | INT | DEFAULT 0, CHECK >= 0 | Current usage count |
| per_customer_limit | INT | DEFAULT 1, CHECK > 0 | Per-customer limit |
| valid_from | TIMESTAMP | NOT NULL | Start date |
| valid_until | TIMESTAMP | NOT NULL | End date |
| is_active | BOOLEAN | DEFAULT TRUE | Active flag |
| applicable_categories | INT[] | | Category IDs (array) |
| applicable_items | INT[] | | Menu item IDs (array) |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Check:** `valid_until > valid_from`

**Indexes:**
- `idx_discount_codes_code` (code)
- `idx_discount_codes_is_active` (is_active)
- `idx_discount_codes_valid_from` (valid_from)
- `idx_discount_codes_valid_until` (valid_until)

---

### 3.4 `discount_usage` ✨ NEW
**Purpose:** Track discount code usage  
**Migration:** `006_create_orders_system.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Usage ID |
| discount_id | INT | NOT NULL, FK → discount_codes.id | Discount code |
| customer_id | UUID | FK → users.id | Customer reference |
| order_id | INT | FK → orders.id | Order reference |
| bill_id | INT | FK → bills.id | Bill reference |
| discount_amount | DECIMAL(10,2) | NOT NULL, CHECK >= 0 | Discount applied |
| used_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Usage timestamp |

**Indexes:**
- `idx_discount_usage_discount_id` (discount_id)
- `idx_discount_usage_customer_id` (customer_id)
- `idx_discount_usage_used_at` (used_at DESC)

---

## 4. CUSTOMER PROFILE & REVIEWS

**Migration:** `007_create_customer_profile_reviews.sql`

### 4.1 `customer_preferences` ✨ NEW
**Purpose:** Customer preferences and settings  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Preference ID |
| customer_id | UUID | NOT NULL UNIQUE, FK → users.id | Customer reference |
| dietary_restrictions | TEXT[] | | Array of restrictions (no_nuts, vegetarian, etc.) |
| allergens | TEXT[] | | Specific allergens |
| favorite_cuisine | TEXT[] | | Favorite cuisines (Italian, Japanese, etc.) |
| spice_level | VARCHAR(20) | CHECK IN ('none', 'mild', 'medium', 'hot', 'extra_hot') | Spice preference |
| notification_enabled | BOOLEAN | DEFAULT TRUE | Notifications on/off |
| email_notifications | BOOLEAN | DEFAULT TRUE | Email notifications |
| sms_notifications | BOOLEAN | DEFAULT FALSE | SMS notifications |
| push_notifications | BOOLEAN | DEFAULT TRUE | Push notifications |
| marketing_emails | BOOLEAN | DEFAULT FALSE | Marketing emails |
| order_updates | BOOLEAN | DEFAULT TRUE | Order update emails |
| promotional_offers | BOOLEAN | DEFAULT TRUE | Promotional emails |
| language | VARCHAR(10) | DEFAULT 'en' | Preferred language |
| currency | VARCHAR(10) | DEFAULT 'USD' | Preferred currency |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_customer_preferences_customer_id` (customer_id)

---

### 4.2 `customer_loyalty` ✨ NEW
**Purpose:** Customer loyalty points and tier  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Loyalty ID |
| customer_id | UUID | NOT NULL UNIQUE, FK → users.id | Customer reference |
| points | INT | DEFAULT 0, CHECK >= 0 | Current points |
| lifetime_points | INT | DEFAULT 0, CHECK >= 0 | Total points earned |
| tier | VARCHAR(20) | DEFAULT 'bronze', CHECK IN ('bronze', 'silver', 'gold', 'platinum', 'diamond') | Loyalty tier |
| tier_expiry_date | DATE | | Tier expiry |
| available_rewards | INT | DEFAULT 0, CHECK >= 0 | Redeemable rewards |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Tier Thresholds:**
- Bronze: 0-99 points
- Silver: 100-499 points
- Gold: 500-999 points
- Platinum: 1000-1999 points
- Diamond: 2000+ points

**Indexes:**
- `idx_customer_loyalty_customer_id` (customer_id)
- `idx_customer_loyalty_tier` (tier)

---

### 4.3 `loyalty_transactions` ✨ NEW
**Purpose:** Loyalty points transaction history  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Transaction ID |
| customer_id | UUID | NOT NULL, FK → users.id | Customer reference |
| order_id | INT | FK → orders.id | Order reference |
| transaction_type | VARCHAR(20) | NOT NULL, CHECK IN ('earned', 'redeemed', 'expired', 'adjusted', 'bonus') | Transaction type |
| points | INT | NOT NULL | Points change (+ or -) |
| description | TEXT | | Transaction note |
| balance_before | INT | NOT NULL, CHECK >= 0 | Balance before |
| balance_after | INT | NOT NULL, CHECK >= 0 | Balance after |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Transaction timestamp |
| expires_at | DATE | | Points expiry (e.g., 1 year) |

**Indexes:**
- `idx_loyalty_transactions_customer_id` (customer_id)
- `idx_loyalty_transactions_order_id` (order_id)
- `idx_loyalty_transactions_created_at` (created_at DESC)

---

### 4.4 `customer_reviews` ✨ NEW
**Purpose:** Customer reviews and ratings  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Review ID |
| order_id | INT | NOT NULL, FK → orders.id | Order reference |
| customer_id | UUID | NOT NULL, FK → users.id | Customer reference |
| rating | INT | NOT NULL, CHECK 1-5 | Overall rating |
| comment | TEXT | | Review text |
| food_quality_rating | INT | CHECK 1-5 | Food quality |
| service_rating | INT | CHECK 1-5 | Service quality |
| ambiance_rating | INT | CHECK 1-5 | Ambiance rating |
| value_rating | INT | CHECK 1-5 | Value for money |
| restaurant_response | TEXT | | Restaurant reply |
| responded_by | UUID | FK → users.id | Responder user |
| responded_by_name | VARCHAR(255) | | Responder name |
| responded_at | TIMESTAMP | | Response timestamp |
| status | VARCHAR(20) | DEFAULT 'published', CHECK IN ('pending', 'published', 'hidden', 'flagged', 'deleted') | Review status |
| is_verified_purchase | BOOLEAN | DEFAULT TRUE | Verified purchase |
| helpful_count | INT | DEFAULT 0, CHECK >= 0 | Helpful votes |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Review timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Unique:** `(order_id, customer_id)` - One review per order per customer

**Indexes:**
- `idx_customer_reviews_order_id` (order_id)
- `idx_customer_reviews_customer_id` (customer_id)
- `idx_customer_reviews_rating` (rating)
- `idx_customer_reviews_status` (status)
- `idx_customer_reviews_created_at` (created_at DESC)

---

### 4.5 `review_photos` ✨ NEW
**Purpose:** Photos uploaded with reviews  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Photo ID |
| review_id | INT | NOT NULL, FK → customer_reviews.id | Review reference |
| url | TEXT | NOT NULL | Photo URL |
| thumbnail_url | TEXT | | Thumbnail URL |
| file_size | INT | | File size (bytes) |
| mime_type | VARCHAR(50) | | MIME type |
| display_order | INT | DEFAULT 0 | Sort order |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Upload timestamp |

**Indexes:**
- `idx_review_photos_review_id` (review_id)

---

### 4.6 `review_items` ✨ NEW
**Purpose:** Individual menu item ratings within a review  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Review item ID |
| review_id | INT | NOT NULL, FK → customer_reviews.id | Review reference |
| menu_item_id | INT | NOT NULL, FK → menu_items.id | Menu item reference |
| item_rating | INT | CHECK 1-5 | Item-specific rating |
| item_comment | TEXT | | Item-specific comment |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_review_items_review_id` (review_id)
- `idx_review_items_menu_item_id` (menu_item_id)

---

### 4.7 `customer_favorite_items` ✨ NEW
**Purpose:** Customer favorite menu items  
**Migration:** `007_create_customer_profile_reviews.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Favorite ID |
| customer_id | UUID | NOT NULL, FK → users.id | Customer reference |
| menu_item_id | INT | NOT NULL, FK → menu_items.id | Menu item reference |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Favorite timestamp |

**Unique:** `(customer_id, menu_item_id)`

**Indexes:**
- `idx_customer_favorite_items_customer_id` (customer_id)
- `idx_customer_favorite_items_menu_item_id` (menu_item_id)

---

## 5. STAFF MANAGEMENT & PERMISSIONS

**Migration:** `008_create_staff_management.sql`

### 5.1 `staff_profiles` ✨ NEW
**Purpose:** Extended staff information  
**Migration:** `008_create_staff_management.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Profile ID |
| user_id | UUID | NOT NULL UNIQUE, FK → users.id | User reference |
| employee_id | VARCHAR(50) | UNIQUE | Employee ID |
| department | VARCHAR(100) | | Department |
| position | VARCHAR(100) | | Job title |
| hire_date | DATE | | Hire date |
| employment_status | VARCHAR(20) | DEFAULT 'active', CHECK IN ('active', 'inactive', 'on_leave', 'terminated') | Employment status |
| shift_type | VARCHAR(20) | CHECK IN ('morning', 'afternoon', 'evening', 'night', 'rotating', 'flexible') | Shift type |
| weekly_hours | INT | CHECK 0-168 | Weekly hours |
| emergency_contact_name | VARCHAR(255) | | Emergency contact |
| emergency_contact_phone | VARCHAR(20) | | Emergency phone |
| emergency_contact_relationship | VARCHAR(50) | | Relationship |
| total_orders_served | INT | DEFAULT 0, CHECK >= 0 | Orders served (waiters) |
| total_orders_prepared | INT | DEFAULT 0, CHECK >= 0 | Orders prepared (kitchen) |
| average_order_time | INT | | Avg time (minutes) |
| average_rating | DECIMAL(3,2) | CHECK 0-5 | Performance rating |
| notes | TEXT | | Additional notes |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_staff_profiles_user_id` (user_id)
- `idx_staff_profiles_employee_id` (employee_id)
- `idx_staff_profiles_employment_status` (employment_status)

---

### 5.2 `waiter_table_assignments` ✨ NEW
**Purpose:** Waiter table assignments  
**Migration:** `008_create_staff_management.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Assignment ID |
| waiter_id | UUID | NOT NULL, FK → users.id | Waiter reference |
| table_id | INT | NOT NULL, FK → tables.id | Table reference |
| assigned_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Assignment timestamp |
| assigned_by | UUID | FK → users.id | Admin/Manager who assigned |
| is_active | BOOLEAN | DEFAULT TRUE | Active assignment |

**Unique:** `(waiter_id, table_id)`

**Indexes:**
- `idx_waiter_table_assignments_waiter_id` (waiter_id)
- `idx_waiter_table_assignments_table_id` (table_id)
- `idx_waiter_table_assignments_is_active` (is_active)

**Note:** Permissions management is handled by the existing `action_control_list` table.

---

### 5.3 `staff_invitations` ✨ NEW
**Purpose:** Email invitations for new staff  
**Migration:** `008_create_staff_management.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Invitation ID |
| email | VARCHAR(255) | NOT NULL | Invitee email |
| role_id | UUID | NOT NULL, FK → roles.id | Assigned role |
| invited_by | UUID | NOT NULL, FK → users.id | Inviter user |
| token | VARCHAR(255) | UNIQUE NOT NULL | Invitation token |
| status | VARCHAR(20) | DEFAULT 'pending', CHECK IN ('pending', 'accepted', 'expired', 'cancelled') | Invitation status |
| expires_at | TIMESTAMP | NOT NULL | Expiration time |
| accepted_at | TIMESTAMP | | Acceptance timestamp |
| accepted_by | UUID | FK → users.id | User who accepted |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_staff_invitations_email` (email)
- `idx_staff_invitations_token` (token)
- `idx_staff_invitations_status` (status)

---

### 5.4 `staff_activity_log` ✨ NEW
**Purpose:** Audit trail of staff actions  
**Migration:** `008_create_staff_management.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Log ID |
| staff_id | UUID | NOT NULL, FK → users.id | Staff user |
| action | VARCHAR(100) | NOT NULL | Action performed (login/logout/create_order/etc.) |
| action_category | VARCHAR(50) | | Category (auth/orders/billing/menu/staff) |
| description | TEXT | | Action description |
| entity_type | VARCHAR(50) | | Affected entity (order/bill/menu_item/staff) |
| entity_id | INT | | Entity ID |
| ip_address | INET | | Request IP |
| user_agent | TEXT | | Browser user agent |
| status | VARCHAR(20) | CHECK IN ('success', 'failed', 'error') | Action result |
| error_message | TEXT | | Error details |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Log timestamp |

**Indexes:**
- `idx_staff_activity_log_staff_id` (staff_id)
- `idx_staff_activity_log_action` (action)
- `idx_staff_activity_log_created_at` (created_at DESC)
- `idx_staff_activity_log_entity` (entity_type, entity_id)

---

## 6. ANALYTICS & REPORTS

**Migration:** `009_create_analytics_reports.sql`

### 6.1 `daily_revenue_snapshots` ✨ NEW
**Purpose:** Pre-calculated daily revenue aggregates  
**Migration:** `009_create_analytics_reports.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Snapshot ID |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| snapshot_date | DATE | NOT NULL | Snapshot date |
| total_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Total revenue |
| subtotal | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Subtotal |
| tax_collected | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Tax collected |
| discounts_given | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Discounts applied |
| service_charges | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Service charges |
| total_orders | INT | DEFAULT 0, CHECK >= 0 | Total orders |
| completed_orders | INT | DEFAULT 0, CHECK >= 0 | Completed orders |
| cancelled_orders | INT | DEFAULT 0, CHECK >= 0 | Cancelled orders |
| average_order_value | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Average order value |
| total_customers | INT | DEFAULT 0, CHECK >= 0 | Total customers |
| new_customers | INT | DEFAULT 0, CHECK >= 0 | New customers |
| returning_customers | INT | DEFAULT 0, CHECK >= 0 | Returning customers |
| cash_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Cash revenue |
| card_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Card revenue |
| ewallet_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | E-wallet revenue |
| stripe_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Stripe revenue |
| tables_used | INT | DEFAULT 0, CHECK >= 0 | Tables used |
| average_table_turnover | DECIMAL(5,2) | | Avg table turnover |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Unique:** `(restaurant_id, snapshot_date)`

**Indexes:**
- `idx_daily_revenue_snapshots_restaurant_date` (restaurant_id, snapshot_date DESC)
- `idx_daily_revenue_snapshots_date` (snapshot_date DESC)

**Purpose:** Fast reporting without heavy aggregation queries. Calculated via cron job.

---

### 6.2 `hourly_revenue_breakdown` ✨ NEW
**Purpose:** Hourly revenue for peak hours analysis  
**Migration:** `009_create_analytics_reports.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Breakdown ID |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| snapshot_date | DATE | NOT NULL | Date |
| hour | INT | NOT NULL, CHECK 0-23 | Hour (0-23) |
| revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Hourly revenue |
| orders_count | INT | DEFAULT 0, CHECK >= 0 | Hourly orders |
| customers_count | INT | DEFAULT 0, CHECK >= 0 | Hourly customers |
| average_order_value | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Avg order value |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Unique:** `(restaurant_id, snapshot_date, hour)`

**Indexes:**
- `idx_hourly_revenue_breakdown_restaurant_date` (restaurant_id, snapshot_date, hour)

---

### 6.3 `item_sales_statistics` ✨ NEW
**Purpose:** Menu item sales statistics by period  
**Migration:** `009_create_analytics_reports.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Statistic ID |
| menu_item_id | INT | NOT NULL, FK → menu_items.id | Menu item reference |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| period_type | VARCHAR(20) | NOT NULL, CHECK IN ('daily', 'weekly', 'monthly', 'all_time') | Period type |
| period_start | DATE | NOT NULL | Period start |
| period_end | DATE | NOT NULL | Period end |
| total_quantity_sold | INT | DEFAULT 0, CHECK >= 0 | Quantity sold |
| total_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Revenue generated |
| orders_count | INT | DEFAULT 0, CHECK >= 0 | Orders containing item |
| average_rating | DECIMAL(3,2) | CHECK 0-5 | Average rating |
| reviews_count | INT | DEFAULT 0, CHECK >= 0 | Review count |
| rank_by_quantity | INT | | Rank by quantity |
| rank_by_revenue | INT | | Rank by revenue |
| percentage_of_total_revenue | DECIMAL(5,2) | | % of total revenue |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Unique:** `(menu_item_id, period_type, period_start)`

**Indexes:**
- `idx_item_sales_statistics_menu_item` (menu_item_id)
- `idx_item_sales_statistics_period` (period_start, period_end)
- `idx_item_sales_statistics_rank_revenue` (rank_by_revenue)

---

### 6.4 `category_sales_statistics` ✨ NEW
**Purpose:** Category sales statistics by period  
**Migration:** `009_create_analytics_reports.sql`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Statistic ID |
| category_id | INT | NOT NULL, FK → menu_categories.id | Category reference |
| restaurant_id | INT | NOT NULL, FK → restaurants.id | Restaurant reference |
| period_type | VARCHAR(20) | NOT NULL, CHECK IN ('daily', 'weekly', 'monthly', 'all_time') | Period type |
| period_start | DATE | NOT NULL | Period start |
| period_end | DATE | NOT NULL | Period end |
| total_quantity_sold | INT | DEFAULT 0, CHECK >= 0 | Quantity sold |
| total_revenue | DECIMAL(12,2) | DEFAULT 0.00, CHECK >= 0 | Revenue generated |
| orders_count | INT | DEFAULT 0, CHECK >= 0 | Orders count |
| percentage_of_total_revenue | DECIMAL(5,2) | | % of total revenue |
| average_order_value | DECIMAL(10,2) | DEFAULT 0.00, CHECK >= 0 | Avg order value |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Unique:** `(category_id, period_type, period_start)`

**Indexes:**
- `idx_category_sales_statistics_category` (category_id)
- `idx_category_sales_statistics_period` (period_start, period_end)

---

## 7. RELATIONSHIPS DIAGRAM

```
┌──────────────┐
│ restaurants  │
└──────┬───────┘
       │
       ├──────► tables
       ├──────► menu_categories ──► menu_items ──► menu_item_photos
       ├──────► modifier_groups ──► modifier_options
       ├──────► orders ──► order_items ──► order_item_modifiers
       │                  │
       │                  ├──► order_timeline
       │                  ├──► order_notes
       │                  ├──► kitchen_alerts
       │                  │
       │                  └──► bills ──► payments
       │                         │
       │                         └──► discount_usage ──► discount_codes
       │
       └──────► daily_revenue_snapshots
                hourly_revenue_breakdown
                item_sales_statistics
                category_sales_statistics

┌──────────┐
│  roles   │──► action_control_list (existing RBAC table)
└────┬─────┘
     │
     ├──────► users ──► customer_preferences
     │              │
     │              ├──► customer_loyalty ──► loyalty_transactions
     │              │
     │              ├──► customer_reviews ──► review_photos
     │              │                      └──► review_items
     │              │
     │              ├──► customer_favorite_items
     │              │
     │              ├──► staff_profiles
     │              │
     │              ├──► waiter_table_assignments
     │              │
     │              └──► staff_activity_log
     │
     └──────► staff_invitations
```

---

## 8. INDEXES SUMMARY

### Performance Indexes
- **Foreign Keys:** All foreign key columns have indexes
- **Status Fields:** All `status` columns have indexes for filtering
- **Timestamps:** `created_at DESC` indexes for sorting
- **Composite:** Restaurant + Date indexes for multi-tenant queries

### Unique Constraints
- Order numbers, bill numbers, payment IDs
- QR tokens, table numbers
- Email addresses, employee IDs
- Review per order per customer

---

## 9. VIEWS SUMMARY

### 9.1 `customer_statistics`
**Purpose:** Aggregated customer metrics  
**Fields:** total_orders, total_spent, average_order_value, loyalty_points, favorite_items

### 9.2 `staff_statistics`
**Purpose:** Aggregated staff performance  
**Fields:** total_orders_served/prepared, revenue_generated, assigned_tables, average_rating

### 9.3 `dashboard_summary`
**Purpose:** Real-time admin dashboard  
**Fields:** revenue_today, orders_today, active_orders, tables_occupied, growth_rate

### 9.4 `kitchen_dashboard`
**Purpose:** Real-time kitchen queue  
**Fields:** pending/preparing/ready counts, average_prep_time, longest_waiting_order

### 9.5 `waiter_dashboard`
**Purpose:** Per-waiter dashboard  
**Fields:** orders_served_today, revenue_generated, active_orders, assigned_tables

---

## 10. TRIGGERS & FUNCTIONS

### Auto-Update Timestamps
- `update_updated_at_column()` - Updates `updated_at` on row modification
- Applied to: orders, order_items, bills, payments, users, profiles, etc.

### Auto-Create Profiles
- `create_customer_profile()` - Creates preferences + loyalty on user registration
- `create_staff_profile()` - Creates staff profile on staff user creation

### Loyalty System
- `update_loyalty_tier()` - Auto-updates tier based on points
- `award_loyalty_points()` - Awards points on order completion (1 point per $1)

### Activity Logging
- `log_staff_login()` - Logs staff login events

### Daily Snapshots
- `calculate_daily_revenue_snapshot(date, restaurant_id)` - Manual/cron function to calculate daily revenue

---

## 📊 SUMMARY STATISTICS

| Category | Tables | New Tables | Indexes | Views | Triggers |
|----------|--------|------------|---------|-------|----------|
| Core | 12 | 0 | ~20 | 0 | 2 |
| Orders | 6 | 4 | ~18 | 0 | 2 |
| Bills & Payments | 4 | 4 | ~12 | 0 | 2 |
| Customer Profile | 7 | 7 | ~15 | 1 | 3 |
| Staff Management | 4 | 4 | ~9 | 1 | 2 |
| Analytics | 4 | 4 | ~8 | 3 | 3 |
| **TOTAL** | **37** | **23** | **~82** | **5** | **14** |

---

## 🎯 NEXT STEPS

1. **Run Migrations:** Execute migrations 006-009 in order
2. **Test Data:** Create seed data for testing
3. **API Development:** Implement backend handlers using these schemas
4. **Frontend Integration:** Connect APIs to React components
5. **Performance Testing:** Monitor query performance and optimize indexes

---

**Last Updated:** January 11, 2026  
**Author:** Smart Restaurant Development Team
