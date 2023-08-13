package auth

import (
	"context"
	"fmt"
	"jwt-auth/internal/models"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	AuthRepo
	Secret string
}

type Authentification interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	History(ctx context.Context, token string) ([]models.AuthAudit, error)
	Clear(token string) error
}

type AuthRepo interface {
	CreateUser(ctx context.Context, username, hashPassword string) error
	StoreToken(ctx context.Context, session *models.Session) error
	GetHistory(ctx context.Context, userId int) ([]models.AuthAudit, error)
	DeleteHistory(ctx context.Context, username string)
	CheckUser(ctx context.Context, username string) (*models.User, error)
	CheckToken(ctx context.Context, token string) (*models.Session, error)
	WrongPassword(ctx context.Context, userId int) error
}

func New(secret string, a AuthRepo) *Authenticator {
	return &Authenticator{a, secret}
}

func (a *Authenticator) SignUp(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	fmt.Println("user Creating... ", password, "    ", string(hash))
	err = a.CreateUser(ctx, username, string(hash))
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (a *Authenticator) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := a.CheckUser(ctx, username)
	if err != nil {
		log.Println("failed to validate user: ", err)
		return "", errors.New("bad request")
	}

	if user.ID == 0 {
		return "", errors.New("wrong login")
	}

	if user.FailedLoginAttempts >= 5 {
		return "", errors.New("User is blocked")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {

		if err = a.WrongPassword(ctx, user.ID); err != nil {
			log.Println("failed to record wrong passord event: ", err)
		}
		return "", errors.New("Wrong email or password")
	}

	// create token
	tokenExpirationTime := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": tokenExpirationTime,
	})

	// store token
	tokenString, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token string")
	}

	session := &models.Session{UserID: user.ID, Token: tokenString, ExpirationTime: tokenExpirationTime}

	err = a.StoreToken(ctx, session)
	if err != nil {
		return "", errors.Wrap(err, "failed to store token")
	}
	// return token
	return tokenString, nil
}

func (a *Authenticator) History(ctx context.Context, token string) ([]models.AuthAudit, error) {
	session, err := a.CheckToken(ctx, token)
	if err != nil {
		log.Println("failed to check token: ", err)
		return nil, errors.New("internal error")
	}

	if session.ID == 0 {
		return nil, errors.New("token expired")
	}

	// get audit history
	history, err := a.GetHistory(ctx, session.UserID)
	if err != nil {
		log.Println("failed to get history: ", err)
	}

	return history, nil
}

func (a *Authenticator) Clear(token string) error {
	// validate token

	// clear history

	return nil
}

// тут пишу бизнес логику (проверки хэша, обращение к бд)
