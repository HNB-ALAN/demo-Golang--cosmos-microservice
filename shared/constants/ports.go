// Package constants provides application constants for USC platform services.
package constants

import "fmt"

// Individual service ports for all 22 services - No conflicts
const (
	// Service 01 - Gateway
	PortGateway        = 8001
	PortGatewayGraphQL = 4000 // GraphQL endpoint
	PortGatewayMetrics = 9001 // Metrics endpoint

	// Service 02 - Auth
	PortAuth        = 8002
	PortAuthMetrics = 9002

	// Service 03 - User
	PortUser        = 8003
	PortUserMetrics = 9003

	// Service 04 - Blockchain Core
	PortBlockchainCore        = 8004
	PortBlockchainCoreMetrics = 9004
	PortBlockchainP2P         = 30303 // P2P networking
	PortBlockchainRPC         = 30301 // Tendermint/Cosmos RPC

	// Service 05 - Wallet
	PortWallet        = 8005
	PortWalletMetrics = 9005

	// Service 06 - Security
	PortSecurity        = 8006
	PortSecurityMetrics = 9006

	// Service 07 - Caching
	PortCaching        = 8007
	PortCachingMetrics = 9007
	PortRedisCluster   = 7000 // Redis cluster communication

	// Service 08 - Monitoring
	PortMonitoring        = 8008
	PortMonitoringMetrics = 9008

	// Service 09 - Social
	PortSocial          = 8009
	PortSocialMetrics   = 9009
	PortSocialWebSocket = 8090 // WebSocket for real-time

	// Service 10 - Bilateral Rewards
	PortBilateralRewards        = 8010
	PortBilateralRewardsMetrics = 9010

	// Service 11 - Content Management
	PortContentManagement        = 8011
	PortContentManagementMetrics = 9011

	// Service 12 - Video
	PortVideo        = 8012
	PortVideoMetrics = 9012
	PortVideoRTMP    = 1935 // RTMP streaming

	// Service 13 - AI
	PortAI        = 8013
	PortAIMetrics = 9013

	// Service 14 - Commerce
	PortCommerce        = 8014
	PortCommerceMetrics = 9014

	// Service 15 - Notification
	PortNotification          = 8015
	PortNotificationMetrics   = 9015
	PortNotificationWebSocket = 8091 // WebSocket for real-time

	// Service 16 - Search
	PortSearch        = 8016
	PortSearchMetrics = 9016

	// Service 17 - Analytics
	PortAnalytics        = 8017
	PortAnalyticsMetrics = 9017

	// Service 18 - Moderation
	PortModeration        = 8018
	PortModerationMetrics = 9018

	// Service 19 - Recommendation
	PortRecommendation        = 8019
	PortRecommendationMetrics = 9019

	// Service 20 - Advertising
	PortAdvertising        = 8020
	PortAdvertisingMetrics = 9020

	// Service 21 - Admin
	PortAdmin        = 8021
	PortAdminMetrics = 9021

	// Service 22 - Kafka Messaging
	PortKafkaMessaging        = 8022
	PortKafkaMessagingMetrics = 9022
)

// Database ports
const (
	PortPostgreSQL = 5432
	PortRedis      = 6379
	PortClickHouse = 9000
	PortInfluxDB   = 8086
	PortQuickwit   = 7280
)

// Monitoring ports
const (
	PortPrometheus   = 9091
	PortGrafana      = 3000
	PortAlertManager = 9093
	PortJaeger       = 16686
	PortZipkin       = 9411
)

// Development ports
const (
	PortDevServer    = 3001
	PortDevProxy     = 3002
	PortDevWebpack   = 3003
	PortDevHotReload = 3004
)

// Testing ports - Updated to avoid conflicts with service ports
const (
	PortTestServer     = 8000
	PortTestDatabase   = 5433
	PortTestRedis      = 6380
	PortTestClickHouse = 9050 // Changed from 9001 to avoid conflict
	PortTestInfluxDB   = 8087
	PortTestQuickwit   = 7281
)

// Port ranges - Updated for individual port allocation
const (
	// Service port range - individual ports per service
	ServicePortStart = 8001
	ServicePortEnd   = 8022

	// Metrics port range - individual metrics ports per service
	MetricsPortStart = 9001
	MetricsPortEnd   = 9022

	// Database port range
	DatabasePortStart = 5000
	DatabasePortEnd   = 6000

	// Monitoring port range
	MonitoringPortStart = 9000
	MonitoringPortEnd   = 9100

	// Development port range
	DevPortStart = 3000
	DevPortEnd   = 4000

	// Testing port range
	TestPortStart = 8000
	TestPortEnd   = 8100
)

