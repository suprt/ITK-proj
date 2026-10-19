package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

func TestRemoveRole_Success(t *testing.T) {
	users, sessions, _, uc := setupRoleTest()

	adminID := domain.UserID("admin")
	targetID := domain.UserID("user")
	createUser(t, users, adminID, []domain.Role{domain.RoleAdmin})
	createUser(t, users, targetID, []domain.Role{domain.RoleUser, domain.RoleTrader})

	_ = sessions.Save(context.Background(), &domain.RefreshSession{
		Token:  "target-session",
		UserID: targetID,
	})

	ctx := auth.WithUserID(context.Background(), adminID)

	err := uc.RemoveRole(ctx, dto.RemoveRoleInput{
		UserID: targetID,
		Role:   domain.RoleTrader,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}


	u, _ := users.GetByID(context.Background(), targetID)
	if u.HasRole(domain.RoleTrader) {
		t.Fatal("expected TRADER to be removed")
	}
	if u.RolesVersion != 2 {
		t.Fatalf("expected RolesVersion 2, got %d", u.RolesVersion)
	}


	if _, err := sessions.GetByToken(context.Background(), "target-session"); !errors.Is(err, er.ErrInvalidToken) {
		t.Fatal("target sessions were not deleted")
	}
}

func TestRemoveRole_Forbidden(t *testing.T) {
	users, _, _, uc := setupRoleTest()

	userID := domain.UserID("user")
	createUser(t, users, userID, []domain.Role{domain.RoleUser, domain.RoleTrader})
	createUser(t, users, "other", []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), userID)

	err := uc.RemoveRole(ctx, dto.RemoveRoleInput{
		UserID: "other",
		Role:   domain.RoleUser,
	})
	if !errors.Is(err, er.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestRemoveRole_Idempotent(t *testing.T) {
	users, _, _, uc := setupRoleTest()

	adminID := domain.UserID("admin")
	targetID := domain.UserID("user")
	createUser(t, users, adminID, []domain.Role{domain.RoleAdmin})
	createUser(t, users, targetID, []domain.Role{domain.RoleUser})

	ctx := auth.WithUserID(context.Background(), adminID)


	err := uc.RemoveRole(ctx, dto.RemoveRoleInput{
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