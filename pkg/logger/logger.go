package logger

import (
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Config holds logger configuration
type Config struct {
	Level  LogLevel
	Format string // "text" or "json"
}

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	level LogLevel
}

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	var level slog.Level
	switch config.Level {
	case DebugLevel:
		level = slog.LevelDebug
	case InfoLevel:
		level = slog.LevelInfo
	case WarnLevel:
		level = slog.LevelWarn
	case ErrorLevel:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
		level:  config.Level,
	}
}

// ParseLevel parses a string level into LogLevel
func ParseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= DebugLevel
}

// WithComponent creates a new logger with a component name prefix
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.Logger.With("component", component),
		level:  l.level,
	}
}

// DebugContext logs a debug message with context
func (l *Logger) DebugContext(msg string, args ...any) {
	l.Debug(msg, args...)
}

// InfoContext logs an info message with context
func (l *Logger) InfoContext(msg string, args ...any) {
	l.Info(msg, args...)
}

// WarnContext logs a warning message with context
func (l *Logger) WarnContext(msg string, args ...any) {
	l.Warn(msg, args...)
}

// ErrorContext logs an error message with context
func (l *Logger) ErrorContext(msg string, args ...any) {
	l.Error(msg, args...)
}