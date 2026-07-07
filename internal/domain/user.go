package entity

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Avatar       string
	Email        string
	Username     string
	PasswordHash string
}

type UserInfo struct {
	Avatar       string
	Email        string
	Username     string
	PasswordHash string
}
