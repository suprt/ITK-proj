package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type LogoutInput struct {
	RefreshToken domain.RefreshToken
}
