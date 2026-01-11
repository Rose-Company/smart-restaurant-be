-- =====================================================
-- 008_create_staff_management.sql
-- STAFF MANAGEMENT & PERMISSIONS SYSTEM
-- =====================================================
-- This migration creates all tables needed for:
-- - Staff Account Management (TASK-021 to TASK-030)
-- - Role-based Access Control (Waiter, Kitchen, Admin)
-- - Staff Assignments & Tracking
-- =====================================================

-- =====================================================
-- 1. ENHANCE EXISTING ROLES TABLE
-- =====================================================

-- Add more detailed role configuration
ALTER TABLE public.roles
ADD COLUMN IF NOT EXISTS display_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS permissions JSONB,
ADD COLUMN IF NOT EXISTS is_staff BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;

-- Update existing roles
UPDATE public.roles SET display_name = 'End User', is_staff = FALSE WHERE name = 'end_user';
UPDATE public.roles SET display_name = 'Administrator', is_staff = TRUE WHERE name = 'admin';
UPDATE public.roles SET display_name = 'Restaurant Owner', is_staff = TRUE WHERE name = 'restaurant_owner';
UPDATE public.roles SET display_name = 'Staff Member', is_staff = TRUE WHERE name = 'staff';

-- Add new staff roles
INSERT INTO public.roles (name, display_name, description, is_staff, is_active) VALUES
('waiter', 'Waiter', 'Waiter/Server staff - handles table service and orders', TRUE, TRUE),
('kitchen_staff', 'Kitchen Staff', 'Kitchen staff - prepares orders', TRUE, TRUE),
('manager', 'Manager', 'Restaurant manager with elevated permissions', TRUE, TRUE),
('cashier', 'Cashier', 'Cashier - handles payments and billing', TRUE, TRUE)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    description = EXCLUDED.description,
    is_staff = EXCLUDED.is_staff;

