package model

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Role - роль пользователя (Admin, Operator, Analyst, Fan)
type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"` // Например: "admin"
}

// User - пользователь системы
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null" json:"username"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	RoleID       uint   `json:"role_id"`
	Role         Role   `gorm:"foreignKey:RoleID" json:"role"`
	AuditLogs    []AuditLog `json:"audit_logs,omitempty"` // Связь с журналом
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