package model

import (
	"time"
	"gorm.io/gorm"
)

// AuditLog - журнал действий (Требование ГОСТ Р ИСО/МЭК 27001)
type AuditLog struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"` // Дублируем имя на случай удаления пользователя
	Action    string    `json:"action"`    // CREATE, UPDATE, DELETE, LOGIN
	Entity    string    `json:"entity_type"` // Match, User
	EntityID  uint      `json:"entity_id"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}