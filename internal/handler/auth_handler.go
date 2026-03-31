package handler

import (
	"BackendFootball/internal/model"
	"BackendFootball/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{service: service.NewAuthService()}
}

// setAuthCookie - установка куки с токеном
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	// HttpOnly - JavaScript не имеет доступа (защита от XSS)
	// Secure - передаётся только по HTTPS (в продакшене)
	// SameSite=Lax - защита от CSRF
	
	c.SetCookie(
		"token",     // имя куки
		token,       // значение
		86400,       // срок жизни (24 часа в секундах)
		"/api",      // путь (только для /api endpoints)
		"",          // домен (пустой = текущий хост)
		false,       // Secure (false для localhost)
		true,        // HttpOnly
	)
}

// Register - регистрация нового пользователя (роль по умолчанию - user)
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

	// Устанавливаем куку с токеном
	h.setAuthCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered",
		"token":   token, // Возвращаем и в JSON для совместимости
		"role":    "user",
	})
}

// RegisterAdmin - регистрация пользователя с ролью (только для админа)
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin operator analyst"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
	}

	token, err := h.service.RegisterWithRole(user, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем куку с токеном
	h.setAuthCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered",
		"token":   token,
		"role":    input.Role,
	})
}

// Login - аутентификация пользователя
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем куку с токеном
	h.setAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
