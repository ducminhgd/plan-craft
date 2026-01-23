package main

import (
	"context"
	"embed"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal"
	"github.com/ducminhgd/plan-craft/internal/services"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Ensure log directory exists
	if err := config.EnsureLogDirectory(config.Cfg.LogPath); err != nil {
		println("Warning: Failed to create log directory:", err.Error())
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create database file service
	dbFileService := services.NewDatabaseFileService()

	// Create menu service and connect it with database file service
	menuService := services.NewMenuService()
	menuService.SetDatabaseFileService(dbFileService)

	// Build menu before wails.Run() - context will be set in callbacks
	appMenu := menuService.BuildApplicationMenu(nil)

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Plan Craft",
		MinWidth:  1024,
		MinHeight: 768,
		Frameless: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Menu:             appMenu,
		OnStartup: func(ctx context.Context) {
			// Update menu service with context for runtime operations
			menuService.SetContext(ctx)

			// Initialize app (including database) and database file service
			app.startup(ctx, dbFileService, menuService)
		},
		Bind: []interface{}{
			app,
		},
		Logger:             logger.NewFileLogger(config.Cfg.LogPath),
		LogLevel:           internal.ConvertWailsLogLevel(config.Cfg.LogLevel),
		LogLevelProduction: logger.ERROR,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
