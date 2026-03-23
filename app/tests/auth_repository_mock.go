package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Store(ctx context.Context, userID string, token string, ttl time.Duration) error {
	args := m.Called(ctx, userID, token, ttl)
	return args.Error(0)
}

func (m *MockAuthRepository) GetByUserID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepository) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
