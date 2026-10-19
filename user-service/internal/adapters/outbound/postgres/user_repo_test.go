package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/postgres"
	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func TestUserRepo_CreateAndByEmail(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		ID:           newID(),
		Email:        "test@local",
		Username:     "test_user",
		Roles:        []domain.Role{domain.RoleUser, domain.RoleTrader},
		RolesVersion: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Читаем обратно
	got, err := repo.GetByEmail(ctx, "test@local")
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.Username, got.Username)
	require.ElementsMatch(t, user.Roles, got.Roles)
}

func TestUserRepo_Create_EmailTaken(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{ID: newID(), Email: "same@local", Username: "u1", Roles: []domain.Role{domain.RoleUser}}
	require.NoError(t, repo.Create(ctx, user))

	err := repo.Create(ctx, &domain.User{ID: newID(), Email: "same@local", Username: "u2", Roles: []domain.Role{domain.RoleUser}})
	require.ErrorIs(t, err, er.ErrEmailTaken)
}

func TestUserRepo_Update(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{ID: newID(), Email: "old@local", Username: "old_user", Roles: []domain.Role{domain.RoleUser}, RolesVersion: 1}
	require.NoError(t, repo.Create(ctx, user))

	user.Email = "new@local"
	user.Username = "new_user"
	user.Roles = []domain.Role{domain.RoleUser, domain.RoleAdmin}
	require.NoError(t, repo.Update(ctx, user))

	got, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "new@local", got.Email)
	require.Equal(t, "new_user", got.Username)
	require.ElementsMatch(t, []domain.Role{domain.RoleUser, domain.RoleAdmin}, got.Roles)
}
