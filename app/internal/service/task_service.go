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
	Create(ctx context.Context, identity *auth.Identity, projectID string, name string, content string) (*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetByProjectID(ctx context.Context, projectID string) ([]models.Task, error)
	Update(ctx context.Context, identity *auth.Identity, id string, name *string, content *string, executiveID *string, auditorID *string) error
	UpdateStatus(ctx context.Context, identity *auth.Identity, id string, status string) error
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

func (s *taskService) Create(ctx context.Context, identity *auth.Identity, projectID string, name string, content string) (*models.Task, error) {
	project, err := s.projectService.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	task := &models.Task{
		ID:          utils.NewUUID(),
		ProjectID:   project.ID,
		Name:        name,
		Content:     content,
		Status:      models.TaskStatusCreated,
		AuditorID:   identity.UserID,
		ExecutiveID: identity.UserID,
		CreatedBy:   identity.UserID,
		UpdatedBy:   identity.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.enrichTaskWithUsers(ctx, task, map[string]*models.UserBrief{})
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

	tasks, err := s.repo.GetByProject(ctx, project)
	if err != nil {
		return nil, err
	}

	userCache := map[string]*models.UserBrief{}
	for i := range tasks {
		s.enrichTaskWithUsers(ctx, &tasks[i], userCache)
	}

	return tasks, nil
}

func (s *taskService) Update(ctx context.Context, identity *auth.Identity, id string, name *string, content *string, executiveID *string, auditorID *string) error {
	task, err := s.repo.GetByID(ctx, id) // 14000130cc0
	if err != nil {
		return err
	}
	if task == nil {
		return myerrors.EntityNotFound("task")
	}

	var changed = false
	if name != nil && *name != "" && task.Name != *name {
		task.Name = *name
		changed = true
	}

	if content != nil && task.Content != *content {
		task.Content = *content
		changed = true
	}

	if executiveID != nil && task.ExecutiveID != *executiveID {
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

	if auditorID != nil && task.AuditorID != *auditorID {
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
		task.UpdatedBy = identity.UserID
		task.UpdatedAt = time.Now()

		err = s.repo.Update(ctx, task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *taskService) UpdateStatus(ctx context.Context, identity *auth.Identity, id string, status string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Println("get task err:", err)
		return err
	}
	if task == nil {
		log.Println("get task err: task not found")
		return myerrors.EntityNotFound("task")
	}

	taskStatus := models.ParseTaskStatus(status)
	if taskStatus == nil {
		log.Println("ParseTaskStatus err: invalid task status")
		return myerrors.InvalidTaskStatus()
	}

	if identity.Role == auth.UserRoleViewer {
		return myerrors.ForbiddenAction("viewer cannot update task")
	}

	if identity.Role == auth.UserRoleEditor &&
		task.Status == models.TaskStatusDone &&
		*taskStatus != models.TaskStatusDone {
		return myerrors.ForbiddenAction("editor cannot change status of done task")
	}

	task.Status = *taskStatus
	task.UpdatedAt = time.Now()
	task.UpdatedBy = identity.UserID

	return s.repo.Update(ctx, task)
}

func (s *taskService) Delete(ctx context.Context, id string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil || task == nil {
		return myerrors.EntityNotFound("task")
	}

	return s.repo.Delete(ctx, task)
}

func (s *taskService) enrichTaskWithUsers(ctx context.Context, task *models.Task, userCache map[string]*models.UserBrief) {
	if task == nil {
		return
	}

	task.ExecutiveUser = s.getUserBriefByID(ctx, task.ExecutiveID, userCache)
	task.AuditorUser = s.getUserBriefByID(ctx, task.AuditorID, userCache)
	task.CreatedByUser = s.getUserBriefByID(ctx, task.CreatedBy, userCache)
	task.UpdatedByUser = s.getUserBriefByID(ctx, task.UpdatedBy, userCache)
}

func (s *taskService) getUserBriefByID(ctx context.Context, userID string, userCache map[string]*models.UserBrief) *models.UserBrief {
	if userID == "" {
		return nil
	}

	if cached, ok := userCache[userID]; ok {
		return cached
	}

	user, err := s.userService.GetById(ctx, userID)
	if err != nil || user == nil {
		userCache[userID] = nil
		return nil
	}

	brief := &models.UserBrief{
		ID:   user.ID,
		Name: user.Name,
	}
	userCache[userID] = brief

	return brief
}
