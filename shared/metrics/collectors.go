package metrics

import (
	"fmt"
	"runtime"
	"time"
)

// SystemMetricsCollector collects system-level metrics
type SystemMetricsCollector struct {
	collector *MetricsCollector
	ticker    *time.Ticker
	done      chan bool
}

// NewSystemMetricsCollector creates a new system metrics collector
func NewSystemMetricsCollector(collector *MetricsCollector, interval time.Duration) *SystemMetricsCollector {
	return &SystemMetricsCollector{
		collector: collector,
		ticker:    time.NewTicker(interval),
		done:      make(chan bool),
	}
}

// Start starts collecting system metrics
func (s *SystemMetricsCollector) Start() {
	go s.collect()
}

// Stop stops collecting system metrics
func (s *SystemMetricsCollector) Stop() {
	s.done <- true
	s.ticker.Stop()
}

// collect collects system metrics periodically
func (s *SystemMetricsCollector) collect() {
	for {
		select {
		case <-s.ticker.C:
			s.collectSystemMetrics()
		case <-s.done:
			return
		}
	}
}

// collectSystemMetrics collects various system metrics
func (s *SystemMetricsCollector) collectSystemMetrics() {
	// Memory metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s.collector.RecordMetric("memory_alloc_bytes", m.Alloc)
	s.collector.RecordMetric("memory_total_alloc_bytes", m.TotalAlloc)
	s.collector.RecordMetric("memory_sys_bytes", m.Sys)
	s.collector.RecordMetric("memory_num_gc", m.NumGC)
	s.collector.RecordMetric("memory_gc_cpu_fraction", m.GCCPUFraction)

	// Goroutine metrics
	s.collector.RecordMetric("goroutines_count", runtime.NumGoroutine())

	// GC metrics
	s.collector.RecordMetric("gc_pause_total_ns", m.PauseTotalNs)
	s.collector.RecordMetric("gc_pause_ns", m.PauseNs[(m.NumGC+255)%256])

	// Heap metrics
	s.collector.RecordMetric("heap_alloc_bytes", m.HeapAlloc)
	s.collector.RecordMetric("heap_sys_bytes", m.HeapSys)
	s.collector.RecordMetric("heap_idle_bytes", m.HeapIdle)
	s.collector.RecordMetric("heap_inuse_bytes", m.HeapInuse)
	s.collector.RecordMetric("heap_released_bytes", m.HeapReleased)
	s.collector.RecordMetric("heap_objects", m.HeapObjects)

	// Stack metrics
	s.collector.RecordMetric("stack_inuse_bytes", m.StackInuse)
	s.collector.RecordMetric("stack_sys_bytes", m.StackSys)

	// MSpan metrics
	s.collector.RecordMetric("mspan_inuse_bytes", m.MSpanInuse)
	s.collector.RecordMetric("mspan_sys_bytes", m.MSpanSys)

	// MCache metrics
	s.collector.RecordMetric("mcache_inuse_bytes", m.MCacheInuse)
	s.collector.RecordMetric("mcache_sys_bytes", m.MCacheSys)

	// BuckHashSys metrics
	s.collector.RecordMetric("buck_hash_sys_bytes", m.BuckHashSys)

	// GCSys metrics
	s.collector.RecordMetric("gc_sys_bytes", m.GCSys)

	// OtherSys metrics
	s.collector.RecordMetric("other_sys_bytes", m.OtherSys)

	// NextGC metrics
	s.collector.RecordMetric("next_gc_bytes", m.NextGC)

	// LastGC metrics
	s.collector.RecordMetric("last_gc_timestamp", m.LastGC)

	// PauseNs metrics (last 256 GC pauses)
	for i, pause := range m.PauseNs {
		if pause > 0 {
			s.collector.RecordMetric(fmt.Sprintf("gc_pause_%d_ns", i), pause)
		}
	}
}

// ApplicationMetricsCollector collects application-level metrics
type ApplicationMetricsCollector struct {
	collector *MetricsCollector
}

// NewApplicationMetricsCollector creates a new application metrics collector
func NewApplicationMetricsCollector(collector *MetricsCollector) *ApplicationMetricsCollector {
	return &ApplicationMetricsCollector{collector: collector}
}

// RecordRequestDuration records request duration
func (a *ApplicationMetricsCollector) RecordRequestDuration(service, method string, duration time.Duration) {
	a.collector.RecordMetric(fmt.Sprintf("request_duration_%s_%s_ms", service, method), duration.Milliseconds())
}

// RecordRequestCount records request count
func (a *ApplicationMetricsCollector) RecordRequestCount(service, method string) {
	a.collector.IncrementMetric(fmt.Sprintf("request_count_%s_%s", service, method))
}

// RecordErrorCount records error count
func (a *ApplicationMetricsCollector) RecordErrorCount(service, method, errorType string) {
	a.collector.IncrementMetric(fmt.Sprintf("error_count_%s_%s_%s", service, method, errorType))
}

// RecordActiveConnections records active connections
func (a *ApplicationMetricsCollector) RecordActiveConnections(service string, count int) {
	a.collector.RecordMetric(fmt.Sprintf("active_connections_%s", service), count)
}

// RecordQueueSize records queue size
func (a *ApplicationMetricsCollector) RecordQueueSize(queue string, size int) {
	a.collector.RecordMetric(fmt.Sprintf("queue_size_%s", queue), size)
}

