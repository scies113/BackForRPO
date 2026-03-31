package model

import (
	"time"
	"gorm.io/gorm"
)

// Match - матч АПЛ
type Match struct {
	gorm.Model
	HomeTeamID   uint      `gorm:"default:0" json:"home_team_id"`
	HomeTeam     string    `gorm:"not null" json:"home_team"` // Для упрощения пока строкой
	AwayTeamID   uint      `gorm:"default:0" json:"away_team_id"`
	AwayTeam     string    `gorm:"not null" json:"away_team"`
	MatchDate    time.Time `gorm:"not null" json:"match_date"`
	HomeScore    int       `gorm:"default:0" json:"home_score"`
	AwayScore    int       `gorm:"default:0" json:"away_score"`
	Status       string    `gorm:"default:scheduled" json:"status"` // scheduled, finished
	Prediction   *Prediction `json:"prediction"` // Связь 1 к 1 с прогнозом
}

// Prediction - прогноз ИИ
type Prediction struct {
	gorm.Model
	MatchID       uint    `gorm:"unique;not null" json:"match_id"`
	HomeWinProb   float32 `json:"home_win_prob"` // Вероятность победы хозяев %
	DrawProb      float32 `json:"draw_prob"` // Вероятность ничьей %
	AwayWinProb   float32 `json:"away_win_prob"` // Вероятность победы гостей %
	IsAccurate    bool    `json:"is_accurate"` // Сбылся ли прогноз (заполняется позже)
}