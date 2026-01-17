// Package cache provides caching utilities for USC platform services.
package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CachePattern represents a caching pattern
type CachePattern string

const (
	// CacheAside pattern - application manages cache
	CacheAside CachePattern = "cache_aside"
	// WriteThrough pattern - write to cache and database simultaneously
	WriteThrough CachePattern = "write_through"
	// WriteBehind pattern - write to cache first, then database asynchronously
	WriteBehind CachePattern = "write_behind"
	// ReadThrough pattern - cache loads data from database on miss
	ReadThrough CachePattern = "read_through"
	// RefreshAhead pattern - cache refreshes data before expiration
	RefreshAhead CachePattern = "refresh_ahead"
)

// CacheManager manages multiple cache instances with different patterns
type CacheManager struct {
	caches   map[string]Cache
	patterns map[string]CachePattern
	mu       sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches:   make(map[string]Cache),
		patterns: make(map[string]CachePattern),
	}
}

// AddCache adds a cache instance with a specific pattern
func (cm *CacheManager) AddCache(name string, cache Cache, pattern CachePattern) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.caches[name] = cache
	cm.patterns[name] = pattern
}

// GetCache returns a cache instance by name
func (cm *CacheManager) GetCache(name string) (Cache, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cache, exists := cm.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache %s not found", name)
	}
	return cache, nil
}

// GetPattern returns the pattern for a cache
func (cm *CacheManager) GetPattern(name string) (CachePattern, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	pattern, exists := cm.patterns[name]
	if !exists {
		return "", fmt.Errorf("pattern for cache %s not found", name)
	}
	return pattern, nil
}

// CacheAsidePattern implements the Cache-Aside pattern
type CacheAsidePattern struct {
	cache Cache
	load  func(string) (interface{}, error)
	ttl   time.Duration
}

// NewCacheAsidePattern creates a new cache-aside pattern
func NewCacheAsidePattern(cache Cache, load func(string) (interface{}, error)) *CacheAsidePattern {
	return &CacheAsidePattern{
		cache: cache,
		load:  load,
		ttl:   5 * time.Minute, // Default TTL
	}
}

// Get retrieves data using cache-aside pattern
func (cap *CacheAsidePattern) Get(ctx context.Context, key string) (string, error) {
	// Try to get from cache first
	value, err := cap.cache.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// Cache miss, load from data source
	loadedValue, err := cap.load(key)
	if err != nil {
		return "", err
	}

	// Store in cache for next time
	cap.cache.Set(ctx, key, loadedValue, cap.ttl)

	// Convert to string
	return fmt.Sprintf("%v", loadedValue), nil
}

// Set stores data using cache-aside pattern
func (cap *CacheAsidePattern) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Store in cache
	return cap.cache.Set(ctx, key, value, expiration)
}

// Delete removes data using cache-aside pattern
func (cap *CacheAsidePattern) Delete(ctx context.Context, key string) error {
	// Remove from cache
	return cap.cache.Delete(ctx, key)
}

// Exists checks if a key exists in cache
func (cap *CacheAsidePattern) Exists(ctx context.Context, key string) (bool, error) {
	return cap.cache.Exists(ctx, key)
}

// Expire sets expiration for a key
func (cap *CacheAsidePattern) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return cap.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (cap *CacheAsidePattern) TTL(ctx context.Context, key string) (time.Duration, error) {
	return cap.cache.TTL(ctx, key)
}

// Increment increments a numeric value
func (cap *CacheAsidePattern) Increment(ctx context.Context, key string) (int64, error) {
	return cap.cache.Increment(ctx, key)
}

// Decrement decrements a numeric value
func (cap *CacheAsidePattern) Decrement(ctx context.Context, key string) (int64, error) {
	return cap.cache.Decrement(ctx, key)
}

// Keys returns all keys matching a pattern
func (cap *CacheAsidePattern) Keys(ctx context.Context, pattern string) ([]string, error) {
	return cap.cache.Keys(ctx, pattern)
}

// Flush removes all keys from cache
func (cap *CacheAsidePattern) Flush(ctx context.Context) error {
	return cap.cache.Flush(ctx)
}

func (cap *CacheAsidePattern) Health(ctx context.Context) error {
	return cap.cache.Health(ctx)
}

// WriteThroughPattern implements the Write-Through pattern
type WriteThroughPattern struct {
	cache Cache
	store func(string, interface{}) error
	ttl   time.Duration
}

// NewWriteThroughPattern creates a new write-through pattern
func NewWriteThroughPattern(cache Cache, store func(string, interface{}) error) *WriteThroughPattern {
	return &WriteThroughPattern{
		cache: cache,
		store: store,
		ttl:   5 * time.Minute, // Default TTL
	}
}

// Set stores data using write-through pattern
func (wtp *WriteThroughPattern) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Store in cache
	if err := wtp.cache.Set(ctx, key, value, expiration); err != nil {
		return err
	}

	// Store in data source
	if err := wtp.store(key, value); err != nil {
		// If store fails, remove from cache
		wtp.cache.Delete(ctx, key)
		return err
	}

	return nil
}

// Delete removes data using write-through pattern
func (wtp *WriteThroughPattern) Delete(ctx context.Context, key string) error {
	// Remove from cache
	if err := wtp.cache.Delete(ctx, key); err != nil {
		return err
	}

	// Remove from data source
	// Note: In a real implementation, you would have a delete function
	// For now, we'll just return the cache delete result
	return nil
}

