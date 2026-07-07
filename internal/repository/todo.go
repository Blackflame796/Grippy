package repository

import (
	entity "Grippy/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToDoRepository struct {
	db *pgxpool.Pool
}

func NewToDoRepository(db *pgxpool.Pool) *ToDoRepository {
	return &ToDoRepository{
		db: db,
	}
}

func (r *ToDoRepository) Create(ctx context.Context, draft *entity.ToDoDraft) (*entity.ToDo, error) {
	query := `
			INSERT INTO todo_schema.todos (user_id, title, description, is_completed)
			VALUES ($1, $2, $3, $4)
			RETURNING id, user_id, title, description, is_completed`
	var todo entity.ToDo

	err := r.db.QueryRow(ctx, query, draft.UserID, draft.Title, draft.Description, false).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *ToDoRepository) Update(ctx context.Context, input *entity.ToDo) (*entity.ToDo, error) {
	query := `
			UPDATE todo_schema.todos
			SET title = $1, description = $2, is_completed = $3
			WHERE id = $4, user_id = $5
			RETURNING id, user_id, title, description, is_completed`
	var todo entity.ToDo
	err := r.db.QueryRow(ctx, query, input.Title, input.Description, input.IsCompleted, input.ID, input.UserID).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *ToDoRepository) Delete(ctx context.Context, id uuid.UUID, claims *entity.Claims) (uuid.UUID, error) {
	query := `
			DELETE FROM todo_schema.todos
			WHERE id = $1, user_id = $2`
	_, err := r.db.Exec(ctx, query, id, claims.ID)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *ToDoRepository) GetByID(ctx context.Context, id uuid.UUID, claims *entity.Claims) (*entity.ToDo, error) {
	query := `
			SELECT id, title, description, is_completed
			FROM todo_schema.todos
			WHERE id = $1, user_id = $2`
	var todo entity.ToDo
	err := r.db.QueryRow(ctx, query, id, claims.ID).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *ToDoRepository) GetAll(ctx context.Context, claims *entity.Claims, limit, offset int) ([]*entity.ToDo, error) {
	query := `
			SELECT id, user_id, title, description, is_completed
			FROM todo_schema.todos
			WHERE user_id = $1
			ORDER BY id
			LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, claims.ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.ToDo
	for rows.Next() {
		var todo entity.ToDo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.IsCompleted)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
