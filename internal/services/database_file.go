package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseFileService handles database file operations (open, save as)
type DatabaseFileService struct {
	ctx             context.Context
	currentDBPath   string
	db              *gorm.DB
	mu              sync.RWMutex
	onDBChanged     func(*gorm.DB) // callback when DB changes
}

// NewDatabaseFileService creates a new DatabaseFileService
func NewDatabaseFileService() *DatabaseFileService {
	return &DatabaseFileService{}
}

// NewDatabaseFileServiceWithDB creates a new DatabaseFileService with an initial database
func NewDatabaseFileServiceWithDB(ctx context.Context, db *gorm.DB, dbPath string) *DatabaseFileService {
	return &DatabaseFileService{
		ctx:           ctx,
		db:            db,
		currentDBPath: dbPath,
	}
}

// SetupDatabaseFileService initializes an existing service with context and database
func SetupDatabaseFileService(svc *DatabaseFileService, ctx context.Context, db *gorm.DB, dbPath string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.ctx = ctx
	svc.db = db
	svc.currentDBPath = dbPath
}

// GetCurrentDatabasePath returns the current database file path
func (s *DatabaseFileService) GetCurrentDatabasePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentDBPath
}

// OpenDatabase opens a SQLite database file using a file dialog
func (s *DatabaseFileService) OpenDatabase() (string, error) {
	if s.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	// Show open file dialog
	filePath, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Open Database",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "SQLite Database",
				Pattern:     "*.db;*.sqlite;*.sqlite3",
			},
			{
				DisplayName: "All Files",
				Pattern:     "*",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}

	// User cancelled the dialog
	if filePath == "" {
		return "", nil
	}

	// Open the new database
	newDB, err := s.openDatabaseFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}

	// Close the old database connection
	s.mu.Lock()
	if s.db != nil {
		if sqlDB, err := s.db.DB(); err == nil {
			sqlDB.Close()
		}
	}
	s.db = newDB
	s.currentDBPath = filePath
	s.mu.Unlock()

	// Notify callback that DB has changed
	if s.onDBChanged != nil {
		s.onDBChanged(newDB)
	}

	return filePath, nil
}

// SaveDatabaseAs saves the current database to a new file
func (s *DatabaseFileService) SaveDatabaseAs() (string, error) {
	if s.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	s.mu.RLock()
	currentPath := s.currentDBPath
	s.mu.RUnlock()

	if currentPath == "" {
		return "", fmt.Errorf("no database is currently open")
	}

	// Get default filename for save dialog
	defaultFilename := filepath.Base(currentPath)

	// Show save file dialog
	filePath, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		Title:           "Save Database As",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "SQLite Database",
				Pattern:     "*.db;*.sqlite;*.sqlite3",
			},
			{
				DisplayName: "All Files",
				Pattern:     "*",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to open save dialog: %w", err)
	}

	// User cancelled the dialog
	if filePath == "" {
		return "", nil
	}

	// Ensure the file has .db extension if no extension provided
	if filepath.Ext(filePath) == "" {
		filePath += ".db"
	}

	// Copy the current database file to the new location
	if err := s.copyDatabaseFile(currentPath, filePath); err != nil {
		return "", fmt.Errorf("failed to copy database: %w", err)
	}

	return filePath, nil
}

// openDatabaseFile opens a SQLite database file with the configured pragmas
func (s *DatabaseFileService) openDatabaseFile(path string) (*gorm.DB, error) {
	cfg := config.Cfg

	// Build DSN with pragma settings
	dsn := fmt.Sprintf("%s?_journal_mode=%s&_synchronous=%s&_foreign_keys=%s&_busy_timeout=%s&cache=shared&_cache_size=%s&_temp_store=%s&_auto_vacuum=%s",
		path,
		cfg.DB.JournalMode,
		cfg.DB.Synchronous,
		cfg.DB.ForeignKeys,
		cfg.DB.BusyTimeout,
		cfg.DB.CacheSize,
		cfg.DB.TempStore,
		cfg.DB.AutoVacuum,
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate entities to ensure schema is up to date
	err = db.AutoMigrate(
		&entities.Client{},
		&entities.HumanResource{},
		&entities.Project{},
		&entities.ProjectResource{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}

// copyDatabaseFile copies a database file to a new location
func (s *DatabaseFileService) copyDatabaseFile(src, dst string) error {
	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Ensure the file is synced to disk
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// OpenGuides opens the guides page in the default browser
func (s *DatabaseFileService) OpenGuides() error {
	if s.ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	// TODO: Replace with actual guides URL when available
	guidesURL := "https://github.com/ducminhgd/plan-craft/wiki"
	runtime.BrowserOpenURL(s.ctx, guidesURL)
	return nil
}
