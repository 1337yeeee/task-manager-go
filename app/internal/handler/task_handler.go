package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/responses"
	"task-manager/internal/service"
	"time"
)

type TaskHandler struct {
	taskService    service.TaskService
	projectService service.ProjectService
}

func NewTaskHandler(taskService service.TaskService, projectService service.ProjectService) *TaskHandler {
	return &TaskHandler{taskService: taskService, projectService: projectService}
}

type taskCreateRequest struct {
	Name    string `json:"name" binding:"required" example:"Implement JWT auth"`
	Content string `json:"content" binding:"required" example:"Add access and refresh token support"`
}

type taskUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty" example:"Implement JWT authentication"`
	Content     *string `json:"content" binding:"omitempty" example:"Update business logic and tests"`
	ExecutiveId *string `json:"executive_id" binding:"omitempty" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	AuditorId   *string `json:"auditor_id" binding:"omitempty" example:"a1f8c8d4-7b2a-4b10-bf18-1d3dcb6f9a77"`
}

type taskUpdateStatusRequest struct {
	Status string `json:"status" binding:"required" example:"in_progress"`
}

type TaskResponse struct {
	ID          string    `json:"id" example:"7c3d9d0a-8dcb-4f7f-9b5e-2a12d4e9a010"`
	ProjectID   string    `json:"project_id" example:"7c3d9d0a-8dcb-4f7f-9b5e-2a12d4e9a001"`
	Name        string    `json:"name" example:"Implement JWT auth"`
	Content     string    `json:"content" example:"Add access and refresh token support"`
	ExecutiveID string    `json:"executive_id" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	AuditorID   string    `json:"auditor_id" example:"a1f8c8d4-7b2a-4b10-bf18-1d3dcb6f9a77"`
	CreatedBy   string    `json:"created_by" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	UpdatedBy   string    `json:"updated_by" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	CreatedAt   time.Time `json:"created_at" example:"2026-04-02T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-04-02T12:30:00Z"`
	Status      string    `json:"status" example:"created"`
}

type TaskDataResponse struct {
	Data TaskResponse `json:"data"`
}

type TasksDataResponse struct {
	Data []TaskResponse `json:"data"`
}

type TasksListResponse struct {
	Tasks []TaskResponse `json:"data"`
}

type OkResponse struct {
	Ok bool `json:"ok" example:"true"`
}

// GetByProject godoc
// @Summary Get tasks by project
// @Description Returns all tasks for a specific project
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} TasksListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id}/tasks [get]
func (h TaskHandler) GetByProject(ctx *gin.Context) {
	projectId := ctx.Param("id")
	project, err := h.projectService.GetByID(ctx, projectId)
	if err != nil {
		log.Println(err)
		responses.NotFound(ctx)
		return
	}

	tasks, err := h.taskService.GetByProjectID(ctx, project.ID)
	if err != nil {
		log.Println(err)
		responses.NotFound(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

// Create godoc
// @Summary Create task in project
// @Description Creates a task for the specified project
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param input body taskCreateRequest true "Task create payload"
// @Success 200 {object} TaskDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id}/tasks [post]
func (h TaskHandler) Create(ctx *gin.Context) {
	projectId := ctx.Param("id")
	project, err := h.projectService.GetByID(ctx, projectId)
	if err != nil {
		log.Println(err)
		responses.NotFound(ctx)
		return
	}

	var req taskCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		responses.Unauthorized(ctx)
		return
	}

	task, err := h.taskService.Create(ctx, identity, project.ID, req.Name, req.Content)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in create"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

// GetByID godoc
// @Summary Get task by ID
// @Description Returns task by identifier
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} TaskDataResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [get]
func (h TaskHandler) GetByID(ctx *gin.Context) {
	taskId := ctx.Param("id")
	task, err := h.taskService.GetByID(ctx, taskId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in get"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

// Update godoc
// @Summary Update task
// @Description Updates task fields by identifier
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param input body taskUpdateRequest true "Task update payload"
// @Success 200 {object} OkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [put]
func (h TaskHandler) Update(ctx *gin.Context) {
	req := taskUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		responses.Unauthorized(ctx)
		return
	}

	err := h.taskService.Update(ctx, identity, id, req.Name, req.Content, req.ExecutiveId, req.AuditorId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// UpdateStatus godoc
// @Summary Update task status
// @Description Updates task status by identifier
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param input body taskUpdateStatusRequest true "Task status payload"
// @Success 200 {object} OkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id}/status [patch]
func (h TaskHandler) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req taskUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		responses.Unauthorized(ctx)
		return
	}

	err := h.taskService.UpdateStatus(ctx, identity, id, req.Status)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

// Delete godoc
// @Summary Delete task
// @Description Deletes task by identifier
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} OkResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h TaskHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.taskService.Delete(ctx, id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in delete"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
