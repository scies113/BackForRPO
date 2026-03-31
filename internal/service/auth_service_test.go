package service

import (
	"BackendFootball/internal/model"
	"testing"
)

func TestSetPassword(t *testing.T) {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	password := "password123"

	// Тест 1: Пароль должен хешироваться
	err := user.SetPassword(password)
	if err != nil {
		t.Fatalf("SetPassword failed: %v", err)
	}

	// Тест 2: Хеш не должен совпадать с исходным паролем
	if user.PasswordHash == password {
		t.Error("Password hash should not equal plain password")
	}

	// Тест 3: CheckPassword должен возвращать true для правильного пароля
	if !user.CheckPassword(password) {
		t.Error("CheckPassword should return true for correct password")
	}

	// Тест 4: CheckPassword должен возвращать false для неправильного пароля
	if user.CheckPassword("wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestCheckPassword(t *testing.T) {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	password := "securepassword"
	user.SetPassword(password)

	// Тест 5: Проверка правильного пароля
	if !user.CheckPassword(password) {
		t.Error("Expected CheckPassword to return true")
	}

	// Тест 6: Проверка неправильного пароля
	if user.CheckPassword("wrongpassword") {
		t.Error("Expected CheckPassword to return false")
	}

	// Тест 7: Проверка пустого пароля
	if user.CheckPassword("") {
		t.Error("Expected CheckPassword to return false for empty password")
	}
}
