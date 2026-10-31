package application

import (
	"context"
	"time"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type ListUsersUseCase struct {
	users outbound.UserRepository
}

func NewListUsersUseCase(users outbound.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{users: users}
}

func (uc *ListUsersUseCase) ListUsers(ctx context.Context, in dto.ListUsersInput) (*dto.ListUsersOutput, error) {
	callerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, er.ErrUnauthenticated
	}

	caller, err := uc.users.GetByID(ctx, callerID)
	if err != nil {
		return nil, err
	}
	if !caller.HasRole(domain.RoleAdmin) {
		return nil, er.ErrForbidden
	}

	limit := int(in.PageSize)
	if limit <= 0 {
		limit = 20
	}

	var afterCreatedAt time.Time
	var afterID domain.UserID
	if in.PageToken != "" {
		afterCreatedAt, afterID, err = decode(in.PageToken)
		if err != nil {
			return nil, er.ErrInvalidPageToken
		}
	}
	users, err := uc.users.List(ctx, limit+1, afterCreatedAt, afterID)
	if err != nil {
		return nil, err
	}

	var nextPageToken string
	if len(users) > limit {
		users = users[:limit]
		last := users[len(users)-1]
		nextPageToken = encode(last.CreatedAt, last.ID)
	}

	return &dto.ListUsersOutput{
		Users:         users,
		NextPageToken: nextPageToken,
	}, nil
}
