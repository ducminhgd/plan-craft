package handlers

import (
	"context"
	"fmt"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// ClientHandler handles client-related operations for Wails bindings
type ClientHandler struct {
	ctx     context.Context
	service *services.ClientService
}

// NewClientHandler creates a new ClientHandler
func NewClientHandler(ctx context.Context, service *services.ClientService) *ClientHandler {
	return &ClientHandler{
		ctx:     ctx,
		service: service,
	}
}

// GetClients retrieves multiple clients with optional query parameters
func (h *ClientHandler) GetClients(params *entities.ClientQueryParams) (*entities.ClientListResponse, error) {
	if h.service == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return h.service.GetClients(h.ctx, params)
}

// GetClient retrieves a single client by ID
func (h *ClientHandler) GetClient(id uint) (*entities.Client, error) {
	if h.service == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return h.service.GetClient(h.ctx, id)
}

// CreateClient creates a new client
func (h *ClientHandler) CreateClient(client *entities.Client) (*entities.Client, error) {
	if h.service == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return h.service.CreateClient(h.ctx, client)
}

// UpdateClient updates an existing client
func (h *ClientHandler) UpdateClient(client *entities.Client) (int64, error) {
	if h.service == nil {
		return 0, fmt.Errorf("client service not initialized")
	}
	return h.service.UpdateClient(h.ctx, client)
}

// DeleteClient deletes a client by ID
func (h *ClientHandler) DeleteClient(id uint) error {
	if h.service == nil {
		return fmt.Errorf("client service not initialized")
	}
	return h.service.DeleteClient(h.ctx, id)
}
