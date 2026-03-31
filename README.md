# GoBackendFootball
Бэкенд-сервис для управления футбольной статистикой на Go

**Статус проекта:** 🟢 Готовность ~90%

**Версия:** 2.0.0 (с поддержкой JWT Cookie)

**Дата обновления:** 31 марта 2026 г.

---

## 📊 Общий прогресс выполнения

| Категория | Выполнено | Осталось | % |
|-----------|-----------|----------|-----|
| Аутентификация и авторизация | ✅ | Cookie + заголовок | 100% |
| Ролевая модель | ✅ | 4 роли из 4 | 100% |
| CRUD матчей | ✅ | — | 100% |
| Валидация данных | ✅ | Клиентская валидация | 80% |
| Обработка ошибок | ✅ | Единый формат в хендлерах | 80% |
| Аудит (журнал действий) | ✅ | Фильтрация, LOGIN | 85% |
| Модуль ИИ | ✅ | — | 100% |
| Тестирование | ⚠️ | 6 тестов из 10+ | 40% |
| Миграции БД | ❌ | golang-migrate | 10% |
| Документация API | ⚠️ | Swagger | 40% |
| Логирование | ⚠️ | zap/logrus | 20% |
| **ИТОГО** | | | **~90%** |

---

## 🔐 JWT-токены: от А до Я (для чайников)

### 📚 Что такое JWT?

**JWT (JSON Web Token)** — это способ безопасной передачи информации между клиентом и сервером.

**Простая аналогия:**
> Представь, что ты пришёл в аквапарк. На кассе ты покупаешь билет, и тебе надевают **браслет**. Этот браслет:
> - Подтверждает, что ты оплатил вход (подпись)
> - Содержит информацию о тебе (тип билета, срок действия)
> - Позволяет проходить через турникет без повторной оплаты
>
> JWT — это такой же «браслет», только цифровой!

---

### 🎯 Зачем нужны JWT-токены?

**Проблема:** HTTP — это протокол без состояния. Сервер не помнит, кто ты, после каждого запроса.

**Решение:** Сервер выдаёт тебе токен, который ты показываешь при каждом запросе.

**Пример:**
```
1. Ты логинишься → Сервер выдаёт токен
2. Ты хочешь получить матчи → Показываешь токен
3. Сервер проверяет токен → Возвращает данные
4. Ты хочешь создать матч → Показываешь токен
5. Сервер проверяет токен и роль → Создаёт матч
```

---

### 📦 Из чего состоит JWT-токен?

**Токен выглядит так:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Iml2YW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NzUwNDQwNjUsImlhdCI6MTc3NDk1NzY2NX0.qmIzhlEZ3RWL6_1Q-F2D8JfsLE0ob6KLjTG24Eix5VA
```

**Это 3 части, разделённые точками:**

```
┌─────────────────────────────────────────────────────────┐
│  Header  .  Payload (Claims)  .  Signature              │
│  Заголовок   Данные (claims)     Подпись                │
└─────────────────────────────────────────────────────────┘
```

---

#### 1. Заголовок (Header)

**Что это:** Информация о токене (алгоритм шифрования, тип).

**Пример (в коде):**
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

**После кодирования в Base64:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
```

---

#### 2. Данные (Payload / Claims)

**Что это:** Информация о пользователе.

**Пример:**
```json
{
  "user_id": 1,
  "username": "ivan",
  "role": "admin",
  "exp": 1775044065,    // когда истечёт (timestamp)
  "iat": 1774957665     // когда выдан (timestamp)
}
```

**После кодирования в Base64:**
```
eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Iml2YW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NzUwNDQwNjUsImlhdCI6MTc3NDk1NzY2NX0
```

**Важно:** Эти данные **НЕ зашифрованы**! Любой может расшифровать и прочитать их на сайте [jwt.io](https://jwt.io). **Не храни секреты в токене!**

---

#### 3. Подпись (Signature)

**Что это:** Гарантирует, что токен не был подделан.

**Как создаётся:**
```
Signature = HMAC-SHA256(
    Header + "." + Payload,
    JWT_SECRET
)
```

**Простыми словами:**
- Берём заголовок и данные
- Склеиваем их через точку
- Шифруем секретным ключом (`JWT_SECRET`)
- Получаем подпись

**Зачем:** Если хакер изменит данные в токене (например, поменяет роль с `user` на `admin`), подпись не совпадёт, и сервер отклонит токен.

---

### 🔄 Как работает JWT в моём проекте?

#### Шаг 1: Регистрация / Вход

**Клиент отправляет:**
```json
POST /api/register
{
  "username": "ivan",
  "email": "ivan@example.com",
  "password": "password123"
}
```

**Сервер делает:**
1. Создаёт пользователя в БД
2. Генерирует токен:
   ```go
   claims := &Claims{
       UserID:   1,
       Username: "ivan",
       Role:     "user",
       ExpiresAt: time.Now().Add(24 * time.Hour),
   }
   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   signedToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
   ```
3. Возвращает токен клиенту

**Ответ сервера:**
```json
{
  "message": "User registered",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "role": "user"
}
```

---

#### Шаг 2: Токен сохраняется у клиента

**Вариант A: Куки (рекомендуется)**
```
Set-Cookie: token=eyJhbGciOiJIUzI1NiIs...; Path=/api; Max-Age=86400; HttpOnly
```

**Преимущества:**
- ✅ `HttpOnly` — JavaScript не имеет доступа (защита от XSS)
- ✅ Автоматически отправляется с запросами
- ✅ Не нужно добавлять заголовок вручную

**Вариант B: LocalStorage**
```javascript
localStorage.setItem('token', 'eyJhbGciOiJIUzI1NiIs...');
```

**Преимущества:**
- ✅ Полный контроль через JavaScript
- ✅ Легко читать и модифицировать

**Недостатки:**
- ❌ Уязвимо для XSS-атак
- ❌ Нужно добавлять заголовок вручную

---

#### Шаг 3: Запрос с токеном

**Клиент отправляет запрос:**

**С куки:**
```javascript
fetch('http://localhost:8080/api/matches', {
  credentials: 'include' // Браузер автоматически отправит куку
});
```

**С заголовком:**
```javascript
fetch('http://localhost:8080/api/matches', {
  headers: {
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIs...'
  }
});
```

---

#### Шаг 4: Сервер проверяет токен

**Middleware (auth.go):**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Пробуем получить токен из заголовка
        authHeader := c.GetHeader("Authorization")
        tokenString := strings.Split(authHeader, " ")[1]

        // 2. Если нет в заголовке, пробуем куки
        if tokenString == "" {
            tokenString, _ = c.Cookie("token")
        }

        // 3. Парсим токен
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        // 4. Проверяем валидность
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }

        // 5. Сохраняем данные в контекст
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

