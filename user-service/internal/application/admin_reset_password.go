package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type AdminResetPasswordUseCase struct {
	users   outbound.UserRepository
	cred    outbound.CredentialsRepository
	session outbound.SessionRepository
	hasher  outbound.PasswordHasher
	uow     outbound.UnitOfWork
}

func NewAdminResetPasswordUseCase(
	users outbound.UserRepository,
	cred outbound.CredentialsRepository,
	session outbound.SessionRepository,
	hasher outbound.PasswordHasher,
	uow outbound.UnitOfWork) *AdminResetPasswordUseCase {
	return &AdminResetPasswordUseCase{
		users:   users,
		cred:    cred,
		session: session,
		hasher:  hasher,
		uow:     uow,
	}
}

const tempPasswordLength = 16

func (uc *AdminResetPasswordUseCase) AdminResetPassword(ctx context.Context,
	in dto.AdminResetPasswordInput) (*dto.AdminResetPasswordOutput, error) {
	callerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, er.ErrUnauthenticated
	}

	tempPass, err := generateTempPass(tempPasswordLength)
	if err != nil {
		return nil, err
	}

	err = uc.uow.Do(ctx, func(ctx context.Context) error {
		caller, err := uc.users.GetByID(ctx, callerID)
		if err != nil {
			if errors.Is(err, er.ErrUserNotFound) {
				return er.ErrUnauthenticated
			}
			return err
		}
		if !caller.HasRole(domain.RoleAdmin) {
			return er.ErrForbidden
		}

		hash, err := uc.hasher.Hash(tempPass)
		if err != nil {
			return err
		}

		if err := uc.cred.Save(ctx, &domain.Credentials{
			UserID:       in.UserID,
			PasswordHash: hash,
		}); err != nil {
			return err
		}

		return uc.session.DeleteByUserID(ctx, in.UserID)
	})
	if err != nil {
		return nil, err
	}

	return &dto.AdminResetPasswordOutput{TemporaryPassword: tempPass}, nil

}

func generateTempPass(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
