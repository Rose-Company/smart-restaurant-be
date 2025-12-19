CREATE SEQUENCE IF NOT EXISTS tables_id_seq;

CREATE TABLE "public"."tables" (
    "id" int4 NOT NULL DEFAULT nextval('tables_id_seq'::regclass),
    "table_number" varchar(50) NOT NULL,
    "capacity" int4 NOT NULL CHECK ((capacity > 0) AND (capacity <= 20)),
    "location" varchar(100),
    "description" text,
    "status" text NOT NULL DEFAULT 'active'::text,
    "qr_token" text,
    "qr_token_created_at" timestamp,
    "qr_token_expires_at" timestamp,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Indices
CREATE UNIQUE INDEX tables_table_number_key ON public.tables USING btree (table_number);

-- Sequence and defined type for orders
CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

-- Table Definition for orders
CREATE TABLE "public"."orders" (
    "id" int4 NOT NULL DEFAULT nextval('orders_id_seq'::regclass),
    "table_id" int4 NOT NULL,
    "order_number" varchar(50) NOT NULL,
    "status" varchar(20) NOT NULL DEFAULT 'pending'::character varying,
    "customer_user_id" uuid,
    "subtotal" numeric(10,2) NOT NULL DEFAULT 0.00,
    "tax" numeric(10,2) NOT NULL DEFAULT 0.00,
    "discount" numeric(10,2) NOT NULL DEFAULT 0.00,
    "total" numeric(10,2) NOT NULL DEFAULT 0.00 CHECK (total >= (0)::numeric),
    "notes" text,
    "special_instructions" text,
    "meta" jsonb,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "accepted_at" timestamp,
    "preparing_at" timestamp,
    "ready_at" timestamp,
    "served_at" timestamp,
    "completed_at" timestamp,
    CONSTRAINT "orders_table_id_fkey" FOREIGN KEY ("table_id") REFERENCES "public"."tables"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- Indices
CREATE UNIQUE INDEX orders_order_number_key ON public.orders USING btree (order_number);

-- Sequence and defined type for order_items
CREATE SEQUENCE IF NOT EXISTS order_items_id_seq;

-- Table Definition for order_items
CREATE TABLE "public"."order_items" (
    "id" int4 NOT NULL DEFAULT nextval('order_items_id_seq'::regclass),
    "order_id" int4 NOT NULL,
    "menu_item_id" int4,
    "item_name" varchar(255) NOT NULL,
    "item_description" text,
    "quantity" int4 NOT NULL DEFAULT 1 CHECK (quantity > 0),
    "unit_price" numeric(10,2) NOT NULL CHECK (unit_price >= (0)::numeric),
    "subtotal" numeric(10,2) NOT NULL CHECK (subtotal >= (0)::numeric),
    "status" varchar(20) NOT NULL DEFAULT 'pending'::character varying,
    "special_instructions" text,
    "meta" jsonb,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "order_items_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "public"."orders"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- Insert sample tables data
INSERT INTO public.tables (table_number, capacity, location, status) VALUES
('T-01', 4, 'Main Hall', 'active'),
('T-02', 2, 'Main Hall', 'occupied'),
('T-03', 6, 'Main Hall', 'active'),
('T-04', 4, 'Main Hall', 'active'),
('T-05', 8, 'VIP', 'active'),
('T-06', 4, 'Patio', 'active'),
('T-07', 2, 'Patio', 'inactive'),
('VIP-01', 10, 'VIP', 'active'),
('VIP-02', 6, 'VIP', 'occupied'),
('P-01', 4, 'Patio', 'active'),
('P-02', 4, 'Patio', 'active'),
('P-03', 2, 'Patio', 'active')
ON CONFLICT (table_number) DO NOTHING;

-- Insert sample orders data for occupied tables
-- Orders for T-02 (occupied)
INSERT INTO public.orders (table_id, order_number, status, total) VALUES
(2, 'ORD-2025-001', 'pending', 156.00),
(2, 'ORD-2025-002', 'processing', 89.50)
ON CONFLICT (order_number) DO NOTHING;

-- Orders for VIP-02 (occupied)
INSERT INTO public.orders (table_id, order_number, status, total) VALUES
(9, 'ORD-2025-003', 'processing', 245.75)
ON CONFLICT (order_number) DO NOTHING;

-- Insert sample order items
INSERT INTO public.order_items (order_id, item_name, quantity, unit_price, subtotal, status) VALUES
(1, 'Grilled Salmon', 2, 45.00, 90.00, 'pending'),
(1, 'Caesar Salad', 1, 12.00, 12.00, 'pending'),
(1, 'Soft Drink', 2, 5.00, 10.00, 'pending'),
(2, 'Beef Steak', 1, 65.00, 65.00, 'processing'),
(2, 'French Fries', 1, 8.50, 8.50, 'processing'),
(3, 'Lobster Pasta', 2, 55.00, 110.00, 'processing'),
(3, 'Garlic Bread', 2, 6.00, 12.00, 'processing'),
(3, 'Red Wine', 1, 30.00, 30.00, 'processing')
ON CONFLICT DO NOTHING;
