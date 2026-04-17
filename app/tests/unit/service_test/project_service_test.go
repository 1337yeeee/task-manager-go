package service_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/service"
	"task-manager/tests"
	"testing"
	"time"
)

func TestProjectService_GetAll_Success(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	// Подготавливаем тестовые данные
	expectedProjects := []models.Project{
		{
			ID:          "project-1",
			Name:        "Project 1",
			Description: "Description 1",
			CreatedBy:   "user-1",
			UpdatedBy:   "user-1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "project-2",
			Name:        "Project 2",
			Description: "Description 2",
			CreatedBy:   "user-2",
			UpdatedBy:   "user-2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Настраиваем мок
	repo.On("GetAll", mock.Anything).Return(expectedProjects, nil)

	// Создаем сервис
	projectService := service.NewProjectService(repo)

	// Вызываем метод
	projects, err := projectService.GetAll(context.Background())

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, projects)
	assert.Equal(t, len(expectedProjects), len(projects))

	// Проверяем содержимое (сравниваем значимые поля)
	for i, expected := range expectedProjects {
		assert.Equal(t, expected.ID, projects[i].ID)
		assert.Equal(t, expected.Name, projects[i].Name)
		assert.Equal(t, expected.Description, projects[i].Description)
		assert.Equal(t, expected.CreatedBy, projects[i].CreatedBy)
	}

	repo.AssertExpectations(t)
}

func TestProjectCreation_Success(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	identity := auth.NewIdentity("user-ID", auth.UserRoleAdmin)

	projectName := "Project Name"
	projectDesc := "Project Description"

	// Перехватываем проект, который сервис передаст в repo.Create
	var capturedProject *models.Project

	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Project")).
		Run(func(args mock.Arguments) {
			// Захватываем проект, который сервис пытается сохранить
			capturedProject = args.Get(1).(*models.Project)
		}).
		Return(nil) // Репозиторий возвращает nil-ошибку

	projectService := service.NewProjectService(repo)

	created, err := projectService.Create(
		context.Background(),
		identity,
		projectName,
		projectDesc,
	)

	assert.NoError(t, err)
	assert.NotNil(t, created)

	repo.AssertExpectations(t)

	assert.NotEmpty(t, capturedProject.ID, "ID must be set")
	assert.Equal(t, projectName, capturedProject.Name)
	assert.Equal(t, projectDesc, capturedProject.Description)
	assert.Equal(t, identity.UserID, capturedProject.CreatedBy)
	assert.Equal(t, identity.UserID, capturedProject.UpdatedBy)
	assert.False(t, capturedProject.CreatedAt.IsZero(), "CreatedAt must be set")
	assert.False(t, capturedProject.UpdatedAt.IsZero(), "UpdatedAt must be set")

	assert.Equal(t, capturedProject, created)
}

func TestGetByID_Success(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	const projectID = "project-1"

	project1 := &models.Project{
		ID:          projectID,
		Name:        "Project 1",
		Description: "Description 1",
		CreatedBy:   "user-1",
		UpdatedBy:   "user-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.On("GetByID", mock.Anything, projectID).Return(project1, nil)

	projectService := service.NewProjectService(repo)
	project, err := projectService.GetByID(context.Background(), projectID)

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, project1, project)
	repo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	repo.On("GetByID", mock.Anything, mock.Anything).Return((*models.Project)(nil), errors.New("not found"))

	projectService := service.NewProjectService(repo)
	project, err := projectService.GetByID(context.Background(), "project-1")

	assert.Error(t, err)
	assert.Nil(t, project)
	repo.AssertExpectations(t)
}

func TestProjectUpdate(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	const projectID = "project-1"
	const projectName = "Project Name"
	newProjectName := "Project Name Upd"
	updTime := time.Now()

	project := &models.Project{
		ID:          projectID,
		Name:        projectName,
		Description: "Description 1",
		CreatedBy:   "user-1",
		UpdatedBy:   "user-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   updTime,
	}

	identity := auth.NewIdentity("user-ID", auth.UserRoleAdmin)

	repo.On("GetByID", mock.Anything, projectID).Return(project, nil)

	repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Project")).Return(nil)

	projectService := service.NewProjectService(repo)

	_, err := projectService.Update(context.Background(), identity, projectID, &newProjectName, nil)

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, projectID, project.ID)
	assert.Equal(t, newProjectName, project.Name)
	assert.Equal(t, identity.UserID, project.UpdatedBy)
	assert.NotEqual(t, identity.UserID, project.CreatedBy)
	assert.NotEqual(t, updTime, project.UpdatedAt)

	repo.AssertExpectations(t)
}

func TestProjectUpdate_NotFound(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	const projectID = "project-1"
	newProjectName := "Project Name Upd"

	identity := auth.NewIdentity("user-ID", auth.UserRoleAdmin)

	repo.On("GetByID", mock.Anything, projectID).Return((*models.Project)(nil), errors.New("not found"))

	projectService := service.NewProjectService(repo)

	project, err := projectService.Update(context.Background(), identity, projectID, &newProjectName, nil)

	assert.Error(t, err)
	assert.Nil(t, project)

	repo.AssertExpectations(t)

}

func TestProjectDelete(t *testing.T) {
	repo := new(tests.MockProjectRepository)

	const projectID = "project-1"
	project1 := &models.Project{
		ID:          projectID,
		Name:        "Project 1",
		Description: "Description 1",
		CreatedBy:   "user-1",
		UpdatedBy:   "user-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.On("GetByID", mock.Anything, projectID).Return(project1, nil)

	repo.On("Delete", mock.Anything, mock.AnythingOfType("*models.Project")).Return(nil)

	projectService := service.NewProjectService(repo)

	err := projectService.Delete(context.Background(), projectID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
