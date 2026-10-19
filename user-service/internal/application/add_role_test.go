package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func setupRoleTest() (*mockUserRepo, *mockSessionRepo, *application.AddRoleUseCase, *application.RemoveRoleUseCase) {
	users := newMockUserRepo()
	sessions := newMockSessionRepo()
	uow := &mockUoW{}

	addUC := application.NewAddRoleUseCase(users, uow)
	removeUC := application.NewRemoveRoleUseCase(users, sessions, uow)
	return users, sessions, addUC, removeUC
}

func createUser(t *testing.T, users *mockUserRepo, id domain.UserID, roles []domain.Role) {
	t.Helper()
	err := users.Create(context.Background(), &domain.User{
		ID:           id,
		Email:        string(id) + "@local",
		Username:     string(id),
		Roles:        roles,
		RolesVersion: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create user %s: %v", id, err)
	}
}

func TestAddRole_Success(t *testing.T) {
	users, _, uc, _ := setupRoleTest()

	adminID := domain.UserID("admin")
	targetID := domain.UserID("user")
	createUser(t, users, adminID, []domain.Role{domain.RoleAdmin})
	createUser(t, users, targetID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), adminID)

	err := uc.AddRole(ctx, dto.AddRoleInput{
		UserID: targetID,
		Role:   domain.RoleTrader,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := users.GetByID(context.Background(), targetID)
	if !u.HasRole(domain.RoleTrader) {
		t.Fatal("expected role TRADER to be added")
	}
	if u.RolesVersion != 2 {
		t.Fatalf("expected RolesVersion 2, got %d", u.RolesVersion)
	}
}

func TestAddRole_Forbidden(t *testing.T) {
	users, _, uc, _ := setupRoleTest()

	userID := domain.UserID("user")
	createUser(t, users, userID, []domain.Role{domain.RoleUser})
	createUser(t, users, "other", []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), userID)

	err := uc.AddRole(ctx, dto.AddRoleInput{
		UserID: "other",
		Role:   domain.RoleAdmin,
	})
	if !errors.Is(err, er.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAddRole_Idempotent(t *testing.T) {
	users, _, uc, _ := setupRoleTest()

	adminID := domain.UserID("admin")
	targetID := domain.UserID("user")
	createUser(t, users, adminID, []domain.Role{domain.RoleAdmin})
	createUser(t, users, targetID, []domain.Role{domain.RoleUser, domain.RoleTrader})

	ctx := auth.WithUserID(context.Background(), adminID)

	err := uc.AddRole(ctx, dto.AddRoleInput{
		UserID: targetID,
		Role:   domain.RoleTrader,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, _ := users.GetByID(context.Background(), targetID)
	if u.RolesVersion != 1 {
		t.Fatalf("expected RolesVersion 1 (no change), got %d", u.RolesVersion)
	}
}
