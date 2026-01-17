// Package constants provides application constants for USC platform services.
package constants

// Configuration keys
const (
	// Server configuration
	ConfigServerHost    = "server.host"
	ConfigServerPort    = "server.port"
	ConfigServerTimeout = "server.timeout"
	ConfigServerTLS     = "server.tls"
	ConfigServerCert    = "server.cert"
	ConfigServerKey     = "server.key"

	// Database configuration
	ConfigDatabaseHost     = "database.host"
	ConfigDatabasePort     = "database.port"
	ConfigDatabaseName     = "database.name"
	ConfigDatabaseUser     = "database.user"
	ConfigDatabasePassword = "database.password"
	ConfigDatabaseSSL      = "database.ssl"
	ConfigDatabaseMaxConns = "database.max_connections"
	ConfigDatabaseTimeout  = "database.timeout"

	// Redis configuration
	ConfigRedisHost     = "redis.host"
	ConfigRedisPort     = "redis.port"
	ConfigRedisPassword = "redis.password"
	ConfigRedisDB       = "redis.db"
	ConfigRedisTimeout  = "redis.timeout"
	ConfigRedisPoolSize = "redis.pool_size"

	// ClickHouse configuration
	ConfigClickHouseHost     = "clickhouse.host"
	ConfigClickHousePort     = "clickhouse.port"
	ConfigClickHouseDatabase = "clickhouse.database"
	ConfigClickHouseUser     = "clickhouse.user"
	ConfigClickHousePassword = "clickhouse.password"
	ConfigClickHouseTimeout  = "clickhouse.timeout"

	// InfluxDB configuration
	ConfigInfluxHost     = "influxdb.host"
	ConfigInfluxPort     = "influxdb.port"
	ConfigInfluxDatabase = "influxdb.database"
	ConfigInfluxUser     = "influxdb.user"
	ConfigInfluxPassword = "influxdb.password"
	ConfigInfluxTimeout  = "influxdb.timeout"

	// Quickwit configuration
	ConfigQuickwitHost     = "quickwit.host"
	ConfigQuickwitPort     = "quickwit.port"
	ConfigQuickwitUser     = "quickwit.user"
	ConfigQuickwitPassword = "quickwit.password"
	ConfigQuickwitTimeout  = "quickwit.timeout"

	// Logging configuration
	ConfigLogLevel  = "log.level"
	ConfigLogFormat = "log.format"
	ConfigLogOutput = "log.output"
	ConfigLogFile   = "log.file"

	// JWT configuration
	ConfigJWTSecret     = "jwt.secret"
	ConfigJWTExpiration = "jwt.expiration"
	ConfigJWTIssuer     = "jwt.issuer"
	ConfigJWTAudience   = "jwt.audience"

	// gRPC configuration
	ConfigGRPCPort    = "grpc.port"
	ConfigGRPCTLS     = "grpc.tls"
	ConfigGRPCCert    = "grpc.cert"
	ConfigGRPCKey     = "grpc.key"
	ConfigGRPCTimeout = "grpc.timeout"
	ConfigGRPCMaxSize = "grpc.max_size"

	// Monitoring configuration
	ConfigPrometheusPort = "prometheus.port"
	ConfigPrometheusPath = "prometheus.path"
	ConfigGrafanaPort    = "grafana.port"
	ConfigGrafanaPath    = "grafana.path"

	// Cache configuration
	ConfigCacheTTL      = "cache.ttl"
	ConfigCacheMaxSize  = "cache.max_size"
	ConfigCacheStrategy = "cache.strategy"

	// Rate limiting configuration
	ConfigRateLimitEnabled = "rate_limit.enabled"
	ConfigRateLimitRPS     = "rate_limit.rps"
	ConfigRateLimitBurst   = "rate_limit.burst"

	// Circuit breaker configuration
	ConfigCircuitBreakerEnabled = "circuit_breaker.enabled"
	ConfigCircuitBreakerFailure = "circuit_breaker.failure_threshold"
	ConfigCircuitBreakerTimeout = "circuit_breaker.timeout"
	ConfigCircuitBreakerReset   = "circuit_breaker.reset_timeout"
)

