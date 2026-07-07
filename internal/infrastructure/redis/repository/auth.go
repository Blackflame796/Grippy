package redis_repository

import (
	entity "Grippy/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{
		client: client,
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *entity.RefreshSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return errors.New("expiration time must be in the future")
	}

	key := "session:" + session.RefreshToken
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*entity.RefreshSession, error) {
	key := "session:" + token
	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var session entity.RefreshSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, token string) error {
	key := "session:" + token
	return r.client.Del(ctx, key).Err()
}
