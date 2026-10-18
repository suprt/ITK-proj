package domain

import "slices"

import "time"

type UserID string

type Role string

const (
	RoleUser   Role = "USER"
	RoleTrader Role = "TRADER"
	RoleAdmin  Role = "ADMIN"
)

type User struct {
	ID           UserID
	Email        string
	Username     string
	Roles        []Role
	RolesVersion int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) HasRole(role Role) bool {
	return slices.Contains(u.Roles, role)
}
