# 🏗 Архитектура GoBackendFootball

**Версия:** 3.0.0  
**Дата:** 20 апреля 2026 г.

Этот документ описывает **архитектуру** и **принцип работы** бэкенд-сервиса для управления футбольной статистикой.

---

## 📋 Оглавление

1. [Общая схема](#общая-схема)
2. [Уровни архитектуры](#уровни-архитектуры)
3. [Потоки данных](#потоки-данных)
4. [База данных](#база-данных)
5. [Аутентификация и авторизация](#аутентификация-и-авторизация)
6. [JWT-токены](#jwt-токены)
7. [Безопасность](#безопасность)
8. [Масштабируемость](#масштабируемость)

---

## 📊 Общая схема

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              КЛИЕНТЫ                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │  React   │  │  Mobile  │  │  cURL    │  │  Postman │               │
│  │  (Web)   │  │   App    │  │  (CLI)   │  │  (Test)  │               │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘               │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ HTTP/HTTPS
                                    │ (JSON + Cookie)
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         HTTP-СЕРВЕР (Gin)                               │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                     Маршрутизатор (Router)                        │  │
│  │  /api/register  /api/login  /api/matches  /api/predict  /api/audit│  │
│  └───────────────────────────────────────────────────────────────────┘  │
│         │                    │                    │                      │
│         ▼                    ▼                    ▼                      │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐               │
│  │ Middleware  │     │  Handler    │     │   Service   │               │
│  │  - JWT      │────▶│  (HTTP)     │────▶│  (Business) │               │
│  │  - Roles    │     │             │     │             │               │
│  └─────────────┘     └─────────────┘     └─────────────┘               │
│                                               │                         │
│                                               ▼                         │
│                                        ┌─────────────┐                  │
│                                        │ Repository  │                  │
│                                        │   (GORM)    │                  │
│                                        └─────────────┘                  │
│                                               │                         │
└───────────────────────────────────────────────┼─────────────────────────┘
                                                │
                                                │ SQL
                                                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        PostgreSQL (БД)                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │  roles   │  │  users   │  │ matches  │  │ audit_   │               │
│  │          │  │          │  │          │  │  logs    │               │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘               │
│  ┌──────────┐  ┌──────────┐                                            │
│  │predictions│ │  teams   │  (опционально)                            │
│  └──────────┘  └──────────┘                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 🏛 Уровни архитектуры

### Уровень 1: Клиенты (Clients)

**Назначение:** Взаимодействие с пользователем, отправка запросов на сервер.

**Типы клиентов:**
- **React (Web)** — будущий фронтенд (localhost:3000)
- **Mobile App** — мобильное приложение (iOS/Android)
- **cURL/Postman** — инструменты для тестирования API

**Что отправляют:**
- HTTP-запросы (GET, POST, PUT, DELETE)
- Заголовки (Authorization, Content-Type)
- Куки (JWT-токен)
- JSON-данные

**Что получают:**
- JSON-ответы
- HTTP-статусы (200, 201, 400, 401, 403, 404, 500)
- Куки (Set-Cookie)

---

### Уровень 2: HTTP-сервер (Gin Framework)

**Файл:** `cmd/api/main.go`

**Назначение:** Приём HTTP-запросов, маршрутизация, middleware.

**Компоненты:**

#### 2.1. Маршрутизатор (Router)

**Группы маршрутов:**
```go
api := r.Group("/api")
{
    // Публичные (без авторизации)
    api.POST("/register", authHandler.Register)
    api.POST("/login", authHandler.Login)

    // Защищённые (требуется токен)
    secure := api.Use(middleware.AuthMiddleware())
    {
        secure.GET("/matches", matchHandler.GetAllMatches)
        secure.GET("/matches/:id", matchHandler.GetMatchByID)
        secure.POST("/matches", matchHandler.CreateMatch, 
                    middleware.RoleMiddleware("admin", "operator"))
        // ...
    }
}
```

#### 2.2. CORS Middleware

**Назначение:** Разрешить запросы с фронтенда (React).

**Настройки:**
```go
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
```

---

### Уровень 3: Middleware (Промежуточное ПО)

**Папка:** `internal/middleware/`

**Назначение:** Обработка запроса ДО передачи в хендлер.

#### 3.1. AuthMiddleware — проверка JWT

**Файл:** `internal/middleware/auth.go`

**Алгоритм:**
```
1. Получить токен из заголовка Authorization: Bearer <token>
2. Если нет в заголовке → попробовать куку "token"
3. Если нет нигде → вернуть 401 Unauthorized
4. Расшифровать токен с помощью JWT_SECRET
5. Проверить подпись и срок действия
6. Извлечь claims (user_id, username, role)
7. Сохранить в контекст Gin: c.Set("userID", ...), c.Set("userRole", ...)
8. Передать управление следующему хендлеру
```

**Код:**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var tokenString string

        // Пробуем заголовок
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

        // Сохраняем в контекст
        c.Set("userID", claims.UserID)
        c.Set("userName", claims.Username)
        c.Set("userRole", claims.Role)

        c.Next()
    }
}
```

#### 3.2. RoleMiddleware — проверка роли

**Алгоритм:**
```
1. Получить роль из контекста (c.Get("userRole"))
2. Сравнить с требуемыми ролями
3. Если совпадает → c.Next()
4. Если нет → вернуть 403 Forbidden
```

**Код:**
```go
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("userRole")

        for _, role := range requiredRoles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
    }
}
```

---

### Уровень 4: Handlers (Обработчики HTTP)

**Папка:** `internal/handler/`

**Назначение:** Приём запроса, валидация, вызов сервиса, возврат ответа.

**Структура хендлера:**
```go
type AuthHandler struct {
    service *service.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{service: service.NewAuthService()}
}
```

#### 4.1. Пример: Register Handler

**Алгоритм:**
```
1. Распарсить JSON из тела запроса
2. Валидировать данные (binding:"required,email,min=6")
3. Создать модель User
4. Вызвать сервис: service.Register(user, password)
5. Установить куку с токеном: c.SetCookie(...)
6. Вернуть JSON-ответ
```

**Код:**
```go
func (h *AuthHandler) Register(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := &model.User{
        Username: input.Username,
        Email:    input.Email,
    }

    token, err := h.service.Register(user, input.Password)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    h.setAuthCookie(c, token)

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered",
        "token":   token,
        "role":    "user",
    })
}
```

**Все хендлеры:**
| Файл | Методы | Эндпоинты |
|------|--------|-----------|
| `auth_handler.go` | Register, RegisterAdmin, Login, Logout, GetMe | POST /register, /login, /logout, /admin/register, GET /me |
| `match_handler.go` | GetAllMatches, GetMatchByID, CreateMatch, UpdateMatch, DeleteMatch | GET/POST/PUT/DELETE /matches |
| `predict_handler.go` | GetPrediction | POST /predict/:id |
| `audit_handler.go` | GetAuditLogs | GET /audit (фильтрация + пагинация) |

---

### Уровень 5: Services (Бизнес-логика)

**Папка:** `internal/service/`

**Назначение:** Реализация бизнес-правил, валидация, аудит.

#### 5.1. Пример: MatchService.CreateMatch

**Алгоритм:**
```
1. Валидация: команды не пустые
2. Валидация: home_team != away_team
3. Валидация: дата не в прошлом
4. Вызов репозитория: repo.Create(match)
5. Запись в аудит: logAudit(userID, "CREATE", "Match", matchID)
6. Вернуть результат
```

**Код:**
```go
func (s *MatchService) CreateMatch(match *model.Match, userID uint, userName string) error {
    // 1. Валидация
    if match.HomeTeam == "" || match.AwayTeam == "" {
        return fmt.Errorf("home_team and away_team are required")
    }

    if match.HomeTeam == match.AwayTeam {
        return errors.ErrSameTeams
    }

    if match.MatchDate.Before(time.Now()) {
        return errors.ErrInvalidDate
    }

    // 2. Сохранение
    if err := s.repo.Create(match); err != nil {
        return err
    }

    // 3. Аудит
    s.logAudit(userID, userName, "CREATE", "Match", match.ID)

    return nil
}
```

**Все сервисы:**
| Файл | Методы | Описание |
|------|--------|----------|
| `auth_service.go` | Register, RegisterWithRole, Login, logAudit | Аутентификация, хеширование, токен, LOGIN аудит |
| `match_service.go` | CreateMatch, GetMatchByID, UpdateMatch, DeleteMatch, GetAllMatches | CRUD матчей, валидация, аудит |
| `predict_service.go` | GetPrediction | Генерация прогноза ИИ, кэширование в БД |
| `audit_service.go` | GetAllLogs, GetLogsFiltered | Журнал аудита + фильтрация + пагинация |

---

### Уровень 6: Repositories (Работа с БД)

**Папка:** `internal/repository/`

**Назначение:** Инкапсуляция SQL-запросов (GORM).

#### 6.1. Пример: MatchRepository

**Код:**
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

func (r *MatchRepository) Update(match *model.Match) error {
    return r.db.Save(match).Error
}

func (r *MatchRepository) Delete(id uint) error {
    return r.db.Delete(&model.Match{}, id).Error
}
```

**Паттерн Repository:**
- Изолирует SQL от бизнес-логики
- Упрощает тестирование (можно замокать)
- Единое место для изменения запросов

---

### Уровень 7: Models (Модели данных)

**Папка:** `internal/model/`

**Назначение:** Описание структуры таблиц и бизнес-объектов.

#### 7.1. User (Пользователь)

```go
type User struct {
    gorm.Model       // ID, CreatedAt, UpdatedAt, DeletedAt
    Username     string `gorm:"unique;not null"`
    Email        string `gorm:"unique;not null"`
    PasswordHash string `gorm:"not null"`
    RoleID       uint
    Role         Role       `gorm:"foreignKey:RoleID"`
    AuditLogs    []AuditLog `gorm:"foreignKey:UserID"`
}
```

#### 7.2. Match (Матч)

```go
type Match struct {
    gorm.Model
    HomeTeam   string      `gorm:"not null" json:"home_team"`
    AwayTeam   string      `gorm:"not null" json:"away_team"`
    MatchDate  time.Time   `gorm:"not null" json:"match_date"`
    HomeScore  int         `gorm:"default:0" json:"home_score"`
    AwayScore  int         `gorm:"default:0" json:"away_score"`
    Status     string      `gorm:"default:scheduled" json:"status"`
    Prediction *Prediction `json:"prediction" gorm:"foreignKey:MatchID"`
}
```

#### 7.3. Prediction (Прогноз ИИ)

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

#### 7.4. AuditLog (Журнал аудита)

```go
type AuditLog struct {
    gorm.Model
    UserID    uint
    UserName  string
    Action    string    // CREATE, UPDATE, DELETE, PREDICT, LOGIN
    Entity    string    // Match, User, Prediction
    EntityID  uint
    Timestamp time.Time `gorm:"autoCreateTime"`
}
```

#### 7.5. Role (Роль)

```go
type Role struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"unique;not null"` // admin, operator, analyst, user
}
```

---

### Уровень 8: Database (PostgreSQL)

**Файл:** `internal/database/postgres.go`

**Назначение:** Подключение к БД, создание таблиц, роли.

#### 8.1. Подключение

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
```

#### 8.2. Создание ролей

```go
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

## 🔄 Потоки данных

### Поток 1: Регистрация пользователя

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Клиент   │     │ Handler  │     │ Service  │     │   Repo   │     │    БД    │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │                │
     │ POST /register │                │                │                │
     │ {username,     │                │                │                │
     │  email,        │                │                │                │
     │  password}     │                │                │                │
     │───────────────▶│                │                │                │
     │                │                │                │                │
     │                │ Валидация      │                │                │
     │                │ (email, min=6) │                │                │
     │                │                │                │                │
     │                │ Register()     │                │                │
     │                │───────────────▶│                │                │
     │                │                │                │                │
     │                │                │ Хеширование    │                │
     │                │                │ (bcrypt)       │                │
     │                │                │                │                │
     │                │                │ Create(user)   │                │
     │                │                │───────────────▶│                │
     │                │                │                │                │
     │                │                │                │ INSERT INTO    │
     │                │                │                │ users ...      │
     │                │                │                │───────────────▶│
     │                │                │                │                │
     │                │                │                │ user_id = 5    │
     │                │                │                │◀───────────────│
     │                │                │                │                │
     │                │                │ GenerateToken()│                │
     │                │                │ (JWT_SECRET)   │                │
     │                │                │                │                │
     │                │ SetCookie()    │                │                │
     │                │◀───────────────│                │                │
     │                │                │                │                │
     │ 201 Created    │                │                │                │
     │ {message,      │                │                │                │
     │  token, role}  │                │                │                │
     │◀───────────────│                │                │                │
     │                │                │                │                │
```

---

### Поток 2: Создание матча (с авторизацией)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Клиент   │     │   CORS   │     │   Auth   │     │  Handler │     │  Service │     │    БД    │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │                │                │                │
     │ POST /matches  │                │                │                │                │
     │ Authorization: │                │                │                │                │
     │ Bearer <token> │                │                │                │                │
     │ {home_team,    │                │                │                │                │
     │  away_team,    │                │                │                │                │
     │  match_date}   │                │                │                │                │
     │───────────────▶│                │                │                │                │
     │                │                │                │                │                │
     │                │ CORS Check     │                │                │                │
     │                │ (Origin,       │                │                │                │
     │                │  Credentials)  │                │                │                │
     │                │───────────────▶│                │                │                │
     │                │                │                │                │                │
     │                │                │ Parse Token    │                │                │
     │                │                │ (JWT_SECRET)   │                │                │
     │                │                │ Check Role     │                │                │
     │                │                │ (admin/operator)              │                │
     │                │                │───────────────▶│                │                │
     │                │                │                │                │                │
     │                │                │                │ CreateMatch()  │                │
     │                │                │                │───────────────▶│                │
     │                │                │                │                │                │
     │                │                │                │                │ Валидация      │
     │                │                │                │                │ (teams, date)  │
     │                │                │                │                │                │
     │                │                │                │                │ Create(match)  │
     │                │                │                │                │───────────────▶│
     │                │                │                │                │                │
     │                │                │                │                │                │ INSERT INTO
     │                │                │                │                │                │ matches ...
     │                │                │                │                │                │───────────▶│
     │                │                │                │                │                │
     │                │                │                │                │                │ match_id = 10
     │                │                │                │                │                │◀───────────│
     │                │                │                │                │                │
     │                │                │                │                │ logAudit()     │
     │                │                │                │                │ (CREATE, Match)│
     │                │                │                │                │───────────────▶│
     │                │                │                │                │                │
     │                │                │                │ 200 OK         │                │
     │                │                │                │ {message}      │                │
     │                │                │                │◀───────────────│                │
     │                │                │                │                │                │
     │ 200 OK         │                │                │                │                │
     │ {message}      │                │                │                │                │
     │◀───────────────│                │                │                │                │
     │                │                │                │                │                │
```

---

## 🗄 База данных

### Схема БД

```sql
-- Роли пользователей
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL  -- admin, operator, analyst, user
);