// Service port mappings - Individual ports for each service
var ServicePortMap = map[string]int{
	ServiceGateway:           PortGateway,
	ServiceAuth:              PortAuth,
	ServiceUser:              PortUser,
	ServiceBlockchainCore:    PortBlockchainCore,
	ServiceWallet:            PortWallet,
	ServiceSecurity:          PortSecurity,
	ServiceCaching:           PortCaching,
	ServiceMonitoring:        PortMonitoring,
	ServiceSocial:            PortSocial,
	ServiceBilateralRewards:  PortBilateralRewards,
	ServiceContentManagement: PortContentManagement,
	ServiceVideo:             PortVideo,
	ServiceAI:                PortAI,
	ServiceCommerce:          PortCommerce,
	ServiceNotification:      PortNotification,
	ServiceSearch:            PortSearch,
	ServiceAnalytics:         PortAnalytics,
	ServiceModeration:        PortModeration,
	ServiceRecommendation:    PortRecommendation,
	ServiceAdvertising:       PortAdvertising,
	ServiceAdmin:             PortAdmin,
	ServiceKafkaMessaging:    PortKafkaMessaging,
}

// Service metrics port mappings
var ServiceMetricsPortMap = map[string]int{
	ServiceGateway:           PortGatewayMetrics,
	ServiceAuth:              PortAuthMetrics,
	ServiceUser:              PortUserMetrics,
	ServiceBlockchainCore:    PortBlockchainCoreMetrics,
	ServiceWallet:            PortWalletMetrics,
	ServiceSecurity:          PortSecurityMetrics,
	ServiceCaching:           PortCachingMetrics,
	ServiceMonitoring:        PortMonitoringMetrics,
	ServiceSocial:            PortSocialMetrics,
	ServiceBilateralRewards:  PortBilateralRewardsMetrics,
	ServiceContentManagement: PortContentManagementMetrics,
	ServiceVideo:             PortVideoMetrics,
	ServiceAI:                PortAIMetrics,
	ServiceCommerce:          PortCommerceMetrics,
	ServiceNotification:      PortNotificationMetrics,
	ServiceSearch:            PortSearchMetrics,
	ServiceAnalytics:         PortAnalyticsMetrics,
	ServiceModeration:        PortModerationMetrics,
	ServiceRecommendation:    PortRecommendationMetrics,
	ServiceAdvertising:       PortAdvertisingMetrics,
	ServiceAdmin:             PortAdminMetrics,
	ServiceKafkaMessaging:    PortKafkaMessagingMetrics,
}

// Individual port allocation for zero conflicts
func GetServicePort(serviceName string) int {
	if port, exists := ServicePortMap[serviceName]; exists {
		return port
	}
	return PortGateway // Default to Gateway port
}

// GetServiceGraphQLPort returns GraphQL port for Gateway
func GetServiceGraphQLPort(serviceName string) int {
	if serviceName == ServiceGateway {
		return PortGatewayGraphQL // 4000
	}
	return 0 // No GraphQL for other services
}

// GetServiceMetricsPort returns individual metrics port for each service
func GetServiceMetricsPort(serviceName string) int {
	if port, exists := ServiceMetricsPortMap[serviceName]; exists {
		return port
	}
	return PortGatewayMetrics // Default to Gateway metrics port
}

// GetServiceWebSocketPort returns WebSocket port for real-time services
func GetServiceWebSocketPort(serviceName string) int {
	if serviceName == ServiceSocial {
		return PortSocialWebSocket // 8090
	}
	if serviceName == ServiceNotification {
		return PortNotificationWebSocket // 8091
	}
	return 0 // No WebSocket for other services
}

// GetServiceSpecialPorts returns special ports for specific services
func GetServiceSpecialPorts(serviceName string) map[string]int {
	ports := make(map[string]int)

	switch serviceName {
	case "service-04-usc-blockchain-core":
		ports["p2p"] = PortBlockchainP2P
		ports["rpc"] = PortBlockchainRPC
	case "service-07-caching":
		ports["cluster"] = PortRedisCluster
	case "service-12-video-service":
		ports["rtmp"] = PortVideoRTMP
	}

	return ports
}

// Database port mappings
var DatabasePortMap = map[string]int{
	"postgresql": PortPostgreSQL,
	"redis":      PortRedis,
	"clickhouse": PortClickHouse,
	"influxdb":   PortInfluxDB,
	"quickwit":   PortQuickwit,
}

