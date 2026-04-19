package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/filters"
	"task-manager/internal/service"
	"time"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

type registerUserRequest struct {
	Name     string         `json:"name" binding:"required" example:"John Doe"`
	Email    string         `json:"email" binding:"required,email" example:"john@example.com"`
	Password string         `json:"password" binding:"required,min=8" example:"strongpassword123"`
	Role     *auth.UserRole `json:"role" binding:"omitempty" swaggertype:"string" example:"viewer"`
}

type updateUserRequest struct {
	Name     *string        `json:"name" binding:"omitempty" example:"John Smith"`
	Email    *string        `json:"email" binding:"omitempty,email" example:"john.smith@example.com"`
	Password *string        `json:"password" binding:"omitempty,min=8" example:"newstrongpassword123"`
	Role     *auth.UserRole `json:"role" binding:"omitempty" swaggertype:"string" example:"editor"`
	IsActive *bool          `json:"is_active" binding:"omitempty" example:"true"`
}

type UserResponse struct {
	ID        string    `json:"id" example:"7c3d9d0a-8dcb-4f7f-9b5e-2a12d4e9a001"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Role      string    `json:"role" example:"viewer" enums:"admin,viewer,editor"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-02T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-04-02T12:30:00Z"`
}

type UserDataResponse struct {
	User UserResponse `json:"user"`
}

type UsersDataResponse struct {
	Users []UserResponse `json:"users"`
}

// Register godoc
// @Summary Register user
// @Description Creates a new user. Endpoint is available only for admin according to route middleware.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body registerUserRequest true "User register payload"
// @Success 200 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/ [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req registerUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(ctx.Request.Context(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in register"})
		log.Default().Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAll godoc
// @Summary Get all users
// @Description Returns all users. Endpoint is available only for admin according to route middleware.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UsersDataResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/ [get]
func (h *UserHandler) GetAll(ctx *gin.Context) {
	filter, err := filters.ApplyUserFilter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.service.GetAll(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in getAll"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

// GetByID godoc
// @Summary Get user by ID
// @Description Returns user by identifier. Endpoint is available only for admin according to route middleware.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserDataResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(ctx *gin.Context) {
	user, err := h.service.GetById(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// Update godoc
// @Summary Update user
// @Description Updates user fields. Current implementation expects user id in request body.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID from route"
// @Param input body updateUserRequest true "User update payload"
// @Success 200 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if identity.Role != auth.UserRoleAdmin {
		req.Role = (*auth.UserRole)(nil)
		req.IsActive = (*bool)(nil)
	}

	id := ctx.Param("id")

	user, err := h.service.Update(ctx.Request.Context(), id, req.Name, req.Email, req.Password, req.Role, req.IsActive)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// Delete godoc
// @Summary Delete user
// @Description Deletes user by identifier. Endpoint is available only for admin according to route middleware.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} OkResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(ctx *gin.Context) {
	userId := ctx.Param("id")
	err := h.service.Delete(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
