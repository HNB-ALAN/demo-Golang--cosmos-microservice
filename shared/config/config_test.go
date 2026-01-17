package config

import (
	"testing"
)

func TestConfig_GetServerAddress(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	address := config.GetServerAddress()
	expected := "localhost:8080"

	if address != expected {
		t.Errorf("Expected server address %s, got %s", expected, address)
	}
}

func TestConfig_GetDatabaseDSN(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	dsn := config.GetDatabaseDSN()
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	if dsn != expected {
		t.Errorf("Expected DSN %s, got %s", expected, dsn)
	}
}

func TestConfig_GetRedisAddress(t *testing.T) {
	config := &Config{
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
	}

	address := config.GetRedisAddress()
	expected := "localhost:6379"

	if address != expected {
		t.Errorf("Expected Redis address %s, got %s", expected, address)
	}
}

func TestConfig_GetClickHouseDSN(t *testing.T) {
	config := &Config{
		ClickHouse: ClickHouseConfig{
			Host:     "localhost",
			Port:     9000,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
		},
	}

	dsn := config.GetClickHouseDSN()
	expected := "http://localhost:9000?username=testuser&password=testpass&database=testdb"

	if dsn != expected {
		t.Errorf("Expected ClickHouse DSN %s, got %s", expected, dsn)
	}
}

func TestConfig_GetInfluxDBURL(t *testing.T) {
	config := &Config{
		InfluxDB: InfluxDBConfig{
			Host: "localhost",
			Port: 8086,
		},
	}

	url := config.GetInfluxDBURL()
	expected := "http://localhost:8086"

	if url != expected {
		t.Errorf("Expected InfluxDB URL %s, got %s", expected, url)
	}
}

func TestConfig_GetQuickwitURL(t *testing.T) {
	config := &Config{
		Quickwit: QuickwitConfig{
			Host: "localhost",
			Port: 7280,
		},
	}

	url := config.GetQuickwitURL()
	expected := "http://localhost:7280"

	if url != expected {
		t.Errorf("Expected Quickwit URL %s, got %s", expected, url)
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	// Test development environment
	config := &Config{
		Service: ServiceConfig{
			Environment: "development",
		},
	}

	if !config.IsDevelopment() {
		t.Error("Expected config to be development environment")
	}

	// Test non-development environment
	config.Service.Environment = "production"
	if config.IsDevelopment() {
		t.Error("Expected config to not be development environment")
	}
}

func TestConfig_IsProduction(t *testing.T) {
	// Test production environment
	config := &Config{
		Service: ServiceConfig{
			Environment: "production",
		},
	}

	if !config.IsProduction() {
		t.Error("Expected config to be production environment")
	}

	// Test non-production environment
	config.Service.Environment = "development"
	if config.IsProduction() {
		t.Error("Expected config to not be production environment")
	}
}

// Benchmark tests
func BenchmarkConfig_GetServerAddress(b *testing.B) {
	config := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.GetServerAddress()
	}
}

func BenchmarkConfig_GetDatabaseDSN(b *testing.B) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.GetDatabaseDSN()
	}
}
