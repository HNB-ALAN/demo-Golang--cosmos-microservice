package config

import (
	"fmt"
	"strconv"
	"strings"
)

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if err := validateServerConfig(&cfg.Server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := validateDatabaseConfig(&cfg.Database); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := validateRedisConfig(&cfg.Redis); err != nil {
		return fmt.Errorf("redis config validation failed: %w", err)
	}

	if err := validateClickHouseConfig(&cfg.ClickHouse); err != nil {
		return fmt.Errorf("clickhouse config validation failed: %w", err)
	}

	if err := validateInfluxDBConfig(&cfg.InfluxDB); err != nil {
		return fmt.Errorf("influxdb config validation failed: %w", err)
	}

	if err := validateQuickwitConfig(&cfg.Quickwit); err != nil {
		return fmt.Errorf("quickwit config validation failed: %w", err)
	}

	if err := validateLogConfig(&cfg.Log); err != nil {
		return fmt.Errorf("log config validation failed: %w", err)
	}

	if err := validateServiceConfig(&cfg.Service); err != nil {
		return fmt.Errorf("service config validation failed: %w", err)
	}

	if err := validateMetricsConfig(&cfg.Metrics); err != nil {
		return fmt.Errorf("metrics config validation failed: %w", err)
	}

	if err := validateAuthConfig(&cfg.Auth); err != nil {
		return fmt.Errorf("auth config validation failed: %w", err)
	}

	if err := validateCacheConfig(&cfg.Cache); err != nil {
		return fmt.Errorf("cache config validation failed: %w", err)
	}

	return nil
}

// validateServerConfig validates server configuration
func validateServerConfig(cfg *ServerConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}

	if cfg.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return fmt.Errorf("invalid server port: %s", cfg.Port)
	}

	if cfg.GRPC.MaxRecvMsgSize <= 0 {
		return fmt.Errorf("grpc max_recv_msg_size must be positive")
	}

	if cfg.GRPC.MaxSendMsgSize <= 0 {
		return fmt.Errorf("grpc max_send_msg_size must be positive")
	}

	if cfg.HTTP.ReadTimeout <= 0 {
		return fmt.Errorf("http read_timeout must be positive")
	}

	if cfg.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("http write_timeout must be positive")
	}

	if cfg.HTTP.IdleTimeout <= 0 {
		return fmt.Errorf("http idle_timeout must be positive")
	}

	return nil
}

// validateDatabaseConfig validates database configuration
func validateDatabaseConfig(cfg *DatabaseConfig) error {
	// Only validate if database is enabled
	if !cfg.Enabled {
		return nil
	}

	if cfg.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.Port)
	}

	if cfg.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}

	if cfg.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
	if !contains(validSSLModes, cfg.SSLMode) {
		return fmt.Errorf("invalid ssl_mode: %s", cfg.SSLMode)
	}

	if cfg.MaxConns <= 0 {
		return fmt.Errorf("max_conns must be positive")
	}

	if cfg.MinConns < 0 {
		return fmt.Errorf("min_conns cannot be negative")
	}

	if cfg.MinConns > cfg.MaxConns {
		return fmt.Errorf("min_conns cannot be greater than max_conns")
	}

	return nil
}

// validateRedisConfig validates Redis configuration
func validateRedisConfig(cfg *RedisConfig) error {
	// Only validate if Redis is enabled
	if !cfg.Enabled {
		return nil
	}

	if cfg.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", cfg.Port)
	}

	if cfg.DB < 0 || cfg.DB > 15 {
		return fmt.Errorf("invalid redis db: %d", cfg.DB)
	}

	if cfg.PoolSize <= 0 {
		return fmt.Errorf("redis pool_size must be positive")
	}

	return nil
}

// validateClickHouseConfig validates ClickHouse configuration
func validateClickHouseConfig(cfg *ClickHouseConfig) error {
	// Only validate if ClickHouse is enabled
	if !cfg.Enabled {
		return nil
	}

	if cfg.Host == "" {
		return fmt.Errorf("clickhouse host cannot be empty")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid clickhouse port: %d", cfg.Port)
	}

	if cfg.User == "" {
		return fmt.Errorf("clickhouse user cannot be empty")
	}

	if cfg.DBName == "" {
		return fmt.Errorf("clickhouse database name cannot be empty")
	}

	return nil
}

