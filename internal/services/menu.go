package services

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MenuService struct {
	dbFileService *DatabaseFileService
}

func NewMenuService() *MenuService {
	return &MenuService{}
}

// SetDatabaseFileService sets the database file service for menu actions
func (m *MenuService) SetDatabaseFileService(dbFileService *DatabaseFileService) {
	m.dbFileService = dbFileService
}

func (m *MenuService) BuildApplicationMenu(ctx context.Context) *menu.Menu {
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Open file", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			m.dbFileService.OpenDatabase()
		}
	})
	fileMenu.AddText("Save as", keys.Combo("s", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			m.dbFileService.SaveDatabaseAs()
		}
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Exit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(ctx)
	})

	// Help menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("Guides", nil, func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			m.dbFileService.OpenGuides()
		}
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
