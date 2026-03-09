package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

type registerUserRequest struct {
	Name     string         `json:"name" binding:"required"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required"`
	Role     *auth.UserRole `json:"role" binding:"omitempty"`
}

type updateUserRequest struct {
	ID       string  `json:"id" binding:"required"`
	Name     *string `json:"name" binding:"omitempty"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty"`
}

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

func (h *UserHandler) GetAll(ctx *gin.Context) {
	users, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in getAll"})
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	user, err := h.service.GetById(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) Update(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Update(ctx.Request.Context(), req.ID, req.Name, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	userId := ctx.Param("id")
	err := h.service.Delete(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