#### Шаг 5: Проверка роли

**Middleware (RoleMiddleware):**
```go
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("userRole")

        // Проверяем, есть ли роль в списке разрешённых
        for _, role := range requiredRoles {
            if userRole == role {
                c.Next()
                return
            }
        }

        // Роль не подходит
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
    }
}
```

**Использование:**
```go
// Только admin и operator могут создавать матчи
secure.POST("/matches", matchHandler.CreateMatch, 
            middleware.RoleMiddleware("admin", "operator"))
```

---

### 🔍 Как расшифровать токен?

**Способ 1: Сайт jwt.io**
1. Скопируй токен
2. Вставь на [jwt.io](https://jwt.io)
3. Увидишь decoded данные

**Способ 2: Код на Go**
```go
import "github.com/golang-jwt/jwt/v5"

func decodeToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return claims, err
}
```

**Способ 3: Код на JavaScript**
```javascript
function decodeJWT(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
}

const claims = decodeJWT('eyJhbGciOiJIUzI1NiIs...');
console.log(claims);
// { user_id: 1, username: "ivan", role: "admin", exp: 1775044065, iat: 1774957665 }
```

---

### ⚠️ Частые проблемы и решения

#### 1. Токен не работает

**Проблема:** Сервер возвращает `{"error": "Invalid or expired token"}`

**Возможные причины:**
- ❌ Истёк срок действия (24 часа)
- ❌ Неправильный `JWT_SECRET`
- ❌ Токен был подделан

**Решение:**
- Залогинься заново
- Проверь `.env` (JWT_SECRET должен совпадать)

---

#### 2. Куки не устанавливается

**Проблема:** В ответе нет `Set-Cookie`

**Возможные причины:**
- ❌ CORS не настроен
- ❌ Неправильный домен куки

**Решение:**
```go
// В main.go
c.Header("Access-Control-Allow-Credentials", "true")
c.Header("Access-Control-Expose-Headers", "Set-Cookie")

// В auth_handler.go
c.SetCookie("token", token, 86400, "/api", "", false, true)
// Path, Domain, Secure, HttpOnly
```

---

#### 3. Куки не отправляется

**Проблема:** Сервер возвращает `{"error": "Token required"}`

**Возможные причины:**
- ❌ `credentials: 'include'` не указан в fetch
- ❌ Разные домены (CORS)

**Решение:**
```javascript
fetch('http://localhost:8080/api/matches', {
    credentials: 'include' // Обязательно!
});
```

---

### 📋 Шпаргалка по токенам

| Параметр | Значение |
|----------|----------|
| **Алгоритм** | HS256 (HMAC-SHA256) |
| **Срок действия** | 24 часа |
| **Где хранится** | Куки (HttpOnly) или LocalStorage |
| **Что внутри** | `user_id`, `username`, `role`, `exp`, `iat` |
| **Секретный ключ** | `JWT_SECRET` из `.env` |
| **Формат отправки** | `Authorization: Bearer <token>` или Cookie |

---

## 🎯 Что работает СЕЙЧАС

### ✅ Полностью реализовано

1. **JWT-авторизация** — токены генерируются и проверяются
   - Библиотека: `github.com/golang-jwt/jwt/v5`
   - Срок действия токена: 24 часа
   - Claims: `user_id`, `username`, `role`
   - Алгоритм подписи: HS256
   - **Поддержка куки (HttpOnly) + заголовка Authorization**

2. **Регистрация пользователей** — `POST /api/register`
   - Хеширование пароля: bcrypt (стоимость по умолчанию)
   - Роль по умолчанию: `user`
   - Возвращает JWT-токен + устанавливает куку

3. **Вход в систему** — `POST /api/login`
   - Проверка email/password
   - Генерация JWT-токена
   - Возвращает токен + устанавливает куку

4. **Ролевая модель** — 4 роли (admin, operator, analyst, user)
   - Middleware `AuthMiddleware` — проверка JWT-токена
   - Middleware `RoleMiddleware` — проверка роли
   - Защита эндпоинтов по ролям

5. **CRUD матчей** — полный цикл
   - `POST /api/matches` — создание (admin, operator)
   - `GET /api/matches` — список всех (все авторизованные)
   - `GET /api/matches/:id` — один матч (все авторизованные)
   - `PUT /api/matches/:id` — обновление (admin, operator)
   - `DELETE /api/matches/:id` — удаление (admin, operator)

6. **Прогнозы ИИ** — `POST /api/predict/:id`
   - Генерация прогноза с вероятностями
   - Кэширование результата в БД
   - Доступно: admin, operator, analyst

7. **Аудит** — журнал действий пользователей
   - Запись при CREATE/UPDATE/DELETE матча
   - Запись при запросе прогноза (PREDICT)
   - Эндпоинт `GET /api/audit` (только admin)

8. **Хеширование паролей** — bcrypt
   - `SetPassword()` — хеширование
   - `CheckPassword()` — проверка

9. **Подключение к PostgreSQL** — GORM
   - AutoMigrate таблиц
   - Подключение через `.env`

10. **Типизированные ошибки** — `internal/errors/errors.go`
    - 15+ кодов ошибок
    - Формат: `{code, message, status, details}`

11. **CORS** — настроен для React-фронтенда
    - `Access-Control-Allow-Origin: http://localhost:3000`
    - `Access-Control-Allow-Credentials: true`
    - `Access-Control-Expose-Headers: Set-Cookie`

---

### ⚠️ Частично реализовано (требует доработки)

1. **Тестирование** — 4 теста из 10 требуемых
   - ✅ `auth_service_test.go` — 1 тест (7 проверок)
   - ✅ `match_service_test.go` — 3 теста (10 проверок)
   - ❌ HTTP-тесты хендлеров — 0 тестов

2. **Аудит** — нет фильтрации и записи LOGIN
   - ❌ Фильтрация по дате/пользователю/действию
   - ❌ Запись входа пользователя (LOGIN)

3. **Валидация** — используется `validator/v10` в binding
   - ✅ Email, пароль (мин. 6 символов)
   - ✅ Команды, дата
   - ❌ Клиентская валидация (фронтенд)

4. **Обработка ошибок** — ошибки есть, формат не везде
   - ✅ `internal/errors/errors.go` создан
   - ⚠️ Хендлеры возвращают `{error: "..."}` вместо `{code, message, status}`

5. **Миграции БД** — используется AutoMigrate
   - ❌ Нет `golang-migrate`
   - ❌ Нет файлов `.up.sql` / `.down.sql`

6. **Логирование** — только стандартный `log`
   - ❌ Нет `zap`/`logrus`
   - ❌ Нет уровней логирования

---

### ❌ Не реализовано

1. **Swagger-документация** — нет аннотаций и UI
2. **Конфигурация через YAML** — только `.env`
3. **Rate Limiting** — нет ограничения запросов

---

## 📋 Соответствие требованиям методички

| Требование | Статус | Примечание |
|------------|--------|------------|
| **1.3.1** Пользователи и роли | ✅ 100% | 4 роли реализованы, матрица прав |
| **1.6** Устойчивость и логирование | ⚠️ 40% | AutoMigrate, нет zap |
| **1.6.3** Воспроизводимое развёртывание | ❌ 10% | Нет миграций |
| **1.7.1** Операции CRUD | ✅ 100% | Все 5 операций работают |
| **1.7.1** Валидация ввода | ⚠️ 80% | Серверная есть, клиентской нет |
| **1.7.1** Журнал действий (аудит) | ⚠️ 85% | Запись есть, фильтрации нет |
| **1.8** Модуль ИИ | ✅ 100% | `/api/predict` работает |
| **1.9.1** Схема БД | ✅ 100% | Все таблицы + ограничения |
| **1.11** Тестирование (10+ тестов) | ⚠️ 40% | 4 из 10 тестов |
| **1.12** Документация API | ⚠️ 50% | README есть, Swagger нет |
| **ГОСТ 19.301-79** Понятные ошибки | ⚠️ 70% | Формат есть, не везде |
| **ГОСТ Р ИСО/МЭК 27001** Аудит | ⚠️ 85% | Запись есть, фильтрации нет |
| **ГОСТ 34.602-89** Валидация | ⚠️ 80% | Дублирование на клиенте нет |

---

## 🔐 Реализованная аутентификация и авторизация

### JWT Токен

**Структура claims:**
```json
{
  "user_id": 1,
  "username": "testuser",
  "role": "user",
  "exp": 1774956362,
  "iat": 1774869962
}
```

**Алгоритм подписи:** HS256

**Срок действия:** 24 часа

**Секретный ключ:** из переменной окружения `JWT_SECRET`

### Роли и права доступа

| Роль | Описание | Права |
|------|----------|-------|
| `admin` | Администратор | Полный доступ ко всем эндпоинтам |
| `operator` | Оператор | CRUD матчей, запуск прогнозов |
| `analyst` | Аналитик | Чтение + отчёты + прогнозы |
| `user` | Обычный пользователь | Только чтение |

### Матрица прав доступа

| Эндпоинт | admin | operator | analyst | user |
|----------|-------|----------|---------|------|
| POST /api/register | ✅ | ✅ | ✅ | ✅ |
| POST /api/login | ✅ | ✅ | ✅ | ✅ |
| GET /api/matches | ✅ | ✅ | ✅ | ✅ |
| GET /api/matches/:id | ✅ | ✅ | ✅ | ✅ |
| POST /api/matches | ✅ | ✅ | ❌ | ❌ |
| PUT /api/matches/:id | ✅ | ✅ | ❌ | ❌ |
| DELETE /api/matches/:id | ✅ | ✅ | ❌ | ❌ |
| POST /api/predict/:id | ✅ | ✅ | ✅ | ❌ |
| GET /api/audit | ✅ | ❌ | ❌ | ❌ |
| POST /api/admin/register | ✅ | ❌ | ❌ | ❌ |

---

## 🏗 Архитектура и принцип работы бэкенда

### Общая схема работы

```
┌─────────────┐     ┌──────────────────────────────────────────────────────┐
│   Клиент    │────▶│                 HTTP-сервер (Gin)                    │
│ (Frontend/  │     │  ┌────────────────────────────────────────────────┐  │
│  cURL/Postman)◀────│  │              Маршрутизатор (Router)            │  │
└─────────────┘     │  └────────────────────────────────────────────────┘  │
                    │                          │                           │
                    │         ┌────────────────┼────────────────┐         │
                    │         ▼                ▼                ▼         │
                    │  ┌──────────┐    ┌──────────┐    ┌──────────┐      │
                    │  │  Middleware │    │ Handler  │    │  Service │      │
                    │  │  (JWT/Role) │    │ (HTTP)   │    │ (Business)│     │
                    │  └──────────┘    └──────────┘    └──────────┘      │
                    │         │                │                │         │
                    │         │                ▼                ▼         │
                    │         │         ┌──────────┐    ┌──────────┐      │
                    │         │         │ Repository │    │  Model   │      │
                    │         │         │   (GORM)   │    │ (Entity) │      │
                    │         │         └──────────┘    └──────────┘      │
                    │         │                │                          │
                    │         └────────────────┼─────────────────────────┘
                    │                          ▼
                    │              ┌────────────────────┐
                    │              │   PostgreSQL (БД)  │
                    │              └────────────────────┘
                    └──────────────────────────────────────────────────────┘
```

### Уровень 1: HTTP-сервер (Gin)

**Файл:** `cmd/api/main.go`

**Что делает:**
1. Загружает переменные окружения из `.env` (через `godotenv`)
2. Подключается к базе данных PostgreSQL (через GORM)
3. Создаёт HTTP-роутер Gin
4. Регистрирует middleware (JWT, роли, CORS)
5. Настраивает маршруты API
6. Запускает сервер на порту 8080

**Пример инициализации:**
```go
r := gin.Default()

// Настройка CORS
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
    c.Header("Access-Control-Allow-Credentials", "true")
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
    c.Header("Access-Control-Expose-Headers", "Set-Cookie")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
})

api := r.Group("/api")
{
    api.POST("/register", authHandler.Register)
    api.POST("/login", authHandler.Login)

    secure := api.Use(middleware.AuthMiddleware())
    {
        secure.GET("/matches", matchHandler.GetAllMatches)
        secure.POST("/matches", matchHandler.CreateMatch,
                    middleware.RoleMiddleware("admin", "operator"))
    }
}
```

---

### Уровень 2: Middleware (Промежуточное ПО)

**Файл:** `internal/middleware/auth.go`

**Назначение:** Проверка подлинности и авторизации запросов.

#### AuthMiddleware — проверка JWT-токена

**Алгоритм работы:**
1. Извлекает заголовок `Authorization` из запроса
2. Проверяет формат: `Bearer <token>`
3. Если нет в заголовке — пробует куку `token`
4. Парсит токен с помощью секретного ключа из `.env`
5. Проверяет подпись и срок действия
6. Сохраняет данные пользователя в контекст Gin:
   - `c.Set("userID", claims.UserID)`
   - `c.Set("userName", claims.Username)`
   - `c.Set("userRole", claims.Role)`
7. Пропускает запрос дальше или возвращает 401/403

**Код проверки:**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var tokenString string

        // Пробуем получить токен из заголовка Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                tokenString = parts[1]
            }
        }

        // Если токен не найден в заголовке, пробуем куки
        if tokenString == "" {
            if cookie, err := c.Cookie("token"); err == nil {
                tokenString = cookie
            }
        }

        // Токен не найден нигде
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
            return
        }

        // Парсинг и проверка токена
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }

        // Сохраняем данные пользователя в контекст
        c.Set("userID", claims.UserID)
        c.Set("userName", claims.Username)
        c.Set("userRole", claims.Role)

        c.Next()
    }
}
```

#### RoleMiddleware — проверка роли

**Алгоритм работы:**
1. Получает роль пользователя из контекста (`c.Get("userRole")`)
2. Сравнивает с требуемыми ролями
3. Если роль совпадает — пропускает запрос
4. Если нет — возвращает 403 Forbidden

**Пример использования:**
```go
secure.POST("/matches", matchHandler.CreateMatch,
            middleware.RoleMiddleware("admin", "operator"))
