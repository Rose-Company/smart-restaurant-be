-- =====================================================
-- 009_create_analytics_reports.sql
-- ANALYTICS & REPORTS SYSTEM
-- =====================================================
-- This migration creates all tables needed for:
-- - Dashboard Analytics (TASK-031)
-- - Reports & Analytics (TASK-032)
-- - Revenue Tracking
-- - Performance Metrics
-- =====================================================

-- =====================================================
-- 1. DAILY REVENUE SNAPSHOTS (For fast reporting)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.daily_revenue_snapshots (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    snapshot_date DATE NOT NULL,
    
    -- Revenue breakdown
    total_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (total_revenue >= 0),
    subtotal DECIMAL(12,2) DEFAULT 0.00 CHECK (subtotal >= 0),
    tax_collected DECIMAL(12,2) DEFAULT 0.00 CHECK (tax_collected >= 0),
    discounts_given DECIMAL(12,2) DEFAULT 0.00 CHECK (discounts_given >= 0),
    service_charges DECIMAL(12,2) DEFAULT 0.00 CHECK (service_charges >= 0),
    
    -- Order statistics
    total_orders INT DEFAULT 0 CHECK (total_orders >= 0),
    completed_orders INT DEFAULT 0 CHECK (completed_orders >= 0),
    cancelled_orders INT DEFAULT 0 CHECK (cancelled_orders >= 0),
    average_order_value DECIMAL(10,2) DEFAULT 0.00 CHECK (average_order_value >= 0),
    
    -- Customer statistics
    total_customers INT DEFAULT 0 CHECK (total_customers >= 0),
    new_customers INT DEFAULT 0 CHECK (new_customers >= 0),
    returning_customers INT DEFAULT 0 CHECK (returning_customers >= 0),
    
    -- Payment method breakdown
    cash_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (cash_revenue >= 0),
    card_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (card_revenue >= 0),
    ewallet_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (ewallet_revenue >= 0),
    stripe_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (stripe_revenue >= 0),
    
    -- Table statistics
    tables_used INT DEFAULT 0 CHECK (tables_used >= 0),
    average_table_turnover DECIMAL(5,2), -- times per day
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT daily_revenue_snapshots_restaurant_id_fkey 
        FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE,
    
    UNIQUE (restaurant_id, snapshot_date)
);

CREATE INDEX idx_daily_revenue_snapshots_restaurant_date ON public.daily_revenue_snapshots(restaurant_id, snapshot_date DESC);
CREATE INDEX idx_daily_revenue_snapshots_date ON public.daily_revenue_snapshots(snapshot_date DESC);

-- =====================================================
-- 2. HOURLY REVENUE BREAKDOWN (For peak hours analysis)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.hourly_revenue_breakdown (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    snapshot_date DATE NOT NULL,
    hour INT NOT NULL CHECK (hour >= 0 AND hour <= 23),
    
    -- Revenue & orders
    revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (revenue >= 0),
    orders_count INT DEFAULT 0 CHECK (orders_count >= 0),
    customers_count INT DEFAULT 0 CHECK (customers_count >= 0),
    average_order_value DECIMAL(10,2) DEFAULT 0.00 CHECK (average_order_value >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT hourly_revenue_breakdown_restaurant_id_fkey 
        FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE,
    
    UNIQUE (restaurant_id, snapshot_date, hour)
);

CREATE INDEX idx_hourly_revenue_breakdown_restaurant_date ON public.hourly_revenue_breakdown(restaurant_id, snapshot_date, hour);

-- =====================================================
-- 3. TOP SELLING ITEMS (For reports)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.item_sales_statistics (
    id SERIAL PRIMARY KEY,
    menu_item_id INT NOT NULL,
    restaurant_id INT NOT NULL,
    
    -- Time period
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('daily', 'weekly', 'monthly', 'all_time')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Sales statistics
    total_quantity_sold INT DEFAULT 0 CHECK (total_quantity_sold >= 0),
    total_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (total_revenue >= 0),
    orders_count INT DEFAULT 0 CHECK (orders_count >= 0),
    average_rating DECIMAL(3,2) CHECK (average_rating >= 0 AND average_rating <= 5),
    reviews_count INT DEFAULT 0 CHECK (reviews_count >= 0),
    
    -- Rankings
    rank_by_quantity INT,
    rank_by_revenue INT,
    percentage_of_total_revenue DECIMAL(5,2),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT item_sales_statistics_menu_item_id_fkey 
        FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id) ON DELETE CASCADE,
    CONSTRAINT item_sales_statistics_restaurant_id_fkey 
        FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE,
    
    UNIQUE (menu_item_id, period_type, period_start)
);

