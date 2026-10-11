package application

import (
	"context"

	"github.com/suprt/trading/user-service/internal/core/dto"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type LogoutUseCase struct {
	sessions outbound.SessionRepository
}

func NewLogoutUseCase(sessions outbound.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{
		sessions: sessions,
	}
}

func (uc *LogoutUseCase) Logout(ctx context.Context, in dto.LogoutInput) error {
	return uc.sessions.Delete(ctx, in.RefreshToken)
}
