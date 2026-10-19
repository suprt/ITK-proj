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

func setupAdminUpdateProfileTest() (*mockUserRepo, *application.AdminUpdateProfileUseCase) {
	users := newMockUserRepo()
	uow := &mockUoW{}
	uc := application.NewAdminUpdateProfileUseCase(users, uow)
	return users, uc
}

func TestAdminUpdateProfile_Success(t *testing.T) {
	users, uc := setupAdminUpdateProfileTest()

	adminID := domain.UserID("admin")
	targetID := domain.UserID("user")
	createUser(t, users, adminID, []domain.Role{domain.RoleAdmin})
	createUser(t, users, targetID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), adminID)

	newEmail := "updated@local"
	newUsername := "updated_user"

	err := uc.AdminUpdateProfile(ctx, dto.AdminUpdateProfileInput{
		UserID:   targetID,
		Email:    &newEmail,
		Username: &newUsername,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := users.GetByID(context.Background(), targetID)
	if u.Email != newEmail || u.Username != newUsername {
		t.Fatal("profile was not updated")
	}
}

func TestAdminUpdateProfile_Forbidden(t *testing.T) {
	users, uc := setupAdminUpdateProfileTest()

	userID := domain.UserID("user")
	createUser(t, users, userID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), userID)

	newEmail := "hack@local"
	err := uc.AdminUpdateProfile(ctx, dto.AdminUpdateProfileInput{
		UserID: userID,
		Email:  &newEmail,
	})
	if !errors.Is(err, er.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
