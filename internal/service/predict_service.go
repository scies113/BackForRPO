package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/errors"
	"BackendFootball/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// getMLPredictURL — возвращает URL ML-сервера из переменной окружения.
// По умолчанию используется http://localhost:8000/predict (для локальной разработки).
// В Docker-окружении устанавливается http://ml:8000/predict через docker-compose.
func getMLPredictURL() string {
	url := os.Getenv("ML_PREDICT_URL")
	if url == "" {
		return "http://localhost:8000/predict"
	}
	return url
}

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

// mlErrorResponse — ответ об ошибке от Python ML-сервера
type mlErrorResponse struct {
	Detail string `json:"detail"`
}

// fetchMLPrediction — отправляет запрос к Python ML-серверу
// и возвращает вероятности исходов матча.
// Возвращает ошибку, если ML-сервер недоступен или вернул ошибку.
func fetchMLPrediction(homeTeam, awayTeam string) (*mlResponse, error) {
	reqBody, err := json.Marshal(mlRequest{
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования запроса к ML-серверу: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(getMLPredictURL(), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ML-сервер недоступен: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Пытаемся извлечь сообщение об ошибке от ML-сервера
		var mlErr mlErrorResponse
		if json.Unmarshal(body, &mlErr) == nil && mlErr.Detail != "" {
			return nil, errors.NewAppError(
				errors.PREDICTION_ERROR,
				fmt.Sprintf("Ошибка ML-модели: %s", mlErr.Detail),
				resp.StatusCode,
			)
		}
		return nil, errors.NewAppError(
			errors.PREDICTION_ERROR,
			fmt.Sprintf("ML-сервер вернул ошибку (HTTP %d)", resp.StatusCode),
			http.StatusBadGateway,
		)
	}

	var mlResp mlResponse
	if err := json.Unmarshal(body, &mlResp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа ML-сервера: %w", err)
	}

	return &mlResp, nil
}

// GetPrediction — получение прогноза для матча.
// Если прогноз уже есть в БД — возвращает его.
// Если нет — запрашивает у Python ML-сервера и сохраняет в БД.
// При недоступности ML-сервера возвращает ошибку (без fallback-значений).
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
		return nil, errors.ErrMatchNotFound
	}

	// Запрашиваем прогноз у Python ML-сервера (XGBoost)
	mlResult, err := fetchMLPrediction(match.HomeTeam, match.AwayTeam)
	if err != nil {
		fmt.Printf("[ML] Ошибка получения прогноза для %s vs %s: %v\n",
			match.HomeTeam, match.AwayTeam, err)
		return nil, err
	}

	// ML-сервер ответил — используем реальные вероятности
	homeWin := mlResult.HomeWinProb * 100 // Переводим из долей в проценты
	draw := mlResult.DrawProb * 100
	awayWin := mlResult.AwayWinProb * 100
	fmt.Printf("[ML] Прогноз получен: %s vs %s → H:%.1f%% D:%.1f%% A:%.1f%%\n",
		match.HomeTeam, match.AwayTeam, homeWin, draw, awayWin)

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
