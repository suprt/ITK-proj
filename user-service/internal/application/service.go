package application

type Service struct {
	*RegisterUseCase
	*LoginUseCase
}

func NewService(register *RegisterUseCase, login *LoginUseCase) *Service {
	return &Service{
		RegisterUseCase: register,
		LoginUseCase:    login,
	}
}
