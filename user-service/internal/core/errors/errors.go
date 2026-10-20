package errors

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid refresh token")
	ErrExpiredToken       = errors.New("expired refresh token")
	ErrWeakPassword       = errors.New("weak password")
	ErrForbidden          = errors.New("forbidden")
	ErrDuplicateID        = errors.New("duplicate ID")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrInvalidRole        = errors.New("invalid role")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrCacheMiss          = errors.New("cache miss")
)
