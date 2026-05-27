package service

import (
	"BackendFootball/internal/database"
	"BackendFootball/internal/errors"
	"BackendFootball/internal/model"
	"BackendFootball/internal/repository"
	"fmt"
	"time"
)

type MatchService struct {
	repo *repository.MatchRepository
}

func NewMatchService() *MatchService {
	return &MatchService{repo: repository.NewMatchRepository()}
}

// CreateMatch - бизнес-процесс создания матча
func (s *MatchService) CreateMatch(match *model.Match, userID uint, userName string) error {
	// 1. Валидация
	if match.HomeTeam == "" || match.AwayTeam == "" {
		return fmt.Errorf("home_team and away_team are required")
	}

	if match.HomeTeam == match.AwayTeam {
		return errors.ErrSameTeams
	}

	// Проверка даты (не в прошлом)
	if match.MatchDate.Before(time.Now()) {
		return errors.ErrInvalidDate
	}

	// 2. Сохранение
	if err := s.repo.Create(match); err != nil {
		return err
	}

	// 3. Аудит (Обязательно по ТЗ!)
	s.logAudit(userID, userName, "CREATE", "Match", match.ID)

	return nil
}

// GetMatchByID - получение матча по ID
func (s *MatchService) GetMatchByID(id uint) (*model.Match, error) {
	match, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.ErrMatchNotFound
	}
	return match, nil
}

// UpdateMatch - бизнес-процесс обновления матча
func (s *MatchService) UpdateMatch(match *model.Match, id uint, userID uint, userName string) error {
	// 1. Проверка существования
	existingMatch, err := s.repo.GetByID(id)
	if err != nil {
		return errors.ErrMatchNotFound
	}

	// 2. Валидация
	if match.HomeTeam == "" || match.AwayTeam == "" {
		return fmt.Errorf("home_team and away_team are required")
	}

	if match.HomeTeam == match.AwayTeam {
		return errors.ErrSameTeams
	}

	// 3. Обновление полей
	existingMatch.HomeTeam = match.HomeTeam
	existingMatch.AwayTeam = match.AwayTeam
	existingMatch.MatchDate = match.MatchDate

	if err := s.repo.Update(existingMatch); err != nil {
		return err
	}

	// 4. Аудит
	s.logAudit(userID, userName, "UPDATE", "Match", id)

	return nil
}

// DeleteMatch - бизнес-процесс удаления матча
func (s *MatchService) DeleteMatch(id uint, userID uint, userName string) error {
	// 1. Проверка существования
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.ErrMatchNotFound
	}

	// 2. Удаление
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 3. Аудит
	s.logAudit(userID, userName, "DELETE", "Match", id)

	return nil
}

// logAudit - вспомогательная функция записи в журнал
func (s *MatchService) logAudit(userID uint, userName, action, entity string, entityID uint) {
	database.DB.Create(&model.AuditLog{
		UserID:    userID,
		UserName:  userName,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Timestamp: time.Now(),
	})
}

// GetAllMatches - получение всех матчей
func (s *MatchService) GetAllMatches() ([]model.Match, error) {
	return s.repo.GetAll()
}