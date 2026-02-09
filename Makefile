.PHONY: help build up down restart logs clean test db-migrate

# Помощь - показать все доступные команды
help:
	@echo "Доступные команды:"
	@echo "  make build      - Собрать Docker образ"
	@echo "  make up         - Запустить все сервисы (PostgreSQL + приложение)"
	@echo "  make down       - Остановить все сервисы"
	@echo "  make restart    - Перезапустить все сервисы"
	@echo "  make logs       - Показать логи приложения"
	@echo "  make logs-db    - Показать логи PostgreSQL"
	@echo "  make clean      - Удалить все контейнеры и volumes"
	@echo "  make db-shell   - Подключиться к PostgreSQL shell"
	@echo "  make app-shell  - Подключиться к shell приложения"

# Собрать Docker образ
build:
	docker-compose build

# Запустить все сервисы
up:
	docker-compose up -d
	@echo "✅ Сервисы запущены!"
	@echo " Приложение доступно на http://localhost:8080"
	@echo " Swagger документация: http://localhost:8080/swagger/"
	@echo " PostgreSQL доступен на localhost:5432"

# Остановить все сервисы
down:
	docker-compose down

# Перезапустить все сервисы
restart:
	docker-compose restart

# Показать логи приложения
logs:
	docker-compose logs -f app

# Показать логи PostgreSQL
logs-db:
	docker-compose logs -f postgres

# Удалить все контейнеры и данные
clean:
	docker-compose down -v
	@echo "✅ Все контейнеры и данные удалены"

# Подключиться к PostgreSQL shell
db-shell:
	docker-compose exec postgres psql -U postgres -d smartbooking

# Подключиться к shell приложения
app-shell:
	docker-compose exec app sh

# Пересобрать и перезапустить
rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "✅ Приложение пересобрано и запущено!"
