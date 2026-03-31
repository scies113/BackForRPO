package repository

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/model"
	"gorm.io/gorm"
)

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{db: database.DB}
}

// Create - создание матча
func (r *MatchRepository) Create(match *model.Match) error {
	return r.db.Create(match).Error
}

// GetAll - получение всех матчей
func (r *MatchRepository) GetAll() ([]model.Match, error) {
	var matches []model.Match
	err := r.db.Preload("Prediction").Find(&matches).Error
	return matches, err
}

// GetByID - поиск по ID
func (r *MatchRepository) GetByID(id uint) (*model.Match, error) {
	var match model.Match
	err := r.db.Preload("Prediction").First(&match, id).Error
	return &match, err
}

// Update - обновление матча
func (r *MatchRepository) Update(match *model.Match) error {
	return r.db.Save(match).Error
}

// Delete - удаление матча
func (r *MatchRepository) Delete(id uint) error {
	return r.db.Delete(&model.Match{}, id).Error
}