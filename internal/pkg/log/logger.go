package log

import (
	"log/slog"
	"os"
)

// Global logger instance
var Log *slog.Logger

// initailize sets up the global logger
func Init() {
	Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(Log)
}

// Optional: convenience function to get the logger
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
