package reflection

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	sharedgrpc "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/logging"
)

// ServiceBlockchainCoreReflectionManager manages gRPC reflection for USC Blockchain Core Service
type ServiceBlockchainCoreReflectionManager struct {
	logger   *logging.Logger
	enabled  bool
	services []string
}

// ServiceBlockchainCoreReflectionConfig represents reflection configuration for USC Blockchain Core Service
type ServiceBlockchainCoreReflectionConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Services []string `mapstructure:"services"` // Specific services to expose (empty = all)
}

// DefaultServiceBlockchainCoreReflectionConfig returns default reflection configuration
func DefaultServiceBlockchainCoreReflectionConfig() *ServiceBlockchainCoreReflectionConfig {
	return &ServiceBlockchainCoreReflectionConfig{
		Enabled:  true,       // Enable reflection by default for development
		Services: []string{}, // Empty means expose all services
	}
}

// NewServiceBlockchainCoreReflectionManager creates a new reflection manager for USC Blockchain Core Service
func NewServiceBlockchainCoreReflectionManager(logger *logging.Logger, config *ServiceBlockchainCoreReflectionConfig) *ServiceBlockchainCoreReflectionManager {
	if config == nil {
		config = DefaultServiceBlockchainCoreReflectionConfig()
	}

	return &ServiceBlockchainCoreReflectionManager{
		logger:   logger,
		enabled:  config.Enabled,
		services: config.Services,
	}
}

// RegisterReflection registers gRPC reflection with the server for USC Blockchain Core Service
func (r *ServiceBlockchainCoreReflectionManager) RegisterReflection(server *grpc.Server) error {
	if !r.enabled {
		r.logger.Debug("gRPC reflection disabled for USC Blockchain Core Service")
		return nil
	}

	// Register reflection using shared library
	sharedgrpc.RegisterReflection(server)

	r.logger.Info("gRPC reflection registered for USC Blockchain Core Service",
		logging.Bool("enabled", r.enabled),
		logging.Strings("services", r.services),
	)

	return nil
}

// RegisterReflectionWithServices registers gRPC reflection with specific services
func (r *ServiceBlockchainCoreReflectionManager) RegisterReflectionWithServices(server *grpc.Server, services ...string) error {
	if !r.enabled {
		r.logger.Debug("gRPC reflection disabled for USC Blockchain Core Service")
		return nil
	}

	// Use provided services or default services
	servicesToRegister := services
	if len(servicesToRegister) == 0 {
		servicesToRegister = r.services
	}

	// Register reflection with specific services using shared library
	sharedgrpc.RegisterReflectionWithServices(server, servicesToRegister...)

	r.logger.Info("gRPC reflection registered with specific services for USC Blockchain Core Service",
		logging.Bool("enabled", r.enabled),
		logging.Strings("services", servicesToRegister),
	)

	return nil
}

// IsEnabled returns whether reflection is enabled
func (r *ServiceBlockchainCoreReflectionManager) IsEnabled() bool {
	return r.enabled
}

// GetServices returns the list of services configured for reflection
func (r *ServiceBlockchainCoreReflectionManager) GetServices() []string {
	return r.services
}

// SetEnabled enables or disables reflection
func (r *ServiceBlockchainCoreReflectionManager) SetEnabled(enabled bool) {
	r.enabled = enabled
	r.logger.Info("gRPC reflection enabled status changed for USC Blockchain Core Service",
		logging.Bool("enabled", enabled),
	)
}

// SetServices sets the list of services for reflection
func (r *ServiceBlockchainCoreReflectionManager) SetServices(services []string) {
	r.services = services
	r.logger.Info("gRPC reflection services updated for USC Blockchain Core Service",
		logging.Strings("services", services),
	)
}

// GetReflectionInfo returns information about reflection configuration
func (r *ServiceBlockchainCoreReflectionManager) GetReflectionInfo() map[string]interface{} {
	return map[string]interface{}{
		"service":       "USC Blockchain Core Service",
		"enabled":       r.enabled,
		"services":      r.services,
		"service_count": len(r.services),
	}
}

// ServiceBlockchainCoreReflectionService provides reflection service functionality for USC Blockchain Core Service
type ServiceBlockchainCoreReflectionService struct {
	manager *ServiceBlockchainCoreReflectionManager
	logger  *logging.Logger
}

// NewServiceBlockchainCoreReflectionService creates a new reflection service
func NewServiceBlockchainCoreReflectionService(logger *logging.Logger, config *ServiceBlockchainCoreReflectionConfig) *ServiceBlockchainCoreReflectionService {
	manager := NewServiceBlockchainCoreReflectionManager(logger, config)

	return &ServiceBlockchainCoreReflectionService{
		manager: manager,
		logger:  logger,
	}
}

// Start starts the reflection service
func (s *ServiceBlockchainCoreReflectionService) Start(server *grpc.Server) error {
	return s.manager.RegisterReflection(server)
}

// StartWithServices starts the reflection service with specific services
func (s *ServiceBlockchainCoreReflectionService) StartWithServices(server *grpc.Server, services ...string) error {
	return s.manager.RegisterReflectionWithServices(server, services...)
}

// Stop stops the reflection service (no-op for reflection)
func (s *ServiceBlockchainCoreReflectionService) Stop(ctx context.Context) error {
	s.logger.Info("gRPC reflection service stopped for USC Blockchain Core Service")
	return nil
}

// HealthCheck performs health check on reflection service
func (s *ServiceBlockchainCoreReflectionService) HealthCheck(ctx context.Context) error {
	if !s.manager.IsEnabled() {
		return fmt.Errorf("reflection service is disabled for USC Blockchain Core Service")
	}

	s.logger.Debug("gRPC reflection service health check passed for USC Blockchain Core Service")
	return nil
}

// GetMetrics returns reflection service metrics
func (s *ServiceBlockchainCoreReflectionService) GetMetrics() map[string]interface{} {
	metrics := s.manager.GetReflectionInfo()
	metrics["status"] = "running"
	return metrics
}

// ServiceBlockchainCoreReflectionController provides control interface for reflection
type ServiceBlockchainCoreReflectionController struct {
	manager *ServiceBlockchainCoreReflectionManager
	logger  *logging.Logger
}

// NewServiceBlockchainCoreReflectionController creates a new reflection controller
func NewServiceBlockchainCoreReflectionController(logger *logging.Logger, config *ServiceBlockchainCoreReflectionConfig) *ServiceBlockchainCoreReflectionController {
	manager := NewServiceBlockchainCoreReflectionManager(logger, config)

	return &ServiceBlockchainCoreReflectionController{
		manager: manager,
		logger:  logger,
	}
}

// EnableReflection enables gRPC reflection
func (c *ServiceBlockchainCoreReflectionController) EnableReflection() {
	c.manager.SetEnabled(true)
	c.logger.Info("gRPC reflection enabled for USC Blockchain Core Service")
}

// DisableReflection disables gRPC reflection
func (c *ServiceBlockchainCoreReflectionController) DisableReflection() {
	c.manager.SetEnabled(false)
	c.logger.Info("gRPC reflection disabled for USC Blockchain Core Service")
}

// UpdateServices updates the list of services for reflection
func (c *ServiceBlockchainCoreReflectionController) UpdateServices(services []string) {
	c.manager.SetServices(services)
	c.logger.Info("gRPC reflection services updated for USC Blockchain Core Service",
		logging.Strings("services", services),
	)
}

// GetStatus returns the current reflection status
func (c *ServiceBlockchainCoreReflectionController) GetStatus() map[string]interface{} {
	return c.manager.GetReflectionInfo()
}
