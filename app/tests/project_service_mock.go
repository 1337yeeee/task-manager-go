package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/service"
)

type MockProjectService struct {
	mock.Mock
}

func NewProjectServiceMock() service.ProjectService {
	return &MockProjectService{}
}

func (p *MockProjectService) GetAll(ctx context.Context) ([]models.Project, error) {
	args := p.Called(ctx)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (p *MockProjectService) Create(ctx context.Context, identity *auth.Identity, name string, desc string) (*models.Project, error) {
	args := p.Called(ctx, identity, name, desc)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (p *MockProjectService) GetByID(ctx context.Context, id string) (*models.Project, error) {
	args := p.Called(ctx, id)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (p *MockProjectService) Update(ctx context.Context, identity *auth.Identity, id string, name *string, desc *string) (*models.Project, error) {
	args := p.Called(ctx, identity, id, name, desc)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (p *MockProjectService) Delete(ctx context.Context, id string) error {
	args := p.Called(ctx, id)
	return args.Error(0)
}
