package repository

import (
	"context"
	"jwt-auth/internal/models"
	"jwt-auth/package/postgres"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

const (
	Success = "successful login"
	Wrong   = "wrong password"
	Block   = "blocked"
)

func NewRepo(pg *postgres.PostgreDB) *AuthRepo {
	return &AuthRepo{pg}
}

type AuthRepo struct {
	*postgres.PostgreDB
}

const (
	sqlCreateUser = `INSERT INTO users (login, password) VALUES ($1,$2) returning id`

	sqlCheckUser = `SELECT id, login, password, login_attempts FROM users WHERE login = $1;`

	sqlSaveToken = `INSERT INTO sessions (user_id, token, expiration_time) VALUES ($1, $2, $3)`

	sqlCreateAuditRecord = `INSERT INTO audit (user_id, time, event) VALUES ($1,$2,$3)`

	sqlAttemptsUpdate = `UPDATE users SET login_attempts = login_attempts + 1 WHERE id = $1 RETURNING login_attempts;`

	sqlCheckToken = `SELECT id, user_id, token, expiration_time FROM sessions WHERE token = $1 `

	sqlAuditHistory = `SELECT id, user_id, time, event FROM audit WHERE user_id = $1`

	sqlDeleteAuditByUserID = `DELETE FROM audit WHERE user_id = $1`
)

// Добавляет нового пользователя.
func (a *AuthRepo) CreateUser(ctx context.Context, username, hashPassword string) error {
	var id int
	err := a.Pool.QueryRow(ctx, sqlCreateUser, username, hashPassword).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	audit := &models.AuthAudit{
		UserID:    id,
		Timestamp: time.Now().Unix(),
		Event:     Success,
	}

	err = a.storeAudit(ctx, audit)
	if err != nil {
		log.Println("failed to make audit record", err)
	}

	return nil
}

// Сохраняет токен при успешной аутентивигации.
func (a *AuthRepo) StoreToken(ctx context.Context, session *models.Session) error {

	_, err := a.Pool.Exec(ctx, sqlSaveToken, session.UserID, session.Token, session.ExpirationTime)
	if err != nil {
		return errors.Wrap(err, "failed to make token saving request")
	}

	audit := &models.AuthAudit{
		UserID:    session.UserID,
		Timestamp: time.Now().Unix(),
		Event:     Success,
	}

	err = a.storeAudit(ctx, audit)
	if err != nil {
		log.Println("failed to make audit record", err)
	}

	return nil
}

// Проверяет, существует ли в базе данных пользователь с таким именем. Возвращет всю информацию о пользователе.
func (a *AuthRepo) CheckUser(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	row := a.Pool.QueryRow(ctx, sqlCheckUser, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.FailedLoginAttempts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read row")
	}

	return &user, nil
}

// Сохраняет информацию о попытках аутентификации пользователя.
func (a *AuthRepo) storeAudit(ctx context.Context, audit *models.AuthAudit) error {
	_, err := a.Pool.Exec(ctx, sqlCreateAuditRecord, audit.UserID, audit.Timestamp, audit.Event)
	if err != nil {
		return errors.Wrap(err, "failed to make audit record")
	}

	return nil
}

// Увеличивает счетчик неудачных попыток аутентификации.
func (a *AuthRepo) WrongPassword(ctx context.Context, userId int) error {
	var attempts int
	err := a.Pool.QueryRow(ctx, sqlAttemptsUpdate, userId).Scan(&attempts)
	if err != nil {
		return errors.Wrap(err, "failed to update attempts")
	}

	if attempts >= 5 {
		err = a.storeAudit(ctx, &models.AuthAudit{
			UserID:    userId,
			Timestamp: time.Now().Unix(),
			Event:     Block,
		})
		if err != nil {
			log.Println("failed to make audit record", err)
		}
	}

	err = a.storeAudit(ctx, &models.AuthAudit{
		UserID:    userId,
		Timestamp: time.Now().Unix(),
		Event:     Wrong,
	})
	if err != nil {
		log.Println("failed to make audit record", err)
	}

	return nil
}

// Проверяет полученный токен на валидность.
func (a *AuthRepo) CheckToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	row := a.Pool.QueryRow(ctx, sqlCheckToken, token)
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.ExpirationTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read row")
	}

	if session.Token != token {
		return nil, nil
	}

	if time.Now().Unix() > session.ExpirationTime {
		return &models.Session{}, nil
	}

	return &session, nil

}

// Получает весь аудит пользователя.
func (a *AuthRepo) GetHistory(ctx context.Context, userId int) ([]models.AuthAudit, error) {
	rows, err := a.Pool.Query(ctx, sqlAuditHistory, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make audit query")
	}

	defer rows.Close()

	var history []models.AuthAudit
	for rows.Next() {
		var auditRecord models.AuthAudit
		err = rows.Scan(&auditRecord.ID, &auditRecord.UserID, &auditRecord.Timestamp, &auditRecord.Event)

		if err != nil {
			return nil, errors.Wrap(err, "failed to read rows")
		}
		history = append(history, auditRecord)
	}

	if err = rows.Err(); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	return history, nil
}

// Очищает аудит пользователя.
func (a *AuthRepo) DeleteHistory(ctx context.Context, userId int) error {
	_, err := a.Pool.Exec(ctx, sqlDeleteAuditByUserID, userId)
	if err != nil {
		return errors.Wrap(err, "failed to make delete query")
	}

	return nil
}
