package services

import (
	"context"
	"path/filepath"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MenuService struct {
	ctx           context.Context
	dbFileService *DatabaseFileService
}

func NewMenuService() *MenuService {
	return &MenuService{}
}

// SetContext sets the runtime context for menu operations
func (m *MenuService) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// SetDatabaseFileService sets the database file service for menu actions
func (m *MenuService) SetDatabaseFileService(dbFileService *DatabaseFileService) {
	m.dbFileService = dbFileService
}

func (m *MenuService) BuildApplicationMenu(ctx context.Context) *menu.Menu {
	// Store context if provided
	if ctx != nil {
		m.ctx = ctx
	}

	appMenu := menu.NewMenu()
	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Open file", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			_, err := m.dbFileService.OpenDatabase()
			if err != nil {
				internal.Logger.Error("failed to open database", "error", err)
				runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "Open file failed",
					Message: err.Error(),
				})
			}
		}
	})
	fileMenu.AddText("Save as", keys.Combo("s", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			_, err := m.dbFileService.SaveDatabaseAs()
			if err != nil {
				internal.Logger.Error("failed to save database", "error", err)
				runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "Save as failed",
					Message: err.Error(),
				})
			}
		}
	})
	fileMenu.AddText("Close", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {
		if m.dbFileService != nil {
			err := m.dbFileService.CloseDatabase()
			if err != nil {
				internal.Logger.Error("failed to close database", "error", err)
				runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "Close database failed",
					Message: err.Error(),
				})
			}
		}
	})
	fileMenu.AddSeparator()

	// Recent Files submenu
	recentMenu := fileMenu.AddSubmenu("Recent Files")
	m.populateRecentFilesMenu(recentMenu)

	fileMenu.AddSeparator()
	fileMenu.AddText("Exit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		m.handleExit()
	})

	// Help menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("Guides", nil, func(_ *menu.CallbackData) {
		m.OpenGuides()
	})
	helpMenu.AddSeparator()
	helpMenu.AddText("About", nil, func(_ *menu.CallbackData) {
		if m.ctx != nil {
			m.ShowAboutDialog(m.ctx)
		}
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

// OpenGuides opens the guides page in the default browser
func (m *MenuService) OpenGuides() error {
	if m.ctx == nil {
		return nil
	}
	// TODO: Replace with actual guides URL when available
	guidesURL := "https://github.com/ducminhgd/plan-craft/wiki"
	runtime.BrowserOpenURL(m.ctx, guidesURL)
	return nil
}

// handleExit handles the exit action with unsaved changes confirmation
func (m *MenuService) handleExit() {
	if m.ctx == nil {
		return
	}

	// Check if in draft mode with unsaved changes
	if m.dbFileService != nil && m.dbFileService.IsMemoryDatabase() {
		hasChanges := m.dbFileService.HasUnsavedChanges()
		if hasChanges {
			// First ask if user wants to save
			result, err := runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
				Type:          runtime.QuestionDialog,
				Title:         "Unsaved Changes",
				Message:       "You have unsaved changes in draft mode. Would you like to save before exiting?",
				DefaultButton: "Yes",
				CancelButton:  "No",
			})
			if err != nil {
				internal.Logger.Error("failed to show exit dialog", "error", err)
				return
			}

			// "Yes" means save first, "No" means exit without saving
			if result == "Yes" {
				// Save the database first
				_, err := m.dbFileService.SaveDatabaseAs()
				if err != nil {
					internal.Logger.Error("failed to save database", "error", err)
					runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
						Type:    runtime.ErrorDialog,
						Title:   "Save Failed",
						Message: err.Error(),
					})
					return
				}
			}
			runtime.Quit(m.ctx)
			return
		}
	}

	// No unsaved changes or not in draft mode - show simple exit confirmation
	result, err := runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Exit Application",
		Message:       "Are you sure you want to exit Plan Craft?",
		DefaultButton: "Yes",
		CancelButton:  "No",
	})
	if err != nil {
		internal.Logger.Error("failed to show exit dialog", "error", err)
		return
	}

	if result == "Yes" {
		runtime.Quit(m.ctx)
	}
}

// populateRecentFilesMenu populates the recent files submenu with up to 10 recent files
func (m *MenuService) populateRecentFilesMenu(recentMenu *menu.Menu) {
	recentFiles := config.GetRecentFiles()

	if len(recentFiles) == 0 {
		recentMenu.AddText("No recent files", nil, nil).Disabled = true
		return
	}

	for _, filePath := range recentFiles {
		// Capture filePath in closure
		path := filePath
		// Show just the filename in the menu, with full path as tooltip would be nice but not supported
		displayName := filepath.Base(path)
		recentMenu.AddText(displayName, nil, func(_ *menu.CallbackData) {
			if m.dbFileService != nil {
				m.dbFileService.OpenDatabasePath(path)
			}
		})
	}
}
