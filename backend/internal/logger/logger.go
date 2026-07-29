package logger

import (
	"context"
	"io"
	"log/slog"
)

// Log is the global logger instance
var Log *slog.Logger

// Init initializes the structured logger.
// For development, it uses TextHandler; for production/other environments, it uses JSONHandler.
func Init(env string, w io.Writer) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "development" || env == "local" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// Info logs messages at Info level
func Info(msg string, args ...any) {
	if Log != nil {
		Log.Info(msg, args...)
	}
}

// Error logs messages at Error level with an optional error object
func Error(msg string, err error, args ...any) {
	if Log != nil {
		if err != nil {
			args = append(args, slog.Any("error", err))
		}
		Log.Error(msg, args...)
	}
}

// Debug logs messages at Debug level
func Debug(msg string, args ...any) {
	if Log != nil {
		Log.Debug(msg, args...)
	}
}

// Warn logs messages at Warning level
func Warn(msg string, args ...any) {
	if Log != nil {
		Log.Warn(msg, args...)
	}
}

// WithContext returns a logger contextualized (placeholder for future trace context propagation)
func WithContext(ctx context.Context) *slog.Logger {
	if Log == nil {
		return slog.Default()
	}
	return Log
}
