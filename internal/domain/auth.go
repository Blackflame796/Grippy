package entity

import (
	"time"

	"github.com/google/uuid"
)

type Claims struct {
	ID        uuid.UUID
	Username  string
	ExpiresAt time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshSession struct {
	ID           string
	UserID       uuid.UUID
	RefreshToken string
	Fingerprint  string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func (s *RefreshSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
