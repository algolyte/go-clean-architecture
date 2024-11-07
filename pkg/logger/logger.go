package logger

import "context"

// Level represents the severity level of a log entry
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Fields represents a map of field keys and values to be included in log entries
type Fields map[string]interface{}

// Logger defines the interface that any logger implementation must satisfy
type Logger interface {
	// Debug logs a message at debug level
	Debug(ctx context.Context, msg string, fields ...Fields)

	// Info logs a message at info level
	Info(ctx context.Context, msg string, fields ...Fields)

	// Warn logs a message at warn level
	Warn(ctx context.Context, msg string, fields ...Fields)

	// Error logs a message at error level
	Error(ctx context.Context, msg string, fields ...Fields)

	// Fatal logs a message at fatal level and terminates the program
	Fatal(ctx context.Context, msg string, fields ...Fields)

	// WithFields creates a new logger with the given fields added to it
	WithFields(fields Fields) Logger

	// WithField creates a new logger with a single field added to it
	WithField(key string, value interface{}) Logger

	// Sync ensures all buffered logs are written
	Sync() error
}
