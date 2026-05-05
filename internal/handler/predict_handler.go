package handler

import (
	"BackendFootball/internal/errors"
	"BackendFootball/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PredictHandler struct {
	service *service.PredictService
}

func NewPredictHandler() *PredictHandler {
	return &PredictHandler{service: service.NewPredictService()}
}

// GetPrediction godoc
// @Summary Получение прогноза ИИ
// @Description Получение прогноза результата матча от модуля ИИ (доступно: admin, operator, analyst)
// @Tags predictions
// @Produce json
// @Param id path int true "ID матча"
// @Success 200 {object} map[string]interface{} "Прогноз матча"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 403 {object} map[string]interface{} "Недостаточно прав"
// @Failure 404 {object} map[string]interface{} "Матч не найден"
// @Security BearerAuth
// @Router /api/predict/{id} [post]
func (h *PredictHandler) GetPrediction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError(
			errors.INVALID_INPUT, "Некорректный ID матча", http.StatusBadRequest,
		))
		return
	}

	// Получаем пользователя из контекста (для аудита)
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	prediction, err := h.service.GetPrediction(uint(id), userID.(uint), userName.(string))
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, http.StatusOK, gin.H{
		"match_id": prediction.MatchID,
		"prediction": gin.H{
			"home_win_prob": prediction.HomeWinProb,
			"draw_prob":     prediction.DrawProb,
			"away_win_prob": prediction.AwayWinProb,
		},
		"model_version": "v1.0",
		"generated_at":  prediction.CreatedAt,
	})
}
