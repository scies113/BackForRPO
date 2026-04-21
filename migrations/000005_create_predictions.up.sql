-- Миграция: создание таблицы прогнозов ИИ
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    match_id INT UNIQUE NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    home_win_prob REAL DEFAULT 0,
    draw_prob REAL DEFAULT 0,
    away_win_prob REAL DEFAULT 0,
    is_accurate BOOLEAN DEFAULT FALSE
);

-- Индекс для soft delete (GORM)
CREATE INDEX IF NOT EXISTS idx_predictions_deleted_at ON predictions(deleted_at);
