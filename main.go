package main

import (
	"context"
	"embed"

	"github.com/ducminhgd/plan-craft/internal/services"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create menu service
	menuService := services.NewMenuService()

	// Menu will be built in OnStartup after context is available
	var appMenu *menu.Menu

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Plan Craft",
		Width:  1024,
		Height: 768,
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
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
