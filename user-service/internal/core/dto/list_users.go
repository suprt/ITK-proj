package dto

import "github.com/suprt/trading/user-service/internal/core/domain"

type ListUsersInput struct {
	PageSize  int32
	PageToken string
}

type ListUsersOutput struct {
	Users         []*domain.User
	NextPageToken string
}
