package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
)

func TestEnhancedGRPCClient(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
		Log: config.LogConfig{
			Level:  "debug",
			Format: "console",
		},
	}

	// Create logger
	logger := logging.NewLogger("test", cfg.Log)

	// Create enhanced client config
	clientConfig := DefaultEnhancedClientConfig()
	clientConfig.Addresses = []string{"localhost:9090", "localhost:9091"}
	clientConfig.EnableMetrics = true
	clientConfig.EnableRetry = true
	clientConfig.EnableCircuitBreaker = true
	clientConfig.EnableLoadBalancer = true

	// Create enhanced client
	client := NewEnhancedGRPCClient(cfg, logger, clientConfig)
	defer client.Close()

	// Test creating a client
	conn, err := client.CreateEnhancedClient("test-client", "localhost:9090")
	if err != nil {
		t.Logf("Expected error creating client (no server running): %v", err)
	} else {
		t.Log("Client created successfully")
		conn.Close()
	}

	// Test metrics
	metrics := client.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be returned")
	} else {
		t.Logf("Metrics: %+v", metrics)
	}

	// Test health check
	ctx := context.Background()
	err = client.HealthCheck(ctx)
	if err != nil {
		t.Logf("Expected health check error (no server running): %v", err)
	}
}

func TestCircuitBreaker(t *testing.T) {
	config := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(config)

	// Test initial state
	if cb.GetState() != StateClosed {
		t.Errorf("Expected initial state to be Closed, got %v", cb.GetState())
	}

	// Test executing a function that fails
	err := cb.Execute(func() error {
		return fmt.Errorf("test error")
	})

	if err == nil {
		t.Error("Expected error to be returned")
	}

	// Test circuit breaker state after failures
	state := cb.GetState()
	t.Logf("Circuit breaker state after failure: %v", state)
}

func TestLoadBalancer(t *testing.T) {
	config := DefaultLoadBalancerConfig()
	addresses := []string{"localhost:9090", "localhost:9091", "localhost:9092"}
	lb := NewLoadBalancer(config, addresses)

	// Test getting next address
	address := lb.GetNextAddress()
	if address == "" {
		t.Error("Expected address to be returned")
	} else {
		t.Logf("Next address: %s", address)
	}

	// Test recording failure (need multiple failures to mark as unhealthy)
	for i := 0; i < config.MaxFailures; i++ {
		lb.RecordFailure(address)
	}
	if lb.IsHealthy(address) {
		t.Error("Expected address to be unhealthy after multiple failures")
	}

	// Test recording success
	lb.RecordSuccess(address)
	if !lb.IsHealthy(address) {
		t.Error("Expected address to be healthy after success")
	}
}

func TestRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()
	if config.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts to be 3, got %d", config.MaxAttempts)
	}

	if config.InitialDelay != 100*time.Millisecond {
		t.Errorf("Expected InitialDelay to be 100ms, got %v", config.InitialDelay)
	}
}

func TestConnectionPoolConfig(t *testing.T) {
	config := DefaultConnectionPoolConfig()
	if config.MaxConnections != 10 {
		t.Errorf("Expected MaxConnections to be 10, got %d", config.MaxConnections)
	}

	if config.MinConnections != 2 {
		t.Errorf("Expected MinConnections to be 2, got %d", config.MinConnections)
	}
}
