package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/model"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// GetAllLogs - получение всех записей журнала аудита
func (s *AuditService) GetAllLogs() ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := database.DB.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}
