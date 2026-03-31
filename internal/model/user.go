package model

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Role - роль пользователя (Admin, Operator, Analyst, Fan)
type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"` // Например: "admin"
}

// User - пользователь системы
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	RoleID       uint
	Role         Role   `gorm:"foreignKey:RoleID"`
	AuditLogs    []AuditLog // Связь с журналом
}

// SetPassword - хеширование пароля (Требование безопасности ГОСТ)
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword - проверка пароля при входе
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}