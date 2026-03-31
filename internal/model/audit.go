package model

import (
	"time"
	"gorm.io/gorm"
)

// AuditLog - журнал действий (Требование ГОСТ Р ИСО/МЭК 27001)
type AuditLog struct {
	gorm.Model
	UserID    uint
	UserName  string // Дублируем имя на случай удаления пользователя
	Action    string // CREATE, UPDATE, DELETE, LOGIN
	Entity    string // Match, User
	EntityID  uint
	Timestamp time.Time `gorm:"autoCreateTime"`
}