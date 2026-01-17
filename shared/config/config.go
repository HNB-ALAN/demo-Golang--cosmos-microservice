package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure
type Config struct {
	Server     ServerConfig            `mapstructure:"server"`
	Database   DatabaseConfig          `mapstructure:"database"`
	Redis      RedisConfig             `mapstructure:"redis"`
	ClickHouse ClickHouseConfig        `mapstructure:"clickhouse"`
	InfluxDB   InfluxDBConfig          `mapstructure:"influxdb"`
	Quickwit   QuickwitConfig          `mapstructure:"quickwit"`
	VectorDB   VectorDBConfig          `mapstructure:"vectordb"`
	MinIO      MinIOConfig             `mapstructure:"minio"`
	BigQuery   BigQueryConfig          `mapstructure:"bigquery"`
	Kafka      KafkaConfig             `mapstructure:"kafka"`
	Log        LogConfig               `mapstructure:"log"`
	Service    ServiceConfig           `mapstructure:"service"`
	Metrics    MetricsConfig           `mapstructure:"metrics"`
	Auth       AuthConfig              `mapstructure:"auth"`
	Cache      CacheConfig             `mapstructure:"cache"`
	Middleware MiddlewareConfig        `mapstructure:"middleware"`
	GraphQL    GraphQLFederationConfig `mapstructure:"graphql"`
	// Consensus, Mempool, and P2P are owned by Cosmos SDK configs per service
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Host    string        `mapstructure:"host"`
	Port    string        `mapstructure:"port"`
	GRPC    GRPCConfig    `mapstructure:"grpc"`
	HTTP    HTTPConfig    `mapstructure:"http"`
	GraphQL GraphQLConfig `mapstructure:"graphql"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int    `mapstructure:"max_conns"`
	MinConns int    `mapstructure:"min_conns"`
	Enabled  bool   `mapstructure:"enabled"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
	Enabled  bool   `mapstructure:"enabled"`
}

// (removed) RocksDBConfig - legacy custom chain storage
// (removed) MongoDBConfig - MongoDB support removed

// ClickHouseConfig contains ClickHouse configuration
type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Secure   bool   `mapstructure:"secure"`
	Enabled  bool   `mapstructure:"enabled"`
}

// InfluxDBConfig contains InfluxDB configuration
type InfluxDBConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Token   string `mapstructure:"token"`
	Org     string `mapstructure:"org"`
	Bucket  string `mapstructure:"bucket"`
	Enabled bool   `mapstructure:"enabled"`
}

// QuickwitConfig contains Quickwit configuration
type QuickwitConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Index    string `mapstructure:"index"`
	Enabled  bool   `mapstructure:"enabled"`
}

// MinIOConfig contains MinIO configuration
type MinIOConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	SSL       bool   `mapstructure:"ssl"`
	Enabled   bool   `mapstructure:"enabled"`
}

// VectorDBConfig contains Vector DB configuration
type VectorDBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	APIKey   string `mapstructure:"api_key"`
	Database string `mapstructure:"database"`
	Enabled  bool   `mapstructure:"enabled"`
}

// BigQueryConfig contains BigQuery configuration
type BigQueryConfig struct {
	ProjectID string `mapstructure:"project_id"`
	Dataset   string `mapstructure:"dataset"`
	Location  string `mapstructure:"location"`
	Enabled   bool   `mapstructure:"enabled"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	Brokers           []string `mapstructure:"brokers"`
	ClientID          string   `mapstructure:"client_id"`
	GroupID           string   `mapstructure:"group_id"`
	SecurityProtocol  string   `mapstructure:"security_protocol"`
	SASLMechanism     string   `mapstructure:"sasl_mechanism"`
	SASLUsername      string   `mapstructure:"sasl_username"`
	SASLPassword      string   `mapstructure:"sasl_password"`
	SSLCAFile         string   `mapstructure:"ssl_ca_file"`
	SSLCertFile       string   `mapstructure:"ssl_cert_file"`
	SSLKeyFile        string   `mapstructure:"ssl_key_file"`
	SSLKeyPassword    string   `mapstructure:"ssl_key_password"`
	SessionTimeout    int      `mapstructure:"session_timeout"`
	HeartbeatInterval int      `mapstructure:"heartbeat_interval"`
	MaxPollRecords    int      `mapstructure:"max_poll_records"`
	AutoOffsetReset   string   `mapstructure:"auto_offset_reset"`
	EnableAutoCommit  bool     `mapstructure:"enable_auto_commit"`
	CompressionType   string   `mapstructure:"compression_type"`
	BatchSize         int      `mapstructure:"batch_size"`
	LingerMs          int      `mapstructure:"linger_ms"`
	Retries           int      `mapstructure:"retries"`
	RequestTimeout    int      `mapstructure:"request_timeout"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	Filename string `mapstructure:"filename"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
	Compress bool   `mapstructure:"compress"`
}

// ServiceConfig contains service-specific configuration
type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Region      string `mapstructure:"region"`
	Instance    string `mapstructure:"instance"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpiry     string `mapstructure:"jwt_expiry"`
	RefreshExpiry string `mapstructure:"refresh_expiry"`
	Issuer        string `mapstructure:"issuer"`
}

