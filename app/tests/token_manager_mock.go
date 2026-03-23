package tests

import (
	"github.com/stretchr/testify/mock"
	"task-manager/internal/auth"
	"task-manager/internal/utils"
)

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateAccessToken(userID string, role auth.UserRole) (string, error) {
	args := m.Called(userID, role)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTokenManager) GenerateRefreshToken(userID string, role auth.UserRole) (string, error) {
	args := m.Called(userID, role)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTokenManager) Parse(tokenStr string) (*utils.Claims, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*utils.Claims), args.Error(1)
}
