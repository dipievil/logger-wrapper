package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dipievil/logger-wrapper/logging/interfaces"
)

// Logger is a wrapper around slog.Logger that provides additional functionality and context management.
type Logger struct {
	base     *slog.Logger
	notifier interfaces.Notifier
	err      error
}

// LoggerConfig holds the configuration for the Logger, including log level, build version, and environment.
type LoggerConfig struct {
	LogLevel     string
	BuildVersion string
	Environment  string
}

// NewLoggerConfig returns a default LoggerConfig with preset values for log level, build version, and environment.
func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{
		LogLevel:     "debug",
		BuildVersion: "dev",
		Environment:  "local",
	}
}

// LoggerOption defines a function type for configuring the Logger with optional parameters.
type LoggerOption func(*Logger)

// WithNotifier is a LoggerOption that sets a Notifier for the Logger, allowing it to send notifications on certain log events.
func WithNotifier(n interfaces.Notifier) LoggerOption {
	return func(l *Logger) {
		if g, ok := n.(*GotifyService); ok {
			isValid, err := g.Validate()
			if err != nil {
				l.err = err
				return
			}

			if !isValid {
				n = nil
			}
		}
		l.notifier = n
	}
}

// NewLogger creates a new Logger instance based on the provided configuration.
func NewLoggerWrapper(config LoggerConfig, opts ...LoggerOption) *Logger {

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

	l := &Logger{
		base: slog.New(handler).With(
			"version", config.BuildVersion,
			"environment", config.Environment,
		),
	}

	for _, opt := range opts {
		opt(l)
	}

	if l.err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", l.err))
	}

	return l
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

// getLevelInfoByString converts a string representation of a log level to the corresponding slog.Level.
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

// Info logs an informational message and notifies Gotify if a notifier is configured.
func (l *Logger) Info(msg string, args ...any) {
	l.base.Info(msg, args...)
}

// InfoNotify logs an informational message and sends a notification without additional context.
func (l *Logger) InfoNotify(msg string) {
	l.Info(msg)
	l.sendNotification(msg)
}

// Error logs an error message with optional arguments.
func (l *Logger) Error(msg string, args ...any) {
	l.base.Error(msg, args...)
}

// Warn logs a warning message with optional arguments.
func (l *Logger) Warn(msg string, args ...any) {
	l.base.Warn(msg, args...)
}

// InfoNotifyf logs a formatted informational message and sends a notification without additional context.
func (l *Logger) InfoNotifyf(msg string, args ...any) {
	l.Infof(msg, args...)
	l.sendNotification(fmt.Sprintf(msg, args...))
}

// Infof logs a formatted informational message and notifies Gotify
func (l *Logger) Infof(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.Info(message)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.Error(message)
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.Debug(message)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	l.Warn(message)
}

// sendNotification sends a notification using the configured Notifier, if available.
func (l *Logger) sendNotification(msg string) {
	if l.notifier == nil {
		return
	}

	go func(m string) {
		if err := l.notifier.Notify(m); err != nil {
			message := fmt.Sprintf("Failed to send notification %s", l.notifier.GetHost())
			l.base.Error(message, "error", err)
		}
	}(msg)
}
