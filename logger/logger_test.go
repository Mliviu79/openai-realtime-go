package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNopLogger(t *testing.T) {
	// Create a NopLogger
	logger := &NopLogger{}

	// Verify that the methods don't panic
	logger.Debugf("test debug message")
	logger.Infof("test info message")
	logger.Warnf("test warning message")
	logger.Errorf("test error message")

	// Verify that WithField returns a logger (no need to check exact instance)
	result := logger.WithField("key", "value")
	if result == nil {
		t.Error("Expected WithField to return a logger, got nil")
	}

	// Verify that WithFields returns a logger (no need to check exact instance)
	result = logger.WithFields(map[string]any{"key": "value"})
	if result == nil {
		t.Error("Expected WithFields to return a logger, got nil")
	}
}

func TestZeroLoggerBasic(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create options with the buffer as output
	opts := LoggerOptions{
		Level:      zerolog.DebugLevel,
		Output:     &buf,
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		NoColor:    true,
		CallerInfo: true,
	}

	// Create a logger
	logger := NewZeroLogger(opts)

	// Log a debug message
	logger.Debugf("test debug message")

	// Verify that the output contains the expected content
	output := buf.String()
	if !strings.Contains(output, "test debug message") {
		t.Errorf("Expected output to contain 'test debug message', got: %s", output)
	}
	if !strings.Contains(output, "debug") {
		t.Errorf("Expected output to contain 'debug' level, got: %s", output)
	}
}

func TestZeroLoggerLevels(t *testing.T) {
	// Test each log level
	levels := []struct {
		level    zerolog.Level
		logFunc  func(l *ZeroLogger, msg string)
		expected string
	}{
		{zerolog.DebugLevel, func(l *ZeroLogger, msg string) { l.Debugf(msg) }, "debug"},
		{zerolog.InfoLevel, func(l *ZeroLogger, msg string) { l.Infof(msg) }, "info"},
		{zerolog.WarnLevel, func(l *ZeroLogger, msg string) { l.Warnf(msg) }, "warn"},
		{zerolog.ErrorLevel, func(l *ZeroLogger, msg string) { l.Errorf(msg) }, "error"},
	}

	for _, test := range levels {
		t.Run(test.expected, func(t *testing.T) {
			// Create a buffer to capture log output
			var buf bytes.Buffer

			// Create options with the buffer as output
			opts := LoggerOptions{
				Level:      test.level,
				Output:     &buf,
				NoColor:    true,
				CallerInfo: false,
			}

			// Create a logger
			logger := NewZeroLogger(opts)

			// Log a message
			test.logFunc(logger, "test message")

			// Verify the output
			output := buf.String()
			if !strings.Contains(output, "test message") {
				t.Errorf("Expected output to contain 'test message', got: %s", output)
			}
			if !strings.Contains(output, test.expected) {
				t.Errorf("Expected output to contain '%s' level, got: %s", test.expected, output)
			}
		})
	}
}

func TestZeroLoggerWithField(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create options with the buffer as output
	opts := LoggerOptions{
		Level:      zerolog.DebugLevel,
		Output:     &buf,
		NoColor:    true,
		CallerInfo: false,
	}

	// Create a logger and add a field
	logger := NewZeroLogger(opts)
	loggerWithField := logger.WithField("testKey", "testValue")

	// Log a message
	loggerWithField.Debugf("test message with field")

	// Parse the JSON output
	var logData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	// Verify the field is present
	value, ok := logData["testKey"]
	if !ok {
		t.Error("Expected log to contain field 'testKey' but it wasn't found")
	} else if value != "testValue" {
		t.Errorf("Expected 'testKey' to be 'testValue', got %v", value)
	}
}

func TestZeroLoggerWithFields(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create options with the buffer as output
	opts := LoggerOptions{
		Level:      zerolog.DebugLevel,
		Output:     &buf,
		NoColor:    true,
		CallerInfo: false,
	}

	// Create a logger and add fields
	logger := NewZeroLogger(opts)
	fields := map[string]any{
		"field1": "value1",
		"field2": float64(42), // JSON unmarshals numbers as float64
	}
	loggerWithFields := logger.WithFields(fields)

	// Log a message
	loggerWithFields.Debugf("test message with fields")

	// Parse the JSON output
	var logData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	// Verify the fields are present
	if value, ok := logData["field1"]; !ok {
		t.Error("Expected log to contain field 'field1' but it wasn't found")
	} else if value != "value1" {
		t.Errorf("Expected 'field1' to be 'value1', got %v", value)
	}

	if value, ok := logData["field2"]; !ok {
		t.Error("Expected log to contain field 'field2' but it wasn't found")
	} else if value != float64(42) {
		t.Errorf("Expected 'field2' to be a float64 with value 42, got %v (type %T)", value, value)
	}
}
