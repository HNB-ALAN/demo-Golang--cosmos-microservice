package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "performance"

	// StoreKey defines the primary store key for the performance module
	StoreKey = ModuleName

	// RouterKey defines the message route for the performance module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the performance module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypePerformanceMetric = "performance_metric"
	EventTypeBenchmarkResult   = "benchmark_result"
	EventTypeOptimization      = "optimization"
	EventTypePerformanceAlert  = "performance_alert"
)

// Event attribute keys
const (
	AttributeKeyMetricID       = "metric_id"
	AttributeKeyMetricName     = "metric_name"
	AttributeKeyValue          = "value"
	AttributeKeyUnit           = "unit"
	AttributeKeyTimestamp      = "timestamp"
	AttributeKeyBenchmarkID    = "benchmark_id"
	AttributeKeyOptimizationID = "optimization_id"
	AttributeKeyAlertID        = "alert_id"
	AttributeKeySeverity       = "severity"
)

// PerformanceMetric represents a performance metric
type PerformanceMetric struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Value       int64             `json:"value"`
	Unit        string            `json:"unit"`
	Timestamp   time.Time         `json:"timestamp"`
	Tags        map[string]string `json:"tags"`
	Description string            `json:"description"`
	Category    string            `json:"category"` // cpu, memory, network, disk, etc.
}

// Benchmark represents a performance benchmark
type Benchmark struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Duration    time.Duration     `json:"duration"`
	Results     map[string]int64  `json:"results"`
	Tags        map[string]string `json:"tags"`
	Status      string            `json:"status"` // running, completed, failed
}

// Optimization represents a performance optimization
type Optimization struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`   // algorithm, configuration, resource
	Impact      string            `json:"impact"` // high, medium, low
	Status      string            `json:"status"` // pending, applied, reverted
	CreatedAt   time.Time         `json:"created_at"`
	AppliedAt   time.Time         `json:"applied_at"`
	RevertedAt  time.Time         `json:"reverted_at"`
	Metrics     map[string]int64  `json:"metrics"`
	Tags        map[string]string `json:"tags"`
}

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // critical, high, medium, low
	Status      string    `json:"status"`   // active, resolved, acknowledged
	MetricID    string    `json:"metric_id"`
	Threshold   int64     `json:"threshold"`
	Condition   string    `json:"condition"` // gt, lt, eq, gte, lte
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  time.Time `json:"resolved_at"`
}

// PerformanceProfile represents a performance profile
type PerformanceProfile struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceName string            `json:"service_name"`
	Metrics     map[string]int64  `json:"metrics"`
	Baselines   map[string]int64  `json:"baselines"`
	Thresholds  map[string]int64  `json:"thresholds"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Tags        map[string]string `json:"tags"`
}

// PerformanceReport represents a performance report
type PerformanceReport struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Duration    time.Duration          `json:"duration"`
	Summary     map[string]int64       `json:"summary"`
	Details     map[string]interface{} `json:"details"`
	CreatedAt   time.Time              `json:"created_at"`
	Tags        map[string]string      `json:"tags"`
}

// GenesisState represents the genesis state of the performance module
type GenesisState struct {
	Metrics       []PerformanceMetric  `json:"metrics"`
	Benchmarks    []Benchmark          `json:"benchmarks"`
	Optimizations []Optimization       `json:"optimizations"`
	Alerts        []PerformanceAlert   `json:"alerts"`
	Profiles      []PerformanceProfile `json:"profiles"`
	Reports       []PerformanceReport  `json:"reports"`
	Params        Params               `json:"params"`
}

// Params represents the parameters for the performance module
type Params struct {
	MaxMetricsPerService int64         `json:"max_metrics_per_service"`
	DefaultRetention     time.Duration `json:"default_retention"`
	AlertCooldown        time.Duration `json:"alert_cooldown"`
	BenchmarkTimeout     time.Duration `json:"benchmark_timeout"`
	OptimizationDelay    time.Duration `json:"optimization_delay"`
}

// DefaultParams returns the default parameters for the performance module
func DefaultParams() Params {
	return Params{
		MaxMetricsPerService: 1000,
		DefaultRetention:     24 * time.Hour,
		AlertCooldown:        5 * time.Minute,
		BenchmarkTimeout:     10 * time.Minute,
		OptimizationDelay:    1 * time.Minute,
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
	if p.BenchmarkTimeout <= 0 {
		return fmt.Errorf("benchmark timeout must be positive")
	}
	if p.OptimizationDelay <= 0 {
		return fmt.Errorf("optimization delay must be positive")
	}
	return nil
}

// Validate validates a performance metric
func (m PerformanceMetric) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("metric ID cannot be empty")
	}
	if m.Name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}
	if m.Value < 0 {
		return fmt.Errorf("metric value cannot be negative")
	}
	if m.Category == "" {
		return fmt.Errorf("metric category cannot be empty")
	}
	return nil
}

// Validate validates a benchmark
func (b Benchmark) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("benchmark ID cannot be empty")
	}
	if b.Name == "" {
		return fmt.Errorf("benchmark name cannot be empty")
	}
	if b.Status != "running" && b.Status != "completed" && b.Status != "failed" {
		return fmt.Errorf("invalid status: %s", b.Status)
	}
	return nil
}

// Validate validates an optimization
func (o Optimization) Validate() error {
	if o.ID == "" {
		return fmt.Errorf("optimization ID cannot be empty")
	}
	if o.Name == "" {
		return fmt.Errorf("optimization name cannot be empty")
	}
	if o.Type != "algorithm" && o.Type != "configuration" && o.Type != "resource" {
		return fmt.Errorf("invalid type: %s", o.Type)
	}
	if o.Impact != "high" && o.Impact != "medium" && o.Impact != "low" {
		return fmt.Errorf("invalid impact: %s", o.Impact)
	}
	if o.Status != "pending" && o.Status != "applied" && o.Status != "reverted" {
		return fmt.Errorf("invalid status: %s", o.Status)
	}
	return nil
}

// Validate validates a performance alert
func (a PerformanceAlert) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("alert ID cannot be empty")
	}
	if a.Name == "" {
		return fmt.Errorf("alert name cannot be empty")
	}
	if a.Severity != "critical" && a.Severity != "high" && a.Severity != "medium" && a.Severity != "low" {
		return fmt.Errorf("invalid severity: %s", a.Severity)
	}
	if a.Status != "active" && a.Status != "resolved" && a.Status != "acknowledged" {
		return fmt.Errorf("invalid status: %s", a.Status)
	}
	return nil
}

// Validate validates a performance profile
func (p PerformanceProfile) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("profile ID cannot be empty")
	}
	if p.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if p.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	return nil
}

// Validate validates a performance report
func (r PerformanceReport) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("report ID cannot be empty")
	}
	if r.Name == "" {
		return fmt.Errorf("report name cannot be empty")
	}
	return nil
}
