-- Seed data for local debugging

-- Categories
INSERT INTO categories (id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Electronics', 'electronics', NULL, 0, TRUE, 1, NOW(), NOW()),
    ('11111111-1111-1111-1111-111111111112', 'Phones', 'phones', '11111111-1111-1111-1111-111111111111', 1, TRUE, 1, NOW(), NOW()),
    ('11111111-1111-1111-1111-111111111113', 'Laptops', 'laptops', '11111111-1111-1111-1111-111111111111', 1, TRUE, 2, NOW(), NOW()),
    ('11111111-1111-1111-1111-111111111114', 'Accessories', 'accessories', '11111111-1111-1111-1111-111111111111', 1, TRUE, 3, NOW(), NOW());

-- Customers
INSERT INTO customers (id, name, email, phone, address, created_at, updated_at) VALUES
    ('22222222-2222-2222-2222-222222222221', 'Alice', 'alice@example.com', '+10000000001', '1 Main St', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'Bob', 'bob@example.com', '+10000000002', '2 Side St', NOW(), NOW());

-- Products
INSERT INTO products (id, name, price, quantity, created_at, updated_at) VALUES
    ('33333333-3333-3333-3333-333333333331', 'iPhone', 999.00, 50, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333332', 'Android Phone', 499.00, 80, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'Laptop Pro', 1499.00, 30, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333334', 'USB-C Cable', 19.99, 200, NOW(), NOW());

-- Product categories
INSERT INTO product_catagories (product_id, catagory_id, created_at, updated_at) VALUES
    ('33333333-3333-3333-3333-333333333331', '11111111-1111-1111-1111-111111111112', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333332', '11111111-1111-1111-1111-111111111112', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111113', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333334', '11111111-1111-1111-1111-111111111114', NOW(), NOW());

-- Orders (dated in previous month to hit the top-5 query)
INSERT INTO orders (id, customer_id, total_price, status, created_at, updated_at) VALUES
    ('44444444-4444-4444-4444-444444444441', '22222222-2222-2222-2222-222222222221', 0, 'new', date_trunc('month', now()) - INTERVAL '10 days', date_trunc('month', now()) - INTERVAL '10 days'),
    ('44444444-4444-4444-4444-444444444442', '22222222-2222-2222-2222-222222222222', 0, 'new', date_trunc('month', now()) - INTERVAL '15 days', date_trunc('month', now()) - INTERVAL '15 days');

-- Order items
INSERT INTO order_items (id, order_id, product_id, quantity, sub_total, created_at, updated_at) VALUES
    ('55555555-5555-5555-5555-555555555551', '44444444-4444-4444-4444-444444444441', '33333333-3333-3333-3333-333333333331', 2, 1998.00, NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555552', '44444444-4444-4444-4444-444444444441', '33333333-3333-3333-3333-333333333334', 3, 59.97, NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555553', '44444444-4444-4444-4444-444444444442', '33333333-3333-3333-3333-333333333332', 4, 1996.00, NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555554', '44444444-4444-4444-4444-444444444442', '33333333-3333-3333-3333-333333333334', 5, 99.95, NOW(), NOW());

-- Update totals based on items
UPDATE orders o
SET total_price = sub.t
FROM (
    SELECT order_id, SUM(sub_total) AS t
    FROM order_items
    GROUP BY order_id
) sub
WHERE o.id = sub.order_id;

