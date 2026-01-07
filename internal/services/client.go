package services

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

// ClientRepository defines the interface for client data operations
type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)
	GetOne(ctx context.Context, id uint) (*entities.Client, error)
	GetMany(ctx context.Context, qParams *entities.ClientQueryParams) ([]*entities.Client, int64, error)
	Update(ctx context.Context, client *entities.Client) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// ClientService handles client business logic
type ClientService struct {
	repo ClientRepository
}

// NewClientService creates a new client service
func NewClientService(repo ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

// CreateClient creates a new client
func (s *ClientService) CreateClient(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	return s.repo.Create(ctx, client)
}

// GetClient retrieves a single client by ID
func (s *ClientService) GetClient(ctx context.Context, id uint) (*entities.Client, error) {
	return s.repo.GetOne(ctx, id)
}

// GetClients retrieves multiple clients with optional query parameters
func (s *ClientService) GetClients(ctx context.Context, params *entities.ClientQueryParams) ([]*entities.Client, int64, error) {
	return s.repo.GetMany(ctx, params)
}

// UpdateClient updates an existing client
func (s *ClientService) UpdateClient(ctx context.Context, client *entities.Client) (int64, error) {
	return s.repo.Update(ctx, client)
}

// DeleteClient deletes a client by ID
func (s *ClientService) DeleteClient(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
