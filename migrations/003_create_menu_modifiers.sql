-- Categories
CREATE TABLE menu_categories (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    display_order INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (restaurant_id, name)
);

-- Items
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    category_id INT NOT NULL,
    name VARCHAR(80) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL CHECK (price > 0),
    prep_time_minutes INT DEFAULT 0 CHECK (prep_time_minutes >= 0 AND prep_time_minutes <= 240),
    status VARCHAR(20) NOT NULL CHECK (status IN ('available', 'unavailable', 'sold_out')),
    is_chef_recommended BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Item photos
CREATE TABLE menu_item_photos (
    id SERIAL PRIMARY KEY,
    menu_item_id INT NOT NULL,
    url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Modifier groups
CREATE TABLE modifier_groups (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    name VARCHAR(80) NOT NULL,
    selection_type VARCHAR(20) NOT NULL CHECK (selection_type IN ('single', 'multiple')),
    is_required BOOLEAN DEFAULT FALSE,
    min_selections INT DEFAULT 0,
    max_selections INT DEFAULT 0,
    display_order INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Modifier options
CREATE TABLE modifier_options (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    name VARCHAR(80) NOT NULL,
    price_adjustment DECIMAL(12,2) DEFAULT 0 CHECK (price_adjustment >= 0),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Attach groups to items 
CREATE TABLE menu_item_modifier_groups (
    id SERIAL PRIMARY KEY,
    menu_item_id INT NOT NULL,
    group_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (menu_item_id, group_id)
);



-- =====================================================
-- ADD restaurant_id TO MENU TABLES (haven't run it yet)
-- =====================================================

-- Already has restaurant_id in menu_categories (from requirement schema)
ALTER TABLE "public"."menu_categories"
ADD CONSTRAINT "menu_categories_restaurant_id_fkey" FOREIGN KEY ("restaurant_id") REFERENCES "public"."restaurants"("id") ON DELETE CASCADE;

-- Already has restaurant_id in menu_items (from requirement schema)
ALTER TABLE "public"."menu_items"
ADD CONSTRAINT "menu_items_restaurant_id_fkey" FOREIGN KEY ("restaurant_id") REFERENCES "public"."restaurants"("id") ON DELETE CASCADE;

-- Add FK for category_id in menu_items
ALTER TABLE "public"."menu_items"
ADD CONSTRAINT "menu_items_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "public"."menu_categories"("id") ON DELETE CASCADE;

-- Add FK for menu_item_photos
ALTER TABLE "public"."menu_item_photos"
ADD CONSTRAINT "menu_item_photos_menu_item_id_fkey" FOREIGN KEY ("menu_item_id") REFERENCES "public"."menu_items"("id") ON DELETE CASCADE;

-- Already has restaurant_id in modifier_groups (from requirement schema)
ALTER TABLE "public"."modifier_groups"
ADD CONSTRAINT "modifier_groups_restaurant_id_fkey" FOREIGN KEY ("restaurant_id") REFERENCES "public"."restaurants"("id") ON DELETE CASCADE;

-- Add FK for modifier_options
ALTER TABLE "public"."modifier_options"
ADD CONSTRAINT "modifier_options_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."modifier_groups"("id") ON DELETE CASCADE;

-- Add FK for menu_item_modifier_groups
ALTER TABLE "public"."menu_item_modifier_groups"
ADD CONSTRAINT "menu_item_modifier_groups_menu_item_id_fkey" FOREIGN KEY ("menu_item_id") REFERENCES "public"."menu_items"("id") ON DELETE CASCADE,
ADD CONSTRAINT "menu_item_modifier_groups_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."modifier_groups"("id") ON DELETE CASCADE;


-- =====================================================
-- SEED DATA FOR MENU MANAGEMENT
-- =====================================================

-- Insert Menu Categories
INSERT INTO public.menu_categories (restaurant_id, name, description, display_order, status) VALUES
(1, 'Appetizers', 'Start your meal with our delicious appetizers', 1, 'active'),
(1, 'Main Courses', 'Our signature main dishes', 2, 'active'),
(1, 'Seafood', 'Fresh seafood selections', 3, 'active'),
(1, 'Steaks & Grills', 'Premium cuts grilled to perfection', 4, 'active'),
(1, 'Pasta & Risotto', 'Italian favorites', 5, 'active'),
(1, 'Desserts', 'Sweet endings to your meal', 6, 'active'),
(1, 'Beverages', 'Drinks and refreshments', 7, 'active'),
(1, 'Kids Menu', 'Special meals for children', 8, 'active')
ON CONFLICT DO NOTHING;

-- Insert Menu Items
INSERT INTO public.menu_items (restaurant_id, category_id, name, description, price, prep_time_minutes, status, is_chef_recommended, is_deleted) VALUES
-- Appetizers (category_id: 1)
(1, 1, 'Caesar Salad', 'Fresh romaine lettuce with parmesan cheese and croutons', 12.00, 10, 'available', false, false),
(1, 1, 'Garlic Bread', 'Toasted bread with garlic butter and herbs', 6.00, 5, 'available', false, false),
(1, 1, 'Bruschetta', 'Grilled bread topped with tomatoes, basil and olive oil', 8.50, 8, 'available', true, false),
(1, 1, 'Chicken Wings', 'Crispy wings with your choice of sauce', 14.00, 15, 'available', false, false),

-- Main Courses (category_id: 2)
(1, 2, 'Grilled Chicken Breast', 'Tender chicken breast with seasonal vegetables', 24.00, 20, 'available', true, false),
(1, 2, 'BBQ Ribs', 'Slow-cooked pork ribs with BBQ sauce', 28.00, 25, 'available', false, false),
(1, 2, 'Lamb Chops', 'Grilled lamb chops with mint sauce', 35.00, 22, 'available', true, false),

-- Seafood (category_id: 3)
(1, 3, 'Grilled Salmon', 'Fresh Atlantic salmon with lemon butter sauce', 45.00, 18, 'available', true, false),
(1, 3, 'Lobster Pasta', 'Fresh lobster with creamy pasta', 55.00, 25, 'available', true, false),
(1, 3, 'Fish & Chips', 'Crispy battered fish with french fries', 22.00, 15, 'available', false, false),
(1, 3, 'Shrimp Scampi', 'Garlic butter shrimp with linguine', 32.00, 18, 'available', false, false),

-- Steaks & Grills (category_id: 4)
(1, 4, 'Ribeye Steak', '12oz premium ribeye steak', 65.00, 25, 'available', true, false),
(1, 4, 'Beef Tenderloin', '8oz tender beef fillet', 58.00, 22, 'available', true, false),
(1, 4, 'T-Bone Steak', '16oz T-bone steak', 70.00, 28, 'available', false, false),
(1, 4, 'Wagyu Burger', 'Premium wagyu beef burger with fries', 28.00, 15, 'sold_out', false, false),

-- Pasta & Risotto (category_id: 5)
(1, 5, 'Spaghetti Carbonara', 'Classic Italian carbonara with bacon', 18.00, 15, 'available', false, false),
(1, 5, 'Seafood Risotto', 'Creamy risotto with mixed seafood', 26.00, 20, 'available', true, false),
(1, 5, 'Penne Arrabbiata', 'Spicy tomato sauce with penne pasta', 16.00, 12, 'available', false, false),

-- Desserts (category_id: 6)
(1, 6, 'Tiramisu', 'Classic Italian coffee-flavored dessert', 9.00, 5, 'available', true, false),
(1, 6, 'Chocolate Lava Cake', 'Warm chocolate cake with vanilla ice cream', 10.00, 8, 'available', true, false),
(1, 6, 'Crème Brûlée', 'French custard with caramelized sugar', 8.50, 5, 'available', false, false),
(1, 6, 'Ice Cream Selection', 'Choose from vanilla, chocolate, or strawberry', 6.00, 2, 'available', false, false),

-- Beverages (category_id: 7)
(1, 7, 'Soft Drink', 'Coke, Sprite, or Fanta', 5.00, 2, 'available', false, false),
(1, 7, 'Fresh Orange Juice', 'Freshly squeezed orange juice', 7.00, 3, 'available', false, false),
(1, 7, 'Coffee', 'Espresso, Americano, or Cappuccino', 6.00, 5, 'available', false, false),
(1, 7, 'Red Wine', 'House red wine', 30.00, 2, 'available', false, false),
(1, 7, 'White Wine', 'House white wine', 28.00, 2, 'available', false, false),

-- Kids Menu (category_id: 8)
(1, 8, 'Kids Chicken Nuggets', 'Crispy chicken nuggets with fries', 12.00, 10, 'available', false, false),
(1, 8, 'Kids Pasta', 'Simple pasta with butter or tomato sauce', 10.00, 8, 'available', false, false),
(1, 8, 'Kids Pizza', 'Small margherita pizza', 11.00, 12, 'unavailable', false, false);

-- Insert Menu Item Photos
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary) VALUES
-- Grilled Salmon photos
(8, '/uploads/menu-items/grilled-salmon-1.jpg', true),
(8, '/uploads/menu-items/grilled-salmon-2.jpg', false),

-- Lobster Pasta photos
(9, '/uploads/menu-items/lobster-pasta-1.jpg', true),

-- Ribeye Steak photos
(12, '/uploads/menu-items/ribeye-steak-1.jpg', true),
(12, '/uploads/menu-items/ribeye-steak-2.jpg', false),

-- Tiramisu photos
(19, '/uploads/menu-items/tiramisu-1.jpg', true),

-- Chocolate Lava Cake photos
(20, '/uploads/menu-items/chocolate-lava-cake-1.jpg', true),

-- Caesar Salad photos
(1, '/uploads/menu-items/caesar-salad-1.jpg', true),

-- BBQ Ribs photos
(6, '/uploads/menu-items/bbq-ribs-1.jpg', true);

-- Insert Modifier Groups
INSERT INTO public.modifier_groups (restaurant_id, name, selection_type, is_required, min_selections, max_selections, display_order, status) VALUES
-- For steaks
(1, 'Steak Temperature', 'single', true, 1, 1, 1, 'active'),
(1, 'Steak Side Dish', 'single', true, 1, 1, 2, 'active'),
(1, 'Extra Toppings', 'multiple', false, 0, 3, 3, 'active'),

-- For beverages
(1, 'Drink Size', 'single', true, 1, 1, 1, 'active'),
(1, 'Ice Preference', 'single', false, 0, 1, 2, 'active'),

-- For pasta
(1, 'Pasta Extras', 'multiple', false, 0, 5, 1, 'active'),

-- For burgers
(1, 'Burger Extras', 'multiple', false, 0, 5, 1, 'active');

-- Insert Modifier Options
INSERT INTO public.modifier_options (group_id, name, price_adjustment, status) VALUES
-- Steak Temperature options (group_id: 1)
(1, 'Rare', 0.00, 'active'),
(1, 'Medium Rare', 0.00, 'active'),
(1, 'Medium', 0.00, 'active'),
(1, 'Medium Well', 0.00, 'active'),
(1, 'Well Done', 0.00, 'active'),

-- Steak Side Dish options (group_id: 2)
(2, 'French Fries', 0.00, 'active'),
(2, 'Mashed Potatoes', 2.00, 'active'),
(2, 'Grilled Vegetables', 3.00, 'active'),
(2, 'Baked Potato', 2.50, 'active'),
(2, 'Coleslaw', 1.50, 'active'),

-- Extra Toppings options (group_id: 3)
(3, 'Fried Egg', 2.00, 'active'),
(3, 'Mushrooms', 3.00, 'active'),
(3, 'Blue Cheese', 3.50, 'active'),
(3, 'Garlic Butter', 2.00, 'active'),

-- Drink Size options (group_id: 4)
(4, 'Small', 0.00, 'active'),
(4, 'Medium', 1.00, 'active'),
(4, 'Large', 2.00, 'active'),

-- Ice Preference options (group_id: 5)
(5, 'No Ice', 0.00, 'active'),
(5, 'Light Ice', 0.00, 'active'),
(5, 'Regular Ice', 0.00, 'active'),
(5, 'Extra Ice', 0.00, 'active'),

-- Pasta Extras options (group_id: 6)
(6, 'Extra Cheese', 2.00, 'active'),
(6, 'Grilled Chicken', 5.00, 'active'),
(6, 'Bacon', 3.00, 'active'),
(6, 'Mushrooms', 2.50, 'active'),
(6, 'Garlic Bread', 3.00, 'active'),

-- Burger Extras options (group_id: 7)
(7, 'Extra Cheese', 1.50, 'active'),
(7, 'Bacon', 2.50, 'active'),
(7, 'Fried Egg', 2.00, 'active'),
(7, 'Avocado', 3.00, 'active'),
(7, 'Jalapeños', 1.00, 'active');

-- Attach Modifier Groups to Menu Items
INSERT INTO public.menu_item_modifier_groups (menu_item_id, group_id) VALUES
-- Ribeye Steak (menu_item_id: 12)
(12, 1), -- Steak Temperature
(12, 2), -- Steak Side Dish
(12, 3), -- Extra Toppings

-- Beef Tenderloin (menu_item_id: 13)
(13, 1), -- Steak Temperature
(13, 2), -- Steak Side Dish
(13, 3), -- Extra Toppings

-- T-Bone Steak (menu_item_id: 14)
(14, 1), -- Steak Temperature
(14, 2), -- Steak Side Dish
(14, 3), -- Extra Toppings

-- Wagyu Burger (menu_item_id: 15)
(15, 7), -- Burger Extras

-- Soft Drink (menu_item_id: 22)
(22, 4), -- Drink Size
(22, 5), -- Ice Preference

-- Fresh Orange Juice (menu_item_id: 23)
(23, 4), -- Drink Size
(23, 5), -- Ice Preference

-- Spaghetti Carbonara (menu_item_id: 16)
(16, 6), -- Pasta Extras

-- Seafood Risotto (menu_item_id: 17)
(17, 6), -- Pasta Extras

-- Penne Arrabbiata (menu_item_id: 18)
(18, 6); -- Pasta Extras


