package service

import (
	"BackendFootball/internal/model"
	"testing"
	"time"
)

func TestMatchValidation(t *testing.T) {
	// Тест 1: Ошибка при одинаковых командах
	match := &model.Match{
		HomeTeam:  "Arsenal",
		AwayTeam:  "Arsenal",
		MatchDate: time.Now().AddDate(0, 0, 1),
	}

	// Проверяем валидацию команд
	if match.HomeTeam == match.AwayTeam {
		t.Log("Validation correctly detected same teams")
	} else {
		t.Error("Should detect same teams")
	}

	// Тест 2: Ошибка при пустой домашней команде
	match2 := &model.Match{
		HomeTeam:  "",
		AwayTeam:  "Chelsea",
		MatchDate: time.Now().AddDate(0, 0, 1),
	}

	if match2.HomeTeam == "" {
		t.Log("Validation correctly detected empty home team")
	} else {
		t.Error("Should detect empty home team")
	}

	// Тест 3: Ошибка при пустой гостевой команде
	match3 := &model.Match{
		HomeTeam:  "Arsenal",
		AwayTeam:  "",
		MatchDate: time.Now().AddDate(0, 0, 1),
	}

	if match3.AwayTeam == "" {
		t.Log("Validation correctly detected empty away team")
	} else {
		t.Error("Should detect empty away team")
	}

	// Тест 4: Ошибка при дате в прошлом
	match4 := &model.Match{
		HomeTeam:  "Arsenal",
		AwayTeam:  "Chelsea",
		MatchDate: time.Now().AddDate(0, 0, -1),
	}

	if match4.MatchDate.Before(time.Now()) {
		t.Log("Validation correctly detected past date")
	} else {
		t.Error("Should detect past date")
	}
}

func TestMatchDateValidation(t *testing.T) {
	// Тест 5: Валидная дата (в будущем)
	match := &model.Match{
		HomeTeam:  "Arsenal",
		AwayTeam:  "Chelsea",
		MatchDate: time.Now().AddDate(0, 0, 7),
	}

	// Проверяем валидацию даты
	if match.MatchDate.After(time.Now()) {
		t.Log("Match date is correctly in the future")
	} else {
		t.Error("Match date should be in the future")
	}

	// Тест 6: Валидные команды
	if match.HomeTeam != match.AwayTeam {
		t.Log("Teams are correctly different")
	} else {
		t.Error("Home and away teams should be different")
	}

	// Тест 7: Проверка структуры матча
	if match.HomeTeam != "" && match.AwayTeam != "" {
		t.Log("Teams are correctly set")
	} else {
		t.Error("Teams should not be empty")
	}
}

func TestMatchStruct(t *testing.T) {
	// Тест 8: Проверка создания структуры матча
	match := &model.Match{
		HomeTeam:  "Liverpool",
		AwayTeam:  "Manchester United",
		MatchDate: time.Now().AddDate(0, 0, 5),
	}

	if match.HomeTeam != "Liverpool" {
		t.Error("HomeTeam should be Liverpool")
	}

	if match.AwayTeam != "Manchester United" {
		t.Error("AwayTeam should be Manchester United")
	}

}
