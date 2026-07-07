package todo_usecase

import (
	entity "Grippy/internal/domain"
	"Grippy/internal/repository"
	"context"

	"github.com/google/uuid"
)

type ToDoUseCase struct {
	repo *repository.ToDoRepository
}

func NewToDoUseCase(r *repository.ToDoRepository) *ToDoUseCase {
	return &ToDoUseCase{repo: r}
}

func (uc *ToDoUseCase) Create(ctx context.Context, userID uuid.UUID, input CreateToDoRequest) (*ToDoResponse, error) {
	todoDraft, err := entity.NewTodoDraft(userID, input.Title, input.Description)
	if err != nil {
		return nil, err
	}

	todo, err := uc.repo.Create(ctx, todoDraft)
	if err != nil {
		return nil, err
	}

	return &ToDoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		UserID:      todo.UserID,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}, nil
}

func (uc *ToDoUseCase) Update(ctx context.Context, input UpdateToDoRequest, claims *entity.Claims) (*ToDoResponse, error) {
	inputToDo := &entity.ToDo{
		ID:          input.ID,
		UserID:      claims.ID,
		Title:       input.Title,
		Description: input.Description,
		IsCompleted: input.IsCompleted,
	}

	todo, err := uc.repo.Update(ctx, inputToDo)
	if err != nil {
		return nil, err
	}

	return &ToDoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		UserID:      todo.UserID,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}, nil
}

func (uc *ToDoUseCase) Delete(ctx context.Context, todoID uuid.UUID, claims *entity.Claims) (uuid.UUID, error) {
	id, err := uc.repo.Delete(ctx, todoID, claims)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (uc *ToDoUseCase) GetToDoByID(ctx context.Context, todoID uuid.UUID, claims *entity.Claims) (*ToDoResponse, error) {
	todo, err := uc.repo.GetByID(ctx, todoID, claims)
	if err != nil {
		return nil, err
	}

	return &ToDoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		UserID:      todo.UserID,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}, nil
}

func (uc *ToDoUseCase) GetAllToDo(ctx context.Context, claims *entity.Claims, limit, offset int) ([]*ToDoResponse, error) {
	todos, err := uc.repo.GetAll(ctx, claims, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*ToDoResponse, 0, len(todos))
	for _, todo := range todos {
		responses = append(responses, &ToDoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			UserID:      todo.UserID,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
		})
	}

	return responses, nil
}
