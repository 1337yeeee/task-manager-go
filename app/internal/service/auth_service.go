package service

import (
	"context"
	"errors"
	"log"
	"task-manager/internal/auth"
	"task-manager/internal/domain/repository"
	"task-manager/internal/myerrors"
	"task-manager/internal/utils"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshToken(ctx context.Context, identity *auth.Identity, token string) (string, string, error)
	Logout(ctx context.Context, identity *auth.Identity) error
}

type authService struct {
	userRepository repository.UserRepository
	authRepository repository.AuthRepository
	tokenManager   utils.TokenManager
}

func NewAuthService(repo repository.UserRepository, authRepo repository.AuthRepository, tokenManager utils.TokenManager) AuthService {
	return &authService{
		userRepository: repo,
		authRepository: authRepo,
		tokenManager:   tokenManager,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		log.Println("error finding user in authService.Login", err)
		return "", "", myerrors.InvalidCredentials()
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		log.Println("error checking password in authService.Login: ", err)
		return "", "", myerrors.InvalidCredentials()
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID, user.Role)
	if err != nil {
		log.Println("error generating tokens in authService.Login: ", err)
		return "", "", myerrors.CouldNotCreateToken()
	}

	err = s.storeToken(ctx, user.ID, refreshToken)
	if err != nil {
		log.Println("error storing token in authService.Login", err)
		return "", "", myerrors.CouldNotCreateToken()
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, identity *auth.Identity, token string) (string, string, error) {
	storedToken, err := s.authRepository.GetByUserID(ctx, identity.UserID)
	if err != nil {
		log.Println("error getting stored refresh token", err)
		return "", "", myerrors.InvalidCredentials()
	}

	if utils.CompareHash(token, storedToken) != 0 {
		log.Println("token from request not equal to token in the store in authService.Login")
		return "", "", myerrors.InvalidCredentials()
	}

	accessToken, refreshToken, err := s.generateTokens(identity.UserID, identity.Role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	err = s.storeToken(ctx, identity.UserID, refreshToken)
	if err != nil {
		log.Println(err)
		return "", "", myerrors.CouldNotCreateToken()
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(ctx context.Context, identity *auth.Identity) error {
	err := s.authRepository.Delete(ctx, identity.UserID)
	if err != nil {
		log.Println("error deleting token in authService.Logout: ", err)
		return errors.New("error in auth.logout")
	}

	return nil
}

func (s *authService) generateTokens(userID string, role auth.UserRole) (string, string, error) {
	accessToken, err := s.tokenManager.GenerateAccessToken(userID, role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(userID, role)
	if err != nil {
		return "", "", myerrors.CouldNotCreateToken()
	}

	return accessToken, refreshToken, nil
}

func (s *authService) storeToken(ctx context.Context, userID string, token string) error {
	tokenHash := utils.Hash(token)

	return s.authRepository.Store(ctx, userID, tokenHash, utils.DefaultRefreshTTL)
}
