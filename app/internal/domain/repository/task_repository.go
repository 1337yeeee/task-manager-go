package repository

import (
	"context"
	"gorm.io/gorm"
	"task-manager/internal/domain/models"
)

type TaskRepository interface {
	GetByProject(ctx context.Context, project *models.Project) ([]models.Task, error)
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, task *models.Task) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return taskRepository{db: db}
}

func (r taskRepository) GetByProject(ctx context.Context, project *models.Project) ([]models.Task, error) {
	var tasks []models.Task
	return tasks, r.db.WithContext(ctx).Where("project_id = ?", project.ID).Order("updated_at desc").Find(&tasks).Error
}

func (r taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r taskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	task := &models.Task{}
	return task, r.db.WithContext(ctx).First(task, "id = ?", id).Error
}

func (r taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Model(task).Updates(task).Error
}

func (r taskRepository) Delete(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Delete(task).Error
}
