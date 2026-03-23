package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"task-manager/internal/config"
	"task-manager/internal/domain/repository"
	"task-manager/internal/handler"
	"task-manager/internal/service"
	"task-manager/internal/utils"
)

type Container struct {
	DB *gorm.DB

	TokenManager utils.TokenManager

	UserRepository    repository.UserRepository
	ProjectRepository repository.ProjectRepository
	TaskRepository    repository.TaskRepository
	AuthRepository    repository.AuthRepository

	AuthService    service.AuthService
	UserService    service.UserService
	ProjectService service.ProjectService
	TaskService    service.TaskService

	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	ProjectHandler *handler.ProjectHandler
	TaskHandler    *handler.TaskHandler
}

func NewContainer(cfg config.Config, db *gorm.DB, redis *redis.Client) *Container {
	c := &Container{DB: db}

	c.TokenManager = utils.NewTokenManager(
		cfg.JWTSecret,
		utils.DefaultAccessTTL,
		utils.DefaultRefreshTTL,
	)

	// repositories
	c.UserRepository = repository.NewUserRepository(db)
	c.ProjectRepository = repository.NewProjectRepository(db)
	c.TaskRepository = repository.NewTaskRepository(db)
	c.AuthRepository = repository.NewAuthRepository(redis)

	// services
	c.AuthService = service.NewAuthService(c.UserRepository, c.AuthRepository, c.TokenManager)
	c.UserService = service.NewUserService(c.UserRepository)
	c.ProjectService = service.NewProjectService(c.ProjectRepository)
	c.TaskService = service.NewTaskService(c.TaskRepository, c.ProjectService, c.UserService)

	// handlers
	c.AuthHandler = handler.NewAuthHandler(c.AuthService)
	c.UserHandler = handler.NewUserHandler(c.UserService)
	c.ProjectHandler = handler.NewProjectHandler(c.ProjectService)
	c.TaskHandler = handler.NewTaskHandler(c.TaskService)

	return c
}
