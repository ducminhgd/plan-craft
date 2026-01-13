package internal

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ducminhgd/plan-craft/config"
	gormslog "github.com/onrik/gorm-slog"
)

var (
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
)

func ConvertSlogLevel(l string) slog.Level {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func NewGORMLogger(cfg *config.Config) *gormslog.Logger {
	return gormslog.New(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     ConvertSlogLevel(cfg.LogLevel),
	})))
}
