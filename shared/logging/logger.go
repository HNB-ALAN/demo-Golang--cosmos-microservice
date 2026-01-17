package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/usc-platform/shared/config"
)

// Logger represents a structured logger
type Logger struct {
	zap *zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(serviceName string, cfg config.LogConfig) *Logger {
	// Configure log level
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	case "panic":
		level = zap.PanicLevel
	}

	// Configure encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	} else {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}

	// Configure output
	var output zapcore.WriteSyncer
	switch cfg.Output {
	case "stdout":
		output = zapcore.AddSync(os.Stdout)
	case "stderr":
		output = zapcore.AddSync(os.Stderr)
	case "file":
		if cfg.Filename != "" {
			file, err := os.OpenFile(cfg.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				// Fallback to stdout if file cannot be opened
				output = zapcore.AddSync(os.Stdout)
			} else {
				output = zapcore.AddSync(file)
			}
		} else {
			output = zapcore.AddSync(os.Stdout)
		}
	default:
		output = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, output, level)

	// Create logger with service name
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Add service name to all logs
	logger = logger.With(zap.String("service", serviceName))

	return &Logger{zap: logger}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(msg string, fields ...Field) {
	l.zap.Panic(msg, fields...)
}

// With creates a new logger with additional fields
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// GetZapLogger returns the underlying zap logger
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zap
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	case "fatal":
		zapLevel = zap.FatalLevel
	case "panic":
		zapLevel = zap.PanicLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// Create new logger with updated level
	// Simplified level change - just create a new logger with the new level
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	newLogger, _ := config.Build()
	l.zap = newLogger
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.zap.Core().Enabled(zap.DebugLevel)
}

// IsInfoEnabled returns true if info logging is enabled
func (l *Logger) IsInfoEnabled() bool {
	return l.zap.Core().Enabled(zap.InfoLevel)
}

// IsWarnEnabled returns true if warn logging is enabled
func (l *Logger) IsWarnEnabled() bool {
	return l.zap.Core().Enabled(zap.WarnLevel)
}

// IsErrorEnabled returns true if error logging is enabled
func (l *Logger) IsErrorEnabled() bool {
	return l.zap.Core().Enabled(zap.ErrorLevel)
}