// Default configuration values
const (
	// Server defaults
	DefaultServerHost    = "127.0.0.1" // Bind to localhost by default for security
	DefaultServerPort    = 8080
	DefaultServerTimeout = 30
	DefaultServerTLS     = true // Enable TLS by default for security

	// Database defaults
	DefaultDatabaseHost     = "localhost"
	DefaultDatabasePort     = 5432
	DefaultDatabaseName     = "usc_platform"
	DefaultDatabaseUser     = "postgres"
	DefaultDatabasePassword = ""   // Must be set via environment variable
	DefaultDatabaseSSL      = true // Enable SSL by default for security
	DefaultDatabaseMaxConns = 100
	DefaultDatabaseTimeout  = 30

	// Redis defaults
	DefaultRedisHost     = "localhost"
	DefaultRedisPort     = 6379
	DefaultRedisPassword = ""
	DefaultRedisDB       = 0
	DefaultRedisTimeout  = 5
	DefaultRedisPoolSize = 10

	// ClickHouse defaults
	DefaultClickHouseHost     = "localhost"
	DefaultClickHousePort     = 9000
	DefaultClickHouseDatabase = "usc_platform"
	DefaultClickHouseUser     = "default"
	DefaultClickHousePassword = ""
	DefaultClickHouseTimeout  = 30

	// InfluxDB defaults
	DefaultInfluxHost     = "localhost"
	DefaultInfluxPort     = 8086
	DefaultInfluxDatabase = "usc_platform"
	DefaultInfluxUser     = "admin"
	DefaultInfluxPassword = "" // Must be set via environment variable
	DefaultInfluxTimeout  = 30

	// Quickwit defaults
	DefaultQuickwitHost     = "localhost"
	DefaultQuickwitPort     = 7280
	DefaultQuickwitUser     = ""
	DefaultQuickwitPassword = "" // Must be set via environment variable
	DefaultQuickwitTimeout  = 30

	// Logging defaults
	DefaultLogFile = ""

	// JWT defaults
	DefaultJWTSecret     = "" // Must be set via environment variable
	DefaultJWTExpiration = 3600
	DefaultJWTAudience   = "usc-platform"

	// gRPC defaults
	DefaultGRPCTLS     = false
	DefaultGRPCTimeout = 30
	DefaultGRPCMaxSize = 4194304

	// Monitoring defaults
	DefaultPrometheusPort = 9091
	DefaultPrometheusPath = "/metrics"
	DefaultGrafanaPort    = 3000
	DefaultGrafanaPath    = "/"

	// Cache defaults
	DefaultCacheMaxSize  = 1000
	DefaultCacheStrategy = "lru"

	// Rate limiting defaults
	DefaultRateLimitEnabled = true
	DefaultRateLimitRPS     = 100
	DefaultRateLimitBurst   = 200

	// Circuit breaker defaults
	DefaultCircuitBreakerEnabled = true
	DefaultCircuitBreakerFailure = 5
	DefaultCircuitBreakerTimeout = 30
	DefaultCircuitBreakerReset   = 60
)

