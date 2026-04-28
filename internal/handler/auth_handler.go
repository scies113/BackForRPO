package handler

import (
	"BackendFootball/internal/errors"
	"BackendFootball/internal/model"
	"BackendFootball/internal/service"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{service: service.NewAuthService()}
}

// isSecureCookie — проверяет переменную окружения COOKIE_SECURE
// В продакшене (HTTPS) должна быть true, для localhost — false
func isSecureCookie() bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SECURE")))
	return val == "true" || val == "1"
}

// setAuthCookie — установка HttpOnly куки с токеном авторизации
// Безопасность:
//   - HttpOnly: true — JavaScript не имеет доступа (защита от XSS)
//   - Secure: env-driven — передаётся только по HTTPS в продакшене
//   - SameSite=Lax — защита от CSRF, разрешает top-level навигацию
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/api",
		MaxAge:   86400, // 24 часа
		HttpOnly: true,  // Блокирует доступ из JS (защита от XSS)
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteLaxMode, // Защита от CSRF
	})
}

// clearAuthCookie — удаление куки авторизации (logout)
func (h *AuthHandler) clearAuthCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		MaxAge:   -1,    // Немедленное удаление
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteLaxMode,
	})
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя с ролью по умолчанию (user)
// @Tags auth
// @Accept json
// @Produce json
// @Param input body object true "Данные пользователя" example({"username":"user1","email":"user@example.com","password":"password123"})
// @Success 201 {object} map[string]interface{} "Пользователь зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 409 {object} map[string]interface{} "Пользователь уже существует"
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		errors.RespondWithError(c, errors.NewAppErrorWithDetails(
			errors.INVALID_INPUT, "Некорректные входные данные", http.StatusBadRequest, err.Error(),
		))
		return
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
	}

	token, err := h.service.Register(user, input.Password)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Устанавливаем безопасную куку с токеном
	h.setAuthCookie(c, token)

	errors.RespondWithSuccess(c, http.StatusCreated, gin.H{
		"message": "User registered",
		"token":   token,
		"role":    "user",
	})
}

// RegisterAdmin godoc
// @Summary Регистрация пользователя с ролью (админ)
// @Description Регистрация пользователя с указанной ролью (только для администратора)
// @Tags auth
// @Accept json
// @Produce json
// @Param input body object true "Данные пользователя с ролью" example({"username":"op1","email":"op@example.com","password":"password123","role":"operator"})
// @Success 201 {object} map[string]interface{} "Пользователь зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 403 {object} map[string]interface{} "Недостаточно прав"
// @Failure 409 {object} map[string]interface{} "Пользователь уже существует"
// @Security BearerAuth
// @Router /api/admin/register [post]
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin operator analyst user"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		errors.RespondWithError(c, errors.NewAppErrorWithDetails(
			errors.INVALID_INPUT, "Некорректные входные данные", http.StatusBadRequest, err.Error(),
		))
		return
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
	}

	token, err := h.service.RegisterWithRole(user, input.Password, input.Role)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Устанавливаем безопасную куку с токеном
	h.setAuthCookie(c, token)

	errors.RespondWithSuccess(c, http.StatusCreated, gin.H{
		"message": "User registered",
		"token":   token,
		"role":    input.Role,
	})
}

// Login godoc
// @Summary Вход в систему
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param input body object true "Данные для входа" example({"email":"user@example.com","password":"password123"})
// @Success 200 {object} map[string]interface{} "Токен авторизации"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Неверные учётные данные"
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		errors.RespondWithError(c, errors.NewAppErrorWithDetails(
			errors.INVALID_INPUT, "Некорректные входные данные", http.StatusBadRequest, err.Error(),
		))
		return
	}

	token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	// Устанавливаем безопасную куку с токеном
	h.setAuthCookie(c, token)

	errors.RespondWithSuccess(c, http.StatusOK, gin.H{
		"token": token,
	})
}

// Logout godoc
// @Summary Выход из системы
// @Description Удаление JWT cookie для деавторизации
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Выход выполнен"
// @Router /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Удаляем куку безопасным способом (с теми же флагами)
	h.clearAuthCookie(c)
	errors.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetMe godoc
// @Summary Получение профиля текущего пользователя
// @Description Возвращает данные авторизованного пользователя из JWT
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Профиль пользователя"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Security BearerAuth
// @Router /api/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")
	userRole, _ := c.Get("userRole")

	errors.RespondWithSuccess(c, http.StatusOK, gin.H{
		"user_id":  userID,
		"username": userName,
		"role":     userRole,
	})
}
