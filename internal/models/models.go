package models

type User struct {
	ID                  int
	Username            string
	Password            string
	FailedLoginAttempts int
}

// Структура сессии
type Session struct {
	ID             int
	UserID         int
	Token          string
	ExpirationTime int64
}

// Структура аудита авторизации
type AuthAudit struct {
	ID        int
	UserID    int
	Timestamp int64
	Event     string
}
