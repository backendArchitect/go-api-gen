package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestLoggerCreation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "debug level text format",
			config: Config{
				Level:  DebugLevel,
				Format: "text",
			},
		},
		{
			name: "info level json format",
			config: Config{
				Level:  InfoLevel,
				Format: "json",
			},
		},
		{
			name: "default config",
			config: Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			if logger == nil {
				t.Error("Expected logger to be created, got nil")
			}
			if logger.Logger == nil {
				t.Error("Expected slog.Logger to be created, got nil")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"WARN", WarnLevel},
		{"error", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"unknown", InfoLevel}, // default
		{"", InfoLevel},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsDebugEnabled(t *testing.T) {
	debugLogger := New(Config{Level: DebugLevel})
	infoLogger := New(Config{Level: InfoLevel})

	if !debugLogger.IsDebugEnabled() {
		t.Error("Expected debug logger to have debug enabled")
	}

	if infoLogger.IsDebugEnabled() {
		t.Error("Expected info logger to have debug disabled")
	}
}

func TestWithComponent(t *testing.T) {
	logger := New(Config{Level: DebugLevel})
	componentLogger := logger.WithComponent("test-component")

	if componentLogger == nil {
		t.Error("Expected component logger to be created")
	}

	if componentLogger.level != logger.level {
		t.Error("Expected component logger to inherit level")
	}
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := &Logger{
		Logger: slog.New(handler),
		level:  DebugLevel,
	}

	// Test different log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Error("Expected debug message in output")
	}
	if !strings.Contains(output, "info message") {
		t.Error("Expected info message in output")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Expected warn message in output")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Expected error message in output")
	}
}

func TestJSONLogging(t *testing.T) {
	var buf bytes.Buffer

	// Create a JSON logger that writes to our buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := &Logger{
		Logger: slog.New(handler),
		level:  InfoLevel,
	}

	logger.Info("test message", "key", "value")

	output := buf.String()
	
	// Verify it's valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v", err)
	}

	// Check for expected fields
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected msg to be 'test message', got %v", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key to be 'value', got %v", logEntry["key"])
	}
}