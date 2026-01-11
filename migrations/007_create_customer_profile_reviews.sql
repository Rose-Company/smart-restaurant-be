-- =====================================================
-- 007_create_customer_profile_reviews.sql
-- CUSTOMER PROFILE & REVIEWS SYSTEM
-- =====================================================
-- This migration creates all tables needed for:
-- - Customer Profile Management (TASK-016 to TASK-019)
-- - Customer Reviews (TASK-020, TASK-009)
-- - Customer Preferences & Settings
-- =====================================================

-- =====================================================
-- 1. ENHANCE EXISTING USERS TABLE FOR CUSTOMERS
-- =====================================================

-- Add customer-specific fields to users table
ALTER TABLE public.users
ADD COLUMN IF NOT EXISTS full_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS avatar_url TEXT,
ADD COLUMN IF NOT EXISTS date_of_birth DATE,
ADD COLUMN IF NOT EXISTS gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other', 'prefer_not_to_say')),
ADD COLUMN IF NOT EXISTS street_address VARCHAR(255),
ADD COLUMN IF NOT EXISTS city VARCHAR(100),
ADD COLUMN IF NOT EXISTS state VARCHAR(100),
ADD COLUMN IF NOT EXISTS postal_code VARCHAR(20),
ADD COLUMN IF NOT EXISTS country VARCHAR(100) DEFAULT 'USA',
ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS profile_completed BOOLEAN DEFAULT FALSE;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_email ON public.users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON public.users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_role ON public.users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON public.users(status);
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON public.users(email_verified);

