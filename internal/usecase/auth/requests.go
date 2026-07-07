package auth_usecase

import (
	"github.com/google/uuid"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	ID uuid.UUID `json:"id"`
}
