package integration

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"task-manager/internal/app"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/handler"
	"task-manager/internal/routes"
	"task-manager/internal/utils"
	"task-manager/tests"
)

type stubAuthService struct{}

func (s stubAuthService) Login(_ context.Context, _ string, _ string) (string, string, auth.UserRole, error) {
	return "access", "refresh", auth.UserRoleViewer, nil
}

func (s stubAuthService) RefreshToken(_ context.Context, identity *auth.Identity, _ string) (string, string, auth.UserRole, error) {
	return "access-next", "refresh-next", identity.Role, nil
}

func (s stubAuthService) Logout(_ context.Context, _ *auth.Identity) error {
	return nil
}

type stubUserService struct{}

func (s stubUserService) Register(_ context.Context, name string, email string, _ string, role *auth.UserRole) (*models.User, error) {
	resultRole := auth.UserRoleViewer
	if role != nil && role.IsValid() {
		resultRole = *role
	}

	return &models.User{
		ID:       "user-created",
		Name:     name,
		Email:    email,
		Role:     resultRole,
		IsActive: true,
	}, nil
}

func (s stubUserService) GetAll(_ context.Context, _ models.UserFilter) ([]models.User, error) {
	return []models.User{
		{
			ID:       "user-1",
			Name:     "User 1",
			Email:    "user1@example.com",
			Role:     auth.UserRoleViewer,
			IsActive: true,
		},
	}, nil
}

func (s stubUserService) GetById(_ context.Context, id string) (*models.User, error) {
	return &models.User{
		ID:       id,
		Name:     "User",
		Email:    "user@example.com",
		Role:     auth.UserRoleViewer,
		IsActive: true,
	}, nil
}

func (s stubUserService) Update(_ context.Context, id string, name *string, email *string, _ *string, role *auth.UserRole, isActive *bool) (*models.User, error) {
	user := &models.User{
		ID:       id,
		Name:     "User",
		Email:    "user@example.com",
		Role:     auth.UserRoleViewer,
		IsActive: true,
	}

	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if role != nil && role.IsValid() {
		user.Role = *role
	}
	if isActive != nil {
		user.IsActive = *isActive
	}

	return user, nil
}

func (s stubUserService) Delete(_ context.Context, _ string) error {
	return nil
}

type stubProjectService struct{}

func (s stubProjectService) GetAll(_ context.Context) ([]models.Project, error) {
	return []models.Project{
		{
			ID:          "project-1",
			Name:        "Project 1",
			Description: "Description",
		},
	}, nil
}

