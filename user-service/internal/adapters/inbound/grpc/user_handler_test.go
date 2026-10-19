package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	userv1 "github.com/suprt/trading/proto/user/v1"
	grpchand "github.com/suprt/trading/user-service/internal/adapters/inbound/grpc"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/dto"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/inbound/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().
		Register(gomock.Any(), dto.RegisterInput{
			Email:    "test@local",
			Username: "test_user",
			Password: "secret123",
		}).
		Return(&dto.RegisterOutput{UserID: "user-1"}, nil)

	h := grpchand.New(svc)

	resp, err := h.Register(context.Background(), &userv1.RegisterRequest{
		Email:    "test@local",
		Username: "test_user",
		Password: "secret123",
	})
	require.NoError(t, err)
	assert.Equal(t, "user-1", resp.UserId)
}

func TestHandler_Register_EmailTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, er.ErrEmailTaken)

	h := grpchand.New(svc)

	_, err := h.Register(context.Background(), &userv1.RegisterRequest{
		Email:    "taken@local",
		Username: "user",
		Password: "secret123",
	})
	require.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().
		Login(gomock.Any(), dto.LoginInput{
			Email:    "test@local",
			Password: "secret123",
		}).
		Return(&dto.LoginOutput{
			UserID:          "user-1",
			AccessToken:     "access",
			RefreshToken:    "refresh",
			AccessExpiresIn: 900,
		}, nil)

	h := grpchand.New(svc)

	resp, err := h.Login(context.Background(), &userv1.LoginRequest{
		Email:    "test@local",
		Password: "secret123",
	})
	require.NoError(t, err)
	assert.Equal(t, "user-1", resp.UserId)
	assert.Equal(t, "access", resp.AccessToken)
	assert.Equal(t, "refresh", resp.RefreshToken)
	assert.Equal(t, int64(900), resp.AccessExpiresIn)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(nil, er.ErrInvalidCredentials)

	h := grpchand.New(svc)

	_, err := h.Login(context.Background(), &userv1.LoginRequest{
		Email:    "test@local",
		Password: "wrong",
	})
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestHandler_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().Logout(gomock.Any(), dto.LogoutInput{RefreshToken: "token"}).Return(nil)

	h := grpchand.New(svc)
	resp, err := h.Logout(context.Background(), &userv1.LogoutRequest{RefreshToken: "token"})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestHandler_Logout_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(er.ErrForbidden)

	h := grpchand.New(svc)
	_, err := h.Logout(context.Background(), &userv1.LogoutRequest{RefreshToken: "token"})

	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestHandler_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().RefreshToken(gomock.Any(), dto.RefreshTokenInput{RefreshToken: "old"}).Return(&dto.RefreshTokenOutput{
		AccessToken:     "new-access",
		RefreshToken:    "new-refresh",
		AccessExpiresIn: 900,
	}, nil)

	h := grpchand.New(svc)
	resp, err := h.RefreshToken(context.Background(), &userv1.RefreshTokenRequest{RefreshToken: "old"})

	require.NoError(t, err)
	assert.Equal(t, "new-access", resp.AccessToken)
	assert.Equal(t, "new-refresh", resp.RefreshToken)
	assert.Equal(t, int64(900), resp.AccessExpiresIn)
}

func TestHandler_RefreshToken_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().RefreshToken(gomock.Any(), gomock.Any()).Return(nil, er.ErrInvalidToken)

	h := grpchand.New(svc)
	_, err := h.RefreshToken(context.Background(), &userv1.RefreshTokenRequest{RefreshToken: "bad"})

	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestHandler_GetSelfProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := &domain.User{ID: "user-1", Email: "a@b.c", Username: "test", Roles: []domain.Role{domain.RoleUser}}
	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().GetSelfProfile(gomock.Any(), dto.GetSelfProfileInput{}).Return(&dto.GetSelfProfileOutput{User: u}, nil)

	h := grpchand.New(svc)
	resp, err := h.GetSelfProfile(context.Background(), &userv1.GetSelfProfileRequest{})

	require.NoError(t, err)
	assert.Equal(t, "user-1", resp.User.Id)
	assert.Equal(t, "test", resp.User.Username)
}

