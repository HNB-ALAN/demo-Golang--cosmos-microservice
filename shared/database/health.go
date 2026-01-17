package database

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents the health status of a database
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// HealthReport represents a comprehensive health report
type HealthReport struct {
	Overall   HealthStatus          `json:"overall"`
	Timestamp time.Time             `json:"timestamp"`
	Databases map[string]HealthInfo `json:"databases"`
	Duration  time.Duration         `json:"duration"`
}

// HealthInfo represents health information for a specific database
type HealthInfo struct {
	Status    HealthStatus  `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// HealthChecker interface for database health checks
type HealthChecker interface {
	Check(ctx context.Context) error
}

// HealthService provides comprehensive health checking
type HealthService struct {
	manager *DatabaseManager
}

// NewHealthService creates a new health service
func NewHealthService(manager *DatabaseManager) *HealthService {
	return &HealthService{
		manager: manager,
	}
}

// CheckAll performs health checks on all databases
func (h *HealthService) CheckAll(ctx context.Context) (*HealthReport, error) {
	start := time.Now()

	report := &HealthReport{
		Overall:   HealthStatusHealthy,
		Timestamp: start,
		Databases: make(map[string]HealthInfo),
	}

	// Check each database
	for name, checker := range h.manager.healthChecks {
		info := h.checkDatabase(ctx, name, checker)
		report.Databases[name] = info

		if info.Status != HealthStatusHealthy {
			report.Overall = HealthStatusUnhealthy
		}
	}

	report.Duration = time.Since(start)
	return report, nil
}

// checkDatabase performs a health check on a specific database
func (h *HealthService) checkDatabase(ctx context.Context, _ string, checker HealthChecker) HealthInfo {
	start := time.Now()

	// Create a timeout context for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := checker.Check(checkCtx)
	duration := time.Since(start)

	info := HealthInfo{
		Duration:  duration,
		Timestamp: start,
	}

	if err != nil {
		info.Status = HealthStatusUnhealthy
		info.Message = err.Error()
	} else {
		info.Status = HealthStatusHealthy
	}

	return info
}

// CheckSpecific performs a health check on a specific database
func (h *HealthService) CheckSpecific(ctx context.Context, name string) (HealthInfo, error) {
	checker, exists := h.manager.healthChecks[name]
	if !exists {
		return HealthInfo{
			Status:   HealthStatusUnknown,
			Message:  fmt.Sprintf("database %s not found", name),
			Duration: 0,
		}, fmt.Errorf("database %s not found", name)
	}

	info := h.checkDatabase(ctx, name, checker)
	return info, nil
}

// IsHealthy returns true if all databases are healthy
func (h *HealthService) IsHealthy(ctx context.Context) bool {
	report, err := h.CheckAll(ctx)
	if err != nil {
		return false
	}

	return report.Overall == HealthStatusHealthy
}

// GetUnhealthyDatabases returns a list of unhealthy databases
func (h *HealthService) GetUnhealthyDatabases(ctx context.Context) ([]string, error) {
	report, err := h.CheckAll(ctx)
	if err != nil {
		return nil, err
	}

	var unhealthy []string
	for name, info := range report.Databases {
		if info.Status != HealthStatusHealthy {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy, nil
}

// WaitForHealthy waits for all databases to become healthy
func (h *HealthService) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for databases to become healthy")
		case <-ticker.C:
			if h.IsHealthy(ctx) {
				return nil
			}
		}
	}
}

// GetHealthSummary returns a summary of database health
func (h *HealthService) GetHealthSummary(ctx context.Context) map[string]interface{} {
	report, err := h.CheckAll(ctx)
	if err != nil {
		return map[string]interface{}{
			"overall": HealthStatusUnknown,
			"error":   err.Error(),
		}
	}

	summary := map[string]interface{}{
		"overall":   report.Overall,
		"timestamp": report.Timestamp,
		"duration":  report.Duration,
		"databases": make(map[string]interface{}),
	}

	for name, info := range report.Databases {
		summary["databases"].(map[string]interface{})[name] = map[string]interface{}{
			"status":   info.Status,
			"message":  info.Message,
			"duration": info.Duration,
		}
	}

	return summary
}
