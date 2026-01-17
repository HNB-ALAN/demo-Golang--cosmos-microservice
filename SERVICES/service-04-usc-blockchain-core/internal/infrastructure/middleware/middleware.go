package middleware

import (
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/middleware"
	"google.golang.org/grpc"
)

// MiddlewareManager manages all middleware for USC Blockchain Core Service
type MiddlewareManager struct {
	logger *logging.Logger
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager(logger *logging.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		logger: logger,
	}
}

// GetUnaryInterceptors returns all unary server interceptors
func (mm *MiddlewareManager) GetUnaryInterceptors() []grpc.UnaryServerInterceptor {
	// Create recovery interceptor
	recoveryInterceptor := middleware.NewGRPCRecoveryInterceptor(middleware.DefaultRecoveryConfig())

	return []grpc.UnaryServerInterceptor{
		// Recovery middleware - handles panics
		recoveryInterceptor.UnaryServerInterceptor(),

		// Add more interceptors as needed based on your requirements
	}
}

// GetStreamInterceptors returns all stream server interceptors
func (mm *MiddlewareManager) GetStreamInterceptors() []grpc.StreamServerInterceptor {
	// Create recovery interceptor
	recoveryInterceptor := middleware.NewGRPCRecoveryInterceptor(middleware.DefaultRecoveryConfig())

	return []grpc.StreamServerInterceptor{
		// Recovery middleware - handles panics
		recoveryInterceptor.StreamServerInterceptor(),

		// Add more interceptors as needed based on your requirements
	}
}

// GetClientUnaryInterceptors returns unary client interceptors
func (mm *MiddlewareManager) GetClientUnaryInterceptors() []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{
		// Add client interceptors as needed based on your requirements
	}
}

// GetClientStreamInterceptors returns stream client interceptors
func (mm *MiddlewareManager) GetClientStreamInterceptors() []grpc.StreamClientInterceptor {
	return []grpc.StreamClientInterceptor{
		// Add client stream interceptors as needed based on your requirements
	}
}

// SetupMiddleware configures middleware for the gRPC server
func (mm *MiddlewareManager) SetupMiddleware(server *grpc.Server) {
	// Middleware is applied via interceptors during server creation
	// This method can be used for additional middleware setup if needed
	mm.logger.Info("Middleware configured for USC Blockchain Core Service")
}
