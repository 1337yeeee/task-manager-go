package routes

import (
	"github.com/gin-gonic/gin"

	"task-manager/internal/app"
	"task-manager/internal/middleware"
)

func SetupRouter(container *app.Container) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	registerAuthRoutes(api, container)
	registerProtectedRoutes(api, container)

	return r
}

func registerAuthRoutes(api *gin.RouterGroup, container *app.Container) {
	authHandler := container.AuthHandler

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}
}

func registerProtectedRoutes(api *gin.RouterGroup, container *app.Container) {
	userHandler := container.UserHandler
	projectHandler := container.ProjectHandler
	taskHandler := container.TaskHandler

	protected := api.Group("")
	protected.Use(middleware.JWTAuthMiddleware(container.TokenManager))

	// Users
	users := protected.Group("/users")
	users.Use(middleware.RequireRole("admin"))
	{
		users.GET("/", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.POST("/", userHandler.Register)
	}

	// Projects
	projects := protected.Group("/projects")
	{
		projects.GET("/", projectHandler.GetAll)
		projects.GET("/:id", projectHandler.GetByID)
		projects.GET("/:id/tasks", taskHandler.GetByProject)
	}

	projectsModerators := protected.Group("/projects")
	projectsModerators.Use(middleware.RequireRolesModerators())
	{
		projectsModerators.POST("/", projectHandler.Create)
		projectsModerators.PUT("/:id", projectHandler.Update)
		projectsModerators.DELETE("/:id", projectHandler.Delete)
		projectsModerators.POST("/:id/tasks", taskHandler.Create)
	}

	// Tasks
	tasks := protected.Group("/tasks")
	{
		tasks.GET("/:id", taskHandler.GetByID)
	}

	tasksModerators := tasks.Group("")
	tasksModerators.Use(middleware.RequireRolesModerators())
	{
		tasksModerators.PUT("/:id", taskHandler.Update)
		tasksModerators.PATCH("/:id/status", taskHandler.UpdateStatus)
		tasksModerators.DELETE("/:id", taskHandler.Delete)
	}
}
