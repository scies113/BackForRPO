# ==========================================
# Этап 1: Сборка (Builder)
# ==========================================
FROM golang:1.25.0-alpine AS builder

# Включаем нужные переменные ДО начала сборки. Добавлен быстрый прокси.
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.io,direct \
    GOSUMDB=off

# Устанавливаем git, так как он нужен для go mod download
RUN apk add --no-cache git

WORKDIR /app

# Сначала копируем только файлы зависимостей для кэширования слоя
COPY go.mod go.sum ./

# Используем BuildKit cache для быстрого скачивания
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Копируем остальной исходный код
COPY . .

# Собираем бинарник с кэшированием сборки.
# Флаги -ldflags="-s -w" вырезают отладочную информацию, делая бинарник еще меньше!
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o main ./cmd/api

# ==========================================
# Этап 2: Финальный образ (Production)
# ==========================================
# Берем пустой и супер-легкий alpine (весит около 5 МБ)
FROM alpine:latest

# Отключаем кэш apk (на всякий случай, если будут устанавливаться пакеты)
RUN rm -rf /var/cache/apk/*

WORKDIR /root/

# Копируем ТОЛЬКО готовый бинарник и статику из первого этапа
COPY --from=builder /app/main .
COPY --from=builder /app/frontend ./frontend
COPY --from=builder /app/.env.example ./.env

# Обозначаем порт
EXPOSE 8080

# Запускаем
CMD ["./main"]