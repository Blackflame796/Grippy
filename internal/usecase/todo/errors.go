package todo_usecase

import "errors"

var (
	ErrTodoNotFound   = errors.New("Todo not found")
	ErrForbidden      = errors.New("You cannot modify this todo")
	ErrInvalidInput   = errors.New("Invalid input data")
	ErrDuplicateTitle = errors.New("A todo with this title already exists")
)
