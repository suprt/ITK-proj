package application

import (
	"context"
	"time"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type UpdateProfileUseCase struct {
	users outbound.UserRepository
	uow   outbound.UnitOfWork
}

func NewUpdateProfileUseCase(
	users outbound.UserRepository,
	uow outbound.UnitOfWork,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		users: users,
		uow:   uow,
	}
}

func (uc *UpdateProfileUseCase) UpdateProfile(ctx context.Context, in dto.UpdateProfileInput) error {
	callerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return er.ErrUnauthenticated
	}

	return uc.uow.Do(ctx, func(ctx context.Context) error {
		user, err := uc.users.GetByID(ctx, callerID)
		if err != nil {
			return err
		}

		if in.Email != nil {
			user.Email = *in.Email
		}
		if in.Username != nil {
			user.Username = *in.Username
		}
		user.UpdatedAt = time.Now()

		return uc.users.Update(ctx, user)
	})
}
