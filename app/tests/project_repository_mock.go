package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/domain/models"
)

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	args := m.Called()
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}
