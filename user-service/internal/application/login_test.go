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

func TestLogin_Success(t *testing.T) {
	users := newMockUserRepo()
	creds := newMockCredentialsRepo()
	hasher := &mockHasher{hash: "hashed", compareErr: nil}
	sessions := newMockSessionRepo()
	tokens := &mockTokenIssuer{}
	uow := &mockUoW{}

	user := &domain.User{ID: "user-1", Email: "test@local", Username: "test_user", Roles: []domain.Role{domain.RoleUser}}
	_ = users.Create(context.Background(), user)
	_ = creds.Save(context.Background(), &domain.Credentials{UserID: "user-1", PasswordHash: "hashed"})

	uc := application.NewLoginUseCase(users, creds, sessions, hasher, tokens, uow, application.LoginConfig{
		RefreshTTL: 168 * time.Hour,
		AccessTTL:  15 * time.Minute,
	})

	out, err := uc.Login(context.Background(), dto.LoginInput{
		Email:    "test@local",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if out.RefreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}

	// Проверяем, что сессия сохранена
	if _, err := sessions.GetByToken(context.Background(), out.RefreshToken); err != nil {
		t.Fatalf("session not saved: %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	uc := application.NewLoginUseCase(newMockUserRepo(), newMockCredentialsRepo(), newMockSessionRepo(),
		&mockHasher{}, &mockTokenIssuer{}, &mockUoW{}, application.LoginConfig{
			RefreshTTL: time.Hour,
			AccessTTL:  time.Minute,
		})

	_, err := uc.Login(context.Background(), dto.LoginInput{
		Email:    "nonexistent@local",
		Password: "secret123",
	})
	if !errors.Is(err, er.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	users := newMockUserRepo()
	creds := newMockCredentialsRepo()
	hasher := &mockHasher{compareErr: errors.New("mismatch")}

	_ = users.Create(context.Background(), &domain.User{ID: "u1", Email: "test@local", Username: "test"})
	_ = creds.Save(context.Background(), &domain.Credentials{UserID: "u1", PasswordHash: "hashed"})

	uc := application.NewLoginUseCase(users, creds, newMockSessionRepo(), hasher, &mockTokenIssuer{}, &mockUoW{}, application.LoginConfig{
		RefreshTTL: time.Hour,
		AccessTTL:  time.Minute,
	})

	_, err := uc.Login(context.Background(), dto.LoginInput{
		Email:    "test@local",
		Password: "wrong",
	})
	if !errors.Is(err, er.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
