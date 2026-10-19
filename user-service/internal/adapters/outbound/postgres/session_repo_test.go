package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/postgres"
	er "github.com/suprt/trading/user-service/internal/core/errors"

	"github.com/suprt/trading/user-service/internal/core/domain"
)

func TestSessionRepo_SaveAndGet(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	userRepo := postgres.NewUserRepo(pool)
	sessionRepo := postgres.NewSessionRepo(pool)
	ctx := context.Background()

	userID := newID()
	require.NoError(t, userRepo.Create(ctx, &domain.User{
		ID: userID, Email: "s@local", Username: "s_user", Roles: []domain.Role{domain.RoleUser},
	}))

	token := domain.RefreshToken("test-token")
	expiresAt := time.Now().Add(time.Hour)

	err := sessionRepo.Save(ctx, &domain.RefreshSession{
		Token: token, UserID: userID, ExpiresAt: expiresAt, IssuedAt: time.Now(),
	})
	require.NoError(t, err)

	// Get
	got, err := sessionRepo.GetByToken(ctx, token)
	require.NoError(t, err)
	require.Equal(t, userID, got.UserID)
}

func TestSessionRepo_DeleteByUserID(t *testing.T) {
	pool, cleanup := setupPostgres(t)
	defer cleanup()

	userRepo := postgres.NewUserRepo(pool)
	sessionRepo := postgres.NewSessionRepo(pool)
	ctx := context.Background()

	user1 := newID()
	user2 := newID()
	require.NoError(t, userRepo.Create(ctx, &domain.User{ID: user1, Email: "u1@local", Username: "u1"}))
	require.NoError(t, userRepo.Create(ctx, &domain.User{ID: user2, Email: "u2@local", Username: "u2"}))

	_ = sessionRepo.Save(ctx, &domain.RefreshSession{Token: "t1", UserID: user1, ExpiresAt: time.Now().Add(time.Hour)})
	_ = sessionRepo.Save(ctx, &domain.RefreshSession{Token: "t2", UserID: user2, ExpiresAt: time.Now().Add(time.Hour)})
	_ = sessionRepo.Save(ctx, &domain.RefreshSession{Token: "t3", UserID: user1, ExpiresAt: time.Now().Add(time.Hour)})

	// Удаляем сессии user1
	require.NoError(t, sessionRepo.DeleteByUserID(ctx, user1))

	// Сессии user1 должны быть удалены
	_, err := sessionRepo.GetByToken(ctx, "t1")
	require.ErrorIs(t, err, er.ErrInvalidToken)
	_, err = sessionRepo.GetByToken(ctx, "t3")
	require.ErrorIs(t, err, er.ErrInvalidToken)

	// Сессия user2 должна остаться
	_, err = sessionRepo.GetByToken(ctx, "t2")
	require.NoError(t, err)
}
