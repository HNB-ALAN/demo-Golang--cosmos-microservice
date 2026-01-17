package database

import (
	"context"
	"testing"
	"time"
)

func TestPoolManager_NewPoolManager(t *testing.T) {
	config := PoolConfig{
		MinSize:             5,
		MaxSize:             20,
		InitialSize:         10,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)

	if manager == nil {
		t.Fatal("Expected pool manager to be created")
	}

	if manager.minSize != config.MinSize {
		t.Errorf("Expected MinSize %d, got %d", config.MinSize, manager.minSize)
	}

	if manager.maxSize != config.MaxSize {
		t.Errorf("Expected MaxSize %d, got %d", config.MaxSize, manager.maxSize)
	}

	if manager.currentSize != config.InitialSize {
		t.Errorf("Expected InitialSize %d, got %d", config.InitialSize, manager.currentSize)
	}
}

func TestPoolManager_GetConnection(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Test getting a connection
	conn, err := manager.GetConnection(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if conn == nil {
		t.Fatal("Expected connection to be returned")
	}

	if !conn.IsActive() {
		t.Error("Expected connection to be active")
	}

	// Verify pool stats
	stats := manager.GetStats()
	if stats.ActiveConnections != 1 {
		t.Errorf("Expected 1 active connection, got %d", stats.ActiveConnections)
	}

	if stats.TotalRequests != 1 {
		t.Errorf("Expected 1 total request, got %d", stats.TotalRequests)
	}
}

func TestPoolManager_ReturnConnection(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Get a connection
	conn, err := manager.GetConnection(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Return the connection
	manager.ReturnConnection(conn)

	// Verify connection is no longer active
	if conn.IsActive() {
		t.Error("Expected connection to be inactive after return")
	}

	// Verify pool stats
	stats := manager.GetStats()
	if stats.ActiveConnections != 0 {
		t.Errorf("Expected 0 active connections, got %d", stats.ActiveConnections)
	}

	if stats.IdleConnections != 1 {
		t.Errorf("Expected 1 idle connection, got %d", stats.IdleConnections)
	}
}

func TestPoolManager_PoolExhaustion(t *testing.T) {
	config := PoolConfig{
		MinSize:             1,
		MaxSize:             2,
		InitialSize:         2,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Exhaust the pool
	conn1, err := manager.GetConnection(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	conn2, err := manager.GetConnection(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Try to get another connection (should fail)
	_, err = manager.GetConnection(ctx)
	if err == nil {
		t.Error("Expected error when pool is exhausted, got nil")
	}

	// Verify stats
	stats := manager.GetStats()
	if stats.ActiveConnections != 2 {
		t.Errorf("Expected 2 active connections, got %d", stats.ActiveConnections)
	}

	if stats.FailedRequests != 1 {
		t.Errorf("Expected 1 failed request, got %d", stats.FailedRequests)
	}

	// Return connections
	manager.ReturnConnection(conn1)
	manager.ReturnConnection(conn2)
}

func TestPoolManager_LoadBasedScaling(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         3,
		AdjustmentInterval:  100 * time.Millisecond, // Short interval for testing
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Simulate high load by getting many connections
	connections := make([]*PooledConnection, 0)
	for i := 0; i < 3; i++ { // Only get 3 connections to avoid pool exhaustion
		conn, err := manager.GetConnection(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		connections = append(connections, conn)
	}

	// Wait for adjustment interval
	time.Sleep(200 * time.Millisecond)

	// Trigger another adjustment
	manager.adjustPoolSize()

	// Verify pool size increased
	stats := manager.GetStats()
	if stats.CurrentSize <= config.InitialSize {
		t.Errorf("Expected pool size to increase, got %d", stats.CurrentSize)
	}

	// Return connections
	for _, conn := range connections {
		manager.ReturnConnection(conn)
	}
}

func TestPoolManager_HealthCheck(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Test healthy pool
	err := manager.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Simulate overloaded pool
	manager.loadFactor = 0.96
	err = manager.HealthCheck(ctx)
	if err == nil {
		t.Error("Expected error for overloaded pool, got nil")
	}

	// Reset load factor
	manager.loadFactor = 0.5

	// Simulate low success rate
	manager.failedRequests = 10
	manager.successfulRequests = 5
	err = manager.HealthCheck(ctx)
	if err == nil {
		t.Error("Expected error for low success rate, got nil")
	}
}

func TestPoolManager_GetStats(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Get some connections
	conn1, _ := manager.GetConnection(ctx)
	conn2, _ := manager.GetConnection(ctx)

	// Return one connection
	manager.ReturnConnection(conn1)

	stats := manager.GetStats()

	// Verify stats
	if stats.MinSize != config.MinSize {
		t.Errorf("Expected MinSize %d, got %d", config.MinSize, stats.MinSize)
	}

	if stats.MaxSize != config.MaxSize {
		t.Errorf("Expected MaxSize %d, got %d", config.MaxSize, stats.MaxSize)
	}

	if stats.ActiveConnections != 1 {
		t.Errorf("Expected 1 active connection, got %d", stats.ActiveConnections)
	}

	if stats.IdleConnections != 1 {
		t.Errorf("Expected 1 idle connection, got %d", stats.IdleConnections)
	}

	if stats.TotalRequests != 2 {
		t.Errorf("Expected 2 total requests, got %d", stats.TotalRequests)
	}

	if stats.SuccessRate != 100.0 {
		t.Errorf("Expected 100%% success rate, got %.2f", stats.SuccessRate)
	}

	// Return remaining connection
	manager.ReturnConnection(conn2)
}

func TestPoolManager_ConcurrentAccess(t *testing.T) {
	config := PoolConfig{
		MinSize:             5,
		MaxSize:             20,
		InitialSize:         10,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	// Test concurrent access
	done := make(chan bool, 10) // Reduce to 10 to avoid pool exhaustion
	for i := 0; i < 10; i++ {
		go func(i int) {
			conn, err := manager.GetConnection(ctx)
			if err != nil {
				// Pool exhaustion is expected in concurrent tests
				done <- false
				return
			}

			// Simulate some work
			time.Sleep(10 * time.Millisecond)

			manager.ReturnConnection(conn)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	successCount := 0
	for i := 0; i < 10; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount < 5 { // Allow some failures due to pool exhaustion
		t.Errorf("Expected at least 5 successful operations, got %d", successCount)
	}
}

func TestMultiPoolManager_NewMultiPoolManager(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)

	if manager == nil {
		t.Fatal("Expected multi-pool manager to be created")
	}

	if len(manager.pools) != 0 {
		t.Errorf("Expected 0 pools initially, got %d", len(manager.pools))
	}
}

func TestMultiPoolManager_GetPool(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)

	// Get a pool
	pool1 := manager.GetPool("database1")
	if pool1 == nil {
		t.Fatal("Expected pool to be created")
	}

	// Get the same pool again
	pool2 := manager.GetPool("database1")
	if pool1 != pool2 {
		t.Error("Expected same pool instance to be returned")
	}

	// Get a different pool
	pool3 := manager.GetPool("database2")
	if pool3 == nil {
		t.Fatal("Expected second pool to be created")
	}

	if pool1 == pool3 {
		t.Error("Expected different pool instances")
	}
}

func TestMultiPoolManager_GetConnection(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)
	ctx := context.Background()

	// Get connection from specific database
	conn, err := manager.GetConnection(ctx, "database1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if conn == nil {
		t.Fatal("Expected connection to be returned")
	}

	// Return connection
	conn.Close()
}

func TestMultiPoolManager_GetStats(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)
	ctx := context.Background()

	// Create some pools and get connections
	conn1, _ := manager.GetConnection(ctx, "database1")
	conn2, _ := manager.GetConnection(ctx, "database2")

	stats := manager.GetStats()

	if len(stats) != 2 {
		t.Errorf("Expected 2 pool stats, got %d", len(stats))
	}

	// Verify each pool has stats
	if _, exists := stats["database1"]; !exists {
		t.Error("Expected database1 stats to exist")
	}

	if _, exists := stats["database2"]; !exists {
		t.Error("Expected database2 stats to exist")
	}

	// Return connections
	conn1.Close()
	conn2.Close()
}

func TestMultiPoolManager_HealthCheck(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)
	ctx := context.Background()

	// Create some pools
	manager.GetPool("database1")
	manager.GetPool("database2")

	// Test health check
	err := manager.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMultiPoolManager_Close(t *testing.T) {
	config := PoolConfig{
		MinSize:             2,
		MaxSize:             10,
		InitialSize:         5,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)
	ctx := context.Background()

	// Create some pools
	manager.GetConnection(ctx, "database1")
	manager.GetConnection(ctx, "database2")

	// Close manager
	err := manager.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify pools are cleared
	stats := manager.GetStats()
	if len(stats) != 0 {
		t.Errorf("Expected 0 pools after close, got %d", len(stats))
	}
}

// Benchmark tests
func BenchmarkPoolManager_GetConnection(b *testing.B) {
	config := PoolConfig{
		MinSize:             10,
		MaxSize:             100,
		InitialSize:         50,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := manager.GetConnection(ctx)
		if err == nil {
			manager.ReturnConnection(conn)
		}
	}
}

func BenchmarkPoolManager_Concurrent(b *testing.B) {
	config := PoolConfig{
		MinSize:             10,
		MaxSize:             100,
		InitialSize:         50,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewPoolManager(config)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := manager.GetConnection(ctx)
			if err == nil {
				manager.ReturnConnection(conn)
			}
		}
	})
}

func BenchmarkMultiPoolManager_GetConnection(b *testing.B) {
	config := PoolConfig{
		MinSize:             10,
		MaxSize:             100,
		InitialSize:         50,
		AdjustmentInterval:  1 * time.Minute,
		LoadThreshold:       0.8,
		HealthCheckInterval: 30 * time.Second,
	}

	manager := NewMultiPoolManager(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := manager.GetConnection(ctx, "database1")
		if err == nil {
			conn.Close()
		}
	}
}
