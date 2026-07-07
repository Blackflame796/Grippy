package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToDo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
}

type ToDoRepository struct {
	db *pgxpool.Pool
}

func NewToDoRepository(db *pgxpool.Pool) *ToDoRepository {
	return &ToDoRepository{
		db: db,
	}
}

func (r *ToDoRepository) Add(ctx context.Context, title string, description string) (*ToDo, error) {
	query := `
			INSERT INTO todo_schema.todos (title, description, is_completed)
			VALUES ($1, $2, $3)
			RETURNING id, title, description, is_completed`
	var todo ToDo
	err := r.db.QueryRow(ctx, query, title, description, false).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *ToDoRepository) Update(ctx context.Context, id uuid.UUID, title string, description string, isCompleted bool) (*ToDo, error) {
	query := `
			UPDATE todo_schema.todos
			SET title = $1, description = $2, is_completed = $3
			WHERE id = $4
			RETURNING id, title, description, is_completed`
	var todo ToDo
	err := r.db.QueryRow(ctx, query, title, description, isCompleted, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *ToDoRepository) Delete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	query := `
			DELETE FROM todo_schema.todos
			WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *ToDoRepository) GetByID(ctx context.Context, id uuid.UUID) (*ToDo, error) {
	query := `
			SELECT id, title, description, is_completed
			FROM todo_schema.todos
			WHERE id = $1`
	var todo ToDo
	err := r.db.QueryRow(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *ToDoRepository) GetAll(ctx context.Context) ([]*ToDo, error) {
	query := `
			SELECT id, title, description, is_completed
			FROM todo_schema.todos`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []*ToDo
	for rows.Next() {
		var todo ToDo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}
