package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Registry manages health checkers for multiple services
type Registry struct {
	services map[string]*Service
	mu       sync.RWMutex
}

// NewRegistry creates a new health registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*Service),
	}
}

// RegisterService registers a health service
func (r *Registry) RegisterService(name, version string) *Service {
	r.mu.Lock()
	defer r.mu.Unlock()

	service := NewService(name, version)
	r.services[name] = service
	return service
}

// UnregisterService unregisters a health service
func (r *Registry) UnregisterService(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, name)
}

// GetService returns a health service by name
func (r *Registry) GetService(name string) (*Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[name]
	return service, exists
}

// GetServiceNames returns the names of all registered services
func (r *Registry) GetServiceNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}

	return names
}

// GetOverallStatus returns the overall status of all services
func (r *Registry) GetOverallStatus(ctx context.Context) *OverallStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	overall := &OverallStatus{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Services:  make(map[string]*HealthStatus),
	}

	// Check each service
	for name, service := range r.services {
		status := service.GetStatus(ctx)
		overall.Services[name] = status

		if status.Status != StatusHealthy {
			overall.Status = StatusUnhealthy
		}
	}

	return overall
}

// IsHealthy returns true if all services are healthy
func (r *Registry) IsHealthy(ctx context.Context) bool {
	overall := r.GetOverallStatus(ctx)
	return overall.Status == StatusHealthy
}

// GetUnhealthyServices returns a list of unhealthy services
func (r *Registry) GetUnhealthyServices(ctx context.Context) []string {
	overall := r.GetOverallStatus(ctx)

	var unhealthy []string
	for name, status := range overall.Services {
		if status.Status != StatusHealthy {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}

// WaitForHealthy waits for all services to become healthy
func (r *Registry) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for services to become healthy")
		case <-ticker.C:
			if r.IsHealthy(ctx) {
				return nil
			}
		}
	}
}

// GetServiceInfo returns information about all services
func (r *Registry) GetServiceInfo() map[string]ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info := make(map[string]ServiceInfo)
	for name, service := range r.services {
		info[name] = service.GetServiceInfo()
	}

	return info
}

// GetCheckSummary returns a summary of all checks across all services
func (r *Registry) GetCheckSummary(ctx context.Context) *CheckSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summary := &CheckSummary{
		TotalChecks:     0,
		HealthyChecks:   0,
		UnhealthyChecks: 0,
		UnknownChecks:   0,
		Services:        make(map[string]ServiceCheckSummary),
	}

	for name, service := range r.services {
		status := service.GetStatus(ctx)
		serviceSummary := ServiceCheckSummary{
			TotalChecks:     len(status.Checks),
			HealthyChecks:   0,
			UnhealthyChecks: 0,
			UnknownChecks:   0,
		}

		for _, result := range status.Checks {
			summary.TotalChecks++
			serviceSummary.TotalChecks++

			switch result.Status {
			case StatusHealthy:
				summary.HealthyChecks++
				serviceSummary.HealthyChecks++
			case StatusUnhealthy:
				summary.UnhealthyChecks++
				serviceSummary.UnhealthyChecks++
			default:
				summary.UnknownChecks++
				serviceSummary.UnknownChecks++
			}
		}

		summary.Services[name] = serviceSummary
	}

	return summary
}

// OverallStatus represents the overall status of all services
type OverallStatus struct {
	Status    Status                   `json:"status"`
	Timestamp time.Time                `json:"timestamp"`
	Services  map[string]*HealthStatus `json:"services"`
}

// CheckSummary represents a summary of all health checks
type CheckSummary struct {
	TotalChecks     int                            `json:"total_checks"`
	HealthyChecks   int                            `json:"healthy_checks"`
	UnhealthyChecks int                            `json:"unhealthy_checks"`
	UnknownChecks   int                            `json:"unknown_checks"`
	Services        map[string]ServiceCheckSummary `json:"services"`
}

// ServiceCheckSummary represents a summary of checks for a specific service
type ServiceCheckSummary struct {
	TotalChecks     int `json:"total_checks"`
	HealthyChecks   int `json:"healthy_checks"`
	UnhealthyChecks int `json:"unhealthy_checks"`
	UnknownChecks   int `json:"unknown_checks"`
}

// GetHealthPercentage returns the percentage of healthy checks
func (s *CheckSummary) GetHealthPercentage() float64 {
	if s.TotalChecks == 0 {
		return 0
	}

	return float64(s.HealthyChecks) / float64(s.TotalChecks) * 100
}

// GetServiceHealthPercentage returns the percentage of healthy checks for a specific service
func (s *CheckSummary) GetServiceHealthPercentage(serviceName string) float64 {
	serviceSummary, exists := s.Services[serviceName]
	if !exists || serviceSummary.TotalChecks == 0 {
		return 0
	}

	return float64(serviceSummary.HealthyChecks) / float64(serviceSummary.TotalChecks) * 100
}
