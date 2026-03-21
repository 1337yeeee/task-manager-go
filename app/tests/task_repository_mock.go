package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/domain/models"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetByProject(ctx context.Context, project *models.Project) ([]models.Task, error) {
	args := m.Called(ctx, project)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}
