package usecases

import (
	"context"

	"github.com/ducminhgd/plan-craft/internal/entities"
)

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) (*entities.Client, error)
	GetOne(ctx context.Context, id uint) (*entities.Client, error)
	GetMany(ctx context.Context, qParams *entities.ClientQueryParams) ([]*entities.Client, int64, error)
	Update(ctx context.Context, client *entities.Client) (int64, error)
	Delete(ctx context.Context, id uint) error
}

// ClientUseCase is the use case for client entities
type ClientUseCase struct {
	repo ClientRepository
}

// NewClientUseCase creates a new client use case
func NewClientUseCase(repo ClientRepository) *ClientUseCase {
	return &ClientUseCase{repo: repo}
}

// Create creates a new client
func (uc *ClientUseCase) CreateAClient(ctx context.Context, client *entities.Client) (*entities.Client, error) {
	return uc.repo.Create(ctx, client)
}

// GetOne gets a client by ID
func (uc *ClientUseCase) GetAClient(ctx context.Context, id uint) (*entities.Client, error) {
	return uc.repo.GetOne(ctx, id)
}

// GetMany gets multiple clients by query parameters
func (uc *ClientUseCase) GetManyClients(ctx context.Context, qParams *entities.ClientQueryParams) ([]*entities.Client, int64, error) {
	return uc.repo.GetMany(ctx, qParams)
}

// Update updates a client
func (uc *ClientUseCase) UpdateAClient(ctx context.Context, client *entities.Client) (int64, error) {
	return uc.repo.Update(ctx, client)
}

// Delete deletes a client
func (uc *ClientUseCase) DeleteAClient(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
