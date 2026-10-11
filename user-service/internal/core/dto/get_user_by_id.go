package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type GetUserByIDInput struct {
	UserID domain.UserID
}

type GetUserByIDOutput struct {
	User *domain.User
}
