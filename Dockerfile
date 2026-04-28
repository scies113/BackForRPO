# ==========================================
# Этап 1: Сборка (Builder)
# ==========================================
FROM golang:1.25.0-alpine AS builder

# Включаем нужные переменные ДО начала сборки
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=off

WORKDIR /app

# Сначала копируем только файлы зависимостей для кэширования слоя
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код
COPY . .

# Собираем бинарник. 
# Флаги -ldflags="-s -w" вырезают отладочную информацию, делая бинарник еще меньше!
RUN go build -ldflags="-s -w" -o main ./cmd/api

# ==========================================
# Этап 2: Финальный образ (Production)
# ==========================================
# Берем пустой и супер-легкий alpine (весит около 5 МБ)
FROM alpine:latest

WORKDIR /root/

# Копируем ТОЛЬКО готовый бинарник и статику из первого этапа
COPY --from=builder /app/main .
COPY --from=builder /app/frontend ./frontend
COPY --from=builder /app/.env.example ./.env

# Обозначаем порт
EXPOSE 8080

# Запускаем
CMD ["./main"]