package service_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/myerrors"
	"task-manager/internal/service"
	"task-manager/tests"
	"testing"
	"time"
)

func TestTaskCreation_Success(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const projectID = "project-id"
	const taskName = "project-name"
	const taskContent = "task-content"

	project := &models.Project{
		ID: projectID,
	}

	projectService.On("GetByID", mock.Anything, projectID).Return(project, nil).Once()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Task")).Return(nil).Once()

	identity := auth.NewIdentity("user-ID", auth.UserRoleAdmin)

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	taskCreated, err := taskService.Create(ctx, identity, projectID, taskName, taskContent)

	assert.NoError(t, err)
	assert.NotNil(t, taskCreated)
	assert.NotEmpty(t, taskCreated.ID)
	assert.Equal(t, projectID, taskCreated.ProjectID)
	assert.Equal(t, taskName, taskCreated.Name)
	assert.Equal(t, taskContent, taskCreated.Content)
	assert.Equal(t, identity.UserID, taskCreated.CreatedBy)
	assert.Equal(t, identity.UserID, taskCreated.UpdatedBy)
	assert.NotEmpty(t, taskCreated.UpdatedAt)
	assert.NotEmpty(t, taskCreated.CreatedAt)
	assert.Equal(t, models.TaskStatusCreated, taskCreated.Status)

	repo.AssertExpectations(t)
}

func TestTaskCreation_ProjectNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const projectID = "project-id"
	const taskName = "project-name"
	const taskContent = "task-content"

	projectService.On("GetByID", mock.Anything, projectID).Return((*models.Project)(nil), errors.New("not found")).Once()

	identity := auth.NewIdentity("user-ID", auth.UserRoleAdmin)

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	taskCreated, err := taskService.Create(ctx, identity, projectID, taskName, taskContent)

	assert.Error(t, err)
	assert.Nil(t, taskCreated)

	repo.AssertExpectations(t)
}

func TestTaskGetByID_Success(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	task := &models.Task{
		ID: taskID,
	}

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	taskGotten, err := taskService.GetByID(ctx, taskID)

	assert.NoError(t, err)
	assert.NotNil(t, taskGotten)
	assert.Equal(t, taskID, taskGotten.ID)

	repo.AssertExpectations(t)
}

func TestTaskGetByID_NotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	repo.On("GetByID", mock.Anything, taskID).Return((*models.Task)(nil), errors.New("not found")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	taskGotten, err := taskService.GetByID(ctx, taskID)

	assert.Error(t, err)
	assert.Nil(t, taskGotten)

	repo.AssertExpectations(t)
}

