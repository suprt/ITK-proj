package application

type Service struct {
	*RegisterUseCase
	*LoginUseCase
	*RefreshTokenUseCase
	*LogoutUseCase
	*GetSelfProfileUseCase
}

func NewService(
	register *RegisterUseCase,
	login *LoginUseCase,
	refreshToken *RefreshTokenUseCase,
	logout *LogoutUseCase,
	getSelfProfile *GetSelfProfileUseCase,
) *Service {
	return &Service{
		RegisterUseCase:       register,
		LoginUseCase:          login,
		RefreshTokenUseCase:   refreshToken,
		LogoutUseCase:         logout,
		GetSelfProfileUseCase: getSelfProfile,
	}
}
