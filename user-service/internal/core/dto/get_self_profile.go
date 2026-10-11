package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type GetSelfProfileInput struct{}

type GetSelfProfileOutput struct {
	User *domain.User
}
