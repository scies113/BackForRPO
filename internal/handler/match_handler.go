package handler

import (
	"BackendFootball/internal/errors"
	"BackendFootball/internal/model"
	"BackendFootball/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{service: service.NewMatchService()}
}

// CreateMatch godoc
// @Summary Создание матча
// @Description Создание нового матча (доступно: admin, operator)
// @Tags matches
// @Accept json
// @Produce json
// @Param input body model.Match true "Данные матча"
// @Success 201 {object} map[string]interface{} "Матч создан"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 403 {object} map[string]interface{} "Недостаточно прав"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка"
// @Security BearerAuth
// @Router /api/matches [post]
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var input model.Match
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.RespondWithError(c, errors.NewAppErrorWithDetails(
			errors.INVALID_INPUT, "Некорректные данные матча", http.StatusBadRequest, err.Error(),
		))
		return
	}

	// Получаем пользователя из контекста (после middleware авторизации)
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.CreateMatch(&input, userID.(uint), userName.(string)); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "Матч создан", "id": input.ID})
}

// GetAllMatches godoc
// @Summary Получение всех матчей
// @Description Получение списка всех матчей с прогнозами
// @Tags matches
// @Produce json
// @Success 200 {array} model.Match "Список матчей"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка"
// @Security BearerAuth
// @Router /api/matches [get]
func (h *MatchHandler) GetAllMatches(c *gin.Context) {
	matches, err := h.service.GetAllMatches()
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}
	errors.RespondWithSuccess(c, http.StatusOK, matches)
}

// GetMatchByID godoc
// @Summary Получение матча по ID
// @Description Получение информации о конкретном матче по его ID
// @Tags matches
// @Produce json
// @Param id path int true "ID матча"
// @Success 200 {object} model.Match "Данные матча"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 404 {object} map[string]interface{} "Матч не найден"
// @Security BearerAuth
// @Router /api/matches/{id} [get]
func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError(
			errors.INVALID_INPUT, "Некорректный ID матча", http.StatusBadRequest,
		))
		return
	}

	match, err := h.service.GetMatchByID(uint(id))
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, http.StatusOK, match)
}

// UpdateMatch godoc
// @Summary Обновление матча
// @Description Обновление данных матча по ID (доступно: admin, operator)
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID матча"
// @Param input body model.Match true "Обновлённые данные матча"
// @Success 200 {object} map[string]interface{} "Матч обновлён"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 403 {object} map[string]interface{} "Недостаточно прав"
// @Failure 404 {object} map[string]interface{} "Матч не найден"
// @Security BearerAuth
// @Router /api/matches/{id} [put]
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError(
			errors.INVALID_INPUT, "Некорректный ID матча", http.StatusBadRequest,
		))
		return
	}

	var input model.Match
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.RespondWithError(c, errors.NewAppErrorWithDetails(
			errors.INVALID_INPUT, "Некорректные данные матча", http.StatusBadRequest, err.Error(),
		))
		return
	}

	// Получаем пользователя из контекста
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.UpdateMatch(&input, uint(id), userID.(uint), userName.(string)); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Матч обновлён"})
}

// DeleteMatch godoc
// @Summary Удаление матча
// @Description Удаление матча по ID (доступно: admin, operator)
// @Tags matches
// @Produce json
// @Param id path int true "ID матча"
// @Success 200 {object} map[string]interface{} "Матч удалён"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 403 {object} map[string]interface{} "Недостаточно прав"
// @Failure 404 {object} map[string]interface{} "Матч не найден"
// @Security BearerAuth
// @Router /api/matches/{id} [delete]
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError(
			errors.INVALID_INPUT, "Некорректный ID матча", http.StatusBadRequest,
		))
		return
	}

	// Получаем пользователя из контекста
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	if err := h.service.DeleteMatch(uint(id), userID.(uint), userName.(string)); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Матч удалён"})
}