func (s stubProjectService) Create(_ context.Context, identity *auth.Identity, name string, desc string) (*models.Project, error) {
	now := time.Now()
	return &models.Project{
		ID:          "project-created",
		Name:        name,
		Description: desc,
		CreatedBy:   identity.UserID,
		UpdatedBy:   identity.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s stubProjectService) GetByID(_ context.Context, id string) (*models.Project, error) {
	return &models.Project{
		ID:          id,
		Name:        "Project",
		Description: "Description",
	}, nil
}

func (s stubProjectService) Update(_ context.Context, identity *auth.Identity, id string, name *string, desc *string) (*models.Project, error) {
	project := &models.Project{
		ID:          id,
		Name:        "Project",
		Description: "Description",
		UpdatedBy:   identity.UserID,
	}
	if name != nil {
		project.Name = *name
	}
	if desc != nil {
		project.Description = *desc
	}
	return project, nil
}

func (s stubProjectService) Delete(_ context.Context, _ string) error {
	return nil
}

type stubTaskService struct{}

func (s stubTaskService) Create(_ context.Context, identity *auth.Identity, projectID string, name string, content string) (*models.Task, error) {
	now := time.Now()
	return &models.Task{
		ID:          "task-created",
		ProjectID:   projectID,
		Name:        name,
		Content:     content,
		Status:      models.TaskStatusCreated,
		ExecutiveID: identity.UserID,
		AuditorID:   identity.UserID,
		CreatedBy:   identity.UserID,
		UpdatedBy:   identity.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s stubTaskService) GetByID(_ context.Context, id string) (*models.Task, error) {
	return &models.Task{
		ID:     id,
		Status: models.TaskStatusCreated,
	}, nil
}

func (s stubTaskService) GetByProjectID(_ context.Context, projectID string) ([]models.Task, error) {
	return []models.Task{
		{
			ID:        "task-1",
			ProjectID: projectID,
			Name:      "Task",
			Status:    models.TaskStatusCreated,
		},
	}, nil
}

func (s stubTaskService) Update(_ context.Context, _ *auth.Identity, _ string, _ *string, _ *string, _ *string, _ *string) error {
	return nil
}

func (s stubTaskService) UpdateStatus(_ context.Context, _ *auth.Identity, _ string, _ string) error {
	return nil
}

func (s stubTaskService) Delete(_ context.Context, _ string) error {
	return nil
}

func buildRouterAndToken(t *testing.T, role auth.UserRole, isActive bool) (http.Handler, string, *tests.MockUserRepository) {
	t.Helper()

	tokenManager := utils.NewTokenManager("test-secret", utils.DefaultAccessTTL, utils.DefaultRefreshTTL)
	accessToken, err := tokenManager.GenerateAccessToken("test-user-id", role)
	require.NoError(t, err)

	userRepo := &tests.MockUserRepository{}
	userRepo.
		On("FindUserByID", mock.Anything, "test-user-id").
		Return(&models.User{
			ID:       "test-user-id",
			Role:     role,
			IsActive: isActive,
		}, nil)

	projectService := stubProjectService{}

	container := &app.Container{
		TokenManager:   tokenManager,
		UserRepository: userRepo,
		AuthHandler:    handler.NewAuthHandler(stubAuthService{}),
		UserHandler:    handler.NewUserHandler(stubUserService{}),
		ProjectHandler: handler.NewProjectHandler(projectService),
		TaskHandler:    handler.NewTaskHandler(stubTaskService{}, projectService),
	}

	return routes.SetupRouter(container), accessToken, userRepo
}

func doAuthorizedRequest(router http.Handler, accessToken string, method string, path string, body string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body == "" {
		bodyReader = bytes.NewReader(nil)
	} else {
		bodyReader = bytes.NewReader([]byte(body))
	}

	req, _ := http.NewRequest(method, path, bodyReader)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestRoutePermissions_Matrix(t *testing.T) {
	type testCase struct {
		name   string
		method string
		path   string
		body   string

		adminExpected  int
		editorExpected int
		viewerExpected int
	}

	cases := []testCase{
		{
			name:           "projects list доступен всем",
			method:         http.MethodGet,
			path:           "/api/projects",
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusOK,
			viewerExpected: http.StatusOK,
		},
		{
			name:           "users list только модераторы",
			method:         http.MethodGet,
			path:           "/api/users",
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusOK,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "register user только admin",
			method:         http.MethodPost,
			path:           "/api/users",
			body:           `{"name":"New User","email":"new@example.com","password":"password123","role":"viewer"}`,
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusForbidden,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "get user by id только admin",
			method:         http.MethodGet,
			path:           "/api/users/user-1",
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusForbidden,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "update user только admin",
			method:         http.MethodPut,
			path:           "/api/users/user-1",
			body:           `{"role":"editor","is_active":true}`,
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusForbidden,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "delete user только admin",
			method:         http.MethodDelete,
			path:           "/api/users/user-1",
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusForbidden,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "create project только модераторы",
			method:         http.MethodPost,
			path:           "/api/projects",
			body:           `{"name":"Project","description":"Desc"}`,
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusOK,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "create task только модераторы",
			method:         http.MethodPost,
			path:           "/api/projects/project-1/tasks",
			body:           `{"name":"Task","content":"Content"}`,
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusOK,
			viewerExpected: http.StatusForbidden,
		},
		{
			name:           "update task status только модераторы",
			method:         http.MethodPatch,
			path:           "/api/tasks/task-1/status",
			body:           `{"status":"done"}`,
			adminExpected:  http.StatusOK,
			editorExpected: http.StatusOK,
			viewerExpected: http.StatusForbidden,
		},
	}

	roles := []struct {
		role        auth.UserRole
		expectedFor func(tc testCase) int
	}{
		{role: auth.UserRoleAdmin, expectedFor: func(tc testCase) int { return tc.adminExpected }},
		{role: auth.UserRoleEditor, expectedFor: func(tc testCase) int { return tc.editorExpected }},
		{role: auth.UserRoleViewer, expectedFor: func(tc testCase) int { return tc.viewerExpected }},
	}

	for _, roleCase := range roles {
		t.Run("role="+string(roleCase.role), func(t *testing.T) {
			router, token, userRepo := buildRouterAndToken(t, roleCase.role, true)

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					w := doAuthorizedRequest(router, token, tc.method, tc.path, tc.body)
					assert.Equal(t, roleCase.expectedFor(tc), w.Code)
				})
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestRoutePermissions_InactiveUserDenied(t *testing.T) {
	router, token, userRepo := buildRouterAndToken(t, auth.UserRoleAdmin, false)

	w := doAuthorizedRequest(router, token, http.MethodGet, "/api/projects", "")

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "user is inactive")
	userRepo.AssertExpectations(t)
}
