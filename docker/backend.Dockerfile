# Используем многоступенчатую сборку для оптимизации размера образа
FROM golang:1.20-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем только директорию backend
COPY backend/ ./backend/
COPY uploads/ ./uploads/

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./backend

# Используем минимальный образ для запуска
FROM alpine:3.18

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Копируем исполняемый файл из сборочного образа
COPY --from=builder /app/main /app/main

# Создаем директории для хранения файлов
RUN mkdir -p /app/uploads/avatars && \
    chmod -R 777 /app/uploads

WORKDIR /app

# Создаем пользователя для запуска приложения
RUN adduser -D -g '' appuser

# Меняем владельца директорий
RUN chown -R appuser:appuser /app/uploads

USER appuser

# Устанавливаем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"] 