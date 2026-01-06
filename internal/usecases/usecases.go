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
