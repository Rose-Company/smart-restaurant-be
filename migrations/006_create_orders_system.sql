-- =====================================================
-- 006_create_orders_system.sql
-- COMPREHENSIVE ORDER MANAGEMENT SYSTEM
-- =====================================================
-- This migration creates all tables needed for:
-- - Order Management (TASK-001 to TASK-009)
-- - Bills & Payments (TASK-010 to TASK-015)
-- - Order Timeline & Status Tracking
-- - Order Item Modifiers
-- =====================================================

-- =====================================================
-- 1. ENHANCE EXISTING ORDERS TABLE
-- =====================================================

-- Add missing fields to orders table
ALTER TABLE public.orders
ADD COLUMN IF NOT EXISTS order_number VARCHAR(50),
ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS customer_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS waiter_id UUID,
ADD COLUMN IF NOT EXISTS kitchen_staff_id UUID,
ADD COLUMN IF NOT EXISTS estimated_ready_time TIMESTAMP,
ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS cancelled_by VARCHAR(50),
ADD COLUMN IF NOT EXISTS cancel_reason TEXT,
ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
ADD COLUMN IF NOT EXISTS source VARCHAR(20) DEFAULT 'qr' CHECK (source IN ('qr', 'waiter', 'admin', 'customer_app'));

-- Update status column to include all required statuses
ALTER TABLE public.orders
DROP CONSTRAINT IF EXISTS orders_status_check,
ADD CONSTRAINT orders_status_check CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'cancelled'));

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_orders_status ON public.orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_table_id ON public.orders(table_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_user_id ON public.orders(customer_user_id);
CREATE INDEX IF NOT EXISTS idx_orders_waiter_id ON public.orders(waiter_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON public.orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON public.orders(order_number);

-- Add foreign key for waiter_id
ALTER TABLE public.orders
DROP CONSTRAINT IF EXISTS orders_waiter_id_fkey,
ADD CONSTRAINT orders_waiter_id_fkey FOREIGN KEY (waiter_id) REFERENCES public.users(id) ON DELETE SET NULL;

-- =====================================================
-- 2. ENHANCE EXISTING ORDER_ITEMS TABLE
-- =====================================================

-- Add missing fields to order_items table
ALTER TABLE public.order_items
ADD COLUMN IF NOT EXISTS modifiers_total DECIMAL(10,2) DEFAULT 0.00 CHECK (modifiers_total >= 0),
ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(10,2) DEFAULT 0.00 CHECK (tax_amount >= 0),
ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(10,2) DEFAULT 0.00 CHECK (discount_amount >= 0),
ADD COLUMN IF NOT EXISTS final_price DECIMAL(10,2) GENERATED ALWAYS AS (subtotal + modifiers_total - discount_amount) STORED;

-- Update status column to include item-level statuses
ALTER TABLE public.order_items
DROP CONSTRAINT IF EXISTS order_items_status_check,
ADD CONSTRAINT order_items_status_check CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'rejected', 'cancelled'));

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON public.order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_menu_item_id ON public.order_items(menu_item_id);
CREATE INDEX IF NOT EXISTS idx_order_items_status ON public.order_items(status);

-- =====================================================
-- 3. ORDER ITEM MODIFIERS (For tracking selected modifiers)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.order_item_modifiers (
    id SERIAL PRIMARY KEY,
    order_item_id INT NOT NULL,
    modifier_group_id INT NOT NULL,
    modifier_group_name VARCHAR(80) NOT NULL,
    modifier_option_id INT NOT NULL,
    modifier_option_name VARCHAR(80) NOT NULL,
    price_adjustment DECIMAL(10,2) DEFAULT 0.00 CHECK (price_adjustment >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT order_item_modifiers_order_item_id_fkey 
        FOREIGN KEY (order_item_id) REFERENCES public.order_items(id) ON DELETE CASCADE,
    CONSTRAINT order_item_modifiers_modifier_group_id_fkey 
        FOREIGN KEY (modifier_group_id) REFERENCES public.modifier_groups(id) ON DELETE RESTRICT,
    CONSTRAINT order_item_modifiers_modifier_option_id_fkey 
        FOREIGN KEY (modifier_option_id) REFERENCES public.modifier_options(id) ON DELETE RESTRICT
);

CREATE INDEX idx_order_item_modifiers_order_item_id ON public.order_item_modifiers(order_item_id);

-- =====================================================
-- 4. ORDER TIMELINE (For tracking status changes)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.order_timeline (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'cancelled')),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50) NOT NULL, -- 'customer', 'waiter', 'kitchen', 'admin', 'system'
    updated_by_id UUID, -- Reference to users.id
    updated_by_name VARCHAR(255),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT order_timeline_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
    CONSTRAINT order_timeline_updated_by_id_fkey 
        FOREIGN KEY (updated_by_id) REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX idx_order_timeline_order_id ON public.order_timeline(order_id);
