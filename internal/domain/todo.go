package entity

import (
	"errors"

	"github.com/google/uuid"
)

type ToDo struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	IsCompleted bool
}

type ToDoDraft struct {
	UserID      uuid.UUID
	Title       string
	Description string
}

func NewTodoDraft(userID uuid.UUID, title, description string) (*ToDoDraft, error) {
	if title == "" {
		return nil, errors.New("Title cannot be empty")
	}
	return &ToDoDraft{
		UserID:      userID,
		Title:       title,
		Description: description,
	}, nil
}
