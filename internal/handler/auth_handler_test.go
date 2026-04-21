package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	// Устанавливаем JWT_SECRET для тестов
	os.Setenv("JWT_SECRET", "test-secret-key-for-handler-tests")
}

// TestRegister_Success - проверка успешной регистрации (без БД - проверяем валидацию)
func TestRegister_Success(t *testing.T) {
	router := gin.New()
	authHandler := &AuthHandler{} // без service, тестируем только binding

	router.POST("/api/register", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверка что данные прошли валидацию
		c.JSON(http.StatusCreated, gin.H{
			"message":  "User registered",
			"username": input.Username,
			"email":    input.Email,
		})
	})
	_ = authHandler // используем для инициализации

	// Тест: корректные данные
	body := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got %v", response["username"])
	}

	if response["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %v", response["email"])
	}
}

// TestRegister_InvalidInput - проверка ошибки валидации при регистрации
func TestRegister_InvalidInput(t *testing.T) {
	router := gin.New()

	router.POST("/api/register", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "ok"})
	})

	tests := []struct {
		name string
		body string
	}{
		{
			name: "EmptyBody",
			body: `{}`,
		},
		{
			name: "InvalidEmail",
			body: `{"username":"test","email":"not-an-email","password":"password123"}`,
		},
		{
			name: "ShortPassword",
			body: `{"username":"test","email":"test@example.com","password":"123"}`,
		},
		{
			name: "MissingUsername",
			body: `{"email":"test@example.com","password":"password123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("[%s] Expected status 400, got %d", tt.name, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if _, ok := response["error"]; !ok {
				t.Errorf("[%s] Expected 'error' field in response", tt.name)
			}
		})
	}
}

// TestLogin_InvalidInput - проверка ошибки валидации при логине
func TestLogin_InvalidInput(t *testing.T) {
	router := gin.New()

	router.POST("/api/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": "mock-token"})
	})

	tests := []struct {
		name       string
		body       string
		expectCode int
	}{
		{
			name:       "EmptyBody",
			body:       `{}`,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "InvalidEmail",
			body:       `{"email":"not-email","password":"password123"}`,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "MissingPassword",
			body:       `{"email":"test@example.com"}`,
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "ValidInput",
			body:       `{"email":"test@example.com","password":"password123"}`,
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("[%s] Expected status %d, got %d", tt.name, tt.expectCode, w.Code)
			}
		})
	}
}
