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

func setupUpdateProfileTest() (*mockUserRepo, *application.UpdateProfileUseCase) {
	users := newMockUserRepo()
	uow := &mockUoW{}
	uc := application.NewUpdateProfileUseCase(users, uow)
	return users, uc
}

func TestUpdateProfile_Success(t *testing.T) {
	users, uc := setupUpdateProfileTest()

	userID := domain.UserID("user")
	createUser(t, users, userID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), userID)

	newUsername := "new_name"
	err := uc.UpdateProfile(ctx, dto.UpdateProfileInput{
		Username: &newUsername,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := users.GetByID(context.Background(), userID)
	if u.Username != newUsername {
		t.Fatal("username was not updated")
	}
}

func TestUpdateProfile_EmailTaken(t *testing.T) {
	users, uc := setupUpdateProfileTest()

	userID := domain.UserID("user")
	existingID := domain.UserID("existing")
	createUser(t, users, userID, []domain.Role{domain.RoleUser})
	createUser(t, users, existingID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), userID)

	takenEmail := string(existingID) + "@local"
	err := uc.UpdateProfile(ctx, dto.UpdateProfileInput{
		Email: &takenEmail,
	})
	if !errors.Is(err, er.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}
