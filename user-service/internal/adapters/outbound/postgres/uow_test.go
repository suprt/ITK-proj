package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/postgres"
	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func TestUoW_RollbackOnError(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	uow := postgres.NewUoW(pool)
	userRepo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	expectedErr := errors.New("boom")

	err := uow.Do(ctx, func(ctx context.Context) error {

		err := userRepo.Create(ctx, &domain.User{
			ID: newID(), Email: "rollback@local", Username: "rollback_user", Roles: []domain.Role{domain.RoleUser},
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		})
		require.NoError(t, err)

		return expectedErr
	})
	require.ErrorIs(t, err, expectedErr)

	_, err = userRepo.GetByEmail(ctx, "rollback@local")
	require.ErrorIs(t, err, er.ErrUserNotFound)
}
