package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// HumanResourceHandler handles human resource-related operations for Wails bindings
type HumanResourceHandler struct {
	ctx     context.Context
	service *services.HumanResourceService
}

// NewHumanResourceHandler creates a new HumanResourceHandler
func NewHumanResourceHandler(ctx context.Context, service *services.HumanResourceService) *HumanResourceHandler {
	return &HumanResourceHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetHumanResources retrieves multiple human resources with optional query parameters
func (h *HumanResourceHandler) GetHumanResources(params *entities.HumanResourceQueryParams) (*entities.HumanResourceListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return h.service.GetHumanResources(h.ctx, params)
}

// GetHumanResource retrieves a single human resource by ID
func (h *HumanResourceHandler) GetHumanResource(id uint) (*entities.HumanResource, error) {
	if h.service == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return h.service.GetHumanResource(h.ctx, id)
}

// CreateHumanResource creates a new human resource
func (h *HumanResourceHandler) CreateHumanResource(humanResource *entities.HumanResource) (*entities.HumanResource, error) {
	if h.service == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return h.service.CreateHumanResource(h.ctx, humanResource)
}

// UpdateHumanResource updates an existing human resource
func (h *HumanResourceHandler) UpdateHumanResource(humanResource *entities.HumanResource) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("human resource service not initialized")
	}
	return h.service.UpdateHumanResource(h.ctx, humanResource)
}

// DeleteHumanResource deletes a human resource by ID
func (h *HumanResourceHandler) DeleteHumanResource(id uint) error {
	if h.service == nil {
		return fmt.Errorf("human resource service not initialized")
	}
	return h.service.DeleteHumanResource(h.ctx, id)
}
