package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// ProjectRoleRepository defines the interface for project role data operations
type ProjectRoleRepository interface {
	Create(ctx context.Context, projectRole *entities.ProjectRole) (*entities.ProjectRole, error)
	GetOne(ctx context.Context, id uint) (*entities.ProjectRole, error)
	GetMany(ctx context.Context, qParams *entities.ProjectRoleQueryParams) ([]*entities.ProjectRole, int64, error)
	Update(ctx context.Context, projectRole *entities.ProjectRole) (int64, error)
	Delete(ctx context.Context, id uint) error
	GetByProjectNameAndLevel(ctx context.Context, projectID uint, name string, level uint) (*entities.ProjectRole, error)
}

// ProjectRoleService handles project role business logic
type ProjectRoleService struct {
	repo ProjectRoleRepository
}

// NewProjectRoleService creates a new project role service
func NewProjectRoleService(repo ProjectRoleRepository) *ProjectRoleService {
	return &ProjectRoleService{repo: repo}
}

// CreateProjectRole creates a new project role
func (s *ProjectRoleService) CreateProjectRole(ctx context.Context, projectRole *entities.ProjectRole) (*entities.ProjectRole, error) {
	return s.repo.Create(ctx, projectRole)
}

// GetProjectRole retrieves a single project role by ID
func (s *ProjectRoleService) GetProjectRole(ctx context.Context, id uint) (*entities.ProjectRole, error) {
	return s.repo.GetOne(ctx, id)
}

// GetProjectRoles retrieves multiple project roles with optional query parameters
func (s *ProjectRoleService) GetProjectRoles(ctx context.Context, params *entities.ProjectRoleQueryParams) (*entities.ProjectRoleListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectRoleListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateProjectRole updates an existing project role
func (s *ProjectRoleService) UpdateProjectRole(ctx context.Context, projectRole *entities.ProjectRole) (int64, error) {
	return s.repo.Update(ctx, projectRole)
}

// DeleteProjectRole deletes a project role by ID
func (s *ProjectRoleService) DeleteProjectRole(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetByProjectNameAndLevel retrieves a project role by project ID, name, and level
func (s *ProjectRoleService) GetByProjectNameAndLevel(ctx context.Context, projectID uint, name string, level uint) (*entities.ProjectRole, error) {
	return s.repo.GetByProjectNameAndLevel(ctx, projectID, name, level)
}

// GetProjectRolesByProject retrieves all project roles for a specific project
func (s *ProjectRoleService) GetProjectRolesByProject(ctx context.Context, projectID uint) (*entities.ProjectRoleListResponse, error) {
	params := &entities.ProjectRoleQueryParams{
		ProjectID: projectID,
	}
	return s.GetProjectRoles(ctx, params)
}
