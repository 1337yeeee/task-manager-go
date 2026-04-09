package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/middleware"
	"task-manager/internal/responses"
	"task-manager/internal/service"
	"task-manager/internal/utils"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// loginRequest represents login payload
type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// TokenResponse represents auth tokens
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"jwt-access-token"`
	TokenType   string `json:"token_type" example:"Bearer"`
}

// ErrorResponse standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// IdentityResponse (optional, но полезно для Swagger)
type IdentityResponse struct {
	UserID string `json:"user_id" example:"uuid"`
	Role   string `json:"role" example:"admin"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access + refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse "validation error"
// @Failure 401 {object} ErrorResponse "invalid credentials"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	access, refresh, err := h.service.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		responses.InvalidCredentials(ctx)
		return
	}

	h.setRefreshTokenCookie(ctx, refresh)
	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken: access,
		TokenType:   "Bearer",
	})
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Requires refresh token in Authorization header (Bearer). Returns new token pair.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TokenResponse
// @Failure 401 {object} ErrorResponse "invalid or missing refresh token"
// @Router /api/refresh [post]
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	identity := auth.FromContext(ctx.Request.Context())

	tokenRaw, exists := ctx.Get(middleware.RefreshTokenContextKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing refresh token"})
		return
	}

	refreshToken := tokenRaw.(string)

	access, refresh, err := h.service.RefreshToken(ctx, identity, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	h.setRefreshTokenCookie(ctx, refresh)
	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken: access,
		TokenType:   "Bearer",
	})
}

// Logout godoc
// @Summary Logout user
// @Description Requires access token. Invalidates session/token.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "empty response"
// @Failure 401 {object} ErrorResponse "unauthorized"
// @Router /api/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	identity := auth.FromContext(ctx.Request.Context())

	err := h.service.Logout(ctx, identity)
	if err != nil {
		responses.InvalidCredentials(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *AuthHandler) setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(
		middleware.RefreshTokenCookieName,
		refreshToken,
		int(utils.DefaultRefreshTTL.Seconds()),
		"/",
		"SameSite",
		false,
		true,
	)
}
