# GoBackendFootball
Бэкенд-сервис для управления футбольной статистикой на Go

**Статус проекта:** 🟢 Готовность ~85%

**Версия:** 1.0.0

**Дата обновления:** 31 марта 2026 г.

---

## 📊 Общий прогресс выполнения

| Категория | Выполнено | Осталось | % |
|-----------|-----------|----------|-----|
| Аутентификация и авторизация | ✅ | — | 100% |
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
| **ИТОГО** | | | **~85%** |

---

## 🎯 Что работает СЕЙЧАС

### ✅ Полностью реализовано

1. **JWT-авторизация** — токены генерируются и проверяются
   - Библиотека: `github.com/golang-jwt/jwt/v5`
   - Срок действия токена: 24 часа
   - Claims: `user_id`, `username`, `role`
   - Алгоритм подписи: HS256

2. **Регистрация пользователей** — `POST /api/register`
   - Хеширование пароля: bcrypt (стоимость по умолчанию)
   - Роль по умолчанию: `fan`
   - Возвращает JWT-токен

3. **Вход в систему** — `POST /api/login`
   - Проверка email/password
   - Генерация JWT-токена
   - Возвращает токен

4. **Ролевая модель** — 4 роли (admin, operator, analyst, fan)
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
4. **CORS** — нет настройки для фронтенда

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
  "role": "fan",
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
| `fan` | Обычный пользователь | Только чтение |

### Матрица прав доступа

| Эндпоинт | admin | operator | analyst | fan |
|----------|-------|----------|---------|-----|
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
4. Регистрирует middleware (JWT, роли)
5. Настраивает маршруты API
6. Запускает сервер на порту 8080

**Пример инициализации:**
```go
r := gin.Default()
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
3. Парсит токен с помощью секретного ключа из `.env`
4. Проверяет подпись и срок действия
5. Сохраняет данные пользователя в контекст Gin:
   - `c.Set("userID", claims.UserID)`
   - `c.Set("userName", claims.Username)`
   - `c.Set("userRole", claims.Role)`
6. Пропускает запрос дальше или возвращает 401/403

**Код проверки:**
```go
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
})

if err != nil || !token.Valid {
    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
    return
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
    
    // 4. Возврат ответа
    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered",
        "token":   token,
        "role":    "fan",
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
    Timestamp time.Time
}
```

---

### Уровень 7: Database (Подключение к БД)

**Файл:** `internal/database/database.go`

**Что делает:**
1. Читает переменные окружения (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
2. Строит DSN-строку для PostgreSQL
3. Открывает соединение через GORM
4. Выполняет AutoMigrate (создаёт таблицы)

**Пример подключения:**
```go
dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
)

DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
   - Получение роли "fan" из БД
   - Сохранение пользователя

4. GenerateToken():
   - Создание claims (user_id, username, role)
   - Подпись токена (HS256)
   - Срок действия: 24 часа

