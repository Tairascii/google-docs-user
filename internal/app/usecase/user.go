package usecase

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-user/internal/app/service/user"
	"github.com/google/uuid"
)

type UserUseCase interface {
	IdByEmail(ctx context.Context, email string) (uuid.UUID, error)
}

type UserUC struct {
	user user.UserService
}

func NewUserUseCase(usr user.UserService) *UserUC {
	return &UserUC{user: usr}
}

func (u *UserUC) IdByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	usr, err := u.user.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return uuid.Nil, ErrUserNotFound
		}
		return uuid.Nil, err
	}
	return usr.ID, nil
}