-- Пользователи
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT REFERENCES roles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Матчи
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team VARCHAR(100) NOT NULL,
    away_team VARCHAR(100) NOT NULL,
    match_date TIMESTAMP NOT NULL,
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'scheduled',  -- scheduled, finished
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Прогнозы ИИ
CREATE TABLE predictions (
    id SERIAL PRIMARY KEY,
    match_id INT UNIQUE REFERENCES matches(id),
    home_win_prob FLOAT NOT NULL,
    draw_prob FLOAT NOT NULL,
    away_win_prob FLOAT NOT NULL,
    is_accurate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Журнал аудита
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    user_name VARCHAR(100),
    action VARCHAR(50) NOT NULL,  -- CREATE, UPDATE, DELETE, PREDICT, LOGIN
    entity VARCHAR(50),           -- Match, User, Prediction
    entity_id INT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Связи между таблицами

```
roles (1) ──────< (N) users
                      │
                      │ (1)
                      │
                      ▼
                    users (1) ──────< (N) audit_logs

matches (1) ──────< (1) predictions

matches (1) ──────< (N) audit_logs (через entity_id)
```

---

## 🔐 Аутентификация и авторизация

### Аутентификация (Authentication)

**Вопрос:** "Кто ты?"

**Реализация:**
1. Пользователь отправляет email/password
2. Сервер находит пользователя в БД
3. Проверяет хеш пароля (bcrypt.CompareHashAndPassword)
4. Генерирует JWT-токен
5. Возвращает токен + устанавливает куку

### Авторизация (Authorization)

**Вопрос:** "Что тебе разрешено?"

**Реализация:**
1. Middleware извлекает роль из токена
2. RoleMiddleware сравнивает с требуемыми ролями
3. Если роль подходит → пропускает запрос
4. Если нет → возвращает 403 Forbidden

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

## 🎫 JWT-токены

### Структура токена

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Iml2YW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3NzUwNDQwNjUsImlhdCI6MTc3NDk1NzY2NX0.qmIzhlEZ3RWL6_1Q-F2D8JfsLE0ob6KLjTG24Eix5VA
│                                         │                                                              │
│              Header (Base64)            │                   Payload (Base64)                           │          Signature (Base64)
│                                         │                                                              │
│  {                                      │  {                                                           │  HMAC-SHA256(
│    "alg": "HS256",                      │    "user_id": 1,                                             │    base64UrlEncode(header) + "." +
│    "typ": "JWT"                         │    "username": "ivan",                                       │    base64UrlEncode(payload),
│  }                                      │    "role": "admin",                                          │    JWT_SECRET
│                                         │    "exp": 1775044065,                                        │  )
│                                         │    "iat": 1774957665                                         │
│                                         │  }                                                           │
```

### Claims (данные в токене)

| Claim | Описание | Пример |
|-------|----------|--------|
| `user_id` | ID пользователя в БД | 1 |
| `username` | Имя пользователя | "ivan" |
| `role` | Роль пользователя | "admin" |
| `exp` | Срок действия (Unix timestamp) | 1775044065 |
| `iat` | Время выдачи (Unix timestamp) | 1774957665 |

### Алгоритм подписи

```
Signature = HMAC-SHA256(
    base64UrlEncode(header) + "." + base64UrlEncode(payload),
    JWT_SECRET
)
```

**Почему это безопасно:**
- Если изменить Payload → подпись не совпадёт
- Если нет JWT_SECRET → нельзя создать валидную подпись
- Сервер проверяет подпись при каждом запросе

---

## 🛡 Безопасность

### 1. Хеширование паролей (bcrypt)

**Стоимость:** DefaultCost (10)

```go
func (u *User) SetPassword(password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hash)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}
```

**Защита:**
- ✅ От утечки паролей (БД украли → пароли зашифрованы)
- ✅ От перебора (bcrypt медленный, ~100ms на хеш)

---

### 2. JWT-токены (HttpOnly Cookie)

**Настройки куки:**
```go
c.SetCookie(
    "token",     // имя
    token,       // значение
    86400,       // срок (24 часа)
    "/api",      // путь
    "",          // домен
    false,       // Secure (false для localhost)
    true,        // HttpOnly (JavaScript не имеет доступа)
)
```

**Защита:**
- ✅ HttpOnly — защита от XSS (JavaScript не читает куку)
- ✅ Path=/api — кука отправляется только на /api
- ✅ Secure=true (в продакшене) — только по HTTPS

---

### 3. Валидация входных данных

**На уровне хендлера:**
```go
var input struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

