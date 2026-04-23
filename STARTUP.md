# 🚀 Как запустить GoBackendFootball с нуля

**Для тех, кто вообще в этом не разбирается** — пошагово, с картинками команд.

---

## 📋 Что должно быть установлено

Перед началом убедись, что у тебя есть:

| Программа | Проверить командой | Скачать |
|-----------|-------------------|---------|
| **Go** | `go version` | https://go.dev/dl/ |
| **Docker Desktop** | `docker --version` | https://www.docker.com/products/docker-desktop/ |
| **Git** | `git --version` | https://git-scm.com/downloads |

---

## 🧹 Сценарий 1: Полная очистка и запуск с нуля (чистая БД)

Используй этот сценарий **перед защитой** — БД будет пустая, красиво покажешь создание данных.

### Шаг 1. Открой терминал

Открой **PowerShell** или **Терминал VS Code** и перейди в папку проекта:

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
```

### Шаг 2. Запусти Docker Desktop

1. Найди **Docker Desktop** в меню Пуск
2. Запусти его
3. **Подожди 1-2 минуты** пока в трее появится иконка кита 🐳
4. Проверь что Docker работает:

```powershell
docker --version
```

Должно вывести что-то типа: `Docker version 27.x.x`

### Шаг 3. Удали старую базу данных (полная очистка)

```powershell
docker-compose down -v
```

**Что делает:** останавливает контейнер PostgreSQL и **удаляет все данные** (флаг `-v` удаляет volumes — это хранилище БД).

### Шаг 4. Запусти чистую базу данных

```powershell
docker-compose up -d db
```

**Что делает:** скачивает (если нужно) и запускает PostgreSQL 15 в Docker-контейнере.

Подожди 5 секунд и проверь:

```powershell
docker ps
```

Ты должен увидеть контейнер со статусом `Up ... (healthy)`:

```
CONTAINER ID   IMAGE         STATUS                    PORTS
xxxxxxxxxxxx   postgres:15   Up 10 seconds (healthy)   0.0.0.0:5433->5432/tcp
```

### Шаг 5. Запусти бэкенд

```powershell
go run cmd/api/main.go
```

**Что делает:** компилирует и запускает Go-сервер. Он автоматически:
- Подключится к PostgreSQL
- Создаст все таблицы (AutoMigrate)
- Создаст роли: admin, operator, analyst, user
- Запустит HTTP-сервер на порту 8080

Ты увидишь:

```
INFO  Сервер запущен       {"port": "8080"}
INFO  Фронтенд доступен   {"url": "http://localhost:8080"}
INFO  Swagger UI           {"url": "http://localhost:8080/swagger/index.html"}
```

### Шаг 6. Открой в браузере

- **Фронтенд:** http://localhost:8080
- **Swagger:** http://localhost:8080/swagger/index.html

**Готово! 🎉** БД чистая, можно демонстрировать.

---

## ▶️ Сценарий 2: Обычный запуск (данные сохраняются)

Используй когда просто хочешь продолжить работу — данные в БД останутся.

```powershell
# 1. Перейди в папку проекта
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"

# 2. Запусти PostgreSQL (если не запущен)
docker-compose up -d db

# 3. Запусти бэкенд
go run cmd/api/main.go
```

---

## 🛑 Как остановить

### Остановить бэкенд
Нажми **Ctrl+C** в терминале где работает `go run`.

### Остановить PostgreSQL
```powershell
docker-compose down
```

### Остановить ВСЁ и удалить данные
```powershell
docker-compose down -v
```

---

## 🔍 Как проверить что всё работает

### Проверка 1: Бэкенд отвечает
Открой в браузере: http://localhost:8080  
Должен увидеть красивый лендинг FootballStats.

### Проверка 2: API работает
Открой: http://localhost:8080/swagger/index.html  
Должен увидеть Swagger UI с документацией.

### Проверка 3: Тесты проходят
```powershell
go test ./internal/... -v
```
Должно быть: `22 тестов, все PASS`.

### Проверка 4: PostgreSQL работает
```powershell
docker ps
```
Должен видеть контейнер `postgres:15` со статусом `healthy`.

---

## 🆘 Частые проблемы

### ❌ «Ошибка подключения к базе данных»

**Причина:** PostgreSQL не запущен.

**Решение:**
```powershell
docker-compose up -d db
# Подожди 5 секунд
go run cmd/api/main.go
```

### ❌ «failed to connect to docker API»

**Причина:** Docker Desktop не запущен.

**Решение:** Открой Docker Desktop из меню Пуск, подожди 1-2 мин.

### ❌ «port 5433 already in use»

**Причина:** Старый контейнер PostgreSQL висит.

**Решение:**
```powershell
docker-compose down
docker-compose up -d db
```

### ❌ «port 8080 already in use»

**Причина:** Старый бэкенд ещё работает.

**Решение:** Закрой предыдущий терминал или нажми Ctrl+C.

### ❌ Запускаю из cmd/api/ — переменные пустые

**Причина:** `.env` файл ищется в текущей папке.

**Решение:** Запускай **ТОЛЬКО из корня проекта:**
```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
go run cmd/api/main.go
```

---

## 📝 Шпаргалка (копируй и вставляй)

### Чистый запуск перед защитой:
```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
docker-compose down -v
docker-compose up -d db
timeout 5
go run cmd/api/main.go
```

### Быстрый запуск:
```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
docker-compose up -d db
go run cmd/api/main.go
```

### Полная остановка:
```powershell
docker-compose down -v
```

---

**Версия документа:** 3.0.0  
**Дата:** 21 апреля 2026 г.
