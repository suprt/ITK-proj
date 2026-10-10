package inbound

import (
	"context"

	"github.com/suprt/trading/user-service/internal/core/dto"
)

type UserService interface {
	Register(ctx context.Context, in dto.RegisterInput) (*dto.RegisterOutput, error)
	Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error)
}
