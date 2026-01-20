package config

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

// MemoryDSN is the special DSN for in-memory SQLite database (draft mode)
// Using file::memory:?cache=shared ensures all connections share the same in-memory database
const MemoryDSN = "file::memory:?cache=shared"

// IsMemoryDSN checks if the given DSN is the memory database DSN
func IsMemoryDSN(dsn string) bool {
	return dsn == MemoryDSN
}

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	Cfg    = Load()
)

// getDefaultLogPath returns the default log path in user's home directory
func getDefaultLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir is not available
		return "logs/plancraft.log"
	}
	return filepath.Join(homeDir, ".plan-craft", "logs", "plancraft.log")
}

// getDefaultDBPath returns the default database path in user's home directory
func getDefaultDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir is not available
		return "data/plancraft.db"
	}
	return filepath.Join(homeDir, ".plan-craft", "data", "plancraft.db")
}

type Config struct {
	AppName    string   `env:"APP_NAME, default=plan-craft"`
	Environemt string   `env:"ENV, default=local"`
	DB         DBConfig `env:", prefix=DB_"`
	LogPath    string   `env:"LOG_PATH"`
	LogLevel   string   `env:"LOG_LEVEL, default=WARN"`
}

type DBConfig struct {
	DSN         string `env:"DSN"`
	JournalMode string `env:"JOURNAL_MODE, default=WAL"`
	Synchronous string `env:"SYNCHRONOUS, default=NORMAL"`
	ForeignKeys string `env:"FOREIGN_KEYS, default=ON"`
	BusyTimeout string `env:"BUSY_TIMEOUT, default=5000"`
	CacheSize   string `env:"CACHE_SIZE, default=-64000"`
	TempStore   string `env:"TEMP_STORE, default=MEMORY"`
	AutoVacuum  string `env:"AUTO_VACUUM, default=INCREMENTAL"`
}

func Load() Config {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		logger.ErrorContext(ctx, "Failed to load config", slog.Any("error", err))
	}

	// Set defaults for paths that need home directory expansion
	if c.LogPath == "" {
		c.LogPath = getDefaultLogPath()
	}
	if c.DB.DSN == "" {
		// Check for persisted database path from previous session
		if lastDBPath := loadLastDatabasePath(); lastDBPath != "" {
			c.DB.DSN = lastDBPath
		} else {
			// Start with memory database (draft mode) when no persisted path exists
			c.DB.DSN = MemoryDSN
		}
	}

	return c
}

// loadLastDatabasePath reads the persisted database path from the settings file.
// Returns empty string if no settings file exists or on error.
func loadLastDatabasePath() string {
	settingsPath := getSettingsFilePath()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}
	path := strings.TrimSpace(string(data))
	// Validate the file still exists before using it
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

// getSettingsFilePath returns the path to the settings file
func getSettingsFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".plan-craft-settings"
	}
	return filepath.Join(homeDir, ".plan-craft", "settings")
}

// EnsureLogDirectory creates the log directory if it doesn't exist
func EnsureLogDirectory(logPath string) error {
	logDir := filepath.Dir(logPath)
	if logDir != "" && logDir != "." {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}
	return nil
}
