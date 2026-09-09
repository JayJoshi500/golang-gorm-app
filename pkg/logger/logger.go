package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New builds a JSON-structured slog.Logger. Using the standard library's
// slog keeps the app dependency-free for something as core as logging.
func New(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: lvl == slog.LevelDebug,
	})

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
