package usecase

type AuthUseCase interface {
	SignIn() error
}

type UseCase struct {
}

func NewAuthUseCase() AuthUseCase {
	return &UseCase{}
}

func (u *UseCase) SignIn() error {
	return nil
}
