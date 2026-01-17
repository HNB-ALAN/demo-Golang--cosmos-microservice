package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/usc-platform/shared/health"
)

// HealthService provides gRPC health checking functionality
type HealthService struct {
	healthService *health.Service
}

// NewHealthService creates a new gRPC health service
func NewHealthService(healthService *health.Service) *HealthService {
	return &HealthService{
		healthService: healthService,
	}
}

// RegisterHealthService registers the health service with a gRPC server
func RegisterHealthService(server *grpc.Server, serviceName, version string) *health.Service {
	healthService := health.NewService(serviceName, version)

	// Register the gRPC health service
	grpc_health_v1.RegisterHealthServer(server, &HealthServer{
		healthService: healthService,
	})

	return healthService
}

// HealthServer implements the gRPC health service
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	healthService *health.Service
}

// Check implements the health check method
func (h *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
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
func (h *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	ctx := stream.Context()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
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
		}
	}
}

// HealthClient provides client-side health checking
type HealthClient struct {
	client grpc_health_v1.HealthClient
}

// NewHealthClient creates a new health client
func NewHealthClient(conn *grpc.ClientConn) *HealthClient {
	return &HealthClient{
		client: grpc_health_v1.NewHealthClient(conn),
	}
}

// Check performs a health check
func (c *HealthClient) Check(ctx context.Context) (*grpc_health_v1.HealthCheckResponse, error) {
	req := &grpc_health_v1.HealthCheckRequest{}
	return c.client.Check(ctx, req)
}

// Watch starts watching health status
func (c *HealthClient) Watch(ctx context.Context) (grpc_health_v1.Health_WatchClient, error) {
	req := &grpc_health_v1.HealthCheckRequest{}
	return c.client.Watch(ctx, req)
}

// IsHealthy checks if the service is healthy
func (c *HealthClient) IsHealthy(ctx context.Context) (bool, error) {
	resp, err := c.Check(ctx)
	if err != nil {
		return false, err
	}

	return resp.Status == grpc_health_v1.HealthCheckResponse_SERVING, nil
}

// WaitForHealthy waits for the service to become healthy
func (c *HealthClient) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for service to become healthy")
		case <-ticker.C:
			healthy, err := c.IsHealthy(ctx)
			if err != nil {
				continue
			}
			if healthy {
				return nil
			}
		}
	}
}

// HealthChecker provides a health checker for gRPC services
type HealthChecker struct {
	client *HealthClient
	name   string
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(conn *grpc.ClientConn, name string) *HealthChecker {
	return &HealthChecker{
		client: NewHealthClient(conn),
		name:   name,
	}
}

// Check performs a health check
func (h *HealthChecker) Check(ctx context.Context) error {
	healthy, err := h.client.IsHealthy(ctx)
	if err != nil {
		return fmt.Errorf("health check failed for %s: %w", h.name, err)
	}

	if !healthy {
		return fmt.Errorf("service %s is not healthy", h.name)
	}

	return nil
}

// Name returns the name of the checker
func (h *HealthChecker) Name() string {
	return h.name
}

// Description returns the description of the checker
func (h *HealthChecker) Description() string {
	return fmt.Sprintf("gRPC health check for %s", h.name)
}
