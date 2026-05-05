# Руководство по ДЕМО — GoBackendFootball

**Версия:** 4.0.0 | **Дата:** 30 апреля 2026 г.

---

## Очистка базы данных (сброс до нуля)

Чтобы начать с «чистой» БД, выполните следующие команды:

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"

# Остановите бэкенд (Ctrl+C), затем:
docker-compose down -v
docker-compose up -d db
timeout 5

# Терминал 1 — ML-модель
cd ml_model
pip install -r requirements.txt   # только при первом запуске
python train.py                   # только при первом запуске
python app.py                     # оставьте работать

# Терминал 2 — Go-бэкенд (из корня проекта)
cd ..
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
# Терминал 1 — База данных
cd "C:\Users\user\Desktop\BackForRPO"
docker-compose down -v
docker-compose up -d db
timeout 5
```

```powershell
# Терминал 2 — ML-модель (оставьте работать)
cd "C:\Users\user\Desktop\BackForRPO\ml_model"
pip install -r requirements.txt   # только при первом запуске
python train.py                   # только при первом запуске
python app.py
```

```powershell
# Терминал 3 — Go-бэкенд (из корня проекта)
cd "C:\Users\user\Desktop\BackForRPO"
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

### Шаг 7. AI Прогнозы (3 мин)

1. Убедитесь, что ML-сервер запущен (`python app.py` в отдельном терминале).
2. Перейдите в "Predictions".
3. Нажмите "Get prediction" для матча.
4. Покажите анимированные прогресс-бары с реальными процентами.

> Прогноз генерируется моделью XGBoost (градиентный бустинг), обученной на ~9400 матчах АПЛ. Go-бэкенд отправляет HTTP-запрос к Python-серверу (FastAPI, порт 8000), получает вероятности и сохраняет в PostgreSQL. Повторный запрос прогноза возвращает кэшированный результат из БД.
>
> Если ML-сервер не запущен, бэкенд автоматически использует запасные значения — приложение не падает.

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

## Структура ML-модели

```
ml_model/
epl_final.csv           # Датасет: ~9400 матчей АПЛ (2000–2024)
dataset.py              # Предобработка: скользящие средние за 5 матчей
model.py                # Обёртка XGBClassifier (градиентный бустинг)
train.py                # Обучение модели и вывод метрик
app.py                  # FastAPI-сервер (порт 8000)
requirements.txt        # Зависимости Python
xgboost_model.pkl       # Обученная модель (создаётся после train.py)
team_encoder.pkl        # Кодировщик названий команд
latest_team_stats.pkl   # Последние 5 матчей каждой команды
feature_cols.json       # Порядок признаков для модели
```

---

**Версия документа:** 4.0.0 | **Дата:** 30 апреля 2026 г.
