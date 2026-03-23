package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"task-manager/internal/config"
)

func NewRedis(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	r := cfg.Redis

	client := redis.NewClient(&redis.Options{
		Addr:         r.Addr,
		Username:     r.User,
		Password:     r.Password,
		DB:           r.DB,
		MaxRetries:   r.MaxRetries,
		DialTimeout:  r.DialTimeout,
		ReadTimeout:  r.Timeout,
		WriteTimeout: r.Timeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return client, nil
}
