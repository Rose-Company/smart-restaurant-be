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



-- Insert sample data
INSERT INTO public.tables (table_number, capacity, location, status) VALUES
('T-01', 4, 'Main Hall', 'Active'),
('T-02', 2, 'Main Hall', 'Occupied'),
('T-03', 6, 'Main Hall', 'Active'),
('T-04', 4, 'Main Hall', 'Active'),
('T-05', 8, 'VIP', 'Active'),
('T-06', 4, 'Patio', 'Active'),
('T-07', 2, 'Patio', 'Inactive'),
('VIP-01', 10, 'VIP', 'Active'),
('VIP-02', 6, 'VIP', 'Occupied'),
('P-01', 4, 'Patio', 'Active'),
('P-02', 4, 'Patio', 'Active'),
('P-03', 2, 'Patio', 'Active')
ON CONFLICT (table_number) DO NOTHING;