func TestHandler_GetPublicProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().GetPublicProfile(gomock.Any(), dto.GetPublicProfileInput{UserID: "target"}).Return(&dto.GetPublicProfileOutput{
		UserID: "target", Username: "target_user",
	}, nil)

	h := grpchand.New(svc)
	resp, err := h.GetPublicProfile(context.Background(), &userv1.GetPublicProfileRequest{UserId: "target"})

	require.NoError(t, err)
	assert.Equal(t, "target", resp.Profile.Id)
	assert.Equal(t, "target_user", resp.Profile.Username)
}

func TestHandler_GetPublicProfile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().GetPublicProfile(gomock.Any(), gomock.Any()).Return(nil, er.ErrUserNotFound)

	h := grpchand.New(svc)
	_, err := h.GetPublicProfile(context.Background(), &userv1.GetPublicProfileRequest{UserId: "nonexistent"})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestHandler_UpdateProfile_PassesPointers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, in dto.UpdateProfileInput) error {
			require.NotNil(t, in.Email)
			assert.Equal(t, "new@local", *in.Email)
			assert.Nil(t, in.Username) // не передан
			return nil
		})

	h := grpchand.New(svc)
	email := "new@local"
	_, err := h.UpdateProfile(context.Background(), &userv1.UpdateProfileRequest{Email: &email})

	require.NoError(t, err)
}

func TestHandler_ChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().ChangePassword(gomock.Any(), dto.ChangePasswordInput{
		CurrentPassword: "old", NewPassword: "new12345",
	}).Return(nil)

	h := grpchand.New(svc)
	_, err := h.ChangePassword(context.Background(), &userv1.ChangePasswordRequest{
		CurrentPassword: "old", NewPassword: "new12345",
	})

	require.NoError(t, err)
}

func TestHandler_ChangePassword_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().ChangePassword(gomock.Any(), gomock.Any()).Return(er.ErrInvalidCredentials)

	h := grpchand.New(svc)
	_, err := h.ChangePassword(context.Background(), &userv1.ChangePasswordRequest{
		CurrentPassword: "wrong", NewPassword: "new",
	})

	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestHandler_GetUserByID_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, er.ErrForbidden)

	h := grpchand.New(svc)
	_, err := h.GetUserByID(context.Background(), &userv1.GetUserByIDRequest{UserId: "target"})

	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestHandler_AdminUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().AdminUpdateProfile(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, in dto.AdminUpdateProfileInput) error {
			assert.Equal(t, domain.UserID("target"), in.UserID)
			require.NotNil(t, in.Username)
			assert.Equal(t, "new_name", *in.Username)
			return nil
		})

	h := grpchand.New(svc)
	username := "new_name"
	_, err := h.AdminUpdateProfile(context.Background(), &userv1.AdminUpdateProfileRequest{
		UserId: "target", Username: &username,
	})

	require.NoError(t, err)
}

func TestHandler_AdminResetPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().AdminResetPassword(gomock.Any(), dto.AdminResetPasswordInput{UserID: "target"}).Return(&dto.AdminResetPasswordOutput{
		TemporaryPassword: "temp123",
	}, nil)

	h := grpchand.New(svc)
	resp, err := h.AdminResetPassword(context.Background(), &userv1.AdminResetPasswordRequest{UserId: "target"})

	require.NoError(t, err)
	assert.Equal(t, "temp123", resp.TemporaryPassword)
}

func TestHandler_AddRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().AddRole(gomock.Any(), dto.AddRoleInput{
		UserID: "target", Role: domain.RoleTrader,
	}).Return(nil)

	h := grpchand.New(svc)
	_, err := h.AddRole(context.Background(), &userv1.AddRoleRequest{
		UserId: "target", Role: userv1.Role_ROLE_TRADER,
	})

	require.NoError(t, err)
}

func TestHandler_RemoveRole_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockUserService(ctrl)
	svc.EXPECT().RemoveRole(gomock.Any(), gomock.Any()).Return(er.ErrForbidden)

	h := grpchand.New(svc)
	_, err := h.RemoveRole(context.Background(), &userv1.RemoveRoleRequest{
		UserId: "target", Role: userv1.Role_ROLE_USER,
	})

	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}
