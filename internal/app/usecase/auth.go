package usecase

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-user/internal/app/service/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	accessSecret  = "yoS0baK1Ya"
	refreshSecret = "NaRU70UzuMaK1"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type AuthUseCase interface {
	SignIn(ctx context.Context, email, password string) (Tokens, error)
	SignUp(ctx context.Context, data SignUpData) (Tokens, error)
}

type AuthUC struct {
	users user.UserService
}

func NewAuthUseCase(users user.UserService) AuthUseCase {
	return &AuthUC{users: users}
}

func generateToken(email, id, secret string, exp int64) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"iat":   time.Now().Unix(),
		"exp":   exp,
	}).SignedString([]byte(secret))
}

func checkPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func (u *AuthUC) SignIn(ctx context.Context, email, password string) (Tokens, error) {
	usr, err := u.users.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return Tokens{}, ErrUserNotFound
		}
		return Tokens{}, err
	}
	if err = checkPassword(password, usr.PasswordHash); err != nil {
		return Tokens{}, ErrUserNotFound
	}

	id := usr.ID.String()
	accessToken, err := generateToken(email, id, accessSecret, time.Now().Add(24*time.Hour).Unix())
	refreshToken, err := generateToken(email, id, refreshSecret, time.Now().Add(7*24*time.Hour).Unix())
	return Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (u *AuthUC) SignUp(ctx context.Context, data SignUpData) (Tokens, error) {
	id, err := u.users.CreateUser(ctx, user.CreateUserData(data))
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return Tokens{}, ErrUserAlreadyExists
		}
		return Tokens{}, err
	}

	accessToken, err := generateToken(data.Email, id.String(), accessSecret, time.Now().Add(24*time.Hour).Unix())
	refreshToken, err := generateToken(data.Email, id.String(), refreshSecret, time.Now().Add(7*24*time.Hour).Unix())
	return Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}
