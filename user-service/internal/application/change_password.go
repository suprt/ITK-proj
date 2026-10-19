package application

import (
	"context"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type ChangePasswordUseCase struct {
	credentials outbound.CredentialsRepository
	sessions    outbound.SessionRepository
	hasher      outbound.PasswordHasher
	uow         outbound.UnitOfWork
}

func NewChangePasswordUseCase(
	credentials outbound.CredentialsRepository,
	sessions outbound.SessionRepository,
	hasher outbound.PasswordHasher,
	uow outbound.UnitOfWork,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		credentials: credentials,
		sessions:    sessions,
		hasher:      hasher,
		uow:         uow,
	}
}

func (uc *ChangePasswordUseCase) ChangePassword(ctx context.Context, in dto.ChangePasswordInput) error {
	callerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return er.ErrUnauthenticated
	}

	return uc.uow.Do(ctx, func(ctx context.Context) error {
		creds, err := uc.credentials.GetByID(ctx, callerID)
		if err != nil {
			return err
		}

		if err := uc.hasher.Compare(creds.PasswordHash, in.CurrentPassword); err != nil {
			return er.ErrInvalidCredentials
		}

		newHash, err := uc.hasher.Hash(in.NewPassword)
		if err != nil {
			return err
		}

		if err := uc.credentials.Save(ctx, &domain.Credentials{
			UserID:       callerID,
			PasswordHash: newHash,
		}); err != nil {
			return err
		}

		return uc.sessions.DeleteByUserID(ctx, callerID) //юзеру требуется перелогинится после смены пароля
	})

}
