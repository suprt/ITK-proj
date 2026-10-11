package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type RemoveRoleInput struct {
	UserID domain.UserID
	Role   domain.Role
}
