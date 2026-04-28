# DEMO Guide - GoBackendFootball

**Version:** 3.0.0 | **Date:** 28 April 2026

---

## Clean Database (Reset to Zero)

To start with a fresh DB, run these commands:

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"

# Stop backend (Ctrl+C), then:
docker-compose down -v
docker-compose up -d db
timeout 5
go run cmd/api/main.go
```

What happens:
- `docker-compose down -v` removes PostgreSQL container and ALL data
- `docker-compose up -d db` creates a new clean container
- `go run cmd/api/main.go` automatically:
  1. Creates all tables (AutoMigrate)
  2. Creates roles: admin, operator, analyst, user
  3. Creates default admin: `admin@demo.com` / `admin123`

Console output:
```
[SEED] Default admin created: email=admin@demo.com password=admin123
INFO  Server started              {"port": "8080"}
```

---

## Step 0: Prepare Roles (BEFORE DEMO)

> Do this once before the presentation. Takes 2 minutes.

### 0.1 Start the system (clean DB)

```powershell
cd "C:\Users\bigge\OneDrive\Рабочий стол\GoBackendFootball"
docker-compose down -v
docker-compose up -d db
timeout 5
go run cmd/api/main.go
```

Admin is created automatically. No SQL commands needed!

### 0.2 Get admin token via Swagger

1. Open http://localhost:8080/swagger/index.html
2. Find `POST /api/login` and click **Try it out**
3. Paste:
```json
{
  "email": "admin@demo.com",
  "password": "admin123"
}
```
4. Click **Execute** - copy the `"token"` from response

### 0.3 Authorize in Swagger

1. Click **Authorize** button (lock icon, top right)
2. In BearerAuth field enter: `Bearer <your_token>`
3. Click **Authorize**, then **Close**

### 0.4 Create users with different roles

Call `POST /api/admin/register` for each role:

**Operator:**
```json
{
  "username": "operator1",
  "email": "operator@demo.com",
  "password": "123456",
  "role": "operator"
}
```

**Analyst:**
```json
{
  "username": "analyst1",
  "email": "analyst@demo.com",
  "password": "123456",
  "role": "analyst"
}
```

**User:**
```json
{
  "username": "user1",
  "email": "user@demo.com",
  "password": "123456",
  "role": "user"
}
```

Each call should return **201 Created**.

### Result: 4 users with different roles

| Role | Email | Password | Permissions |
|------|-------|----------|-------------|
| admin | admin@demo.com | admin123 | Everything + audit |
| operator | operator@demo.com | 123456 | CRUD matches + predictions |
| analyst | analyst@demo.com | 123456 | Predictions only |
| user | user@demo.com | 123456 | View only |

---

## Demo Scenario (15 minutes)

### Step 1. Landing Page (1 min)

1. Open http://localhost:8080
2. Show: hero section, feature cards, API Docs button

> Frontend is built with vanilla HTML/CSS/JS and served directly by Go server - no separate Node.js needed.

---

### Step 2. Register a regular user (2 min)

1. Click "Start" -> login page
2. Switch to "Registration" tab
3. Fill: Name `testuser`, Email `test@demo.com`, Password `123456`
4. Click "Register"

> Registration sends POST to /api/register. Password is hashed with bcrypt. Server sets JWT in HttpOnly Cookie.

---

### Step 3. Dashboard (1 min)

1. After registration - auto redirect to Dashboard
2. Show: greeting with name, role "User", navigation cards

> Dashboard loads profile via GET /api/me. Token is sent automatically in Cookie.

---

### Step 4. Matches - view without permissions (2 min)

1. Click "Matches"
2. "Add match" button is hidden - user role has no write permissions!

> Notice the "Add match" button is missing. User role can only view. This is role-based access control in action.

---

### Step 5. Switch to admin (2 min)

1. Click "Logout"
2. Login with admin: `admin@demo.com` / `admin123`
3. Dashboard shows role "Administrator"

---

### Step 6. CRUD Matches (3 min)

1. Go to "Matches" - now "Add match" button is visible
2. Create match: Arsenal vs Chelsea, tomorrow, 0:0, Scheduled
3. Create another: Real Madrid vs Barcelona
4. Edit first match - change score to 2:1
5. Delete second match

> Each operation sends HTTP request: POST, PUT, DELETE. Server validates data. Every action is logged in audit.

---

### Step 7. AI Predictions (2 min)

1. Go to "Predictions"
2. Click "Get prediction" for a match
3. Show animated progress bars

> Prediction is generated on backend and cached in DB. Currently uses a stub - in production this would be an ML model.

---

### Step 8. Audit Log (2 min)

1. Go to "Audit"
2. Show entries: LOGIN, CREATE, UPDATE, DELETE, PREDICT
3. Use filters: by action CREATE, by date, reset

> Audit log records all user actions. Supports filtering and pagination.

---

## Bonus Points

### Swagger UI
Open http://localhost:8080/swagger/index.html - show all endpoints.

### Tests
```bash
go test ./internal/... -v
```
> 22 tests - unit tests for services, middleware, and HTTP handler tests.

### Migrations
```bash
go run cmd/migrate/main.go -direction=up
```
> Migrations via golang-migrate. 5 pairs of up/down files.

### Logs
Requests are logged via zap - show logs during demo.

---

## Frontend Structure

```
frontend/
index.html          # Landing page
login.html          # Login / Registration
dashboard.html      # Dashboard
matches.html        # CRUD matches
predict.html        # AI Predictions
audit.html          # Audit log (admin only)
css/styles.css      # Dark theme, animations
js/api.js           # Fetch wrapper for API
js/auth.js          # Auth + roles
js/app.js           # UI utilities (toast, dates)
```

---

**Document version:** 3.0.0 | **Date:** 28 April 2026
