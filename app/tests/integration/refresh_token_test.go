package integration

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain/models"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"
	"task-manager/internal/service"
	"task-manager/internal/utils"
	"task-manager/tests"
	"testing"
)

func TestRefreshToken_Success(t *testing.T) {
	r, tokenManager, identity := setupAuthTest()

	userRepoMock := &tests.MockUserRepository{}
	authRepoMock := &tests.MockAuthRepository{}

	const UserEmail = "user-email"
	const UserPassword = "password"
	userPasswordHash, _ := utils.HashPassword(UserPassword)

	user := &models.User{
		ID:       identity.UserID,
		Email:    UserEmail,
		Password: userPasswordHash,
	}

	userRepoMock.On("FindUserByEmail", mock.Anything, UserEmail).Return(user, nil)
	authRepoMock.On("Store", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	ctx := context.Background()

	authService := service.NewAuthService(userRepoMock, authRepoMock, tokenManager)
	access, refresh, err := authService.Login(ctx, UserEmail, UserPassword)
	if err != nil {
		t.Fatal(err)
	}

	authRepoMock.On("GetByUserID", mock.Anything, identity.UserID).Return(refresh, nil)

	authHandler := handler.NewAuthHandler(authService)

	r.Use(middleware.JWTRefreshMiddleware(tokenManager))
	r.POST("/refresh", authHandler.Refresh)

	req, _ := http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+refresh)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "access_token")
	assert.Contains(t, response, "refresh_token")
	assert.Contains(t, response, "token_type")
	assert.NotEqual(t, access, response["access_token"])
	assert.NotEqual(t, refresh, response["refresh_token"])

	userRepoMock.AssertExpectations(t)
	authRepoMock.AssertExpectations(t)
}

func TestRefreshToken_Fail_AccessTokenProvided(t *testing.T) {
	r, tokenManager, identity := setupAuthTest()

	userRepoMock := &tests.MockUserRepository{}
	authRepoMock := &tests.MockAuthRepository{}

	const UserEmail = "user-email"
	const UserPassword = "password"
	userPasswordHash, _ := utils.HashPassword(UserPassword)

	user := &models.User{
		ID:       identity.UserID,
		Email:    UserEmail,
		Password: userPasswordHash,
	}

	userRepoMock.On("FindUserByEmail", mock.Anything, UserEmail).Return(user, nil)
	authRepoMock.On("Store", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	ctx := context.Background()

	authService := service.NewAuthService(userRepoMock, authRepoMock, tokenManager)
	access, _, err := authService.Login(ctx, UserEmail, UserPassword)
	if err != nil {
		t.Fatal(err)
	}

	authHandler := handler.NewAuthHandler(authService)

	r.Use(middleware.JWTRefreshMiddleware(tokenManager))
	r.POST("/refresh", authHandler.Refresh)

	req, _ := http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+access)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "invalid refresh token", response["error"])

	userRepoMock.AssertExpectations(t)
	authRepoMock.AssertExpectations(t)
}
