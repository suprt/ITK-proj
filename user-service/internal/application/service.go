package application

type Service struct {
	*RegisterUseCase
	*LoginUseCase
	*RefreshTokenUseCase
	*LogoutUseCase
	*GetSelfProfileUseCase
	*GetUserByIDUseCase
	*AddRoleUseCase
	*RemoveRoleUseCase
	*GetPublicProfileUseCase
	*UpdateProfileUseCase
}

func NewService(
	register *RegisterUseCase,
	login *LoginUseCase,
	refreshToken *RefreshTokenUseCase,
	logout *LogoutUseCase,
	getSelfProfile *GetSelfProfileUseCase,
	getUserByID *GetUserByIDUseCase,
	addRole *AddRoleUseCase,
	removeRole *RemoveRoleUseCase,
	getPublicProfile *GetPublicProfileUseCase,
	updateProfile *UpdateProfileUseCase,
) *Service {
	return &Service{
		RegisterUseCase:         register,
		LoginUseCase:            login,
		RefreshTokenUseCase:     refreshToken,
		LogoutUseCase:           logout,
		GetSelfProfileUseCase:   getSelfProfile,
		GetUserByIDUseCase:      getUserByID,
		AddRoleUseCase:          addRole,
		RemoveRoleUseCase:       removeRole,
		GetPublicProfileUseCase: getPublicProfile,
		UpdateProfileUseCase:    updateProfile,
	}
}