// CacheConfig contains cache configuration
type CacheConfig struct {
	TTL     int  `mapstructure:"ttl"`
	MaxSize int  `mapstructure:"max_size"`
	Enabled bool `mapstructure:"enabled"`
}

// MiddlewareConfig contains middleware configuration
type MiddlewareConfig struct {
	TimeoutSeconds  int  `mapstructure:"timeout_seconds"`
	RateLimitPerSec int  `mapstructure:"rate_limit_per_sec"`
	EnableRecovery  bool `mapstructure:"enable_recovery"`
	EnableLogging   bool `mapstructure:"enable_logging"`
	EnableMetrics   bool `mapstructure:"enable_metrics"`
	EnableAuth      bool `mapstructure:"enable_auth"`
	EnableRateLimit bool `mapstructure:"enable_rate_limit"`
}

// ConsensusConfig contains consensus engine configuration
type ConsensusConfig struct {
	Type             string  `mapstructure:"type"`
	BlockTime        string  `mapstructure:"block_time"`
	ValidatorTimeout string  `mapstructure:"validator_timeout"`
	MinStakeAmount   int64   `mapstructure:"min_stake_amount"`
	MaxValidators    int     `mapstructure:"max_validators"`
	SlashingEnabled  bool    `mapstructure:"slashing_enabled"`
	SlashingPenalty  float64 `mapstructure:"slashing_penalty"`
}

// (removed) MempoolConfig - handled by Cosmos SDK

// (removed) P2PConfig - handled by Cosmos SDK

// GRPCConfig contains gRPC configuration
type GRPCConfig struct {
	MaxRecvMsgSize int  `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize int  `mapstructure:"max_send_msg_size"`
	KeepAlive      bool `mapstructure:"keep_alive"`
}

// HTTPConfig contains HTTP configuration
type HTTPConfig struct {
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
	IdleTimeout  int `mapstructure:"idle_timeout"`
}

// GraphQLConfig contains GraphQL configuration
type GraphQLConfig struct {
	Port          string `mapstructure:"port"`
	Enabled       bool   `mapstructure:"enabled"`
	Playground    bool   `mapstructure:"playground"`
	Introspection bool   `mapstructure:"introspection"`
}

// GraphQLFederationConfig contains GraphQL federation configuration
type GraphQLFederationConfig struct {
	Federation GraphQLFederation `mapstructure:"federation"`
}

// GraphQLFederation contains federation settings
type GraphQLFederation struct {
	GatewayURL   string               `mapstructure:"gateway_url"`
	GRPCServices []GraphQLGRPCService `mapstructure:"grpc_services"`
}

// GraphQLGRPCService represents a gRPC service for GraphQL federation
type GraphQLGRPCService struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath, serviceName string) (*Config, error) {
	cfg := &Config{}

	// Set default values
	setDefaults(cfg, serviceName)

	// Initialize viper
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/usc")

	// Set environment variable prefix
	v.SetEnvPrefix("USC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file if exists
	if configPath != "" {
		v.SetConfigFile(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults and environment variables
	}

	// Unmarshal config
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// GetGraphQLServerAddress returns the GraphQL server address
func (c *Config) GetGraphQLServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.GraphQL.Port)
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)
}

// GetRedisAddress returns the Redis address
func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetClickHouseDSN returns the ClickHouse connection string
func (c *Config) GetClickHouseDSN() string {
	protocol := "http"
	if c.ClickHouse.Secure {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d?username=%s&password=%s&database=%s",
		protocol, c.ClickHouse.Host, c.ClickHouse.Port, c.ClickHouse.User, c.ClickHouse.Password, c.ClickHouse.DBName)
}

// GetInfluxDBURL returns the InfluxDB URL
func (c *Config) GetInfluxDBURL() string {
	return fmt.Sprintf("http://%s:%d", c.InfluxDB.Host, c.InfluxDB.Port)
}

// GetQuickwitURL returns the Quickwit URL
func (c *Config) GetQuickwitURL() string {
	return fmt.Sprintf("http://%s:%d", c.Quickwit.Host, c.Quickwit.Port)
}

// GetMinIOURL returns the MinIO URL
func (c *Config) GetMinIOURL() string {
	protocol := "http"
	if c.MinIO.SSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, c.MinIO.Host, c.MinIO.Port)
}

// GetVectorDBURL returns the VectorDB URL
func (c *Config) GetVectorDBURL() string {
	return fmt.Sprintf("http://%s:%d", c.VectorDB.Host, c.VectorDB.Port)
}

// GetKafkaAddress returns the Kafka broker address
func (c *Config) GetKafkaAddress() string {
	if len(c.Kafka.Brokers) > 0 {
		return c.Kafka.Brokers[0]
	}
	return "localhost:9092"
}

// GetKafkaBrokers returns all Kafka brokers
func (c *Config) GetKafkaBrokers() []string {
	if len(c.Kafka.Brokers) > 0 {
		return c.Kafka.Brokers
	}
	return []string{"localhost:9092"}
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Service.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Service.Environment == "production"
}
