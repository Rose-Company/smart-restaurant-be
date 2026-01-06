-- =====================================================
-- AUTHENTICATION TABLES
-- =====================================================

-- Create Roles table
CREATE TABLE IF NOT EXISTS public.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Users table
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone_number VARCHAR(20),
    role UUID NOT NULL,
    status VARCHAR(20) DEFAULT 'active' ,
    is_active BOOLEAN DEFAULT TRUE,
    provider VARCHAR(50) DEFAULT 'local',
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_fkey FOREIGN KEY (role) REFERENCES public.roles(id) ON DELETE RESTRICT
);

-- Create OTP table
CREATE TABLE IF NOT EXISTS public.otps (
    id SERIAL PRIMARY KEY,
    target VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verify_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create OTP Attempts table
CREATE TABLE IF NOT EXISTS public.otp_attempts (
    id SERIAL PRIMARY KEY,
    otp_id INT NOT NULL,
    value VARCHAR(6) NOT NULL,
    is_success BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT otp_attempts_otp_id_fkey FOREIGN KEY (otp_id) REFERENCES public.otps(id) ON DELETE CASCADE
);


-- =====================================================
-- INSERT DEFAULT ROLES
-- =====================================================

INSERT INTO public.roles (name, description) VALUES
('end_user', 'Regular user/customer'),
('admin', 'Administrator with full access'),
('restaurant_owner', 'Restaurant owner account'),
('staff', 'Restaurant staff member')
ON CONFLICT (name) DO NOTHING;
