package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	UserID         domain.UserID
	AccessToken    domain.AccessToken
	RefreshToken    domain.RefreshToken
	AccessExpiresIn int64
}
