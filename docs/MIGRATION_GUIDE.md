# 🚀 DATABASE MIGRATION EXECUTION GUIDE

## 📋 MIGRATION ORDER

Execute migrations in this **exact order** to ensure proper foreign key relationships:

```bash
# 1. Core tables (existing)
001_create_tables.sql                    # ✅ Tables, Orders, Order Items
002_create_restaurant.sql                # ✅ Restaurants
003_create_menu_modifiers.sql            # ✅ Menu Categories, Items, Modifiers
004_import_menu_photos.sql               # ✅ Menu photos import
005_create_authentication.sql            # ✅ Roles, Users, OTPs

# 2. New migrations (run in order)
006_create_orders_system.sql             # ⭐ Order Management, Bills, Payments
007_create_customer_profile_reviews.sql  # ⭐ Customer Profiles, Reviews, Loyalty
008_create_staff_management.sql          # ⭐ Staff Profiles, RBAC (via action_control_list)
009_create_analytics_reports.sql         # ⭐ Analytics, Reports, Dashboards
```

---

## ✅ WHAT EACH MIGRATION DOES

### Migration 006: Order Management System
**Dependencies:** 001, 002, 005  
**Creates:**
- ✨ Enhanced `orders` table with new fields (waiter_id, priority, source, etc.)
- ✨ Enhanced `order_items` table with modifiers support
- ✨ `order_item_modifiers` - Track selected modifiers (e.g., "Medium Rare", "Extra Cheese")
- ✨ `order_timeline` - Status change history for order tracking
- ✨ `order_notes` - Special requests, allergies, kitchen notes
- ✨ `kitchen_alerts` - Alerts from kitchen to waiters
- ✨ `bills` - Bills generated from orders
- ✨ `payments` - Payment transactions (Stripe, cash, card, e-wallet)
- ✨ `discount_codes` - Promotional discount codes
- ✨ `discount_usage` - Track discount usage per customer

**Use Cases:**
- ✅ TASK-001 to TASK-009: Order creation, status updates, timeline tracking
- ✅ TASK-010 to TASK-015: Bill generation, payment processing, discount validation
- ✅ Customer order tracking with real-time status
- ✅ Kitchen order queue management
- ✅ Waiter order assignment

---

### Migration 007: Customer Profile & Reviews
**Dependencies:** 001, 005, 006  
**Creates:**
- ✨ Enhanced `users` table with profile fields (avatar, address, gender, etc.)
- ✨ `customer_preferences` - Dietary restrictions, notifications, language
- ✨ `customer_loyalty` - Loyalty points, tier (bronze → diamond)
- ✨ `loyalty_transactions` - Points earned/redeemed history
- ✨ `customer_reviews` - Order reviews with ratings (1-5 stars)
- ✨ `review_photos` - Photos uploaded with reviews
- ✨ `review_items` - Individual menu item ratings
- ✨ `customer_favorite_items` - Quick reorder favorites
- ✨ `customer_statistics` VIEW - Aggregated customer metrics

**Use Cases:**
- ✅ TASK-016 to TASK-020: Customer profile management, reviews
- ✅ Customer loyalty program (1 point per $1 spent)
- ✅ Dietary restrictions and allergen management
- ✅ Review system with restaurant responses
- ✅ Favorite items for quick reordering

**Auto-Triggers:**
- ✅ Auto-create customer preferences + loyalty on user registration
- ✅ Auto-update loyalty tier when points change
- ✅ Auto-award loyalty points on order completion

---

### Migration 008: Staff Management
**Dependencies:** 001, 002, 005  
**Creates:**
- ✨ Enhanced `roles` table with staff flags
- ✨ New roles: `waiter`, `kitchen_staff`, `manager`, `cashier`
- ✨ `staff_profiles` - Employment info, performance metrics
- ✨ `waiter_table_assignments` - Assign tables to waiters
- ✨ `staff_invitations` - Email invitations for new staff
- ✨ `staff_activity_log` - Audit trail of staff actions
- ✨ `staff_statistics` VIEW - Staff performance metrics

**Note:** RBAC permissions managed via existing `action_control_list` table.

**Use Cases:**
- ✅ TASK-021 to TASK-030: Staff account management, role assignments
- ✅ Role-based access control (RBAC via action_control_list)
- ✅ Waiter table assignment
- ✅ Staff performance tracking
- ✅ Activity logging for security

**Auto-Triggers:**
- ✅ Auto-create staff profile on staff user creation
- ✅ Auto-log staff login events

---

### Migration 009: Analytics & Reports
**Dependencies:** 001, 002, 005, 006, 007  
**Creates:**
- ✨ `daily_revenue_snapshots` - Pre-calculated daily metrics (fast reporting)
- ✨ `hourly_revenue_breakdown` - Hourly revenue for peak hours analysis
- ✨ `item_sales_statistics` - Top selling items by period
- ✨ `category_sales_statistics` - Category performance by period
- ✨ `dashboard_summary` VIEW - Real-time admin dashboard
- ✨ `kitchen_dashboard` VIEW - Real-time kitchen queue stats
- ✨ `waiter_dashboard` VIEW - Per-waiter performance dashboard

**Use Cases:**
- ✅ TASK-031 to TASK-032: Dashboard analytics, revenue reports
- ✅ Peak hours analysis (which hours are busiest)
- ✅ Top selling items report
- ✅ Sales by category report
- ✅ Real-time dashboard for admin, kitchen, waiters

**Functions:**
- ✅ `calculate_daily_revenue_snapshot(date, restaurant_id)` - Manual/cron job to calculate daily metrics

---

## 🎯 QUICK START

