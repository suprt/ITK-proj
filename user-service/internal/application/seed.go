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
	Username string
	Password string
	Roles    []domain.Role
}

type Seed struct {
	user   outbound.UserRepository
	cred   outbound.CredentialsRepository
	hasher outbound.PasswordHasher
	ids    outbound.IDGenerator
	uow    outbound.UnitOfWork
}

func NewSeed(
	user outbound.UserRepository,
	cred outbound.CredentialsRepository,
	hasher outbound.PasswordHasher,
	ids outbound.IDGenerator,
	uow outbound.UnitOfWork,
) *Seed {
	return &Seed{
		user:   user,
		cred:   cred,
		hasher: hasher,
		ids:    ids,
		uow:    uow,
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
		ID:           domain.UserID(s.ids.New()),
		Email:        in.Email,
		Username:     in.Username,
		Roles:        in.Roles,
		RolesVersion: 1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	creds := &domain.Credentials{
		UserID:       user.ID,
		PasswordHash: hash,
	}
	err = s.uow.Do(ctx, func(ctx context.Context) error {
		if err := s.user.Create(ctx, user); err != nil {
			return err
		}
		return s.cred.Save(ctx, creds)
	})
	return err
}
