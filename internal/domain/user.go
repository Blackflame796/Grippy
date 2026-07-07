package entity

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
}

type UserInfo struct {
	Email        string
	Username     string
	PasswordHash string
}
