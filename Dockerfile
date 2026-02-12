# Многоступенчатая сборка для оптимизации размера образа

# Этап 1: Сборка приложения
FROM golang:1.25.3-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Устанавливаем swag для генерации документации
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем весь исходный код
COPY . .

# Генерируем Swagger документацию
RUN swag init

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smartbooking .

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем CA сертификаты и wget для healthcheck
RUN apk --no-cache add ca-certificates wget

# Создаем рабочую директорию
WORKDIR /root/

# Создаем директорию для uploads
RUN mkdir -p /uploads

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/smartbooking .

# Копируем файлы миграций
COPY --from=builder /app/migrations ./migrations

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./smartbooking"]
