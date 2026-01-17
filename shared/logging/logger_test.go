package logging

import (
	"testing"

	"github.com/usc-platform/shared/config"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		service string
		config  config.LogConfig
		wantErr bool
	}{
		{
			name:    "default config",
			service: "test-service",
			config:  config.LogConfig{},
			wantErr: false,
		},
		{
			name:    "debug level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "debug",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name:    "info level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "info",
				Format: "console",
			},
			wantErr: false,
		},
		{
			name:    "warn level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "warn",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name:    "error level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "error",
				Format: "console",
			},
			wantErr: false,
		},
		{
			name:    "fatal level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "fatal",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name:    "panic level",
			service: "test-service",
			config: config.LogConfig{
				Level:  "panic",
				Format: "console",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.service, tt.config)
			if logger == nil {
				t.Errorf("NewLogger() returned nil")
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test basic info logging
	logger.Info("test message")
	logger.Info("test message with fields", String("key", "value"))
}

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{Level: "debug"})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test debug logging
	logger.Debug("debug message")
	logger.Debug("debug message with fields", String("key", "value"))
}

func TestLogger_Warn(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test warn logging
	logger.Warn("warn message")
	logger.Warn("warn message with fields", String("key", "value"))
}

func TestLogger_Error(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test error logging
	logger.Error("error message")
	logger.Error("error message with fields", String("key", "value"))
}

func TestLogger_Fatal(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test fatal logging (but don't actually call Fatal as it would exit)
	// logger.Fatal("fatal message")
}

func TestLogger_Panic(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test panic logging (but don't actually call Panic as it would panic)
	// logger.Panic("panic message")
}

func TestLogger_With(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test with fields
	childLogger := logger.With(String("parent", "value"))
	if childLogger == nil {
		t.Errorf("With() returned nil")
	}

	childLogger.Info("child message")
}

func TestLogger_Sync(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test sync
	err := logger.Sync()
	if err != nil {
		t.Errorf("Sync() error = %v", err)
	}
}

func TestLogger_Fields(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	// Test various field types
	logger.Info("test with various fields",
		String("string", "value"),
		Int("int", 42),
		Int64("int64", 123456789),
		Float64("float64", 3.14),
		Bool("bool", true),
	)
}

func TestLogger_JSONFormat(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{
		Level:  "info",
		Format: "json",
	})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	logger.Info("json format test", String("key", "value"))
}

func TestLogger_ConsoleFormat(t *testing.T) {
	logger := NewLogger("test-service", config.LogConfig{
		Level:  "info",
		Format: "console",
	})
	if logger == nil {
		t.Fatalf("NewLogger() returned nil")
	}

	logger.Info("console format test", String("key", "value"))
}
