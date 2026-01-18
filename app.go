package main

import (
	"context"
	"log"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal/handlers"
	"github.com/ducminhgd/plan-craft/internal/infrastructures"
	"github.com/ducminhgd/plan-craft/internal/repositories"
	"github.com/ducminhgd/plan-craft/internal/services"
	"gorm.io/gorm"
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
func (a *App) startup(ctx context.Context, dbFileService *services.DatabaseFileService) {
	a.ctx = ctx

	// Initialize database
	db, err := infrastructures.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Determine if this is a memory database (draft mode)
	isMemory := config.IsMemoryDSN(config.Cfg.DB.DSN)

	// Setup the database file service with context and database
	services.SetupDatabaseFileService(dbFileService, ctx, db, config.Cfg.DB.DSN, isMemory)

	// Wire dependencies with the initial database
	a.wireHandlers(ctx, db)

	// Set up callback to re-wire handlers when database changes
	services.SetOnDBChanged(dbFileService, func(newDB *gorm.DB) {
		a.wireHandlers(ctx, newDB)
	})
}

// wireHandlers creates all repositories, services, and handlers for the given database
func (a *App) wireHandlers(ctx context.Context, db *gorm.DB) {
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

	// Update handlers container with new handlers
	a.Handlers = handlers.NewHandlers(clientHandler, hrHandler, projectHandler, projectResourceHandler)
}