-- =====================================================
-- 2. CUSTOMER PREFERENCES TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.customer_preferences (
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL UNIQUE,
    
    -- Dietary restrictions
    dietary_restrictions TEXT[], -- Array of: 'no_nuts', 'no_shellfish', 'vegetarian', 'vegan', 'gluten_free', 'dairy_free', 'halal', 'kosher'
    allergens TEXT[], -- Array of specific allergens
    
    -- Food preferences
    favorite_cuisine TEXT[], -- Array of: 'Italian', 'Japanese', 'American', 'Mexican', 'Chinese', 'Thai', etc.
    spice_level VARCHAR(20) CHECK (spice_level IN ('none', 'mild', 'medium', 'hot', 'extra_hot')),
    
    -- Notification preferences
    notification_enabled BOOLEAN DEFAULT TRUE,
    email_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    push_notifications BOOLEAN DEFAULT TRUE,
    marketing_emails BOOLEAN DEFAULT FALSE,
    order_updates BOOLEAN DEFAULT TRUE,
    promotional_offers BOOLEAN DEFAULT TRUE,
    
    -- Display preferences
    language VARCHAR(10) DEFAULT 'en',
    currency VARCHAR(10) DEFAULT 'USD',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT customer_preferences_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_customer_preferences_customer_id ON public.customer_preferences(customer_id);

-- =====================================================
-- 3. CUSTOMER LOYALTY POINTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.customer_loyalty (
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL UNIQUE,
    
    -- Points system
    points INT DEFAULT 0 CHECK (points >= 0),
    lifetime_points INT DEFAULT 0 CHECK (lifetime_points >= 0),
    tier VARCHAR(20) DEFAULT 'bronze' CHECK (tier IN ('bronze', 'silver', 'gold', 'platinum', 'diamond')),
    
    -- Tier thresholds (points needed)
    -- bronze: 0-99, silver: 100-499, gold: 500-999, platinum: 1000-1999, diamond: 2000+
    tier_expiry_date DATE,
    
    -- Rewards
    available_rewards INT DEFAULT 0 CHECK (available_rewards >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT customer_loyalty_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_customer_loyalty_customer_id ON public.customer_loyalty(customer_id);
CREATE INDEX idx_customer_loyalty_tier ON public.customer_loyalty(tier);

-- =====================================================
-- 4. LOYALTY POINTS TRANSACTIONS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.loyalty_transactions (
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL,
    order_id INT,
    
    -- Transaction details
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('earned', 'redeemed', 'expired', 'adjusted', 'bonus')),
    points INT NOT NULL, -- Positive for earned/bonus, negative for redeemed/expired
    description TEXT,
    
    -- Balance tracking
    balance_before INT NOT NULL CHECK (balance_before >= 0),
    balance_after INT NOT NULL CHECK (balance_after >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at DATE, -- Points expiry (e.g., 1 year from earned)
    
    CONSTRAINT loyalty_transactions_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT loyalty_transactions_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE SET NULL
);

CREATE INDEX idx_loyalty_transactions_customer_id ON public.loyalty_transactions(customer_id);
CREATE INDEX idx_loyalty_transactions_order_id ON public.loyalty_transactions(order_id);
CREATE INDEX idx_loyalty_transactions_created_at ON public.loyalty_transactions(created_at DESC);

-- =====================================================
-- 5. CUSTOMER REVIEWS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.customer_reviews (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    customer_id UUID NOT NULL,
    
    -- Review content
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    
    -- Review details
    food_quality_rating INT CHECK (food_quality_rating >= 1 AND food_quality_rating <= 5),
    service_rating INT CHECK (service_rating >= 1 AND service_rating <= 5),
    ambiance_rating INT CHECK (ambiance_rating >= 1 AND ambiance_rating <= 5),
    value_rating INT CHECK (value_rating >= 1 AND value_rating <= 5),
    
    -- Restaurant response
    restaurant_response TEXT,
    responded_by UUID,
    responded_by_name VARCHAR(255),
    responded_at TIMESTAMP,
    
    -- Status & visibility
    status VARCHAR(20) DEFAULT 'published' CHECK (status IN ('pending', 'published', 'hidden', 'flagged', 'deleted')),
    is_verified_purchase BOOLEAN DEFAULT TRUE,
    helpful_count INT DEFAULT 0 CHECK (helpful_count >= 0),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT customer_reviews_order_id_fkey 
        FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
    CONSTRAINT customer_reviews_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT customer_reviews_responded_by_fkey 
        FOREIGN KEY (responded_by) REFERENCES public.users(id) ON DELETE SET NULL,
    
    -- Each order can only be reviewed once per customer
    UNIQUE (order_id, customer_id)
);

CREATE INDEX idx_customer_reviews_order_id ON public.customer_reviews(order_id);
CREATE INDEX idx_customer_reviews_customer_id ON public.customer_reviews(customer_id);
CREATE INDEX idx_customer_reviews_rating ON public.customer_reviews(rating);
CREATE INDEX idx_customer_reviews_status ON public.customer_reviews(status);
CREATE INDEX idx_customer_reviews_created_at ON public.customer_reviews(created_at DESC);

-- =====================================================
-- 6. REVIEW PHOTOS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.review_photos (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_size INT, -- in bytes
    mime_type VARCHAR(50),
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT review_photos_review_id_fkey 
        FOREIGN KEY (review_id) REFERENCES public.customer_reviews(id) ON DELETE CASCADE
);

CREATE INDEX idx_review_photos_review_id ON public.review_photos(review_id);

-- =====================================================
-- 7. REVIEW ITEMS (Which menu items were reviewed)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.review_items (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    item_rating INT CHECK (item_rating >= 1 AND item_rating <= 5),
    item_comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT review_items_review_id_fkey 
        FOREIGN KEY (review_id) REFERENCES public.customer_reviews(id) ON DELETE CASCADE,
    CONSTRAINT review_items_menu_item_id_fkey 
        FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id) ON DELETE CASCADE
);

CREATE INDEX idx_review_items_review_id ON public.review_items(review_id);
CREATE INDEX idx_review_items_menu_item_id ON public.review_items(menu_item_id);

-- =====================================================
-- 8. FAVORITE ITEMS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.customer_favorite_items (
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL,
    menu_item_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT customer_favorite_items_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT customer_favorite_items_menu_item_id_fkey 
        FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id) ON DELETE CASCADE,
    
    UNIQUE (customer_id, menu_item_id)
);

CREATE INDEX idx_customer_favorite_items_customer_id ON public.customer_favorite_items(customer_id);
CREATE INDEX idx_customer_favorite_items_menu_item_id ON public.customer_favorite_items(menu_item_id);

-- =====================================================
-- 9. CUSTOMER STATISTICS VIEW (For fast queries)
-- =====================================================

CREATE OR REPLACE VIEW public.customer_statistics AS
SELECT 
    u.id AS customer_id,
    u.email,
    u.full_name,
    u.created_at AS customer_since,
    
    -- Order statistics
    COUNT(DISTINCT o.id) AS total_orders,
    COALESCE(SUM(o.total), 0) AS total_spent,
    COALESCE(AVG(o.total), 0) AS average_order_value,
    MAX(o.created_at) AS last_order_date,
    
    -- Review statistics
    COUNT(DISTINCT cr.id) AS total_reviews,
    COALESCE(AVG(cr.rating), 0) AS average_rating_given,
    
    -- Loyalty
    COALESCE(cl.points, 0) AS loyalty_points,
    COALESCE(cl.tier, 'bronze') AS loyalty_tier,
    
    -- Favorite items (top 3)
    ARRAY_AGG(DISTINCT mi.name ORDER BY mi.name LIMIT 3) FILTER (WHERE mi.name IS NOT NULL) AS favorite_items
    
FROM public.users u
LEFT JOIN public.orders o ON u.id = o.customer_user_id
LEFT JOIN public.customer_reviews cr ON u.id = cr.customer_id
LEFT JOIN public.customer_loyalty cl ON u.id = cl.customer_id
LEFT JOIN public.customer_favorite_items cfi ON u.id = cfi.customer_id
LEFT JOIN public.menu_items mi ON cfi.menu_item_id = mi.id
WHERE u.role = (SELECT id FROM public.roles WHERE name = 'end_user')
GROUP BY u.id, u.email, u.full_name, u.created_at, cl.points, cl.tier;

-- =====================================================
-- TRIGGERS FOR AUTO-UPDATE TIMESTAMPS
-- =====================================================

-- Apply trigger to customer_preferences
DROP TRIGGER IF EXISTS update_customer_preferences_updated_at ON public.customer_preferences;
CREATE TRIGGER update_customer_preferences_updated_at
    BEFORE UPDATE ON public.customer_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to customer_loyalty
DROP TRIGGER IF EXISTS update_customer_loyalty_updated_at ON public.customer_loyalty;
CREATE TRIGGER update_customer_loyalty_updated_at
    BEFORE UPDATE ON public.customer_loyalty
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to customer_reviews
DROP TRIGGER IF EXISTS update_customer_reviews_updated_at ON public.customer_reviews;
CREATE TRIGGER update_customer_reviews_updated_at
    BEFORE UPDATE ON public.customer_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to users
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- FUNCTION: Auto-create customer preferences & loyalty on user registration
-- =====================================================

CREATE OR REPLACE FUNCTION create_customer_profile()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create profile for end_user role
    IF NEW.role = (SELECT id FROM public.roles WHERE name = 'end_user') THEN
        -- Create customer preferences
        INSERT INTO public.customer_preferences (customer_id)
        VALUES (NEW.id)
        ON CONFLICT (customer_id) DO NOTHING;
        
        -- Create customer loyalty
        INSERT INTO public.customer_loyalty (customer_id)
        VALUES (NEW.id)
        ON CONFLICT (customer_id) DO NOTHING;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-create profile on user insert
DROP TRIGGER IF EXISTS create_customer_profile_trigger ON public.users;
CREATE TRIGGER create_customer_profile_trigger
    AFTER INSERT ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION create_customer_profile();

-- =====================================================
-- FUNCTION: Calculate loyalty tier based on points
-- =====================================================

CREATE OR REPLACE FUNCTION update_loyalty_tier()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.points >= 2000 THEN
        NEW.tier = 'diamond';
    ELSIF NEW.points >= 1000 THEN
        NEW.tier = 'platinum';
    ELSIF NEW.points >= 500 THEN
        NEW.tier = 'gold';
    ELSIF NEW.points >= 100 THEN
        NEW.tier = 'silver';
    ELSE
        NEW.tier = 'bronze';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update tier on points change
DROP TRIGGER IF EXISTS update_loyalty_tier_trigger ON public.customer_loyalty;
CREATE TRIGGER update_loyalty_tier_trigger
    BEFORE UPDATE OF points ON public.customer_loyalty
    FOR EACH ROW
    EXECUTE FUNCTION update_loyalty_tier();

-- =====================================================
-- FUNCTION: Award loyalty points on order completion
-- =====================================================

CREATE OR REPLACE FUNCTION award_loyalty_points()
RETURNS TRIGGER AS $$
DECLARE
    points_earned INT;
    current_balance INT;
BEGIN
    -- Only award points when order is completed
    IF NEW.status = 'completed' AND OLD.status != 'completed' AND NEW.customer_user_id IS NOT NULL THEN
        -- Calculate points: 1 point per $1 spent (rounded down)
        points_earned := FLOOR(NEW.total);
        
        -- Get current balance
        SELECT points INTO current_balance
        FROM public.customer_loyalty
        WHERE customer_id = NEW.customer_user_id;
        
        IF current_balance IS NULL THEN
            current_balance := 0;
        END IF;
        
        -- Update loyalty points
        UPDATE public.customer_loyalty
        SET 
            points = points + points_earned,
            lifetime_points = lifetime_points + points_earned
        WHERE customer_id = NEW.customer_user_id;
        
        -- Record transaction
        INSERT INTO public.loyalty_transactions (
            customer_id, 
            order_id, 
            transaction_type, 
            points, 
            description,
            balance_before,
            balance_after,
            expires_at
        ) VALUES (
            NEW.customer_user_id,
            NEW.id,
            'earned',
            points_earned,
            'Points earned from order ' || NEW.order_number,
            current_balance,
            current_balance + points_earned,
            CURRENT_DATE + INTERVAL '1 year'
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to award points on order completion
DROP TRIGGER IF EXISTS award_loyalty_points_trigger ON public.orders;
CREATE TRIGGER award_loyalty_points_trigger
    AFTER UPDATE ON public.orders
    FOR EACH ROW
    EXECUTE FUNCTION award_loyalty_points();

-- =====================================================
-- COMMENTS FOR DOCUMENTATION
-- =====================================================

COMMENT ON TABLE public.customer_preferences IS 'Customer preferences including dietary restrictions, notifications, and display settings';
COMMENT ON TABLE public.customer_loyalty IS 'Customer loyalty points and tier information';
COMMENT ON TABLE public.loyalty_transactions IS 'Audit trail of all loyalty points transactions';
COMMENT ON TABLE public.customer_reviews IS 'Customer reviews and ratings for orders';
COMMENT ON TABLE public.review_photos IS 'Photos uploaded with customer reviews';
COMMENT ON TABLE public.review_items IS 'Individual menu item ratings within a review';
COMMENT ON TABLE public.customer_favorite_items IS 'Customer favorite menu items for quick reordering';
COMMENT ON VIEW public.customer_statistics IS 'Aggregated customer statistics for analytics and profile display';
