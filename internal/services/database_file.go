package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal/entities"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OnDBChangedFunc is a callback type for when the database changes
type OnDBChangedFunc func(db *gorm.DB)

// DatabaseFileService handles database file operations (open, save as)
type DatabaseFileService struct {
	ctx           context.Context
	currentDBPath string
	isMemoryDB    bool // Tracks if current DB is memory/draft mode
	db            *gorm.DB
	mu            sync.RWMutex
	onDBChanged   OnDBChangedFunc
}

// NewDatabaseFileService creates a new DatabaseFileService
func NewDatabaseFileService() *DatabaseFileService {
	return &DatabaseFileService{}
}

// SetupDatabaseFileService initializes an existing service with context and database
func SetupDatabaseFileService(svc *DatabaseFileService, ctx context.Context, db *gorm.DB, dbPath string, isMemory bool) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.ctx = ctx
	svc.db = db
	svc.currentDBPath = dbPath
	svc.isMemoryDB = isMemory
}

// SetOnDBChanged sets a callback that will be invoked when the database is switched
func SetOnDBChanged(svc *DatabaseFileService, callback OnDBChangedFunc) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.onDBChanged = callback
}

// GetCurrentDatabasePath returns the current database file path
func (s *DatabaseFileService) GetCurrentDatabasePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentDBPath
}

// IsMemoryDatabase returns true if the current database is in-memory (draft mode)
func (s *DatabaseFileService) IsMemoryDatabase() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isMemoryDB
}

// HasUnsavedChanges returns true if in draft mode and there are records in the database.
// This is used to determine if the user should be prompted before exiting.
func (s *DatabaseFileService) HasUnsavedChanges() bool {
	s.mu.RLock()
	isMemory := s.isMemoryDB
	db := s.db
	s.mu.RUnlock()

	// If not in memory/draft mode, no "unsaved" changes to worry about
	if !isMemory {
		return false
	}

	if db == nil {
		return false
	}

	// Check if any of the tracked entities have records
	var count int64

	// Check clients
	if err := db.Model(&entities.Client{}).Count(&count).Error; err == nil && count > 0 {
		return true
	}

	// Check human resources
	if err := db.Model(&entities.HumanResource{}).Count(&count).Error; err == nil && count > 0 {
		return true
	}

	// Check projects
	if err := db.Model(&entities.Project{}).Count(&count).Error; err == nil && count > 0 {
		return true
	}

	// Check project resources
	if err := db.Model(&entities.ProjectResource{}).Count(&count).Error; err == nil && count > 0 {
		return true
	}

	return false
}

// OpenDatabase opens a SQLite database file using a file dialog.
// It switches to the new database and re-wires all dependencies via the onDBChanged callback.
// Returns the selected file path.
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

	// Validate that the file exists
	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file")
	}

	// Open the new database
	newDB, err := s.openDatabaseFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}

	// Close the old database connection
	s.mu.Lock()
	oldDB := s.db
	s.db = newDB
	s.currentDBPath = filePath
	s.isMemoryDB = false // Opening a file means we're no longer in draft mode
	callback := s.onDBChanged
	s.mu.Unlock()

	if oldDB != nil {
		if sqlDB, err := oldDB.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// Persist the selected database path for next app launch
	if err := s.saveLastDatabasePath(filePath); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("Warning: failed to save database path: %v\n", err)
	}

	// Notify the app to re-wire dependencies with the new database
	if callback != nil {
		callback(newDB)
	}

	return filePath, nil
}

// saveLastDatabasePath persists the database path for the next app launch
func (s *DatabaseFileService) saveLastDatabasePath(dbPath string) error {
	settings := config.Settings{
		LastDatabasePath: dbPath,
	}
	if err := config.SaveSettings(settings); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}
	return nil
}

