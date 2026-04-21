-- Миграция: создание таблицы матчей
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    home_team_id INT DEFAULT 0,
    home_team VARCHAR(100) NOT NULL,
    away_team_id INT DEFAULT 0,
    away_team VARCHAR(100) NOT NULL,
    match_date TIMESTAMP WITH TIME ZONE NOT NULL,
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'scheduled'
);

-- Индекс для soft delete (GORM)
CREATE INDEX IF NOT EXISTS idx_matches_deleted_at ON matches(deleted_at);