CREATE INDEX idx_item_sales_statistics_menu_item ON public.item_sales_statistics(menu_item_id);
CREATE INDEX idx_item_sales_statistics_period ON public.item_sales_statistics(period_start, period_end);
CREATE INDEX idx_item_sales_statistics_rank_revenue ON public.item_sales_statistics(rank_by_revenue);

-- =====================================================
-- 4. CATEGORY SALES STATISTICS
-- =====================================================

CREATE TABLE IF NOT EXISTS public.category_sales_statistics (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    restaurant_id INT NOT NULL,
    
    -- Time period
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('daily', 'weekly', 'monthly', 'all_time')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Sales statistics
    total_quantity_sold INT DEFAULT 0 CHECK (total_quantity_sold >= 0),
    total_revenue DECIMAL(12,2) DEFAULT 0.00 CHECK (total_revenue >= 0),
    orders_count INT DEFAULT 0 CHECK (orders_count >= 0),
    percentage_of_total_revenue DECIMAL(5,2),
    average_order_value DECIMAL(10,2) DEFAULT 0.00 CHECK (average_order_value >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT category_sales_statistics_category_id_fkey 
        FOREIGN KEY (category_id) REFERENCES public.menu_categories(id) ON DELETE CASCADE,
    CONSTRAINT category_sales_statistics_restaurant_id_fkey 
        FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE,
    
    UNIQUE (category_id, period_type, period_start)
);

CREATE INDEX idx_category_sales_statistics_category ON public.category_sales_statistics(category_id);
CREATE INDEX idx_category_sales_statistics_period ON public.category_sales_statistics(period_start, period_end);

-- =====================================================
-- 5. DASHBOARD SUMMARY VIEW (Real-time aggregated data)
-- =====================================================

CREATE OR REPLACE VIEW public.dashboard_summary AS
SELECT 
    -- Today's revenue
    (SELECT COALESCE(SUM(total), 0) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed') AS revenue_today,
    
    -- Today's orders
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE) AS orders_today,
    
    -- Today's customers
    (SELECT COUNT(DISTINCT customer_user_id) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND customer_user_id IS NOT NULL) AS customers_today,
    
    -- Active orders
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed', 'preparing', 'ready')) AS active_orders,
    
    -- Pending orders (kitchen queue)
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed')) AS pending_orders,
    
    -- Preparing orders
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status = 'preparing') AS preparing_orders,
    
    -- Ready orders
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status = 'ready') AS ready_orders,
    
    -- Completed orders today
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed') AS completed_orders_today,
    
    -- Cancelled orders today
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND status = 'cancelled') AS cancelled_orders_today,
    
    -- Average order value today
    (SELECT COALESCE(AVG(total), 0) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed') AS average_order_value_today,
    
    -- Tables occupied
    (SELECT COUNT(DISTINCT table_id) 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed', 'preparing', 'ready', 'served')) AS tables_occupied,
    
    -- Total tables
    (SELECT COUNT(*) 
     FROM public.tables 
     WHERE status = 'active') AS total_tables,
    
    -- Growth rate (compare to yesterday)
    CASE 
        WHEN (SELECT COALESCE(SUM(total), 0) FROM public.orders WHERE DATE(created_at) = CURRENT_DATE - INTERVAL '1 day' AND status = 'completed') > 0 THEN
            ((SELECT COALESCE(SUM(total), 0) FROM public.orders WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed') - 
             (SELECT COALESCE(SUM(total), 0) FROM public.orders WHERE DATE(created_at) = CURRENT_DATE - INTERVAL '1 day' AND status = 'completed')) /
            (SELECT COALESCE(SUM(total), 0) FROM public.orders WHERE DATE(created_at) = CURRENT_DATE - INTERVAL '1 day' AND status = 'completed') * 100
        ELSE 0
    END AS revenue_growth_rate_vs_yesterday;

-- =====================================================
-- 6. KITCHEN DASHBOARD VIEW (For KDS)
-- =====================================================

CREATE OR REPLACE VIEW public.kitchen_dashboard AS
SELECT 
    -- Queue counts
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed')) AS pending_orders_count,
    
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status = 'preparing') AS preparing_orders_count,
    
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE status = 'ready') AS ready_orders_count,
    
    -- Today's kitchen stats
    (SELECT COUNT(*) 
     FROM public.orders 
     WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed') AS completed_orders_today,
    
    -- Average preparation time (in minutes)
    (SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (ready_at - preparing_at)) / 60), 0)
     FROM public.orders 
     WHERE ready_at IS NOT NULL AND preparing_at IS NOT NULL 
     AND DATE(created_at) = CURRENT_DATE) AS average_prep_time_minutes_today,
    
    -- Longest waiting order
    (SELECT order_number 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed', 'preparing')
     ORDER BY created_at ASC 
     LIMIT 1) AS longest_waiting_order_number,
    
    (SELECT EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - created_at)) / 60 
     FROM public.orders 
     WHERE status IN ('pending', 'confirmed', 'preparing')
     ORDER BY created_at ASC 
     LIMIT 1) AS longest_waiting_time_minutes;

