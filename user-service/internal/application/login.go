package application

import (
	"context"
	"errors"
	"time"

	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type LoginConfig struct {
	RefreshTTL time.Duration
}

type LoginUseCase struct {
	users      outbound.UserRepository
	cred       outbound.CredentialsRepository
	sessions   outbound.SessionRepository
	hasher     outbound.PasswordHasher
	tokens     outbound.TokenIssuer
	refreshTTL time.Duration
}

func NewLoginUseCase(
	users outbound.UserRepository,
	cred outbound.CredentialsRepository,
	sessions outbound.SessionRepository,
	hasher outbound.PasswordHasher,
	tokens outbound.TokenIssuer,
	cfg LoginConfig,
) *LoginUseCase {
	return &LoginUseCase{
		users:      users,
		cred:       cred,
		sessions:   sessions,
		hasher:     hasher,
		tokens:     tokens,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (l *LoginUseCase) Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error) {
	user, err := l.users.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, er.ErrUserNotFound) {
			return nil, er.ErrInvalidCredentials
		}
		return nil, err
	}

	creds, err := l.cred.GetByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, er.ErrUserNotFound) {
			return nil, er.ErrInvalidCredentials
		}
		return nil, err
	}
	if err := l.hasher.Compare(creds.PasswordHash, in.Password); err != nil {
		return nil, er.ErrInvalidCredentials
	}

	access, err := l.tokens.IssueAccess(user.ID, user.Roles)
	if err != nil {
		return nil, err
	}
	refresh, err := l.tokens.NewRefresh()
	if err != nil {
		return nil, err
	}

	session := &domain.RefreshSession{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(l.refreshTTL),
	}
	if err := l.sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		UserID:       user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
