package metrics

import (
	"context"
	"time"
)

// Middleware provides metrics middleware functionality
type Middleware struct {
	collector *MetricsCollector
}

// NewMiddleware creates a new metrics middleware
func NewMiddleware(collector *MetricsCollector) *Middleware {
	return &Middleware{collector: collector}
}

// HTTPMetricsMiddleware creates middleware for HTTP metrics
func (m *Middleware) HTTPMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract request information from context
			_ = m.getHTTPMethod(ctx)
			_ = m.getHTTPPath(ctx)

			// Execute next handler
			next(ctx)

			// Record metrics
			_ = time.Since(start)
			_ = m.getHTTPStatusCode(ctx)
			_ = m.getHTTPRequestSize(ctx)
			_ = m.getHTTPResponseSize(ctx)

			m.collector.IncrementMetric("http_requests_total")
		}
	}
}

// GRPCMetricsMiddleware creates middleware for gRPC metrics
func (m *Middleware) GRPCMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract request information from context
			_ = m.getGRPCMethod(ctx)

			// Execute next handler
			next(ctx)

			// Record metrics
			_ = time.Since(start)
			_ = m.getGRPCStatus(ctx)

			m.collector.IncrementMetric("grpc_requests_total")
		}
	}
}

// DatabaseMetricsMiddleware creates middleware for database metrics
func (m *Middleware) DatabaseMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract database information from context
			_ = m.getDatabaseName(ctx)
			_ = m.getTableName(ctx)
			_ = m.getDatabaseOperation(ctx)

			// Execute next handler
			next(ctx)

			// Record metrics
			_ = time.Since(start)

			m.collector.IncrementMetric("database_queries_total")
		}
	}
}

// CacheMetricsMiddleware creates middleware for cache metrics
func (m *Middleware) CacheMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			start := time.Now()

			// Extract cache information from context
			_ = m.getCacheName(ctx)
			_ = m.getCacheOperation(ctx)

			// Execute next handler
			next(ctx)

			// Record metrics
			_ = time.Since(start)
			_ = m.getCacheHit(ctx)

			m.collector.IncrementMetric("cache_operations_total")
		}
	}
}

// BusinessMetricsMiddleware creates middleware for business metrics
func (m *Middleware) BusinessMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			// Execute next handler
			next(ctx)

			// Record business metrics
			eventType := m.getBusinessEventType(ctx)
			entityType := m.getBusinessEntityType(ctx)

			if eventType != "" && entityType != "" {
				m.collector.IncrementMetric("business_events_total")
			}
		}
	}
}

// UserMetricsMiddleware creates middleware for user metrics
func (m *Middleware) UserMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			// Execute next handler
			next(ctx)

			// Record user metrics
			userID := m.getUserID(ctx)
			action := m.getUserAction(ctx)
			resource := m.getUserResource(ctx)

			if userID != "" && action != "" && resource != "" {
				m.collector.IncrementMetric("user_actions_total")
			}
		}
	}
}

// ErrorMetricsMiddleware creates middleware for error metrics
func (m *Middleware) ErrorMetricsMiddleware() func(next func(context.Context)) func(context.Context) {
	return func(next func(context.Context)) func(context.Context) {
		return func(ctx context.Context) {
			// Execute next handler
			next(ctx)

			// Record error metrics if error occurred
			if err := m.getError(ctx); err != nil {
				_ = m.getErrorType(ctx)
				_ = m.getServiceName(ctx)
				_ = m.getComponentName(ctx)

				m.collector.IncrementMetric("errors_total")
			}
		}
	}
}

// Helper methods to extract information from context
// These would typically be implemented based on your specific context structure

func (m *Middleware) getHTTPMethod(ctx context.Context) string {
	if method, ok := ctx.Value("http_method").(string); ok {
		return method
	}
	return "unknown"
}

func (m *Middleware) getHTTPPath(ctx context.Context) string {
	if path, ok := ctx.Value("http_path").(string); ok {
		return path
	}
	return "unknown"
}

func (m *Middleware) getHTTPStatusCode(ctx context.Context) string {
	if code, ok := ctx.Value("http_status_code").(int); ok {
		return string(rune(code))
	}
	return "unknown"
}

func (m *Middleware) getHTTPRequestSize(ctx context.Context) int64 {
	if size, ok := ctx.Value("http_request_size").(int64); ok {
		return size
	}
	return 0
}

func (m *Middleware) getHTTPResponseSize(ctx context.Context) int64 {
	if size, ok := ctx.Value("http_response_size").(int64); ok {
		return size
	}
	return 0
}

func (m *Middleware) getGRPCMethod(ctx context.Context) string {
	if method, ok := ctx.Value("grpc_method").(string); ok {
		return method
	}
	return "unknown"
}

func (m *Middleware) getGRPCStatus(ctx context.Context) string {
	if status, ok := ctx.Value("grpc_status").(string); ok {
		return status
	}
	return "unknown"
}

func (m *Middleware) getDatabaseName(ctx context.Context) string {
	if db, ok := ctx.Value("database_name").(string); ok {
		return db
	}
	return "unknown"
}

func (m *Middleware) getTableName(ctx context.Context) string {
	if table, ok := ctx.Value("table_name").(string); ok {
		return table
	}
	return "unknown"
}

func (m *Middleware) getDatabaseOperation(ctx context.Context) string {
	if op, ok := ctx.Value("database_operation").(string); ok {
		return op
	}
	return "unknown"
}

func (m *Middleware) getCacheName(ctx context.Context) string {
	if cache, ok := ctx.Value("cache_name").(string); ok {
		return cache
	}
	return "unknown"
}

func (m *Middleware) getCacheOperation(ctx context.Context) string {
	if op, ok := ctx.Value("cache_operation").(string); ok {
		return op
	}
	return "unknown"
}

func (m *Middleware) getCacheHit(ctx context.Context) bool {
	if hit, ok := ctx.Value("cache_hit").(bool); ok {
		return hit
	}
	return false
}

func (m *Middleware) getBusinessEventType(ctx context.Context) string {
	if eventType, ok := ctx.Value("business_event_type").(string); ok {
		return eventType
	}
	return ""
}

func (m *Middleware) getBusinessEntityType(ctx context.Context) string {
	if entityType, ok := ctx.Value("business_entity_type").(string); ok {
		return entityType
	}
	return ""
}

func (m *Middleware) getUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

func (m *Middleware) getUserAction(ctx context.Context) string {
	if action, ok := ctx.Value("user_action").(string); ok {
		return action
	}
	return ""
}

func (m *Middleware) getUserResource(ctx context.Context) string {
	if resource, ok := ctx.Value("user_resource").(string); ok {
		return resource
	}
	return ""
}

func (m *Middleware) getError(ctx context.Context) error {
	if err, ok := ctx.Value("error").(error); ok {
		return err
	}
	return nil
}

func (m *Middleware) getErrorType(ctx context.Context) string {
	if errorType, ok := ctx.Value("error_type").(string); ok {
		return errorType
	}
	return "unknown"
}

func (m *Middleware) getServiceName(ctx context.Context) string {
	if service, ok := ctx.Value("service_name").(string); ok {
		return service
	}
	return "unknown"
}

func (m *Middleware) getComponentName(ctx context.Context) string {
	if component, ok := ctx.Value("component_name").(string); ok {
		return component
	}
	return "unknown"
}
