package infrastructures

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ducminhgd/plan-craft/config"
	"github.com/ducminhgd/plan-craft/internal/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDatabase initializes the database connection and runs migrations
func InitializeDatabase() (*gorm.DB, error) {
	cfg := config.Cfg

	// Only create directory for file-based databases (not memory)
	if !config.IsMemoryDSN(cfg.DB.DSN) {
		dbDir := filepath.Dir(cfg.DB.DSN)
		if dbDir != "" && dbDir != "." {
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
		}
	}

	// Build DSN with pragma settings
	// Memory database uses different settings (no WAL, shared cache is already in MemoryDSN)
	var dsn string
	if config.IsMemoryDSN(cfg.DB.DSN) {
		// MemoryDSN already has ?cache=shared, so use & for additional params
		dsn = fmt.Sprintf("%s&_journal_mode=MEMORY&_synchronous=%s&_foreign_keys=%s&_busy_timeout=%s&_cache_size=%s&_temp_store=%s",
			cfg.DB.DSN,
			cfg.DB.Synchronous,
			cfg.DB.ForeignKeys,
			cfg.DB.BusyTimeout,
			cfg.DB.CacheSize,
			cfg.DB.TempStore,
		)
	} else {
		dsn = fmt.Sprintf("%s?_journal_mode=%s&_synchronous=%s&_foreign_keys=%s&_busy_timeout=%s&cache=shared&_cache_size=%s&_temp_store=%s&_auto_vacuum=%s",
			cfg.DB.DSN,
			cfg.DB.JournalMode,
			cfg.DB.Synchronous,
			cfg.DB.ForeignKeys,
			cfg.DB.BusyTimeout,
			cfg.DB.CacheSize,
			cfg.DB.TempStore,
			cfg.DB.AutoVacuum,
		)
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate entities
	err = db.AutoMigrate(
		&entities.Client{},
		&entities.HumanResource{},
		&entities.Project{},
		&entities.ProjectResource{},
		&entities.ProjectRole{},
		&entities.Milestone{},
		&entities.Task{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}
