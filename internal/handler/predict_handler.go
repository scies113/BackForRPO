package handler

import (
	"BackendFootball/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PredictHandler struct {
	service *service.PredictService
}

func NewPredictHandler() *PredictHandler {
	return &PredictHandler{service: service.NewPredictService()}
}

// GetPrediction - получение прогноза ИИ для матча
func (h *PredictHandler) GetPrediction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Получаем пользователя из контекста (для аудита)
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	prediction, err := h.service.GetPrediction(uint(id), userID.(uint), userName.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"match_id": prediction.MatchID,
		"prediction": gin.H{
			"home_win": prediction.HomeWinProb,
			"draw":     prediction.DrawProb,
			"away_win": prediction.AwayWinProb,
		},
		"model_version": "v1.0",
		"generated_at":  prediction.CreatedAt,
	})
}
