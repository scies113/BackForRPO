# GoBackendFootball
Бэкенд-сервис для управления футбольной статистикой на Go с использованием микросервиса машинного обучения.

**Статус проекта:** 🟢 Готовность 100% (Финальная версия)

---

## 📊 Общий прогресс и ключевые фичи
- **Аутентификация и авторизация:** Безопасные JWT-токены передаются через HttpOnly Cookie + Header.
- **Строгая Ролевая модель (RBAC):**
  - `admin` — полный доступ к системе, просмотр журнала аудита, доступ к Swagger UI.
  - `operator` — управление расписанием (CRUD матчей).
  - `analyst` — создание матчей и запуск ML-прогнозов.
  - `user` (болельщик) — только просмотр матчей и готовых прогнозов.
- **Модуль ИИ:** Генерация прогнозов исходов матчей (Python, FastAPI, XGBoost) с кэшированием результатов в PostgreSQL.
- **Аудит (журнал действий):** Логирование критически важных действий (LOGIN, CREATE, UPDATE, DELETE, PREDICT) по ГОСТ.
- **Документация API:** Автогенерируемый Swagger UI (доступен строго для администраторов).
- **Инфраструктура:** Полная контейнеризация через Docker Compose (Go API + FastAPI ML + PostgreSQL 15).

---

## 📁 Структура проекта и новые артефакты

В корне проекта теперь находится папка **`document/`**, в которой аккуратно сложены все сгенерированные диаграммы и листинги для вашего отчёта:

```text
GoBackendFootball/
├── cmd/api/main.go              # Точка входа, маршрутизация, CORS, Swagger UI, RBAC
├── document/                    # 🆕 Артефакты, схемы и диаграммы для отчёта
│   ├── activity_diagram.txt       # Диаграмма регистрации и авторизации
│   ├── deployment_diagram.txt     # Архитектура развёртывания (Docker)
│   ├── er_diagram.txt             # Схема БД (GORM)
│   ├── ml_training_activity.txt   # Процесс обучения ML
│   ├── sequence_diagram.txt       # Процесс генерации прогноза
│   └── ЛИСТИНГИ_КОДА.md           # Актуальные листинги кода
├── frontend/                    # Vanilla JS фронтенд
├── internal/                    # Бизнес-логика (handlers, services, models, middleware)
├── migrations/                  # SQL-миграции базы данных (golang-migrate)
├── ml_model/                    # Python-сервис машинного обучения
├── docker-compose.yml           # Конфигурация Docker
└── README.md                    # Этот файл
```

---

## 🔐 Актуальная матрица прав доступа (RBAC)

Все защищенные эндпоинты закрыты с помощью `AuthMiddleware` и `RoleMiddleware`.

| Эндпоинт | admin | operator | analyst | user (fan) |
|----------|-------|----------|---------|------------|
| `POST /api/register` | ✅ | ✅ | ✅ | ✅ |
| `POST /api/login` | ✅ | ✅ | ✅ | ✅ |
| `GET /swagger/*any` | ✅ | ❌ | ❌ | ❌ |
| `GET /api/matches` | ✅ | ✅ | ✅ | ✅ |
| `POST /api/matches` | ✅ | ✅ | ✅ | ❌ |
| `PUT /api/matches/:id` | ✅ | ✅ | ✅ | ❌ |
| `DELETE /api/matches/:id`| ✅ | ✅ | ✅ | ❌ |
| `POST /api/predict/:id`| ✅ | ❌ | ✅ | ❌ |
| `GET /api/audit` | ✅ | ❌ | ❌ | ❌ |
| `POST /api/admin/register`| ✅ | ❌ | ❌ | ❌ |

---

## 🚀 Быстрый старт (Docker Compose)

Проект полностью настроен для запуска в контейнерах. Команда поднимет PostgreSQL, Python ML-сервис и Go-бэкенд в единой виртуальной сети.

```bash
# Сборка и запуск всех сервисов в фоновом режиме
docker-compose up --build -d
```

### Доступные адреса после запуска:
- **Основной сайт (Фронтенд):** [http://localhost:8080](http://localhost:8080)
- **API Бэкенда:** `http://localhost:8080/api/...`
- **Защищенный Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) *(Потребуется авторизация под учётной записью `admin`)*
- **ML-микросервис:** Работает во внутренней сети Docker на порту `8000`.

---

## 🧪 Тестирование
В проекте реализовано 22 автоматических теста (handlers, services, middlewares), которые полностью покрывают логику RBAC и валидации.

Запуск тестов локально:
```bash
go test ./internal/... -v
```
