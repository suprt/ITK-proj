package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID       domain.UserID
	AccessToken  domain.AccessToken
	RefreshToken domain.RefreshToken
}
