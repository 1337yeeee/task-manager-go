package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/middleware"
	"task-manager/internal/responses"
	"task-manager/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		responses.InvalidCredentials(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	identity := auth.FromContext(ctx.Request.Context())

	tokenRaw, exists := ctx.Get(middleware.RefreshTokenContextKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	refreshToken := tokenRaw.(string)

	access, refresh, err := h.service.RefreshToken(ctx, identity, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	identity := auth.FromContext(ctx.Request.Context())

	err := h.service.Logout(ctx, identity)
	if err != nil {
		responses.InvalidCredentials(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