// SaveDatabaseAs saves the current database to a new file and switches to it
func (s *DatabaseFileService) SaveDatabaseAs() (string, error) {
	if s.ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	s.mu.RLock()
	isMemory := s.isMemoryDB
	db := s.db
	s.mu.RUnlock()

	if db == nil {
		return "", fmt.Errorf("database connection not available")
	}

	// Get default filename for save dialog
	defaultFilename := "plancraft.db"

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

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(filePath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	if isMemory {
		// For memory database, we need to create the new file and copy data
		// Open the new file database first
		newDB, err := s.openDatabaseFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to create database file: %w", err)
		}

		// Copy all data from memory to new file using SQLite backup
		if err := s.copyDatabaseContent(db, newDB); err != nil {
			// Close the new DB on error
			if sqlDB, closeErr := newDB.DB(); closeErr == nil {
				sqlDB.Close()
			}
			return "", fmt.Errorf("failed to copy data: %w", err)
		}

		// Switch to new database
		s.mu.Lock()
		oldDB := s.db
		s.db = newDB
		s.currentDBPath = filePath
		s.isMemoryDB = false
		callback := s.onDBChanged
		s.mu.Unlock()

		// Close old memory database
		if oldDB != nil {
			if sqlDB, err := oldDB.DB(); err == nil {
				sqlDB.Close()
			}
		}

		// Persist the database path for next app launch
		if err := s.saveLastDatabasePath(filePath); err != nil {
			fmt.Printf("Warning: failed to save database path: %v\n", err)
		}

		// Notify the app to re-wire dependencies with the new database
		if callback != nil {
			callback(newDB)
		}
	} else {
		// For file-based database, use VACUUM INTO to create a consistent copy
		// VACUUM INTO creates a new database file with all data from the current database.
		// It's atomic and handles WAL mode correctly (available since SQLite 3.27.0).
		if err := db.Exec("VACUUM INTO ?", filePath).Error; err != nil {
			return "", fmt.Errorf("failed to save database: %w", err)
		}

		// Open and switch to the new file
		newDB, err := s.openDatabaseFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to open saved database: %w", err)
		}

		// Switch to new database
		s.mu.Lock()
		oldDB := s.db
		s.db = newDB
		s.currentDBPath = filePath
		s.isMemoryDB = false
		callback := s.onDBChanged
		s.mu.Unlock()

		// Close old database
		if oldDB != nil {
			if sqlDB, err := oldDB.DB(); err == nil {
				sqlDB.Close()
			}
		}

		// Persist the database path for next app launch
		if err := s.saveLastDatabasePath(filePath); err != nil {
			fmt.Printf("Warning: failed to save database path: %v\n", err)
		}

		// Notify the app to re-wire dependencies with the new database
		if callback != nil {
			callback(newDB)
		}
	}

	return filePath, nil
}

// copyDatabaseContent copies all table data from source to destination database
func (s *DatabaseFileService) copyDatabaseContent(srcDB, dstDB *gorm.DB) error {
	// Get the underlying sql.DB connections
	srcSqlDB, err := srcDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get source sql.DB: %w", err)
	}
	dstSqlDB, err := dstDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get destination sql.DB: %w", err)
	}

	// Get the raw connections for backup
	srcConn, err := srcSqlDB.Conn(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to get source connection: %w", err)
	}
	defer srcConn.Close()

	dstConn, err := dstSqlDB.Conn(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to get destination connection: %w", err)
	}
	defer dstConn.Close()

	// Copy data table by table using INSERT INTO ... SELECT
	// Get list of tables from source
	var tables []string
	rows, err := srcSqlDB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	// Copy each table's data
	for _, table := range tables {
		// Get all rows from source table
		sourceRows, err := srcSqlDB.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			continue // Table might not have data or might be auto-generated
		}

		cols, err := sourceRows.Columns()
		if err != nil {
			sourceRows.Close()
			continue
		}

		// Prepare values slice
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Build insert statement
		placeholders := make([]string, len(cols))
		for i := range placeholders {
			placeholders[i] = "?"
		}

		insertSQL := fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, joinStrings(placeholders, ", "))

		// Copy rows
		for sourceRows.Next() {
			if err := sourceRows.Scan(valuePtrs...); err != nil {
				continue
			}
			if _, err := dstSqlDB.Exec(insertSQL, values...); err != nil {
				// Ignore duplicate key errors, continue with other rows
				continue
			}
		}
		sourceRows.Close()
	}

	return nil
}

// joinStrings joins strings with a separator (helper function)
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
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
