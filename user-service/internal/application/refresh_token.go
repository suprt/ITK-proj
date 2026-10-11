package application

import (
	"context"
	"time"

	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type RefreshTokenConfig struct {
	RefreshTTL time.Duration
}

type RefreshTokenUseCase struct {
	users      outbound.UserRepository
	sessions   outbound.SessionRepository
	tokens     outbound.TokenIssuer
	refreshTTL time.Duration
}

func NewRefreshTokenUseCase(
	users outbound.UserRepository,
	sessions outbound.SessionRepository,
	tokens outbound.TokenIssuer,
	cfg RefreshTokenConfig) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		users:      users,
		sessions:   sessions,
		tokens:     tokens,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (r *RefreshTokenUseCase) RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error) {
	session, err := r.sessions.GetByToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, er.ErrInvalidToken
	}

	if session.Expired(time.Now()) {
		err = r.sessions.Delete(ctx, in.RefreshToken)
		if err != nil {
			return nil, err
		}
		return nil, er.ErrExpiredToken
	}

	user, err := r.users.GetByID(ctx, session.UserID)
	if err != nil {
		_ = r.sessions.Delete(ctx, in.RefreshToken)
		return nil, er.ErrInvalidToken
	}

	if err := r.sessions.Delete(ctx, in.RefreshToken); err != nil {
		return nil, err
	}

	newRefresh, err := r.tokens.NewRefresh()
	if err != nil {
		return nil, err
	}

	newSession := &domain.RefreshSession{
		Token:     newRefresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(r.refreshTTL),
	}
	if err := r.sessions.Save(ctx, newSession); err != nil {
		return nil, err
	}

	access, err := r.tokens.IssueAccess(user.ID, user.Roles)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenOutput{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}