// Get retrieves data using write-through pattern
func (wtp *WriteThroughPattern) Get(ctx context.Context, key string) (string, error) {
	return wtp.cache.Get(ctx, key)
}

// Exists checks if a key exists in cache
func (wtp *WriteThroughPattern) Exists(ctx context.Context, key string) (bool, error) {
	return wtp.cache.Exists(ctx, key)
}

// Expire sets expiration for a key
func (wtp *WriteThroughPattern) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return wtp.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (wtp *WriteThroughPattern) TTL(ctx context.Context, key string) (time.Duration, error) {
	return wtp.cache.TTL(ctx, key)
}

// Increment increments a numeric value
func (wtp *WriteThroughPattern) Increment(ctx context.Context, key string) (int64, error) {
	return wtp.cache.Increment(ctx, key)
}

// Decrement decrements a numeric value
func (wtp *WriteThroughPattern) Decrement(ctx context.Context, key string) (int64, error) {
	return wtp.cache.Decrement(ctx, key)
}

// Keys returns all keys matching a pattern
func (wtp *WriteThroughPattern) Keys(ctx context.Context, pattern string) ([]string, error) {
	return wtp.cache.Keys(ctx, pattern)
}

// Flush removes all keys from cache
func (wtp *WriteThroughPattern) Flush(ctx context.Context) error {
	return wtp.cache.Flush(ctx)
}

func (wtp *WriteThroughPattern) Health(ctx context.Context) error {
	return wtp.cache.Health(ctx)
}

// WriteBehindPattern implements the Write-Behind pattern
type WriteBehindPattern struct {
	cache     Cache
	store     func(string, interface{}) error
	batchSize int
	queue     chan writeOperation
	ttl       time.Duration
}

// writeOperation represents a write operation
type writeOperation struct {
	key   string
	value interface{}
}

// NewWriteBehindPattern creates a new write-behind pattern
func NewWriteBehindPattern(cache Cache, store func(string, interface{}) error, batchSize int) *WriteBehindPattern {
	wbp := &WriteBehindPattern{
		cache:     cache,
		store:     store,
		batchSize: batchSize,
		queue:     make(chan writeOperation, 1000), // Buffer for 1000 operations
		ttl:       5 * time.Minute,                 // Default TTL
	}

	// Start background writer
	go wbp.backgroundWriter()

	return wbp
}

// Set stores data using write-behind pattern
func (wbp *WriteBehindPattern) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Store in cache immediately
	if err := wbp.cache.Set(ctx, key, value, expiration); err != nil {
		return err
	}

	// Queue for background write
	select {
	case wbp.queue <- writeOperation{key: key, value: value}:
		return nil
	default:
		// Queue is full, handle error
		return fmt.Errorf("write queue is full")
	}
}

// Get retrieves data using write-behind pattern
func (wbp *WriteBehindPattern) Get(ctx context.Context, key string) (string, error) {
	return wbp.cache.Get(ctx, key)
}

// Delete removes data using write-behind pattern
func (wbp *WriteBehindPattern) Delete(ctx context.Context, key string) error {
	// Remove from cache
	return wbp.cache.Delete(ctx, key)
}

// Exists checks if a key exists in cache
func (wbp *WriteBehindPattern) Exists(ctx context.Context, key string) (bool, error) {
	return wbp.cache.Exists(ctx, key)
}

// Expire sets expiration for a key
func (wbp *WriteBehindPattern) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return wbp.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (wbp *WriteBehindPattern) TTL(ctx context.Context, key string) (time.Duration, error) {
	return wbp.cache.TTL(ctx, key)
}

// Increment increments a numeric value
func (wbp *WriteBehindPattern) Increment(ctx context.Context, key string) (int64, error) {
	return wbp.cache.Increment(ctx, key)
}

// Decrement decrements a numeric value
func (wbp *WriteBehindPattern) Decrement(ctx context.Context, key string) (int64, error) {
	return wbp.cache.Decrement(ctx, key)
}

// Keys returns all keys matching a pattern
func (wbp *WriteBehindPattern) Keys(ctx context.Context, pattern string) ([]string, error) {
	return wbp.cache.Keys(ctx, pattern)
}

// Flush removes all keys from cache
func (wbp *WriteBehindPattern) Flush(ctx context.Context) error {
	return wbp.cache.Flush(ctx)
}

func (wbp *WriteBehindPattern) Health(ctx context.Context) error {
	return wbp.cache.Health(ctx)
}

