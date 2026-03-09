package service

import (
	"context"
	"log"
	"task-manager/internal/auth"
	"task-manager/internal/config"
	"task-manager/internal/domain/repository"
	"task-manager/internal/myerrors"
	"task-manager/internal/utils"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
}

type authService struct {
	repo repository.UserRepository
	cfg  config.Config
}

func NewAuthService(repo repository.UserRepository, cfg config.Config) AuthService {
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", myerrors.InvalidCredentials()
	}
	log.Println("User: ", user)

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", "", myerrors.InvalidCredentials()
	}

	log.Println(user.ID, user.Role)

	accessToken, refreshToken, err := s.generateTokens(user.ID, user.Role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	return accessToken, refreshToken, nil
}

func (s *authService) generateTokens(userID string, role auth.UserRole) (string, string, error) {
	tokenManager := utils.NewTokenManager(
		s.cfg.JWTSecret,
		utils.DefaultAccessTTL,
		utils.DefaultRefreshTTL,
	)

	accessToken, err := tokenManager.GenerateAccessToken(userID, role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	refreshToken, err := tokenManager.GenerateRefreshToken(userID, role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	return accessToken, refreshToken, nil
}
