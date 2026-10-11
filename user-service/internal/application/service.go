package application

type Service struct {
	*RegisterUseCase
	*LoginUseCase
	*RefreshTokenUseCase
	*LogoutUseCase
	*GetSelfProfileUseCase
	*GetUserByIDUseCase
}

func NewService(
	register *RegisterUseCase,
	login *LoginUseCase,
	refreshToken *RefreshTokenUseCase,
	logout *LogoutUseCase,
	getSelfProfile *GetSelfProfileUseCase,
	getUserByID *GetUserByIDUseCase,
) *Service {
	return &Service{
		RegisterUseCase:       register,
		LoginUseCase:          login,
		RefreshTokenUseCase:   refreshToken,
		LogoutUseCase:         logout,
		GetSelfProfileUseCase: getSelfProfile,
		GetUserByIDUseCase:    getUserByID,
	}
}