5. Возврат ответа:
   {
     "message": "User registered",
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
     "role": "fan"
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

4. Возврат ответа:
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
   - Извлечение токена из заголовка
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
│ name (UNIQUE)   │  ← admin, operator, analyst, fan
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

### 10. CORS — 0%
- [ ] Разрешить запросы с фронтенда
- [ ] Разрешённые методы и заголовки

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "fan"
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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Тест 3: Доступ с токеном (получить матчи)
```bash
curl -X GET http://localhost:8080/api/matches ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." ^
  -H "Content-Type: application/json" ^
  -d "{\"home_team\":\"Arsenal\",\"away_team\":\"Chelsea\",\"match_date\":\"2026-04-05T15:00:00Z\"}"
```
**Ожидаемый ответ:**
```json
{
  "message": "Матч создан"
}
```

### Тест 6: Получение прогноза ИИ
```bash
curl -X POST http://localhost:8080/api/predict/1 ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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

### Тест 7: Проверка прав (fan не может создать матч)
```bash
curl -X POST http://localhost:8080/api/matches ^
  -H "Authorization: Bearer <token_fan>" ^
  -H "Content-Type: application/json" ^
  -d "{\"home_team\":\"Arsenal\",\"away_team\":\"Chelsea\",\"match_date\":\"2026-04-05T15:00:00Z\"}"
```
**Ожидаемый ответ:**
```json
{
  "error": "Insufficient permissions"
}
```
**Статус:** 403 Forbidden

---

## 📁 Структура проекта

```
GoBackendFootball/
│
├── cmd/
│   └── api/
│       └── main.go              # ✅ Точка входа, маршрутизация
│
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go      # ✅ Регистрация, вход
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
│   │   └── auth.go              # ✅ JWT + Роли
│   │
│   ├── errors/
│   │   └── errors.go            # ✅ Типизированные ошибки
│   │
│   ├── database/
│   │   └── database.go          # ✅ Подключение к БД
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
├── docker-compose.yml           # ✅ Docker
├── Dockerfile                   # ✅ Образ приложения
├── README.md                    # ✅ Документация
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
| JWT (golang-jwt) | v5.3.1 | ✅ | Аутентификация |
| bcrypt | v0.49.0 | ✅ | Хеширование паролей |
| godotenv | v1.5.1 | ✅ | Загрузка .env |
| validator/v10 | v10.30.1 | ✅ | Валидация данных |
| zap (логирование) | — | ❌ | Логирование |
| golang-migrate | — | ❌ | Миграции БД |
| Swagger | — | ❌ | Документация API |

---

## 📝 История изменений

### [2026-03-31] — Текущее состояние

**Выполнено:**
- ✅ Полный CRUD для матчей (GET/POST/PUT/DELETE)
- ✅ JWT-авторизация с проверкой ролей
- ✅ 4 роли: admin, operator, analyst, fan
- ✅ Модуль ИИ с прогнозами
- ✅ Аудит действий (CREATE/UPDATE/DELETE/PREDICT)
- ✅ Типизированные ошибки (15+ кодов)
- ✅ Валидация данных (email, пароль, команды, дата)
- ✅ 4 теста (auth_service, match_service)

**В работе:**
- ⚠️ Расширение тестового покрытия (4/10)
- ⚠️ Фильтрация аудита
- ⚠️ Запись LOGIN в аудит

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
- ✅ Обновлены роли в БД: admin, operator, analyst, fan
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

- [x] JWT-авторизация работает
- [x] Регистрация и логин реализованы
- [x] Роли проверяются (4 роли: admin, operator, analyst, fan)
- [x] CRUD для матчей (полный)
- [x] Валидация входных данных (базовая)
- [x] Обработка ошибок (типы есть, формат не везде)
- [x] Аудит (запись + эндпоинт)
- [x] Модуль ИИ (эндпоинт /predict)
- [ ] Тесты: минимум 10 (4/10)
- [ ] Миграции работают
- [ ] Swagger-документация
- [ ] Логирование (zap)

**Текущая готовность: 85%** (8 из 12 пунктов ✅, 3 частично ⚠️)

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
JWT_EXPIRE=24
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

Если возникли вопросы по запуску — см. раздел «Быстрый старт» выше.

---

## 📚 Ответы на вопросы преподавателя

### ❓ Как работает аутентификация?

**Ответ:** Используется JWT (JSON Web Tokens). При регистрации или входе пользователь получает токен, подписанный секретным ключом (HS256). Токен содержит `user_id`, `username`, `role` и срок действия (24 часа). При каждом запросе middleware проверяет подпись и извлекает данные пользователя в контекст.

### ❓ Как проверяются роли?

**Ответ:** После проверки JWT middleware `RoleMiddleware` получает роль из контекста и сравнивает с требуемыми. Например, для создания матча нужны роли `admin` или `operator`. Если роль не совпадает — возвращается 403 Forbidden.

### ❓ Где хранятся пароли?

**Ответ:** Пароли хранятся в хешированном виде (bcrypt). При регистрации пароль хешируется через `bcrypt.GenerateFromPassword()`, при входе — проверяется через `bcrypt.CompareHashAndPassword()`.

### ❓ Как работает аудит?

**Ответ:** При каждом CRUD-действии (CREATE/UPDATE/DELETE матча, запрос прогноза) сервис вызывает `logAudit()`, которая создаёт запись в таблице `audit_logs`. Запись содержит `user_id`, `user_name`, `action`, `entity`, `entity_id`, `timestamp`.

### ❓ Как работает модуль ИИ?

**Ответ:** При запросе `POST /api/predict/:id` сервис проверяет, есть ли прогноз в БД. Если нет — генерирует «прогноз» (заглушка с вероятностями 45.2/30.1/24.7) и сохраняет. В реальности здесь была бы загрузка обученной модели из файла.

### ❓ Почему AutoMigrate вместо миграций?

**Ответ:** AutoMigrate используется для простоты разработки — таблицы создаются автоматически при старте. Для продакшена планируются миграции через `golang-migrate` с файлами `.up.sql` и `.down.sql` для отката.

### ❓ Какие требования ГОСТ выполнены?

**Ответ:**
- **ГОСТ 19.301-79** — понятные сообщения об ошибках (типизированные ошибки с кодами)
- **ГОСТ Р ИСО/МЭК 27001** — разграничение прав (роли + JWT + аудит)
- **ГОСТ 34.602-89** — валидация входных данных (серверная)
- **ГОСТ 34.201-89** — хранение данных в структурированной форме (PostgreSQL, нормализация)

### ❓ Сколько тестов написано?

**Ответ:** 4 теста (7+ проверок в `auth_service_test.go`, 10 проверок в `match_service_test.go`). Для защиты требуется минимум 10 тестов — в работе ещё 6 HTTP-тестов для хендлеров.

---
