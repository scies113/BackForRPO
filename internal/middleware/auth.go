package middleware

import (
	"BackendFootball/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

//Claims - структура токена
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware - проверка JWT токена
// Токен может быть передан:
// 1. В заголовке Authorization: Bearer <token>
// 2. В куки с именем "token"
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Пробуем получить токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Проверка формата "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// Если токен не найден в заголовке, пробуем куки
		if tokenString == "" {
			if cookie, err := c.Cookie("token"); err == nil {
				tokenString = cookie
			}
		}

		// Токен не найден нигде
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		// Парсинг и проверка токена
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контекст
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.Username)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RoleMiddleware - проверка роли (например, только оператор может создавать матчи)
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		// Проверка, есть ли роль пользователя в списке разрешённых
		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// GenerateToken - генерация JWT токена
func GenerateToken(user *model.User, roleName string) (string, error) {
	expireHours := 24

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}