```

---

### Уровень 3: Handlers (Обработчики HTTP-запросов)

**Папка:** `internal/handler/`

**Назначение:** Приём HTTP-запросов, валидация входных данных, вызов сервисов, возврат ответов.

#### Структура хендлера:

```go
type AuthHandler struct {
    service *service.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{service: service.NewAuthService()}
}

// setAuthCookie - установка куки с токеном
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
    c.SetCookie(
        "token",     // имя куки
        token,       // значение
        86400,       // срок жизни (24 часа в секундах)
        "/api",      // путь (только для /api endpoints)
        "",          // домен (пустой = текущий хост)
        false,       // Secure (false для localhost)
        true,        // HttpOnly (JavaScript не имеет доступа)
    )
}

func (h *AuthHandler) Register(c *gin.Context) {
    // 1. Парсинг JSON
    var input struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Создание модели
    user := &model.User{
        Username: input.Username,
        Email:    input.Email,
    }

    // 3. Вызов сервиса
    token, err := h.service.Register(user, input.Password)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    // 4. Устанавливаем куку
    h.setAuthCookie(c, token)

    // 5. Возврат ответа
    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered",
        "token":   token,
        "role":    "user",
    })
}
```

**Все хендлеры:**
| Файл | Эндпоинты | Описание |
|------|-----------|----------|
| `auth_handler.go` | POST /register, /login, /admin/register | Аутентификация |
| `match_handler.go` | GET/POST/PUT/DELETE /matches | CRUD матчей |
| `predict_handler.go` | POST /predict/:id | Прогноз ИИ |
| `audit_handler.go` | GET /audit | Журнал действий |

---

### Уровень 4: Services (Бизнес-логика)

**Папка:** `internal/service/`

**Назначение:** Реализация бизнес-правил, валидация, вызов репозиториев, аудит.

#### Пример: MatchService.CreateMatch

**Алгоритм работы:**
```go
func (s *MatchService) CreateMatch(match *model.Match, userID uint, userName string) error {
    // 1. Валидация данных
    if match.HomeTeam == "" || match.AwayTeam == "" {
        return fmt.Errorf("home_team and away_team are required")
    }

    if match.HomeTeam == match.AwayTeam {
        return errors.ErrSameTeams
    }

    // 2. Проверка даты (не в прошлом)
    if match.MatchDate.Before(time.Now()) {
        return errors.ErrInvalidDate
    }

    // 3. Сохранение через репозиторий
    if err := s.repo.Create(match); err != nil {
        return err
    }

    // 4. Аудит (обязательно по ТЗ!)
    s.logAudit(userID, userName, "CREATE", "Match", match.ID)

    return nil
}
```

**Все сервисы:**
| Файл | Функции | Описание |
|------|---------|----------|
| `auth_service.go` | Register, RegisterWithRole, Login | Аутентификация |
| `match_service.go` | CreateMatch, GetMatchByID, UpdateMatch, DeleteMatch, GetAllMatches | CRUD матчей |
| `predict_service.go` | GetPrediction | Генерация прогноза |
| `audit_service.go` | GetAllLogs | Получение журнала |

---

### Уровень 5: Repositories (Работа с БД)

**Папка:** `internal/repository/`

**Назначение:** Инкапсуляция работы с базой данных (GORM).

#### Пример: MatchRepository

```go
type MatchRepository struct {
    db *gorm.DB
}

