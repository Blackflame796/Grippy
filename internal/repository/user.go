package repository

import (
	entity "Grippy/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, username, password_hash FROM user_schema.users WHERE email = $1 LIMIT 1`

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, input *entity.UserInfo) (*entity.User, error) {
	query := `
			INSERT INTO user_schema.users (email, username, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, email, username`
	var user entity.User

	err := r.db.QueryRow(ctx, query, input.Email, input.Username, input.PasswordHash).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserInfo(ctx context.Context, input *entity.User) (*entity.User, error) {
	query := `
			UPDATE user_schema.users
			SET email = $1, username = $2
			WHERE id = $3
			RETURNING id, email, username
	`

	var user entity.User

	err := r.db.QueryRow(ctx, query, input.Email, input.Username, input.ID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT id, email, username, password_hash FROM user_schema.users WHERE username = $1 LIMIT 1`

	var user entity.User
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*entity.User, string, error) {
	query := `
			WITH old_user AS (
				SELECT avatar FROM user_schema.users WHERE id = $2
			)
			UPDATE user_schema.users
			SET avatar = $1
			WHERE id = $2
			RETURNING id, avatar, email, username, (SELECT avatar FROM old_user) AS old_avatar
		`

	var user entity.User
	var oldAvatar *string

	err := r.db.QueryRow(ctx, query, avatarURL, userID).Scan(
		&user.ID,
		&user.Avatar,
		&user.Email,
		&user.Username,
		&oldAvatar,
	)
	if err != nil {
		return nil, "", err
	}

	var oldAvatarStr string
	if oldAvatar != nil {
		oldAvatarStr = *oldAvatar
	}

	return &user, oldAvatarStr, nil
}
