package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(environment, level string) *slog.Logger {
	options := &slog.HandlerOptions{Level: parseLevel(level)}
	if environment == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, options)).With("service", "gopa")
	}
	return slog.New(slog.NewTextHandler(os.Stdout, options)).With("service", "gopa")
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
