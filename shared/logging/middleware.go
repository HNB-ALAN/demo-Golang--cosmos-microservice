package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// LoggingLoggingMiddleware provides logging middleware functionality
type LoggingLoggingMiddleware struct {
	logger *Logger
}

// NewLoggingLoggingMiddleware creates a new logging middleware
func NewLoggingLoggingMiddleware(logger *Logger) *LoggingLoggingMiddleware {
	return &LoggingLoggingMiddleware{logger: logger}
}

// LoggingContextKey is the key used to store logger in context
type LoggingContextKey string

const (
	LoggerKey LoggingContextKey = "logger"
)

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// FromContext retrieves a logger from the context
func FromContext(ctx context.Context) (*Logger, bool) {
	logger, ok := ctx.Value(LoggerKey).(*Logger)
	return logger, ok
}

// MustFromContext retrieves a logger from the context or panics
func MustFromContext(ctx context.Context) *Logger {
	logger, ok := FromContext(ctx)
	if !ok {
		panic("logger not found in context")
	}
	return logger
}

// RequestLoggingLoggingMiddleware creates a middleware that logs HTTP requests
func (m *LoggingLoggingMiddleware) RequestLoggingLoggingMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract request information from context
			requestID := m.getRequestID(ctx)
			method := m.getMethod(ctx)
			path := m.getPath(ctx)

			// Create request-scoped logger
			requestLogger := m.logger.With(
				RequestID(requestID),
				RequestMethod(method),
				RequestPath(path),
			)

			// Add logger to context
			ctx = WithLogger(ctx, requestLogger)

			// Log request start
			requestLogger.Info("request started")

			// Execute next handler
			next(ctx)

			// Log request completion
			duration := time.Since(start)
			statusCode := m.getStatusCode(ctx)

			requestLogger.Info("request completed",
				Duration("duration", duration),
				Int("status_code", statusCode),
			)
		}
	}
}

// DatabaseLoggingLoggingMiddleware creates a middleware that logs database operations
func (m *LoggingLoggingMiddleware) DatabaseLoggingLoggingMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract database information from context
			dbName := m.getDatabaseName(ctx)
			tableName := m.getTableName(ctx)
			queryType := m.getQueryType(ctx)

			// Create database-scoped logger
			dbLogger := m.logger.With(
				DatabaseName(dbName),
				TableName(tableName),
				QueryType(queryType),
			)

			// Add logger to context
			ctx = WithLogger(ctx, dbLogger)

			// Execute next handler
			next(ctx)

			// Log database operation completion
			duration := time.Since(start)
			rowsAffected := m.getRowsAffected(ctx)

			dbLogger.Info("database operation completed",
				QueryDuration(duration),
				Int64("rows_affected", rowsAffected),
			)
		}
	}
}

// ErrorLoggingLoggingMiddleware creates a middleware that logs errors
func (m *LoggingLoggingMiddleware) ErrorLoggingLoggingMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					// Log panic
					m.logger.Error("panic recovered",
						Any("panic", r),
						Stack("stack"),
					)
					panic(r) // Re-panic
				}
			}()

			// Execute next handler
			next(ctx)

			// Check for errors in context
			if err := m.getError(ctx); err != nil {
				errorCode := m.getErrorCode(ctx)
				errorType := m.getErrorType(ctx)

				m.logger.Error("error occurred",
					Error(err),
					ErrorCode(errorCode),
					ErrorType(errorType),
				)
			}
		}
	}
}

// PerformanceLoggingLoggingMiddleware creates a middleware that logs performance metrics
func (m *LoggingLoggingMiddleware) PerformanceLoggingLoggingMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Execute next handler
			next(ctx)

			// Log performance metrics
			duration := time.Since(start)
			memoryUsage := m.getMemoryUsage(ctx)
			cpuUsage := m.getCPUUsage(ctx)

			m.logger.Info("performance metrics",
				Duration("duration", duration),
				MemoryUsage(memoryUsage),
				CPUUsage(cpuUsage),
			)
		}
	}
}

// Helper methods to extract information from context
// These would typically be implemented based on your specific context structure

func (m *LoggingLoggingMiddleware) getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getMethod(ctx context.Context) string {
	if method, ok := ctx.Value("method").(string); ok {
		return method
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getPath(ctx context.Context) string {
	if path, ok := ctx.Value("path").(string); ok {
		return path
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getStatusCode(ctx context.Context) int {
	if code, ok := ctx.Value("status_code").(int); ok {
		return code
	}
	return 0
}

func (m *LoggingLoggingMiddleware) getDatabaseName(ctx context.Context) string {
	if db, ok := ctx.Value("database_name").(string); ok {
		return db
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getTableName(ctx context.Context) string {
	if table, ok := ctx.Value("table_name").(string); ok {
		return table
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getQueryType(ctx context.Context) string {
	if queryType, ok := ctx.Value("query_type").(string); ok {
		return queryType
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getRowsAffected(ctx context.Context) int64 {
	if rows, ok := ctx.Value("rows_affected").(int64); ok {
		return rows
	}
	return 0
}

func (m *LoggingLoggingMiddleware) getError(ctx context.Context) error {
	if err, ok := ctx.Value("error").(error); ok {
		return err
	}
	return nil
}

func (m *LoggingLoggingMiddleware) getErrorCode(ctx context.Context) string {
	if code, ok := ctx.Value("error_code").(string); ok {
		return code
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getErrorType(ctx context.Context) string {
	if errorType, ok := ctx.Value("error_type").(string); ok {
		return errorType
	}
	return "unknown"
}

func (m *LoggingLoggingMiddleware) getMemoryUsage(ctx context.Context) int64 {
	if usage, ok := ctx.Value("memory_usage").(int64); ok {
		return usage
	}
	return 0
}

func (m *LoggingLoggingMiddleware) getCPUUsage(ctx context.Context) float64 {
	if usage, ok := ctx.Value("cpu_usage").(float64); ok {
		return usage
	}
	return 0.0
}

// StructuredLoggerFromContext creates a structured logger from context
func StructuredLoggerFromContext(ctx context.Context) *StructuredLogger {
	logger, ok := FromContext(ctx)
	if !ok {
		// Fallback to default logger
		logger = &Logger{zap: zap.NewNop()}
	}
	return NewStructuredLogger(logger)
}
