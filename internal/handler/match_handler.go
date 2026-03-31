package handler

import (
	"BackendFootball/internal/model"
	"BackendFootball/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{service: service.NewMatchService()}
}

// CreateMatch - HTTP обработчик создания матча
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var input model.Match
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя из контекста (после middleware авторизации)
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.CreateMatch(&input, userID.(uint), userName.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Матч создан"})
}

// GetAllMatches - получение списка матчей
func (h *MatchHandler) GetAllMatches(c *gin.Context) {
	matches, err := h.service.GetAllMatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch"})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetMatchByID - получение одного матча по ID
func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	match, err := h.service.GetMatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// UpdateMatch - обновление матча
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var input model.Match
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя из контекста
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.UpdateMatch(&input, uint(id), userID.(uint), userName.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Матч обновлён"})
}

// DeleteMatch - удаление матча
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Получаем пользователя из контекста
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.DeleteMatch(uint(id), userID.(uint), userName.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Матч удалён"})
}