CREATE INDEX idx_order_timeline_timestamp ON public.order_timeline(timestamp DESC);

-- =====================================================
-- 5. BILLS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.bills (
    id SERIAL PRIMARY KEY,
    bill_number VARCHAR(50) UNIQUE NOT NULL,
    order_id INT NOT NULL,
    restaurant_id INT,
    table_id INT,
    customer_id UUID,
    
    -- Amounts
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (subtotal >= 0),
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (tax_amount >= 0),
    tax_rate DECIMAL(5,2) DEFAULT 10.00 CHECK (tax_rate >= 0 AND tax_rate <= 100),
    discount_amount DECIMAL(10,2) DEFAULT 0.00 CHECK (discount_amount >= 0),
    discount_code VARCHAR(50),
    service_charge DECIMAL(10,2) DEFAULT 0.00 CHECK (service_charge >= 0),
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    
    -- Bill metadata
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
    bill_type VARCHAR(20) DEFAULT 'request' CHECK (bill_type IN ('request', 'generated', 'manual')),
    payment_method VARCHAR(20) CHECK (payment_method IN ('cash', 'card', 'ewallet', 'stripe', 'other')),
    
    -- Tracking
    requested_by VARCHAR(50), -- 'customer', 'waiter', 'system'
    requested_by_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP,
    
    CONSTRAINT bills_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE RESTRICT,
    CONSTRAINT bills_restaurant_id_fkey 
        FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE,
    CONSTRAINT bills_table_id_fkey 
        FOREIGN KEY (table_id) REFERENCES public.tables(id) ON DELETE SET NULL,
    CONSTRAINT bills_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE SET NULL,
    CONSTRAINT bills_requested_by_id_fkey 
        FOREIGN KEY (requested_by_id) REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX idx_bills_bill_number ON public.bills(bill_number);
CREATE INDEX idx_bills_order_id ON public.bills(order_id);
CREATE INDEX idx_bills_status ON public.bills(status);
CREATE INDEX idx_bills_created_at ON public.bills(created_at DESC);
CREATE INDEX idx_bills_customer_id ON public.bills(customer_id);

-- =====================================================
-- 6. PAYMENTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.payments (
    id SERIAL PRIMARY KEY,
    payment_id VARCHAR(100) UNIQUE NOT NULL,
    bill_id INT NOT NULL,
    order_id INT NOT NULL,
    
    -- Payment details
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    method VARCHAR(20) NOT NULL CHECK (method IN ('cash', 'card', 'ewallet', 'stripe', 'other')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'succeeded', 'failed', 'cancelled', 'refunded')),
    
    -- Cash payment details
    received_amount DECIMAL(10,2) CHECK (received_amount >= 0),
    change_amount DECIMAL(10,2) CHECK (change_amount >= 0),
    
    -- Stripe payment details
    stripe_payment_intent_id VARCHAR(255),
    stripe_charge_id VARCHAR(255),
    stripe_payment_method_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    receipt_url TEXT,
    
    -- Error tracking
    error_code VARCHAR(50),
    error_message TEXT,
    decline_reason VARCHAR(255),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    failed_at TIMESTAMP,
    refunded_at TIMESTAMP,
    
    -- Metadata
    metadata JSONB,
    
    CONSTRAINT payments_bill_id_fkey 
        FOREIGN KEY (bill_id) REFERENCES public.bills(id) ON DELETE RESTRICT,
    CONSTRAINT payments_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE RESTRICT
);

