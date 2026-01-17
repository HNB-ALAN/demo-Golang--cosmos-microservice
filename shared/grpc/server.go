package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/health"
	"github.com/usc-platform/shared/logging"
)

// Server represents a gRPC server with shared functionality
type Server struct {
	config        *config.Config
	logger        *logging.Logger
	grpcServer    *grpc.Server
	healthService *health.Service
	startTime     time.Time
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.Config, logger *logging.Logger) *Server {
	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(logger)),
		grpc.StreamInterceptor(StreamServerInterceptor(logger)),
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.Server.GRPC.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)

	// Create health service
	healthService := health.NewService(cfg.Service.Name, cfg.Service.Version)

	return &Server{
		config:        cfg,
		logger:        logger,
		grpcServer:    grpcServer,
		healthService: healthService,
		startTime:     time.Now(),
	}
}

// RegisterHealthService registers the health service
func (s *Server) RegisterHealthService(serviceName, version string) {
	s.healthService = health.NewService(serviceName, version)

	// Register health checkers
	s.healthService.RegisterCheck("server", &ServerHealthChecker{
		server: s,
	})

	// Register gRPC health service
	grpc_health_v1.RegisterHealthServer(s.grpcServer, &HealthServer{
		healthService: s.healthService,
	})
}

// RegisterReflection registers gRPC reflection
func (s *Server) RegisterReflection() {
	reflection.Register(s.grpcServer)
}

// RegisterService registers a gRPC service
func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.grpcServer.RegisterService(desc, impl)
}

// Start starts the gRPC server
func (s *Server) Start() error {
	address := s.config.GetServerAddress()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	s.logger.Info("Starting gRPC server",
		logging.String("address", address),
		logging.String("service", s.config.Service.Name),
		logging.String("version", s.config.Service.Version),
	)

	// Start server in a goroutine
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("Failed to serve gRPC server", logging.Error(err))
		}
	}()

	return nil
}

// Stop stops the gRPC server gracefully
func (s *Server) Stop() error {
	s.logger.Info("Stopping gRPC server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("gRPC server shutdown timeout, forcing stop")
		s.grpcServer.Stop()
		return fmt.Errorf("server shutdown timeout")
	}
}

// GetHealthService returns the health service
func (s *Server) GetHealthService() *health.Service {
	return s.healthService
}

// GetServer returns the underlying gRPC server
func (s *Server) GetServer() *grpc.Server {
	return s.grpcServer
}

// IsHealthy returns true if the server is healthy
func (s *Server) IsHealthy(ctx context.Context) bool {
	return s.healthService.IsHealthy(ctx)
}

// GetStatus returns the server status
func (s *Server) GetStatus(ctx context.Context) *health.HealthStatus {
	return s.healthService.GetStatus(ctx)
}

// ServerHealthChecker implements health checking for the server
type ServerHealthChecker struct {
	server *Server
}

// Check performs a health check on the server
func (s *ServerHealthChecker) Check(ctx context.Context) error {
	// Check if server is running
	if s.server.grpcServer == nil {
		return fmt.Errorf("gRPC server is not initialized")
	}

	// Check uptime
	uptime := time.Since(s.server.startTime)
	if uptime < 0 {
		return fmt.Errorf("server start time is invalid")
	}

	return nil
}

// Name returns the name of the checker
func (s *ServerHealthChecker) Name() string {
	return "server"
}

// Description returns the description of the checker
func (s *ServerHealthChecker) Description() string {
	return "gRPC server health check"
}

// GRPCHealthServer implements the gRPC health service
type GRPCHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	healthService *health.Service
}

// Check implements the health check method
func (h *GRPCHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	status := h.healthService.GetStatus(ctx)

	var grpcStatus grpc_health_v1.HealthCheckResponse_ServingStatus
	switch status.Status {
	case health.StatusHealthy:
		grpcStatus = grpc_health_v1.HealthCheckResponse_SERVING
	case health.StatusUnhealthy:
		grpcStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	default:
		grpcStatus = grpc_health_v1.HealthCheckResponse_UNKNOWN
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpcStatus,
	}, nil
}

// Watch implements the health watch method
func (h *GRPCHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// Simple implementation - in production, you might want to implement proper watching
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			status := h.healthService.GetStatus(ctx)

			var grpcStatus grpc_health_v1.HealthCheckResponse_ServingStatus
			switch status.Status {
			case health.StatusHealthy:
				grpcStatus = grpc_health_v1.HealthCheckResponse_SERVING
			case health.StatusUnhealthy:
				grpcStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
			default:
				grpcStatus = grpc_health_v1.HealthCheckResponse_UNKNOWN
			}

			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: grpcStatus,
			}); err != nil {
				return err
			}

			// Use non-blocking wait with context cancellation
			select {
			case <-time.After(5 * time.Second):
				// Continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