**На уровне сервиса:**
```go
if match.HomeTeam == match.AwayTeam {
    return errors.ErrSameTeams
}

if match.MatchDate.Before(time.Now()) {
    return errors.ErrInvalidDate
}
```

**Защита:**
- ✅ От SQL-инъекций (GORM параметризирует запросы)
- ✅ От некорректных данных (валидация на сервере)

---

### 4. Аудит (журнал действий)

**Записывается:**
- CREATE/UPDATE/DELETE матча
- PREDICT (запрос прогноза)
- LOGIN (вход пользователя)

**Таблица audit_logs:**
```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    user_name VARCHAR(100),
    action VARCHAR(50),  -- CREATE, UPDATE, DELETE, PREDICT, LOGIN
    entity VARCHAR(50),  -- Match, User, Prediction
    entity_id INT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Защита:**
- ✅ Отслеживание всех действий
- ✅ Возможность расследования инцидентов

---

## 📈 Масштабируемость

### Текущая архитектура

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Клиент   │────▶│  Server  │────▶│   PostgreSQL  │
│          │     │  (Go)    │     │   (БД)        │
└──────────┘     └──────────┘     └──────────┘
```

### Горизонтальное масштабирование

```
┌──────────┐
│ Клиент   │
└────┬─────┘
     │
     ▼
┌──────────────────┐
│  Load Balancer   │  (Nginx, HAProxy)
└────────┬─────────┘
         │
    ┌────┴────┬────────────┐
    ▼         ▼            ▼
┌──────┐  ┌──────┐    ┌──────┐
│ Srv1 │  │ Srv2 │    │ SrvN │  (Go, stateless)
└──┬───┘  └──┬───┘    └──┬───┘
   │         │           │
   └─────────┴───────────┘
             │
             ▼
      ┌──────────────┐
      │  PostgreSQL  │  (Master-Slave репликация)
      └──────────────┘
```

