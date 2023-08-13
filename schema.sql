DROP TABLE IF EXISTS users cascade;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS audit;

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    login_attempts INT DEFAULT 0
);

-- Создание таблицы сессий
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    token VARCHAR(255) NOT NULL,
    expiration_time BIGINT NOT NULL
);

-- Создание таблицы аудита авторизации
CREATE TABLE IF NOT EXISTS audit (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    time BIGINT NOT NULL,
    event VARCHAR(20) NOT NULL
);
