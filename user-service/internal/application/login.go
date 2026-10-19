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
	AccessTTL  time.Duration
}

type LoginUseCase struct {
	users      outbound.UserRepository
	cred       outbound.CredentialsRepository
	sessions   outbound.SessionRepository
	hasher     outbound.PasswordHasher
	tokens     outbound.TokenIssuer
	uow        outbound.UnitOfWork
	refreshTTL time.Duration
	accessTTL  time.Duration
}

func NewLoginUseCase(
	users outbound.UserRepository,
	cred outbound.CredentialsRepository,
	sessions outbound.SessionRepository,
	hasher outbound.PasswordHasher,
	tokens outbound.TokenIssuer,
	uow outbound.UnitOfWork,
	cfg LoginConfig,
) *LoginUseCase {
	return &LoginUseCase{
		users:      users,
		cred:       cred,
		sessions:   sessions,
		hasher:     hasher,
		tokens:     tokens,
		uow:        uow,
		refreshTTL: cfg.RefreshTTL,
		accessTTL:  cfg.AccessTTL,
	}
}

func (l *LoginUseCase) Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error) {
	var user *domain.User
	var access domain.AccessToken
	var refresh domain.RefreshToken

	err := l.uow.Do(ctx, func(ctx context.Context) error {
		var err error
		user, err = l.users.GetByEmail(ctx, in.Email)
		if err != nil {
			if errors.Is(err, er.ErrUserNotFound) {
				return er.ErrInvalidCredentials
			}
			return err
		}

		creds, err := l.cred.GetByID(ctx, user.ID)
		if err != nil {
			if errors.Is(err, er.ErrUserNotFound) {
				return er.ErrInvalidCredentials
			}
			return err
		}

		if err := l.hasher.Compare(creds.PasswordHash, in.Password); err != nil {
			return er.ErrInvalidCredentials
		}

		access, err = l.tokens.IssueAccess(user.ID, user.Roles, user.RolesVersion)
		if err != nil {
			return err
		}

		refresh, err = l.tokens.NewRefresh()
		if err != nil {
			return err
		}

		now := time.Now()
		session := &domain.RefreshSession{
			Token:     refresh,
			UserID:    user.ID,
			ExpiresAt: now.Add(l.refreshTTL),
			IssuedAt:  now,
		}

		return l.sessions.Save(ctx, session)
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		UserID:          user.ID,
		AccessToken:     access,
		RefreshToken:    refresh,
		AccessExpiresIn: int64(l.accessTTL.Seconds()),
	}, nil
}
