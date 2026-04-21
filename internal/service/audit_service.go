package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/model"
	"strconv"
	"time"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// AuditFilter - параметры фильтрации аудита
type AuditFilter struct {
	UserID   string
	Action   string
	DateFrom string
	DateTo   string
	Page     int
	Limit    int
}

// AuditResponse - ответ с пагинацией
type AuditResponse struct {
	Data       []model.AuditLog `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// GetAllLogs - получение всех записей журнала аудита
func (s *AuditService) GetAllLogs() ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := database.DB.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

// GetLogsFiltered - получение записей с фильтрацией и пагинацией
func (s *AuditService) GetLogsFiltered(filter AuditFilter) (*AuditResponse, error) {
	// Значения по умолчанию
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	query := database.DB.Model(&model.AuditLog{})

	// Фильтр по user_id
	if filter.UserID != "" {
		if uid, err := strconv.ParseUint(filter.UserID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uint(uid))
		}
	}

	// Фильтр по действию (CREATE, UPDATE, DELETE, LOGIN, PREDICT)
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	// Фильтр по дате (от)
	if filter.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("timestamp >= ?", t)
		}
	}

	// Фильтр по дате (до)
	if filter.DateTo != "" {
		if t, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			// Включаем весь день date_to
			query = query.Where("timestamp < ?", t.AddDate(0, 0, 1))
		}
	}

	// Подсчёт общего количества записей
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Пагинация
	offset := (filter.Page - 1) * filter.Limit
	var logs []model.AuditLog
	err := query.Order("timestamp DESC").Offset(offset).Limit(filter.Limit).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	// Вычисление количества страниц
	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit != 0 {
		totalPages++
	}

	return &AuditResponse{
		Data:       logs,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}
