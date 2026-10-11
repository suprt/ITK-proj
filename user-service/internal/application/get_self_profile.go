package application

import (
	"context"
	"errors"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type GetSelfProfileUseCase struct {
	users outbound.UserRepository
}

func NewGetSelfProfileUseCase(users outbound.UserRepository) *GetSelfProfileUseCase {
	return &GetSelfProfileUseCase{users: users}
}

func (uc *GetSelfProfileUseCase) GetSelfProfile(ctx context.Context, _ dto.GetSelfProfileInput) (*dto.GetSelfProfileOutput, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, er.ErrUnauthenticated
	}

	user, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, er.ErrUserNotFound) {
			// Если по токену не находит юзера -> токен не валидный
			return nil, er.ErrInvalidToken
		}
		return nil, err
	}

	return &dto.GetSelfProfileOutput{User: user}, nil
}
