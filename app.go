package main

import (
	"context"
	"log"

	"github.com/ducminhgd/plan-craft/internal/handlers"
	"github.com/ducminhgd/plan-craft/internal/infrastructures"
	"github.com/ducminhgd/plan-craft/internal/repositories"
	"github.com/ducminhgd/plan-craft/internal/services"
)

// App struct
type App struct {
	ctx context.Context
	*handlers.Handlers
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

	// Wire dependencies: repository → service → handler
	clientRepo := repositories.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(ctx, clientService)

	hrRepo := repositories.NewHRRepository(db)
	hrService := services.NewHumanResourceService(hrRepo)
	hrHandler := handlers.NewHumanResourceHandler(ctx, hrService)

	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(ctx, projectService)

	projectResourceRepo := repositories.NewProjectResourceRepository(db)
	projectResourceService := services.NewProjectResourceService(projectResourceRepo)
	projectResourceHandler := handlers.NewProjectResourceHandler(ctx, projectResourceService)

	// Initialize handlers container with all handlers
	a.Handlers = handlers.NewHandlers(clientHandler, hrHandler, projectHandler, projectResourceHandler)
}
