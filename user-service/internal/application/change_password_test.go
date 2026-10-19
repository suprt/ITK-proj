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

func setupChangePasswordTest() (*mockCredentialsRepo, *mockSessionRepo, *application.ChangePasswordUseCase) {
	creds := newMockCredentialsRepo()
	sessions := newMockSessionRepo()
	hasher := &mockHasher{}
	uow := &mockUoW{}

	uc := application.NewChangePasswordUseCase(creds, sessions, hasher, uow)
	return creds, sessions, uc
}

func ctxWithUser(id domain.UserID) context.Context {
	return auth.WithUserID(context.Background(), id)
}

func TestChangePassword_Success(t *testing.T) {
	creds, sessions, uc := setupChangePasswordTest()
	userID := domain.UserID("user-1")

	_ = creds.Save(context.Background(), &domain.Credentials{
		UserID:       userID,
		PasswordHash: "old-hash",
	})


	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:  "session-1",
		UserID: userID,
	})

	ctx := ctxWithUser(userID)


	hasher := &mockHasher{hash: "new-hash"}
	uc = application.NewChangePasswordUseCase(creds, sessions, hasher, &mockUoW{})

	err := uc.ChangePassword(ctx, dto.ChangePasswordInput{
		CurrentPassword: "old",
		NewPassword:     "new",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}


	updatedCreds, _ := creds.GetByID(context.Background(), userID)
	if updatedCreds.PasswordHash != "new-hash" {
		t.Fatalf("expected new-hash, got %s", updatedCreds.PasswordHash)
	}

	// Проверяем, что сессии удалены
	if _, err := sessions.GetByToken(context.Background(), "session-1"); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("sessions were not deleted")
	}
}

func TestChangePassword_InvalidCurrentPassword(t *testing.T) {
	creds, sessions, _ := setupChangePasswordTest()
	userID := domain.UserID("user-1")

	_ = creds.Save(context.Background(), &domain.Credentials{
		UserID:       userID,
		PasswordHash: "old-hash",
	})

	hasher := &mockHasher{compareErr: errors.New("mismatch")}
	uc := application.NewChangePasswordUseCase(creds, sessions, hasher, &mockUoW{})

	err := uc.ChangePassword(ctxWithUser(userID), dto.ChangePasswordInput{
		CurrentPassword: "wrong",
		NewPassword:     "new",
	})
	if !errors.Is(err, er.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestChangePassword_NoAuth(t *testing.T) {
	_, _, uc := setupChangePasswordTest()

	err := uc.ChangePassword(context.Background(), dto.ChangePasswordInput{
		CurrentPassword: "old",
		NewPassword:     "new",
	})
	if !errors.Is(err, er.ErrUnauthenticated) {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}
