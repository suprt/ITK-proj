package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type AdminUpdateProfileInput struct {
	UserID   domain.UserID
	Email    *string
	Username *string
}
