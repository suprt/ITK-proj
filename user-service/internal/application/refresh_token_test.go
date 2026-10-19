package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func setupRefreshTest() (*mockUserRepo, *mockSessionRepo, *mockTokenIssuer, *application.RefreshTokenUseCase) {
	users := newMockUserRepo()
	sessions := newMockSessionRepo()
	tokens := &mockTokenIssuer{}
	uow := &mockUoW{}

	uc := application.NewRefreshTokenUseCase(users, sessions, tokens, uow, application.RefreshTokenConfig{
		RefreshTTL: 168 * time.Hour,
		AccessTTL:  15 * time.Minute,
	})

	return users, sessions, tokens, uc
}

func TestRefreshToken_Success(t *testing.T) {
	users, sessions, tokens, uc := setupRefreshTest()

	userID := domain.UserID("user-1")
	_ = users.Create(context.Background(), &domain.User{ID: userID, Email: "test@local"})

	oldToken := domain.RefreshToken("old-refresh")
	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:     oldToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour),
	})

	tokens.access = "new-access"
	tokens.refresh = "new-refresh"

	out, err := uc.RefreshToken(context.Background(), dto.RefreshTokenInput{
		RefreshToken: oldToken,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.AccessToken != "new-access" {
		t.Fatalf("expected new-access, got %s", out.AccessToken)
	}
	if out.RefreshToken != "new-refresh" {
		t.Fatalf("expected new-refresh, got %s", out.RefreshToken)
	}

	if _, err := sessions.GetByToken(context.Background(), oldToken); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("old session was not deleted")
	}

	if _, err := sessions.GetByToken(context.Background(), "new-refresh"); err != nil {
		t.Fatalf("new session not found: %v", err)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	_, _, _, uc := setupRefreshTest()

	_, err := uc.RefreshToken(context.Background(), dto.RefreshTokenInput{
		RefreshToken: "nonexistent",
	})
	if !errors.Is(err, er.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestRefreshToken_Expired(t *testing.T) {
	users, sessions, _, uc := setupRefreshTest()

	userID := domain.UserID("user-1")
	_ = users.Create(context.Background(), &domain.User{ID: userID, Email: "test@local"})

	oldToken := domain.RefreshToken("expired-token")
	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:     oldToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(-time.Hour), // протух
	})

	_, err := uc.RefreshToken(context.Background(), dto.RefreshTokenInput{
		RefreshToken: oldToken,
	})
	if !errors.Is(err, er.ErrExpiredToken) {
		t.Fatalf("expected ErrExpiredToken, got %v", err)
	}

	if _, err := sessions.GetByToken(context.Background(), oldToken); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("expired session was not deleted")
	}
}

func TestRefreshToken_UserDeleted(t *testing.T) {

	_, sessions, _, uc := setupRefreshTest()

	oldToken := domain.RefreshToken("orphan-token")
	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:     oldToken,
		UserID:    "deleted-user",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	_, err := uc.RefreshToken(context.Background(), dto.RefreshTokenInput{
		RefreshToken: oldToken,
	})

	if !errors.Is(err, er.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	// Сессия должна быть удалена
	if _, err := sessions.GetByToken(context.Background(), oldToken); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("orphan session was not deleted")
	}
}
