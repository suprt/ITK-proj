package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type AdminResetPasswordInput struct {
	UserID domain.UserID
}

type AdminResetPasswordOutput struct {
	TemporaryPassword string
}
