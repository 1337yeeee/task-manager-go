// @title Task Manager API
// @version 1.0
// @description API for task manager
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"task-manager/internal/app"
	"task-manager/internal/auth"
	"task-manager/internal/config"
	"task-manager/internal/database"
	"task-manager/internal/server"

	_ "task-manager/docs"
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

	commandExecuted, err := runCommandIfRequested(container)
	if err != nil {
		log.Fatalf("command error: %v", err)
	}
	if commandExecuted {
		return
	}

	srv := server.New(cfg, container)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runCommandIfRequested(container *app.Container) (bool, error) {
	if len(os.Args) < 2 {
		return false, nil
	}

	switch os.Args[1] {
	case "add-admin":
		return true, runAddAdminCommand(container)
	default:
		return false, fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func runAddAdminCommand(container *app.Container) error {
	reader := bufio.NewReader(os.Stdin)

	email, err := promptLine(reader, "Email: ")
	if err != nil {
		return err
	}

	if email == "" {
		return errors.New("email cannot be empty")
	}

	password, err := promptPassword("Password: ")
	if err != nil {
		return err
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	role := auth.UserRoleAdmin
	user, err := container.UserService.Register(
		context.Background(),
		"Administrator",
		email,
		password,
		&role,
	)
	if err != nil {
		return err
	}

	fmt.Printf("admin created: %s (%s)\n", user.Email, user.ID)
	return nil
}

func promptLine(reader *bufio.Reader, label string) (string, error) {
	fmt.Print(label)

	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return strings.TrimSpace(value), nil
}

func promptPassword(label string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return promptLine(reader, label)
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
