package inbound

import (
	"context"

	"github.com/suprt/trading/user-service/internal/core/dto"
)

type UserService interface {
	Register(ctx context.Context, in dto.RegisterInput) (*dto.RegisterOutput, error)
	Login(ctx context.Context, in dto.LoginInput) (*dto.LoginOutput, error)
	RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (*dto.RefreshTokenOutput, error)
	Logout(ctx context.Context, in dto.LogoutInput) error
	GetSelfProfile(ctx context.Context, in dto.GetSelfProfileInput) (*dto.GetSelfProfileOutput, error)
	GetUserByID(ctx context.Context, in dto.GetUserByIDInput) (*dto.GetUserByIDOutput, error)
	AddRole(ctx context.Context, in dto.AddRoleInput) error
	RemoveRole(ctx context.Context, in dto.RemoveRoleInput) error
	GetPublicProfile(ctx context.Context, in dto.GetPublicProfileInput) (*dto.GetPublicProfileOutput, error)
}
