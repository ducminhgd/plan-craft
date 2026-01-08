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
	ctx           context.Context
	clientService *services.ClientService
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
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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
