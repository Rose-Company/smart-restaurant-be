-- =====================================================
-- RESTAURANT TABLE
-- =====================================================

CREATE SEQUENCE IF NOT EXISTS restaurants_id_seq;

CREATE TABLE "public"."restaurants" (
    "id" int4 NOT NULL DEFAULT nextval('restaurants_id_seq'::regclass),
    "name" varchar(255) NOT NULL,
    "description" text,
    "address" text,
    "phone" varchar(20),
    "email" varchar(255),
    "logo_url" text,
    "status" varchar(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- =====================================================
-- ADD restaurant_id TO EXISTING TABLES
-- =====================================================

-- Add restaurant_id to tables
ALTER TABLE "public"."tables" 
ADD COLUMN "restaurant_id" int4 


CREATE INDEX idx_tables_restaurant ON public.tables(restaurant_id);

-- Add restaurant_id to orders
ALTER TABLE "public"."orders" 
ADD COLUMN "restaurant_id" int4 

CREATE INDEX idx_orders_restaurant ON public.orders(restaurant_id);


-- =====================================================
-- SEED DATA FOR RESTAURANT
-- =====================================================

INSERT INTO public.restaurants (name, description, address, phone, email, status) VALUES
('Smart Restaurant', 'Modern dining experience with QR ordering', '123 Main Street, City Center', '+1234567890', 'info@smartrestaurant.com', 'active')
ON CONFLICT DO NOTHING;
