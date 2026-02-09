-- Расширение для систем бронирования бань, бассейнов, футбольных площадок

-- Таблица категорий ресурсов
CREATE TABLE IF NOT EXISTS resource_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_slug ON resource_categories(slug);

CREATE TRIGGER set_timestamp_categories
    BEFORE UPDATE ON resource_categories
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

COMMENT ON TABLE resource_categories IS 'Категории ресурсов (баня, бассейн, площадка и тд)';

-- Добавляем новые поля в таблицу resources
ALTER TABLE resources ADD COLUMN IF NOT EXISTS category_id INT REFERENCES resource_categories(id) ON DELETE SET NULL;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS city VARCHAR(100);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS latitude DECIMAL(10, 8);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS longitude DECIMAL(11, 8);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS amenities TEXT[];
ALTER TABLE resources ADD COLUMN IF NOT EXISTS rules TEXT;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS price_per_hour DECIMAL(10, 2);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

CREATE INDEX IF NOT EXISTS idx_resources_category ON resources(category_id);
CREATE INDEX IF NOT EXISTS idx_resources_city ON resources(city);
CREATE INDEX IF NOT EXISTS idx_resources_is_active ON resources(is_active);
CREATE INDEX IF NOT EXISTS idx_resources_location ON resources(latitude, longitude);

COMMENT ON COLUMN resources.amenities IS 'Удобства: душ, парковка, wi-fi и тд';
COMMENT ON COLUMN resources.price_per_hour IS 'Базовая цена за час';

-- Таблица фотографий ресурсов
CREATE TABLE IF NOT EXISTS resource_photos (
    id BIGSERIAL PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    storage_key VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    width INT,
    height INT,
    is_primary BOOLEAN DEFAULT false,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_file_size CHECK (file_size > 0 AND file_size <= 10485760)
);

CREATE INDEX IF NOT EXISTS idx_photos_resource ON resource_photos(resource_id);
CREATE INDEX IF NOT EXISTS idx_photos_is_primary ON resource_photos(resource_id, is_primary);
CREATE INDEX IF NOT EXISTS idx_photos_order ON resource_photos(resource_id, display_order);

CREATE TRIGGER set_timestamp_photos
    BEFORE UPDATE ON resource_photos
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

COMMENT ON TABLE resource_photos IS 'Фотографии ресурсов';
COMMENT ON COLUMN resource_photos.storage_key IS 'Путь в S3/Cloud Storage';
COMMENT ON COLUMN resource_photos.is_primary IS 'Главная фотография для превью';

-- Таблица тарифов (гибкая ценовая политика)
CREATE TABLE IF NOT EXISTS resource_pricing (
    id BIGSERIAL PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    day_of_week INT CHECK (day_of_week BETWEEN 0 AND 6),
    time_from TIME,
    time_to TIME,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pricing_resource ON resource_pricing(resource_id);
CREATE INDEX IF NOT EXISTS idx_pricing_active ON resource_pricing(resource_id, is_active);

CREATE TRIGGER set_timestamp_pricing
    BEFORE UPDATE ON resource_pricing
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

COMMENT ON TABLE resource_pricing IS 'Тарифы и цены на бронирование';
COMMENT ON COLUMN resource_pricing.day_of_week IS '0=воскресенье, 1=понедельник, ..., 6=суббота';

-- Таблица расписания работы
CREATE TABLE IF NOT EXISTS resource_schedules (
    id BIGSERIAL PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    is_closed BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_schedule_time CHECK (close_time > open_time OR is_closed = true),
    UNIQUE(resource_id, day_of_week)
);

CREATE INDEX IF NOT EXISTS idx_schedules_resource ON resource_schedules(resource_id);

CREATE TRIGGER set_timestamp_schedules
    BEFORE UPDATE ON resource_schedules
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

COMMENT ON TABLE resource_schedules IS 'Расписание работы ресурсов';

-- Таблица отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    booking_id INT REFERENCES bookings(id) ON DELETE SET NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(booking_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_resource ON reviews(resource_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(resource_id, rating);

CREATE TRIGGER set_timestamp_reviews
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

COMMENT ON TABLE reviews IS 'Отзывы пользователей о ресурсах';
COMMENT ON COLUMN reviews.is_verified IS 'Отзыв от реального клиента с подтвержденным бронированием';

-- Обновляем цену в bookings (цена на момент бронирования)
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS total_price DECIMAL(10, 2);
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS notes TEXT;

COMMENT ON COLUMN bookings.total_price IS 'Итоговая цена бронирования';
COMMENT ON COLUMN bookings.notes IS 'Комментарий клиента к бронированию';