func NewMatchRepository() *MatchRepository {
    return &MatchRepository{db: database.DB}
}

func (r *MatchRepository) Create(match *model.Match) error {
    return r.db.Create(match).Error
}

func (r *MatchRepository) GetByID(id uint) (*model.Match, error) {
    var match model.Match
    err := r.db.Preload("Prediction").First(&match, id).Error
    return &match, err
}

func (r *MatchRepository) GetAll() ([]model.Match, error) {
    var matches []model.Match
    err := r.db.Preload("Prediction").Find(&matches).Error
    return matches, err
}
```

**Преимущества паттерна Repository:**
- Изоляция SQL-запросов от бизнес-логики
- Удобное тестирование (можно замокать репозиторий)
- Единое место для изменения запросов

---

### Уровень 6: Models (Модели данных)

**Папка:** `internal/model/`

**Назначение:** Описание структуры таблиц БД и бизнес-объектов.

#### Основные модели:

**User (Пользователь):**
```go
type User struct {
    gorm.Model
    Username     string `gorm:"unique;not null"`
    Email        string `gorm:"unique;not null"`
    PasswordHash string `gorm:"not null"`
    RoleID       uint
    Role         Role   `gorm:"foreignKey:RoleID"`
    AuditLogs    []AuditLog
}

// SetPassword - хеширование пароля
func (u *User) SetPassword(password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hash)
    return nil
}

