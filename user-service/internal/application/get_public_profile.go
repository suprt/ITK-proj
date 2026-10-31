package application

import (
	"context"

	"github.com/suprt/trading/user-service/internal/core/dto"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type GetPublicProfileUseCase struct {
	users outbound.UserRepository
}

func NewGetPublicProfileUseCase(users outbound.UserRepository) *GetPublicProfileUseCase {
	return &GetPublicProfileUseCase{users: users}
}

func (uc *GetPublicProfileUseCase) GetPublicProfile(ctx context.Context,
	in dto.GetPublicProfileInput) (*dto.GetPublicProfileOutput, error) {
	user, err := uc.users.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.GetPublicProfileOutput{
		UserID:    user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
