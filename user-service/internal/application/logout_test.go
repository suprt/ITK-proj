package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func setupLogoutTest() (*mockSessionRepo, *application.LogoutUseCase) {
	sessions := newMockSessionRepo()
	uow := &mockUoW{}
	uc := application.NewLogoutUseCase(sessions, uow)
	return sessions, uc
}

func TestLogout_Success(t *testing.T) {
	sessions, uc := setupLogoutTest()

	userID := domain.UserID("user")
	token := domain.RefreshToken("valid-token")

	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:  token,
		UserID: userID,
	})

	ctx := auth.WithUserID(context.Background(), userID)

	err := uc.Logout(ctx, dto.LogoutInput{RefreshToken: token})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := sessions.GetByToken(context.Background(), token); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("session was not deleted")
	}
}

func TestLogout_Forbidden(t *testing.T) {
	sessions, uc := setupLogoutTest()

	ownerID := domain.UserID("owner")
	attackerID := domain.UserID("attacker")
	token := domain.RefreshToken("valid-token")

	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:  token,
		UserID: ownerID,
	})

	ctx := auth.WithUserID(context.Background(), attackerID)

	err := uc.Logout(ctx, dto.LogoutInput{RefreshToken: token})
	if !errors.Is(err, er.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