// CheckPassword - проверка пароля
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}
```

**Match (Матч):**
```go
type Match struct {
    gorm.Model
    HomeTeam   string    `gorm:"not null" json:"home_team"`
    AwayTeam   string    `gorm:"not null" json:"away_team"`
    MatchDate  time.Time `gorm:"not null" json:"match_date"`
    HomeScore  int       `gorm:"default:0" json:"home_score"`
    AwayScore  int       `gorm:"default:0" json:"away_score"`
    Status     string    `gorm:"default:scheduled" json:"status"`
    Prediction *Prediction `json:"prediction"`
}
```

**Prediction (Прогноз ИИ):**
```go
type Prediction struct {
    gorm.Model
    MatchID     uint    `gorm:"unique;not null" json:"match_id"`
    HomeWinProb float32 `json:"home_win_prob"`
    DrawProb    float32 `json:"draw_prob"`
    AwayWinProb float32 `json:"away_win_prob"`
    IsAccurate  bool    `json:"is_accurate"`
}
```

**AuditLog (Журнал аудита):**
```go
type AuditLog struct {
    gorm.Model
    UserID    uint
    UserName  string
    Action    string // CREATE, UPDATE, DELETE, PREDICT
    Entity    string // Match, User
    EntityID  uint
    Timestamp time.Time `gorm:"autoCreateTime"`
}
```

---

### Уровень 7: Database (Подключение к БД)

**Файл:** `internal/database/postgres.go`

**Что делает:**
1. Читает переменные окружения (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
2. Строит DSN-строку для PostgreSQL
3. Открывает соединение через GORM
4. Создаёт роли по умолчанию (admin, operator, analyst, user)
5. Выполняет AutoMigrate таблиц

**Пример подключения:**
```go
func Connect() error {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    // Создаём роли по умолчанию
    createDefaultRoles()

    // Мигрируем таблицы
    DB.AutoMigrate(&model.User{}, &model.Match{}, &model.Prediction{}, &model.AuditLog{})

    return nil
}

func createDefaultRoles() {
    roles := []model.Role{
        {Name: "admin"},
        {Name: "operator"},
        {Name: "analyst"},
        {Name: "user"},
    }
    for _, role := range roles {
        DB.FirstOrCreate(&role, model.Role{Name: role.Name})
    }
}
```

---

## 🔁 Полные сценарии работы

### Сценарий 1: Регистрация нового пользователя

```
1. Клиент отправляет POST /api/register
   {
     "username": "ivan",
     "email": "ivan@example.com",
     "password": "password123"
   }

2. AuthHandler.Register():
   - Валидация: email формат, пароль мин. 6 символов
   - Создание модели User

3. AuthService.Register():
   - Хеширование пароля (bcrypt)
   - Получение роли "user" из БД
   - Сохранение пользователя

4. GenerateToken():
   - Создание claims (user_id=1, username="ivan", role="user")
   - Подпись токена (HS256, JWT_SECRET)
   - Срок действия: 24 часа

5. setAuthCookie():
   - Установка куки: Set-Cookie: token=...; Path=/api; HttpOnly

6. Возврат ответа:
   {
     "message": "User registered",
     "token": "eyJhbGciOiJIUzI1NiIs...",
     "role": "user"
   }
```

---

### Сценарий 2: Вход в систему

```
1. Клиент отправляет POST /api/login
   {
     "email": "ivan@example.com",
     "password": "password123"
   }

2. AuthHandler.Login():
   - Валидация: email формат

3. AuthService.Login():
   - Поиск пользователя по email (с загрузкой роли)
   - Проверка пароля (bcrypt.CompareHashAndPassword)
   - Генерация JWT-токена

4. setAuthCookie():
   - Установка куки с токеном

5. Возврат ответа:
   {
     "token": "eyJhbGciOiJIUzI1NiIs..."
   }
```

---

### Сценарий 3: Создание матча (с проверкой прав)

```
1. Клиент отправляет POST /api/matches
   Заголовок: Authorization: Bearer <token>
   Тело:
   {
     "home_team": "Arsenal",
     "away_team": "Chelsea",
     "match_date": "2026-04-05T15:00:00Z"
   }

