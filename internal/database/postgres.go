package database

import (
	"fmt"
	"log"
	"os"

	"BackendFootball/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect - подключение к PostgreSQL
func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Явно создаём таблицу roles (без AutoMigrate, чтобы не создавать лишние данные)
	err = DB.Table("roles").AutoMigrate(&model.Role{})
	if err != nil {
		return err
	}

	// Создаём роли по умолчанию
	createDefaultRoles()

	// Автоматическая миграция остальных таблиц
	err = DB.AutoMigrate(&model.User{}, &model.Match{}, &model.Prediction{}, &model.AuditLog{})
	if err != nil {
		return err
	}

	// Автоматически создаём дефолтного admin, если его ещё нет
	SeedDefaultAdmin()

	return nil
}

// createDefaultRoles - создание ролей по умолчанию (admin, operator, analyst, user)
func createDefaultRoles() {
	roles := []model.Role{
		{Name: "admin"},
		{Name: "operator"},
		{Name: "analyst"},
		{Name: "user"},
	}
	for _, role := range roles {
		DB.FirstOrCreate(&role, model.Role{Name: role.Name})
	}
}

// SeedDefaultAdmin - создаёт дефолтного admin при старте, если ни одного нет.
// Логин/пароль берётся из .env: ADMIN_EMAIL, ADMIN_PASSWORD.
// Если не заданы — используются значения по умолчанию.
func SeedDefaultAdmin() {
	var adminRole model.Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return // Роли ещё нет — пропускаем
	}

	// Проверяем: есть ли хоть один admin-пользователь?
	var count int64
	DB.Model(&model.User{}).Where("role_id = ?", adminRole.ID).Count(&count)
	if count > 0 {
		return // Admin уже существует — ничего не делаем
	}

	// Берём учётные данные из .env, иначе дефолт
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	username := os.Getenv("ADMIN_USERNAME")
	if email == "" {
		email = "admin@demo.com"
	}
	if password == "" {
		password = "admin123"
	}
	if username == "" {
		username = "admin"
	}

	admin := &model.User{
		Username: username,
		Email:    email,
		RoleID:   adminRole.ID,
	}
	if err := admin.SetPassword(password); err != nil {
		log.Printf("[SEED] Ошибка хеширования пароля admin: %v", err)
		return
	}

	if err := DB.Create(admin).Error; err != nil {
		log.Printf("[SEED] Ошибка создания дефолтного admin: %v", err)
		return
	}

	log.Printf("[SEED] ✅ Дефолтный admin создан: email=%s password=%s", email, password)
}