-- =====================================================
-- 2. STAFF PROFILES TABLE (Extended user info for staff)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.staff_profiles (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    
    -- Employment info
    employee_id VARCHAR(50) UNIQUE,
    department VARCHAR(100),
    position VARCHAR(100),
    hire_date DATE,
    employment_status VARCHAR(20) DEFAULT 'active' CHECK (employment_status IN ('active', 'inactive', 'on_leave', 'terminated')),
    
    -- Work schedule
    shift_type VARCHAR(20) CHECK (shift_type IN ('morning', 'afternoon', 'evening', 'night', 'rotating', 'flexible')),
    weekly_hours INT CHECK (weekly_hours >= 0 AND weekly_hours <= 168),
    
    -- Emergency contact
    emergency_contact_name VARCHAR(255),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relationship VARCHAR(50),
    
    -- Performance tracking
    total_orders_served INT DEFAULT 0 CHECK (total_orders_served >= 0),
    total_orders_prepared INT DEFAULT 0 CHECK (total_orders_prepared >= 0),
    average_order_time INT, -- in minutes
    average_rating DECIMAL(3,2) CHECK (average_rating >= 0 AND average_rating <= 5),
    
    -- Notes
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT staff_profiles_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_staff_profiles_user_id ON public.staff_profiles(user_id);
CREATE INDEX idx_staff_profiles_employee_id ON public.staff_profiles(employee_id);
CREATE INDEX idx_staff_profiles_employment_status ON public.staff_profiles(employment_status);

-- =====================================================
-- 3. WAITER TABLE ASSIGNMENTS
-- =====================================================

CREATE TABLE IF NOT EXISTS public.waiter_table_assignments (
    id SERIAL PRIMARY KEY,
    waiter_id UUID NOT NULL,
    table_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID,
    is_active BOOLEAN DEFAULT TRUE,
    
    CONSTRAINT waiter_table_assignments_waiter_id_fkey 
        FOREIGN KEY (waiter_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT waiter_table_assignments_table_id_fkey 
        FOREIGN KEY (table_id) REFERENCES public.tables(id) ON DELETE CASCADE,
    CONSTRAINT waiter_table_assignments_assigned_by_fkey 
        FOREIGN KEY (assigned_by) REFERENCES public.users(id) ON DELETE SET NULL,
    
    UNIQUE (waiter_id, table_id)
);

CREATE INDEX idx_waiter_table_assignments_waiter_id ON public.waiter_table_assignments(waiter_id);
CREATE INDEX idx_waiter_table_assignments_table_id ON public.waiter_table_assignments(table_id);
CREATE INDEX idx_waiter_table_assignments_is_active ON public.waiter_table_assignments(is_active);

-- =====================================================
-- 4. STAFF PERMISSIONS
-- =====================================================
-- Note: Using existing action_control_list table for permissions management
-- No additional permissions tables needed

-- =====================================================
-- 5. ROLE PERMISSIONS
-- =====================================================
-- Note: Using existing action_control_list table for role-permission mapping
-- No additional role_permissions table needed

-- =====================================================
-- 4. STAFF INVITATIONS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS public.staff_invitations (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    role_id UUID NOT NULL,
    invited_by UUID NOT NULL,
    
    -- Invitation details
    token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'expired', 'cancelled')),
    expires_at TIMESTAMP NOT NULL,
    
    -- Acceptance tracking
    accepted_at TIMESTAMP,
    accepted_by UUID,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT staff_invitations_role_id_fkey 
        FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE,
    CONSTRAINT staff_invitations_invited_by_fkey 
        FOREIGN KEY (invited_by) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT staff_invitations_accepted_by_fkey 
        FOREIGN KEY (accepted_by) REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX idx_staff_invitations_email ON public.staff_invitations(email);
CREATE INDEX idx_staff_invitations_token ON public.staff_invitations(token);
CREATE INDEX idx_staff_invitations_status ON public.staff_invitations(status);

-- =====================================================
-- 5. STAFF ACTIVITY LOG (Audit trail)
-- =====================================================

CREATE TABLE IF NOT EXISTS public.staff_activity_log (
    id SERIAL PRIMARY KEY,
    staff_id UUID NOT NULL,
    
    -- Activity details
    action VARCHAR(100) NOT NULL, -- 'login', 'logout', 'create_order', 'update_order', 'process_payment', etc.
    action_category VARCHAR(50), -- 'auth', 'orders', 'billing', 'menu', 'staff'
    description TEXT,
    
    -- Context
    entity_type VARCHAR(50), -- 'order', 'bill', 'menu_item', 'staff', etc.
    entity_id INT,
    
    -- Request metadata
    ip_address INET,
    user_agent TEXT,
    
    -- Result
    status VARCHAR(20) CHECK (status IN ('success', 'failed', 'error')),
    error_message TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT staff_activity_log_staff_id_fkey 
        FOREIGN KEY (staff_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_staff_activity_log_staff_id ON public.staff_activity_log(staff_id);
CREATE INDEX idx_staff_activity_log_action ON public.staff_activity_log(action);
CREATE INDEX idx_staff_activity_log_created_at ON public.staff_activity_log(created_at DESC);
CREATE INDEX idx_staff_activity_log_entity ON public.staff_activity_log(entity_type, entity_id);

-- =====================================================
-- 6. STAFF STATISTICS VIEW
-- =====================================================

CREATE OR REPLACE VIEW public.staff_statistics AS
SELECT 
    u.id AS staff_id,
    u.email,
    u.full_name,
    r.name AS role_name,
    r.display_name AS role_display_name,
    sp.employee_id,
    sp.employment_status,
    sp.hire_date,
    
    -- Order statistics (for waiters)
    COUNT(DISTINCT o_waiter.id) AS total_orders_served,
    COALESCE(SUM(o_waiter.total), 0) AS revenue_generated,
    COALESCE(AVG(o_waiter.total), 0) AS average_order_value,
    
    -- Kitchen statistics (for kitchen staff)
    COUNT(DISTINCT o_kitchen.id) AS total_orders_prepared,
    
    -- Table assignments (for waiters)
    COUNT(DISTINCT wta.table_id) FILTER (WHERE wta.is_active = TRUE) AS assigned_tables_count,
    ARRAY_AGG(DISTINCT t.table_number ORDER BY t.table_number) FILTER (WHERE wta.is_active = TRUE) AS assigned_table_numbers,
    
    -- Performance
    sp.average_rating,
    
    -- Last activity
    MAX(sal.created_at) AS last_activity_at,
    u.last_login_at
    
FROM public.users u
INNER JOIN public.roles r ON u.role = r.id
LEFT JOIN public.staff_profiles sp ON u.id = sp.user_id
LEFT JOIN public.orders o_waiter ON u.id = o_waiter.waiter_id
LEFT JOIN public.orders o_kitchen ON u.id = o_kitchen.kitchen_staff_id
LEFT JOIN public.waiter_table_assignments wta ON u.id = wta.waiter_id
LEFT JOIN public.tables t ON wta.table_id = t.id
LEFT JOIN public.staff_activity_log sal ON u.id = sal.staff_id
WHERE r.is_staff = TRUE
GROUP BY u.id, u.email, u.full_name, r.name, r.display_name, 
         sp.employee_id, sp.employment_status, sp.hire_date, sp.average_rating, u.last_login_at;

-- =====================================================
-- TRIGGERS FOR AUTO-UPDATE TIMESTAMPS
-- =====================================================

-- Apply trigger to staff_profiles
DROP TRIGGER IF EXISTS update_staff_profiles_updated_at ON public.staff_profiles;
CREATE TRIGGER update_staff_profiles_updated_at
    BEFORE UPDATE ON public.staff_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- FUNCTION: Auto-create staff profile on staff user creation
-- =====================================================

CREATE OR REPLACE FUNCTION create_staff_profile()
RETURNS TRIGGER AS $$
DECLARE
    role_name_var VARCHAR(50);
BEGIN
    -- Get role name
    SELECT r.name INTO role_name_var
    FROM public.roles r
    WHERE r.id = NEW.role;
    
    -- Only create staff profile for staff roles
    IF role_name_var IN ('admin', 'waiter', 'kitchen_staff', 'manager', 'cashier', 'staff', 'restaurant_owner') THEN
        INSERT INTO public.staff_profiles (user_id, hire_date, employment_status)
        VALUES (NEW.id, CURRENT_DATE, 'active')
        ON CONFLICT (user_id) DO NOTHING;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-create staff profile
DROP TRIGGER IF EXISTS create_staff_profile_trigger ON public.users;
CREATE TRIGGER create_staff_profile_trigger
    AFTER INSERT ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION create_staff_profile();

-- =====================================================
-- FUNCTION: Log staff login activity
-- =====================================================

CREATE OR REPLACE FUNCTION log_staff_login()
RETURNS TRIGGER AS $$
BEGIN
    -- Only log for staff users
    IF EXISTS (SELECT 1 FROM public.roles WHERE id = NEW.role AND is_staff = TRUE) THEN
        INSERT INTO public.staff_activity_log (
            staff_id, 
            action, 
            action_category, 
            description, 
            status
        ) VALUES (
            NEW.id,
            'login',
            'auth',
            'Staff member logged in',
            'success'
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to log login
DROP TRIGGER IF EXISTS log_staff_login_trigger ON public.users;
CREATE TRIGGER log_staff_login_trigger
    AFTER UPDATE OF last_login_at ON public.users
    FOR EACH ROW
    WHEN (NEW.last_login_at IS DISTINCT FROM OLD.last_login_at)
    EXECUTE FUNCTION log_staff_login();

-- =====================================================
-- COMMENTS FOR DOCUMENTATION
-- =====================================================

COMMENT ON TABLE public.staff_profiles IS 'Extended staff information including employment details and performance metrics';
COMMENT ON TABLE public.waiter_table_assignments IS 'Tracks which tables are assigned to which waiters';
COMMENT ON TABLE public.staff_invitations IS 'Email invitations sent to new staff members';
COMMENT ON TABLE public.staff_activity_log IS 'Audit trail of all staff actions in the system';
COMMENT ON VIEW public.staff_statistics IS 'Aggregated staff statistics for performance tracking and management';
