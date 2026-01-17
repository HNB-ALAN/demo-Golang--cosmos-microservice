// Package constants provides application constants for USC platform services.
package constants

// Application constants
const (
	// Application info
	AppName        = "USC Platform"
	AppVersion     = "1.0.0"
	AppDescription = "USC Social Media Platform"

	// Environment
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"

	// Default values
	DefaultPageSize   = 20
	MaxPageSize       = 100
	DefaultTimeout    = 30
	DefaultRetryCount = 3
	DefaultRetryDelay = 1

	// Cache
	DefaultCacheTTL  = 3600 // 1 hour
	DefaultCacheSize = 1000
	CacheKeyPrefix   = "usc:"

	// Database
	DefaultMaxConns    = 25
	DefaultMinConns    = 5
	DefaultConnTimeout = 10

	// HTTP
	DefaultHTTPPort     = 8080
	DefaultHTTPSPort    = 8443
	DefaultReadTimeout  = 30
	DefaultWriteTimeout = 30
	DefaultIdleTimeout  = 120

	// gRPC
	DefaultGRPCPort    = 9090
	DefaultMaxRecvSize = 4 * 1024 * 1024 // 4MB
	DefaultMaxSendSize = 4 * 1024 * 1024 // 4MB

	// Logging
	DefaultLogLevel  = "info"
	DefaultLogFormat = "json"
	DefaultLogOutput = "stdout"

	// Metrics
	DefaultMetricsPort = 9090
	DefaultMetricsPath = "/metrics"

	// Health
	DefaultHealthPort = 8081
	DefaultHealthPath = "/health"

	// JWT
	DefaultJWTExpiry     = "24h"
	DefaultRefreshExpiry = "168h" // 7 days
	DefaultJWTIssuer     = "usc-platform"

	// Rate limiting
	DefaultRateLimit  = 100
	DefaultRateWindow = 60 // seconds

	// File upload
	MaxFileSize       = 10 * 1024 * 1024 // 10MB
	AllowedImageTypes = "image/jpeg,image/png,image/gif,image/webp"
	AllowedVideoTypes = "video/mp4,video/webm,video/ogg"
	AllowedAudioTypes = "audio/mp3,audio/wav,audio/ogg"

	// Validation
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MinEmailLength    = 5
	MaxEmailLength    = 254

	// Pagination
	MinPage = 1
	MaxPage = 10000

	// Search
	DefaultSearchLimit   = 20
	MaxSearchLimit       = 100
	MinSearchQueryLength = 2
	MaxSearchQueryLength = 100

	// Content
	MaxPostLength    = 2000
	MaxCommentLength = 500
	MaxBioLength     = 160
	MaxTitleLength   = 100

	// Notifications
	MaxNotificationCount   = 100
	DefaultNotificationTTL = 86400 // 24 hours

	// Social
	MaxFriendsCount   = 5000
	MaxFollowersCount = 10000
	MaxFollowingCount = 5000

	// Security
	MaxLoginAttempts = 5
	LockoutDuration  = 900   // 15 minutes
	SessionTimeout   = 86400 // 24 hours

	// API
	DefaultAPIVersion     = "v1"
	MaxRequestSize        = 10 * 1024 * 1024 // 10MB
	DefaultRequestTimeout = 30

	// WebSocket
	DefaultWSReadBufferSize   = 1024
	DefaultWSWriteBufferSize  = 1024
	DefaultWSHandshakeTimeout = 10

	// Background jobs
	DefaultJobTimeout = 300 // 5 minutes
	MaxJobRetries     = 3
	DefaultJobDelay   = 1

	// Monitoring
	DefaultHealthCheckInterval = 30
	DefaultMetricsInterval     = 60

	// Storage
	DefaultStorageProvider = "local"
	MaxStorageQuota        = 1024 * 1024 * 1024 // 1GB

	// Analytics
	DefaultAnalyticsInterval = 3600 // 1 hour
	MaxAnalyticsRetention    = 90   // days

	// Backup
	DefaultBackupInterval = 86400 // 24 hours
	MaxBackupRetention    = 30    // days

	// Migration
	DefaultMigrationTimeout = 300 // 5 minutes
	MaxMigrationRetries     = 3

	// Testing
	DefaultTestTimeout = 30
	DefaultTestRetries = 3
)
