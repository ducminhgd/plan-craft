package internal

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ducminhgd/plan-craft/config"
	gormslog "github.com/onrik/gorm-slog"
	wlr "github.com/wailsapp/wails/v2/pkg/logger"
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

func ConvertWailsLogLevel(l string) wlr.LogLevel {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return wlr.DEBUG
	case "INFO":
		return wlr.INFO
	case "WARN":
		return wlr.WARNING
	case "ERROR":
		return wlr.ERROR
	default:
		return wlr.DEBUG
	}
}