**Stateless-серверы:**
- JWT-токены не хранятся на сервере
- Любой сервер может обработать любой запрос
- Легко добавить новые серверы

---

## 📊 Метрики производительности

### Время отклика (в среднем)

| Эндпоинт | Время (ms) |
|----------|------------|
| POST /api/register | ~150ms (bcrypt хеширование) |
| POST /api/login | ~150ms (bcrypt проверка) |
| GET /api/matches | ~10ms |
| GET /api/matches/:id | ~5ms |
| POST /api/matches | ~20ms |
| POST /api/predict/:id | ~15ms |
| GET /api/audit | ~10ms |

### Нагрузка (потенциальная)

| Метрика | Значение |
|---------|----------|
| RPS (Requests Per Second) | ~1000 RPS (один сервер) |
| Concurrent connections | ~10,000 |
| Database connections | ~100 (connection pool) |

---

## 📁 Структура проекта

```
GoBackendFootball/
│
├── cmd/
│   ├── api/
│   │   └── main.go              # Точка входа, маршрутизация, CORS, Swagger UI
│   └── migrate/
│       └── main.go              # CLI миграций (up/down/force)
│
├── internal/
│   ├── handler/                 # Обработчики HTTP + Swagger-аннотации
│   │   ├── auth_handler.go      # Register, Login, Logout, GetMe
│   │   ├── auth_handler_test.go # Тесты регистрации и логина
│   │   ├── match_handler.go     # CRUD матчей
│   │   ├── match_handler_test.go# Тесты авторизации и ролей
│   │   ├── predict_handler.go
│   │   └── audit_handler.go     # Фильтрация + пагинация
│   │
│   ├── service/                 # Бизнес-логика
│   │   ├── auth_service.go      # + LOGIN аудит
│   │   ├── match_service.go
│   │   ├── predict_service.go
│   │   └── audit_service.go     # + фильтрация + пагинация
│   │
│   ├── repository/              # Работа с БД
│   │   └── match_repository.go
│   │
│   ├── model/                   # Модели данных
│   │   ├── user.go
│   │   ├── match.go
│   │   └── audit.go
│   │
│   ├── middleware/              # Промежуточное ПО
│   │   └── auth.go              # JWT + Роли
│   │
│   ├── errors/                  # Типизированные ошибки
│   │   └── errors.go            # + RespondWithError/Success
│   │
│   ├── logger/                  # Логирование
│   │   └── logger.go            # Zap (dev/prod)
│   │
│   └── database/                # Подключение к БД
│       └── postgres.go
│
├── migrations/                  # golang-migrate (5 пар)
│   ├── 000001_create_roles.up/down.sql
│   ├── 000002_create_users.up/down.sql
│   ├── 000003_create_matches.up/down.sql
│   ├── 000004_create_audit_logs.up/down.sql
│   └── 000005_create_predictions.up/down.sql
│
├── ml_model/                    # Модуль машинного обучения (Python FastAPI)
│   ├── app.py                   # FastAPI сервер
│   ├── model.py                 # Обертка XGBoost
│   ├── train.py                 # Скрипт обучения
│   ├── dataset.py               # Обработка признаков
│   ├── epl_final.csv            # Исторический датасет АПЛ
│   ├── Dockerfile               # Docker-образ ML-сервиса
│   ├── requirements.txt         # Зависимости Python
│   └── TEAMS.md                 # Список поддерживаемых команд
│
├── docs/                        # Swagger (автогенерация)
│   ├── docs.go / swagger.json / swagger.yaml
│
├── .env                         # Переменные окружения
├── docker-compose.yml           # Docker-конфигурация (Go + DB + ML)
├── Dockerfile                   # Образ Go-бэкенда
├── README.md                    # Документация
├── ARCHITECTURE.md              # Этот файл
├── GUIDE.md                     # Инструкция по запуску API
├── STARTUP.md                   # Инструкция для новичков
├── DEMO.md                      # Сценарий демонстрации
├── PROGRESS.md                  # Прогресс разработки
├── go.mod                       # Зависимости Go
└── go.sum                       # Версии зависимостей
```

---

## 📚 Дополнительные документы

- **README.md** — общая информация, прогресс, Быстрый старт
- **GUIDE.md** — подробная инструкция по API и тестированию
- **STARTUP.md** — пошаговый запуск для абсолютных новичков
- **DEMO.md** — сценарий презентации (как показывать проект)
- **PROGRESS.md** — текущий прогресс разработки

---

**Версия документа:** 4.0.0  
**Дата:** 5 мая 2026 г.