func TestTaskGetByProjectID_Success(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const projectID = "project-id"

	project := &models.Project{
		ID: projectID,
	}

	tasks := []models.Task{
		{
			ID:        "task-id-1",
			Name:      "task-name-1",
			Content:   "task-content-1",
			Status:    models.TaskStatusCreated,
			ProjectID: projectID,
		},
		{
			ID:        "task-id-2",
			Name:      "task-name-2",
			Content:   "task-content-2",
			Status:    models.TaskStatusCreated,
			ProjectID: projectID,
		},
	}

	projectService.On("GetByID", mock.Anything, projectID).Return(project, nil).Once()
	repo.On("GetByProject", mock.Anything, mock.AnythingOfType("*models.Project")).Return(tasks, nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	tasksGotten, err := taskService.GetByProjectID(ctx, projectID)

	assert.NoError(t, err)
	assert.NotNil(t, tasksGotten)
	assert.Equal(t, len(tasks), len(tasksGotten))

	repo.AssertExpectations(t)
}

func TestTaskGetByProjectID_ProjectNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const projectID = "project-id"

	projectService.On("GetByID", mock.Anything, projectID).Return((*models.Project)(nil), myerrors.EntityNotFound("project")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	tasksGotten, err := taskService.GetByProjectID(ctx, projectID)

	assert.Error(t, err)
	assert.Nil(t, tasksGotten)

	repo.AssertExpectations(t)
}

func TestTaskUpdate_Table(t *testing.T) {
	const taskID = "task-id"
	const projectID = "project-id"

	const taskNameOld = "task-name-old"
	const taskNameNew = "task-name-new"

	const taskContentOld = "task-content-old"
	const taskContentNew = "task-content-new"

	const oldIdentityID = "user-ID-1"
	const newIdentityID = "user-ID-2"

	timestamp := time.Now().Add(-time.Hour)

	identityOld := auth.NewIdentity(oldIdentityID, auth.UserRoleAdmin)
	identityNew := auth.NewIdentity(newIdentityID, auth.UserRoleAdmin)

	stringPtr := func(s string) *string {
		return &s
	}

	createTask := func() *models.Task {
		return &models.Task{
			ID:          taskID,
			ProjectID:   projectID,
			Name:        taskNameOld,
			Content:     taskContentOld,
			ExecutiveID: identityOld.UserID,
			AuditorID:   identityOld.UserID,
			UpdatedBy:   identityOld.UserID,
			CreatedBy:   identityOld.UserID,
			UpdatedAt:   timestamp,
			CreatedAt:   timestamp,
		}
	}

	setupMocks := func(
		mockRepo *tests.MockTaskRepository,
		mockUserService *tests.MockUserService,
		task *models.Task,
		needUpdate bool,
	) {
		mockRepo.On("GetByID", mock.Anything, task.ID).
			Return(task, nil)

		if needUpdate {
			mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Task")).
				Return(nil).Once()
		}

		mockUserService.
			On("GetById", mock.Anything, oldIdentityID).
			Return(&models.User{ID: identityOld.UserID}, nil).
			Maybe()

		mockUserService.
			On("GetById", mock.Anything, newIdentityID).
			Return(&models.User{ID: identityNew.UserID}, nil).
			Maybe()
	}

	currentTests := []struct {
		name string

		updateName      *string
		updateContent   *string
		updateExecutive *string
		updateAuditor   *string

		expectedName      string
		expectedContent   string
		expectedExecutive string
		expectedAuditor   string

		isUpdated bool
	}{
		{
			name: "Update ALL",

			updateName:      stringPtr(taskNameNew),
			updateContent:   stringPtr(taskContentNew),
			updateExecutive: stringPtr(identityNew.UserID),
			updateAuditor:   stringPtr(identityNew.UserID),

			expectedName:      taskNameNew,
			expectedContent:   taskContentNew,
			expectedExecutive: identityNew.UserID,
			expectedAuditor:   identityNew.UserID,

			isUpdated: true,
		},
		{
			name: "Update only name",

			updateName:      stringPtr(taskNameNew),
			updateContent:   nil,
			updateExecutive: nil,
			updateAuditor:   nil,

			expectedName:      taskNameNew,
			expectedContent:   taskContentOld,
			expectedExecutive: identityOld.UserID,
			expectedAuditor:   identityOld.UserID,

			isUpdated: true,
		},
		{
			name: "Update only content",

			updateName:      nil,
			updateContent:   stringPtr(taskContentNew),
			updateExecutive: nil,
			updateAuditor:   nil,

			expectedName:      taskNameOld,
			expectedContent:   taskContentNew,
			expectedExecutive: identityOld.UserID,
			expectedAuditor:   identityOld.UserID,

			isUpdated: true,
		},
		{
			name: "Update only Executive",

			updateName:      nil,
			updateContent:   nil,
			updateExecutive: stringPtr(identityNew.UserID),
			updateAuditor:   nil,

			expectedName:      taskNameOld,
			expectedContent:   taskContentOld,
			expectedExecutive: identityNew.UserID,
			expectedAuditor:   identityOld.UserID,

			isUpdated: true,
		},
		{
			name: "Update only Auditor",

			updateName:      nil,
			updateContent:   nil,
			updateExecutive: nil,
			updateAuditor:   stringPtr(identityNew.UserID),

			expectedName:      taskNameOld,
			expectedContent:   taskContentOld,
			expectedExecutive: identityOld.UserID,
			expectedAuditor:   identityNew.UserID,

			isUpdated: true,
		},
		{
			name: "Update only name but no nil pointers",

			updateName:      stringPtr(taskNameNew),
			updateContent:   stringPtr(taskContentOld),
			updateExecutive: stringPtr(identityOld.UserID),
			updateAuditor:   stringPtr(identityOld.UserID),

			expectedName:      taskNameNew,
			expectedContent:   taskContentOld,
			expectedExecutive: identityOld.UserID,
			expectedAuditor:   identityOld.UserID,

			isUpdated: true,
		},
		{
			name: "Nothing to update",

			updateName:      nil,
			updateContent:   nil,
			updateExecutive: nil,
			updateAuditor:   nil,

			expectedName:      taskNameOld,
			expectedContent:   taskContentOld,
			expectedExecutive: identityOld.UserID,
			expectedAuditor:   identityOld.UserID,

			isUpdated: false,
		},
	}

	ctx := context.Background()

	for _, tt := range currentTests {
		t.Run(tt.name, func(t *testing.T) {
			task := createTask()

			mockRepo := new(tests.MockTaskRepository)
			mockUserService := new(tests.MockUserService)
			mockProjectService := new(tests.MockProjectService)

			setupMocks(mockRepo, mockUserService, task, tt.isUpdated)

			taskService := service.NewTaskService(mockRepo, mockProjectService, mockUserService)

			err := taskService.Update(ctx, identityNew, taskID, tt.updateName, tt.updateContent, tt.updateExecutive, tt.updateAuditor)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedName, task.Name)
			assert.Equal(t, tt.expectedContent, task.Content)
			assert.Equal(t, tt.expectedExecutive, task.ExecutiveID)
			assert.Equal(t, tt.expectedAuditor, task.AuditorID)

			assert.Equal(t, taskID, task.ID)
			assert.Equal(t, projectID, task.ProjectID)
			assert.Equal(t, identityOld.UserID, task.CreatedBy)
			assert.Equal(t, timestamp, task.CreatedAt)

			if tt.isUpdated {
				assert.Equal(t, identityNew.UserID, task.UpdatedBy)
				assert.NotEqual(t, timestamp, task.UpdatedAt)
			} else {
				mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			}

			mockRepo.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestTaskUpdate_TaskNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"
	const taskNameNew = "task-name-new"
	const taskContentNew = "task-content-new"

	identityNew := auth.NewIdentity("user-ID-2", auth.UserRoleAdmin)

	userAuditor := &models.User{
		ID:   "user-auditor",
		Name: "Auditor",
	}

	repo.On("GetByID", mock.Anything, taskID).Return((*models.Task)(nil), myerrors.EntityNotFound("task")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	newTaskName := taskNameNew
	newTaskContent := taskContentNew
	newAuditorID := userAuditor.ID

	err := taskService.Update(ctx, identityNew, taskID, &newTaskName, &newTaskContent, nil, &newAuditorID)

	assert.Error(t, err)

	repo.AssertExpectations(t)
	projectService.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestTaskUpdate_UserNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"
	const projectID = "project-id"

	const taskNameOld = "task-name-old"
	const taskNameNew = "task-name-new"

	const taskContentOld = "task-content-old"
	const taskContentNew = "task-content-new"

	var timestamp = time.Now()

	identityOld := auth.NewIdentity("user-ID-1", auth.UserRoleAdmin)
	identityNew := auth.NewIdentity("user-ID-2", auth.UserRoleAdmin)

	userAuditor := &models.User{
		ID:   "user-auditor",
		Name: "Auditor",
	}

	task := &models.Task{
		ID:          taskID,
		ProjectID:   projectID,
		Name:        taskNameOld,
		Content:     taskContentOld,
		ExecutiveID: identityOld.UserID,
		AuditorID:   identityOld.UserID,
		CreatedBy:   identityOld.UserID,
		UpdatedBy:   identityOld.UserID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()
	userService.On("GetById", mock.Anything, mock.Anything).Return((*models.User)(nil), myerrors.EntityNotFound("user")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	newTaskName := taskNameNew
	newTaskContent := taskContentNew
	newAuditorID := userAuditor.ID

	err := taskService.Update(ctx, identityNew, taskID, &newTaskName, &newTaskContent, nil, &newAuditorID)

	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_Success(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"
	const projectID = "project-id"

	var timestamp = time.Now()

	identityOld := auth.NewIdentity("user-ID-1", auth.UserRoleAdmin)
	identityNew := auth.NewIdentity("user-ID-2", auth.UserRoleAdmin)

	task := &models.Task{
		ID:        taskID,
		ProjectID: projectID,
		CreatedBy: identityOld.UserID,
		UpdatedBy: identityOld.UserID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Status:    models.TaskStatusCreated,
	}

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Task")).Return(nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	err := taskService.UpdateStatus(ctx, identityNew, taskID, models.TaskStatusDone.String())

	assert.NoError(t, err)
	assert.Equal(t, task.ID, taskID)
	assert.Equal(t, task.ProjectID, projectID)
	assert.Equal(t, task.Status, models.TaskStatusDone)
	assert.Equal(t, task.UpdatedBy, identityNew.UserID)
	assert.Equal(t, task.CreatedBy, identityOld.UserID)
	assert.NotEqual(t, task.UpdatedAt, timestamp)
	assert.Equal(t, task.CreatedAt, timestamp)

	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_WrongStatus(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"
	const projectID = "project-id"
	const wrongStatus = "wrong-status"

	var timestamp = time.Now()

	identityOld := auth.NewIdentity("user-ID-1", auth.UserRoleAdmin)
	identityNew := auth.NewIdentity("user-ID-2", auth.UserRoleAdmin)

	task := &models.Task{
		ID:        taskID,
		ProjectID: projectID,
		CreatedBy: identityOld.UserID,
		UpdatedBy: identityOld.UserID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Status:    models.TaskStatusCreated,
	}

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	err := taskService.UpdateStatus(ctx, identityNew, taskID, wrongStatus)

	assert.Error(t, err)
	assert.Equal(t, err, myerrors.InvalidTaskStatus())
	assert.Equal(t, task.ID, taskID)
	assert.Equal(t, task.ProjectID, projectID)
	assert.NotEqual(t, task.Status.String(), wrongStatus)
	assert.Equal(t, task.Status, models.TaskStatusCreated)
	assert.Equal(t, task.UpdatedBy, identityOld.UserID)
	assert.Equal(t, task.CreatedBy, identityOld.UserID)
	assert.Equal(t, task.UpdatedAt, timestamp)
	assert.Equal(t, task.CreatedAt, timestamp)

	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_ViewerForbidden(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	task := &models.Task{
		ID:     taskID,
		Status: models.TaskStatusCreated,
	}

	viewerIdentity := auth.NewIdentity("viewer-id", auth.UserRoleViewer)

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	err := taskService.UpdateStatus(context.Background(), viewerIdentity, taskID, models.TaskStatusDone.String())

	assert.Error(t, err)
	assert.Equal(t, "viewer cannot update task", err.Error())
	assert.Equal(t, models.TaskStatusCreated, task.Status)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_EditorCannotMoveFromDone(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	task := &models.Task{
		ID:     taskID,
		Status: models.TaskStatusDone,
	}

	editorIdentity := auth.NewIdentity("editor-id", auth.UserRoleEditor)

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	err := taskService.UpdateStatus(context.Background(), editorIdentity, taskID, models.TaskStatusInProgress.String())

	assert.Error(t, err)
	assert.Equal(t, "editor cannot change status of done task", err.Error())
	assert.Equal(t, models.TaskStatusDone, task.Status)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_EditorCanMoveToDone(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	task := &models.Task{
		ID:     taskID,
		Status: models.TaskStatusAudit,
	}

	editorIdentity := auth.NewIdentity("editor-id", auth.UserRoleEditor)

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()
	repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Task")).Return(nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	err := taskService.UpdateStatus(context.Background(), editorIdentity, taskID, models.TaskStatusDone.String())

	assert.NoError(t, err)
	assert.Equal(t, models.TaskStatusDone, task.Status)
	assert.Equal(t, editorIdentity.UserID, task.UpdatedBy)
	repo.AssertExpectations(t)
}

func TestTaskUpdateStatus_TaskNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	identityNew := auth.NewIdentity("user-ID-1", auth.UserRoleAdmin)

	const taskID = "task-id"

	repo.On("GetByID", mock.Anything, taskID).Return((*models.Task)(nil), myerrors.EntityNotFound("task")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	err := taskService.UpdateStatus(ctx, identityNew, taskID, models.TaskStatusDone.String())

	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestTaskDelete_Success(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"
	const projectID = "project-id"

	var timestamp = time.Now()

	identityOld := auth.NewIdentity("user-ID-1", auth.UserRoleAdmin)

	task := &models.Task{
		ID:        taskID,
		ProjectID: projectID,
		CreatedBy: identityOld.UserID,
		UpdatedBy: identityOld.UserID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Status:    models.TaskStatusCreated,
	}

	repo.On("GetByID", mock.Anything, taskID).Return(task, nil).Once()
	repo.On("Delete", mock.Anything, mock.AnythingOfType("*models.Task")).Return(nil).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	err := taskService.Delete(ctx, taskID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTaskDelete_TaskNotFound(t *testing.T) {
	repo := new(tests.MockTaskRepository)
	projectService := new(tests.MockProjectService)
	userService := new(tests.MockUserService)

	const taskID = "task-id"

	repo.On("GetByID", mock.Anything, taskID).Return((*models.Task)(nil), errors.New("NOT FOUND")).Once()

	taskService := service.NewTaskService(repo, projectService, userService)
	ctx := context.Background()

	err := taskService.Delete(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, err, myerrors.EntityNotFound("task"))
	repo.AssertExpectations(t)
}
