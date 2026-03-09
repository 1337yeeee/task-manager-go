package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"task-manager/internal/service"
)

type ProjectHandler struct {
	service service.ProjectService
}

type createProjectRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type updateProjectRequest struct {
	ID   string  `json:"id" binding:"required"`
	Name *string `json:"name" binding:"omitempty"`
	Desc *string `json:"desc" binding:"omitempty"`
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: projectService}
}

func (h *ProjectHandler) GetAll(ctx *gin.Context) {
	projects, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in getAll"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": projects})
}

func (h *ProjectHandler) Create(ctx *gin.Context) {
	var req createProjectRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	project, err := h.service.Create(ctx, req.Name, req.Desc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in create"})
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

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

func (h *ProjectHandler) Update(ctx *gin.Context) {
	var req updateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	project, err := h.service.Update(ctx, req.ID, req.Name, req.Desc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in update"})
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": project})
}

func (h *ProjectHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.service.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error in delete"})
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true})
}
