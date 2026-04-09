# logger-wrapper

A small Go package built on top of `log/slog` to provide structured logging, consistent environment metadata, audit events, and optional notifier integration.

## Table of contents

- [logger-wrapper](#logger-wrapper)
  - [Table of contents](#table-of-contents)
  - [Installation](#installation)
  - [Package import](#package-import)
  - [Quick start](#quick-start)
  - [Configuration](#configuration)
  - [Available methods](#available-methods)
    - [Logging methods](#logging-methods)
    - [Formatted logging methods](#formatted-logging-methods)
    - [Audit and context methods](#audit-and-context-methods)
    - [Other methods](#other-methods)
  - [Structured logging examples](#structured-logging-examples)
  - [Notifier support](#notifier-support)
    - [Using Gotify](#using-gotify)
    - [Using your own notifier](#using-your-own-notifier)
  - [Accessing the underlying `slog.Logger`](#accessing-the-underlying-sloglogger)
  - [Development](#development)
  - [Requirements](#requirements)
  - [License](#license)

## Installation

```bash
go get github.com/dipievil/logger-wrapper
```

## Package import

```go
import "github.com/dipievil/logger-wrapper/logging"
```

## Quick start

```go
package main

import "github.com/dipievil/logger-wrapper/logging"

func main() {
  cfg := logging.LoggerConfig{
    LogLevel:     "info",
    BuildVersion: "1.0.0",
    Environment:  "production",
  }

  logger := logging.NewLoggerWrapper(cfg)

  logger.Info("application started", "service", "api")  
  logger.Warn("cache miss", "key", "user:42")
  logger.Error("request failed", "status", 500)
}
```

By default, logs are emitted as JSON to `stdout`.

## Configuration

Use `logging.LoggerConfig` to configure the logger instance:

| Field | Description | Default |
| --- | --- | --- |
| `LogLevel` | Log verbosity: `debug`, `info`, `warn`, `error` | `debug` |
| `BuildVersion` | Version metadata added to every log line | `dev` |
| `Environment` | Environment metadata added to every log line | `local` |

You can also start from the defaults:

```go
cfg := logging.NewLoggerConfig()
cfg.LogLevel = "warn"

logger := logging.NewLoggerWrapper(cfg)
```

## Available methods

The wrapper exposes the following methods:

### Logging methods

- **Debug(msg string, args ...any)**
    Used for detailed debugging information, typically only useful during development or troubleshooting.
- **Info(msg string, args ...any)**
    Used for general informational messages that highlight the progress of the application at a coarse-grained level.
- **Warn(msg string, args ...any)**
    Used for potentially harmful situations or important events that are not errors but may require attention.
- **Error(msg string, args ...any)**
    Used for error events that might still allow the application to continue running.
- **Debugf(msg string, args ...any)**
    Used for detailed debugging information with formatted output, typically only useful during development or troubleshooting.

### Formatted logging methods

- **Infof(msg string, args ...any)**
    Used for general informational messages with formatted output that highlight the progress of the application at a coarse-grained level.
- **Warnf(msg string, args ...any)**
    Used for potentially harmful situations or important events with formatted output that are not errors but may require attention.
- **Errorf(msg string, args ...any)**
    Used for error events with formatted output that might still allow the application to continue running.

### Audit and context methods

- **Audit(ctx context.Context, action string, args ...any)**
    Used for logging audit events with a specific action and structured context. The `ctx` can be used to pass additional metadata or correlation IDs.

### Other methods

- **With(args ...any) *slog.Logger***
    Used to create a child logger with additional context fields. This is useful for adding request-specific or component-specific metadata to logs without modifying the global logger configuration.

- **Base() *slog.Logger***
    Used to access the underlying `slog.Logger` instance for advanced use cases or direct logging. Use this when you need to leverage features of `slog` that are not exposed by the wrapper or when you want to create custom loggers with specific configurations.

## Structured logging examples

- ***Structured context***

```go
logger.Info("user authenticated", "user_id", 42, "role", "admin")
```

- ***Audit events***

```go
ctx := context.Background()
logger.Audit(ctx, "user.updated", "user_id", 42, "level", "info")
```

## Notifier support

The package supports optional notifications through `WithNotifier`.

### Using Gotify

```go
package main

import "github.com/dipievil/logger-wrapper/logging"

func main() {
  cfg := logging.NewLoggerConfig()

  notifier := &logging.GotifyService{
    URL:   "https://gotify.example.com",
    Token: "your-token",
    Title: "My Service",
  }

  logger := logging.NewLoggerWrapper(cfg, logging.WithNotifier(notifier))
  logger.Info("deployment finished")
}
```

> `Info()` and `Infof()` trigger notifications when a notifier is configured.

### Using your own notifier

Implement the `Notifier` interface from `github.com/dipievil/logger-wrapper/logging/interfaces`:

```go
package main

import (
  "fmt"

  "github.com/dipievil/logger-wrapper/logging"
  logginginterfaces "github.com/dipievil/logger-wrapper/logging/interfaces"
)

type StdoutNotifier struct{}

var _ logginginterfaces.Notifier = (*StdoutNotifier)(nil)

func (n *StdoutNotifier) Notify(message string) error {
  fmt.Println("notification:", message)
  return nil
}

func main() {
  logger := logging.NewLoggerWrapper(
    logging.NewLoggerConfig(),
    logging.WithNotifier(&StdoutNotifier{}),
  )

  logger.Info("background job completed")
}
```

Notifications are dispatched asynchronously so they do not block the main logging flow.

## Accessing the underlying `slog.Logger`

If you need direct access to the standard logger:

```go
base := logger.Base()
base.Info("native slog call")
```

Or create a child logger with additional fields:

```go
requestLogger := logger.With("request_id", "abc-123")
requestLogger.Info("processing request")
```

## Development

Run the test suite locally with:

```bash
go test ./...
```

## Requirements

- Go `1.26+`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.