// Package metrics provides metrics collection utilities for USC platform services.
package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics provides Prometheus metrics collection
type PrometheusMetrics struct {
	registry  prometheus.Registerer
	namespace string
	subsystem string
}

// NewPrometheusMetrics creates a new Prometheus metrics instance
func NewPrometheusMetrics(namespace, subsystem string) *PrometheusMetrics {
	return &PrometheusMetrics{
		registry:  prometheus.DefaultRegisterer,
		namespace: namespace,
		subsystem: subsystem,
	}
}

// PrometheusCounter represents a Prometheus counter metric
type PrometheusCounter struct {
	counter prometheus.Counter
}

// NewPrometheusCounter creates a new Prometheus counter metric
func (pm *PrometheusMetrics) NewPrometheusCounter(name, help string, labels []string) *PrometheusCounter {
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      name,
			Help:      help,
			Namespace: pm.namespace,
			Subsystem: pm.subsystem,
		},
		labels,
	)

	return &PrometheusCounter{
		counter: counter.WithLabelValues(),
	}
}

// Increment increments the counter by 1
func (c *PrometheusCounter) Increment() {
	c.counter.Inc()
}

// Add adds the given value to the counter
func (c *PrometheusCounter) Add(value float64) {
	c.counter.Add(value)
}

// PrometheusGauge represents a Prometheus gauge metric
type PrometheusGauge struct {
	gauge prometheus.Gauge
}

// NewPrometheusGauge creates a new Prometheus gauge metric
func (pm *PrometheusMetrics) NewPrometheusGauge(name, help string, labels []string) *PrometheusGauge {
	gauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      name,
			Help:      help,
			Namespace: pm.namespace,
			Subsystem: pm.subsystem,
		},
		labels,
	)

	return &PrometheusGauge{
		gauge: gauge.WithLabelValues(),
	}
}

// Set sets the gauge value
func (g *PrometheusGauge) Set(value float64) {
	g.gauge.Set(value)
}

// Add adds the given value to the gauge
func (g *PrometheusGauge) Add(value float64) {
	g.gauge.Add(value)
}

// Sub subtracts the given value from the gauge
func (g *PrometheusGauge) Sub(value float64) {
	g.gauge.Sub(value)
}

// Inc increments the gauge by 1
func (g *PrometheusGauge) Inc() {
	g.gauge.Inc()
}

// Dec decrements the gauge by 1
func (g *PrometheusGauge) Dec() {
	g.gauge.Dec()
}

// PrometheusHistogram represents a Prometheus histogram metric
type PrometheusHistogram struct {
	histogram *prometheus.HistogramVec
}

// NewPrometheusHistogram creates a new Prometheus histogram metric
func (pm *PrometheusMetrics) NewPrometheusHistogram(name, help string, buckets []float64, labels []string) *PrometheusHistogram {
	histogram := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      name,
			Help:      help,
			Namespace: pm.namespace,
			Subsystem: pm.subsystem,
			Buckets:   buckets,
		},
		labels,
	)

	return &PrometheusHistogram{
		histogram: histogram,
	}
}

// Observe records a value in the histogram
func (h *PrometheusHistogram) Observe(value float64) {
	h.histogram.WithLabelValues().Observe(value)
}

// PrometheusTimer represents a Prometheus timer metric
type PrometheusTimer struct {
	histogram *prometheus.HistogramVec
	start     time.Time
}

// NewPrometheusTimer creates a new Prometheus timer metric
func (pm *PrometheusMetrics) NewPrometheusTimer(name, help string, buckets []float64, labels []string) *PrometheusTimer {
	histogram := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      name,
			Help:      help,
			Namespace: pm.namespace,
			Subsystem: pm.subsystem,
			Buckets:   buckets,
		},
		labels,
	)

	return &PrometheusTimer{
		histogram: histogram,
		start:     time.Now(),
	}
}

// Stop stops the timer and records the duration
func (t *PrometheusTimer) Stop() {
	t.histogram.WithLabelValues().Observe(time.Since(t.start).Seconds())
}

// Duration records a duration value
func (t *PrometheusTimer) Duration(duration time.Duration) {
	t.histogram.WithLabelValues().Observe(duration.Seconds())
}

// PrometheusMetricsCollector represents a Prometheus metrics collector
type PrometheusMetricsCollector struct {
	prometheus *PrometheusMetrics
	custom     *CustomMetrics
}

// NewPrometheusMetricsCollector creates a new Prometheus metrics collector
func NewPrometheusMetricsCollector(namespace, subsystem string) *PrometheusMetricsCollector {
	return &PrometheusMetricsCollector{
		prometheus: NewPrometheusMetrics(namespace, subsystem),
		custom:     NewCustomMetrics(),
	}
}

// GetMetrics returns all metrics
func (pmc *PrometheusMetricsCollector) GetMetrics(ctx context.Context) map[string]interface{} {
	return pmc.custom.GetAll()
}

// Reset resets all metrics
func (pmc *PrometheusMetricsCollector) Reset(ctx context.Context) {
	pmc.custom.Clear()
}
