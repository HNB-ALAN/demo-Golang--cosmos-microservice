package constants

import (
	"fmt"
	"testing"
)

func TestPortAllocation(t *testing.T) {
	// Test no port conflicts
	usedPorts := make(map[int]string)

	// Check main service ports
	for serviceName, port := range ServicePortMap {
		if existingService, exists := usedPorts[port]; exists {
			t.Errorf("Port conflict: %s and %s both use port %d", serviceName, existingService, port)
		}
		usedPorts[port] = serviceName
	}

	// Check metrics ports
	for serviceName, port := range ServiceMetricsPortMap {
		if existingService, exists := usedPorts[port]; exists {
			t.Errorf("Port conflict: %s and %s both use metrics port %d", serviceName, existingService, port)
		}
		usedPorts[port] = serviceName + "-metrics"
	}

	// Test special ports don't conflict
	specialPorts := map[int]string{
		PortGatewayGraphQL:        "Gateway GraphQL",
		PortSocialWebSocket:       "Social WebSocket",
		PortNotificationWebSocket: "Notification WebSocket",
		PortRedisCluster:          "Redis Cluster",
		PortVideoRTMP:             "Video RTMP",
	}

	for port, name := range specialPorts {
		if existingService, exists := usedPorts[port]; exists {
			t.Errorf("Port conflict: %s and %s both use port %d", name, existingService, port)
		}
		usedPorts[port] = name
	}

	fmt.Printf("✅ Total unique ports allocated: %d\n", len(usedPorts))
}

func TestServicePortRange(t *testing.T) {
	// Test main service ports are in range 8001-8022
	for serviceName, port := range ServicePortMap {
		if port < 8001 || port > 8022 {
			t.Errorf("Service %s port %d is outside expected range 8001-8022", serviceName, port)
		}
	}

	// Test metrics ports are in range 9001-9022
	for serviceName, port := range ServiceMetricsPortMap {
		if port < 9001 || port > 9022 {
			t.Errorf("Service %s metrics port %d is outside expected range 9001-9022", serviceName, port)
		}
	}
}

func TestPortValidation(t *testing.T) {
	// Test service port validation
	validPorts := []int{8001, 8005, 8010, 8022, 4000, 8090, 8091, 7000, 1935}
	for _, port := range validPorts {
		if !IsValidServicePort(port) {
			t.Errorf("Port %d should be valid but validation failed", port)
		}
	}

	// Test invalid ports
	invalidPorts := []int{8000, 8023, 7999, 8100}
	for _, port := range invalidPorts {
		if IsValidServicePort(port) {
			t.Errorf("Port %d should be invalid but validation passed", port)
		}
	}
}

func TestServiceFunctions(t *testing.T) {
	// Test GetServicePort
	gatewayPort := GetServicePort(ServiceGateway)
	if gatewayPort != PortGateway {
		t.Errorf("Expected Gateway port %d, got %d", PortGateway, gatewayPort)
	}

	// Test GetServiceGraphQLPort
	graphqlPort := GetServiceGraphQLPort(ServiceGateway)
	if graphqlPort != PortGatewayGraphQL {
		t.Errorf("Expected Gateway GraphQL port %d, got %d", PortGatewayGraphQL, graphqlPort)
	}

	// Test non-Gateway service has no GraphQL
	authGraphQL := GetServiceGraphQLPort(ServiceAuth)
	if authGraphQL != 0 {
		t.Errorf("Auth service should not have GraphQL port, got %d", authGraphQL)
	}

	// Test WebSocket ports
	socialWS := GetServiceWebSocketPort(ServiceSocial)
	if socialWS != PortSocialWebSocket {
		t.Errorf("Expected Social WebSocket port %d, got %d", PortSocialWebSocket, socialWS)
	}

	notificationWS := GetServiceWebSocketPort(ServiceNotification)
	if notificationWS != PortNotificationWebSocket {
		t.Errorf("Expected Notification WebSocket port %d, got %d", PortNotificationWebSocket, notificationWS)
	}
}

func TestPortPrintFunction(t *testing.T) {
	// Test that PrintPortAllocation doesn't panic
	allocation := PrintPortAllocation()
	if len(allocation) == 0 {
		t.Error("PrintPortAllocation returned empty string")
	}
	fmt.Println(allocation)
}
