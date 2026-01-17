package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Checker represents a basic health checker
type Checker struct {
	name        string
	description string
	checkFunc   func(ctx context.Context) error
}

// NewChecker creates a new health checker
func NewChecker(name, description string, checkFunc func(ctx context.Context) error) *Checker {
	return &Checker{
		name:        name,
		description: description,
		checkFunc:   checkFunc,
	}
}

// Check performs the health check
func (c *Checker) Check(ctx context.Context) error {
	if c.checkFunc == nil {
		return fmt.Errorf("check function not implemented")
	}

	return c.checkFunc(ctx)
}

// Name returns the name of the checker
func (c *Checker) Name() string {
	return c.name
}

// Description returns the description of the checker
func (c *Checker) Description() string {
	return c.description
}

// DatabaseChecker represents a database health checker
type DatabaseChecker struct {
	name        string
	description string
	pingFunc    func(ctx context.Context) error
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(name, description string, pingFunc func(ctx context.Context) error) *DatabaseChecker {
	return &DatabaseChecker{
		name:        name,
		description: description,
		pingFunc:    pingFunc,
	}
}

// Check performs the database health check
func (d *DatabaseChecker) Check(ctx context.Context) error {
	if d.pingFunc == nil {
		return fmt.Errorf("ping function not implemented")
	}

	// Create a timeout context for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return d.pingFunc(pingCtx)
}

// Name returns the name of the checker
func (d *DatabaseChecker) Name() string {
	return d.name
}

// Description returns the description of the checker
func (d *DatabaseChecker) Description() string {
	return d.description
}

// HTTPChecker represents an HTTP health checker
type HTTPChecker struct {
	name        string
	description string
	url         string
	timeout     time.Duration
}

// NewHTTPChecker creates a new HTTP health checker
func NewHTTPChecker(name, description, url string, timeout time.Duration) *HTTPChecker {
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &HTTPChecker{
		name:        name,
		description: description,
		url:         url,
		timeout:     timeout,
	}
}

// Check performs the HTTP health check
func (h *HTTPChecker) Check(ctx context.Context) error {
	// Create a timeout context for the HTTP request
	_, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// This would typically make an HTTP request
	// For now, we'll just return a placeholder error
	return fmt.Errorf("HTTP health check not implemented for %s", h.url)
}

// Name returns the name of the checker
func (h *HTTPChecker) Name() string {
	return h.name
}

// Description returns the description of the checker
func (h *HTTPChecker) Description() string {
	return h.description
}

// CompositeChecker represents a composite health checker
type CompositeChecker struct {
	name        string
	description string
	checkers    []HealthChecker
}

// NewCompositeChecker creates a new composite health checker
func NewCompositeChecker(name, description string, checkers ...HealthChecker) *CompositeChecker {
	return &CompositeChecker{
		name:        name,
		description: description,
		checkers:    checkers,
	}
}

// Check performs all health checks
func (c *CompositeChecker) Check(ctx context.Context) error {
	for _, checker := range c.checkers {
		if err := checker.Check(ctx); err != nil {
			return fmt.Errorf("check %s failed: %w", checker.Name(), err)
		}
	}

	return nil
}

// Name returns the name of the checker
func (c *CompositeChecker) Name() string {
	return c.name
}

// Description returns the description of the checker
func (c *CompositeChecker) Description() string {
	return c.description
}

// AddChecker adds a checker to the composite checker
func (c *CompositeChecker) AddChecker(checker HealthChecker) {
	c.checkers = append(c.checkers, checker)
}

// RemoveChecker removes a checker from the composite checker
func (c *CompositeChecker) RemoveChecker(name string) {
	for i, checker := range c.checkers {
		if checker.Name() == name {
			c.checkers = append(c.checkers[:i], c.checkers[i+1:]...)
			break
		}
	}
}

// GetCheckers returns all checkers
func (c *CompositeChecker) GetCheckers() []HealthChecker {
	return c.checkers
}

// PeriodicChecker represents a periodic health checker
type PeriodicChecker struct {
	name        string
	description string
	checker     HealthChecker
	interval    time.Duration
	lastResult  CheckResult
	mu          sync.RWMutex
}

// NewPeriodicChecker creates a new periodic health checker
func NewPeriodicChecker(name, description string, checker HealthChecker, interval time.Duration) *PeriodicChecker {
	return &PeriodicChecker{
		name:        name,
		description: description,
		checker:     checker,
		interval:    interval,
	}
}

// Check performs the health check
func (p *PeriodicChecker) Check(ctx context.Context) error {
	p.mu.RLock()
	lastResult := p.lastResult
	p.mu.RUnlock()

	// If the last check was recent and successful, return cached result
	if time.Since(lastResult.Timestamp) < p.interval && lastResult.Status == StatusHealthy {
		return nil
	}

	// Perform the actual check
	err := p.checker.Check(ctx)

	// Update the cached result
	p.mu.Lock()
	p.lastResult = CheckResult{
		Name:        p.name,
		Description: p.description,
		Status:      StatusHealthy,
		Timestamp:   time.Now(),
	}

	if err != nil {
		p.lastResult.Status = StatusUnhealthy
		p.lastResult.Error = err.Error()
	}
	p.mu.Unlock()

	return err
}

// Name returns the name of the checker
func (p *PeriodicChecker) Name() string {
	return p.name
}

// Description returns the description of the checker
func (p *PeriodicChecker) Description() string {
	return p.description
}

// GetLastResult returns the last check result
func (p *PeriodicChecker) GetLastResult() CheckResult {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastResult
}
