package application

import (
	"context"
	"errors"
	"time"

	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type SeedUserInput struct {
	Email    string
	Password string
	Roles    []domain.Role
}

type Seed struct {
	user   outbound.UserRepository
	cred   outbound.CredentialsRepository
	hasher outbound.PasswordHasher
	ids    outbound.IDGenerator
}

func NewSeed(
	user outbound.UserRepository,
	cred outbound.CredentialsRepository,
	hasher outbound.PasswordHasher,
	ids outbound.IDGenerator,
) *Seed {
	return &Seed{
		user:   user,
		cred:   cred,
		hasher: hasher,
		ids:    ids,
	}
}

func (s *Seed) SeedUser(ctx context.Context, in SeedUserInput) error {
	_, err := s.user.GetByEmail(ctx, in.Email)
	if err == nil {
		return nil
	}
	if !errors.Is(err, er.ErrUserNotFound) {
		return err
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &domain.User{
		ID:        domain.UserID(s.ids.New()),
		Email:     in.Email,
		Roles:     in.Roles,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.user.Create(ctx, user); err != nil {
		return err
	}

	creds := &domain.Credentials{
		UserID:       user.ID,
		PasswordHash: hash,
	}
	if err := s.cred.Save(ctx, creds); err != nil {
		return err
	}
	return nil
}
