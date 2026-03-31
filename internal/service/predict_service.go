package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/model"
	"time"
)

type PredictService struct{}

func NewPredictService() *PredictService {
	return &PredictService{}
}

// GetPrediction - получение прогноза для матча
// Если прогноз уже есть в БД - возвращает его
// Если нет - создаёт новый на основе "модели" (заглушка с вероятностями)
func (s *PredictService) GetPrediction(matchID uint, userID uint, userName string) (*model.Prediction, error) {
	var prediction model.Prediction

	// Проверяем, есть ли уже прогноз для этого матча
	result := database.DB.Where("match_id = ?", matchID).First(&prediction)

	if result.Error == nil {
		// Прогноз найден, возвращаем его
		return &prediction, nil
	}

	// Прогноза нет, создаём новый
	// Получаем матч для анализа
	var match model.Match
	if err := database.DB.First(&match, matchID).Error; err != nil {
		return nil, err
	}

	// "Модель ИИ" - заглушка с фиксированными вероятностями
	// В реальности здесь была бы загрузка обученной модели
	prediction = model.Prediction{
		MatchID:     matchID,
		HomeWinProb: 45.2, // Вероятность победы хозяев (%)
		DrawProb:    30.1, // Вероятность ничьей (%)
		AwayWinProb: 24.7, // Вероятность победы гостей (%)
		IsAccurate:  false,
	}

	// Сохраняем прогноз в БД
	if err := database.DB.Create(&prediction).Error; err != nil {
		return nil, err
	}

	// Аудит - запись о запросе к ИИ-модулю
	database.DB.Create(&model.AuditLog{
		UserID:   userID,
		UserName: userName,
		Action:   "PREDICT",
		Entity:   "Prediction",
		EntityID: prediction.ID,
		Timestamp: time.Now(),
	})

	return &prediction, nil
}