// Environment variables
const (
	EnvServerHost    = "SERVER_HOST"
	EnvServerPort    = "SERVER_PORT"
	EnvServerTimeout = "SERVER_TIMEOUT"
	EnvServerTLS     = "SERVER_TLS"
	EnvServerCert    = "SERVER_CERT"
	EnvServerKey     = "SERVER_KEY"

	EnvDatabaseHost     = "DATABASE_HOST"
	EnvDatabasePort     = "DATABASE_PORT"
	EnvDatabaseName     = "DATABASE_NAME"
	EnvDatabaseUser     = "DATABASE_USER"
	EnvDatabasePassword = "DATABASE_PASSWORD"
	EnvDatabaseSSL      = "DATABASE_SSL"
	EnvDatabaseMaxConns = "DATABASE_MAX_CONNECTIONS"
	EnvDatabaseTimeout  = "DATABASE_TIMEOUT"

	EnvRedisHost     = "REDIS_HOST"
	EnvRedisPort     = "REDIS_PORT"
	EnvRedisPassword = "REDIS_PASSWORD"
	EnvRedisDB       = "REDIS_DB"
	EnvRedisTimeout  = "REDIS_TIMEOUT"
	EnvRedisPoolSize = "REDIS_POOL_SIZE"

	EnvClickHouseHost     = "CLICKHOUSE_HOST"
	EnvClickHousePort     = "CLICKHOUSE_PORT"
	EnvClickHouseDatabase = "CLICKHOUSE_DATABASE"
	EnvClickHouseUser     = "CLICKHOUSE_USER"
	EnvClickHousePassword = "CLICKHOUSE_PASSWORD"
	EnvClickHouseTimeout  = "CLICKHOUSE_TIMEOUT"

	EnvInfluxHost     = "INFLUX_HOST"
	EnvInfluxPort     = "INFLUX_PORT"
	EnvInfluxDatabase = "INFLUX_DATABASE"
	EnvInfluxUser     = "INFLUX_USER"
	EnvInfluxPassword = "INFLUX_PASSWORD"
	EnvInfluxTimeout  = "INFLUX_TIMEOUT"

	EnvQuickwitHost     = "QUICKWIT_HOST"
	EnvQuickwitPort     = "QUICKWIT_PORT"
	EnvQuickwitUser     = "QUICKWIT_USER"
	EnvQuickwitPassword = "QUICKWIT_PASSWORD"
	EnvQuickwitTimeout  = "QUICKWIT_TIMEOUT"

	EnvLogLevel  = "LOG_LEVEL"
	EnvLogFormat = "LOG_FORMAT"
	EnvLogOutput = "LOG_OUTPUT"
	EnvLogFile   = "LOG_FILE"

	EnvJWTSecret     = "JWT_SECRET"
	EnvJWTExpiration = "JWT_EXPIRATION"
	EnvJWTIssuer     = "JWT_ISSUER"
	EnvJWTAudience   = "JWT_AUDIENCE"

	EnvGRPCPort    = "GRPC_PORT"
	EnvGRPCTLS     = "GRPC_TLS"
	EnvGRPCCert    = "GRPC_CERT"
	EnvGRPCKey     = "GRPC_KEY"
	EnvGRPCTimeout = "GRPC_TIMEOUT"
	EnvGRPCMaxSize = "GRPC_MAX_SIZE"

	EnvPrometheusPort = "PROMETHEUS_PORT"
	EnvPrometheusPath = "PROMETHEUS_PATH"
	EnvGrafanaPort    = "GRAFANA_PORT"
	EnvGrafanaPath    = "GRAFANA_PATH"

	EnvCacheTTL      = "CACHE_TTL"
	EnvCacheMaxSize  = "CACHE_MAX_SIZE"
	EnvCacheStrategy = "CACHE_STRATEGY"

	EnvRateLimitEnabled = "RATE_LIMIT_ENABLED"
	EnvRateLimitRPS     = "RATE_LIMIT_RPS"
	EnvRateLimitBurst   = "RATE_LIMIT_BURST"

	EnvCircuitBreakerEnabled = "CIRCUIT_BREAKER_ENABLED"
	EnvCircuitBreakerFailure = "CIRCUIT_BREAKER_FAILURE_THRESHOLD"
	EnvCircuitBreakerTimeout = "CIRCUIT_BREAKER_TIMEOUT"
	EnvCircuitBreakerReset   = "CIRCUIT_BREAKER_RESET_TIMEOUT"
)
