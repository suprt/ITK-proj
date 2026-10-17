package application

import (
	"context"
	"errors"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
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
	callerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return er.ErrUnauthenticated
	}

	session, err := uc.sessions.GetByToken(ctx, in.RefreshToken)
	if err != nil {
		if errors.Is(err, er.ErrInvalidToken) {
			return nil
		}
		return err
	}

	if session.UserID != callerID {
		return er.ErrForbidden
	}

	return uc.sessions.Delete(ctx, in.RefreshToken)
}
