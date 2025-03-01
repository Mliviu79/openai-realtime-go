// Package logger provides logging functionality for the OpenAI Realtime API client.
// It defines a flexible logging interface with multiple implementations to suit different needs.
//
// The package offers:
//   - A standard Logger interface with common logging methods (Debugf, Infof, Warnf, Errorf)
//   - Field-based logging for structured logging with key-value pairs
//   - Multiple implementations including ZeroLogger (using zerolog) and NopLogger (no-op logger)
//   - Simple configuration options for controlling log levels, formatting, and output
//
// Example usage:
//
//	// Create a zero logger with debug level enabled
//	zeroLogger := logger.NewZeroLogger(logger.LoggerOptions{
//		Level: zerolog.DebugLevel,
//	})
//
//	// Log a debug message with a field
//	zeroLogger.WithField("session_id", "sess_123").Debugf("Connected to session")
//
//	// Or use the default logger
//	logger.Default.Infof("Using default logger")
//
// For performance-critical code or when logging is not needed, use the NopLogger:
//
//	client.SetLogger(logger.Nop)
//
// The logger package is designed to be easily extensible if additional logging backends are needed.
package logger

import (
	"io"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Logger defines the logging interface used throughout the application
type Logger interface {
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

// ZeroLogger implements Logger using zerolog
type ZeroLogger struct {
	log zerolog.Logger
	// Cache for file paths to avoid repeated processing
	pathCache *sync.Map
	// Fields to include with every log message
	fields map[string]any
}

// LoggerOptions contains configuration options for creating a logger
type LoggerOptions struct {
	Level      zerolog.Level
	Output     io.Writer
	TimeFormat string
	NoColor    bool
	CallerInfo bool
}

// DefaultLoggerOptions returns the default logger options
func DefaultLoggerOptions() LoggerOptions {
	return LoggerOptions{
		Level:      zerolog.DebugLevel,
		Output:     os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		CallerInfo: true,
	}
}

// NewZeroLogger creates a new ZeroLogger with the specified options
func NewZeroLogger(opts LoggerOptions) *ZeroLogger {
	// Set up the output writer
	var output io.Writer = opts.Output
	if output == nil {
		output = os.Stdout
	}

	// Add color if enabled and output is a terminal
	if !opts.NoColor && isTerminal(output) {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: opts.TimeFormat,
			NoColor:    opts.NoColor,
		}
	}

	// Create the logger
	logger := zerolog.New(output).
		Level(opts.Level).
		With().
		Timestamp().
		Logger()

	return &ZeroLogger{
		log:       logger,
		pathCache: &sync.Map{},
		fields:    make(map[string]any),
	}
}

// isTerminal checks if the writer is a terminal
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return f == os.Stdout || f == os.Stderr
	}
	return false
}

// getCallerInfo returns the caller's file and line number
// Uses a cache to avoid repeated path processing for the same files
func (l *ZeroLogger) getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(3) // Skip 3 frames to get the actual caller
	if !ok {
		return "unknown", 0
	}

	// Check if we've already processed this path
	if cachedPath, ok := l.pathCache.Load(file); ok {
		return cachedPath.(string), line
	}

	// Only use the last two parts of the path for brevity
	dir, fileName := filepath.Split(file)
	parentDir := filepath.Base(dir)
	shortPath := filepath.Join(parentDir, fileName)

	// Store in cache for future use
	l.pathCache.Store(file, shortPath)

	return shortPath, line
}

// createEvent creates a new log event with fields and caller info
func (l *ZeroLogger) createEvent(level zerolog.Level) *zerolog.Event {
	event := l.log.WithLevel(level)

	// Add fields
	for k, v := range l.fields {
		event = event.Interface(k, v)
	}

	// Add caller info
	file, line := l.getCallerInfo()
	event = event.Str("file", file).Int("line", line)

	return event
}

func (l *ZeroLogger) Debugf(format string, v ...any) {
	l.createEvent(zerolog.DebugLevel).Msgf(format, v...)
}

func (l *ZeroLogger) Infof(format string, v ...any) {
	l.createEvent(zerolog.InfoLevel).Msgf(format, v...)
}

func (l *ZeroLogger) Warnf(format string, v ...any) {
	l.createEvent(zerolog.WarnLevel).Msgf(format, v...)
}

func (l *ZeroLogger) Errorf(format string, v ...any) {
	l.createEvent(zerolog.ErrorLevel).Msgf(format, v...)
}

// WithField returns a new logger with the field added to the logger's context
func (l *ZeroLogger) WithField(key string, value any) Logger {
	newLogger := &ZeroLogger{
		log:       l.log,
		pathCache: l.pathCache,
		fields:    make(map[string]any, len(l.fields)+1),
	}

	// Use maps.Copy instead of manual loop
	maps.Copy(newLogger.fields, l.fields)

	// Add new field
	newLogger.fields[key] = value

	return newLogger
}

// WithFields returns a new logger with the fields added to the logger's context
func (l *ZeroLogger) WithFields(fields map[string]any) Logger {
	newLogger := &ZeroLogger{
		log:       l.log,
		pathCache: l.pathCache,
		fields:    make(map[string]any, len(l.fields)+len(fields)),
	}

	// Use maps.Copy instead of manual loops
	maps.Copy(newLogger.fields, l.fields)
	maps.Copy(newLogger.fields, fields)

	return newLogger
}

// NopLogger is a logger that does nothing
type NopLogger struct{}

func (l NopLogger) Debugf(format string, v ...any)          {}
func (l NopLogger) Infof(format string, v ...any)           {}
func (l NopLogger) Warnf(format string, v ...any)           {}
func (l NopLogger) Errorf(format string, v ...any)          {}
func (l NopLogger) WithField(key string, value any) Logger  { return l }
func (l NopLogger) WithFields(fields map[string]any) Logger { return l }

// Default loggers
var (
	Nop     = &NopLogger{}
	Default = NewZeroLogger(DefaultLoggerOptions())
)

// SetGlobalLevel sets the global logging level
func SetGlobalLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}
