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

func setupAdminResetTest() (*mockUserRepo, *mockCredentialsRepo, *mockSessionRepo, *mockHasher, *application.AdminResetPasswordUseCase) {
	users := newMockUserRepo()
	creds := newMockCredentialsRepo()
	sessions := newMockSessionRepo()
	hasher := &mockHasher{hash: "temp-hash"}
	uow := &mockUoW{}

	uc := application.NewAdminResetPasswordUseCase(users, creds, sessions, hasher, uow)
	return users, creds, sessions, hasher, uc
}

func TestAdminResetPassword_Success(t *testing.T) {
	users, creds, sessions, _, uc := setupAdminResetTest()

	adminID := domain.UserID("admin-1")
	targetID := domain.UserID("user-1")

	_ = users.Create(context.Background(), &domain.User{
		ID:       adminID,
		Email:    "admin@local",
		Username: "admin",
		Roles:    []domain.Role{domain.RoleAdmin},
	})

	_ = users.Create(context.Background(), &domain.User{
		ID:       targetID,
		Email:    "user@local",
		Username: "user",
	})

	_ = creds.Save(context.Background(), &domain.Credentials{
		UserID:       targetID,
		PasswordHash: "old-hash",
	})

	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:  "target-session",
		UserID: targetID,
	})

	ctx := auth.WithUserID(context.Background(), adminID)

	out, err := uc.AdminResetPassword(ctx, dto.AdminResetPasswordInput{
		UserID: targetID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.TemporaryPassword == "" {
		t.Fatal("expected non-empty temporary password")
	}

	updatedCreds, _ := creds.GetByID(context.Background(), targetID)
	if updatedCreds.PasswordHash != "temp-hash" {
		t.Fatalf("expected temp-hash, got %s", updatedCreds.PasswordHash)
	}

	if _, err := sessions.GetByToken(context.Background(), "target-session"); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("target sessions were not deleted")
	}
}

func TestAdminResetPassword_Forbidden(t *testing.T) {
	users, _, _, _, uc := setupAdminResetTest()

	userID := domain.UserID("user-1")

	_ = users.Create(context.Background(), &domain.User{
		ID:    userID,
		Email: "user@local",
		Roles: []domain.Role{domain.RoleUser},
	})

	ctx := auth.WithUserID(context.Background(), userID)

	_, err := uc.AdminResetPassword(ctx, dto.AdminResetPasswordInput{
		UserID: userID,
	})
	if !errors.Is(err, er.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAdminResetPassword_TargetNotFound(t *testing.T) {
	users, _, _, _, uc := setupAdminResetTest()

	adminID := domain.UserID("admin-1")

	_ = users.Create(context.Background(), &domain.User{
		ID:    adminID,
		Email: "admin@local",
		Roles: []domain.Role{domain.RoleAdmin},
	})

	ctx := auth.WithUserID(context.Background(), adminID)

	_, err := uc.AdminResetPassword(ctx, dto.AdminResetPasswordInput{
		UserID: "nonexistent",
	})
	if !errors.Is(err, er.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