// validateInfluxDBConfig validates InfluxDB configuration
func validateInfluxDBConfig(cfg *InfluxDBConfig) error {
	// Only validate if InfluxDB is enabled
	if !cfg.Enabled {
		return nil
	}

	if cfg.Host == "" {
		return fmt.Errorf("influxdb host cannot be empty")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid influxdb port: %d", cfg.Port)
	}

	if cfg.Token == "" {
		return fmt.Errorf("influxdb token cannot be empty")
	}

	if cfg.Org == "" {
		return fmt.Errorf("influxdb org cannot be empty")
	}

	if cfg.Bucket == "" {
		return fmt.Errorf("influxdb bucket cannot be empty")
	}

	return nil
}

// validateQuickwitConfig validates Quickwit configuration
func validateQuickwitConfig(cfg *QuickwitConfig) error {
	// Only validate if Quickwit is enabled
	if !cfg.Enabled {
		return nil
	}

	if cfg.Host == "" {
		return fmt.Errorf("quickwit host cannot be empty")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid quickwit port: %d", cfg.Port)
	}

	if cfg.Index == "" {
		return fmt.Errorf("quickwit index cannot be empty")
	}

	return nil
}

// validateLogConfig validates logging configuration
func validateLogConfig(cfg *LogConfig) error {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	if !contains(validLevels, cfg.Level) {
		return fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	validFormats := []string{"json", "text", "console"}
	if !contains(validFormats, cfg.Format) {
		return fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	validOutputs := []string{"stdout", "stderr", "file"}
	if !contains(validOutputs, cfg.Output) {
		return fmt.Errorf("invalid log output: %s", cfg.Output)
	}

	if cfg.Output == "file" && cfg.Filename == "" {
		return fmt.Errorf("log filename is required when output is file")
	}

	if cfg.MaxSize <= 0 {
		return fmt.Errorf("log max_size must be positive")
	}

	if cfg.MaxAge <= 0 {
		return fmt.Errorf("log max_age must be positive")
	}

	return nil
}

// validateServiceConfig validates service configuration
func validateServiceConfig(cfg *ServiceConfig) error {
	if cfg.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if cfg.Version == "" {
		return fmt.Errorf("service version cannot be empty")
	}

	validEnvironments := []string{"development", "staging", "production"}
	if !contains(validEnvironments, cfg.Environment) {
		return fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	return nil
}

// validateMetricsConfig validates metrics configuration
func validateMetricsConfig(cfg *MetricsConfig) error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid metrics port: %d", cfg.Port)
	}

	if cfg.Path == "" {
		return fmt.Errorf("metrics path cannot be empty")
	}

	if !strings.HasPrefix(cfg.Path, "/") {
		return fmt.Errorf("metrics path must start with /")
	}

	return nil
}

// validateAuthConfig validates authentication configuration
func validateAuthConfig(cfg *AuthConfig) error {
	if cfg.JWTSecret == "" {
		return fmt.Errorf("jwt secret cannot be empty")
	}

	if len(cfg.JWTSecret) < 32 {
		return fmt.Errorf("jwt secret must be at least 32 characters")
	}

	if cfg.JWTExpiry == "" {
		return fmt.Errorf("jwt expiry cannot be empty")
	}

	if cfg.RefreshExpiry == "" {
		return fmt.Errorf("refresh expiry cannot be empty")
	}

	if cfg.Issuer == "" {
		return fmt.Errorf("jwt issuer cannot be empty")
	}

	return nil
}

// validateCacheConfig validates cache configuration
func validateCacheConfig(cfg *CacheConfig) error {
	if cfg.TTL <= 0 {
		return fmt.Errorf("cache ttl must be positive")
	}

	if cfg.MaxSize <= 0 {
		return fmt.Errorf("cache max_size must be positive")
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
