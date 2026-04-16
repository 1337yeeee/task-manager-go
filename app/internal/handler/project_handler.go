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

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: projectService}
}

// createProjectRequest represents project creation payload
type createProjectRequest struct {
	Name string `json:"name" binding:"required" example:"Task Manager"`
	Desc string `json:"desc" example:"Backend service for task management"`
}

// updateProjectRequest represents project update payload
type updateProjectRequest struct {
	Name *string `json:"name" binding:"omitempty" example:"Task Manager API"`
	Desc *string `json:"desc" binding:"omitempty" example:"Updated project description"`
}

// ProjectResponse represents project model in API docs
type ProjectResponse struct {
	ID        string    `json:"id" example:"7c3d9d0a-8dcb-4f7f-9b5e-2a12d4e9a001"`
	Name      string    `json:"name" example:"Task Manager"`
	Desc      string    `json:"desc" example:"Backend service for task management"`
	CreatedBy string    `json:"created_by" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	UpdatedBy string    `json:"updated_by" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-02T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-04-02T12:30:00Z"`
}

// ProjectDataResponse wraps single project response
type ProjectDataResponse struct {
	Data ProjectResponse `json:"data"`
}

// ProjectsDataResponse wraps project list response
type ProjectsDataResponse struct {
	Data []ProjectResponse `json:"data"`
}

// SuccessResponse is a generic success response
type SuccessResponse struct {
	Data bool `json:"data" example:"true"`
}

// GetAll godoc
// @Summary Get all projects
// @Description Returns all projects available to the authenticated user
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProjectsDataResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/ [get]
func (h *ProjectHandler) GetAll(ctx *gin.Context) {
	log.Println("GetAll Projects")
	projects, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in getAll"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": projects})
}

// Create godoc
// @Summary Create project
// @Description Creates a new project. Requires access token and moderator/editor/admin permissions according to middleware.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body createProjectRequest true "Project payload"
// @Success 200 {object} ProjectDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/ [post]
func (h *ProjectHandler) Create(ctx *gin.Context) {
	var req createProjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		responses.Unauthorized(ctx)
		return
	}

	project, err := h.service.Create(ctx, identity, req.Name, req.Desc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in create"})
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

// GetByID godoc
// @Summary Get project by ID
// @Description Returns a project by its identifier
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} ProjectDataResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	project, err := h.service.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in getByID"})
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

// Update godoc
// @Summary Update project
// @Description Updates a project. Current implementation expects project id in request body.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID from route"
// @Param input body updateProjectRequest true "Project update payload"
// @Success 200 {object} ProjectDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(ctx *gin.Context) {
	projectId := ctx.Param("id")
	var req updateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	identity := auth.FromContext(ctx.Request.Context())
	if identity == nil {
		responses.Unauthorized(ctx)
		return
	}

	project, err := h.service.Update(ctx, identity, projectId, req.Name, req.Desc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

// Delete godoc
// @Summary Delete project
// @Description Deletes a project by its identifier
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.service.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in delete"})
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true})
}
