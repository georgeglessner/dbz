-- Seed data for e-commerce database
-- Sample customers, products, and orders

-- Insert sample customers
INSERT INTO customers (email, password_hash, first_name, last_name, phone) VALUES
('alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Alice', 'Johnson', '555-0101'),
('bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Bob', 'Smith', '555-0102'),
('carol@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Carol', 'Williams', '555-0103'),
('david@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'David', 'Brown', '555-0104')
ON CONFLICT (email) DO NOTHING;

-- Insert categories
INSERT INTO categories (name, slug, description) VALUES
('Electronics', 'electronics', 'Gadgets and electronic devices'),
('Clothing', 'clothing', 'Apparel and fashion items'),
('Books', 'books', 'Physical and digital books'),
('Home & Garden', 'home-garden', 'Home decor and gardening supplies')
ON CONFLICT (slug) DO NOTHING;

-- Insert subcategories
INSERT INTO categories (name, slug, description, parent_id) VALUES
('Smartphones', 'smartphones', 'Mobile phones and accessories', 1),
('Laptops', 'laptops', 'Notebook computers and accessories', 1),
('Men''s Clothing', 'mens-clothing', 'Apparel for men', 2),
('Women''s Clothing', 'womens-clothing', 'Apparel for women', 2),
('Fiction', 'fiction', 'Fiction books and novels', 3),
('Non-Fiction', 'non-fiction', 'Educational and reference books', 3)
ON CONFLICT (slug) DO NOTHING;

-- Insert sample products
INSERT INTO products (sku, name, slug, description, price, quantity, weight, status, category_id) VALUES
('PHONE-001', 'SmartPhone Pro X', 'smartphone-pro-x',
'Latest flagship smartphone with 6.7" OLED display, 256GB storage, and advanced camera system. Features 5G connectivity and all-day battery life.',
999.99, 50, 0.5, 'active', (SELECT id FROM categories WHERE slug = 'smartphones')),

('PHONE-002', 'Budget Phone Plus', 'budget-phone-plus',
'Affordable smartphone with great features. 6.1" display, 128GB storage, dual camera setup, and 2-day battery life.',
299.99, 100, 0.4, 'active', (SELECT id FROM categories WHERE slug = 'smartphones')),

('LAPTOP-001', 'ProBook 15"', 'probook-15',
'High-performance laptop for professionals. Intel i7 processor, 16GB RAM, 512GB SSD, and dedicated graphics card. Perfect for development and design work.',
1499.99, 25, 2.1, 'active', (SELECT id FROM categories WHERE slug = 'laptops')),

('LAPTOP-002', 'UltraBook Air', 'ultrabook-air',
'Ultra-portable laptop weighing just 2.5 lbs. 13" retina display, 8GB RAM, 256GB SSD. 12-hour battery life makes it perfect for travel.',
899.99, 40, 1.8, 'active', (SELECT id FROM categories WHERE slug = 'laptops')),

('SHIRT-001', 'Classic Cotton T-Shirt', 'classic-cotton-t-shirt',
'Comfortable 100% cotton t-shirt available in multiple colors. Pre-shrunk and machine washable.',
29.99, 200, 0.3, 'active', (SELECT id FROM categories WHERE slug = 'mens-clothing')),

('BOOK-001', 'The Art of Programming', 'art-of-programming',
'Comprehensive guide to software development principles, algorithms, and best practices. A must-read for any serious developer.',
49.99, 75, 1.2, 'active', (SELECT id FROM categories WHERE slug = 'fiction')),

('BOOK-002', 'Learning SQL', 'learning-sql',
'Beginner-friendly introduction to SQL and database design. Covers PostgreSQL, MySQL, and SQLite with practical examples.',
34.99, 100, 0.9, 'active', (SELECT id FROM categories WHERE slug = 'fiction'))
ON CONFLICT (sku) DO NOTHING;

-- Insert product images
INSERT INTO product_images (product_id, url, alt_text, position, is_primary) VALUES
(1, 'https://example.com/images/phone-pro-x-1.jpg', 'SmartPhone Pro X front view', 1, TRUE),
(1, 'https://example.com/images/phone-pro-x-2.jpg', 'SmartPhone Pro X back view', 2, FALSE),
(2, 'https://example.com/images/budget-phone-plus-1.jpg', 'Budget Phone Plus', 1, TRUE),
(3, 'https://example.com/images/probook-15-1.jpg', 'ProBook 15 inch laptop', 1, TRUE),
(4, 'https://example.com/images/ultrabook-air-1.jpg', 'UltraBook Air laptop', 1, TRUE),
(5, 'https://example.com/images/tshirt-1.jpg', 'Classic Cotton T-Shirt', 1, TRUE),
(6, 'https://example.com/images/book-programming-1.jpg', 'The Art of Programming book', 1, TRUE),
(7, 'https://example.com/images/book-sql-1.jpg', 'Learning SQL book', 1, TRUE);

-- Insert addresses for customers
INSERT INTO addresses (customer_id, type, first_name, last_name, address1, city, state, postal_code, country, is_default) VALUES
(1, 'shipping', 'Alice', 'Johnson', '123 Main St', 'New York', 'NY', '10001', 'USA', TRUE),
(1, 'billing', 'Alice', 'Johnson', '123 Main St', 'New York', 'NY', '10001', 'USA', TRUE),
(2, 'shipping', 'Bob', 'Smith', '456 Oak Ave', 'Los Angeles', 'CA', '90001', 'USA', TRUE),
(3, 'shipping', 'Carol', 'Williams', '789 Pine Rd', 'Chicago', 'IL', '60601', 'USA', TRUE),
(4, 'shipping', 'David', 'Brown', '321 Elm St', 'Houston', 'TX', '77001', 'USA', TRUE);

-- Insert sample orders
INSERT INTO orders (order_number, customer_id, status, payment_status, fulfillment_status, subtotal, tax_amount, shipping_amount, total_amount, shipping_address_id, billing_address_id) VALUES
('ORD-2024-001', 1, 'delivered', 'paid', 'fulfilled', 999.99, 80.00, 15.00, 1094.99, 1, 2),
('ORD-2024-002', 2, 'shipped', 'paid', 'fulfilled', 1199.98, 96.00, 0.00, 1295.98, 3, 3),
('ORD-2024-003', 3, 'processing', 'paid', 'unfulfilled', 79.98, 6.40, 10.00, 96.38, 4, 4),
('ORD-2024-004', 1, 'pending', 'pending', 'unfulfilled', 34.99, 2.80, 5.00, 42.79, 1, 2),
('ORD-2024-005', 4, 'cancelled', 'refunded', 'unfulfilled', 299.99, 24.00, 15.00, 338.99, 5, 5);

-- Insert order items
INSERT INTO order_items (order_id, product_id, quantity, price, total) VALUES
(1, 1, 1, 999.99, 999.99),  -- Alice bought a smartphone
(2, 1, 1, 999.99, 999.99),  -- Bob bought a smartphone
(2, 5, 2, 29.99, 59.99),    -- Bob bought 2 t-shirts
(3, 5, 2, 29.99, 59.99),    -- Carol bought 2 t-shirts
(3, 7, 1, 19.99, 19.99),    -- Carol bought a book (on sale)
(4, 7, 1, 34.99, 34.99),    -- Alice bought another book
(5, 2, 1, 299.99, 299.99);  -- David bought a budget phone (then cancelled)

-- Insert reviews
INSERT INTO reviews (product_id, customer_id, rating, title, body, is_verified_purchase, status) VALUES
(1, 1, 5, 'Amazing phone!', 'Best smartphone I have ever owned. Camera is incredible and battery lasts all day.', TRUE, 'approved'),
(1, 2, 4, 'Great but expensive', 'Excellent phone, but the price is quite high. Still worth it for the features.', TRUE, 'approved'),
(3, 1, 5, 'Perfect for work', 'This laptop handles all my development work perfectly. Fast and reliable.', FALSE, 'approved'),
(6, 3, 5, 'Must-read for developers', 'Incredible depth of knowledge. Changed how I think about programming.', TRUE, 'approved'),
(7, 4, 4, 'Good SQL introduction', 'Clear explanations and practical examples. Perfect for beginners.', TRUE, 'pending');

-- Insert inventory logs
INSERT INTO inventory_logs (product_id, quantity_change, reason, reference_id, notes) VALUES
(1, 50, 'restock', NULL, 'Initial inventory'),
(2, 100, 'restock', NULL, 'Initial inventory'),
(3, 25, 'restock', NULL, 'Initial inventory'),
(4, 40, 'restock', NULL, 'Initial inventory'),
(5, 200, 'restock', NULL, 'Initial inventory'),
(1, -1, 'sale', 1, 'Order ORD-2024-001'),
(1, -1, 'sale', 2, 'Order ORD-2024-002'),
(5, -2, 'sale', 2, 'Order ORD-2024-002'),
(5, -2, 'sale', 3, 'Order ORD-2024-003'),
(2, -1, 'sale', 5, 'Order ORD-2024-005 - then cancelled');

-- Add some items to shopping carts
INSERT INTO cart_items (customer_id, product_id, quantity) VALUES
(1, 3, 1),  -- Alice has a laptop in cart
(2, 6, 1),  -- Bob has a book in cart
(4, 4, 1);  -- David has a laptop in cart
