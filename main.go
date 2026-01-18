package main

import (
	"context"
	"embed"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal"
	"github.com/ducminhgd/plan-craft/internal/services"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Ensure log directory exists
	if err := config.EnsureLogDirectory(cfg.LogPath); err != nil {
		println("Warning: Failed to create log directory:", err.Error())
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create menu service
	menuService := services.NewMenuService()

	// Menu will be built in OnStartup after context is available
	var appMenu *menu.Menu

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
			app.startup(ctx)
			appMenu = menuService.BuildApplicationMenu(ctx)
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
