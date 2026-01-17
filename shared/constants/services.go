// Package constants provides application constants for USC platform services.
package constants

// Service names for all 22 services (matching K8s deployment names)
const (
	ServiceGateway           = "service-01-gateway"
	ServiceAuth              = "service-02-auth"
	ServiceUser              = "service-03-user"
	ServiceBlockchainCore    = "service-04-usc-blockchain-core"
	ServiceWallet            = "service-05-usc-wallet"
	ServiceSecurity          = "service-06-security"
	ServiceCaching           = "service-07-caching"
	ServiceMonitoring        = "service-08-monitoring"
	ServiceSocial            = "service-09-social"
	ServiceBilateralRewards  = "service-10-usc-bilateral-rewards"
	ServiceContentManagement = "service-11-content-management"
	ServiceVideo             = "service-12-video-service"
	ServiceAI                = "service-13-ai-service"
	ServiceCommerce          = "service-14-commerce-service"
	ServiceNotification      = "service-15-notification-service"
	ServiceSearch            = "service-16-search-service"
	ServiceAnalytics         = "service-17-analytics-service"
	ServiceModeration        = "service-18-moderation-service"
	ServiceRecommendation    = "service-19-recommendation-service"
	ServiceAdvertising       = "service-20-advertising-service"
	ServiceAdmin             = "service-21-admin-service"
	ServiceKafkaMessaging    = "service-22-kafka-messaging-service"
)

// Service categories
const (
	CategoryCore         = "core"
	CategoryUser         = "user"
	CategoryBlockchain   = "blockchain"
	CategoryTransaction  = "transaction"
	CategoryNotification = "notification"
	CategoryAudit        = "audit"
	CategoryCompliance   = "compliance"
	CategoryRisk         = "risk"
	CategoryContent      = "content"
	CategoryDocument     = "document"
	CategoryStorage      = "storage"
	CategoryAnalytics    = "analytics"
	CategoryReporting    = "reporting"
	CategoryIntegration  = "integration"
	CategoryExternal     = "external"
	CategoryWorkflow     = "workflow"
	CategoryTask         = "task"
	CategoryMonitoring   = "monitoring"
	CategoryHealth       = "health"
)

// Service status
const (
	StatusActive      = "active"
	StatusInactive    = "inactive"
	StatusMaintenance = "maintenance"
	StatusDeprecated  = "deprecated"
)

// Service versions
const (
	VersionV1 = "v1"
	VersionV2 = "v2"
	VersionV3 = "v3"
)

// Service endpoints
const (
	EndpointHealth    = "/health"
	EndpointMetrics   = "/metrics"
	EndpointStatus    = "/status"
	EndpointVersion   = "/version"
	EndpointDocs      = "/docs"
	EndpointSwagger   = "/swagger"
	EndpointAPI       = "/api"
	EndpointGRPC      = "/grpc"
	EndpointWebSocket = "/ws"
	EndpointWebhook   = "/webhook"
)
