package repository

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthRepository interface {
	Store(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetByUserID(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, userID string) error
}

type authRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) AuthRepository {
	return &authRepository{redis: redis}
}

func (r *authRepository) Store(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.redis.Set(ctx, userID, token, ttl).Err()
}

func (r *authRepository) GetByUserID(ctx context.Context, userID string) (string, error) {
	val, err := r.redis.Get(ctx, userID).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *authRepository) Delete(ctx context.Context, userID string) error {
	return r.redis.Del(ctx, userID).Err()
}