-- =====================================================
-- 7. WAITER DASHBOARD VIEW
-- =====================================================

CREATE OR REPLACE VIEW public.waiter_dashboard AS
SELECT 
    u.id AS waiter_id,
    u.full_name AS waiter_name,
    
    -- Orders served today
    COUNT(DISTINCT o.id) FILTER (WHERE DATE(o.created_at) = CURRENT_DATE) AS orders_served_today,
    
    -- Revenue generated today
    COALESCE(SUM(o.total) FILTER (WHERE DATE(o.created_at) = CURRENT_DATE AND o.status = 'completed'), 0) AS revenue_generated_today,
    
    -- Active orders
    COUNT(DISTINCT o.id) FILTER (WHERE o.status IN ('pending', 'confirmed', 'preparing', 'ready', 'served')) AS active_orders,
    
    -- Assigned tables
    COUNT(DISTINCT wta.table_id) AS assigned_tables_count,
    
    -- Occupied tables
    COUNT(DISTINCT wta.table_id) FILTER (
        WHERE EXISTS (
            SELECT 1 FROM public.orders o2 
            WHERE o2.table_id = wta.table_id 
            AND o2.status IN ('pending', 'confirmed', 'preparing', 'ready', 'served')
        )
    ) AS occupied_tables_count
    
FROM public.users u
LEFT JOIN public.orders o ON u.id = o.waiter_id
LEFT JOIN public.waiter_table_assignments wta ON u.id = wta.waiter_id AND wta.is_active = TRUE
WHERE u.role = (SELECT id FROM public.roles WHERE name = 'waiter')
GROUP BY u.id, u.full_name;

-- =====================================================
-- TRIGGERS FOR AUTO-UPDATE TIMESTAMPS
-- =====================================================

DROP TRIGGER IF EXISTS update_daily_revenue_snapshots_updated_at ON public.daily_revenue_snapshots;
CREATE TRIGGER update_daily_revenue_snapshots_updated_at
    BEFORE UPDATE ON public.daily_revenue_snapshots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_item_sales_statistics_updated_at ON public.item_sales_statistics;
CREATE TRIGGER update_item_sales_statistics_updated_at
    BEFORE UPDATE ON public.item_sales_statistics
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_category_sales_statistics_updated_at ON public.category_sales_statistics;
CREATE TRIGGER update_category_sales_statistics_updated_at
    BEFORE UPDATE ON public.category_sales_statistics
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- FUNCTION: Calculate daily revenue snapshot (Run via cron job)
-- =====================================================

CREATE OR REPLACE FUNCTION calculate_daily_revenue_snapshot(target_date DATE, target_restaurant_id INT)
RETURNS VOID AS $$
DECLARE
    v_total_revenue DECIMAL(12,2);
    v_subtotal DECIMAL(12,2);
    v_tax_collected DECIMAL(12,2);
    v_discounts_given DECIMAL(12,2);
    v_service_charges DECIMAL(12,2);
    v_total_orders INT;
    v_completed_orders INT;
    v_cancelled_orders INT;
    v_average_order_value DECIMAL(10,2);
    v_total_customers INT;
    v_new_customers INT;
    v_returning_customers INT;
    v_cash_revenue DECIMAL(12,2);
    v_card_revenue DECIMAL(12,2);
    v_ewallet_revenue DECIMAL(12,2);
    v_stripe_revenue DECIMAL(12,2);
    v_tables_used INT;
