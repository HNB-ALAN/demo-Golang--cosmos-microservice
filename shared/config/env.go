package config

import (
	"os"
	"strconv"
	"strings"
)

// EnvManager handles environment variable operations
type EnvManager struct {
	prefix string
}

// NewEnvManager creates a new environment manager
func NewEnvManager(prefix string) *EnvManager {
	return &EnvManager{
		prefix: prefix,
	}
}

// GetString gets a string environment variable
func (em *EnvManager) GetString(key string, defaultValue string) string {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		return value
	}
	return defaultValue
}

// GetInt gets an integer environment variable
func (em *EnvManager) GetInt(key string, defaultValue int) int {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetBool gets a boolean environment variable
func (em *EnvManager) GetBool(key string, defaultValue bool) bool {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetFloat64 gets a float64 environment variable
func (em *EnvManager) GetFloat64(key string, defaultValue float64) float64 {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// GetStringSlice gets a string slice environment variable
func (em *EnvManager) GetStringSlice(key string, defaultValue []string) []string {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// GetIntSlice gets an integer slice environment variable
func (em *EnvManager) GetIntSlice(key string, defaultValue []int) []int {
	fullKey := em.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		parts := strings.Split(value, ",")
		result := make([]int, 0, len(parts))
		for _, part := range parts {
			if intValue, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				result = append(result, intValue)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// Set sets an environment variable
func (em *EnvManager) Set(key, value string) error {
	fullKey := em.getFullKey(key)
	return os.Setenv(fullKey, value)
}

// Unset unsets an environment variable
func (em *EnvManager) Unset(key string) error {
	fullKey := em.getFullKey(key)
	return os.Unsetenv(fullKey)
}

// Exists checks if an environment variable exists
func (em *EnvManager) Exists(key string) bool {
	fullKey := em.getFullKey(key)
	_, exists := os.LookupEnv(fullKey)
	return exists
}

// GetAll gets all environment variables with the prefix
func (em *EnvManager) GetAll() map[string]string {
	result := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], em.prefix) {
			key := strings.TrimPrefix(parts[0], em.prefix+"_")
			result[key] = parts[1]
		}
	}
	return result
}

// getFullKey returns the full environment variable key with prefix
func (em *EnvManager) getFullKey(key string) string {
	if em.prefix == "" {
		return key
	}
	return em.prefix + "_" + strings.ToUpper(key)
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv(cfg *Config, prefix string) error {
	em := NewEnvManager(prefix)

	// Load server config
	cfg.Server.Host = em.GetString("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = em.GetString("SERVER_PORT", cfg.Server.Port)
	cfg.Server.GRPC.MaxRecvMsgSize = em.GetInt("GRPC_MAX_RECV_MSG_SIZE", cfg.Server.GRPC.MaxRecvMsgSize)
	cfg.Server.GRPC.MaxSendMsgSize = em.GetInt("GRPC_MAX_SEND_MSG_SIZE", cfg.Server.GRPC.MaxSendMsgSize)
	cfg.Server.GRPC.KeepAlive = em.GetBool("GRPC_KEEP_ALIVE", cfg.Server.GRPC.KeepAlive)
	cfg.Server.HTTP.ReadTimeout = em.GetInt("HTTP_READ_TIMEOUT", cfg.Server.HTTP.ReadTimeout)
	cfg.Server.HTTP.WriteTimeout = em.GetInt("HTTP_WRITE_TIMEOUT", cfg.Server.HTTP.WriteTimeout)
	cfg.Server.HTTP.IdleTimeout = em.GetInt("HTTP_IDLE_TIMEOUT", cfg.Server.HTTP.IdleTimeout)

	// Load database config
	cfg.Database.Host = em.GetString("POSTGRES_HOST", cfg.Database.Host)
	cfg.Database.Port = em.GetInt("POSTGRES_PORT", cfg.Database.Port)
	cfg.Database.User = em.GetString("POSTGRES_USER", cfg.Database.User)
	cfg.Database.Password = em.GetString("POSTGRES_PASSWORD", cfg.Database.Password)
	cfg.Database.DBName = em.GetString("POSTGRES_DB", cfg.Database.DBName)
	cfg.Database.SSLMode = em.GetString("POSTGRES_SSLMODE", cfg.Database.SSLMode)
	cfg.Database.MaxConns = em.GetInt("POSTGRES_MAX_CONNS", cfg.Database.MaxConns)
	cfg.Database.MinConns = em.GetInt("POSTGRES_MIN_CONNS", cfg.Database.MinConns)

	// Load Redis config
	cfg.Redis.Host = em.GetString("REDIS_HOST", cfg.Redis.Host)
	cfg.Redis.Port = em.GetInt("REDIS_PORT", cfg.Redis.Port)
	cfg.Redis.Password = em.GetString("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = em.GetInt("REDIS_DB", cfg.Redis.DB)
	cfg.Redis.PoolSize = em.GetInt("REDIS_POOL_SIZE", cfg.Redis.PoolSize)

	// Load ClickHouse config
	cfg.ClickHouse.Host = em.GetString("CLICKHOUSE_HOST", cfg.ClickHouse.Host)
	cfg.ClickHouse.Port = em.GetInt("CLICKHOUSE_PORT", cfg.ClickHouse.Port)
	cfg.ClickHouse.User = em.GetString("CLICKHOUSE_USER", cfg.ClickHouse.User)
	cfg.ClickHouse.Password = em.GetString("CLICKHOUSE_PASSWORD", cfg.ClickHouse.Password)
	cfg.ClickHouse.DBName = em.GetString("CLICKHOUSE_DB", cfg.ClickHouse.DBName)
	cfg.ClickHouse.Secure = em.GetBool("CLICKHOUSE_SECURE", cfg.ClickHouse.Secure)

	// Load InfluxDB config
	cfg.InfluxDB.Host = em.GetString("INFLUXDB_HOST", cfg.InfluxDB.Host)
	cfg.InfluxDB.Port = em.GetInt("INFLUXDB_PORT", cfg.InfluxDB.Port)
	cfg.InfluxDB.Token = em.GetString("INFLUXDB_TOKEN", cfg.InfluxDB.Token)
	cfg.InfluxDB.Org = em.GetString("INFLUXDB_ORG", cfg.InfluxDB.Org)
	cfg.InfluxDB.Bucket = em.GetString("INFLUXDB_BUCKET", cfg.InfluxDB.Bucket)

	// Load Quickwit config
	cfg.Quickwit.Host = em.GetString("QUICKWIT_HOST", cfg.Quickwit.Host)
	cfg.Quickwit.Port = em.GetInt("QUICKWIT_PORT", cfg.Quickwit.Port)
	cfg.Quickwit.User = em.GetString("QUICKWIT_USER", cfg.Quickwit.User)
	cfg.Quickwit.Password = em.GetString("QUICKWIT_PASSWORD", cfg.Quickwit.Password)
	cfg.Quickwit.Index = em.GetString("QUICKWIT_INDEX", cfg.Quickwit.Index)

	// Load logging config
	cfg.Log.Level = em.GetString("LOG_LEVEL", cfg.Log.Level)
	cfg.Log.Format = em.GetString("LOG_FORMAT", cfg.Log.Format)
	cfg.Log.Output = em.GetString("LOG_OUTPUT", cfg.Log.Output)
	cfg.Log.Filename = em.GetString("LOG_FILENAME", cfg.Log.Filename)
	cfg.Log.MaxSize = em.GetInt("LOG_MAX_SIZE", cfg.Log.MaxSize)
	cfg.Log.MaxAge = em.GetInt("LOG_MAX_AGE", cfg.Log.MaxAge)
	cfg.Log.Compress = em.GetBool("LOG_COMPRESS", cfg.Log.Compress)

	// Load service config
	cfg.Service.Name = em.GetString("SERVICE_NAME", cfg.Service.Name)
	cfg.Service.Version = em.GetString("SERVICE_VERSION", cfg.Service.Version)
	cfg.Service.Environment = em.GetString("ENVIRONMENT", cfg.Service.Environment)
	cfg.Service.Region = em.GetString("REGION", cfg.Service.Region)
	cfg.Service.Instance = em.GetString("INSTANCE", cfg.Service.Instance)

	// Load metrics config
	cfg.Metrics.Enabled = em.GetBool("METRICS_ENABLED", cfg.Metrics.Enabled)
	cfg.Metrics.Port = em.GetInt("METRICS_PORT", cfg.Metrics.Port)
	cfg.Metrics.Path = em.GetString("METRICS_PATH", cfg.Metrics.Path)

	// Load auth config
	cfg.Auth.JWTSecret = em.GetString("JWT_SECRET", cfg.Auth.JWTSecret)
	cfg.Auth.JWTExpiry = em.GetString("JWT_EXPIRY", cfg.Auth.JWTExpiry)
	cfg.Auth.RefreshExpiry = em.GetString("REFRESH_EXPIRY", cfg.Auth.RefreshExpiry)
	cfg.Auth.Issuer = em.GetString("JWT_ISSUER", cfg.Auth.Issuer)

	// Load cache config
	cfg.Cache.TTL = em.GetInt("CACHE_TTL", cfg.Cache.TTL)
	cfg.Cache.MaxSize = em.GetInt("CACHE_MAX_SIZE", cfg.Cache.MaxSize)
	cfg.Cache.Enabled = em.GetBool("CACHE_ENABLED", cfg.Cache.Enabled)

	return nil
}
