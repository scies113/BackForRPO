# Руководство по ДЕМО — GoBackendFootball

**Версия:** 3.0.0 | **Дата:** 28 апреля 2026 г.

---

## Очистка базы данных (сброс до нуля)

Чтобы начать с «чистой» БД, выполните следующие команды:

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"

# Остановите бэкенд (Ctrl+C), затем:
docker-compose down -v
docker-compose up
timeout 5
go run cmd/api/main.go
```

Что произойдет:
- `docker-compose down -v` удалит контейнер PostgreSQL и ВСЕ данные.
- `docker-compose up -d db` создаст новый чистый контейнер.
- `go run cmd/api/main.go` автоматически:
  1. Создаст все таблицы (AutoMigrate).
  2. Создаст роли: admin, operator, analyst, user.
  3. Создаст администратора по умолчанию: `admin@demo.com` / `admin123`.

Вывод в консоли:
```
[SEED] Default admin created: email=admin@demo.com password=admin123
INFO  Server started              {"port": "8080"}
```

---

## Шаг 0: Подготовка ролей (ПЕРЕД ДЕМО)

> Сделайте это один раз перед презентацией. Занимает 2 минуты.

### 0.1 Запуск системы (чистая БД)

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
docker-compose down -v
docker-compose up -d db
timeout 5
go run cmd/api/main.go
```

Админ создается автоматически. SQL-команды не нужны!

### 0.2 Получение токена админа через Swagger

1. Откройте http://localhost:8080/swagger/index.html
2. Найдите `POST /api/login` и нажмите **Try it out**.
3. Вставьте:
```json
{
  "email": "admin@demo.com",
  "password": "admin123"
}
```
4. Нажмите **Execute** — скопируйте `"token"` из ответа.

### 0.3 Авторизация в Swagger

1. Нажмите кнопку **Authorize** (иконка замка справа вверху).
2. В поле BearerAuth введите: `Bearer <ваш_токен>`.
3. Нажмите **Authorize**, затем **Close**.

### 0.4 Создание пользователей с разными ролями

Вызовите `POST /api/admin/register` для каждой роли:

**Оператор (Operator):**
```json
{
  "username": "operator1",
  "email": "operator@demo.com",
  "password": "123456",
  "role": "operator"
}
```

**Аналитик (Analyst):**
```json
{
  "username": "analyst1",
  "email": "analyst@demo.com",
  "password": "123456",
  "role": "analyst"
}
```

**Пользователь (User):**
```json
{
  "username": "user1",
  "email": "user@demo.com",
  "password": "123456",
  "role": "user"
}
```

Каждый вызов должен возвращать **201 Created**.

### Итог: 4 пользователя с разными ролями

| Роль | Email | Пароль | Права доступа |
|------|-------|--------|---------------|
| admin | admin@demo.com | admin123 | Всё + аудит |
| operator | operator@demo.com | 123456 | CRUD матчей + прогнозы |
| analyst | analyst@demo.com | 123456 | Только прогнозы |
| user | user@demo.com | 123456 | Только просмотр |

---

## Сценарий демо (15 минут)

### Шаг 1. Главная страница (1 мин)

1. Откройте http://localhost:8080
2. Покажите: главный блок (hero section), карточки функций, кнопку API Docs.

> Фронтенд построен на чистом HTML/CSS/JS и раздается напрямую Go-сервером — отдельный Node.js не требуется.

---

### Шаг 2. Регистрация обычного пользователя (2 мин)

1. Нажмите "Start" -> страница входа.
2. Переключитесь на вкладку "Registration".
3. Заполните: Имя `testuser`, Email `test@demo.com`, Пароль `123456`.
4. Нажмите "Register".

> Регистрация отправляет POST на /api/register. Пароль хешируется с помощью bcrypt. Сервер устанавливает JWT в HttpOnly Cookie.

---

### Шаг 3. Личный кабинет (Dashboard) (1 мин)

1. После регистрации — автопереход в Dashboard.
2. Покажите: приветствие с именем, роль "User", навигационные карточки.

> Dashboard загружает профиль через GET /api/me. Токен отправляется автоматически в Cookie.

---

### Шаг 4. Матчи — просмотр без прав (2 мин)

1. Нажмите "Matches".
2. Кнопка "Add match" скрыта — у роли user нет прав на запись!

> Обратите внимание, что кнопка "Add match" отсутствует. Роль User может только просматривать. Это пример работы RBAC (управления доступом на основе ролей).

---

### Шаг 5. Переключение на админа (2 мин)

1. Нажмите "Logout".
2. Войдите как админ: `admin@demo.com` / `admin123`.
3. Dashboard показывает роль "Administrator".

---

### Шаг 6. CRUD матчей (3 мин)

1. Перейдите в "Matches" — теперь кнопка "Add match" видна.
2. Создайте матч: Arsenal vs Chelsea, завтра, 0:0, Scheduled.
3. Создайте еще один: Real Madrid vs Barcelona.
4. Отредактируйте первый матч — измените счет на 2:1.
5. Удалите второй матч.

> Каждая операция отправляет HTTP-запрос: POST, PUT, DELETE. Сервер валидирует данные. Каждое действие записывается в аудит.

---

### Шаг 7. AI Прогнозы (2 мин)

1. Перейдите в "Predictions".
2. Нажмите "Get prediction" для матча.
3. Покажите анимированные прогресс-бары.

> Прогноз генерируется на бэкенде и кэшируется в БД. Сейчас используется заглушка — в продакшене здесь была бы ML-модель.

---

### Шаг 8. Журнал аудита (2 мин)

1. Перейдите в "Audit".
2. Покажите записи: LOGIN, CREATE, UPDATE, DELETE, PREDICT.
3. Используйте фильтры: по действию CREATE, по дате, сброс.

> Журнал аудита фиксирует все действия пользователей. Поддерживает фильтрацию и пагинацию.

---

## Дополнительные моменты

### Swagger UI
Откройте http://localhost:8080/swagger/index.html — покажите все эндпоинты.

### Тесты
```bash
go test ./internal/... -v
```
> 22 теста — юнит-тесты для сервисов, middleware и тесты HTTP-хендлеров.

### Миграции
```bash
go run cmd/migrate/main.go -direction=up
```
> Миграции через golang-migrate. 5 пар файлов up/down.

### Логи
Запросы логируются через zap — покажите логи во время демо.

---

## Структура фронтенда

```
frontend/
index.html          # Главная страница
login.html          # Вход / Регистрация
dashboard.html      # Личный кабинет
matches.html        # CRUD матчей
predict.html        # AI Прогнозы
audit.html          # Журнал аудита (только админ)
css/styles.css      # Темная тема, анимации
js/api.js           # Обертка Fetch для API
js/auth.js          # Авторизация и роли
js/app.js           # Утилиты UI (уведомления, даты)
```

---

**Версия документа:** 3.0.0 | **Дата:** 28 апреля 2026 г.
