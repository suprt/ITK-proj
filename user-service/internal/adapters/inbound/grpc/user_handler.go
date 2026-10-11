package grpc

import (
	"context"

	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/user-service/internal/core/domain"
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
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
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
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
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

func (h *Handler) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	out, err := h.userService.RefreshToken(ctx, dto.RefreshTokenInput{
		RefreshToken: domain.RefreshToken(req.GetRefreshToken()),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.RefreshTokenResponse{
		AccessToken:  out.AccessToken.String(),
		RefreshToken: out.RefreshToken.String(),
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	err := h.userService.Logout(ctx, dto.LogoutInput{
		RefreshToken: domain.RefreshToken(req.GetRefreshToken()),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.LogoutResponse{}, nil
}

func (h *Handler) GetSelfProfile(ctx context.Context, _ *userv1.GetSelfProfileRequest) (*userv1.GetSelfProfileResponse, error) {
	out, err := h.userService.GetSelfProfile(ctx, dto.GetSelfProfileInput{})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.GetSelfProfileResponse{
		User: userToProto(out.User),
	}, nil
}

func (h *Handler) GetUserByID(ctx context.Context, request *userv1.GetUserByIDRequest) (*userv1.GetUserByIDResponse, error) {
	out, err := h.userService.GetUserByID(ctx, dto.GetUserByIDInput{
		UserID: domain.UserID(request.GetUserId()),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &userv1.GetUserByIDResponse{
		User: userToProto(out.User),
	}, nil
}
