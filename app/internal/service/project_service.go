package service

import (
	"context"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/domain/repository"
	"task-manager/internal/utils"
	"time"
)

type ProjectService interface {
	GetAll(ctx context.Context) ([]models.Project, error)
	Create(ctx context.Context, identity *auth.Identity, name string, desc string) (*models.Project, error)
	GetByID(ctx context.Context, id string) (*models.Project, error)
	Update(ctx context.Context, identity *auth.Identity, id string, name *string, desc *string) (*models.Project, error)
	Delete(ctx context.Context, id string) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return projectService{repo: repo}
}

func (s projectService) GetAll(ctx context.Context) ([]models.Project, error) {
	projects, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s projectService) Create(ctx context.Context, identity *auth.Identity, name string, desc string) (*models.Project, error) {
	project := models.Project{
		Name: name,
		Desc: desc,
	}

	project.ID = utils.NewUUID()
	project.CreatedBy = identity.UserID
	project.UpdatedBy = identity.UserID
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	err := s.repo.Create(ctx, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s projectService) GetByID(ctx context.Context, id string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s projectService) Update(ctx context.Context, identity *auth.Identity, id string, name *string, desc *string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var changed = false

	if name != nil && project.Name != *name {
		project.Name = *name
		changed = true
	}

	if desc != nil && project.Desc != *desc {
		project.Desc = *desc
		changed = true
	}

	if changed {
		project.UpdatedBy = identity.UserID
		project.UpdatedAt = time.Now()
		err = s.repo.Update(ctx, project)
		if err != nil {
			return nil, err
		}
	}
	return project, nil
}

func (s projectService) Delete(ctx context.Context, id string) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, project)
}