// backgroundWriter processes write operations in the background
func (wbp *WriteBehindPattern) backgroundWriter() {
	batch := make([]writeOperation, 0, wbp.batchSize)
	ticker := time.NewTicker(5 * time.Second) // Flush every 5 seconds

	for {
		select {
		case op := <-wbp.queue:
			batch = append(batch, op)
			if len(batch) >= wbp.batchSize {
				wbp.flushBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				wbp.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// flushBatch flushes a batch of write operations
func (wbp *WriteBehindPattern) flushBatch(batch []writeOperation) {
	for _, op := range batch {
		if err := wbp.store(op.key, op.value); err != nil {
			// Log error in real implementation
			continue
		}
	}
}

// ReadThroughPattern implements the Read-Through pattern
type ReadThroughPattern struct {
	cache Cache
	load  func(string) (interface{}, error)
	ttl   time.Duration
}

// NewReadThroughPattern creates a new read-through pattern
func NewReadThroughPattern(cache Cache, load func(string) (interface{}, error)) *ReadThroughPattern {
	return &ReadThroughPattern{
		cache: cache,
		load:  load,
		ttl:   5 * time.Minute, // Default TTL
	}
}

// Get retrieves data using read-through pattern
func (rtp *ReadThroughPattern) Get(ctx context.Context, key string) (string, error) {
	// Try to get from cache first
	value, err := rtp.cache.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// Cache miss, load from data source
	loadedValue, err := rtp.load(key)
	if err != nil {
		return "", err
	}

	// Store in cache
	rtp.cache.Set(ctx, key, loadedValue, rtp.ttl)

	// Convert to string
	return fmt.Sprintf("%v", loadedValue), nil
}

// Set stores data using read-through pattern
func (rtp *ReadThroughPattern) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rtp.cache.Set(ctx, key, value, expiration)
}

// Delete removes data using read-through pattern
func (rtp *ReadThroughPattern) Delete(ctx context.Context, key string) error {
	return rtp.cache.Delete(ctx, key)
}

// Exists checks if a key exists in cache
func (rtp *ReadThroughPattern) Exists(ctx context.Context, key string) (bool, error) {
	return rtp.cache.Exists(ctx, key)
}

// Expire sets expiration for a key
func (rtp *ReadThroughPattern) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rtp.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (rtp *ReadThroughPattern) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rtp.cache.TTL(ctx, key)
}

// Increment increments a numeric value
func (rtp *ReadThroughPattern) Increment(ctx context.Context, key string) (int64, error) {
	return rtp.cache.Increment(ctx, key)
}

// Decrement decrements a numeric value
func (rtp *ReadThroughPattern) Decrement(ctx context.Context, key string) (int64, error) {
	return rtp.cache.Decrement(ctx, key)
}

// Keys returns all keys matching a pattern
func (rtp *ReadThroughPattern) Keys(ctx context.Context, pattern string) ([]string, error) {
	return rtp.cache.Keys(ctx, pattern)
}

// Flush removes all keys from cache
func (rtp *ReadThroughPattern) Flush(ctx context.Context) error {
	return rtp.cache.Flush(ctx)
}

func (rtp *ReadThroughPattern) Health(ctx context.Context) error {
	return rtp.cache.Health(ctx)
}

// RefreshAheadPattern implements the Refresh-Ahead pattern
type RefreshAheadPattern struct {
	cache      Cache
	load       func(string) (interface{}, error)
	refreshTTL time.Duration
	refreshCh  chan string
	ttl        time.Duration
}

// NewRefreshAheadPattern creates a new refresh-ahead pattern
func NewRefreshAheadPattern(cache Cache, load func(string) (interface{}, error), refreshTTL time.Duration) *RefreshAheadPattern {
	rap := &RefreshAheadPattern{
		cache:      cache,
		load:       load,
		refreshTTL: refreshTTL,
		refreshCh:  make(chan string, 1000),
		ttl:        5 * time.Minute, // Default TTL
	}

	// Start background refresher
	go rap.backgroundRefresher()

	return rap
}

// Get retrieves data using refresh-ahead pattern
func (rap *RefreshAheadPattern) Get(ctx context.Context, key string) (string, error) {
	// Try to get from cache first
	value, err := rap.cache.Get(ctx, key)
	if err == nil {
		// Check if we need to refresh
		ttl, err := rap.cache.TTL(ctx, key)
		if err == nil && ttl < rap.refreshTTL {
			// Queue for refresh
			select {
			case rap.refreshCh <- key:
			default:
				// Queue is full, ignore
			}
		}
		return value, nil
	}

	// Cache miss, load from data source
	loadedValue, err := rap.load(key)
	if err != nil {
		return "", err
	}

	// Store in cache
	rap.cache.Set(ctx, key, loadedValue, rap.ttl)

	// Convert to string
	return fmt.Sprintf("%v", loadedValue), nil
}

// Set stores data using refresh-ahead pattern
func (rap *RefreshAheadPattern) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rap.cache.Set(ctx, key, value, expiration)
}

// Delete removes data using refresh-ahead pattern
func (rap *RefreshAheadPattern) Delete(ctx context.Context, key string) error {
	return rap.cache.Delete(ctx, key)
}

// Exists checks if a key exists in cache
func (rap *RefreshAheadPattern) Exists(ctx context.Context, key string) (bool, error) {
	return rap.cache.Exists(ctx, key)
}

// Expire sets expiration for a key
func (rap *RefreshAheadPattern) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rap.cache.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (rap *RefreshAheadPattern) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rap.cache.TTL(ctx, key)
}

// Increment increments a numeric value
func (rap *RefreshAheadPattern) Increment(ctx context.Context, key string) (int64, error) {
	return rap.cache.Increment(ctx, key)
}

// Decrement decrements a numeric value
func (rap *RefreshAheadPattern) Decrement(ctx context.Context, key string) (int64, error) {
	return rap.cache.Decrement(ctx, key)
}

// Keys returns all keys matching a pattern
func (rap *RefreshAheadPattern) Keys(ctx context.Context, pattern string) ([]string, error) {
	return rap.cache.Keys(ctx, pattern)
}

// Flush removes all keys from cache
func (rap *RefreshAheadPattern) Flush(ctx context.Context) error {
	return rap.cache.Flush(ctx)
}

func (rap *RefreshAheadPattern) Health(ctx context.Context) error {
	return rap.cache.Health(ctx)
}

