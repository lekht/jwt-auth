-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    failed_attempts INT DEFAULT 0
);

-- Создание таблицы сессий
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expiration_time TIMESTAMP
);

-- Создание таблицы аудита авторизации
CREATE TABLE IF NOT EXISTS audit (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(id),
    timestamp TIMESTAMP,
    event VARCHAR(50)
);
