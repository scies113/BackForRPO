# 📘 Полный гайд по GoBackendFootball

**Версия:** 3.0.0  
**Дата:** 20 апреля 2026 г.

Этот документ — **полное руководство** по запуску, тестированию и использованию бэкенд-сервиса.

---

## 📋 Оглавление

1. [Требования](#требования)
2. [Быстрый старт](#быстрый-старт)
3. [Запуск через Docker](#запуск-через-docker)
4. [Локальный запуск](#локальный-запуск)
5. [Тестирование API](#тестирование-api)
6. [JWT-токены](#jwt-токены)
7. [Работа с куки](#работа-с-куки)
8. [База данных](#база-данных)
9. [Частые ошибки](#частые-ошибки)
10. [Разработка](#разработка)
11. [Деплой](#деплой)

---

## 🛠 Требования

### Обязательные

| Программа | Версия | Зачем |
|-----------|--------|-------|
| **Go** | 1.21+ | Язык программирования |
| **Docker** | 20.10+ | Контейнеризация БД |
| **Docker Compose** | 2.0+ | Управление контейнерами |
| **Git** | 2.0+ | Система контроля версий |

### Опциональные

| Программа | Версия | Зачем |
|-----------|--------|-------|
| **Postman** | Любая | Тестирование API |
| **cURL** | Любая | Тестирование из CLI |
| **pgAdmin** | Любая | Управление БД |
| **VS Code** | Любая | Редактор кода |

---

## 🚀 Быстрый старт

### Шаг 1: Клонирование репозитория

```bash
git clone <URL_РЕПОЗИТОРИЯ>
cd GoBackendFootball
```

### Шаг 2: Настройка переменных окружения

Создай файл `.env` в корне проекта:

```env
# База данных
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234567890
DB_NAME=football_stats
DB_PORT=5433

# Сервер
PORT=8080

# JWT
JWT_SECRET=your-secret-key-change-in-production
```

### Шаг 3: Запуск базы данных

```bash
docker-compose up -d db
```

### Шаг 4: Запуск сервера

```bash
go run cmd/api/main.go
```

### Шаг 5: Проверка

```bash
curl http://localhost:8080/api/matches
```

**Ожидаемый ответ:**
```json
{"error": "Token required"}
```

✅ **Готово!** Сервер работает.

---

## 🐳 Запуск через Docker

### Вариант 1: Только база данных (рекомендуется для разработки)

```bash
# Запустить БД
docker-compose up -d db

# Проверить статус
docker-compose ps

# Остановить БД
docker-compose down
```

### Вариант 2: База данных + бэкенд

```bash
# Запустить всё
docker-compose up --build

# Запустить в фоновом режиме
docker-compose up -d

# Посмотреть логи
docker-compose logs -f backend

# Остановить всё
docker-compose down
```

### Вариант 3: Полное пересоздание

```bash
# Удалить контейнеры и тома (БД будет очищена!)
docker-compose down -v

# Запустить заново
docker-compose up --build
```

---

## 💻 Локальный запуск

### Шаг 1: Установка Go

1. Скачай с https://go.dev/dl/
2. Установи (версия 1.21 или новее)
3. Проверь:
   ```bash
   go version
   ```

### Шаг 2: Установка зависимостей

```bash
cd GoBackendFootball
go mod download
```

### Шаг 3: Запуск базы данных

```bash
docker-compose up -d db
```

### Шаг 4: Настройка .env

Проверь, что файл `.env` существует:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234567890
DB_NAME=football_stats
DB_PORT=5433
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

### Шаг 5: Запуск сервера

```bash
# Windows (PowerShell)
$env:DB_HOST="localhost"; $env:DB_PORT="5433"; $env:DB_PASSWORD="1234567890"; go run cmd/api/main.go

# Windows (cmd)
set DB_HOST=localhost
set DB_PORT=5433
set DB_PASSWORD=1234567890
go run cmd/api/main.go

# Linux/Mac
export DB_HOST=localhost
export DB_PORT=5433
export DB_PASSWORD=1234567890
go run cmd/api/main.go
```

### Шаг 6: Сборка бинарника

```bash
# Собрать
go build -o bin/api.exe cmd/api/main.go

# Запустить
./bin/api.exe
```

---

## 🧪 Тестирование API

### Через cURL

#### 1. Регистрация нового пользователя

```bash
curl -X POST http://localhost:8080/api/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}"
```

**Ожидаемый ответ:**
```json
{
  "message": "User registered",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "role": "user"
}
```

#### 2. Вход в систему

```bash
curl -X POST http://localhost:8080/api/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
```

**Ожидаемый ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 3. Получить все матчи (с токеном)

```bash
curl -X GET http://localhost:8080/api/matches ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Ожидаемый ответ:**
```json
[]
```

#### 4. Получить все матчи (без токена)

```bash
curl -X GET http://localhost:8080/api/matches
```

**Ожидаемый ответ:**
```json
{"error": "Token required"}
```

#### 5. Создать матч (только admin/operator)

```bash
curl -X POST http://localhost:8080/api/matches ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." ^
  -H "Content-Type: application/json" ^
  -d "{\"home_team\":\"Arsenal\",\"away_team\":\"Chelsea\",\"match_date\":\"2026-04-05T15:00:00Z\"}"
```

**Ожидаемый ответ:**
```json
{"message": "Матч создан"}
```

#### 6. Получить прогноз ИИ

```bash
curl -X POST http://localhost:8080/api/predict/1 ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Ожидаемый ответ:**
```json
{
  "match_id": 1,
  "prediction": {
    "home_win": 45.2,
    "draw": 30.1,
    "away_win": 24.7
  },
  "model_version": "v1.0",
  "generated_at": "2026-03-31T12:00:00Z"
}
```

#### 7. Получить журнал аудита (только admin)

```bash
curl -X GET http://localhost:8080/api/audit ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Ожидаемый ответ:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "user_name": "admin",
    "action": "CREATE",
    "entity": "Match",
    "entity_id": 1,
    "timestamp": "2026-03-31T12:00:00Z"
  }
]
```

---

### Через Postman

#### Шаг 1: Создание коллекции

1. Открой Postman
2. Создай новую коллекцию "Football API"
3. Добавь переменную коллекции `base_url` = `http://localhost:8080`

#### Шаг 2: Создание запросов

**Регистрация:**
- Method: `POST`
- URL: `{{base_url}}/api/register`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**Логин:**
- Method: `POST`
- URL: `{{base_url}}/api/login`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

**Получить матчи:**
- Method: `GET`
- URL: `{{base_url}}/api/matches`
- Headers: `Authorization: Bearer {{token}}`

#### Шаг 3: Сохранение токена

После логина скопируй токен из ответа и сохрани в переменную `token`.

Или используй тесты Postman:
```javascript
// В тесте логина
const response = pm.response.json();
pm.collectionVariables.set('token', response.token);
```

---

### Через браузер (с куки)

#### Шаг 1: Логин через fetch

Открой консоль браузера (F12) и выполни:

```javascript
fetch('http://localhost:8080/api/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({
    email: 'test@example.com',
    password: 'password123'
  })
})
.then(r => r.json())
.then(console.log);
```

#### Шаг 2: Запрос с куки

```javascript
fetch('http://localhost:8080/api/matches', {
  credentials: 'include'
})
.then(r => r.json())
.then(console.log);
```

**Важно:** Браузер автоматически отправит куку с токеном.

---

## 🎫 JWT-токены

### Что такое JWT?

**JWT (JSON Web Token)** — способ безопасной передачи данных между клиентом и сервером.

**Структура:**
```
Header.Payload.Signature
```

**Пример токена:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Iml2YW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NzUwNDQwNjUsImlhdCI6MTc3NDk1NzY2NX0.qmIzhlEZ3RWL6_1Q-F2D8JfsLE0ob6KLjTG24Eix5VA
```

### Из чего состоит токен

#### 1. Header (Заголовок)

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

**Что означает:**
- `alg: HS256` — алгоритм подписи HMAC-SHA256
- `typ: JWT` — тип токена

#### 2. Payload (Данные)

```json
{
  "user_id": 1,
  "username": "ivan",
  "role": "admin",
  "exp": 1775044065,
  "iat": 1774957665
}
```

**Что означает:**
- `user_id` — ID пользователя в БД
- `username` — имя пользователя
- `role` — роль (admin, operator, analyst, user)
- `exp` — срок действия (Unix timestamp)
- `iat` — время выдачи (Unix timestamp)

#### 3. Signature (Подпись)

```
Signature = HMAC-SHA256(
    base64UrlEncode(header) + "." + base64UrlEncode(payload),
    JWT_SECRET
)
```

**Зачем:** Гарантия, что токен не был подделан.

---

### Как работает токен в моём проекте

#### Шаг 1: Генерация токена

**Файл:** `internal/middleware/auth.go`

```go
func GenerateToken(user *model.User, roleName string) (string, error) {
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     roleName,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

**Что происходит:**
1. Создаём claims с данными пользователя
2. Устанавливаем срок действия (24 часа)
3. Создаём JWT с алгоритмом HS256
4. Подписываем токен секретным ключом
5. Возвращаем строку токена

---

#### Шаг 2: Сохранение токена

**Вариант A: Куки (рекомендуется)**

```go
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
    c.SetCookie(
        "token",     // имя куки
        token,       // значение
        86400,       // срок (24 часа в секундах)
        "/api",      // путь
        "",          // домен
        false,       // Secure (false для localhost)
        true,        // HttpOnly (JavaScript не имеет доступа)
    )
}
```

**Преимущества:**
- ✅ HttpOnly — защита от XSS
- ✅ Автоматическая отправка с запросами
- ✅ Не нужно добавлять заголовок вручную

---

**Вариант B: Заголовок Authorization**

Токен возвращается в JSON-ответе:

```json
{
  "message": "User registered",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "role": "user"
}
```

Клиент сохраняет токен (localStorage, sessionStorage) и отправляет с запросами:

```javascript
fetch('/api/matches', {
  headers: {
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIs...'
  }
});
```

---

#### Шаг 3: Проверка токена

**Файл:** `internal/middleware/auth.go`

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var tokenString string

        // Пробуем заголовок Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                tokenString = parts[1]
            }
        }

        // Пробуем куки
        if tokenString == "" {
            if cookie, err := c.Cookie("token"); err == nil {
                tokenString = cookie
            }
        }

        // Токен не найден
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
            return
        }

        // Проверка токена
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }

        // Сохраняем данные в контекст
        c.Set("userID", claims.UserID)
        c.Set("userName", claims.Username)
        c.Set("userRole", claims.Role)

        c.Next()
    }
}
```

**Что происходит:**
1. Извлекаем токен (из заголовка или куки)
2. Расшифровываем подпись с помощью `JWT_SECRET`
3. Проверяем, совпадает ли подпись
4. Проверяем срок действия (`exp`)
5. Извлекаем `user_id`, `username`, `role`
6. Сохраняем в контекст Gin (доступно в хендлерах)

---

### Как расшифровать токен

#### Способ 1: Сайт jwt.io

1. Скопируй токен
2. Открой https://jwt.io
3. Вставь токен в поле "Encoded"
4. Увидишь decoded данные

**Пример:**
```json
{
  "user_id": 1,
  "username": "ivan",
  "role": "admin",
  "exp": 1775044065,
  "iat": 1774957665
}
```

⚠️ **Важно:** Не вставляй реальные токены с продакшена на публичные сайты!

---

#### Способ 2: Код на Go

```go
package main

import (
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

func decodeToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte("JWT_SECRET"), nil
    })
    return claims, err
}

func main() {
    token := "eyJhbGciOiJIUzI1NiIs..."
    claims, err := decodeToken(token)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("User ID: %d, Username: %s, Role: %s\n", 
        claims.UserID, claims.Username, claims.Role)
}
```

---

#### Способ 3: Код на JavaScript

```javascript
function decodeJWT(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
}

const token = "eyJhbGciOiJIUzI1NiIs...";
const claims = decodeJWT(token);
console.log(claims);
// { user_id: 1, username: "ivan", role: "admin", exp: 1775044065, iat: 1774957665 }
```

---

### Срок действия токена

**Текущий срок:** 24 часа

**Как изменить:**

**Файл:** `internal/middleware/auth.go`

```go
// Было (24 часа):
ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),

// Стало (7 дней):
ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),

// Стало (1 час):
ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
```

---

### Частые проблемы с токенами

#### 1. `{"error": "Token required"}`

**Причина:** Токен не передан в запросе.

**Решение:**
```bash
# Добавь заголовок
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# Или используй куки
curl -b cookies.txt
```

---

#### 2. `{"error": "Invalid or expired token"}`

**Причина:** Токен истёк или повреждён.

**Решение:**
- Залогинься заново
- Проверь `JWT_SECRET` в `.env`

---

#### 3. Куки не устанавливается

**Причина:** CORS не настроен или неправильный домен.

**Решение:**
```go
// В main.go
c.Header("Access-Control-Allow-Credentials", "true")
c.Header("Access-Control-Expose-Headers", "Set-Cookie")

// В auth_handler.go
c.SetCookie("token", token, 86400, "/api", "", false, true)
```

---

## 🍪 Работа с куки

### Как установить куку

**Сервер (Go):**
```go
c.SetCookie(
    "token",     // имя
    token,       // значение
    86400,       // срок (секунды)
    "/api",      // путь
    "",          // домен
    false,       // Secure
    true,        // HttpOnly
)
```

### Как отправить куку

**Клиент (JavaScript):**
```javascript
fetch('http://localhost:8080/api/matches', {
    credentials: 'include' // Обязательно!
});
```

**Клиент (cURL):**
```bash
# Сохранить куки
curl -c cookies.txt -X POST http://localhost:8080/api/login ...

# Использовать куки
curl -b cookies.txt -X GET http://localhost:8080/api/matches
```

### Настройки куки

| Параметр | Значение | Описание |
|----------|----------|----------|
| `Name` | `token` | Имя куки |
| `Value` | JWT-токен | Значение куки |
| `Max-Age` | `86400` | Срок жизни (24 часа) |
| `Path` | `/api` | Путь (кука отправляется только на /api) |
| `Domain` | `` (пусто) | Домен (пустой = текущий хост) |
| `Secure` | `false` | Только HTTPS (true для продакшена) |
| `HttpOnly` | `true` | JavaScript не имеет доступа |

---

## 🗄 База данных

### Подключение

**Файл:** `.env`

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234567890
DB_NAME=football_stats
DB_PORT=5433
```

### Как подключиться к БД

#### Через Docker

```bash
# Запустить БД
docker-compose up -d db

# Подключиться через psql
docker exec -it gobackendfootball-db-1 psql -U postgres -d football_stats

# Посмотреть таблицы
\dt

# Выйти
\q
```

#### Через pgAdmin

1. Открой pgAdmin
2. Создай новое подключение:
   - Host: `localhost`
   - Port: `5433`
   - Database: `football_stats`
   - Username: `postgres`
   - Password: `1234567890`

### Просмотр данных

```sql
-- Посмотреть всех пользователей
SELECT id, username, email, role_id FROM users;

-- Посмотреть все матчи
SELECT * FROM matches;

-- Посмотреть прогнозы
SELECT * FROM predictions;

-- Посмотреть журнал аудита
SELECT * FROM audit_logs ORDER BY timestamp DESC;

-- Посмотреть роли
SELECT * FROM roles;
```

### Очистка БД

```sql
-- Удалить все матчи
DELETE FROM matches;

-- Удалить всех пользователей (кроме админа)
DELETE FROM users WHERE username != 'admin';

-- Удалить журнал аудита
DELETE FROM audit_logs;

-- Сбросить автоинкремент ID
ALTER SEQUENCE matches_id_seq RESTART WITH 1;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
```

---

## ❌ Частые ошибки

### 1. Порт 8080 занят

**Ошибка:**
```
listen tcp :8080: bind: Only one usage of each socket address is normally permitted
```

**Решение:**
```bash
# Найти процесс
netstat -ano | findstr :8080

# Убить процесс
taskkill /F /PID <PID>

# Или изменить порт в .env
PORT=8081
```

---

### 2. БД не подключается

**Ошибка:**
```
failed to connect to database: dial tcp: lookup localhost: no such host
```

**Решение:**
```bash
# Проверь, что БД запущена
docker-compose ps

# Перезапусти БД
docker-compose down
docker-compose up -d db

# Проверь .env
DB_HOST=localhost
DB_PORT=5433
```

---

### 3. Токен не работает

**Ошибка:**
```
{"error": "Invalid or expired token"}
```

**Решение:**
- Залогинься заново
- Проверь `JWT_SECRET` в `.env` (должен совпадать)
- Проверь срок действия токена (24 часа)

---

### 4. Куки не устанавливается

**Ошибка:**
```
Set-Cookie заголовок отсутствует в ответе
```

**Решение:**
```go
// Проверь CORS в main.go
c.Header("Access-Control-Allow-Credentials", "true")
c.Header("Access-Control-Expose-Headers", "Set-Cookie")

// Проверь setAuthCookie в auth_handler.go
c.SetCookie("token", token, 86400, "/api", "", false, true)
```

---

### 5. Роли не работают

**Ошибка:**
```
{"error": "Insufficient permissions"}
```

**Решение:**
```sql
-- Проверь роль пользователя
SELECT u.username, r.name 
FROM users u 
JOIN roles r ON u.role_id = r.id 
WHERE u.username = 'testuser';

-- Если роль неправильная, обнови
UPDATE users SET role_id = 1 WHERE username = 'testuser';
-- 1 = admin, 2 = operator, 3 = analyst, 4 = user
```

---

## 💻 Разработка

### Запуск тестов

```bash
# Все тесты (22 теста)
go test ./internal/... -v

# Тесты конкретного пакета
go test ./internal/middleware/... -v
go test ./internal/service/... -v
go test ./internal/handler/... -v

# Тесты с покрытием
go test ./internal/... -v -cover
```

### Форматирование кода

```bash
# Отформатировать все файлы
go fmt ./...

# Проверить ошибки
go vet ./...
```

### Сборка проекта

```bash
# Собрать бинарник
go build -o bin/api.exe cmd/api/main.go

# Собрать для Linux
GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go

# Собрать для macOS
GOOS=darwin GOARCH=amd64 go build -o bin/api-macos cmd/api/main.go
```

### Добавление зависимостей

```bash
# Установить пакет
go get github.com/some/package

# Обновить зависимости
go mod tidy

# Проверить зависимости
go mod verify
```

---

## 🚀 Деплой

### Деплой на сервер (Ubuntu)

#### Шаг 1: Подготовка сервера

```bash
# Установить Go
sudo apt update
sudo apt install golang-go

# Установить Docker
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker

# Установить Docker Compose
sudo apt install docker-compose
```

#### Шаг 2: Клонирование проекта

```bash
git clone <URL_РЕПОЗИТОРИЯ>
cd GoBackendFootball
```

#### Шаг 3: Настройка .env

```bash
nano .env
```

```env
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=<СЛОЖНЫЙ_ПАРОЛЬ>
DB_NAME=football_stats
DB_PORT=5432
PORT=8080
JWT_SECRET=<СЛОЖНЫЙ_СЕКРЕТ>
```

#### Шаг 4: Запуск

```bash
docker-compose up -d --build
```

#### Шаг 5: Проверка

```bash
# Посмотреть логи
docker-compose logs -f backend

# Проверить статус
docker-compose ps

# Протестировать API
curl http://localhost:8080/api/matches
```

---

### Деплой через Docker Hub

#### Шаг 1: Сборка образа

```bash
docker build -t yourusername/go-backend-football:latest .
```

#### Шаг 2: Публикация

```bash
docker login
docker push yourusername/go-backend-football:latest
```

#### Шаг 3: Использование на сервере

```bash
docker pull yourusername/go-backend-football:latest

docker run -d \
  -p 8080:8080 \
  --env-file .env \
  --link db \
  yourusername/go-backend-football:latest
```

---

## 📊 Мониторинг

### Логи сервера

```bash
# Docker Compose
docker-compose logs -f backend

# Локально
# Логи пишутся в stdout
```

### Метрики

| Метрика | Как посмотреть |
|---------|----------------|
| RPS (Requests Per Second) | Nginx logs, Prometheus |
| Время отклика | Логи Gin, Jaeger |
| Ошибки | Логи, Sentry |
| Использование памяти | `docker stats`, `top` |
| Подключения к БД | `pg_stat_activity` |

---

## 📚 Дополнительные ресурсы

- [Документация Go](https://go.dev/doc/)
- [Документация Gin](https://gin-gonic.com/docs/)
- [Документация GORM](https://gorm.io/docs/)
- [JWT.io](https://jwt.io/) — расшифровать токен
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Docker Docs](https://docs.docker.com/)

---

## 📞 Поддержка

Если возникли вопросы:

1. Проверь этот гайд
2. Посмотри логи (`docker-compose logs`)
3. Проверь `.env`
4. Убедись, что БД запущена

---

**Версия документа:** 3.0.0  
**Дата:** 20 апреля 2026 г.
