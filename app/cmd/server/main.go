package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"

	"task-manager/internal/app"
	"task-manager/internal/config"
	"task-manager/internal/database"
	"task-manager/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := WaitForDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	redisClient, err := WaitForRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	container := app.NewContainer(cfg, db, redisClient)

	srv := server.New(cfg, container)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func WaitForDB(cfg config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = database.NewPostgres(cfg)
		if err == nil {
			return db, nil
		}

		log.Println("waiting for database...")
		time.Sleep(2 * time.Second)
	}

	return nil, err
}

func WaitForRedis(cfg config.Config) (*redis.Client, error) {
	var client *redis.Client
	var err error
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		client, err = database.NewRedis(ctx, cfg)
		if err == nil {
			return client, nil
		}

		log.Println("waiting for redis...")
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