CREATE INDEX idx_payments_payment_id ON public.payments(payment_id);
CREATE INDEX idx_payments_bill_id ON public.payments(bill_id);
CREATE INDEX idx_payments_order_id ON public.payments(order_id);
CREATE INDEX idx_payments_status ON public.payments(status);
CREATE INDEX idx_payments_stripe_payment_intent_id ON public.payments(stripe_payment_intent_id);
CREATE INDEX idx_payments_created_at ON public.payments(created_at DESC);

-- =====================================================
-- 7. DISCOUNT CODES TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.discount_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    
    -- Discount configuration
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value DECIMAL(10,2) NOT NULL CHECK (discount_value > 0),
    max_discount_amount DECIMAL(10,2) CHECK (max_discount_amount >= 0),
    min_order_amount DECIMAL(10,2) DEFAULT 0.00 CHECK (min_order_amount >= 0),
    
    -- Usage limits
    usage_limit INT CHECK (usage_limit > 0),
    usage_count INT DEFAULT 0 CHECK (usage_count >= 0),
    per_customer_limit INT DEFAULT 1 CHECK (per_customer_limit > 0),
    
    -- Validity
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Applicability
    applicable_categories INT[], -- Array of category IDs
    applicable_items INT[], -- Array of menu item IDs
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT discount_codes_valid_dates_check 
        CHECK (valid_until > valid_from)
);

CREATE INDEX idx_discount_codes_code ON public.discount_codes(code);
CREATE INDEX idx_discount_codes_is_active ON public.discount_codes(is_active);
CREATE INDEX idx_discount_codes_valid_from ON public.discount_codes(valid_from);
CREATE INDEX idx_discount_codes_valid_until ON public.discount_codes(valid_until);

-- =====================================================
-- 8. DISCOUNT USAGE TRACKING
-- =====================================================

CREATE TABLE IF NOT EXISTS public.discount_usage (
    id SERIAL PRIMARY KEY,
    discount_id INT NOT NULL,
    customer_id UUID,
    order_id INT,
    bill_id INT,
    discount_amount DECIMAL(10,2) NOT NULL CHECK (discount_amount >= 0),
    used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT discount_usage_discount_id_fkey 
        FOREIGN KEY (discount_id) REFERENCES public.discount_codes(id) ON DELETE CASCADE,
    CONSTRAINT discount_usage_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE SET NULL,
    CONSTRAINT discount_usage_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
    CONSTRAINT discount_usage_bill_id_fkey 
        FOREIGN KEY (bill_id) REFERENCES public.bills(id) ON DELETE CASCADE
);

CREATE INDEX idx_discount_usage_discount_id ON public.discount_usage(discount_id);
CREATE INDEX idx_discount_usage_customer_id ON public.discount_usage(customer_id);
CREATE INDEX idx_discount_usage_used_at ON public.discount_usage(used_at DESC);