// RecordCacheHitRate records cache hit rate
func (a *ApplicationMetricsCollector) RecordCacheHitRate(cache string, hitRate float64) {
	a.collector.RecordMetric(fmt.Sprintf("cache_hit_rate_%s", cache), hitRate)
}

// RecordDatabaseConnectionPool records database connection pool metrics
func (a *ApplicationMetricsCollector) RecordDatabaseConnectionPool(database string, active, idle, max int) {
	a.collector.RecordMetric(fmt.Sprintf("db_pool_active_%s", database), active)
	a.collector.RecordMetric(fmt.Sprintf("db_pool_idle_%s", database), idle)
	a.collector.RecordMetric(fmt.Sprintf("db_pool_max_%s", database), max)
}

// RecordBusinessMetric records business-specific metrics
func (a *ApplicationMetricsCollector) RecordBusinessMetric(metric string, value interface{}) {
	a.collector.RecordMetric(fmt.Sprintf("business_%s", metric), value)
}

// RecordUserMetric records user-specific metrics
func (a *ApplicationMetricsCollector) RecordUserMetric(userID, metric string, value interface{}) {
	a.collector.RecordMetric(fmt.Sprintf("user_%s_%s", userID, metric), value)
}

// RecordServiceMetric records service-specific metrics
func (a *ApplicationMetricsCollector) RecordServiceMetric(service, metric string, value interface{}) {
	a.collector.RecordMetric(fmt.Sprintf("service_%s_%s", service, metric), value)
}

// DatabaseMetricsCollector collects database-specific metrics
type DatabaseMetricsCollector struct {
	collector *MetricsCollector
}

// NewDatabaseMetricsCollector creates a new database metrics collector
func NewDatabaseMetricsCollector(collector *MetricsCollector) *DatabaseMetricsCollector {
	return &DatabaseMetricsCollector{collector: collector}
}

// RecordQueryDuration records query duration
func (d *DatabaseMetricsCollector) RecordQueryDuration(database, table, operation string, duration time.Duration) {
	d.collector.RecordMetric(fmt.Sprintf("db_query_duration_%s_%s_%s_ms", database, table, operation), duration.Milliseconds())
}

// RecordQueryCount records query count
func (d *DatabaseMetricsCollector) RecordQueryCount(database, table, operation string) {
	d.collector.IncrementMetric(fmt.Sprintf("db_query_count_%s_%s_%s", database, table, operation))
}

// RecordConnectionCount records connection count
func (d *DatabaseMetricsCollector) RecordConnectionCount(database string, count int) {
	d.collector.RecordMetric(fmt.Sprintf("db_connections_%s", database), count)
}

// RecordTransactionCount records transaction count
func (d *DatabaseMetricsCollector) RecordTransactionCount(database string, count int) {
	d.collector.RecordMetric(fmt.Sprintf("db_transactions_%s", database), count)
}

// RecordLockWaitTime records lock wait time
func (d *DatabaseMetricsCollector) RecordLockWaitTime(database string, waitTime time.Duration) {
	d.collector.RecordMetric(fmt.Sprintf("db_lock_wait_%s_ms", database), waitTime.Milliseconds())
}

// RecordDeadlockCount records deadlock count
func (d *DatabaseMetricsCollector) RecordDeadlockCount(database string) {
	d.collector.IncrementMetric(fmt.Sprintf("db_deadlocks_%s", database))
}

// RecordSlowQueryCount records slow query count
func (d *DatabaseMetricsCollector) RecordSlowQueryCount(database string) {
	d.collector.IncrementMetric(fmt.Sprintf("db_slow_queries_%s", database))
}

// CacheMetricsCollector collects cache-specific metrics
type CacheMetricsCollector struct {
	collector *MetricsCollector
}

// NewCacheMetricsCollector creates a new cache metrics collector
func NewCacheMetricsCollector(collector *MetricsCollector) *CacheMetricsCollector {
	return &CacheMetricsCollector{collector: collector}
}

// RecordCacheOperation records cache operation
func (c *CacheMetricsCollector) RecordCacheOperation(cache, operation string, duration time.Duration) {
	c.collector.RecordMetric(fmt.Sprintf("cache_operation_%s_%s_ms", cache, operation), duration.Milliseconds())
}

// RecordCacheHit records cache hit
func (c *CacheMetricsCollector) RecordCacheHit(cache string) {
	c.collector.IncrementMetric(fmt.Sprintf("cache_hits_%s", cache))
}

// RecordCacheMiss records cache miss
func (c *CacheMetricsCollector) RecordCacheMiss(cache string) {
	c.collector.IncrementMetric(fmt.Sprintf("cache_misses_%s", cache))
}

// RecordCacheSize records cache size
func (c *CacheMetricsCollector) RecordCacheSize(cache string, size int) {
	c.collector.RecordMetric(fmt.Sprintf("cache_size_%s", cache), size)
}

// RecordCacheEviction records cache eviction
func (c *CacheMetricsCollector) RecordCacheEviction(cache string) {
	c.collector.IncrementMetric(fmt.Sprintf("cache_evictions_%s", cache))
}

// RecordCacheExpiration records cache expiration
func (c *CacheMetricsCollector) RecordCacheExpiration(cache string) {
	c.collector.IncrementMetric(fmt.Sprintf("cache_expirations_%s", cache))
}

// RecordCacheMemoryUsage records cache memory usage
func (c *CacheMetricsCollector) RecordCacheMemoryUsage(cache string, usage int64) {
	c.collector.RecordMetric(fmt.Sprintf("cache_memory_%s_bytes", cache), usage)
}
