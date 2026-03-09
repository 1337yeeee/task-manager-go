package repository

import (
	"context"
	"gorm.io/gorm"
	"task-manager/internal/domain/models"
)

type ProjectRepository interface {
	GetAll(ctx context.Context) ([]models.Project, error)
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r projectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	var projects []models.Project
	return projects, r.db.WithContext(ctx).Find(&projects).Error
}

func (r projectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r projectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	return &project, r.db.WithContext(ctx).First(&project, "id = ?", id).Error
}

func (r projectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r projectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}
