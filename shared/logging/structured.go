package logging

import (
	"context"
	"time"
)

// StructuredLogger provides structured logging capabilities
type StructuredLogger struct {
	logger *Logger
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(logger *Logger) *StructuredLogger {
	return &StructuredLogger{logger: logger}
}

// LogRequest logs an HTTP/gRPC request
func (s *StructuredLogger) LogRequest(ctx context.Context, method, path string, duration time.Duration, statusCode int, fields ...Field) {
	allFields := []Field{
		RequestMethod(method),
		RequestPath(path),
		Duration("duration", duration),
		Int("status_code", statusCode),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("request completed", allFields...)
}

// LogDatabaseQuery logs a database query
func (s *StructuredLogger) LogDatabaseQuery(ctx context.Context, dbName, tableName, queryType string, duration time.Duration, rowsAffected int64, fields ...Field) {
	allFields := []Field{
		DatabaseName(dbName),
		TableName(tableName),
		QueryType(queryType),
		QueryDuration(duration),
		Int64("rows_affected", rowsAffected),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("database query executed", allFields...)
}

// LogUserAction logs a user action
func (s *StructuredLogger) LogUserAction(ctx context.Context, userID, action, resource string, fields ...Field) {
	allFields := []Field{
		UserID(userID),
		String("action", action),
		String("resource", resource),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("user action", allFields...)
}

// LogError logs an error with context
func (s *StructuredLogger) LogError(ctx context.Context, err error, errorCode, errorType string, fields ...Field) {
	allFields := []Field{
		Error(err),
		ErrorCode(errorCode),
		ErrorType(errorType),
	}
	allFields = append(allFields, fields...)

	s.logger.Error("error occurred", allFields...)
}

// LogPerformance logs performance metrics
func (s *StructuredLogger) LogPerformance(ctx context.Context, operation string, duration time.Duration, memoryUsage int64, cpuUsage float64, fields ...Field) {
	allFields := []Field{
		String("operation", operation),
		Duration("duration", duration),
		MemoryUsage(memoryUsage),
		CPUUsage(cpuUsage),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("performance metrics", allFields...)
}

// LogSecurity logs security events
func (s *StructuredLogger) LogSecurity(ctx context.Context, event, severity string, fields ...Field) {
	allFields := []Field{
		String("security_event", event),
		String("severity", severity),
	}
	allFields = append(allFields, fields...)

	s.logger.Warn("security event", allFields...)
}

// LogBusinessEvent logs business events
func (s *StructuredLogger) LogBusinessEvent(ctx context.Context, event, entityType, entityID string, fields ...Field) {
	allFields := []Field{
		String("business_event", event),
		String("entity_type", entityType),
		String("entity_id", entityID),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("business event", allFields...)
}

// LogSystemEvent logs system events
func (s *StructuredLogger) LogSystemEvent(ctx context.Context, event, component string, fields ...Field) {
	allFields := []Field{
		String("system_event", event),
		String("component", component),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("system event", allFields...)
}

// LogAudit logs audit events
func (s *StructuredLogger) LogAudit(ctx context.Context, action, resource, userID string, fields ...Field) {
	allFields := []Field{
		String("audit_action", action),
		String("audit_resource", resource),
		UserID(userID),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("audit event", allFields...)
}

// LogMetrics logs application metrics
func (s *StructuredLogger) LogMetrics(ctx context.Context, metricName string, value float64, tags map[string]string, fields ...Field) {
	allFields := []Field{
		String("metric_name", metricName),
		Float64("metric_value", value),
	}

	// Add tags as fields
	for key, value := range tags {
		allFields = append(allFields, String("tag_"+key, value))
	}

	allFields = append(allFields, fields...)

	s.logger.Info("metric recorded", allFields...)
}

// LogStartup logs service startup
func (s *StructuredLogger) LogStartup(ctx context.Context, serviceName, version, environment string, fields ...Field) {
	allFields := []Field{
		ServiceName(serviceName),
		ServiceVersion(version),
		Environment(environment),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("service started", allFields...)
}

// LogShutdown logs service shutdown
func (s *StructuredLogger) LogShutdown(ctx context.Context, serviceName string, uptime time.Duration, fields ...Field) {
	allFields := []Field{
		ServiceName(serviceName),
		Duration("uptime", uptime),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("service shutting down", allFields...)
}

// LogHealthCheck logs health check results
func (s *StructuredLogger) LogHealthCheck(ctx context.Context, component string, healthy bool, duration time.Duration, fields ...Field) {
	allFields := []Field{
		String("component", component),
		Bool("healthy", healthy),
		Duration("duration", duration),
	}
	allFields = append(allFields, fields...)

	if healthy {
		s.logger.Debug("health check passed", allFields...)
	} else {
		s.logger.Warn("health check failed", allFields...)
	}
}

// LogConfiguration logs configuration changes
func (s *StructuredLogger) LogConfiguration(ctx context.Context, configKey, oldValue, newValue string, fields ...Field) {
	allFields := []Field{
		String("config_key", configKey),
		String("old_value", oldValue),
		String("new_value", newValue),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("configuration changed", allFields...)
}

// LogCache logs cache operations
func (s *StructuredLogger) LogCache(ctx context.Context, operation, key string, hit bool, duration time.Duration, fields ...Field) {
	allFields := []Field{
		String("cache_operation", operation),
		String("cache_key", key),
		Bool("cache_hit", hit),
		Duration("duration", duration),
	}
	allFields = append(allFields, fields...)

	s.logger.Debug("cache operation", allFields...)
}

// LogExternalService logs external service calls
func (s *StructuredLogger) LogExternalService(ctx context.Context, serviceName, endpoint string, duration time.Duration, statusCode int, fields ...Field) {
	allFields := []Field{
		String("external_service", serviceName),
		String("endpoint", endpoint),
		Duration("duration", duration),
		Int("status_code", statusCode),
	}
	allFields = append(allFields, fields...)

	s.logger.Info("external service call", allFields...)
}
