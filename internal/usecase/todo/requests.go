package todo_usecase

import (
	"github.com/google/uuid"
)

type CreateToDoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateToDoRequest struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
}

type DeleteToDoRequest struct {
	ID uuid.UUID `json:"id"`
}
