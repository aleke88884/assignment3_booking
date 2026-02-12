-- Add owner users and photos for resources

-- Add owner users
INSERT INTO users (name, email, password, role) VALUES
('Владелец Алексей', 'owner1@smartbooking.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'owner'),
('Владелец Марина', 'owner2@smartbooking.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'owner'),
('Владелец Дмитрий', 'owner3@smartbooking.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'owner')
ON CONFLICT (email) DO NOTHING;

-- Assign owners to existing resources
DO $$
DECLARE
    owner1_id INT;
    owner2_id INT;
    owner3_id INT;
BEGIN
    -- Get owner IDs
    SELECT id INTO owner1_id FROM users WHERE email = 'owner1@smartbooking.com';
    SELECT id INTO owner2_id FROM users WHERE email = 'owner2@smartbooking.com';
    SELECT id INTO owner3_id FROM users WHERE email = 'owner3@smartbooking.com';

    -- Assign owners to resources
    UPDATE resources SET owner_id = owner1_id WHERE name LIKE '%Баня%' OR name LIKE '%сауна%';
    UPDATE resources SET owner_id = owner2_id WHERE name LIKE '%Бассейн%' OR name LIKE '%площадка%';
    UPDATE resources SET owner_id = owner3_id WHERE name LIKE '%онференц%' OR name LIKE '%оворкинг%' OR name LIKE '%тудия%' OR name LIKE '%ереговорная%' OR name LIKE '%екционный%';
END $$;

-- Add placeholder photos for resources
DO $$
DECLARE
    resource_rec RECORD;
    photo_counter INT := 1;
    category_name TEXT;
    photo_url TEXT;
BEGIN
    FOR resource_rec IN SELECT r.id, r.name, rc.name as cat_name
                       FROM resources r
                       LEFT JOIN resource_categories rc ON r.category_id = rc.id
    LOOP
        category_name := COALESCE(resource_rec.cat_name, 'general');

        -- Determine photo URL based on category
        CASE
            WHEN category_name LIKE '%Баня%' OR category_name LIKE '%сауна%' THEN
                photo_url := 'https://placehold.co/800x600/8B4513/FFF?text=Sauna+' || resource_rec.id;
            WHEN category_name LIKE '%Бассейн%' THEN
                photo_url := 'https://placehold.co/800x600/4682B4/FFF?text=Pool+' || resource_rec.id;
            WHEN category_name LIKE '%площадка%' THEN
                photo_url := 'https://placehold.co/800x600/228B22/FFF?text=Sports+' || resource_rec.id;
            WHEN category_name LIKE '%онференц%' THEN
                photo_url := 'https://placehold.co/800x600/2F4F4F/FFF?text=Conference+' || resource_rec.id;
            WHEN category_name LIKE '%оворкинг%' THEN
                photo_url := 'https://placehold.co/800x600/4B0082/FFF?text=Coworking+' || resource_rec.id;
            WHEN category_name LIKE '%тудия%' THEN
                photo_url := 'https://placehold.co/800x600/FF69B4/FFF?text=Studio+' || resource_rec.id;
            ELSE
                photo_url := 'https://placehold.co/800x600/808080/FFF?text=Resource+' || resource_rec.id;
        END CASE;

        -- Insert primary photo
        INSERT INTO resource_photos (resource_id, url, storage_key, file_name, file_size, mime_type, width, height, is_primary, display_order)
        VALUES (
            resource_rec.id,
            photo_url,
            'placeholder/resource_' || resource_rec.id || '_1.jpg',
            'resource_' || resource_rec.id || '_1.jpg',
            102400,
            'image/jpeg',
            800,
            600,
            true,
            0
        )
        ON CONFLICT DO NOTHING;

        -- Insert additional photos (2-3 per resource)
        FOR i IN 1..2 LOOP
            INSERT INTO resource_photos (resource_id, url, storage_key, file_name, file_size, mime_type, width, height, is_primary, display_order)
            VALUES (
                resource_rec.id,
                'https://placehold.co/800x600/A9A9A9/000?text=Photo+' || (i + 1),
                'placeholder/resource_' || resource_rec.id || '_' || (i + 1) || '.jpg',
                'resource_' || resource_rec.id || '_' || (i + 1) || '.jpg',
                102400,
                'image/jpeg',
                800,
                600,
                false,
                i
            )
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Add some sample reviews
DO $$
DECLARE
    resource_rec RECORD;
    user_id INT;
BEGIN
    -- Get a regular user for reviews
    SELECT id INTO user_id FROM users WHERE role = 'user' LIMIT 1;

    IF user_id IS NOT NULL THEN
        -- Add reviews for some resources
        FOR resource_rec IN SELECT id FROM resources LIMIT 5 LOOP
            INSERT INTO reviews (resource_id, user_id, rating, comment, is_verified)
            VALUES (
                resource_rec.id,
                user_id,
                4 + (RANDOM() * 1)::INT, -- Random rating 4-5
                'Отличное место! Всё понравилось, обязательно вернёмся.',
                true
            )
            ON CONFLICT DO NOTHING;
        END LOOP;
    END IF;
END $$;

-- Update bookings with prices based on resource prices
UPDATE bookings b
SET total_price = r.price_per_hour *
    EXTRACT(EPOCH FROM (b.end_time - b.start_time)) / 3600
FROM resources r
WHERE b.resource_id = r.id
  AND b.total_price IS NULL
  AND r.price_per_hour IS NOT NULL;
