package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// ProjectResourceHandler handles project resource-related operations for Wails bindings
type ProjectResourceHandler struct {
	ctx     context.Context
	service *services.ProjectResourceService
}

// NewProjectResourceHandler creates a new ProjectResourceHandler
func NewProjectResourceHandler(ctx context.Context, service *services.ProjectResourceService) *ProjectResourceHandler {
	return &ProjectResourceHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetProjectResources retrieves multiple project resources with optional query parameters
func (h *ProjectResourceHandler) GetProjectResources(params *entities.ProjectResourceQueryParams) (*entities.ProjectResourceListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project resource service not initialized")
	}
	return h.service.GetProjectResources(h.ctx, params)
}

// GetProjectResource retrieves a single project resource by ID
func (h *ProjectResourceHandler) GetProjectResource(id uint) (*entities.ProjectResource, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project resource service not initialized")
	}
	return h.service.GetProjectResource(h.ctx, id)
}

// CreateProjectResource creates a new project resource allocation
func (h *ProjectResourceHandler) CreateProjectResource(projectResource *entities.ProjectResource) (*entities.ProjectResource, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project resource service not initialized")
	}
	return h.service.CreateProjectResource(h.ctx, projectResource)
}

// UpdateProjectResource updates an existing project resource allocation
func (h *ProjectResourceHandler) UpdateProjectResource(projectResource *entities.ProjectResource) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("project resource service not initialized")
	}
	return h.service.UpdateProjectResource(h.ctx, projectResource)
}

// DeleteProjectResource deletes a project resource allocation by ID
func (h *ProjectResourceHandler) DeleteProjectResource(id uint) error {
	if h.service == nil {
		return fmt.Errorf("project resource service not initialized")
	}
	return h.service.DeleteProjectResource(h.ctx, id)
}

// GetByProjectAndResource retrieves a project resource by project ID and human resource ID
func (h *ProjectResourceHandler) GetByProjectAndResource(projectID, humanResourceID uint) (*entities.ProjectResource, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project resource service not initialized")
	}
	return h.service.GetByProjectAndResource(h.ctx, projectID, humanResourceID)
}
