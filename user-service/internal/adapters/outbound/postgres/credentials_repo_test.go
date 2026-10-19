package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/postgres"
	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func TestCredentialsRepo_SaveAndGet(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	userRepo := postgres.NewUserRepo(pool)
	credsRepo := postgres.NewCredentialsRepo(pool)
	ctx := context.Background()

	userID := newID()
	require.NoError(t, userRepo.Create(ctx, &domain.User{
		ID: userID, Email: "c@local", Username: "c_user", Roles: []domain.Role{domain.RoleUser},
	}))

	// Save
	err := credsRepo.Save(ctx, &domain.Credentials{UserID: userID, PasswordHash: "hash1"})
	require.NoError(t, err)

	// Get
	creds, err := credsRepo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "hash1", creds.PasswordHash)

	// Save (ON CONFLICT UPDATE)
	err = credsRepo.Save(ctx, &domain.Credentials{UserID: userID, PasswordHash: "hash2"})
	require.NoError(t, err)

	creds, err = credsRepo.GetByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "hash2", creds.PasswordHash)
}

func TestCredentialsRepo_NotFound(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	credsRepo := postgres.NewCredentialsRepo(pool)
	ctx := context.Background()

	_, err := credsRepo.GetByID(ctx, "00000000-0000-0000-0000-000000000099")
	require.ErrorIs(t, err, er.ErrUserNotFound)
}
