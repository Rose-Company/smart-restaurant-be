-- Import Menu Item Photos
-- This script inserts photo URLs into menu_item_photos table based on menu item names

-- Nhóm Khai Vị (Appetizers)
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1467003909585-2f8a72700288?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Grilled Salmon Supreme'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1573140247632-f8fd74997d5c?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Garlic Bread'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1572656631137-7935297eff55?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Bruschetta'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1567620832903-9fc6debc209f?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Chicken Wings'
ON CONFLICT DO NOTHING;

-- Nhóm Món Chính (Main Course - Gà, Sườn, Cá)
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1532550907401-a500c9a57435?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Grilled Chicken Breast'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1544025162-d76694265947?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'BBQ Ribs'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1602847213180-50e43a80df31?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Lamb Chops'
ON CONFLICT DO NOTHING;

-- Nhóm Hải Sản & Steak
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1485921325833-c519f76c4927?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name LIKE 'Grilled Salmon%' AND name != 'Grilled Salmon Supreme'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1551183053-bf91a1d81141?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Lobster Pasta'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1524339102455-6f4021ca9c1b?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Fish & Chips'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.pbase.com/i35/35/603435/1/140255375.z6A8iL1E.ShrimpScampi1.jpg', TRUE, NOW()
FROM public.menu_items WHERE name = 'Shrimp Scampi'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1546241072-48010ad28c2c?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Ribeye Steak'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1558030006-450675393462?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Beef Tenderloin'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1594041680534-e8c8cdebd659?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'T-Bone Steak'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1586816001966-79b736744398?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Wagyu Burger'
ON CONFLICT DO NOTHING;

-- Nhóm Mì Ý & Risotto
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1612459284970-e8f027596582?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Spaghetti Carbonara'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1534422298391-e4f8c170db0a?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Seafood Risotto'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Penne Arrabbiata'
ON CONFLICT DO NOTHING;

-- Nhóm Tráng Miệng (Desserts)
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1571877227200-a0d98ea607e9?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Tiramisu'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1624353365286-3f8d62daad51?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Chocolate Lava Cake'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1470124182917-cc6e71b22ecc?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Crème Brûlée'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1501443762994-82bd5dace89a?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Ice Cream Selection'
ON CONFLICT DO NOTHING;

-- Nhóm Đồ Uống (Beverages)
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1581636625402-29b2a704ef13?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Soft Drink'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1613478223719-2ab802602423?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Fresh Orange Juice'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1541167760496-162955ed8a9f?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Coffee'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1510812431401-41d2bd2722f3?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Red Wine'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1559158068-930a7af4c744?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'White Wine'
ON CONFLICT DO NOTHING;

-- Nhóm Trẻ Em (Kids Menu)
INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1562967914-608f82629710?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Kids Chicken Nuggets'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1473093226795-af9932fe5856?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Kids Pasta'
ON CONFLICT DO NOTHING;

INSERT INTO public.menu_item_photos (menu_item_id, url, is_primary, created_at)
SELECT id, 'https://images.unsplash.com/photo-1513104890138-7c749659a591?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80', TRUE, NOW()
FROM public.menu_items WHERE name = 'Kids Pizza'
ON CONFLICT DO NOTHING;

-- Verify inserted photos
SELECT COUNT(*) as total_photos FROM public.menu_item_photos;

-- Show all inserted photos with item names
SELECT 
    mip.id,
    mi.name as menu_item_name,
    mip.url,
    mip.is_primary,
    mip.created_at
FROM public.menu_item_photos mip
JOIN public.menu_items mi ON mip.menu_item_id = mi.id
ORDER BY mi.id;
