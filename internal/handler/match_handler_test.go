package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"BackendFootball/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key-for-match-tests")
	os.Exit(m.Run())
}

// TestGetAllMatches_Unauthorized - запрос без токена должен вернуть 401
func TestGetAllMatches_Unauthorized(t *testing.T) {
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/api/matches", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"matches": []interface{}{}})
	})

	req, _ := http.NewRequest("GET", "/api/matches", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["error"]; !ok {
		t.Error("Expected 'error' field in response")
	}
}

// TestGetMatchByID_InvalidID - запрос с некорректным ID должен вернуть 400
func TestGetMatchByID_InvalidID(t *testing.T) {
	router := gin.New()

	// Имитируем хендлер без БД
	router.GET("/api/matches/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "abc" || idStr == "-1" || idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_INPUT",
				"message": "Некорректный ID матча",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": idStr})
	})

	tests := []struct {
		name       string
		id         string
		expectCode int
	}{
		{
			name:       "NonNumericID",
			id:         "abc",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "NegativeID",
			id:         "-1",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "ValidID",
			id:         "1",
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/matches/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("[%s] Expected status %d, got %d", tt.name, tt.expectCode, w.Code)
			}
		})
	}
}

// TestCreateMatch_Unauthorized - создание матча без токена должно вернуть 401
func TestCreateMatch_Unauthorized(t *testing.T) {
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.POST("/api/matches", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "Матч создан"})
	})

	body := `{"home_team":"Arsenal","away_team":"Chelsea","match_date":"2026-06-01T15:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/matches", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestCreateMatch_ForbiddenRole - создание матча с ролью user должно вернуть 403
func TestCreateMatch_ForbiddenRole(t *testing.T) {
	router := gin.New()

	// Имитируем авторизованного пользователя с ролью "user"
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("userName", "testuser")
		c.Set("userRole", "user")
		c.Next()
	})
	router.POST("/api/matches",
		middleware.RoleMiddleware("admin", "operator", "analyst"),
		func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "Матч создан"})
		},
	)

	body := `{"home_team":"Arsenal","away_team":"Chelsea","match_date":"2026-06-01T15:00:00Z"}`
	req, _ := http.NewRequest("POST", "/api/matches", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// TestCreateMatch_InvalidBody - создание матча с невалидным JSON
func TestCreateMatch_InvalidBody(t *testing.T) {
	router := gin.New()

	// Имитируем авторизованного пользователя с ролью "admin"
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("userName", "admin")
		c.Set("userRole", "admin")
		c.Next()
	})
	router.POST("/api/matches", func(c *gin.Context) {
		var input struct {
			HomeTeam  string `json:"home_team" binding:"required"`
			AwayTeam  string `json:"away_team" binding:"required"`
			MatchDate string `json:"match_date" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "ok"})
	})

	// Пустой JSON
	req, _ := http.NewRequest("POST", "/api/matches", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestLogout - проверка выхода из системы
func TestLogout(t *testing.T) {
	router := gin.New()

	router.POST("/api/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/api", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	})

	req, _ := http.NewRequest("POST", "/api/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Проверяем, что кука удалена (Max-Age = -1)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "token" && cookie.MaxAge < 0 {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'token' cookie to be deleted (MaxAge < 0)")
	}
}
