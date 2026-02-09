-- Add test users with easy credentials for development/testing
-- IMPORTANT: These are test credentials only! Do not use in production.

-- Password hashes (bcrypt with cost 10):
-- "admin" -> $2a$10$7Z9O8K5q5Z5Q5Q5Q5Q5Q5eFQ5Q5Q5Q5Q5Q5Q5Q5Q5Q5Q5Q5Q5Q5Q5
-- "owner1" -> $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- "owner2" -> $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- "owner3" -> $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- "user1" -> $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
-- "password" -> $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

-- Note: For simplicity, all test passwords use the same hash (password: "password")
-- In a real application, you'd generate unique hashes

-- Insert Admin Users
INSERT INTO users (name, email, password, role, created_at, updated_at)
VALUES
    ('Admin', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Super Admin', 'superadmin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

-- Insert Owner Users
INSERT INTO users (name, email, password, role, created_at, updated_at)
VALUES
    ('Owner One', 'owner1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Owner Two', 'owner2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Owner Three', 'owner3', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

-- Insert Regular Users
INSERT INTO users (name, email, password, role, created_at, updated_at)
VALUES
    ('User One', 'user1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('User Two', 'user2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('User Three', 'user3', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('User Four', 'user4', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('User Five', 'user5', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

-- Add test resources owned by owner users
-- First, get the IDs of the owner users we just created
DO $$
DECLARE
    owner1_id INT;
    owner2_id INT;
    owner3_id INT;
BEGIN
    -- Get owner IDs
    SELECT id INTO owner1_id FROM users WHERE email = 'owner1';
    SELECT id INTO owner2_id FROM users WHERE email = 'owner2';
    SELECT id INTO owner3_id FROM users WHERE email = 'owner3';

    -- Insert resources for owner1
    IF owner1_id IS NOT NULL THEN
        INSERT INTO resources (name, description, capacity, owner_id, category_id, address, city, latitude, longitude, price_per_hour, is_active, created_at, updated_at)
        VALUES
            ('Downtown Office Suite', 'Modern office space in the heart of downtown', 15, owner1_id, 1, '123 Main St', 'New York', 40.7128, -74.0060, 50.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Rooftop Conference Room', 'Stunning views for important meetings', 25, owner1_id, 2, '456 Sky Tower', 'New York', 40.7489, -73.9680, 75.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Creative Workshop Space', 'Perfect for brainstorming and team building', 20, owner1_id, 3, '789 Innovation Ave', 'New York', 40.7614, -73.9776, 60.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Insert resources for owner2
    IF owner2_id IS NOT NULL THEN
        INSERT INTO resources (name, description, capacity, owner_id, category_id, address, city, latitude, longitude, price_per_hour, is_active, created_at, updated_at)
        VALUES
            ('Beachside Villa', 'Luxury accommodation with ocean views', 8, owner2_id, 4, '101 Ocean Drive', 'Miami', 25.7617, -80.1918, 150.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Mountain Cabin Retreat', 'Peaceful getaway in nature', 6, owner2_id, 4, '202 Pine Ridge', 'Denver', 39.7392, -104.9903, 80.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Urban Loft Apartment', 'Stylish downtown living space', 4, owner2_id, 4, '303 City Center', 'Chicago', 41.8781, -87.6298, 45.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Insert resources for owner3
    IF owner3_id IS NOT NULL THEN
        INSERT INTO resources (name, description, capacity, owner_id, category_id, address, city, latitude, longitude, price_per_hour, is_active, created_at, updated_at)
        VALUES
            ('Professional Recording Studio', 'State-of-the-art audio equipment', 5, owner3_id, 5, '404 Sound Wave Blvd', 'Los Angeles', 34.0522, -118.2437, 100.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Photography Studio', 'Natural lighting and backdrops', 10, owner3_id, 5, '505 Lens Lane', 'Los Angeles', 34.0522, -118.2437, 65.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Event Hall', 'Large space for weddings and parties', 200, owner3_id, 6, '606 Celebration St', 'Las Vegas', 36.1699, -115.1398, 200.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            ('Small Meeting Pod', 'Quick huddle space for teams', 4, owner3_id, 2, '707 Quick Meet Dr', 'San Francisco', 37.7749, -122.4194, 25.00, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Add test bookings for these resources
DO $$
DECLARE
    user1_id INT;
    user2_id INT;
    user3_id INT;
    res1_id INT;
    res2_id INT;
    res3_id INT;
BEGIN
    -- Get user IDs
    SELECT id INTO user1_id FROM users WHERE email = 'user1';
    SELECT id INTO user2_id FROM users WHERE email = 'user2';
    SELECT id INTO user3_id FROM users WHERE email = 'user3';

    -- Get some resource IDs from owner1
    SELECT id INTO res1_id FROM resources WHERE name = 'Downtown Office Suite' LIMIT 1;
    SELECT id INTO res2_id FROM resources WHERE name = 'Rooftop Conference Room' LIMIT 1;
    SELECT id INTO res3_id FROM resources WHERE name = 'Beachside Villa' LIMIT 1;

    -- Insert bookings if we have valid IDs
    IF user1_id IS NOT NULL AND res1_id IS NOT NULL THEN
        INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, total_price, notes, created_at, updated_at)
        VALUES
            (user1_id, res1_id, CURRENT_TIMESTAMP + INTERVAL '1 day', CURRENT_TIMESTAMP + INTERVAL '1 day' + INTERVAL '3 hours', 'confirmed', 150.00, 'Team meeting', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            (user1_id, res1_id, CURRENT_TIMESTAMP + INTERVAL '7 days', CURRENT_TIMESTAMP + INTERVAL '7 days' + INTERVAL '4 hours', 'confirmed', 200.00, 'Client presentation', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    IF user2_id IS NOT NULL AND res2_id IS NOT NULL THEN
        INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, total_price, notes, created_at, updated_at)
        VALUES
            (user2_id, res2_id, CURRENT_TIMESTAMP + INTERVAL '2 days', CURRENT_TIMESTAMP + INTERVAL '2 days' + INTERVAL '2 hours', 'confirmed', 150.00, 'Board meeting', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            (user2_id, res2_id, CURRENT_TIMESTAMP + INTERVAL '5 days', CURRENT_TIMESTAMP + INTERVAL '5 days' + INTERVAL '3 hours', 'pending', 225.00, 'Strategy session', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    IF user3_id IS NOT NULL AND res3_id IS NOT NULL THEN
        INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, total_price, notes, created_at, updated_at)
        VALUES
            (user3_id, res3_id, CURRENT_TIMESTAMP + INTERVAL '3 days', CURRENT_TIMESTAMP + INTERVAL '3 days' + INTERVAL '24 hours', 'confirmed', 3600.00, 'Weekend getaway', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
            (user3_id, res3_id, CURRENT_TIMESTAMP + INTERVAL '14 days', CURRENT_TIMESTAMP + INTERVAL '14 days' + INTERVAL '48 hours', 'pending', 7200.00, 'Family vacation', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Add test reviews for some resources
DO $$
DECLARE
    user1_id INT;
    user2_id INT;
    user3_id INT;
    res1_id INT;
    res2_id INT;
    res3_id INT;
BEGIN
    -- Get user IDs
    SELECT id INTO user1_id FROM users WHERE email = 'user1';
    SELECT id INTO user2_id FROM users WHERE email = 'user2';
    SELECT id INTO user3_id FROM users WHERE email = 'user3';

    -- Get resource IDs
    SELECT id INTO res1_id FROM resources WHERE name = 'Downtown Office Suite' LIMIT 1;
    SELECT id INTO res2_id FROM resources WHERE name = 'Rooftop Conference Room' LIMIT 1;
    SELECT id INTO res3_id FROM resources WHERE name = 'Beachside Villa' LIMIT 1;

    -- Insert reviews
    IF user1_id IS NOT NULL AND res1_id IS NOT NULL THEN
        INSERT INTO reviews (user_id, resource_id, rating, comment, created_at, updated_at)
        VALUES
            (user1_id, res1_id, 5, 'Amazing space! Perfect for our team meetings. Very professional setup.', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    IF user2_id IS NOT NULL AND res2_id IS NOT NULL THEN
        INSERT INTO reviews (user_id, resource_id, rating, comment, created_at, updated_at)
        VALUES
            (user2_id, res2_id, 4, 'Great location and views. The room was a bit warm but overall excellent.', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    IF user3_id IS NOT NULL AND res3_id IS NOT NULL THEN
        INSERT INTO reviews (user_id, resource_id, rating, comment, created_at, updated_at)
        VALUES
            (user3_id, res3_id, 5, 'Absolutely stunning! The villa exceeded all expectations. Will definitely book again!', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Summary of test accounts created
COMMENT ON TABLE users IS 'Test accounts:
Admin: email=admin, password=password (role=admin)
Admin: email=superadmin, password=password (role=admin)
Owner: email=owner1, password=password (has 3 resources)
Owner: email=owner2, password=password (has 3 resources)
Owner: email=owner3, password=password (has 4 resources)
Users: email=user1/user2/user3/user4/user5, password=password
All passwords are "password" (bcrypt hashed)';
