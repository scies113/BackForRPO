package middleware

import (
	"BackendFootball/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupTestRouter(handlerFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/test", handlerFunc)
	return r
}

func TestGenerateToken(t *testing.T) {
	// Устанавливаем секретный ключ для тестов
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Создаём пользователя с явным ID (как будто он из БД)
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		RoleID:   1,
	}
	// Устанавливаем ID вручную (в реальном сценарии он присваивается при сохранении в БД)
	user.ID = 1

	token, err := GenerateToken(user, "fan")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Тест 2: Токен не должен быть пустым
	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Тест 3: Парсинг и проверка токена
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Тест 4: Токен должен быть валидным
	if !parsedToken.Valid {
		t.Error("Token should be valid")
	}

	// Тест 5: Проверка claims
	if claims.UserID != 1 {
		t.Errorf("Expected user_id 1, got %d", claims.UserID)
	}

	if claims.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", claims.Username)
	}

	if claims.Role != "fan" {
		t.Errorf("Expected role 'fan', got %s", claims.Role)
	}

	// Тест 6: Проверка срока действия
	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("Token should not be expired yet")
	}
}

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Тест 7: Проверка middleware с правильной ролью
	t.Run("ValidRole", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		// Устанавливаем роль в контекст
		c.Set("userRole", "admin")

		// Создаём финальный хендлер
		finalHandler := func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		}

		// Вызываем middleware: RoleMiddleware возвращает gin.HandlerFunc
		handler := RoleMiddleware("admin", "operator")
		handler(c)
		
		// Если middleware пропустило, вызываем финальный хендлер
		if w.Code == http.StatusOK {
			finalHandler(c)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	// Тест 8: Проверка middleware с неправильной ролью
	t.Run("InvalidRole", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		// Устанавливаем роль в контекст
		c.Set("userRole", "fan")

		// Вызываем middleware
		handler := RoleMiddleware("admin", "operator")
		handler(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	// Тест 9: Проверка middleware без роли
	t.Run("NoRole", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		// Вызываем middleware (роль не установлена)
		handler := RoleMiddleware("admin")
		handler(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}
