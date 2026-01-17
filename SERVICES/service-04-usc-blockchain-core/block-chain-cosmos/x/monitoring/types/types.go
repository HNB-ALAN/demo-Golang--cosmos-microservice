package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "monitoring"

	// StoreKey defines the primary store key for the monitoring module
	StoreKey = ModuleName

	// RouterKey defines the message route for the monitoring module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the monitoring module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeMetricCreated  = "metric_created"
	EventTypeMetricUpdated  = "metric_updated"
	EventTypeAlertTriggered = "alert_triggered"
	EventTypeAlertResolved  = "alert_resolved"
)

// Event attribute keys
const (
	AttributeKeyMetricID   = "metric_id"
	AttributeKeyMetricName = "metric_name"
	AttributeKeyValue      = "value"
	AttributeKeyTimestamp  = "timestamp"
	AttributeKeyAlertID    = "alert_id"
	AttributeKeyAlertType  = "alert_type"
	AttributeKeySeverity   = "severity"
)

// Metric represents a performance metric
type Metric struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Value       int64             `json:"value"`
	Unit        string            `json:"unit"`
	Timestamp   time.Time         `json:"timestamp"`
	Tags        map[string]string `json:"tags"`
	Description string            `json:"description"`
}

// Alert represents a monitoring alert
type Alert struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // low, medium, high, critical
	Status      string    `json:"status"`   // active, resolved, acknowledged
	MetricID    string    `json:"metric_id"`
	Threshold   int64     `json:"threshold"`
	Condition   string    `json:"condition"` // gt, lt, eq, gte, lte
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PerformanceData represents performance monitoring data
type PerformanceData struct {
	ID          string            `json:"id"`
	ServiceName string            `json:"service_name"`
	MetricName  string            `json:"metric_name"`
	Value       int64             `json:"value"`
	Unit        string            `json:"unit"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata"`
}

// SystemHealth represents overall system health status
type SystemHealth struct {
	ID         string            `json:"id"`
	Status     string            `json:"status"` // healthy, warning, critical, down
	Score      int64             `json:"score"`  // 0-100
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentHealth `json:"components"`
	Summary    string            `json:"summary"`
}

// ComponentHealth represents health of individual components
type ComponentHealth struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Score     int64     `json:"score"`
	LastCheck time.Time `json:"last_check"`
	Message   string    `json:"message"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	ID              string           `json:"id"`
	ServiceName     string           `json:"service_name"`
	Enabled         bool             `json:"enabled"`
	CheckInterval   time.Duration    `json:"check_interval"`
	AlertThresholds map[string]int64 `json:"alert_thresholds"`
	RetentionPeriod time.Duration    `json:"retention_period"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// GenesisState represents the genesis state of the monitoring module
type GenesisState struct {
	Metrics          []Metric           `json:"metrics"`
	Alerts           []Alert            `json:"alerts"`
	PerformanceData  []PerformanceData  `json:"performance_data"`
	SystemHealth     []SystemHealth     `json:"system_health"`
	MonitoringConfig []MonitoringConfig `json:"monitoring_config"`
	Params           Params             `json:"params"`
}

// Params represents the parameters for the monitoring module
type Params struct {
	MaxMetricsPerService int64         `json:"max_metrics_per_service"`
	DefaultRetention     time.Duration `json:"default_retention"`
	AlertCooldown        time.Duration `json:"alert_cooldown"`
	HealthCheckInterval  time.Duration `json:"health_check_interval"`
}

// DefaultParams returns the default parameters for the monitoring module
func DefaultParams() Params {
	return Params{
		MaxMetricsPerService: 1000,
		DefaultRetention:     24 * time.Hour,
		AlertCooldown:        5 * time.Minute,
		HealthCheckInterval:  30 * time.Second,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxMetricsPerService <= 0 {
		return fmt.Errorf("max metrics per service must be positive")
	}
	if p.DefaultRetention <= 0 {
		return fmt.Errorf("default retention must be positive")
	}
	if p.AlertCooldown <= 0 {
		return fmt.Errorf("alert cooldown must be positive")
	}
	if p.HealthCheckInterval <= 0 {
		return fmt.Errorf("health check interval must be positive")
	}
	return nil
}

// Validate validates a metric
func (m Metric) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("metric ID cannot be empty")
	}
	if m.Name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}
	if m.Value < 0 {
		return fmt.Errorf("metric value cannot be negative")
	}
	return nil
}

// Validate validates an alert
func (a Alert) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("alert ID cannot be empty")
	}
	if a.Name == "" {
		return fmt.Errorf("alert name cannot be empty")
	}
	if a.Severity != "low" && a.Severity != "medium" && a.Severity != "high" && a.Severity != "critical" {
		return fmt.Errorf("invalid severity: %s", a.Severity)
	}
	if a.Status != "active" && a.Status != "resolved" && a.Status != "acknowledged" {
		return fmt.Errorf("invalid status: %s", a.Status)
	}
	return nil
}

// Validate validates performance data
func (p PerformanceData) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("performance data ID cannot be empty")
	}
	if p.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if p.MetricName == "" {
		return fmt.Errorf("metric name cannot be empty")
	}
	return nil
}

// Validate validates system health
func (s SystemHealth) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("system health ID cannot be empty")
	}
	if s.Status != "healthy" && s.Status != "warning" && s.Status != "critical" && s.Status != "down" {
		return fmt.Errorf("invalid status: %s", s.Status)
	}
	if s.Score < 0 || s.Score > 100 {
		return fmt.Errorf("score must be between 0 and 100")
	}
	return nil
}

// Validate validates monitoring config
func (c MonitoringConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("config ID cannot be empty")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if c.CheckInterval <= 0 {
		return fmt.Errorf("check interval must be positive")
	}
	return nil
}
