package logging

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

type MockNotifier struct {
	notified chan string
}

func newMockNotifier() *MockNotifier {
	return &MockNotifier{notified: make(chan string, 1)}
}

func (m *MockNotifier) Notify(message string) error {
	m.notified <- message
	return nil
}

type MockLogger struct {
	base *slog.Logger
}

func TestLoggerWithNotifier(t *testing.T) {
	mockNotifier := newMockNotifier()

	logger := NewLoggerWrapper(NewLoggerConfig(), WithNotifier(mockNotifier))

	logger.Info("Test log message")

	select {
	case got := <-mockNotifier.notified:
		if got != "Test log message" {
			t.Errorf("Expected notifier to receive 'Test log message', but got '%s'", got)
		} else {
			t.Log("Notifier received the expected message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for notifier message")
	}
}

func TestLoggerWithCustomNotifier(t *testing.T) {
	mockNotifier := newMockNotifier()

	logger := NewLoggerWrapper(NewLoggerConfig(), WithNotifier(mockNotifier))

	logger.Infof("Test log message with args: %d", 42)

	select {
	case got := <-mockNotifier.notified:
		expected := "Test log message with args: 42"
		if got != expected {
			t.Errorf("Expected notifier to receive '%s', but got '%s'", expected, got)
		} else {
			t.Log("Notifier received the expected formatted message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for notifier message")
	}
}

func TestLoggerWithInvalidNotifier(t *testing.T) {
	logger := NewLoggerWrapper(NewLoggerConfig(), WithNotifier(nil))

	logger.Info("Test log message with invalid notifier")

	// Since the notifier is nil, we expect no panic and no notification sent.
	t.Log("Logger handled nil notifier without crashing")
}

func TestLoggerWithoutNotifier(t *testing.T) {
	logger := NewLoggerWrapper(NewLoggerConfig())

	logger.base.Info("Test log message without notifier")
}

func TestLoggerWithCustomLevel(t *testing.T) {
	loggerConfig := NewLoggerConfig()
	loggerConfig.LogLevel = "debug"

	logger := NewLoggerWrapper(loggerConfig)

	ctx := context.Background()

	if logger.base.Handler().Enabled(ctx, slog.LevelDebug) {
		t.Log("Debug level is enabled as expected")
	} else {
		t.Error("Expected debug level to be enabled, but it is not")
	}
}

func TestLoggerWithInvalidLevel(t *testing.T) {
	loggerConfig := NewLoggerConfig()
	loggerConfig.LogLevel = "invalid"

	logger := NewLoggerWrapper(loggerConfig)

	ctx := context.Background()

	if logger.base.Handler().Enabled(ctx, slog.LevelInfo) {
		t.Log("Info level is enabled as expected for invalid log level")
	} else {
		t.Error("Expected info level to be enabled for invalid log level, but it is not")
	}
}