// backgroundRefresher refreshes data in the background
func (rap *RefreshAheadPattern) backgroundRefresher() {
	for key := range rap.refreshCh {
		// Load fresh data
		value, err := rap.load(key)
		if err != nil {
			continue
		}

		// Store in cache
		rap.cache.Set(context.Background(), key, value, rap.ttl)
	}
}

// CacheFactory creates cache instances with different patterns
type CacheFactory struct {
	patterns map[CachePattern]func(Cache, interface{}) Cache
}

// NewCacheFactory creates a new cache factory
func NewCacheFactory() *CacheFactory {
	cf := &CacheFactory{
		patterns: make(map[CachePattern]func(Cache, interface{}) Cache),
	}

	// Register default patterns
	cf.RegisterPattern(CacheAside, func(cache Cache, config interface{}) Cache {
		load := config.(func(string) (interface{}, error))
		return NewCacheAsidePattern(cache, load)
	})

	cf.RegisterPattern(WriteThrough, func(cache Cache, config interface{}) Cache {
		store := config.(func(string, interface{}) error)
		return NewWriteThroughPattern(cache, store)
	})

	cf.RegisterPattern(WriteBehind, func(cache Cache, config interface{}) Cache {
		configMap := config.(map[string]interface{})
		store := configMap["store"].(func(string, interface{}) error)
		batchSize := configMap["batch_size"].(int)
		return NewWriteBehindPattern(cache, store, batchSize)
	})

	cf.RegisterPattern(ReadThrough, func(cache Cache, config interface{}) Cache {
		load := config.(func(string) (interface{}, error))
		return NewReadThroughPattern(cache, load)
	})

	cf.RegisterPattern(RefreshAhead, func(cache Cache, config interface{}) Cache {
		configMap := config.(map[string]interface{})
		load := configMap["load"].(func(string) (interface{}, error))
		refreshTTL := configMap["refresh_ttl"].(time.Duration)
		return NewRefreshAheadPattern(cache, load, refreshTTL)
	})

	return cf
}

// RegisterPattern registers a new cache pattern
func (cf *CacheFactory) RegisterPattern(pattern CachePattern, factory func(Cache, interface{}) Cache) {
	cf.patterns[pattern] = factory
}

// CreateCache creates a cache instance with the specified pattern
func (cf *CacheFactory) CreateCache(pattern CachePattern, baseCache Cache, config interface{}) (Cache, error) {
	factory, exists := cf.patterns[pattern]
	if !exists {
		return nil, fmt.Errorf("pattern %s not registered", pattern)
	}

	return factory(baseCache, config), nil
}

// MultiTierCacheConfig represents multi-tier cache configuration
type MultiTierCacheConfig struct {
	L1Config L1CacheConfig `mapstructure:"l1"`
	L2Config L2CacheConfig `mapstructure:"l2"`
	L4Config L4CacheConfig `mapstructure:"l4"`
}

