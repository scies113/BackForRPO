package handler

import (
	"BackendFootball/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{service: service.NewAuditService()}
}

// GetAuditLogs - получение журнала действий (только для admin)
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	logs, err := h.service.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
