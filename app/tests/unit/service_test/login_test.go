package service_test

import (
	"context"
	"task-manager/tests"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/service"
	"task-manager/internal/utils"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(tests.MockUserRepository)
	authRepo := new(tests.MockAuthRepository)

	password := "password123"
	hash, _ := utils.HashPassword(password)

	user := &models.User{
		ID:       "user-id",
		Email:    "test@test.com",
		Password: hash,
		Role:     auth.UserRole("admin"),
	}

	userRepo.On("FindUserByEmail", mock.Anything, user.Email).
		Return(user, nil)

	authRepo.On("Store", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	tokenManager := utils.NewTokenManager("test-secret", utils.DefaultAccessTTL, utils.DefaultRefreshTTL)

	authService := service.NewAuthService(userRepo, authRepo, tokenManager)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		user.Email,
		password,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	userRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := new(tests.MockUserRepository)
	authRepo := new(tests.MockAuthRepository)

	hash, _ := utils.HashPassword("correct-password")

	user := &models.User{
		ID:       "user-id",
		Email:    "test@test.com",
		Password: hash,
		Role:     auth.UserRole("admin"),
	}

	userRepo.On("FindUserByEmail", mock.Anything, user.Email).
		Return(user, nil)

	tokenManager := utils.NewTokenManager("test-secret", utils.DefaultAccessTTL, utils.DefaultRefreshTTL)

	authService := service.NewAuthService(userRepo, authRepo, tokenManager)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		user.Email,
		"wrong-password",
	)

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	userRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := new(tests.MockUserRepository)
	authRepo := new(tests.MockAuthRepository)

	userRepo.On("FindUserByEmail", mock.Anything, "test@test.com").
		Return(&models.User{}, assert.AnError)

	tokenManager := utils.NewTokenManager("test-secret", utils.DefaultAccessTTL, utils.DefaultRefreshTTL)

	authService := service.NewAuthService(userRepo, authRepo, tokenManager)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		"test@test.com",
		"password",
	)

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	userRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}
