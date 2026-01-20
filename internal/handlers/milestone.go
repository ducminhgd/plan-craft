package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// MilestoneHandler handles milestone-related operations for Wails bindings
type MilestoneHandler struct {
	ctx     context.Context
	service *services.MilestoneService
}

// NewMilestoneHandler creates a new MilestoneHandler
func NewMilestoneHandler(ctx context.Context, service *services.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetMilestones retrieves multiple milestones with optional query parameters
func (h *MilestoneHandler) GetMilestones(params *entities.MilestoneQueryParams) (*entities.MilestoneListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("milestone service not initialized")
	}
	return h.service.GetMilestones(h.ctx, params)
}

// GetMilestone retrieves a single milestone by ID
func (h *MilestoneHandler) GetMilestone(id uint) (*entities.Milestone, error) {
	if h.service == nil {
		return nil, fmt.Errorf("milestone service not initialized")
	}
	return h.service.GetMilestone(h.ctx, id)
}

// CreateMilestone creates a new milestone
func (h *MilestoneHandler) CreateMilestone(milestone *entities.Milestone) (*entities.Milestone, error) {
	if h.service == nil {
		return nil, fmt.Errorf("milestone service not initialized")
	}
	return h.service.CreateMilestone(h.ctx, milestone)
}

// UpdateMilestone updates an existing milestone
func (h *MilestoneHandler) UpdateMilestone(milestone *entities.Milestone) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("milestone service not initialized")
	}
	return h.service.UpdateMilestone(h.ctx, milestone)
}

// DeleteMilestone deletes a milestone by ID
func (h *MilestoneHandler) DeleteMilestone(id uint) error {
	if h.service == nil {
		return fmt.Errorf("milestone service not initialized")
	}
	return h.service.DeleteMilestone(h.ctx, id)
}
