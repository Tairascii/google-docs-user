package user

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-user/internal/app/service/user/repo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrOnPassword        = errors.New("password error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidId         = errors.New("invalid id")
	ErrOnQuery           = errors.New("query error")
)

type UserService interface {
	CreateUser(ctx context.Context, data CreateUserData) (uuid.UUID, error)
	UserByEmail(ctx context.Context, email string) (User, error)
	UserByID(ctx context.Context) (User, error)
}

type Service struct {
	repo repo.UserRepo
}

func New(rp repo.UserRepo) *Service {
	return &Service{repo: rp}
}

func (s *Service) CreateUser(ctx context.Context, data CreateUserData) (uuid.UUID, error) {
	passHash, err := hashPassword(data.Password)
	if err != nil {
		return uuid.Nil, errors.Join(ErrOnPassword, err)
	}

	id, err := s.repo.CreateUser(ctx, repo.CreateUserData{
		Name:          data.Name,
		Email:         data.Email,
		Password:      passHash,
		ProfilePicUrl: data.ProfilePicUrl,
	})

	if errors.Is(err, repo.ErrUserAlreadyExists) {
		return uuid.Nil, ErrUserAlreadyExists
	}
	return id, err
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *Service) UserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return User(user), nil
}

func (s *Service) UserByID(ctx context.Context) (User, error) {
	id, ok := ctx.Value("id").(string)
	if !ok {
		return User{}, ErrInvalidId
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		return User{}, errors.Join(ErrInvalidId, err)
	}

	user, err := s.repo.UserByID(ctx, idParsed)
	if err != nil {
		return User{}, errors.Join(ErrOnQuery, err)
	}

	return User(user), nil
}
