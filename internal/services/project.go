package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) (*entities.Project, error)
	GetOne(ctx context.Context, id uint) (*entities.Project, error)
	GetMany(ctx context.Context, qParams *entities.ProjectQueryParams) ([]*entities.Project, int64, error)
	Update(ctx context.Context, project *entities.Project) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// ProjectService handles project business logic
type ProjectService struct {
	repo ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, project *entities.Project) (*entities.Project, error) {
	return s.repo.Create(ctx, project)
}

// GetProject retrieves a single project by ID
func (s *ProjectService) GetProject(ctx context.Context, id uint) (*entities.Project, error) {
	return s.repo.GetOne(ctx, id)
}

// GetProjects retrieves multiple projects with optional query parameters
func (s *ProjectService) GetProjects(ctx context.Context, params *entities.ProjectQueryParams) (*entities.ProjectListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, project *entities.Project) (int64, error) {
	return s.repo.Update(ctx, project)
}

// DeleteProject deletes a project by ID
func (s *ProjectService) DeleteProject(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
