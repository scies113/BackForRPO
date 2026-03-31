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

	// Автоматическая миграция таблиц (создание схемы)
	err = DB.AutoMigrate(&model.Role{}, &model.User{}, &model.Match{}, &model.Prediction{}, &model.AuditLog{})
	if err != nil {
		return err
	}

	return nil
}