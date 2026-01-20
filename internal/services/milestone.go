package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// MilestoneRepository defines the interface for milestone data operations
type MilestoneRepository interface {
	Create(ctx context.Context, milestone *entities.Milestone) (*entities.Milestone, error)
	GetOne(ctx context.Context, id uint) (*entities.Milestone, error)
	GetMany(ctx context.Context, qParams *entities.MilestoneQueryParams) ([]*entities.Milestone, int64, error)
	Update(ctx context.Context, milestone *entities.Milestone) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// MilestoneService handles milestone business logic
type MilestoneService struct {
	repo MilestoneRepository
}

// NewMilestoneService creates a new milestone service
func NewMilestoneService(repo MilestoneRepository) *MilestoneService {
	return &MilestoneService{repo: repo}
}

// CreateMilestone creates a new milestone
func (s *MilestoneService) CreateMilestone(ctx context.Context, milestone *entities.Milestone) (*entities.Milestone, error) {
	return s.repo.Create(ctx, milestone)
}

// GetMilestone retrieves a single milestone by ID
func (s *MilestoneService) GetMilestone(ctx context.Context, id uint) (*entities.Milestone, error) {
	return s.repo.GetOne(ctx, id)
}

// GetMilestones retrieves multiple milestones with optional query parameters
func (s *MilestoneService) GetMilestones(ctx context.Context, params *entities.MilestoneQueryParams) (*entities.MilestoneListResponse, error) {
	data, total, err := s.repo.GetMany(ctx, params)
	if err != nil {
		return nil, err
	}
	return &entities.MilestoneListResponse{
		Data:  data,
		Total: total,
	}, nil
}

// UpdateMilestone updates an existing milestone
func (s *MilestoneService) UpdateMilestone(ctx context.Context, milestone *entities.Milestone) (int64, error) {
	return s.repo.Update(ctx, milestone)
}

// DeleteMilestone deletes a milestone by ID
func (s *MilestoneService) DeleteMilestone(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
