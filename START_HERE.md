# 🚀 Быстрый старт SmartBooking

## Что было исправлено

✅ **Фронтенд работает стабильно** - исправлена проблема с API_URL в admin.js и owner-dashboard.js
✅ **Добавлены тестовые данные** - 3 owner пользователя, ресурсы с категориями
✅ **Добавлены placeholder фотографии** - для всех ресурсов
✅ **Добавлена навигация** - Admin Panel и Owner Dashboard теперь доступны из меню
✅ **Добавлены отзывы** - тестовые отзывы для ресурсов

## Шаги для запуска

### 1. Запустите Docker Desktop
Убедитесь что Docker Desktop запущен на вашем компьютере.

### 2. Очистите старые данные (если были)
```bash
docker-compose down -v
```

### 3. Запустите проект
```bash
docker-compose up -d --build
```

Ожидайте 30-60 секунд пока все сервисы запустятся и миграции применятся.

### 4. Проверьте что всё запущено
```bash
docker-compose ps
```

Все сервисы должны быть в статусе "Up".

### 5. Откройте приложение
- **Фронтенд**: http://localhost
- **Backend API**: http://localhost:8080
- **Swagger Docs**: http://localhost:8080/swagger/
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
- **pgAdmin**: http://localhost:5050 (admin@smartbooking.com/admin)

## Тестовые аккаунты

### Администратор
- Email: `admin@smartbooking.com`
- Password: `password123`
- Доступ: Admin Panel

### Владельцы (Owner)
- Email: `owner1@smartbooking.com` / Password: `password123`
- Email: `owner2@smartbooking.com` / Password: `password123`
- Email: `owner3@smartbooking.com` / Password: `password123`
- Доступ: Owner Dashboard

### Обычные пользователи
- Email: `john@example.com` / Password: `password123`
- Email: `jane@example.com` / Password: `password123`

## Что показать на защите

### 1. Главная страница (/)
- Категории ресурсов (Баня, Бассейн, Спортивная площадка, и т.д.)
- Красивый дизайн с иконками

### 2. Страница ресурсов (/resources.html)
- Список всех ресурсов с фотографиями
- Фильтры по категориям и городам
- Возможность бронирования

### 3. Вход как обычный пользователь (john@example.com)
- Забронировать ресурс
- Посмотреть "Мои брони" (/bookings.html)
- Отменить бронирование

### 4. Вход как Owner (owner1@smartbooking.com)
- Owner Dashboard (/owner-dashboard.html)
- Статистика: ресурсы, бронирования, доход, рейтинг
- Список своих ресурсов
- Список бронирований для своих ресурсов

### 5. Вход как Admin (admin@smartbooking.com)
- Admin Panel (/admin.html)
- Overview: статистика системы
- Управление всеми бронированиями
- Управление всеми ресурсами
- Управление пользователями
- Управление категориями

### 6. API Documentation
- Swagger UI: http://localhost:8080/swagger/
- 35+ endpoints
- RESTful API

## Особенности проекта

### Backend (Go)
- Clean Architecture (Handler → Service → Repository)
- PostgreSQL база данных
- MinIO (S3-compatible) для хранения фото
- JWT-подобная аутентификация (упрощенная)
- Middleware для логирования
- Background worker для статистики

### Frontend (Vanilla JS)
- Роль-based навигация
- Admin Panel для администраторов
- Owner Dashboard для владельцев
- Адаптивный дизайн
- Динамическая загрузка данных

### Database
- 8 миграций
- Категории, ресурсы, фото, бронирования
- Расписание, тарифы, отзывы
- Связи owner → resources

### Features
- Multi-role система (admin, owner, user)
- Booking система с проверкой доступности
- Photo upload с MinIO/S3
- Reviews и рейтинги
- Ценообразование (базовая цена + тарифы)
- Расписание работы ресурсов

## Troubleshooting

### Порты заняты
Если порты 80, 8080, 5432, 9000, 9001 или 5050 заняты:
```bash
# Проверьте что использует порты
lsof -i :80
lsof -i :8080
lsof -i :5432

# Остановите конфликтующие сервисы или измените порты в docker-compose.yml
```

### Фронтенд не работает
1. Проверьте что nginx контейнер запущен: `docker-compose ps`
2. Проверьте логи: `docker-compose logs nginx`
3. Убедитесь что бэкенд доступен: `curl http://localhost:8080/health`

### Backend ошибки
```bash
# Посмотрите логи
docker-compose logs app

# Проверьте подключение к БД
docker-compose logs postgres
```

### Миграции не применились
```bash
# Пересоздайте с чистыми volumes
docker-compose down -v
docker-compose up -d --build
```

## Команды для защиты

```bash
# Посмотреть статистику в логах (каждые 30 сек)
docker-compose logs -f app

# Проверить БД
docker-compose exec postgres psql -U postgres -d smartbooking -c "SELECT COUNT(*) FROM users;"
docker-compose exec postgres psql -U postgres -d smartbooking -c "SELECT COUNT(*) FROM resources;"
docker-compose exec postgres psql -U postgres -d smartbooking -c "SELECT COUNT(*) FROM bookings;"

# Посмотреть все таблицы
docker-compose exec postgres psql -U postgres -d smartbooking -c "\dt"
```

## Структура проекта
```
assignment3/
├── main.go                    # Entry point
├── config/                    # Configuration
├── internal/
│   ├── handler/              # HTTP handlers (controllers)
│   ├── service/              # Business logic
│   ├── repository/           # Database access
│   ├── models/               # Data models
│   ├── middleware/           # HTTP middleware
│   ├── storage/              # S3/MinIO storage
│   └── database/             # DB connection
├── migrations/               # SQL migrations (8 files)
├── frontend/                 # Frontend (HTML/CSS/JS)
│   ├── index.html           # Homepage
│   ├── resources.html       # Resources list
│   ├── bookings.html        # User bookings
│   ├── owner-dashboard.html # Owner dashboard
│   ├── admin.html           # Admin panel
│   └── js/
│       ├── app.js           # Main JS
│       ├── admin.js         # Admin panel JS
│       └── owner-dashboard.js # Owner dashboard JS
├── docker-compose.yml        # Docker services
└── nginx/                    # Nginx config

5 Services:
- app (Go backend)
- postgres (Database)
- minio (S3 storage)
- nginx (Frontend)
- pgadmin (DB admin tool)
```

## Успешной защиты! 🎓
