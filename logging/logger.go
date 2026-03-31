package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	base *slog.Logger
}

type LoggerConfig struct {
	LogLevel     string
	BuildVersion string
	Environment  string
}

func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{
		LogLevel:     "debug",
		BuildVersion: "dev",
		Environment:  "local",
	}
}

// NewLogger creates a new Logger instance based on the provided configuration.
func NewLoggerWrapper(config LoggerConfig) *Logger {

	if config.LogLevel == "" {
		config.LogLevel = "debug"
	}

	if config.BuildVersion == "" {
		config.BuildVersion = "dev"
	}

	if config.Environment == "" {
		config.Environment = "local"
	}

	level := getLevelInfoByString(config.LogLevel)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	logger = logger.With("version", config.BuildVersion)
	logger = logger.With("environment", config.Environment)

	return &Logger{
		base: logger,
	}
}

// Base returns the underlying slog.Logger instance for advanced usage.
func (l *Logger) Base() *slog.Logger {
	return l.base
}

// With creates a new logger with additional context fields.
func (l *Logger) With(args ...any) *slog.Logger {
	return l.base.With(args...)
}

// Audit logs an audit event with the given action and arguments.
func (l *Logger) Audit(ctx context.Context, action string, args ...any) {

	// If args has level, use it, otherwise default to info
	level := slog.LevelInfo
	for i := 0; i < len(args)-1; i += 2 {
		if args[i] == "level" {
			if lvl, ok := args[i+1].(string); ok {
				level = getLevelInfoByString(lvl)
			}
			break
		}
	}

	l.base.Log(ctx, level, "audit", append([]any{"action", action}, args...)...)
}

func getLevelInfoByString(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// Debug logs a debug message with optional arguments.
func (l *Logger) Debug(msg string, args ...any) {
	l.base.Debug(msg, args...)
}

// Info logs an informational message with optional arguments.
func (l *Logger) Info(msg string, args ...any) {
	l.base.Info(msg, args...)
}

// Error logs an error message with optional arguments.
func (l *Logger) Error(msg string, args ...any) {
	l.base.Error(msg, args...)
}

// Warn logs a warning message with optional arguments.
func (l *Logger) Warn(msg string, args ...any) {
	l.base.Warn(msg, args...)
}

// Infof logs a formatted informational message.
func (l *Logger) Infof(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.base.Info(message)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.base.Error(message)
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.base.Debug(message)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.base.Warn(message)
}
