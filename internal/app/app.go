package app

import (
	"github.com/Tairascii/google-docs-user/internal/app/usecase"
)

type UseCase struct {
	Auth usecase.AuthUseCase
	User usecase.UserUseCase
}

type DI struct {
	UseCase UseCase
}
