-- Password for all test users: password123
INSERT INTO users (name, email, password, role)
VALUES (
        'Admin User',
        'admin@smartbooking.com',
        '$2a$10$XVkjArT5QKlksmLDAnahNeRspvJ1lTwMjVGjyfPy5o2oNd2FkiaAC',
        'admin'
    ),
    (
        'John Doe',
        'john@example.com',
        '$2a$10$XVkjArT5QKlksmLDAnahNeRspvJ1lTwMjVGjyfPy5o2oNd2FkiaAC',
        'user'
    ),
    (
        'Jane Smith',
        'jane@example.com',
        '$2a$10$XVkjArT5QKlksmLDAnahNeRspvJ1lTwMjVGjyfPy5o2oNd2FkiaAC',
        'user'
    ),
    (
        'Bob Wilson',
        'bob@example.com',
        '$2a$10$XVkjArT5QKlksmLDAnahNeRspvJ1lTwMjVGjyfPy5o2oNd2FkiaAC',
        'user'
    ) ON CONFLICT (email) DO NOTHING;
INSERT INTO resources (name, description, capacity)
VALUES (
        'Конференц-зал А',
        'Большой зал для совещаний с проектором и доской',
        20
    ),
    (
        'Конференц-зал Б',
        'Средний зал для встреч команды',
        10
    ),
    (
        'Переговорная 1',
        'Маленькая комната для 1-на-1 встреч',
        4
    ),
    ('Переговорная 2', 'Маленькая комната с ТВ', 4),
    (
        'Коворкинг-зона',
        'Открытое пространство для совместной работы',
        30
    ),
    (
        'Лекционный зал',
        'Зал для презентаций и тренингов',
        50
    ),
    (
        'VIP переговорная',
        'Премиум комната для важных встреч',
        8
    ),
    (
        'Креативная студия',
        'Пространство для мозговых штурмов',
        12
    ) ON CONFLICT DO NOTHING;
INSERT INTO bookings (
        user_id,
        resource_id,
        start_time,
        end_time,
        status
    )
VALUES (
        2,
        1,
        CURRENT_TIMESTAMP + INTERVAL '1 day',
        CURRENT_TIMESTAMP + INTERVAL '1 day' + INTERVAL '2 hours',
        'confirmed'
    ),
    (
        3,
        2,
        CURRENT_TIMESTAMP + INTERVAL '2 days',
        CURRENT_TIMESTAMP + INTERVAL '2 days' + INTERVAL '1 hour',
        'confirmed'
    ),
    (
        4,
        3,
        CURRENT_TIMESTAMP + INTERVAL '3 days',
        CURRENT_TIMESTAMP + INTERVAL '3 days' + INTERVAL '30 minutes',
        'pending'
    ),
    (
        2,
        4,
        CURRENT_TIMESTAMP + INTERVAL '4 days',
        CURRENT_TIMESTAMP + INTERVAL '4 days' + INTERVAL '1 hour',
        'confirmed'
    ),
    (
        3,
        5,
        CURRENT_TIMESTAMP + INTERVAL '5 days',
        CURRENT_TIMESTAMP + INTERVAL '5 days' + INTERVAL '3 hours',
        'confirmed'
    ) ON CONFLICT DO NOTHING;
INSERT INTO notifications (
        user_id,
        title,
        message,
        type,
        related_entity_type,
        related_entity_id
    )
VALUES (
        2,
        'Бронирование подтверждено',
        'Ваше бронирование Конференц-зала А успешно подтверждено',
        'success',
        'booking',
        1
    ),
    (
        3,
        'Напоминание',
        'Ваша встреча в Конференц-зале Б через 1 час',
        'info',
        'booking',
        2
    ),
    (
        4,
        'Ожидает подтверждения',
        'Ваше бронирование ожидает подтверждения администратором',
        'warning',
        'booking',
        3
    ) ON CONFLICT DO NOTHING;