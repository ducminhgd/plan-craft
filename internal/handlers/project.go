package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// ProjectHandler handles project-related operations for Wails bindings
type ProjectHandler struct {
	ctx     context.Context
	service *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(ctx context.Context, service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetProjects retrieves multiple projects with optional query parameters
func (h *ProjectHandler) GetProjects(params *entities.ProjectQueryParams) (*entities.ProjectListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project service not initialized")
	}
	return h.service.GetProjects(h.ctx, params)
}

// GetProject retrieves a single project by ID
func (h *ProjectHandler) GetProject(id uint) (*entities.Project, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project service not initialized")
	}
	return h.service.GetProject(h.ctx, id)
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(project *entities.Project) (*entities.Project, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project service not initialized")
	}
	return h.service.CreateProject(h.ctx, project)
}

// UpdateProject updates an existing project
func (h *ProjectHandler) UpdateProject(project *entities.Project) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("project service not initialized")
	}
	return h.service.UpdateProject(h.ctx, project)
}

// DeleteProject deletes a project by ID
func (h *ProjectHandler) DeleteProject(id uint) error {
	if h.service == nil {
		return fmt.Errorf("project service not initialized")
	}
	return h.service.DeleteProject(h.ctx, id)
}
