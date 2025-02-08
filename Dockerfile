# Используем официальный образ Go в качестве базового
FROM golang:1.23-alpine AS builder

# Установим необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы приложения
COPY auth-service .

# Сборка приложения
RUN go build -o main ./cmd/server

# Минимальный образ для запуска
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем конфигурационные файлы
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/main .
COPY --from=builder /app/docs /app/docs
COPY .env .env
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations

# Экспонируем порт
EXPOSE 8000

# Устанавливаем команду запуска
CMD ["./main"]