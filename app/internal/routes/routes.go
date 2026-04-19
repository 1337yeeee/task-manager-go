package routes

import (
	"github.com/gin-gonic/gin"

	"task-manager/internal/app"
	"task-manager/internal/middleware"
)

func SetupRouter(container *app.Container) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	registerSwaggerRoute(api)
	registerAuthRoutes(api, container)
	registerProtectedRoutes(api, container)

	return r
}

func registerAuthRoutes(api *gin.RouterGroup, container *app.Container) {
	authHandler := container.AuthHandler

	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)

	refresh := api.Group("")
	refresh.Use(middleware.JWTRefreshMiddleware(container.TokenManager))
	refresh.POST("/refresh", authHandler.Refresh)

	logout := api.Group("")
	logout.Use(middleware.JWTAccessMiddleware(container.TokenManager, container.UserRepository))
	logout.Any("/logout", authHandler.Logout)
}

func registerProtectedRoutes(api *gin.RouterGroup, container *app.Container) {
	userHandler := container.UserHandler
	projectHandler := container.ProjectHandler
	taskHandler := container.TaskHandler

	protected := api.Group("")
	protected.Use(middleware.JWTAccessMiddleware(container.TokenManager, container.UserRepository))

	// Users
	usersRead := protected.Group("/users")
	usersRead.Use(middleware.RequireRolesModerators())
	{
		usersRead.GET("", userHandler.GetAll)
	}

	usersAdmin := protected.Group("/users")
	usersAdmin.Use(middleware.RequireRole("admin"))
	{
		usersAdmin.GET("/:id", userHandler.GetByID)
		usersAdmin.PUT("/:id", userHandler.Update)
		usersAdmin.DELETE("/:id", userHandler.Delete)
		usersAdmin.POST("", userHandler.Register)
	}

	// Projects
	projects := protected.Group("/projects")
	{
		projects.GET("", projectHandler.GetAll)
		projects.GET("/:id", projectHandler.GetByID)
		projects.GET("/:id/tasks", taskHandler.GetByProject)
	}

	projectsModerators := protected.Group("/projects")
	projectsModerators.Use(middleware.RequireRolesModerators())
	{
		projectsModerators.POST("", projectHandler.Create)
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
