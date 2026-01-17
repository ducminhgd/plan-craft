package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// ProjectResourceRepository defines the interface for project resource data operations
type ProjectResourceRepository interface {
	Create(ctx context.Context, projectResource *entities.ProjectResource) (*entities.ProjectResource, error)
	GetOne(ctx context.Context, id uint) (*entities.ProjectResource, error)
	GetMany(ctx context.Context, qParams *entities.ProjectResourceQueryParams) ([]*entities.ProjectResource, int64, error)
	Update(ctx context.Context, projectResource *entities.ProjectResource) (int64, error)
	Delete(ctx context.Context, id uint) error
	GetByProjectAndResource(ctx context.Context, projectID, humanResourceID uint) (*entities.ProjectResource, error)
}

// ProjectResourceService handles project resource business logic
type ProjectResourceService struct {
	repo ProjectResourceRepository
}

// NewProjectResourceService creates a new project resource service
func NewProjectResourceService(repo ProjectResourceRepository) *ProjectResourceService {
	return &ProjectResourceService{repo: repo}
}

// CreateProjectResource creates a new project resource allocation
func (s *ProjectResourceService) CreateProjectResource(ctx context.Context, projectResource *entities.ProjectResource) (*entities.ProjectResource, error) {
	return s.repo.Create(ctx, projectResource)
}

// GetProjectResource retrieves a single project resource by ID
func (s *ProjectResourceService) GetProjectResource(ctx context.Context, id uint) (*entities.ProjectResource, error) {
	return s.repo.GetOne(ctx, id)
}

// GetProjectResources retrieves multiple project resources with optional query parameters
func (s *ProjectResourceService) GetProjectResources(ctx context.Context, params *entities.ProjectResourceQueryParams) (*entities.ProjectResourceListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectResourceListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateProjectResource updates an existing project resource allocation
func (s *ProjectResourceService) UpdateProjectResource(ctx context.Context, projectResource *entities.ProjectResource) (int64, error) {
	return s.repo.Update(ctx, projectResource)
}

// DeleteProjectResource deletes a project resource allocation by ID
func (s *ProjectResourceService) DeleteProjectResource(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetByProjectAndResource retrieves a project resource by project ID and human resource ID
func (s *ProjectResourceService) GetByProjectAndResource(ctx context.Context, projectID, humanResourceID uint) (*entities.ProjectResource, error) {
	return s.repo.GetByProjectAndResource(ctx, projectID, humanResourceID)
}
