package handler

import (
	"BackendFootball/internal/errors"
	"BackendFootball/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{service: service.NewAuditService()}
}

// GetAuditLogs godoc
// @Summary Получение журнала аудита
// @Description Получение записей журнала аудита с фильтрацией и пагинацией (только для admin)
// @Tags audit
// @Accept json
// @Produce json
// @Param user_id query string false "Фильтр по ID пользователя"
// @Param action query string false "Фильтр по действию (CREATE, UPDATE, DELETE, LOGIN, PREDICT)"
// @Param date_from query string false "Дата начала (формат: 2006-01-02)"
// @Param date_to query string false "Дата окончания (формат: 2006-01-02)"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(50)
// @Success 200 {object} service.AuditResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/audit [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// Считываем query-параметры фильтрации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	filter := service.AuditFilter{
		UserID:   c.Query("user_id"),
		Action:   c.Query("action"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
	}

	result, err := h.service.GetLogsFiltered(filter)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, 200, result)
}
