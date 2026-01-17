package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Service represents a health checking service
type Service struct {
	name      string
	version   string
	checks    map[string]HealthChecker
	mu        sync.RWMutex
	startTime time.Time
	lastCheck time.Time
}

// HealthChecker interface for health checks
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
	Description() string
}

// NewService creates a new health service
func NewService(name, version string) *Service {
	return &Service{
		name:      name,
		version:   version,
		checks:    make(map[string]HealthChecker),
		startTime: time.Now(),
	}
}

// RegisterCheck registers a health checker
func (s *Service) RegisterCheck(name string, checker HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[name] = checker
}

// UnregisterCheck unregisters a health checker
func (s *Service) UnregisterCheck(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.checks, name)
}

// GetStatus returns the overall health status
func (s *Service) GetStatus(ctx context.Context) *HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := &HealthStatus{
		Service:   s.name,
		Version:   s.version,
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Uptime:    time.Since(s.startTime),
		Checks:    make(map[string]CheckResult),
	}

	// Perform all health checks
	for name, checker := range s.checks {
		result := s.performCheck(ctx, name, checker)
		status.Checks[name] = result

		if result.Status != StatusHealthy {
			status.Status = StatusUnhealthy
		}
	}

	s.lastCheck = time.Now()
	return status
}

// performCheck performs a single health check
func (s *Service) performCheck(ctx context.Context, name string, checker HealthChecker) CheckResult {
	start := time.Now()

	// Create a timeout context for the check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := checker.Check(checkCtx)
	duration := time.Since(start)

	result := CheckResult{
		Name:        name,
		Description: checker.Description(),
		Status:      StatusHealthy,
		Duration:    duration,
		Timestamp:   start,
	}

	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
	}

	return result
}

// IsHealthy returns true if all checks are healthy
func (s *Service) IsHealthy(ctx context.Context) bool {
	status := s.GetStatus(ctx)
	return status.Status == StatusHealthy
}

// GetUnhealthyChecks returns a list of unhealthy checks
func (s *Service) GetUnhealthyChecks(ctx context.Context) []string {
	status := s.GetStatus(ctx)

	var unhealthy []string
	for name, result := range status.Checks {
		if result.Status != StatusHealthy {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}

// WaitForHealthy waits for all checks to become healthy
func (s *Service) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for service to become healthy")
		case <-ticker.C:
			if s.IsHealthy(ctx) {
				return nil
			}
		}
	}
}

// GetServiceInfo returns basic service information
func (s *Service) GetServiceInfo() ServiceInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return ServiceInfo{
		Name:      s.name,
		Version:   s.version,
		StartTime: s.startTime,
		LastCheck: s.lastCheck,
		Uptime:    time.Since(s.startTime),
		Checks:    len(s.checks),
	}
}

// GetCheckNames returns the names of all registered checks
func (s *Service) GetCheckNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.checks))
	for name := range s.checks {
		names = append(names, name)
	}

	return names
}

// GetCheckInfo returns information about a specific check
func (s *Service) GetCheckInfo(name string) (CheckInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	checker, exists := s.checks[name]
	if !exists {
		return CheckInfo{}, false
	}

	return CheckInfo{
		Name:        name,
		Description: checker.Description(),
	}, true
}
