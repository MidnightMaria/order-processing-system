-- Seed users
INSERT INTO users (name, email)
SELECT
    'User ' || i,
    'user' || i || '@example.com'
FROM generate_series(1, 1000) AS i;


-- Seed products
INSERT INTO products (name, price, stock)
SELECT
    'Product ' || i,
    (RANDOM() * 100000)::NUMERIC(12,2),
    (RANDOM() * 100)::INTEGER
FROM generate_series(1, 10000) AS i;


-- Seed orders
INSERT INTO orders (user_id, status, total_amount)
SELECT
    ((i - 1) % 1000) + 1,
    'pending',
    0
FROM generate_series(1, 10000) AS i;


-- Seed order items
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
SELECT
    ((i - 1) % 10000) + 1,
    ((i - 1) % 10000) + 1,
    1,
    100
FROM generate_series(1, 50000) AS i;