package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// ProjectRoleHandler handles project role-related operations for Wails bindings
type ProjectRoleHandler struct {
	ctx     context.Context
	service *services.ProjectRoleService
}

// NewProjectRoleHandler creates a new ProjectRoleHandler
func NewProjectRoleHandler(ctx context.Context, service *services.ProjectRoleService) *ProjectRoleHandler {
	return &ProjectRoleHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetProjectRoles retrieves multiple project roles with optional query parameters
func (h *ProjectRoleHandler) GetProjectRoles(params *entities.ProjectRoleQueryParams) (*entities.ProjectRoleListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project role service not initialized")
	}
	return h.service.GetProjectRoles(h.ctx, params)
}

// GetProjectRole retrieves a single project role by ID
func (h *ProjectRoleHandler) GetProjectRole(id uint) (*entities.ProjectRole, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project role service not initialized")
	}
	return h.service.GetProjectRole(h.ctx, id)
}

// GetProjectRolesByProject retrieves all project roles for a specific project
func (h *ProjectRoleHandler) GetProjectRolesByProject(projectID uint) (*entities.ProjectRoleListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project role service not initialized")
	}
	return h.service.GetProjectRolesByProject(h.ctx, projectID)
}

// CreateProjectRole creates a new project role
func (h *ProjectRoleHandler) CreateProjectRole(projectRole *entities.ProjectRole) (*entities.ProjectRole, error) {
	if h.service == nil {
		return nil, fmt.Errorf("project role service not initialized")
	}
	return h.service.CreateProjectRole(h.ctx, projectRole)
}

// UpdateProjectRole updates an existing project role
func (h *ProjectRoleHandler) UpdateProjectRole(projectRole *entities.ProjectRole) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("project role service not initialized")
	}
	return h.service.UpdateProjectRole(h.ctx, projectRole)
}

// DeleteProjectRole deletes a project role by ID
func (h *ProjectRoleHandler) DeleteProjectRole(id uint) error {
	if h.service == nil {
		return fmt.Errorf("project role service not initialized")
	}
	return h.service.DeleteProjectRole(h.ctx, id)
}