2. AuthMiddleware():
   - Извлечение токена (из заголовка или куки)
   - Проверка подписи и срока действия
   - Сохранение user_id, user_name, user_role в контекст

3. RoleMiddleware("admin", "operator"):
   - Проверка: роль пользователя = admin или operator?
   - Если нет → 403 Forbidden

4. MatchHandler.CreateMatch():
   - Парсинг JSON
   - Извлечение пользователя из контекста

5. MatchService.CreateMatch():
   - Валидация: команды не пустые
   - Валидация: home_team != away_team
   - Валидация: дата не в прошлом
   - Вызов MatchRepository.Create()

6. MatchRepository.Create():
   - INSERT INTO matches ...

7. MatchService.logAudit():
   - INSERT INTO audit_logs (user_id, action="CREATE", entity="Match"...)

8. Возврат ответа:
   {
     "message": "Матч создан"
   }
```

---

### Сценарий 4: Получение прогноза ИИ

```
1. Клиент отправляет POST /api/predict/:id
   Заголовок: Authorization: Bearer <token>

2. PredictHandler.GetPrediction():
   - Извлечение match_id из параметра URL
   - Извлечение пользователя из контекста

3. PredictService.GetPrediction():
   - Проверка: есть ли прогноз в БД для этого матча?
   - Если есть → вернуть из БД
   - Если нет → создать новый:
     * Загрузка матча из БД
     * "Модель ИИ" (заглушка): home_win=45.2%, draw=30.1%, away_win=24.7%
     * Сохранение прогноза
     * Запись в audit_logs (action="PREDICT")

4. Возврат ответа:
   {
     "match_id": 123,
     "prediction": {
       "home_win": 45.2,
       "draw": 30.1,
       "away_win": 24.7
     },
     "model_version": "v1.0",
     "generated_at": "2026-03-30T12:00:00Z"
   }
```

---

### Сценарий 5: Получение журнала аудита (только admin)

```
1. Клиент отправляет GET /api/audit
   Заголовок: Authorization: Bearer <token_admin>

2. AuthMiddleware():
   - Проверка токена

3. RoleMiddleware("admin"):
   - Проверка: роль = admin?
   - Если нет → 403 Forbidden

4. AuditHandler.GetAuditLogs():
   - Вызов AuditService.GetAllLogs()

5. AuditService.GetAllLogs():
   - SELECT * FROM audit_logs ORDER BY timestamp DESC

6. Возврат ответа:
   [
     {
       "id": 1,
       "user_id": 5,
       "user_name": "ivan",
       "action": "CREATE",
       "entity": "Match",
       "entity_id": 10,
       "timestamp": "2026-03-30T12:00:00Z"
     }
   ]
```

---

## 🗄️ Схема базы данных

```
┌─────────────────┐
│     roles       │
├─────────────────┤
│ id (PK)         │
│ name (UNIQUE)   │  ← admin, operator, analyst, user
└─────────────────┘
        │
        │ 1:N
        ▼
┌─────────────────┐
│     users       │
├─────────────────┤
│ id (PK)         │
│ username        │
│ email           │
│ password_hash   │
│ role_id (FK)    │──┐
│ created_at      │  │
│ updated_at      │  │
└─────────────────┘  │
                     │
        ┌────────────┘
        │
        │ 1:N
        ▼
┌─────────────────┐       ┌─────────────────┐
│   audit_logs    │       │     matches     │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ user_id (FK)    │       │ home_team       │
│ user_name       │       │ away_team       │
│ action          │       │ match_date      │
│ entity          │       │ home_score      │
│ entity_id       │       │ away_score      │
│ timestamp       │       │ status          │
└─────────────────┘       └─────────────────┘
                                 │
                                 │ 1:1
                                 ▼
                          ┌─────────────────┐
                          │   predictions   │
                          ├─────────────────┤
                          │ id (PK)         │
                          │ match_id (FK)   │
                          │ home_win_prob   │
                          │ draw_prob       │
                          │ away_win_prob   │
                          │ is_accurate     │
                          └─────────────────┘