// L1CacheConfig represents L1 (in-memory) cache configuration
type L1CacheConfig struct {
	MaxSize         int           `mapstructure:"max_size"`
	DefaultTTL      time.Duration `mapstructure:"default_ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// L2CacheConfig represents L2 (Redis) cache configuration
type L2CacheConfig struct {
	Address    string        `mapstructure:"address"`
	Password   string        `mapstructure:"password"`
	DB         int           `mapstructure:"db"`
	DefaultTTL time.Duration `mapstructure:"default_ttl"`
	MaxRetries int           `mapstructure:"max_retries"`
}

// L4CacheConfig represents L4 (CDN) cache configuration
type L4CacheConfig struct {
	CDNURL      string        `mapstructure:"cdn_url"`
	DefaultTTL  time.Duration `mapstructure:"default_ttl"`
	MaxFileSize int64         `mapstructure:"max_file_size"`
}

// DefaultMultiTierCacheConfig returns default multi-tier cache configuration
func DefaultMultiTierCacheConfig() MultiTierCacheConfig {
	return MultiTierCacheConfig{
		L1Config: L1CacheConfig{
			MaxSize:         1000,
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
		L2Config: L2CacheConfig{
			Address:    "localhost:6379",
			Password:   "",
			DB:         0,
			DefaultTTL: 30 * time.Minute,
			MaxRetries: 3,
		},
		L4Config: L4CacheConfig{
			CDNURL:      "",
			DefaultTTL:  24 * time.Hour,
			MaxFileSize: 10 * 1024 * 1024, // 10MB
		},
	}
}

// MultiTierCache implements a three-tier caching system
type MultiTierCache struct {
	config MultiTierCacheConfig
	l1     *L1Cache
	l2     *L2Cache
	l4     *L4Cache
	mu     sync.RWMutex
}

// NewMultiTierCache creates a new multi-tier cache
func NewMultiTierCache(config MultiTierCacheConfig) (*MultiTierCache, error) {
	// Initialize L1 cache (in-memory)
	l1, err := NewL1Cache(config.L1Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 cache: %w", err)
	}

	// Initialize L2 cache (Redis)
	l2, err := NewL2Cache(config.L2Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create L2 cache: %w", err)
	}

	// Initialize L4 cache (CDN)
	l4, err := NewL4Cache(config.L4Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create L4 cache: %w", err)
	}

	return &MultiTierCache{
		config: config,
		l1:     l1,
		l2:     l2,
		l4:     l4,
	}, nil
}

// Get retrieves a value from the multi-tier cache
func (mtc *MultiTierCache) Get(ctx context.Context, key string) (string, error) {
	// Try L1 cache first
	value, err := mtc.l1.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// Try L2 cache
	value, err = mtc.l2.Get(ctx, key)
	if err == nil {
		// Store in L1 for faster access
		mtc.l1.Set(ctx, key, value, mtc.config.L1Config.DefaultTTL)
		return value, nil
	}

	// Try L4 cache
	value, err = mtc.l4.Get(ctx, key)
	if err == nil {
		// Store in L2 and L1
		mtc.l2.Set(ctx, key, value, mtc.config.L2Config.DefaultTTL)
		mtc.l1.Set(ctx, key, value, mtc.config.L1Config.DefaultTTL)
		return value, nil
	}

	return "", fmt.Errorf("key %s not found in any cache tier", key)
}

// Set stores a value in the multi-tier cache
func (mtc *MultiTierCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Store in all tiers
	if err := mtc.l1.Set(ctx, key, value, expiration); err != nil {
		return fmt.Errorf("failed to set L1 cache: %w", err)
	}

	if err := mtc.l2.Set(ctx, key, value, expiration); err != nil {
		return fmt.Errorf("failed to set L2 cache: %w", err)
	}

	if err := mtc.l4.Set(ctx, key, value, expiration); err != nil {
		return fmt.Errorf("failed to set L4 cache: %w", err)
	}

	return nil
}

// Delete removes a value from all cache tiers
func (mtc *MultiTierCache) Delete(ctx context.Context, key string) error {
	var errors []error

	if err := mtc.l1.Delete(ctx, key); err != nil {
		errors = append(errors, fmt.Errorf("failed to delete from L1: %w", err))
	}

	if err := mtc.l2.Delete(ctx, key); err != nil {
		errors = append(errors, fmt.Errorf("failed to delete from L2: %w", err))
	}

	if err := mtc.l4.Delete(ctx, key); err != nil {
		errors = append(errors, fmt.Errorf("failed to delete from L4: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors deleting from cache: %v", errors)
	}

	return nil
}

// Close closes the multi-tier cache
func (mtc *MultiTierCache) Close() error {
	// In a real implementation, you would close connections here
	// For now, we just clear the data
	mtc.mu.Lock()
	defer mtc.mu.Unlock()

	mtc.l1.data = make(map[string]multiTierCacheItem)
	mtc.l2.data = make(map[string]multiTierCacheItem)
	mtc.l4.data = make(map[string]multiTierCacheItem)

	return nil
}

// L1Cache implements in-memory cache
type L1Cache struct {
	config L1CacheConfig
	data   map[string]multiTierCacheItem
	mu     sync.RWMutex
}

type multiTierCacheItem struct {
	value      string
	expiration time.Time
}

// NewL1Cache creates a new L1 cache
func NewL1Cache(config L1CacheConfig) (*L1Cache, error) {
	cache := &L1Cache{
		config: config,
		data:   make(map[string]multiTierCacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache, nil
}

// Get retrieves a value from L1 cache
func (l1 *L1Cache) Get(ctx context.Context, key string) (string, error) {
	l1.mu.RLock()
	defer l1.mu.RUnlock()

	item, exists := l1.data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	if time.Now().After(item.expiration) {
		return "", fmt.Errorf("key %s expired", key)
	}

	return item.value, nil
}

// Set stores a value in L1 cache
func (l1 *L1Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	l1.mu.Lock()
	defer l1.mu.Unlock()

	// Check size limit
	if len(l1.data) >= l1.config.MaxSize {
		// Remove oldest item (simple implementation)
		for k := range l1.data {
			delete(l1.data, k)
			break
		}
	}

	l1.data[key] = multiTierCacheItem{
		value:      fmt.Sprintf("%v", value),
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// Delete removes a value from L1 cache
func (l1 *L1Cache) Delete(ctx context.Context, key string) error {
	l1.mu.Lock()
	defer l1.mu.Unlock()

	delete(l1.data, key)
	return nil
}

// Health performs health check on L1 cache
func (l1 *L1Cache) Health(ctx context.Context) error {
	l1.mu.RLock()
	defer l1.mu.RUnlock()

	// Simple health check - just verify cache is accessible
	if l1.data == nil {
		return fmt.Errorf("L1 cache data is nil")
	}

	return nil
}

// cleanup removes expired items
func (l1 *L1Cache) cleanup() {
	ticker := time.NewTicker(l1.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l1.mu.Lock()
		now := time.Now()
		for key, item := range l1.data {
			if now.After(item.expiration) {
				delete(l1.data, key)
			}
		}
		l1.mu.Unlock()
	}
}

// L2Cache implements Redis cache
type L2Cache struct {
	config L2CacheConfig
	// In a real implementation, this would be a Redis client
	// For now, we'll use a simple in-memory store
	data map[string]multiTierCacheItem
	mu   sync.RWMutex
}

// NewL2Cache creates a new L2 cache
func NewL2Cache(config L2CacheConfig) (*L2Cache, error) {
	return &L2Cache{
		config: config,
		data:   make(map[string]multiTierCacheItem),
	}, nil
}

// Get retrieves a value from L2 cache
func (l2 *L2Cache) Get(ctx context.Context, key string) (string, error) {
	l2.mu.RLock()
	defer l2.mu.RUnlock()

	item, exists := l2.data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	if time.Now().After(item.expiration) {
		return "", fmt.Errorf("key %s expired", key)
	}

	return item.value, nil
}

// Set stores a value in L2 cache
func (l2 *L2Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	l2.mu.Lock()
	defer l2.mu.Unlock()

	l2.data[key] = multiTierCacheItem{
		value:      fmt.Sprintf("%v", value),
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// Delete removes a value from L2 cache
func (l2 *L2Cache) Delete(ctx context.Context, key string) error {
	l2.mu.Lock()
	defer l2.mu.Unlock()

	delete(l2.data, key)
	return nil
}

// Health performs health check on L2 cache
func (l2 *L2Cache) Health(ctx context.Context) error {
	l2.mu.RLock()
	defer l2.mu.RUnlock()

	// Simple health check - just verify cache is accessible
	if l2.data == nil {
		return fmt.Errorf("L2 cache data is nil")
	}

	return nil
}

// L4Cache implements CDN cache
type L4Cache struct {
	config L4CacheConfig
	// In a real implementation, this would be a CDN client
	// For now, we'll use a simple in-memory store
	data map[string]multiTierCacheItem
	mu   sync.RWMutex
}

// NewL4Cache creates a new L4 cache
func NewL4Cache(config L4CacheConfig) (*L4Cache, error) {
	return &L4Cache{
		config: config,
		data:   make(map[string]multiTierCacheItem),
	}, nil
}

// Get retrieves a value from L4 cache
func (l4 *L4Cache) Get(ctx context.Context, key string) (string, error) {
	l4.mu.RLock()
	defer l4.mu.RUnlock()

	item, exists := l4.data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	if time.Now().After(item.expiration) {
		return "", fmt.Errorf("key %s expired", key)
	}

	return item.value, nil
}

// Set stores a value in L4 cache
func (l4 *L4Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	l4.mu.Lock()
	defer l4.mu.Unlock()

	l4.data[key] = multiTierCacheItem{
		value:      fmt.Sprintf("%v", value),
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// Delete removes a value from L4 cache
func (l4 *L4Cache) Delete(ctx context.Context, key string) error {
	l4.mu.Lock()
	defer l4.mu.Unlock()

	delete(l4.data, key)
	return nil
}

// Health performs health check on L4 cache
func (l4 *L4Cache) Health(ctx context.Context) error {
	l4.mu.RLock()
	defer l4.mu.RUnlock()

	// Simple health check - just verify cache is accessible
	if l4.data == nil {
		return fmt.Errorf("L4 cache data is nil")
	}

	return nil
}

// Enhanced Multi-Tier Cache Features

// CacheStatistics represents cache performance statistics
type CacheStatistics struct {
	L1Hits      int64 `json:"l1_hits"`
	L1Misses    int64 `json:"l1_misses"`
	L2Hits      int64 `json:"l2_hits"`
	L2Misses    int64 `json:"l2_misses"`
	L4Hits      int64 `json:"l4_hits"`
	L4Misses    int64 `json:"l4_misses"`
	TotalHits   int64 `json:"total_hits"`
	TotalMisses int64 `json:"total_misses"`
	Writes      int64 `json:"writes"`
	Deletes     int64 `json:"deletes"`
	Errors      int64 `json:"errors"`
}

// MultiTierCacheMetrics represents multi-tier cache metrics and monitoring
type MultiTierCacheMetrics struct {
	mu sync.RWMutex

	// Statistics
	stats CacheStatistics

	// Performance tracking
	responseTimes map[string][]time.Duration
	errorRates    map[string]int64

	// Health monitoring
	healthChecks map[string]time.Time
	lastHealth   time.Time
}

// NewMultiTierCacheMetrics creates a new multi-tier cache metrics instance
func NewMultiTierCacheMetrics() *MultiTierCacheMetrics {
	return &MultiTierCacheMetrics{
		responseTimes: make(map[string][]time.Duration),
		errorRates:    make(map[string]int64),
		healthChecks:  make(map[string]time.Time),
	}
}

// RecordHit records a cache hit
func (cm *MultiTierCacheMetrics) RecordHit(tier string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	switch tier {
	case "l1":
		cm.stats.L1Hits++
	case "l2":
		cm.stats.L2Hits++
	case "l4":
		cm.stats.L4Hits++
	}
	cm.stats.TotalHits++
}

// RecordMiss records a cache miss
func (cm *MultiTierCacheMetrics) RecordMiss(tier string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	switch tier {
	case "l1":
		cm.stats.L1Misses++
	case "l2":
		cm.stats.L2Misses++
	case "l4":
		cm.stats.L4Misses++
	}
	cm.stats.TotalMisses++
}

// RecordWrite records a cache write
func (cm *MultiTierCacheMetrics) RecordWrite() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.stats.Writes++
}

// RecordDelete records a cache delete
func (cm *MultiTierCacheMetrics) RecordDelete() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.stats.Deletes++
}

// RecordError records a cache error
func (cm *MultiTierCacheMetrics) RecordError(tier string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.stats.Errors++
	cm.errorRates[tier]++
}

// RecordResponseTime records response time for a tier
func (cm *MultiTierCacheMetrics) RecordResponseTime(tier string, duration time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.responseTimes[tier] == nil {
		cm.responseTimes[tier] = make([]time.Duration, 0, 100)
	}

	cm.responseTimes[tier] = append(cm.responseTimes[tier], duration)

	// Keep only last 100 measurements
	if len(cm.responseTimes[tier]) > 100 {
		cm.responseTimes[tier] = cm.responseTimes[tier][1:]
	}
}

// GetStats returns current statistics
func (cm *MultiTierCacheMetrics) GetStats() CacheStatistics {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.stats
}

// GetHitRate returns overall hit rate
func (cm *MultiTierCacheMetrics) GetHitRate() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total == 0 {
		return 0
	}

	return float64(cm.stats.TotalHits) / float64(total) * 100
}

// GetTierHitRates returns hit rates for each tier
func (cm *MultiTierCacheMetrics) GetTierHitRates() map[string]float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	rates := make(map[string]float64)

	// L1 hit rate
	l1Total := cm.stats.L1Hits + cm.stats.L1Misses
	if l1Total > 0 {
		rates["l1"] = float64(cm.stats.L1Hits) / float64(l1Total) * 100
	}

	// L2 hit rate
	l2Total := cm.stats.L2Hits + cm.stats.L2Misses
	if l2Total > 0 {
		rates["l2"] = float64(cm.stats.L2Hits) / float64(l2Total) * 100
	}

	// L4 hit rate
	l4Total := cm.stats.L4Hits + cm.stats.L4Misses
	if l4Total > 0 {
		rates["l4"] = float64(cm.stats.L4Hits) / float64(l4Total) * 100
	}

	return rates
}

// GetAverageResponseTime returns average response time for a tier
func (cm *MultiTierCacheMetrics) GetAverageResponseTime(tier string) time.Duration {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	times := cm.responseTimes[tier]
	if len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	return total / time.Duration(len(times))
}

// GetErrorRate returns error rate for a tier
func (cm *MultiTierCacheMetrics) GetErrorRate(tier string) float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	errors := cm.errorRates[tier]
	total := cm.stats.TotalHits + cm.stats.TotalMisses

	if total == 0 {
		return 0
	}

	return float64(errors) / float64(total) * 100
}

// EnhancedMultiTierCache represents an enhanced multi-tier cache with advanced features
type EnhancedMultiTierCache struct {
	*MultiTierCache
	metrics *MultiTierCacheMetrics
	config  EnhancedMultiTierConfig
	mu      sync.RWMutex
}

// EnhancedMultiTierConfig represents enhanced multi-tier cache configuration
type EnhancedMultiTierConfig struct {
	MultiTierCacheConfig

	// Advanced features
	EnableMetrics     bool          `mapstructure:"enable_metrics"`
	EnableHealthCheck bool          `mapstructure:"enable_health_check"`
	HealthInterval    time.Duration `mapstructure:"health_interval"`
	EnableAutoWarmup  bool          `mapstructure:"enable_auto_warmup"`
	WarmupInterval    time.Duration `mapstructure:"warmup_interval"`
	EnableCompression bool          `mapstructure:"enable_compression"`
	EnableEncryption  bool          `mapstructure:"enable_encryption"`

	// Performance tuning
	MaxRetries              int           `mapstructure:"max_retries"`
	RetryDelay              time.Duration `mapstructure:"retry_delay"`
	CircuitBreakerThreshold int           `mapstructure:"circuit_breaker_threshold"`
}

// DefaultEnhancedMultiTierConfig returns default enhanced configuration
func DefaultEnhancedMultiTierConfig() EnhancedMultiTierConfig {
	return EnhancedMultiTierConfig{
		MultiTierCacheConfig:    DefaultMultiTierCacheConfig(),
		EnableMetrics:           true,
		EnableHealthCheck:       true,
		HealthInterval:          30 * time.Second,
		EnableAutoWarmup:        false,
		WarmupInterval:          5 * time.Minute,
		EnableCompression:       false,
		EnableEncryption:        false,
		MaxRetries:              3,
		RetryDelay:              100 * time.Millisecond,
		CircuitBreakerThreshold: 10,
	}
}

// NewEnhancedMultiTierCache creates a new enhanced multi-tier cache
func NewEnhancedMultiTierCache(config EnhancedMultiTierConfig) (*EnhancedMultiTierCache, error) {
	// Create base multi-tier cache
	baseCache, err := NewMultiTierCache(config.MultiTierCacheConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create base multi-tier cache: %w", err)
	}

	// Create metrics
	metrics := NewMultiTierCacheMetrics()

	return &EnhancedMultiTierCache{
		MultiTierCache: baseCache,
		metrics:        metrics,
		config:         config,
	}, nil
}

// Get retrieves a value with enhanced features
func (emtc *EnhancedMultiTierCache) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	defer func() {
		emtc.metrics.RecordResponseTime("get", time.Since(start))
	}()

	// Try L1 cache first
	value, err := emtc.l1.Get(ctx, key)
	if err == nil {
		emtc.metrics.RecordHit("l1")
		return value, nil
	}
	emtc.metrics.RecordMiss("l1")

	// Try L2 cache
	value, err = emtc.l2.Get(ctx, key)
	if err == nil {
		emtc.metrics.RecordHit("l2")
		// Store in L1 for faster access
		emtc.l1.Set(ctx, key, value, emtc.config.L1Config.DefaultTTL)
		return value, nil
	}
	emtc.metrics.RecordMiss("l2")

	// Try L4 cache
	value, err = emtc.l4.Get(ctx, key)
	if err == nil {
		emtc.metrics.RecordHit("l4")
		// Store in upper tiers for faster access
		emtc.l2.Set(ctx, key, value, emtc.config.L2Config.DefaultTTL)
		emtc.l1.Set(ctx, key, value, emtc.config.L1Config.DefaultTTL)
		return value, nil
	}
	emtc.metrics.RecordMiss("l4")

	return "", ErrCacheMiss
}

// Set stores a value with enhanced features
func (emtc *EnhancedMultiTierCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	defer func() {
		emtc.metrics.RecordResponseTime("set", time.Since(start))
		emtc.metrics.RecordWrite()
	}()

	// Set in all tiers
	if err := emtc.l1.Set(ctx, key, value, expiration); err != nil {
		emtc.metrics.RecordError("l1")
		return fmt.Errorf("L1 cache set failed: %w", err)
	}

	if err := emtc.l2.Set(ctx, key, value, expiration); err != nil {
		emtc.metrics.RecordError("l2")
		return fmt.Errorf("L2 cache set failed: %w", err)
	}

	if err := emtc.l4.Set(ctx, key, value, expiration); err != nil {
		emtc.metrics.RecordError("l4")
		return fmt.Errorf("L4 cache set failed: %w", err)
	}

	return nil
}

// Delete removes a value with enhanced features
func (emtc *EnhancedMultiTierCache) Delete(ctx context.Context, key string) error {
	start := time.Now()
	defer func() {
		emtc.metrics.RecordResponseTime("delete", time.Since(start))
		emtc.metrics.RecordDelete()
	}()

	// Delete from all tiers
	if err := emtc.l1.Delete(ctx, key); err != nil {
		emtc.metrics.RecordError("l1")
		return fmt.Errorf("L1 cache delete failed: %w", err)
	}

	if err := emtc.l2.Delete(ctx, key); err != nil {
		emtc.metrics.RecordError("l2")
		return fmt.Errorf("L2 cache delete failed: %w", err)
	}

	if err := emtc.l4.Delete(ctx, key); err != nil {
		emtc.metrics.RecordError("l4")
		return fmt.Errorf("L4 cache delete failed: %w", err)
	}

	return nil
}

// GetMetrics returns cache metrics
func (emtc *EnhancedMultiTierCache) GetMetrics() *MultiTierCacheMetrics {
	return emtc.metrics
}

// GetPerformanceReport returns a comprehensive performance report
func (emtc *EnhancedMultiTierCache) GetPerformanceReport() map[string]interface{} {
	emtc.mu.RLock()
	defer emtc.mu.RUnlock()

	report := make(map[string]interface{})

	// Basic statistics
	report["stats"] = emtc.metrics.GetStats()

	// Hit rates
	report["hit_rates"] = emtc.metrics.GetTierHitRates()
	report["overall_hit_rate"] = emtc.metrics.GetHitRate()

	// Response times
	responseTimes := make(map[string]time.Duration)
	responseTimes["l1"] = emtc.metrics.GetAverageResponseTime("l1")
	responseTimes["l2"] = emtc.metrics.GetAverageResponseTime("l2")
	responseTimes["l4"] = emtc.metrics.GetAverageResponseTime("l4")
	report["response_times"] = responseTimes

	// Error rates
	errorRates := make(map[string]float64)
	errorRates["l1"] = emtc.metrics.GetErrorRate("l1")
	errorRates["l2"] = emtc.metrics.GetErrorRate("l2")
	errorRates["l4"] = emtc.metrics.GetErrorRate("l4")
	report["error_rates"] = errorRates

	// Configuration
	report["config"] = emtc.config

	return report
}

// HealthCheck performs comprehensive health check
func (emtc *EnhancedMultiTierCache) HealthCheck(ctx context.Context) error {
	if !emtc.config.EnableHealthCheck {
		return nil
	}

	emtc.mu.Lock()
	defer emtc.mu.Unlock()

	var errors []error

	// Check L1 cache health
	if err := emtc.l1.Health(ctx); err != nil {
		errors = append(errors, fmt.Errorf("L1 cache health check failed: %w", err))
	} else {
		emtc.metrics.healthChecks["l1"] = time.Now()
	}

	// Check L2 cache health
	if err := emtc.l2.Health(ctx); err != nil {
		errors = append(errors, fmt.Errorf("L2 cache health check failed: %w", err))
	} else {
		emtc.metrics.healthChecks["l2"] = time.Now()
	}

	// Check L4 cache health
	if err := emtc.l4.Health(ctx); err != nil {
		errors = append(errors, fmt.Errorf("L4 cache health check failed: %w", err))
	} else {
		emtc.metrics.healthChecks["l4"] = time.Now()
	}

	emtc.metrics.lastHealth = time.Now()

	if len(errors) > 0 {
		return fmt.Errorf("health check errors: %v", errors)
	}

	return nil
}

// Warmup performs cache warmup
func (emtc *EnhancedMultiTierCache) Warmup(ctx context.Context, keys []string) error {
	if !emtc.config.EnableAutoWarmup {
		return nil
	}

	emtc.mu.Lock()
	defer emtc.mu.Unlock()

	for _, key := range keys {
		// Try to get value from L4 and populate upper tiers
		value, err := emtc.l4.Get(ctx, key)
		if err == nil {
			// Populate L2 and L1
			emtc.l2.Set(ctx, key, value, emtc.config.L2Config.DefaultTTL)
			emtc.l1.Set(ctx, key, value, emtc.config.L1Config.DefaultTTL)
		}
	}

	return nil
}

// Close closes the enhanced multi-tier cache
func (emtc *EnhancedMultiTierCache) Close() error {
	emtc.mu.Lock()
	defer emtc.mu.Unlock()

	// Close base cache
	if err := emtc.MultiTierCache.Close(); err != nil {
		return fmt.Errorf("failed to close base cache: %w", err)
	}

	return nil
}
