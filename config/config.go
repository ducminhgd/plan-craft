package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/sethvargo/go-envconfig"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	Cfg    = Load()
)

type Config struct {
	AppName    string   `env:"APP_NAME, default=plan-craft"`
	Environemt string   `env:"ENV, default=local"`
	DB         DBConfig `env:", prefix=DB_"`
	LogPath    string   `env:"LOG_PATH, default=logs/plancraft.log"`
	LogLevel   string   `env:"LOG_LEVEL, default=WARN"`
}

type DBConfig struct {
	DSN         string `env:"DSN, default=data/plancraft.db"`
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
	return c
}
