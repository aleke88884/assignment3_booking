-- Seed данные для категорий и цен

-- Добавляем категории
INSERT INTO resource_categories (name, slug, description, icon) VALUES
('Баня и сауна', 'bath-sauna', 'Традиционные бани, финские сауны, хаммамы', '🧖'),
('Бассейн', 'pool', 'Крытые и открытые бассейны', '🏊'),
('Спортивная площадка', 'sports-field', 'Футбольные, баскетбольные, волейбольные площадки', '⚽'),
('Конференц-зал', 'conference', 'Залы для мероприятий и встреч', '🏢'),
('Коворкинг', 'coworking', 'Рабочие пространства', '💼'),
('Студия', 'studio', 'Фото, танцевальные, йога студии', '🎨')
ON CONFLICT (slug) DO NOTHING;

-- Обновляем существующие ресурсы с категориями и ценами
UPDATE resources SET
    category_id = (SELECT id FROM resource_categories WHERE slug = 'conference'),
    city = 'Алматы',
    address = 'ул. Абая 123',
    amenities = ARRAY['Wi-Fi', 'Проектор', 'Доска', 'Кондиционер'],
    price_per_hour = 5000.00,
    is_active = true
WHERE name LIKE '%онференц-зал%';

UPDATE resources SET
    category_id = (SELECT id FROM resource_categories WHERE slug = 'conference'),
    city = 'Алматы',
    address = 'ул. Абая 123',
    amenities = ARRAY['Wi-Fi', 'ТВ'],
    price_per_hour = 2500.00,
    is_active = true
WHERE name LIKE '%ереговорная%';

UPDATE resources SET
    category_id = (SELECT id FROM resource_categories WHERE slug = 'coworking'),
    city = 'Алматы',
    address = 'ул. Абая 123',
    amenities = ARRAY['Wi-Fi', 'Кофе', 'Принтер'],
    price_per_hour = 1500.00,
    is_active = true
WHERE name LIKE '%оворкинг%';

UPDATE resources SET
    category_id = (SELECT id FROM resource_categories WHERE slug = 'conference'),
    city = 'Алматы',
    address = 'ул. Абая 123',
    amenities = ARRAY['Wi-Fi', 'Проектор', 'Микрофон', 'Кондиционер'],
    price_per_hour = 8000.00,
    is_active = true
WHERE name LIKE '%екционный%';

UPDATE resources SET
    category_id = (SELECT id FROM resource_categories WHERE slug = 'studio'),
    city = 'Алматы',
    address = 'ул. Абая 123',
    amenities = ARRAY['Wi-Fi', 'Доска', 'Маркеры'],
    price_per_hour = 3500.00,
    is_active = true
WHERE name LIKE '%тудия%';

-- Добавляем новые ресурсы для бань и бассейнов
INSERT INTO resources (name, description, capacity, category_id, city, address, amenities, price_per_hour, is_active) VALUES
('Баня "Алтын"', 'Классическая русская баня с бассейном и комнатой отдыха', 8,
    (SELECT id FROM resource_categories WHERE slug = 'bath-sauna'),
    'Алматы', 'мкр. Самал-2, д. 15',
    ARRAY['Бассейн', 'Душ', 'Комната отдыха', 'Парковка', 'Wi-Fi', 'Караоке'],
    12000.00, true),

('Финская сауна "Nord"', 'Современная финская сауна с панорамным видом', 6,
    (SELECT id FROM resource_categories WHERE slug = 'bath-sauna'),
    'Алматы', 'ул. Достык 120',
    ARRAY['Душ', 'Джакузи', 'Комната отдыха', 'Парковка'],
    15000.00, true),

('Бассейн "Aqua Center"', 'Крытый бассейн 25м с тренажерным залом', 20,
    (SELECT id FROM resource_categories WHERE slug = 'pool'),
    'Алматы', 'пр. Райымбека 250',
    ARRAY['Душевые', 'Раздевалки', 'Парковка', 'Кафе'],
    8000.00, true),

('Футбольное поле "Arena"', 'Поле с искусственным покрытием и освещением', 14,
    (SELECT id FROM resource_categories WHERE slug = 'sports-field'),
    'Алматы', 'мкр. Мамыр, д. 5',
    ARRAY['Раздевалки', 'Душ', 'Парковка', 'Освещение', 'Мячи в аренду'],
    18000.00, true),

('Баскетбольная площадка "Hoops"', 'Крытая площадка с профессиональным покрытием', 10,
    (SELECT id FROM resource_categories WHERE slug = 'sports-field'),
    'Алматы', 'ул. Жандосова 98',
    ARRAY['Раздевалки', 'Душ', 'Парковка', 'Трибуны'],
    10000.00, true)
ON CONFLICT DO NOTHING;

-- Добавляем расписание работы (для новых ресурсов)
DO $$
DECLARE
    resource_rec RECORD;
BEGIN
    FOR resource_rec IN
        SELECT id FROM resources
        WHERE category_id IN (
            SELECT id FROM resource_categories
            WHERE slug IN ('bath-sauna', 'pool', 'sports-field')
        )
    LOOP
        INSERT INTO resource_schedules (resource_id, day_of_week, open_time, close_time, is_closed)
        VALUES
            (resource_rec.id, 1, '09:00', '23:00', false),
            (resource_rec.id, 2, '09:00', '23:00', false),
            (resource_rec.id, 3, '09:00', '23:00', false),
            (resource_rec.id, 4, '09:00', '23:00', false),
            (resource_rec.id, 5, '09:00', '23:00', false),
            (resource_rec.id, 6, '10:00', '02:00', false),
            (resource_rec.id, 0, '10:00', '02:00', false)
        ON CONFLICT (resource_id, day_of_week) DO NOTHING;
    END LOOP;
END $$;

-- Добавляем тарифы (будни дешевле, выходные дороже)
DO $$
DECLARE
    resource_rec RECORD;
BEGIN
    FOR resource_rec IN
        SELECT id, price_per_hour, name FROM resources
        WHERE category_id IN (
            SELECT id FROM resource_categories
            WHERE slug IN ('bath-sauna', 'pool', 'sports-field')
        )
    LOOP
        -- Тариф для будних дней
        INSERT INTO resource_pricing (resource_id, name, price, duration_minutes, day_of_week, time_from, time_to, is_active)
        SELECT resource_rec.id, 'Будний день', resource_rec.price_per_hour, 60, day, '09:00', '18:00', true
        FROM generate_series(1, 5) AS day
        ON CONFLICT DO NOTHING;

        -- Тариф для вечера будних дней (дороже)
        INSERT INTO resource_pricing (resource_id, name, price, duration_minutes, day_of_week, time_from, time_to, is_active)
        SELECT resource_rec.id, 'Вечер будни', resource_rec.price_per_hour * 1.3, 60, day, '18:00', '23:00', true
        FROM generate_series(1, 5) AS day
        ON CONFLICT DO NOTHING;

        -- Тариф для выходных (дороже)
        INSERT INTO resource_pricing (resource_id, name, price, duration_minutes, day_of_week, time_from, time_to, is_active)
        SELECT resource_rec.id, 'Выходной', resource_rec.price_per_hour * 1.5, 60, day, '10:00', '02:00', true
        FROM generate_series(0, 0) AS day
        UNION ALL
        SELECT resource_rec.id, 'Выходной', resource_rec.price_per_hour * 1.5, 60, day, '10:00', '02:00', true
        FROM generate_series(6, 6) AS day
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;
