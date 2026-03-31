package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/errors"
	"BackendFootball/internal/middleware"
	"BackendFootball/internal/model"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register - регистрация пользователя с ролью по умолчанию (fan)
func (s *AuthService) Register(user *model.User, password string) (string, error) {
	// Хеширование пароля
	if err := user.SetPassword(password); err != nil {
		return "", err
	}

	// Получение роли "fan" (роль по умолчанию)
	var role model.Role
	result := database.DB.Where("name = ?", "fan").First(&role)
	if result.Error != nil {
		// Если роль не найдена, используем роль с наименьшим ID
		database.DB.Order("id ASC").First(&role)
	}

	user.RoleID = role.ID

	// Сохранение пользователя
	if err := database.DB.Create(user).Error; err != nil {
		return "", errors.ErrUserExists
	}

	// Генерация токена
	token, err := middleware.GenerateToken(user, role.Name)
	if err != nil {
		return "", errors.ErrInternalError
	}

	return token, nil
}

// RegisterWithRole - регистрация пользователя с указанной ролью (только для админа)
func (s *AuthService) RegisterWithRole(user *model.User, password, roleName string) (string, error) {
	// Хеширование пароля
	if err := user.SetPassword(password); err != nil {
		return "", err
	}

	// Получение роли
	var role model.Role
	result := database.DB.Where("name = ?", roleName).First(&role)
	if result.Error != nil {
		return "", errors.ErrInvalidRole
	}
	user.RoleID = role.ID

	// Сохранение пользователя
	if err := database.DB.Create(user).Error; err != nil {
		return "", errors.ErrUserExists
	}

	// Генерация токена
	token, err := middleware.GenerateToken(user, role.Name)
	if err != nil {
		return "", errors.ErrInternalError
	}

	return token, nil
}

// Login - аутентификация
func (s *AuthService) Login(email, password string) (string, error) {
	var user model.User
	var role model.Role

	// Поиск пользователя по email
	result := database.DB.Where("email = ?", email).Preload("Role").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", errors.ErrUserNotFound
		}
		return "", errors.ErrInvalidCredentials
	}

	// Проверка пароля
	if !user.CheckPassword(password) {
		return "", errors.ErrInvalidCredentials
	}

	// Загрузка роли
	database.DB.Where("id = ?", user.RoleID).First(&role)

	// Генерация токена
	token, err := middleware.GenerateToken(&user, role.Name)
	if err != nil {
		return "", errors.ErrInternalError
	}

	return token, nil
}

// createDefaultRoles - создание ролей по умолчанию
func (s *AuthService) createDefaultRoles() {
	roles := []model.Role{
		{Name: "admin"},
		{Name: "operator"},
		{Name: "analyst"},
		{Name: "fan"},
	}

	for _, role := range roles {
		database.DB.FirstOrCreate(&role, model.Role{Name: role.Name})
	}
}
