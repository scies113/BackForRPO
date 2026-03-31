# Используйте ту же версию, что в go.mod
FROM golang:1.25.0

# Добавьте переменные окружения для го-модулей
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

WORKDIR /app

# Копируем только модульные файлы сначала (кэширование!)
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем
RUN go build -o main ./cmd/api

EXPOSE 8080
CMD ["./main"]

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off
ENV CGO_ENABLED=0
ENV GOOS=linux