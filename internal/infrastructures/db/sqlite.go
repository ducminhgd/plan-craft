package db

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/ducminhgd/plan-craft/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB is the global GORM database instance
	DB *gorm.DB
)

// InitializeDatabase initializes the database connection
func InitializeDatabase() error {
	cfg := config.Cfg

	// Ensure database directory exists
	dbPath := cfg.DB.DSN
	if dbPath == "" {
		dbPath = "data/plancraft.db"
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Build SQLite DSN with optimizations
	dsn := buildSQLiteDSN(cfg.DB)

	// Configure GORM logger
	gormLogger := configureLogger(cfg.LogLevel)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true, // Improve performance
		PrepareStmt:            true, // Cache prepared statements
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	// Assign to global variable
	DB = db

	slog.Info("Database initialized successfully",
		slog.String("path", dbPath),
		slog.String("driver", "sqlite"),
	)

	return nil
}

// buildSQLiteDSN builds an optimized SQLite DSN from configuration
func buildSQLiteDSN(dbConfig config.DBConfig) string {
	// SQLite connection string with optimizations
	// Reference: https://www.sqlite.org/pragma.html

	dbPath := dbConfig.DSN
	if dbPath == "" {
		dbPath = "data/plancraft.db"
	}

	// Build parameters only if they are set
	var params []string

	if dbConfig.JournalMode != "" {
		params = append(params, fmt.Sprintf("_journal_mode=%s", dbConfig.JournalMode))
	}
	if dbConfig.Synchronous != "" {
		params = append(params, fmt.Sprintf("_synchronous=%s", dbConfig.Synchronous))
	}
	if dbConfig.ForeignKeys != "" {
		params = append(params, fmt.Sprintf("_foreign_keys=%s", dbConfig.ForeignKeys))
	}
	if dbConfig.BusyTimeout != "" {
		params = append(params, fmt.Sprintf("_busy_timeout=%s", dbConfig.BusyTimeout))
	}
	if dbConfig.CacheSize != "" {
		params = append(params, fmt.Sprintf("_cache_size=%s", dbConfig.CacheSize))
	}
	if dbConfig.TempStore != "" {
		params = append(params, fmt.Sprintf("_temp_store=%s", dbConfig.TempStore))
	}
	if dbConfig.AutoVacuum != "" {
		params = append(params, fmt.Sprintf("_auto_vacuum=%s", dbConfig.AutoVacuum))
	}

	// Build DSN with parameters
	dsn := dbPath
	for i, param := range params {
		if i == 0 {
			dsn += "?" + param
		} else {
			dsn += "&" + param
		}
	}

	return dsn
}

// configureLogger configures GORM logger based on log level
func configureLogger(logLevel string) logger.Interface {
	var level logger.LogLevel

	switch logLevel {
	case "ERROR":
		level = logger.Error
	case "WARN":
		level = logger.Warn
	case "INFO":
		level = logger.Info
	case "DEBUG":
		level = logger.Info
	default:
		level = logger.Warn
	}

	return logger.New(
		slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
			ParameterizedQueries:      false, // Log SQL with parameters for debugging
		},
	)
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	slog.Info("Database connection closed")
	return nil
}

// HealthCheck checks the database connection health
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
