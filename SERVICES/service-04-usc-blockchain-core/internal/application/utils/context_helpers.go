package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	// CorrelationIDKey is the context key for correlation ID
	CorrelationIDKey = "correlation_id"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
	// XCorrelationIDHeader is the gRPC metadata header key for correlation ID
	XCorrelationIDHeader = "x-correlation-id"
	// XRequestIDHeader is the gRPC metadata header key for request ID
	XRequestIDHeader = "x-request-id"
)

// GetCorrelationID extracts correlation ID from context
// Supports both context.Value and gRPC metadata
func GetCorrelationID(ctx context.Context) string {
	// Try context.Value first (for HTTP/standard context)
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok && id != "" {
		return id
	}

	// Try gRPC metadata (incoming)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(XCorrelationIDHeader)
		if len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}

	// Try gRPC metadata (outgoing)
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		vals := md.Get(XCorrelationIDHeader)
		if len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}

	return ""
}

// GetRequestID extracts request ID from context
// Supports both context.Value and gRPC metadata
func GetRequestID(ctx context.Context) string {
	// Try context.Value first (for HTTP/standard context)
	if id, ok := ctx.Value(RequestIDKey).(string); ok && id != "" {
		return id
	}

	// Try gRPC metadata (incoming)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(XRequestIDHeader)
		if len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}

	// Try gRPC metadata (outgoing)
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		vals := md.Get(XRequestIDHeader)
		if len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}

	return ""
}

// GetObservabilityContext extracts all observability IDs from context
// Returns a map with correlation_id and request_id if available
func GetObservabilityContext(ctx context.Context) map[string]string {
	result := make(map[string]string)

	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		result["correlation_id"] = correlationID
	}

	if requestID := GetRequestID(ctx); requestID != "" {
		result["request_id"] = requestID
	}

	return result
}
