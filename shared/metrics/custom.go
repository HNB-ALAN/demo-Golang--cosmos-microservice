package metrics

import (
	"sync"
	"time"
)

// CustomMetrics provides custom metrics collection
type CustomMetrics struct {
	mu      sync.RWMutex
	metrics map[string]interface{}
}

// NewCustomMetrics creates a new custom metrics instance
func NewCustomMetrics() *CustomMetrics {
	return &CustomMetrics{
		metrics: make(map[string]interface{}),
	}
}

// Set sets a custom metric value
func (c *CustomMetrics) Set(name string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics[name] = value
}

// Get gets a custom metric value
func (c *CustomMetrics) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.metrics[name]
	return value, exists
}

// Increment increments a numeric metric
func (c *CustomMetrics) Increment(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, exists := c.metrics[name]; exists {
		if intVal, ok := value.(int); ok {
			c.metrics[name] = intVal + 1
		} else if floatVal, ok := value.(float64); ok {
			c.metrics[name] = floatVal + 1
		}
	} else {
		c.metrics[name] = 1
	}
}

// Decrement decrements a numeric metric
func (c *CustomMetrics) Decrement(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, exists := c.metrics[name]; exists {
		if intVal, ok := value.(int); ok {
			c.metrics[name] = intVal - 1
		} else if floatVal, ok := value.(float64); ok {
			c.metrics[name] = floatVal - 1
		}
	} else {
		c.metrics[name] = -1
	}
}

// Add adds a value to a numeric metric
func (c *CustomMetrics) Add(name string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, exists := c.metrics[name]; exists {
		if floatVal, ok := existing.(float64); ok {
			c.metrics[name] = floatVal + value
		} else if intVal, ok := existing.(int); ok {
			c.metrics[name] = intVal + int(value)
		}
	} else {
		c.metrics[name] = value
	}
}

// GetAll returns all metrics
func (c *CustomMetrics) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]interface{})
	for key, value := range c.metrics {
		result[key] = value
	}
	return result
}

// Clear clears all metrics
func (c *CustomMetrics) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make(map[string]interface{})
}

// Timer represents a custom timer
type Timer struct {
	start time.Time
	name  string
}

// NewTimer creates a new timer
func NewTimer(name string) *Timer {
	return &Timer{
		start: time.Now(),
		name:  name,
	}
}

// Stop stops the timer and returns the duration
func (t *Timer) Stop() time.Duration {
	return time.Since(t.start)
}

// Duration returns the current duration
func (t *Timer) Duration() time.Duration {
	return time.Since(t.start)
}

// Counter represents a custom counter
type Counter struct {
	mu    sync.RWMutex
	value int64
	name  string
}

// NewCounter creates a new counter
func NewCounter(name string) *Counter {
	return &Counter{name: name}
}

// Increment increments the counter
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Add adds a value to the counter
func (c *Counter) Add(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += value
}

// Get returns the current value
func (c *Counter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Reset resets the counter
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// Gauge represents a custom gauge
type Gauge struct {
	mu    sync.RWMutex
	value float64
	name  string
}

// NewGauge creates a new gauge
func NewGauge(name string) *Gauge {
	return &Gauge{name: name}
}

// Set sets the gauge value
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// Get returns the current value
func (g *Gauge) Get() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Add adds a value to the gauge
func (g *Gauge) Add(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += value
}

// Sub subtracts a value from the gauge
func (g *Gauge) Sub(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value -= value
}

// Histogram represents a custom histogram
type Histogram struct {
	mu      sync.RWMutex
	buckets map[string]int64
	name    string
}

// NewHistogram creates a new histogram
func NewHistogram(name string) *Histogram {
	return &Histogram{
		buckets: make(map[string]int64),
		name:    name,
	}
}

// Observe records a value in the histogram
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Simple bucket logic - can be customized
	bucket := "0-1"
	if value > 1 && value <= 5 {
		bucket = "1-5"
	} else if value > 5 && value <= 10 {
		bucket = "5-10"
	} else if value > 10 && value <= 50 {
		bucket = "10-50"
	} else if value > 50 && value <= 100 {
		bucket = "50-100"
	} else if value > 100 {
		bucket = "100+"
	}

	h.buckets[bucket]++
}

// GetBuckets returns the histogram buckets
func (h *Histogram) GetBuckets() map[string]int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]int64)
	for bucket, count := range h.buckets {
		result[bucket] = count
	}
	return result
}

// Reset resets the histogram
func (h *Histogram) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.buckets = make(map[string]int64)
}

// MetricsRegistry manages custom metrics
type MetricsRegistry struct {
	mu      sync.RWMutex
	metrics map[string]interface{}
}

// NewMetricsRegistry creates a new metrics registry
func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{
		metrics: make(map[string]interface{}),
	}
}

// Register registers a metric
func (r *MetricsRegistry) Register(name string, metric interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[name] = metric
}

// Get gets a metric
func (r *MetricsRegistry) Get(name string) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metric, exists := r.metrics[name]
	return metric, exists
}

// Unregister unregisters a metric
func (r *MetricsRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.metrics, name)
}

// GetAll returns all metrics
func (r *MetricsRegistry) GetAll() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]interface{})
	for name, metric := range r.metrics {
		result[name] = metric
	}
	return result
}

// Clear clears all metrics
func (r *MetricsRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = make(map[string]interface{})
}

// MetricsCollector provides a collection of custom metrics
type MetricsCollector struct {
	registry *MetricsRegistry
	custom   *CustomMetrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		registry: NewMetricsRegistry(),
		custom:   NewCustomMetrics(),
	}
}

// GetRegistry returns the metrics registry
func (c *MetricsCollector) GetRegistry() *MetricsRegistry {
	return c.registry
}

// GetCustom returns the custom metrics
func (c *MetricsCollector) GetCustom() *CustomMetrics {
	return c.custom
}

// RecordMetric records a custom metric
func (c *MetricsCollector) RecordMetric(name string, value interface{}) {
	c.custom.Set(name, value)
}

// IncrementMetric increments a custom metric
func (c *MetricsCollector) IncrementMetric(name string) {
	c.custom.Increment(name)
}

// DecrementMetric decrements a custom metric
func (c *MetricsCollector) DecrementMetric(name string) {
	c.custom.Decrement(name)
}

// AddMetric adds a value to a custom metric
func (c *MetricsCollector) AddMetric(name string, value float64) {
	c.custom.Add(name, value)
}

// GetMetric gets a custom metric value
func (c *MetricsCollector) GetMetric(name string) (interface{}, bool) {
	return c.custom.Get(name)
}

// GetAllMetrics returns all custom metrics
func (c *MetricsCollector) GetAllMetrics() map[string]interface{} {
	return c.custom.GetAll()
}

// ClearMetrics clears all custom metrics
func (c *MetricsCollector) ClearMetrics() {
	c.custom.Clear()
}
