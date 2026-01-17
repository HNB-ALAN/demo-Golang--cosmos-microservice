package health

import (
	"time"
)

// Status represents the overall health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    time.Duration          `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      Status        `json:"status"`
	Duration    time.Duration `json:"duration"`
	Timestamp   time.Time     `json:"timestamp"`
	Error       string        `json:"error,omitempty"`
}

// ServiceInfo represents basic service information
type ServiceInfo struct {
	Name      string        `json:"name"`
	Version   string        `json:"version"`
	StartTime time.Time     `json:"start_time"`
	LastCheck time.Time     `json:"last_check"`
	Uptime    time.Duration `json:"uptime"`
	Checks    int           `json:"checks"`
}

// CheckInfo represents information about a health check
type CheckInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// IsHealthy returns true if the status is healthy
func (s Status) IsHealthy() bool {
	return s == StatusHealthy
}

// IsUnhealthy returns true if the status is unhealthy
func (s Status) IsUnhealthy() bool {
	return s == StatusUnhealthy
}

// IsUnknown returns true if the status is unknown
func (s Status) IsUnknown() bool {
	return s == StatusUnknown
}

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// GetOverallStatus returns the overall status based on individual check results
func GetOverallStatus(checks map[string]CheckResult) Status {
	if len(checks) == 0 {
		return StatusUnknown
	}

	for _, result := range checks {
		if result.Status != StatusHealthy {
			return StatusUnhealthy
		}
	}

	return StatusHealthy
}

// GetUnhealthyChecks returns a list of unhealthy checks
func GetUnhealthyChecks(checks map[string]CheckResult) []string {
	var unhealthy []string

	for name, result := range checks {
		if result.Status != StatusHealthy {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}

// GetHealthyChecks returns a list of healthy checks
func GetHealthyChecks(checks map[string]CheckResult) []string {
	var healthy []string

	for name, result := range checks {
		if result.Status == StatusHealthy {
			healthy = append(healthy, name)
		}
	}

	return healthy
}

// GetTotalDuration returns the total duration of all checks
func GetTotalDuration(checks map[string]CheckResult) time.Duration {
	var total time.Duration

	for _, result := range checks {
		total += result.Duration
	}

	return total
}

// GetAverageDuration returns the average duration of all checks
func GetAverageDuration(checks map[string]CheckResult) time.Duration {
	if len(checks) == 0 {
		return 0
	}

	total := GetTotalDuration(checks)
	return total / time.Duration(len(checks))
}

// GetLongestCheck returns the check with the longest duration
func GetLongestCheck(checks map[string]CheckResult) (string, time.Duration) {
	var longestName string
	var longestDuration time.Duration

	for name, result := range checks {
		if result.Duration > longestDuration {
			longestName = name
			longestDuration = result.Duration
		}
	}

	return longestName, longestDuration
}

// GetShortestCheck returns the check with the shortest duration
func GetShortestCheck(checks map[string]CheckResult) (string, time.Duration) {
	if len(checks) == 0 {
		return "", 0
	}

	var shortestName string
	var shortestDuration time.Duration
	first := true

	for name, result := range checks {
		if first || result.Duration < shortestDuration {
			shortestName = name
			shortestDuration = result.Duration
			first = false
		}
	}

	return shortestName, shortestDuration
}