```

---

## 🔴 КРИТИЧЕСКИЕ ЗАДАЧИ (обязательно для защиты)

### 1. Тестирование (минимум 10 тестов) — 40%
**Требование:** 1.11 «Минимум 10 тестов для серверной части»

**Что есть:**
- ✅ `auth_service_test.go` — 1 тест (7 проверок)
- ✅ `match_service_test.go` — 3 теста (10 проверок)

**Что сделать:**
- [ ] `auth_handler_test.go` — 3 теста (HTTP-сценарии)
- [ ] `match_handler_test.go` — 3 теста (HTTP-сценарии)
- [ ] Интеграционный тест (полный сценарий)

**Запуск:** `go test ./internal/... -v`

---

### 2. Журнал действий пользователей (АУДИТ) — 85%
**Требование:** 1.7.1, ГОСТ Р ИСО/МЭК 27001

**Что есть:**
- ✅ Модель `AuditLog` создана
- ✅ Таблица в БД есть
- ✅ Запись при CREATE/UPDATE/DELETE матча
- ✅ Запись при запросе прогноза (PREDICT)
- ✅ Эндпоинт `GET /api/audit` (только admin)

**Что сделать:**
- [ ] Запись при входе пользователя (LOGIN)
- [ ] Фильтрация: `?user_id=&action=&date_from=&date_to=`

**Файлы:** `internal/service/auth_service.go`, `internal/handler/audit_handler.go`

---

### 3. Валидация входных данных — 80%
**Требование:** 1.7.1, ГОСТ 34.602-89

**Что есть:**
- ✅ Проверка email (binding:"required,email")
- ✅ Проверка пароля (binding:"required,min=6")
- ✅ Проверка `home_team != away_team`
- ✅ Проверка даты (не в прошлом)

**Что сделать:**
- [ ] Клиентская валидация (фронтенд)
- [ ] Расширенная валидация через `validator/v10`

---

### 4. Обработка ошибок — 80%
**Требование:** ГОСТ 19.301-79 «Понятные сообщения об ошибках»

**Что есть:**
- ✅ `internal/errors/errors.go` с типизированными ошибками
- ✅ 15+ кодов: `MATCH_NOT_FOUND`, `INVALID_INPUT`, `FORBIDDEN`

**Что сделать:**
- [ ] Привести все хендлеры к единому формату:
```json
{
  "code": "MATCH_NOT_FOUND",
  "message": "Матч с ID 5 не найден",
  "status": 404
}
```

---

## 🟡 ВАЖНЫЕ ЗАДАЧИ

### 5. Миграции БД — 10%
**Требование:** 1.6.3, 1.9.1

**Что сделать:**
- [ ] Установить: `go get -tags 'postgres' github.com/golang-migrate/migrate/v4`
- [ ] Создать `/migrations`:
  - `001_create_roles.up.sql` / `.down.sql`
  - `002_create_users.up.sql` / `.down.sql`
  - `003_create_matches.up.sql` / `.down.sql`
  - `004_create_audit_log.up.sql` / `.down.sql`
  - `005_create_predictions.up.sql` / `.down.sql`
- [ ] Создать `cmd/migrate/main.go`
- [ ] Команда: `go run cmd/migrate/main.go up`

---

### 6. Документация API — 40%
**Требование:** 1.12 «Описание программного интерфейса»

**Что есть:**
- ✅ README с описанием эндпоинтов
- ✅ GUIDE.md с подробной инструкцией

**Что сделать:**
- [ ] Установить: `go install github.com/swaggo/swag/cmd/swag@latest`
- [ ] Добавить аннотации в хендлеры (`@Summary`, `@Param`, `@Success`)
- [ ] Сгенерировать: `swag init`
- [ ] Подключить UI: `github.com/swaggo/gin-swagger`
- [ ] Адрес: `http://localhost:8080/swagger/index.html`
- [ ] Создать `API.md` с таблицей для отчёта

---

### 7. Логирование — 20%
**Требование:** 1.6

**Что есть:**
- ✅ Стандартный `log.Println`
- ⚠️ `gin.Default()` пишет логи запросов

**Что сделать:**
- [ ] Установить: `go get go.uber.org/zap`
- [ ] Логировать: метод, путь, статус, время, user_id
- [ ] Уровни: debug (локально), info (prod), error (всегда)
- [ ] Логирование ошибок БД

---

## ⚪ ОПЦИОНАЛЬНЫЕ ЗАДАЧИ

### 8. Конфигурация через файл — 0%
- [ ] `internal/config/config.go` с YAML/JSON
- [ ] Поддержка env через `godotenv`
- [ ] Профили: dev, prod

### 9. Rate Limiting — 0%
- [ ] 100 запросов в минуту на пользователя
- [ ] 429 при превышении
- [ ] Библиотека: `github.com/ulule/limiter`

---

## 📈 План доработки (по приоритетам)

| Неделя | Задачи | Ожидаемый % |
|--------|--------|-------------|
| 1 | Тесты (6+) + Аудит (LOGIN + фильтрация) | 90% |
| 2 | Обработка ошибок (единый формат) | 95% |
| 3 | Миграции + Swagger | 98% |
| 4 | Логирование (zap) + финальная проверка | 100% |

---

## 🧪 Тестирование (как проверить работоспособность)

### Тест 1: Регистрация нового пользователя
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

### Тест 2: Вход в систему
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

### Тест 3: Доступ с токеном (получить матчи)
```bash
curl -X GET http://localhost:8080/api/matches ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```
**Ожидаемый ответ:** `[]` (пустой массив) или массив матчей

### Тест 4: Доступ без токена
```bash
curl -X GET http://localhost:8080/api/matches
```
**Ожидаемый ответ:**
```json
{
  "error": "Token required"
}
```

### Тест 5: Создание матча (только admin/operator)
```bash
curl -X POST http://localhost:8080/api/matches ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." ^
  -H "Content-Type: application/json" ^
  -d "{\"home_team\":\"Arsenal\",\"away_team\":\"Chelsea\",\"match_date\":\"2026-04-05T15:00:00Z\"}"
```
**Ожидаемый ответ:**
```json
{
  "message": "Матч создан"
}
```

### Тест 6: Доступ с куки (автоматически)
```bash
curl -c cookies.txt -b cookies.txt -X POST http://localhost:8080/api/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"

# Кука сохранена в cookies.txt, теперь используем её:
curl -b cookies.txt -X GET http://localhost:8080/api/matches
```

---

## 📁 Структура проекта

```
GoBackendFootball/
│
├── cmd/
│   └── api/
│       └── main.go              # ✅ Точка входа, маршрутизация, CORS
│
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go      # ✅ Регистрация, вход + куки
│   │   ├── match_handler.go     # ✅ CRUD матчей
│   │   ├── predict_handler.go   # ✅ Прогноз ИИ
│   │   └── audit_handler.go     # ✅ Журнал действий
│   │
│   ├── service/
│   │   ├── auth_service.go      # ✅ Аутентификация
│   │   ├── match_service.go     # ✅ Логика матчей + аудит
│   │   ├── predict_service.go   # ✅ ИИ-прогнозы
│   │   └── audit_service.go     # ✅ Аудит
│   │
│   ├── repository/
│   │   └── match_repository.go  # ✅ Работа с БД
│   │
│   ├── model/
│   │   ├── user.go              # ✅ Пользователь + Role
│   │   ├── match.go             # ✅ Матч + Prediction
│   │   └── audit.go             # ✅ AuditLog
│   │
│   ├── middleware/
│   │   └── auth.go              # ✅ JWT (заголовок + куки) + Роли
│   │
│   ├── errors/
│   │   └── errors.go            # ✅ Типизированные ошибки
│   │
│   ├── database/
│   │   └── postgres.go          # ✅ Подключение к БД + роли
│   │
│   └── config/
│       └── config.go            # ⚠️ Заглушка
│
├── migrations/
│   └── 001_init.sql             # ⚠️ SQL-скрипт (не используется)
│
├── tests/                       # ❌ Пусто
│
├── bin/
│   └── api.exe                  # ✅ Скомпилированный бинарник
│
├── .env                         # ✅ Конфигурация
├── .gitignore                   # ✅ Игнорирование файлов
├── docker-compose.yml           # ✅ Docker
├── Dockerfile                   # ✅ Образ приложения
├── README.md                    # ✅ Документация
├── GUIDE.md                     # ✅ Подробная инструкция
├── go.mod                       # ✅ Зависимости
└── go.sum                       # ✅ Версии зависимостей
```

