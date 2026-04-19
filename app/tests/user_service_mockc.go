package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/domain/repository"
	"task-manager/internal/service"
)

type MockUserService struct {
	mock.Mock
}

func NewUserServiceMock(repo repository.ProjectRepository) service.UserService {
	return &MockUserService{}
}

func (m *MockUserService) Register(ctx context.Context, name string, email string, password string, role *auth.UserRole) (*models.User, error) {
	args := m.Called(ctx, name, email, password, role)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetById(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, ID string, name *string, email *string, password *string, role *auth.UserRole, isActive *bool) (*models.User, error) {
	args := m.Called(ctx, ID, name, email, password, role, isActive)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
