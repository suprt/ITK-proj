package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func TestRegister_Success(t *testing.T) {
	users := newMockUserRepo()
	creds := newMockCredentialsRepo()
	hasher := &mockHasher{hash: "hashed"}
	ids := &mockIDGen{id: "user-1"}
	uow := &mockUoW{}

	uc := application.NewRegisterUseCase(users, creds, hasher, ids, uow)

	out, err := uc.Register(context.Background(), dto.RegisterInput{
		Email:    "test@local",
		Username: "test_user",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserID != "user-1" {
		t.Fatalf("expected user-1, got %s", out.UserID)
	}

	// Проверяем, что юзер в репозитории
	u, err := users.GetByEmail(context.Background(), "test@local")
	if err != nil {
		t.Fatalf("user not found: %v", err)
	}
	if u.Username != "test_user" {
		t.Fatalf("expected test_user, got %s", u.Username)
	}

	if _, err := creds.GetByID(context.Background(), "user-1"); err != nil {
		t.Fatalf("credentials not saved: %v", err)
	}
}

func TestRegister_EmailTaken(t *testing.T) {
	users := newMockUserRepo()
	_ = users.Create(context.Background(), &domain.User{
		ID:       "existing",
		Email:    "taken@local",
		Username: "existing_user",
	})

	uc := application.NewRegisterUseCase(users, newMockCredentialsRepo(), &mockHasher{}, &mockIDGen{id: "new"}, &mockUoW{})

	_, err := uc.Register(context.Background(), dto.RegisterInput{
		Email:    "taken@local",
		Username: "new_user",
		Password: "secret123",
	})
	if !errors.Is(err, er.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestRegister_UsernameTaken(t *testing.T) {
	users := newMockUserRepo()
	_ = users.Create(context.Background(), &domain.User{
		ID:       "existing",
		Email:    "existing@local",
		Username: "taken_user",
	})

	uc := application.NewRegisterUseCase(users, newMockCredentialsRepo(), &mockHasher{}, &mockIDGen{id: "new"}, &mockUoW{})

	_, err := uc.Register(context.Background(), dto.RegisterInput{
		Email:    "new@local",
		Username: "taken_user",
		Password: "secret123",
	})
	if !errors.Is(err, er.ErrUsernameTaken) {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}
