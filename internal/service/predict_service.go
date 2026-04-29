package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// mlPredictURL — адрес Python-сервера с ML-моделью
const mlPredictURL = "http://localhost:8000/predict"

type PredictService struct{}

func NewPredictService() *PredictService {
	return &PredictService{}
}

// mlRequest — тело запроса к Python ML-серверу
type mlRequest struct {
	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`
}

// mlResponse — ответ от Python ML-сервера
type mlResponse struct {
	HomeWinProb float32 `json:"home_win_prob"`
	DrawProb    float32 `json:"draw_prob"`
	AwayWinProb float32 `json:"away_win_prob"`
}

// fetchMLPrediction — отправляет запрос к Python ML-серверу
// и возвращает вероятности исходов матча.
// В случае ошибки возвращает nil — вызывающий код использует запасные значения.
func fetchMLPrediction(homeTeam, awayTeam string) *mlResponse {
	reqBody, err := json.Marshal(mlRequest{
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	})
	if err != nil {
		fmt.Printf("[ML] Ошибка формирования запроса: %v\n", err)
		return nil
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(mlPredictURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("[ML] ML-сервер недоступен: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ML] ML-сервер вернул статус %d\n", resp.StatusCode)
		return nil
	}

	var mlResp mlResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResp); err != nil {
		fmt.Printf("[ML] Ошибка декодирования ответа: %v\n", err)
		return nil
	}

	return &mlResp
}

// GetPrediction — получение прогноза для матча.
// Если прогноз уже есть в БД — возвращает его.
// Если нет — запрашивает у Python ML-сервера и сохраняет в БД.
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

	// Запрашиваем прогноз у Python ML-сервера (XGBoost)
	var homeWin, draw, awayWin float32

	mlResult := fetchMLPrediction(match.HomeTeam, match.AwayTeam)
	if mlResult != nil {
		// ML-сервер ответил — используем реальные вероятности
		homeWin = mlResult.HomeWinProb * 100 // Переводим из долей в проценты
		draw = mlResult.DrawProb * 100
		awayWin = mlResult.AwayWinProb * 100
		fmt.Printf("[ML] Прогноз получен: %s vs %s → H:%.1f%% D:%.1f%% A:%.1f%%\n",
			match.HomeTeam, match.AwayTeam, homeWin, draw, awayWin)
	} else {
		// ML-сервер недоступен — используем запасные значения
		homeWin = 45.2
		draw = 30.1
		awayWin = 24.7
		fmt.Printf("[ML] Используются запасные значения для %s vs %s\n",
			match.HomeTeam, match.AwayTeam)
	}

	prediction = model.Prediction{
		MatchID:     matchID,
		HomeWinProb: homeWin,
		DrawProb:    draw,
		AwayWinProb: awayWin,
		IsAccurate:  false,
	}

	// Сохраняем прогноз в БД
	if err := database.DB.Create(&prediction).Error; err != nil {
		return nil, err
	}

	// Аудит — запись о запросе к ИИ-модулю
	database.DB.Create(&model.AuditLog{
		UserID:    userID,
		UserName:  userName,
		Action:    "PREDICT",
		Entity:    "Prediction",
		EntityID:  prediction.ID,
		Timestamp: time.Now(),
	})

	return &prediction, nil
}