BEGIN
    -- Calculate aggregates
    SELECT 
        COALESCE(SUM(o.total), 0),
        COALESCE(SUM(o.subtotal), 0),
        COALESCE(SUM(o.tax), 0),
        COALESCE(SUM(o.discount), 0),
        COALESCE(SUM(b.service_charge), 0),
        COUNT(*),
        COUNT(*) FILTER (WHERE o.status = 'completed'),
        COUNT(*) FILTER (WHERE o.status = 'cancelled'),
        COALESCE(AVG(o.total) FILTER (WHERE o.status = 'completed'), 0),
        COUNT(DISTINCT o.customer_user_id),
        COUNT(DISTINCT o.customer_user_id) FILTER (WHERE DATE(u.created_at) = target_date),
        COUNT(DISTINCT o.customer_user_id) FILTER (WHERE DATE(u.created_at) < target_date),
        COALESCE(SUM(p.amount) FILTER (WHERE p.method = 'cash' AND p.status = 'succeeded'), 0),
        COALESCE(SUM(p.amount) FILTER (WHERE p.method = 'card' AND p.status = 'succeeded'), 0),
        COALESCE(SUM(p.amount) FILTER (WHERE p.method = 'ewallet' AND p.status = 'succeeded'), 0),
        COALESCE(SUM(p.amount) FILTER (WHERE p.method = 'stripe' AND p.status = 'succeeded'), 0),
        COUNT(DISTINCT o.table_id)
    INTO 
        v_total_revenue, v_subtotal, v_tax_collected, v_discounts_given, v_service_charges,
        v_total_orders, v_completed_orders, v_cancelled_orders, v_average_order_value,
        v_total_customers, v_new_customers, v_returning_customers,
        v_cash_revenue, v_card_revenue, v_ewallet_revenue, v_stripe_revenue,
        v_tables_used
    FROM public.orders o
    LEFT JOIN public.bills b ON o.id = b.order_id
    LEFT JOIN public.payments p ON b.id = p.bill_id
    LEFT JOIN public.users u ON o.customer_user_id = u.id
    WHERE DATE(o.created_at) = target_date
    AND (target_restaurant_id IS NULL OR o.restaurant_id = target_restaurant_id);
    
    -- Insert or update snapshot
    INSERT INTO public.daily_revenue_snapshots (
        restaurant_id, snapshot_date,
        total_revenue, subtotal, tax_collected, discounts_given, service_charges,
        total_orders, completed_orders, cancelled_orders, average_order_value,
        total_customers, new_customers, returning_customers,
        cash_revenue, card_revenue, ewallet_revenue, stripe_revenue,
        tables_used
    ) VALUES (
        COALESCE(target_restaurant_id, 1), target_date,
        v_total_revenue, v_subtotal, v_tax_collected, v_discounts_given, v_service_charges,
        v_total_orders, v_completed_orders, v_cancelled_orders, v_average_order_value,
        v_total_customers, v_new_customers, v_returning_customers,
        v_cash_revenue, v_card_revenue, v_ewallet_revenue, v_stripe_revenue,
        v_tables_used
    )
    ON CONFLICT (restaurant_id, snapshot_date) 
    DO UPDATE SET
        total_revenue = EXCLUDED.total_revenue,
        subtotal = EXCLUDED.subtotal,
        tax_collected = EXCLUDED.tax_collected,
        discounts_given = EXCLUDED.discounts_given,
        service_charges = EXCLUDED.service_charges,
        total_orders = EXCLUDED.total_orders,
        completed_orders = EXCLUDED.completed_orders,
        cancelled_orders = EXCLUDED.cancelled_orders,
        average_order_value = EXCLUDED.average_order_value,
        total_customers = EXCLUDED.total_customers,
        new_customers = EXCLUDED.new_customers,
        returning_customers = EXCLUDED.returning_customers,
        cash_revenue = EXCLUDED.cash_revenue,
        card_revenue = EXCLUDED.card_revenue,
        ewallet_revenue = EXCLUDED.ewallet_revenue,
        stripe_revenue = EXCLUDED.stripe_revenue,
        tables_used = EXCLUDED.tables_used,
        updated_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- COMMENTS FOR DOCUMENTATION
-- =====================================================

COMMENT ON TABLE public.daily_revenue_snapshots IS 'Pre-calculated daily revenue snapshots for fast reporting';
COMMENT ON TABLE public.hourly_revenue_breakdown IS 'Hourly revenue breakdown for peak hours analysis';
COMMENT ON TABLE public.item_sales_statistics IS 'Menu item sales statistics aggregated by time period';
COMMENT ON TABLE public.category_sales_statistics IS 'Category sales statistics aggregated by time period';
COMMENT ON VIEW public.dashboard_summary IS 'Real-time dashboard summary for admin';
COMMENT ON VIEW public.kitchen_dashboard IS 'Real-time kitchen dashboard for KDS';
COMMENT ON VIEW public.waiter_dashboard IS 'Real-time waiter dashboard showing assigned orders and tables';

-- =====================================================
-- SEED DATA: Calculate today's snapshot
-- =====================================================

SELECT calculate_daily_revenue_snapshot(CURRENT_DATE, 1);
