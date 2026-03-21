package service_test

import (
	"context"
	"task-manager/tests"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task-manager/internal/auth"
	"task-manager/internal/config"
	"task-manager/internal/domain/models"
	"task-manager/internal/service"
	"task-manager/internal/utils"
)

func TestAuthService_Login_Success(t *testing.T) {
	repo := new(tests.MockUserRepository)

	password := "password123"
	hash, _ := utils.HashPassword(password)

	user := &models.User{
		ID:       "user-id",
		Email:    "test@test.com",
		Password: hash,
		Role:     auth.UserRole("admin"),
	}

	repo.On("FindUserByEmail", mock.Anything, user.Email).
		Return(user, nil)

	cfg := config.Config{
		JWTSecret: "secret",
	}

	authService := service.NewAuthService(repo, cfg)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		user.Email,
		password,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	repo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := new(tests.MockUserRepository)

	hash, _ := utils.HashPassword("correct-password")

	user := &models.User{
		ID:       "user-id",
		Email:    "test@test.com",
		Password: hash,
		Role:     auth.UserRole("admin"),
	}

	repo.On("FindUserByEmail", mock.Anything, user.Email).
		Return(user, nil)

	cfg := config.Config{
		JWTSecret: "secret",
	}

	authService := service.NewAuthService(repo, cfg)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		user.Email,
		"wrong-password",
	)

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := new(tests.MockUserRepository)

	repo.On("FindUserByEmail", mock.Anything, "test@test.com").
		Return(&models.User{}, assert.AnError)

	cfg := config.Config{
		JWTSecret: "secret",
	}

	authService := service.NewAuthService(repo, cfg)

	accessToken, refreshToken, err := authService.Login(
		context.Background(),
		"test@test.com",
		"password",
	)

	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
}
