package grpc

import (
	"context"

	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/user-service/internal/core/dto"
	"github.com/suprt/trading/user-service/internal/core/ports/inbound"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	userService inbound.UserService
}

func New(userService inbound.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	out, err := h.userService.Register(ctx, dto.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.RegisterResponse{
		UserId:       string(out.UserID),
		AccessToken:  out.AccessToken.String(),
		RefreshToken: out.RefreshToken.String(),
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	out, err := h.userService.Login(ctx, dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.LoginResponse{
		UserId:       string(out.UserID),
		AccessToken:  out.AccessToken.String(),
		RefreshToken: out.RefreshToken.String(),
	}, nil
}
