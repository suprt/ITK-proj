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

type RegisterUseCase struct {
	users       outbound.UserRepository
	credentials outbound.CredentialsRepository
	hasher      outbound.PasswordHasher
	ids         outbound.IDGenerator
}

func NewRegisterUseCase(
	users outbound.UserRepository,
	credentials outbound.CredentialsRepository,
	hasher outbound.PasswordHasher,
	ids outbound.IDGenerator,
) *RegisterUseCase {
	return &RegisterUseCase{
		users:       users,
		credentials: credentials,
		hasher:      hasher,
		ids:         ids,
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

	if err := r.users.Create(ctx, user); err != nil {
		return nil, err
	}

	creds := &domain.Credentials{UserID: user.ID, PasswordHash: hash}
	if err := r.credentials.Save(ctx, creds); err != nil {
		return nil, err
	}

	return &dto.RegisterOutput{
		UserID: user.ID,
	}, nil
}
