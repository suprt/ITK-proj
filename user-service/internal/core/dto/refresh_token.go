package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type RefreshTokenInput struct {
	RefreshToken domain.RefreshToken
}

type RefreshTokenOutput struct {
	AccessToken    domain.AccessToken
	RefreshToken   domain.RefreshToken
	AccessExpireIn int64
}