### Step 1: Run All Migrations
```bash
# Navigate to migrations folder
cd /Users/mainhatnam/Documents/2025\ -\ 2026\ FIT/Web\ Application/Final/app-api/migrations

# Run existing migrations (if not already run)
psql -U your_username -d your_database -f 001_create_tables.sql
psql -U your_username -d your_database -f 002_create_restaurant.sql
psql -U your_username -d your_database -f 003_create_menu_modifiers.sql
psql -U your_username -d your_database -f 004_import_menu_photos.sql
psql -U your_username -d your_database -f 005_create_authentication.sql

# Run new migrations
psql -U your_username -d your_database -f 006_create_orders_system.sql
psql -U your_username -d your_database -f 007_create_customer_profile_reviews.sql
psql -U your_username -d your_database -f 008_create_staff_management.sql
psql -U your_username -d your_database -f 009_create_analytics_reports.sql
```

### Step 2: Verify Installation
```sql
-- Check tables created
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;

-- Should return ~39 tables

-- Check views created
SELECT table_name 
FROM information_schema.views 
WHERE table_schema = 'public';

-- Should return: customer_statistics, staff_statistics, dashboard_summary, kitchen_dashboard, waiter_dashboard

-- Check triggers
SELECT trigger_name, event_object_table, action_statement 
FROM information_schema.triggers 
WHERE trigger_schema = 'public';

-- Should return ~14 triggers
```

### Step 3: Test Data (Optional)
```sql
-- Create test customer
INSERT INTO users (email, password, first_name, last_name, full_name, phone_number, role, status)
SELECT 
    'test.customer@example.com',
    '$2a$10$...', -- hashed password
    'Test',
    'Customer',
    'Test Customer',
    '+1234567890',
    r.id,
    'active'
FROM roles r
WHERE r.name = 'end_user';

-- Create test waiter
INSERT INTO users (email, password, full_name, phone_number, role, status)
SELECT 
    'test.waiter@restaurant.com',
    '$2a$10$...', -- hashed password
    'Test Waiter',
    '+1234567891',
    r.id,
    'active'
FROM roles r
WHERE r.name = 'waiter';

-- Verify customer preferences + loyalty auto-created
SELECT * FROM customer_preferences;
SELECT * FROM customer_loyalty;

-- Verify staff profile auto-created
SELECT * FROM staff_profiles;
```

---

## 🔧 ROLLBACK (IF NEEDED)

### Rollback Order (Reverse of execution)
```bash
# Drop in reverse order
psql -U your_username -d your_database -c "DROP TABLE IF EXISTS category_sales_statistics, item_sales_statistics, hourly_revenue_breakdown, daily_revenue_snapshots CASCADE;"
psql -U your_username -d your_database -c "DROP TABLE IF EXISTS staff_activity_log, staff_invitations, waiter_table_assignments, staff_profiles CASCADE;"
psql -U your_username -d your_database -c "DROP TABLE IF EXISTS customer_favorite_items, review_items, review_photos, customer_reviews, loyalty_transactions, customer_loyalty, customer_preferences CASCADE;"
psql -U your_username -d your_database -c "DROP TABLE IF EXISTS discount_usage, discount_codes, payments, bills, kitchen_alerts, order_notes, order_timeline, order_item_modifiers CASCADE;"

# Drop views
psql -U your_username -d your_database -c "DROP VIEW IF EXISTS customer_statistics, staff_statistics, dashboard_summary, kitchen_dashboard, waiter_dashboard CASCADE;"

# Drop functions
psql -U your_username -d your_database -c "DROP FUNCTION IF EXISTS update_updated_at_column, create_customer_profile, create_staff_profile, update_loyalty_tier, award_loyalty_points, log_staff_login, calculate_daily_revenue_snapshot CASCADE;"
```

---

## 📊 DATABASE SCHEMA STATISTICS

| Metric | Count |
|--------|-------|
| Total Tables | 37 |
| New Tables (006-009) | 23 |
| Views | 5 |
| Indexes | ~82 |
| Triggers | 14 |
| Functions | 7 |
| Roles | 8 (end_user, admin, waiter, kitchen_staff, manager, cashier, restaurant_owner, staff) |

---

## ✅ POST-MIGRATION CHECKLIST

- [ ] All migrations executed successfully
- [ ] 39 tables exist in database
- [ ] 5 views created
- [ ] 14 triggers active
- [ ] Test customer auto-creates preferences + loyalty
- [ ] Test staff auto-creates staff profile
- [ ] Loyalty points auto-awarded on order completion
- [ ] Daily revenue snapshot function works
- [ ] All foreign keys enforced
- [ ] Indexes created for performance

---

## 🐛 TROUBLESHOOTING

### Error: "relation already exists"
**Solution:** Migration already run. Check existing tables:
```sql
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
```

### Error: "foreign key constraint violation"
**Solution:** Run migrations in order (001 → 009). Dependencies must exist first.

### Error: "function does not exist"
**Solution:** Ensure `update_updated_at_column()` function exists (created in migration 006).

### Performance Issues
**Solution:** 
1. Check indexes: `SELECT * FROM pg_indexes WHERE schemaname = 'public';`
2. Run ANALYZE: `ANALYZE VERBOSE;`
3. Check query plans: `EXPLAIN ANALYZE SELECT ...`

---

## 📚 RELATED DOCUMENTATION

- [DATABASE_SCHEMA_DOCUMENTATION.md](./DATABASE_SCHEMA_DOCUMENTATION.md) - Full schema reference
- [GOOGLE_OAUTH_GET_TOKEN.md](../docs/GOOGLE_OAUTH_GET_TOKEN.md) - OAuth setup
- [API Examples](../docs/) - API usage examples

---

**Last Updated:** January 11, 2026  
**Migration Version:** v2.0  
**Database:** PostgreSQL 14+
