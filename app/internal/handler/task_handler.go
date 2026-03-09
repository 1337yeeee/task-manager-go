package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task-manager/internal/responses"
	"task-manager/internal/service"
)

type TaskHandler struct {
	taskService    service.TaskService
	projectService service.ProjectService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: service}
}

type taskCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type taskUpdateRequest struct {
	ID          string  `json:"id" binding:"required"`
	Name        *string `json:"name" binding:"omitempty"`
	Content     *string `json:"content" binding:"omitempty"`
	ExecutiveId *string `json:"executive_id" binding:"omitempty"`
	AuditorId   *string `json:"auditor_id" binding:"omitempty"`
}

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
		"tasks": tasks,
	})
}

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

	task, err := h.taskService.Create(ctx, project.ID, req.Name, req.Content)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in create"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

func (h TaskHandler) GetByID(ctx *gin.Context) {
	taskId := ctx.Param("task_id")
	task, err := h.taskService.GetByID(ctx, taskId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in get"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

func (h TaskHandler) Update(ctx *gin.Context) {
	req := taskUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.taskService.Update(ctx, req.ID, req.Name, req.Content, req.ExecutiveId, req.AuditorId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h TaskHandler) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Param("status")

	err := h.taskService.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h TaskHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.taskService.Delete(ctx, id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in delete"})
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
