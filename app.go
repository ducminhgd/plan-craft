package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/ducminhgd/plan-craft/internal/infrastructures"
	"github.com/ducminhgd/plan-craft/internal/repositories"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// App struct
type App struct {
	ctx                  context.Context
	clientService        *services.ClientService
	humanResourceService *services.HumanResourceService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	db, err := infrastructures.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Wire dependencies: repository → service
	clientRepo := repositories.NewClientRepository(db)
	a.clientService = services.NewClientService(clientRepo)

	hrRepo := repositories.NewHRRepository(db)
	a.humanResourceService = services.NewHumanResourceService(hrRepo)
}

// Client service wrapper methods for Wails bindings

// GetClients retrieves multiple clients with optional query parameters
func (a *App) GetClients(params *entities.ClientQueryParams) (*entities.ClientListResponse, error) {
	if a.clientService == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return a.clientService.GetClients(a.ctx, params)
}

// GetClient retrieves a single client by ID
func (a *App) GetClient(id uint) (*entities.Client, error) {
	if a.clientService == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return a.clientService.GetClient(a.ctx, id)
}

// CreateClient creates a new client
func (a *App) CreateClient(client *entities.Client) (*entities.Client, error) {
	if a.clientService == nil {
		return nil, fmt.Errorf("client service not initialized")
	}
	return a.clientService.CreateClient(a.ctx, client)
}

// UpdateClient updates an existing client
func (a *App) UpdateClient(client *entities.Client) (int64, error) {
	if a.clientService == nil {
		return 0, fmt.Errorf("client service not initialized")
	}
	return a.clientService.UpdateClient(a.ctx, client)
}

// DeleteClient deletes a client by ID
func (a *App) DeleteClient(id uint) error {
	if a.clientService == nil {
		return fmt.Errorf("client service not initialized")
	}
	return a.clientService.DeleteClient(a.ctx, id)
}

// Human Resource service wrapper methods for Wails bindings

// GetHumanResources retrieves multiple human resources with optional query parameters
func (a *App) GetHumanResources(params *entities.HumanResourceQueryParams) (*entities.HumanResourceListResponse, error) {
	if a.humanResourceService == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return a.humanResourceService.GetHumanResources(a.ctx, params)
}

// GetHumanResource retrieves a single human resource by ID
func (a *App) GetHumanResource(id uint) (*entities.HumanResource, error) {
	if a.humanResourceService == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return a.humanResourceService.GetHumanResource(a.ctx, id)
}

// CreateHumanResource creates a new human resource
func (a *App) CreateHumanResource(humanResource *entities.HumanResource) (*entities.HumanResource, error) {
	if a.humanResourceService == nil {
		return nil, fmt.Errorf("human resource service not initialized")
	}
	return a.humanResourceService.CreateHumanResource(a.ctx, humanResource)
}

// UpdateHumanResource updates an existing human resource
func (a *App) UpdateHumanResource(humanResource *entities.HumanResource) (int64, error) {
	if a.humanResourceService == nil {
		return 0, fmt.Errorf("human resource service not initialized")
	}
	return a.humanResourceService.UpdateHumanResource(a.ctx, humanResource)
}

// DeleteHumanResource deletes a human resource by ID
func (a *App) DeleteHumanResource(id uint) error {
	if a.humanResourceService == nil {
		return fmt.Errorf("human resource service not initialized")
	}
	return a.humanResourceService.DeleteHumanResource(a.ctx, id)
}
