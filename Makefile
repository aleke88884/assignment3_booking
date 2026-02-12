.PHONY: help start stop restart clean logs status db test rebuild

# По умолчанию показываем помощь
.DEFAULT_GOAL := help

# Помощь - показать все доступные команды
help:
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "🚀 SmartBooking - Команды для управления проектом"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "📦 Основные команды:"
	@echo "  make start      - Запустить проект (первый раз или после остановки)"
	@echo "  make stop       - Остановить проект"
	@echo "  make restart    - Перезапустить проект"
	@echo "  make clean      - Полная очистка и перезапуск с нуля"
	@echo ""
	@echo "📊 Мониторинг:"
	@echo "  make status     - Показать статус всех контейнеров"
	@echo "  make logs       - Показать логи backend"
	@echo "  make logs-all   - Показать логи всех сервисов"
	@echo ""
	@echo "🗄️  База данных:"
	@echo "  make db         - Подключиться к PostgreSQL"
	@echo "  make db-stats   - Показать статистику БД"
	@echo ""
	@echo "🔧 Разработка:"
	@echo "  make rebuild    - Пересобрать всё заново"
	@echo "  make test       - Проверить что всё работает"
	@echo ""
	@echo "🌐 После запуска откройте:"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Backend:  http://localhost:8080"
	@echo "  Swagger:  http://localhost:8080/swagger/"
	@echo ""

# Запустить проект
start:
	@echo "🚀 Запуск SmartBooking..."
	@docker-compose up -d --build
	@echo ""
	@echo "⏳ Ждём запуска сервисов (15 секунд)..."
	@sleep 15
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✅ SmartBooking запущен!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "🌐 Откройте в браузере: http://localhost:3000"
	@echo ""
	@echo "👤 Тестовые аккаунты (пароль: password123):"
	@echo "   admin@smartbooking.com  - Админ"
	@echo "   owner1@smartbooking.com - Владелец"
	@echo "   john@example.com        - Пользователь"
	@echo ""

# Остановить проект
stop:
	@echo "🛑 Остановка SmartBooking..."
	@docker-compose down
	@echo "✅ Проект остановлен"

# Перезапустить проект
restart:
	@echo "🔄 Перезапуск SmartBooking..."
	@docker-compose restart
	@echo "✅ Проект перезапущен"

# Полная очистка и перезапуск
clean:
	@echo "🧹 Очистка всех данных и перезапуск..."
	@docker-compose down -v
	@echo "🔨 Сборка и запуск с нуля..."
	@docker-compose up -d --build
	@echo ""
	@echo "⏳ Ждём запуска (15 секунд)..."
	@sleep 15
	@echo ""
	@echo "✅ Проект очищен и запущен заново!"
	@echo "🌐 Откройте: http://localhost:3000"

# Показать статус контейнеров
status:
	@echo "📊 Статус контейнеров:"
	@docker-compose ps

# Показать логи backend
logs:
	@echo "📝 Логи backend (Ctrl+C для выхода):"
	@docker-compose logs -f app

# Показать логи всех сервисов
logs-all:
	@echo "📝 Логи всех сервисов (Ctrl+C для выхода):"
	@docker-compose logs -f

# Подключиться к PostgreSQL
db:
	@echo "🗄️  Подключение к PostgreSQL..."
	@docker-compose exec postgres psql -U postgres -d smartbooking

# Показать статистику БД
db-stats:
	@echo "📊 Статистика базы данных:"
	@docker-compose exec postgres psql -U postgres -d smartbooking -c "\
		SELECT 'Users' as table_name, COUNT(*) as count FROM users \
		UNION ALL SELECT 'Resources', COUNT(*) FROM resources \
		UNION ALL SELECT 'Bookings', COUNT(*) FROM bookings \
		UNION ALL SELECT 'Categories', COUNT(*) FROM resource_categories \
		UNION ALL SELECT 'Photos', COUNT(*) FROM resource_photos \
		UNION ALL SELECT 'Reviews', COUNT(*) FROM reviews;"

# Пересобрать всё заново
rebuild:
	@echo "🔨 Пересборка проекта..."
	@docker-compose down
	@docker-compose build --no-cache
	@docker-compose up -d
	@echo ""
	@echo "⏳ Ждём запуска (15 секунд)..."
	@sleep 15
	@echo "✅ Проект пересобран и запущен!"

# Проверить что всё работает
test:
	@echo "🧪 Проверка работоспособности..."
	@echo ""
	@echo "1️⃣ Проверка контейнеров:"
	@docker-compose ps
	@echo ""
	@echo "2️⃣ Проверка backend health:"
	@curl -s http://localhost:8080/health | grep -q "ok" && echo "✅ Backend работает" || echo "❌ Backend не отвечает"
	@echo ""
	@echo "3️⃣ Проверка frontend:"
	@curl -s http://localhost:3000 | grep -q "SmartBooking" && echo "✅ Frontend работает" || echo "❌ Frontend не отвечает"
	@echo ""
	@echo "4️⃣ Проверка базы данных:"
	@docker-compose exec -T postgres psql -U postgres -d smartbooking -c "SELECT COUNT(*) FROM users;" > /dev/null 2>&1 && echo "✅ База данных работает" || echo "❌ База данных не отвечает"
	@echo ""
	@echo "Если всё ✅ - проект работает отлично!"
