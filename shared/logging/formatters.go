package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

// Formatter defines the interface for log formatters
type Formatter interface {
	Format(entry *LogEntry) ([]byte, error)
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
}

// JSONFormatter formats logs as JSON
type JSONFormatter struct {
	PrettyPrint bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(prettyPrint bool) *JSONFormatter {
	return &JSONFormatter{PrettyPrint: prettyPrint}
}

// Format formats a log entry as JSON
func (f *JSONFormatter) Format(entry *LogEntry) ([]byte, error) {
	if f.PrettyPrint {
		return json.MarshalIndent(entry, "", "  ")
	}
	return json.Marshal(entry)
}

// TextFormatter formats logs as text
type TextFormatter struct {
	TimestampFormat string
	FullTimestamp   bool
}

// NewTextFormatter creates a new text formatter
func NewTextFormatter(timestampFormat string, fullTimestamp bool) *TextFormatter {
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	return &TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   fullTimestamp,
	}
}

// Format formats a log entry as text
func (f *TextFormatter) Format(entry *LogEntry) ([]byte, error) {
	timestamp := entry.Timestamp.Format(f.TimestampFormat)

	// Base format
	format := fmt.Sprintf("%s [%s] %s", timestamp, entry.Level, entry.Message)

	// Add caller if present
	if entry.Caller != "" {
		format += fmt.Sprintf(" (%s)", entry.Caller)
	}

	// Add fields if present
	if len(entry.Fields) > 0 {
		fieldsJSON, err := json.Marshal(entry.Fields)
		if err == nil {
			format += fmt.Sprintf(" %s", string(fieldsJSON))
		}
	}

	// Add stack trace if present
	if entry.Stack != "" {
		format += fmt.Sprintf("\n%s", entry.Stack)
	}

	return []byte(format), nil
}

// ConsoleFormatter formats logs for console output with colors
type ConsoleFormatter struct {
	TimestampFormat string
	FullTimestamp   bool
	ColorOutput     bool
}

// NewConsoleFormatter creates a new console formatter
func NewConsoleFormatter(timestampFormat string, fullTimestamp, colorOutput bool) *ConsoleFormatter {
	if timestampFormat == "" {
		timestampFormat = "15:04:05"
	}

	return &ConsoleFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   fullTimestamp,
		ColorOutput:     colorOutput,
	}
}

// Format formats a log entry for console output
func (f *ConsoleFormatter) Format(entry *LogEntry) ([]byte, error) {
	timestamp := entry.Timestamp.Format(f.TimestampFormat)

	// Color codes
	var levelColor, resetColor string
	if f.ColorOutput {
		resetColor = "\033[0m"
		switch entry.Level {
		case "DEBUG":
			levelColor = "\033[36m" // Cyan
		case "INFO":
			levelColor = "\033[32m" // Green
		case "WARN":
			levelColor = "\033[33m" // Yellow
		case "ERROR":
			levelColor = "\033[31m" // Red
		case "FATAL":
			levelColor = "\033[35m" // Magenta
		case "PANIC":
			levelColor = "\033[35m" // Magenta
		default:
			levelColor = "\033[37m" // White
		}
	}

	// Base format
	format := fmt.Sprintf("%s %s[%s]%s %s",
		timestamp, levelColor, entry.Level, resetColor, entry.Message)

	// Add caller if present
	if entry.Caller != "" {
		format += fmt.Sprintf(" %s(%s)%s", levelColor, entry.Caller, resetColor)
	}

	// Add fields if present
	if len(entry.Fields) > 0 {
		fieldsJSON, err := json.Marshal(entry.Fields)
		if err == nil {
			format += fmt.Sprintf(" %s%s%s", levelColor, string(fieldsJSON), resetColor)
		}
	}

	// Add stack trace if present
	if entry.Stack != "" {
		format += fmt.Sprintf("\n%s", entry.Stack)
	}

	return []byte(format), nil
}

// CustomFormatter allows custom formatting logic
type CustomFormatter struct {
	FormatFunc func(entry *LogEntry) ([]byte, error)
}

// NewCustomFormatter creates a new custom formatter
func NewCustomFormatter(formatFunc func(entry *LogEntry) ([]byte, error)) *CustomFormatter {
	return &CustomFormatter{FormatFunc: formatFunc}
}

// Format formats a log entry using the custom function
func (f *CustomFormatter) Format(entry *LogEntry) ([]byte, error) {
	return f.FormatFunc(entry)
}

// LogEntryBuilder helps build log entries
type LogEntryBuilder struct {
	entry *LogEntry
}

// NewLogEntryBuilder creates a new log entry builder
func NewLogEntryBuilder() *LogEntryBuilder {
	return &LogEntryBuilder{
		entry: &LogEntry{
			Timestamp: time.Now(),
			Fields:    make(map[string]interface{}),
		},
	}
}

// WithTimestamp sets the timestamp
func (b *LogEntryBuilder) WithTimestamp(timestamp time.Time) *LogEntryBuilder {
	b.entry.Timestamp = timestamp
	return b
}

// WithLevel sets the log level
func (b *LogEntryBuilder) WithLevel(level string) *LogEntryBuilder {
	b.entry.Level = level
	return b
}

// WithMessage sets the message
func (b *LogEntryBuilder) WithMessage(message string) *LogEntryBuilder {
	b.entry.Message = message
	return b
}

// WithField adds a field
func (b *LogEntryBuilder) WithField(key string, value interface{}) *LogEntryBuilder {
	b.entry.Fields[key] = value
	return b
}

// WithFields adds multiple fields
func (b *LogEntryBuilder) WithFields(fields map[string]interface{}) *LogEntryBuilder {
	for key, value := range fields {
		b.entry.Fields[key] = value
	}
	return b
}

// WithCaller sets the caller
func (b *LogEntryBuilder) WithCaller(caller string) *LogEntryBuilder {
	b.entry.Caller = caller
	return b
}

// WithStack sets the stack trace
func (b *LogEntryBuilder) WithStack(stack string) *LogEntryBuilder {
	b.entry.Stack = stack
	return b
}

// Build builds the log entry
func (b *LogEntryBuilder) Build() *LogEntry {
	return b.entry
}

// ZapCoreFormatter adapts a Formatter to work with zapcore
type ZapCoreFormatter struct {
	formatter Formatter
}

// NewZapCoreFormatter creates a new zapcore formatter
func NewZapCoreFormatter(formatter Formatter) *ZapCoreFormatter {
	return &ZapCoreFormatter{formatter: formatter}
}

// Format formats a zapcore entry
func (f *ZapCoreFormatter) Format(entry zapcore.Entry, fields []zapcore.Field) ([]byte, error) {
	// Convert zapcore entry to LogEntry
	logEntry := &LogEntry{
		Timestamp: entry.Time,
		Level:     entry.Level.CapitalString(),
		Message:   entry.Message,
		Caller:    entry.Caller.String(),
		Stack:     entry.Stack,
		Fields:    make(map[string]interface{}),
	}

	// Convert zapcore fields to map
	for _, field := range fields {
		logEntry.Fields[field.Key] = field.Interface
	}

	return f.formatter.Format(logEntry)
}