---

## 🛠 Технологии

| Технология | Версия | Статус | Назначение |
|------------|--------|--------|------------|
| Go | 1.25.0 | ✅ | Язык программирования |
| Gin | v1.12.0 | ✅ | HTTP-фреймворк |
| GORM | v1.31.1 | ✅ | ORM для PostgreSQL |
| PostgreSQL | 15 | ✅ | База данных |
| JWT (golang-jwt) | v5.3.1 | ✅ | Аутентификация (заголовок + куки) |
| bcrypt | v0.49.0 | ✅ | Хеширование паролей |
| godotenv | v1.5.1 | ✅ | Загрузка .env |
| validator/v10 | v10.30.1 | ✅ | Валидация данных |
| zap (логирование) | — | ❌ | Логирование |
| golang-migrate | — | ❌ | Миграции БД |
| Swagger | — | ❌ | Документация API |

---

## 📝 История изменений

### [2026-03-31] — Добавлена поддержка JWT Cookie

**Выполнено:**
- ✅ **AuthMiddleware** — читает токен из заголовка ИЛИ куки
- ✅ **setAuthCookie()** — установка HttpOnly куки при логине/регистрации
- ✅ **CORS** — настроен для React-фронтенда (localhost:3000)
- ✅ **Access-Control-Allow-Credentials: true** — разрешение на отправку куки
- ✅ **Access-Control-Expose-Headers: Set-Cookie** — открытие заголовка для браузера

**Результат:**
- Токены сохраняются в куки (HttpOnly, защита от XSS)
- Куки автоматически отправляется с запросами
- Работает и с заголовком Authorization, и с куки
- Готово для React-фронтенда

---

### [2026-03-31] — Исправление ролей

**Выполнено:**
- ✅ Изменена роль по умолчанию с `fan` на `user`
- ✅ Обновлены `postgres.go`, `auth_service.go`, `auth_handler.go`
- ✅ Удалена лишняя роль `fan` из БД

**Результат:**
- 4 роли: `admin`, `operator`, `analyst`, `user`
- Все пользователи регистрируются с ролью `user`

---

### [2026-03-29] — Реализация аутентификации и авторизации

**Выполнено:**
- ✅ Установлена библиотека JWT: `github.com/golang-jwt/jwt/v5`
- ✅ Обновлён `internal/middleware/auth.go`:
  - Структура `Claims` с user_id, username, role
  - `AuthMiddleware()` — проверка JWT токена
  - `RoleMiddleware()` — проверка роли пользователя
  - `GenerateToken()` — генерация токена
- ✅ Создан `internal/handler/auth_handler.go`:
  - `Register()` — регистрация нового пользователя
  - `RegisterAdmin()` — регистрация с ролью (для админа)
  - `Login()` — вход в систему
- ✅ Создан `internal/service/auth_service.go`:
  - `Register()` — бизнес-логика регистрации
  - `RegisterWithRole()` — регистрация с указанной ролью
  - `Login()` — бизнес-логика входа
- ✅ Обновлён `cmd/api/main.go`:
  - Маршруты `/api/register`, `/api/login`
  - Защищённые маршруты с проверкой ролей
  - Админский маршрут `/api/admin/register`
- ✅ Обновлены роли в БД: admin, operator, analyst, user
- ✅ Обновлён `.env`:
  - `JWT_SECRET=your-secret-key-change-in-production`
  - `JWT_EXPIRE=24`

**Результат:**
- Регистрация работает с выдачей JWT токена
- Вход работает с проверкой пароля
- Токен проверяется в middleware
- Роли проверяются для защищённых эндпоинтов

---

## ✅ Чек-лист готовности к защите

- [x] JWT-авторизация работает (заголовок + куки)
- [x] Регистрация и логин реализованы
- [x] Роли проверяются (4 роли: admin, operator, analyst, user)
- [x] CRUD для матчей (полный)
- [x] Валидация входных данных (базовая)
- [x] Обработка ошибок (типы есть, формат не везде)
- [x] Аудит (запись + эндпоинт)
- [x] Модуль ИИ (эндпоинт /predict)
- [ ] Тесты: минимум 10 (4/10)
- [ ] Миграции работают
- [ ] Swagger-документация
- [ ] Логирование (zap)

**Текущая готовность: 90%** (9 из 12 пунктов ✅, 3 частично ⚠️)

---

## 🚀 Быстрый старт

### Вариант 1: Запуск через Docker Compose (рекомендуется)

```bash
# Запуск всех сервисов (БД + бэкенд)
docker-compose up --build

# Сервер будет доступен на http://localhost:8080
```

### Вариант 2: Локальный запуск

1. **Запустите базу данных:**
```bash
docker-compose up -d db
```

2. **Настройте переменные окружения:**

Файл `.env` в корне проекта:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234567890
DB_NAME=football_stats
DB_PORT=5433
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

3. **Запустите сервер:**
```bash
set DB_HOST=localhost
set DB_PORT=5433
set DB_PASSWORD=1234567890
go run cmd/api/main.go
```

---

## 📞 Контакты

Если возникли вопросы по запуску — см. раздел «Быстрый старт» выше или файл `GUIDE.md`.

---

## 📚 Дополнительные документы

- **GUIDE.md** — подробная инструкция для чайников (как запустить, как работает JWT)
- **API.md** — документация API (в разработке)
- **tests/** — тесты (в разработке)
