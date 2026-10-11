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

type RegisterConfig struct {
	RefreshTTL time.Duration
}

type RegisterUseCase struct {
	users       outbound.UserRepository
	credentials outbound.CredentialsRepository
	sessions    outbound.SessionRepository
	hasher      outbound.PasswordHasher
	tokens      outbound.TokenIssuer
	ids         outbound.IDGenerator
	refreshTTL  time.Duration
}

func NewRegisterUseCase(
	users outbound.UserRepository,
	credentials outbound.CredentialsRepository,
	sessions outbound.SessionRepository,
	hasher outbound.PasswordHasher,
	tokens outbound.TokenIssuer,
	ids outbound.IDGenerator,
	cfg RegisterConfig,
) *RegisterUseCase {
	return &RegisterUseCase{
		users:       users,
		credentials: credentials,
		sessions:    sessions,
		hasher:      hasher,
		tokens:      tokens,
		ids:         ids,
		refreshTTL:  cfg.RefreshTTL,
	}
}

func (r *RegisterUseCase) Register(ctx context.Context, in dto.RegisterInput) (*dto.RegisterOutput, error) {

	_, err := r.users.GetByEmail(ctx, in.Email)
	if err == nil {
		return nil, er.ErrEmailTaken
	}
	if !errors.Is(err, er.ErrUserNotFound) {
		return nil, err
	}

	hash, err := r.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &domain.User{
		ID:           domain.UserID(r.ids.New()),
		Email:        in.Email,
		Roles:        []domain.Role{domain.RoleUser}, //всегда создаётся с User, другие роли могут быть выданы через GrantRole
		RolesVersion: 1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// UUIDv4 коллизия крайне маловероятна, но не невозможна
	if err := r.createUserWithRetry(ctx, user); err != nil {
		return nil, err
	}

	creds := &domain.Credentials{UserID: user.ID, PasswordHash: hash}
	if err := r.credentials.Save(ctx, creds); err != nil {
		return nil, err
	}

	access, err := r.tokens.IssueAccess(user.ID, user.Roles, user.RolesVersion)
	if err != nil {
		return nil, err
	}
	refresh, err := r.tokens.NewRefresh()
	if err != nil {
		return nil, err
	}

	session := &domain.RefreshSession{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(r.refreshTTL),
	}

	if err := r.sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	return &dto.RegisterOutput{
		UserID:       user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

const maxIDAttempts = 3

func (r *RegisterUseCase) createUserWithRetry(ctx context.Context, u *domain.User) error {
	var err error
	for range maxIDAttempts {
		u.ID = domain.UserID(r.ids.New())
		err = r.users.Create(ctx, u)
		if err == nil {
			return nil
		}
		if !errors.Is(err, er.ErrDuplicateID) {
			return err
		}
	}
	return err
}
