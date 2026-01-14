package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// HumanResourceRepository defines the interface for human resource data operations
type HumanResourceRepository interface {
	Create(ctx context.Context, humanResource *entities.HumanResource) (*entities.HumanResource, error)
	GetOne(ctx context.Context, id uint) (*entities.HumanResource, error)
	GetMany(ctx context.Context, qParams *entities.HumanResourceQueryParams) ([]*entities.HumanResource, int64, error)
	Update(ctx context.Context, humanResource *entities.HumanResource) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// HumanResourceService handles human resource business logic
type HumanResourceService struct {
	repo HumanResourceRepository
}

// NewHumanResourceService creates a new human resource service
func NewHumanResourceService(repo HumanResourceRepository) *HumanResourceService {
	return &HumanResourceService{repo: repo}
}

// CreateHumanResource creates a new human resource
func (s *HumanResourceService) CreateHumanResource(ctx context.Context, humanResource *entities.HumanResource) (*entities.HumanResource, error) {
	return s.repo.Create(ctx, humanResource)
}

// GetHumanResource retrieves a single human resource by ID
func (s *HumanResourceService) GetHumanResource(ctx context.Context, id uint) (*entities.HumanResource, error) {
	return s.repo.GetOne(ctx, id)
}

// GetHumanResources retrieves multiple human resources with optional query parameters
func (s *HumanResourceService) GetHumanResources(ctx context.Context, params *entities.HumanResourceQueryParams) (*entities.HumanResourceListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.HumanResourceListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateHumanResource updates an existing human resource
func (s *HumanResourceService) UpdateHumanResource(ctx context.Context, humanResource *entities.HumanResource) (int64, error) {
	return s.repo.Update(ctx, humanResource)
}

// DeleteHumanResource deletes a human resource by ID
func (s *HumanResourceService) DeleteHumanResource(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
