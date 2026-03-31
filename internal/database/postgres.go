package database

import (
	"fmt"
	"BackendFootball/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
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

	return nil
}

// createDefaultRoles - создание ролей по умолчанию (admin, operator, analyst, user)
func createDefaultRoles() {
	roles := []model.Role{
		{Name: "admin"},
		{Name: "operator"},
		{Name: "analyst"},
		{Name: "user"}, // user = fan (обычный пользователь)
	}

	for _, role := range roles {
		// Используем FirstOrCreate, чтобы не дублировать роли
		DB.FirstOrCreate(&role, model.Role{Name: role.Name})
	}
}