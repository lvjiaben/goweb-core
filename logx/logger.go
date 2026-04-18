package logx

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}

	if strings.EqualFold(cfg.Format, "json") {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
