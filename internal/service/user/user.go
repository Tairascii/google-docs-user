package user

import (
	"errors"
	"github.com/Tairascii/google-docs-user/internal/service/user/repo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrOnPassword        = errors.New("password error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserService interface {
	CreateUser(name, email, password, picUrl string) error
	GetUser(email, password string) (User, error)
}

type Service struct {
	repo repo.UserRepo
}

func New(rp repo.UserRepo) *Service {
	return &Service{repo: rp}
}

func (s *Service) CreateUser(name, email, password, picUrl string) error {
	passHash, err := hashPassword(password)
	if err != nil {
		return errors.Join(ErrOnPassword, err)
	}

	err = s.repo.CreateUser(repo.CreateUserData{
		Name:          name,
		Email:         email,
		Password:      passHash,
		ProfilePicUrl: picUrl,
	})

	if errors.Is(err, repo.ErrUserAlreadyExists) {
		return ErrUserAlreadyExists
	}
	return err
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *Service) GetUser(email, password string) (User, error) {
	passHash, err := hashPassword(password)
	if err != nil {
		return User{}, errors.Join(ErrOnPassword, err)
	}

	user, err := s.repo.GetUser(email, passHash)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return User(user), nil
}
