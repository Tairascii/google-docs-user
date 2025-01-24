package usecase

import (
	"github.com/Tairascii/google-docs-user/internal/service/user"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	accessSecret  = "yoS0baK1Ya"
	refreshSecret = "NaRU70UzuMaK1"
)

type AuthUseCase interface {
	SignIn(email, password string) (Tokens, error)
	SignUp() error
}

type UseCase struct {
	users user.UserService
}

func NewAuthUseCase() AuthUseCase {
	return &UseCase{}
}

func (u *UseCase) SignIn(email, password string) (Tokens, error) {
	usr, err := u.users.GetUser(email, password)
	if err != nil {
		return Tokens{}, err
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    usr.ID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString([]byte(accessSecret))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    usr.ID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}).SignedString([]byte(refreshSecret))

	return Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (u *UseCase) SignUp() error {
	return nil
}
