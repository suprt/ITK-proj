package grpc

import (
	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func roleToProto(r domain.Role) userv1.Role {
	switch r {
	case domain.RoleUser:
		return userv1.Role_ROLE_USER
	case domain.RoleTrader:
		return userv1.Role_ROLE_TRADER
	case domain.RoleAdmin:
		return userv1.Role_ROLE_ADMIN
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}

func userToProto(u *domain.User) *userv1.User {
	roles := make([]userv1.Role, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = roleToProto(role)
	}
	return &userv1.User{
		Id:        string(u.ID),
		Email:     u.Email,
		Roles:     roles,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
