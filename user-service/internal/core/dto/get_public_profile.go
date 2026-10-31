package dto

import (
	"time"

	"github.com/suprt/trading/user-service/internal/core/domain"
)

type GetPublicProfileInput struct {
	UserID domain.UserID
}

type GetPublicProfileOutput struct {
	UserID    domain.UserID
	Username  string
	CreatedAt time.Time
}
