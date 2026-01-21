package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

type Config struct {
	AppName     string   `yaml:"app_name"`
	Environment string   `yaml:"environment"`
	DB          DBConfig `yaml:"database"`
	LogPath     string   `yaml:"log_path"`
	LogLevel    string   `yaml:"log_level"`
}

type DBConfig struct {
	DSN         string `yaml:"dsn"`
	JournalMode string `yaml:"journal_mode"`
	Synchronous string `yaml:"synchronous"`
	ForeignKeys string `yaml:"foreign_keys"`
	BusyTimeout string `yaml:"busy_timeout"`
	CacheSize   string `yaml:"cache_size"`
	TempStore   string `yaml:"temp_store"`
	AutoVacuum  string `yaml:"auto_vacuum"`
}

// getConfigFilePaths returns the list of paths to search for the config file
func getConfigFilePaths() []string {
	paths := []string{
		"config.yaml",
		"config.yml",
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths,
			filepath.Join(homeDir, ".plan-craft", "config.yaml"),
			filepath.Join(homeDir, ".plan-craft", "config.yml"),
		)
	}

	return paths
}

// Load loads the configuration from YAML file
func Load() Config {
	c := Config{
		AppName:     "plan-craft",
		Environment: "local",
		LogLevel:    "WARN",
		DB: DBConfig{
			JournalMode: "WAL",
			Synchronous: "NORMAL",
			ForeignKeys: "ON",
			BusyTimeout: "5000",
			CacheSize:   "-64000",
			TempStore:   "MEMORY",
			AutoVacuum:  "INCREMENTAL",
		},
	}

	// Try to load config from file
	configPaths := getConfigFilePaths()
	for _, configPath := range configPaths {
		data, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		if err := yaml.Unmarshal(data, &c); err != nil {
			logger.Error("Failed to parse config file", slog.String("path", configPath), slog.Any("error", err))
			continue
		}

		// Successfully loaded config
		break
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

	// Ensure DB config defaults are set if not specified in YAML
	if c.DB.JournalMode == "" {
		c.DB.JournalMode = "WAL"
	}
	if c.DB.Synchronous == "" {
		c.DB.Synchronous = "NORMAL"
	}
	if c.DB.ForeignKeys == "" {
		c.DB.ForeignKeys = "ON"
	}
	if c.DB.BusyTimeout == "" {
		c.DB.BusyTimeout = "5000"
	}
	if c.DB.CacheSize == "" {
		c.DB.CacheSize = "-64000"
	}
	if c.DB.TempStore == "" {
		c.DB.TempStore = "MEMORY"
	}
	if c.DB.AutoVacuum == "" {
		c.DB.AutoVacuum = "INCREMENTAL"
	}

	return c
}

// Settings represents the application settings stored in YAML format
type Settings struct {
	LastDatabasePath string `yaml:"last_database_path"`
}

// loadLastDatabasePath reads the persisted database path from the settings file.
// Returns empty string if no settings file exists or on error.
func loadLastDatabasePath() string {
	settingsPath := GetSettingsFilePath()
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return ""
	}

	path := settings.LastDatabasePath
	// Validate the file still exists before using it
	if path != "" {
		if _, err := os.Stat(path); err != nil {
			return ""
		}
	}
	return path
}

// GetSettingsFilePath returns the path to the settings file
func GetSettingsFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".plan-craft-settings.yaml"
	}
	return filepath.Join(homeDir, ".plan-craft", "settings.yaml")
}

// SaveSettings saves the application settings to the YAML settings file
func SaveSettings(settings Settings) error {
	settingsPath := GetSettingsFilePath()

	// Create settings directory if it doesn't exist
	settingsDir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(&settings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
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
