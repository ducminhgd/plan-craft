package services

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MenuService struct{}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func (m *MenuService) BuildApplicationMenu(ctx context.Context) *menu.Menu {
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Open", nil, func(_ *menu.CallbackData) {
		// No action for now
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Exit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(ctx)
	})

	// Help menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("Guides", nil, func(_ *menu.CallbackData) {
		// No action for now
	})
	helpMenu.AddSeparator()
	helpMenu.AddText("About", nil, func(_ *menu.CallbackData) {
		m.ShowAboutDialog(ctx)
	})

	return appMenu
}

func (m *MenuService) ShowAboutDialog(ctx context.Context) {
	runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "About Plan Craft",
		Message: "Plan Craft v1.0.0\n\nA desktop project management and estimation tool.",
	})
}
