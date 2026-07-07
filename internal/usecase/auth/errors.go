package auth_usecase

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("Username already taken")
	ErrUserAlreadyExists     = errors.New("User with this email already exists")
	ErrInvalidCredentials    = errors.New("Invalid email or password")
	ErrInvalidToken          = errors.New("Invalid or expired token")
	ErrSessionNotFound       = errors.New("Refresh session not found")
)