-- =====================================================
-- 9. ORDER NOTES/ALERTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.order_notes (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    note TEXT NOT NULL,
    note_type VARCHAR(20) DEFAULT 'general' CHECK (note_type IN ('general', 'allergy', 'special_request', 'kitchen_note', 'alert')),
    created_by VARCHAR(50) NOT NULL, -- 'customer', 'waiter', 'kitchen', 'admin'
    created_by_id UUID,
    created_by_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT order_notes_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
    CONSTRAINT order_notes_created_by_id_fkey 
        FOREIGN KEY (created_by_id) REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX idx_order_notes_order_id ON public.order_notes(order_id);
CREATE INDEX idx_order_notes_created_at ON public.order_notes(created_at DESC);

-- =====================================================
-- 10. KITCHEN ALERTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.kitchen_alerts (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    alert_type VARCHAR(50) NOT NULL CHECK (alert_type IN ('order_ready', 'item_ready', 'delay_warning', 'special_request', 'urgent')),
    message TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    sent_to VARCHAR(50) NOT NULL, -- 'waiter', 'kitchen', 'all_waiters', 'specific_waiter'
    waiter_id UUID,
    status VARCHAR(20) DEFAULT 'sent' CHECK (status IN ('sent', 'acknowledged', 'resolved', 'dismissed')),
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    acknowledged_at TIMESTAMP,
    resolved_at TIMESTAMP,
    
    CONSTRAINT kitchen_alerts_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
    CONSTRAINT kitchen_alerts_waiter_id_fkey 
        FOREIGN KEY (waiter_id) REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX idx_kitchen_alerts_order_id ON public.kitchen_alerts(order_id);
CREATE INDEX idx_kitchen_alerts_status ON public.kitchen_alerts(status);
CREATE INDEX idx_kitchen_alerts_waiter_id ON public.kitchen_alerts(waiter_id);
CREATE INDEX idx_kitchen_alerts_sent_at ON public.kitchen_alerts(sent_at DESC);

-- =====================================================
-- 11. SEED DATA - DISCOUNT CODES
-- =====================================================

INSERT INTO public.discount_codes (code, description, discount_type, discount_value, max_discount_amount, min_order_amount, usage_limit, valid_from, valid_until, is_active) VALUES
('SAVE5', 'Get $5 off on orders above $50', 'fixed_amount', 5.00, 5.00, 50.00, 100, '2026-01-01 00:00:00', '2026-12-31 23:59:59', TRUE),
('SAVE10', 'Get $10 off on orders above $100', 'fixed_amount', 10.00, 10.00, 100.00, 50, '2026-01-01 00:00:00', '2026-12-31 23:59:59', TRUE),
('PERCENT10', 'Get 10% off on all orders', 'percentage', 10.00, 20.00, 50.00, 200, '2026-01-01 00:00:00', '2026-06-30 23:59:59', TRUE),
('PERCENT20', 'Get 20% off on orders above $200', 'percentage', 20.00, 50.00, 200.00, 100, '2026-01-01 00:00:00', '2026-03-31 23:59:59', TRUE),
('WELCOME15', 'Welcome discount - 15% off first order', 'percentage', 15.00, 30.00, 0.00, NULL, '2026-01-01 00:00:00', '2026-12-31 23:59:59', TRUE)
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- TRIGGERS FOR AUTO-UPDATE TIMESTAMPS
-- =====================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to orders
DROP TRIGGER IF EXISTS update_orders_updated_at ON public.orders;
CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON public.orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to order_items
DROP TRIGGER IF EXISTS update_order_items_updated_at ON public.order_items;
CREATE TRIGGER update_order_items_updated_at
    BEFORE UPDATE ON public.order_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to bills
DROP TRIGGER IF EXISTS update_bills_updated_at ON public.bills;
CREATE TRIGGER update_bills_updated_at
    BEFORE UPDATE ON public.bills
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to payments
DROP TRIGGER IF EXISTS update_payments_updated_at ON public.payments;
CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON public.payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to discount_codes
DROP TRIGGER IF EXISTS update_discount_codes_updated_at ON public.discount_codes;
CREATE TRIGGER update_discount_codes_updated_at
    BEFORE UPDATE ON public.discount_codes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- COMMENTS FOR DOCUMENTATION
-- =====================================================

COMMENT ON TABLE public.orders IS 'Main orders table - stores customer orders from QR, waiter, or admin';
COMMENT ON TABLE public.order_items IS 'Order line items - individual menu items in an order';
COMMENT ON TABLE public.order_item_modifiers IS 'Selected modifiers for each order item (e.g., steak temperature, extras)';
COMMENT ON TABLE public.order_timeline IS 'Audit trail of order status changes for tracking and customer visibility';
COMMENT ON TABLE public.bills IS 'Bills generated from orders - supports discounts, taxes, service charges';
COMMENT ON TABLE public.payments IS 'Payment transactions - supports cash, card, e-wallet, Stripe';
COMMENT ON TABLE public.discount_codes IS 'Promotional discount codes with usage limits and validity periods';
COMMENT ON TABLE public.discount_usage IS 'Tracks which customers used which discount codes';
COMMENT ON TABLE public.order_notes IS 'Additional notes, special requests, and allergy information for orders';
COMMENT ON TABLE public.kitchen_alerts IS 'Alerts sent from kitchen to waiters (e.g., order ready, delays)';
