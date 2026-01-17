package database

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PoolManager manages database connection pools with dynamic sizing
type PoolManager struct {
	mu sync.RWMutex

	// Pool configuration
	minSize     int
	maxSize     int
	currentSize int

	// Performance tracking
	activeConnections int
	idleConnections   int
	waitingRequests   int

	// Load balancing
	loadFactor         float64
	lastAdjustment     time.Time
	adjustmentInterval time.Duration

	// Metrics
	totalRequests      int64
	successfulRequests int64
	failedRequests     int64

	// Health monitoring
	healthChecks   map[string]time.Time
	unhealthyPools map[string]bool
}

// PoolConfig represents pool configuration
type PoolConfig struct {
	MinSize             int           `mapstructure:"min_size"`
	MaxSize             int           `mapstructure:"max_size"`
	InitialSize         int           `mapstructure:"initial_size"`
	AdjustmentInterval  time.Duration `mapstructure:"adjustment_interval"`
	LoadThreshold       float64       `mapstructure:"load_threshold"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
}

// NewPoolManager creates a new pool manager
func NewPoolManager(config PoolConfig) *PoolManager {
	return &PoolManager{
		minSize:            config.MinSize,
		maxSize:            config.MaxSize,
		currentSize:        config.InitialSize,
		loadFactor:         0.0,
		lastAdjustment:     time.Now(),
		adjustmentInterval: config.AdjustmentInterval,
		healthChecks:       make(map[string]time.Time),
		unhealthyPools:     make(map[string]bool),
	}
}

// GetConnection gets a connection from the pool
func (p *PoolManager) GetConnection(ctx context.Context) (*PooledConnection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.totalRequests++
	p.waitingRequests++

	// Check if we need to adjust pool size
	if time.Since(p.lastAdjustment) > p.adjustmentInterval {
		p.adjustPoolSize()
		p.lastAdjustment = time.Now()
	}

	// Try to get an idle connection
	if p.idleConnections > 0 {
		p.idleConnections--
		p.activeConnections++
		p.waitingRequests--
		p.successfulRequests++

		return &PooledConnection{
			pool:    p,
			created: time.Now(),
			active:  true,
		}, nil
	}

	// Create new connection if under limit
	if p.activeConnections < p.currentSize {
		p.activeConnections++
		p.waitingRequests--
		p.successfulRequests++

		return &PooledConnection{
			pool:    p,
			created: time.Now(),
			active:  true,
		}, nil
	}

	// Wait for connection to become available
	p.waitingRequests--
	p.failedRequests++

	return nil, fmt.Errorf("no connections available, pool size: %d, active: %d", p.currentSize, p.activeConnections)
}

// ReturnConnection returns a connection to the pool
func (p *PoolManager) ReturnConnection(conn *PooledConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !conn.active {
		return
	}

	conn.active = false
	p.activeConnections--

	// Return to idle pool if we have space
	if p.idleConnections < p.currentSize/2 {
		p.idleConnections++
	} else {
		// Close excess connections
		p.currentSize--
	}
}

// adjustPoolSize adjusts the pool size based on load
func (p *PoolManager) adjustPoolSize() {
	// Calculate load factor
	if p.currentSize > 0 {
		p.loadFactor = float64(p.activeConnections+p.waitingRequests) / float64(p.currentSize)
	}

	// Increase pool size if load is high
	if p.loadFactor > 0.8 && p.currentSize < p.maxSize {
		newSize := p.currentSize * 2
		if newSize > p.maxSize {
			newSize = p.maxSize
		}
		p.currentSize = newSize
	}

	// Decrease pool size if load is low
	if p.loadFactor < 0.3 && p.currentSize > p.minSize {
		newSize := p.currentSize / 2
		if newSize < p.minSize {
			newSize = p.minSize
		}
		p.currentSize = newSize
	}
}

// GetStats returns pool statistics
func (p *PoolManager) GetStats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return PoolStats{
		MinSize:            p.minSize,
		MaxSize:            p.maxSize,
		CurrentSize:        p.currentSize,
		ActiveConnections:  p.activeConnections,
		IdleConnections:    p.idleConnections,
		WaitingRequests:    p.waitingRequests,
		LoadFactor:         p.loadFactor,
		TotalRequests:      p.totalRequests,
		SuccessfulRequests: p.successfulRequests,
		FailedRequests:     p.failedRequests,
		SuccessRate:        p.calculateSuccessRate(),
	}
}

// calculateSuccessRate calculates the success rate
func (p *PoolManager) calculateSuccessRate() float64 {
	if p.totalRequests == 0 {
		return 0
	}
	return float64(p.successfulRequests) / float64(p.totalRequests) * 100
}

// HealthCheck performs health check on the pool
func (p *PoolManager) HealthCheck(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if pool is healthy
	if p.loadFactor > 0.95 {
		return fmt.Errorf("pool overloaded: load factor %.2f", p.loadFactor)
	}

	if p.failedRequests > 0 && p.calculateSuccessRate() < 90 {
		return fmt.Errorf("pool unhealthy: success rate %.2f%%", p.calculateSuccessRate())
	}

	return nil
}

// PoolStats represents pool statistics
type PoolStats struct {
	MinSize            int     `json:"min_size"`
	MaxSize            int     `json:"max_size"`
	CurrentSize        int     `json:"current_size"`
	ActiveConnections  int     `json:"active_connections"`
	IdleConnections    int     `json:"idle_connections"`
	WaitingRequests    int     `json:"waiting_requests"`
	LoadFactor         float64 `json:"load_factor"`
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	SuccessRate        float64 `json:"success_rate"`
}

// PooledConnection represents a pooled database connection
type PooledConnection struct {
	pool    *PoolManager
	created time.Time
	active  bool
}

// Close closes the connection and returns it to the pool
func (c *PooledConnection) Close() error {
	if c.active {
		c.pool.ReturnConnection(c)
	}
	return nil
}

// IsActive returns true if the connection is active
func (c *PooledConnection) IsActive() bool {
	return c.active
}

// GetAge returns the age of the connection
func (c *PooledConnection) GetAge() time.Duration {
	return time.Since(c.created)
}

// MultiPoolManager manages multiple database pools
type MultiPoolManager struct {
	mu     sync.RWMutex
	pools  map[string]*PoolManager
	config PoolConfig
}

// NewMultiPoolManager creates a new multi-pool manager
func NewMultiPoolManager(config PoolConfig) *MultiPoolManager {
	return &MultiPoolManager{
		pools:  make(map[string]*PoolManager),
		config: config,
	}
}

// GetPool gets or creates a pool for the given database
func (m *MultiPoolManager) GetPool(database string) *PoolManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	if pool, exists := m.pools[database]; exists {
		return pool
	}

	pool := NewPoolManager(m.config)
	m.pools[database] = pool
	return pool
}

// GetConnection gets a connection from the specified database pool
func (m *MultiPoolManager) GetConnection(ctx context.Context, database string) (*PooledConnection, error) {
	pool := m.GetPool(database)
	return pool.GetConnection(ctx)
}

// GetStats returns statistics for all pools
func (m *MultiPoolManager) GetStats() map[string]PoolStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]PoolStats)
	for name, pool := range m.pools {
		stats[name] = pool.GetStats()
	}
	return stats
}

// HealthCheck performs health check on all pools
func (m *MultiPoolManager) HealthCheck(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, pool := range m.pools {
		if err := pool.HealthCheck(ctx); err != nil {
			return fmt.Errorf("pool %s health check failed: %w", name, err)
		}
	}
	return nil
}

// Close closes all pools
func (m *MultiPoolManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close all pools
	for name, pool := range m.pools {
		// In a real implementation, you'd close actual connections
		_ = pool
		_ = name
	}

	m.pools = make(map[string]*PoolManager)
	return nil
}
