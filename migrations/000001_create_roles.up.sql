-- Миграция: создание таблицы ролей
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Вставка ролей по умолчанию
INSERT INTO roles (name) VALUES ('admin'), ('operator'), ('analyst'), ('user')
ON CONFLICT (name) DO NOTHING;
