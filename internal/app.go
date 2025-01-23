package app

import "github.com/Tairascii/google-docs-user/internal/usecase"

type UseCase struct {
	Auth usecase.AuthUseCase
}

type DI struct {
	UseCase UseCase
}
