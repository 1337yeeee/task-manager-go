package service

import (
	"context"
	"log"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/domain/repository"
	"task-manager/internal/myerrors"
	"task-manager/internal/utils"
	"time"
)

type TaskService interface {
	Create(ctx context.Context, projectID string, name string, content string) (*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetByProjectID(ctx context.Context, projectID string) ([]models.Task, error)
	Update(ctx context.Context, id string, name *string, content *string, executiveID *string, auditorID *string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

type taskService struct {
	repo           repository.TaskRepository
	projectService ProjectService
	userService    UserService
}

type CreateTaskModel struct {
	ProjectID string
	Name      string
	Content   string
}

func NewTaskService(repo repository.TaskRepository, projectService ProjectService, userService UserService) TaskService {
	return &taskService{
		repo:           repo,
		projectService: projectService,
		userService:    userService,
	}
}

func (s *taskService) Create(ctx context.Context, projectID string, name string, content string) (*models.Task, error) {
	var task models.Task

	task.ID = utils.NewUUID()
	task.ProjectID = projectID
	task.Name = name
	task.Content = content
	task.Status = models.TaskStatusCreated

	identity := auth.IdentityFromContext(ctx)
	if identity == nil {
		log.Println("identity not found in context")
		return nil, myerrors.IdentityNotFoundInContext()
	}

	task.CreatedBy = identity.UserID
	task.UpdatedBy = identity.UserID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err := s.repo.Create(ctx, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *taskService) GetByID(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetByProjectID(ctx context.Context, projectID string) ([]models.Task, error) {
	project, err := s.projectService.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, myerrors.EntityNotFound("project")
	}

	task, err := s.repo.GetByProject(ctx, *project)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) Update(ctx context.Context, id string, name *string, content *string, executiveID *string, auditorID *string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return myerrors.EntityNotFound("task")
	}

	var changed = false
	if name != nil && *name != "" {
		task.Name = *name
		changed = true
	}

	if content != nil {
		task.Content = *content
		changed = true
	}

	if executiveID != nil {
		executive, err := s.userService.GetById(ctx, *executiveID)
		if err != nil {
			return err
		}
		if executive == nil {
			return myerrors.EntityNotFound("executive")
		}
		task.ExecutiveID = executive.ID
		changed = true
	}

	if auditorID != nil {
		auditor, err := s.userService.GetById(ctx, *auditorID)
		if err != nil {
			return err
		}
		if auditor == nil {
			return myerrors.EntityNotFound("auditor")
		}
		task.AuditorID = auditor.ID
		changed = true
	}

	if changed {
		identity := auth.IdentityFromContext(ctx)
		if identity == nil {
			log.Println("identity not found in context")
			return myerrors.IdentityNotFoundInContext()
		}
		task.UpdatedBy = identity.UserID
		task.UpdatedAt = time.Now()

		err = s.repo.Update(ctx, task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *taskService) UpdateStatus(ctx context.Context, id string, status string) error {
	identity := auth.IdentityFromContext(ctx)
	if identity == nil {
		log.Println("identity not found in context")
		return myerrors.IdentityNotFoundInContext()
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return myerrors.EntityNotFound("task")
	}

	taskStatus := models.ParseTaskStatus(status)
	if taskStatus == nil {
		return myerrors.InvalidTaskStatus()
	}

	task.Status = *taskStatus
	task.UpdatedAt = time.Now()
	task.UpdatedBy = identity.UserID

	return s.repo.Update(ctx, task)
}

func (s *taskService) Delete(ctx context.Context, id string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return myerrors.EntityNotFound("task")
	}

	return s.repo.Delete(ctx, task)
}
