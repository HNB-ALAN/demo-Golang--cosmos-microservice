package config

import (
	"os"
	"strconv"
)

// setDefaults sets default values for configuration
func setDefaults(cfg *Config, serviceName string) {
	// Service defaults
	cfg.Service.Name = getEnvOrDefault("SERVICE_NAME", serviceName)
	cfg.Service.Version = getEnvOrDefault("SERVICE_VERSION", "1.0.0")
	cfg.Service.Environment = getEnvOrDefault("ENVIRONMENT", "development")
	cfg.Service.Region = getEnvOrDefault("REGION", "us-east-1")
	cfg.Service.Instance = getEnvOrDefault("INSTANCE", "default")

	// Server defaults
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", "0.0.0.0")
	cfg.Server.Port = getEnvOrDefault("SERVER_PORT", "8080")
	cfg.Server.GRPC.MaxRecvMsgSize = getEnvIntOrDefault("GRPC_MAX_RECV_MSG_SIZE", 4*1024*1024) // 4MB
	cfg.Server.GRPC.MaxSendMsgSize = getEnvIntOrDefault("GRPC_MAX_SEND_MSG_SIZE", 4*1024*1024) // 4MB
	cfg.Server.GRPC.KeepAlive = getEnvBoolOrDefault("GRPC_KEEP_ALIVE", true)
	cfg.Server.HTTP.ReadTimeout = getEnvIntOrDefault("HTTP_READ_TIMEOUT", 30)
	cfg.Server.HTTP.WriteTimeout = getEnvIntOrDefault("HTTP_WRITE_TIMEOUT", 30)
	cfg.Server.HTTP.IdleTimeout = getEnvIntOrDefault("HTTP_IDLE_TIMEOUT", 120)

	// Database defaults
	cfg.Database.Host = getEnvOrDefault("POSTGRES_HOST", "localhost")
	cfg.Database.Port = getEnvIntOrDefault("POSTGRES_PORT", 5432)
	cfg.Database.User = getEnvOrDefault("POSTGRES_USER", "postgres")
	cfg.Database.Password = getEnvOrDefault("POSTGRES_PASSWORD", "password")
	cfg.Database.DBName = getEnvOrDefault("POSTGRES_DB", "usc_social_media")
	cfg.Database.SSLMode = getEnvOrDefault("POSTGRES_SSLMODE", "disable")
	cfg.Database.MaxConns = getEnvIntOrDefault("POSTGRES_MAX_CONNS", 25)
	cfg.Database.MinConns = getEnvIntOrDefault("POSTGRES_MIN_CONNS", 5)

	// Redis defaults
	cfg.Redis.Host = getEnvOrDefault("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnvIntOrDefault("REDIS_PORT", 6379)
	cfg.Redis.Password = getEnvOrDefault("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvIntOrDefault("REDIS_DB", 0)
	cfg.Redis.PoolSize = getEnvIntOrDefault("REDIS_POOL_SIZE", 10)

	// ClickHouse defaults
	cfg.ClickHouse.Host = getEnvOrDefault("CLICKHOUSE_HOST", "localhost")
	cfg.ClickHouse.Port = getEnvIntOrDefault("CLICKHOUSE_PORT", 9000)
	cfg.ClickHouse.User = getEnvOrDefault("CLICKHOUSE_USER", "default")
	cfg.ClickHouse.Password = getEnvOrDefault("CLICKHOUSE_PASSWORD", "")
	cfg.ClickHouse.DBName = getEnvOrDefault("CLICKHOUSE_DB", "usc_analytics")
	cfg.ClickHouse.Secure = getEnvBoolOrDefault("CLICKHOUSE_SECURE", false)

	// InfluxDB defaults
	cfg.InfluxDB.Host = getEnvOrDefault("INFLUXDB_HOST", "localhost")
	cfg.InfluxDB.Port = getEnvIntOrDefault("INFLUXDB_PORT", 8086)
	cfg.InfluxDB.Token = getEnvOrDefault("INFLUXDB_TOKEN", "")
	cfg.InfluxDB.Org = getEnvOrDefault("INFLUXDB_ORG", "usc")
	cfg.InfluxDB.Bucket = getEnvOrDefault("INFLUXDB_BUCKET", "metrics")

	// Quickwit defaults
	cfg.Quickwit.Host = getEnvOrDefault("QUICKWIT_HOST", "localhost")
	cfg.Quickwit.Port = getEnvIntOrDefault("QUICKWIT_PORT", 7280)
	cfg.Quickwit.User = getEnvOrDefault("QUICKWIT_USER", "")
	cfg.Quickwit.Password = getEnvOrDefault("QUICKWIT_PASSWORD", "")
	cfg.Quickwit.Index = getEnvOrDefault("QUICKWIT_INDEX", "usc_index")

	// Kafka defaults
	cfg.Kafka.Brokers = []string{getEnvOrDefault("KAFKA_BROKERS", "localhost:9092")}
	cfg.Kafka.ClientID = getEnvOrDefault("KAFKA_CLIENT_ID", serviceName)
	cfg.Kafka.GroupID = getEnvOrDefault("KAFKA_GROUP_ID", serviceName+"-group")
	cfg.Kafka.SecurityProtocol = getEnvOrDefault("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT")
	cfg.Kafka.SASLMechanism = getEnvOrDefault("KAFKA_SASL_MECHANISM", "PLAIN")
	cfg.Kafka.SASLUsername = getEnvOrDefault("KAFKA_SASL_USERNAME", "")
	cfg.Kafka.SASLPassword = getEnvOrDefault("KAFKA_SASL_PASSWORD", "")
	cfg.Kafka.SSLCAFile = getEnvOrDefault("KAFKA_SSL_CA_FILE", "")
	cfg.Kafka.SSLCertFile = getEnvOrDefault("KAFKA_SSL_CERT_FILE", "")
	cfg.Kafka.SSLKeyFile = getEnvOrDefault("KAFKA_SSL_KEY_FILE", "")
	cfg.Kafka.SSLKeyPassword = getEnvOrDefault("KAFKA_SSL_KEY_PASSWORD", "")
	cfg.Kafka.SessionTimeout = getEnvIntOrDefault("KAFKA_SESSION_TIMEOUT", 30000)      // 30 seconds
	cfg.Kafka.HeartbeatInterval = getEnvIntOrDefault("KAFKA_HEARTBEAT_INTERVAL", 3000) // 3 seconds
	cfg.Kafka.MaxPollRecords = getEnvIntOrDefault("KAFKA_MAX_POLL_RECORDS", 500)
	cfg.Kafka.AutoOffsetReset = getEnvOrDefault("KAFKA_AUTO_OFFSET_RESET", "latest")
	cfg.Kafka.EnableAutoCommit = getEnvBoolOrDefault("KAFKA_ENABLE_AUTO_COMMIT", true)
	cfg.Kafka.CompressionType = getEnvOrDefault("KAFKA_COMPRESSION_TYPE", "snappy")
	cfg.Kafka.BatchSize = getEnvIntOrDefault("KAFKA_BATCH_SIZE", 16384) // 16KB
	cfg.Kafka.LingerMs = getEnvIntOrDefault("KAFKA_LINGER_MS", 5)       // 5ms
	cfg.Kafka.Retries = getEnvIntOrDefault("KAFKA_RETRIES", 3)
	cfg.Kafka.RequestTimeout = getEnvIntOrDefault("KAFKA_REQUEST_TIMEOUT", 30000) // 30 seconds

	// Logging defaults
	cfg.Log.Level = getEnvOrDefault("LOG_LEVEL", "info")
	cfg.Log.Format = getEnvOrDefault("LOG_FORMAT", "json")
	cfg.Log.Output = getEnvOrDefault("LOG_OUTPUT", "stdout")
	cfg.Log.Filename = getEnvOrDefault("LOG_FILENAME", "")
	cfg.Log.MaxSize = getEnvIntOrDefault("LOG_MAX_SIZE", 100) // MB
	cfg.Log.MaxAge = getEnvIntOrDefault("LOG_MAX_AGE", 30)    // days
	cfg.Log.Compress = getEnvBoolOrDefault("LOG_COMPRESS", true)

	// Metrics defaults
	cfg.Metrics.Enabled = getEnvBoolOrDefault("METRICS_ENABLED", true)
	cfg.Metrics.Port = getEnvIntOrDefault("METRICS_PORT", 9090)
	cfg.Metrics.Path = getEnvOrDefault("METRICS_PATH", "/metrics")

	// Auth defaults
	cfg.Auth.JWTSecret = getEnvOrDefault("JWT_SECRET", "usc-platform-jwt-secret-key-2024-default-very-long-secret-key-for-production-use-only")
	cfg.Auth.JWTExpiry = getEnvOrDefault("JWT_EXPIRY", "24h")
	cfg.Auth.RefreshExpiry = getEnvOrDefault("REFRESH_EXPIRY", "168h") // 7 days
	cfg.Auth.Issuer = getEnvOrDefault("JWT_ISSUER", "usc-platform")

	// Cache defaults
	cfg.Cache.TTL = getEnvIntOrDefault("CACHE_TTL", 3600) // 1 hour
	cfg.Cache.MaxSize = getEnvIntOrDefault("CACHE_MAX_SIZE", 1000)
	cfg.Cache.Enabled = getEnvBoolOrDefault("CACHE_ENABLED", true)
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets environment variable as int or returns default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBoolOrDefault gets environment variable as bool or returns default value
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