// Monitoring port mappings
var MonitoringPortMap = map[string]int{
	"prometheus":   PortPrometheus,
	"grafana":      PortGrafana,
	"alertmanager": PortAlertManager,
	"jaeger":       PortJaeger,
	"zipkin":       PortZipkin,
}

// GetDatabasePort returns the port for a database
func GetDatabasePort(databaseName string) int {
	if port, exists := DatabasePortMap[databaseName]; exists {
		return port
	}
	return PortPostgreSQL // Default database port
}

// GetMonitoringPort returns the port for a monitoring service
func GetMonitoringPort(monitoringName string) int {
	if port, exists := MonitoringPortMap[monitoringName]; exists {
		return port
	}
	return PortPrometheus // Default monitoring port
}

// IsValidServicePort checks if a port is valid for services
func IsValidServicePort(port int) bool {
	// Check main service ports (8001-8022)
	if port >= 8001 && port <= 8022 {
		return true
	}
	// Check special ports
	return port == PortGatewayGraphQL ||
		port == PortSocialWebSocket || port == PortNotificationWebSocket ||
		port == PortBlockchainP2P || port == PortBlockchainRPC ||
		port == PortRedisCluster || port == PortVideoRTMP
}

// IsValidMetricsPort checks if a port is valid for metrics
func IsValidMetricsPort(port int) bool {
	// Check metrics ports (9001-9022)
	return port >= 9001 && port <= 9022
}

// IsValidDatabasePort checks if a port is valid for databases
func IsValidDatabasePort(port int) bool {
	return port >= DatabasePortStart && port <= DatabasePortEnd
}

// IsValidMonitoringPort checks if a port is valid for monitoring services
func IsValidMonitoringPort(port int) bool {
	return port >= MonitoringPortStart && port <= MonitoringPortEnd
}

// IsValidDevPort checks if a port is valid for development
func IsValidDevPort(port int) bool {
	return port >= DevPortStart && port <= DevPortEnd
}

// IsValidTestPort checks if a port is valid for testing
func IsValidTestPort(port int) bool {
	return port >= TestPortStart && port <= TestPortEnd
}

// GetAllServicePorts returns all service main ports for documentation
func GetAllServicePorts() map[string]int {
	return ServicePortMap
}

// GetAllMetricsPorts returns all metrics ports for documentation
func GetAllMetricsPorts() map[string]int {
	return ServiceMetricsPortMap
}

// PrintPortAllocation prints formatted port allocation table
func PrintPortAllocation() string {
	result := "=== USC PLATFORM PORT ALLOCATION ===\n"
	result += "Service                    Main Port  Metrics Port  Special Ports\n"
	result += "--------------------------------------------------------\n"

	services := []struct{ name, key string }{
		{"Gateway", ServiceGateway},
		{"Auth", ServiceAuth},
		{"User", ServiceUser},
		{"Blockchain Core", ServiceBlockchainCore},
		{"Wallet", ServiceWallet},
		{"Security", ServiceSecurity},
		{"Caching", ServiceCaching},
		{"Monitoring", ServiceMonitoring},
		{"Social", ServiceSocial},
		{"Bilateral Rewards", ServiceBilateralRewards},
		{"Content Management", ServiceContentManagement},
		{"Video", ServiceVideo},
		{"AI", ServiceAI},
		{"Commerce", ServiceCommerce},
		{"Notification", ServiceNotification},
		{"Search", ServiceSearch},
		{"Analytics", ServiceAnalytics},
		{"Moderation", ServiceModeration},
		{"Recommendation", ServiceRecommendation},
		{"Advertising", ServiceAdvertising},
		{"Admin", ServiceAdmin},
		{"Kafka Messaging", ServiceKafkaMessaging},
	}

	for _, svc := range services {
		mainPort := GetServicePort(svc.key)
		metricsPort := GetServiceMetricsPort(svc.key)
		special := ""

		switch svc.key {
		case ServiceGateway:
			special = "GraphQL:4000"
		case ServiceSocial:
			special = "WebSocket:8090"
		case ServiceNotification:
			special = "WebSocket:8091"
		case ServiceBlockchainCore:
			special = "P2P:30303,RPC:30301"
		case ServiceCaching:
			special = "Cluster:7000"
		case ServiceVideo:
			special = "RTMP:1935"
		}

		result += fmt.Sprintf("%-25s %d       %d          %s\n",
			svc.name, mainPort, metricsPort, special)
	}

	return result
}
