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
	ctx           context.Context
	dbFileService *services.DatabaseFileService
	menuService   *services.MenuService
	*handlers.Handlers
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context, dbFileService *services.DatabaseFileService, menuService *services.MenuService) {
	a.ctx = ctx
	a.dbFileService = dbFileService
	a.menuService = menuService

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

// GetCurrentDatabasePath returns the current database file path
func (a *App) GetCurrentDatabasePath() string {
	return a.dbFileService.GetCurrentDatabasePath()
}

// IsMemoryDatabase returns true if the current database is in-memory (draft mode)
func (a *App) IsMemoryDatabase() bool {
	return a.dbFileService.IsMemoryDatabase()
}

// HasUnsavedChanges returns true if in draft mode and there are records in the database
func (a *App) HasUnsavedChanges() bool {
	return a.dbFileService.HasUnsavedChanges()
}

// OpenDatabase opens a SQLite database file using a file dialog
func (a *App) OpenDatabase() (string, error) {
	return a.dbFileService.OpenDatabase()
}

// SaveDatabaseAs saves the current database to a new file and switches to it
func (a *App) SaveDatabaseAs() (string, error) {
	return a.dbFileService.SaveDatabaseAs()
}

// OpenGuides opens the guides page in the default browser
func (a *App) OpenGuides() error {
	return a.menuService.OpenGuides()
}

// CloseDatabase closes the current database and switches to draft mode
func (a *App) CloseDatabase() error {
	return a.dbFileService.CloseDatabase()
}

// GetRecentFiles returns the list of recent database files
func (a *App) GetRecentFiles() []string {
	return config.GetRecentFiles()
}

// OpenRecentFile opens a specific database file from the recent files list
func (a *App) OpenRecentFile(filePath string) error {
	return a.dbFileService.OpenDatabasePath(filePath)
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

	projectRoleRepo := repositories.NewProjectRoleRepository(db)
	projectRoleService := services.NewProjectRoleService(projectRoleRepo)
	projectRoleHandler := handlers.NewProjectRoleHandler(ctx, projectRoleService)

	milestoneRepo := repositories.NewMilestoneRepository(db)
	milestoneService := services.NewMilestoneService(milestoneRepo)
	milestoneHandler := handlers.NewMilestoneHandler(ctx, milestoneService)

	taskRepo := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(ctx, taskService)

	// Update handlers container with new handlers
	a.Handlers = handlers.NewHandlers(clientHandler, hrHandler, projectHandler, projectResourceHandler, projectRoleHandler, milestoneHandler, taskHandler)
}
