package interfaces

import "context"

type Logger interface {
	// Debug logs a message at the debug level with optional arguments.
	Debug(msg string, args ...any)

	// Debugf logs a formatted message at the debug level with optional arguments.
	Debugf(msg string, args ...any)

	// Info logs a message at the info level with optional arguments and sends a notification if a notifier is configured.
	Info(msg string, args ...any)

	// Infof logs a formatted message at the info level with optional arguments and sends a notification if a notifier is configured.
	Infof(msg string, args ...any)

	// Error logs a message at the error level with optional arguments.
	Error(msg string, args ...any)

	// Errorf logs a formatted message at the error level with optional arguments.
	Errorf(msg string, args ...any)

	// Warn logs a message at the warning level with optional arguments.
	Warn(msg string, args ...any)

	// Warnf logs a formatted message at the warning level with optional arguments.
	Warnf(msg string, args ...any)

	// Audit logs an audit message with the specified action and optional arguments.
	Audit(ctx context.Context, action string, args ...any)
}
