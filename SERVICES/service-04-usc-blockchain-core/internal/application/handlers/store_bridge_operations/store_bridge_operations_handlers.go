package store_bridge_operations

import (
	"context"

	"service-04/internal/application/business/store_bridge_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles store bridge operations gRPC requests
type Handlers struct {
	proto.UnimplementedStoreBridgeOperationsServiceServer
	service *store_bridge_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new store bridge operations handlers
func NewHandlers(service *store_bridge_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// DeployStoreBridge deploys a new store bridge
func (h *Handlers) DeployStoreBridge(ctx context.Context, req *proto.DeployStoreBridgeRequest) (*proto.DeployStoreBridgeResponse, error) {
	return h.service.DeployStoreBridge(ctx, req)
}

// RegisterStoreNetwork registers a new store network
func (h *Handlers) RegisterStoreNetwork(ctx context.Context, req *proto.RegisterStoreNetworkRequest) (*proto.RegisterStoreNetworkResponse, error) {
	return h.service.RegisterStoreNetwork(ctx, req)
}

// BridgeStoreTokenToUSC bridges store tokens to USC
func (h *Handlers) BridgeStoreTokenToUSC(ctx context.Context, req *proto.BridgeStoreTokenToUSCRequest) (*proto.BridgeStoreTokenToUSCResponse, error) {
	return h.service.BridgeStoreTokenToUSC(ctx, req)
}

// BridgeUSCToStoreToken bridges USC to store tokens
func (h *Handlers) BridgeUSCToStoreToken(ctx context.Context, req *proto.BridgeUSCToStoreTokenRequest) (*proto.BridgeUSCToStoreTokenResponse, error) {
	return h.service.BridgeUSCToStoreToken(ctx, req)
}

// GetStoreBridgeMetrics retrieves store bridge metrics
func (h *Handlers) GetStoreBridgeMetrics(ctx context.Context, req *proto.GetStoreBridgeMetricsRequest) (*proto.GetStoreBridgeMetricsResponse, error) {
	return h.service.GetStoreBridgeMetrics(ctx, req)
}

// ValidateStoreBridge validates a store bridge
func (h *Handlers) ValidateStoreBridge(ctx context.Context, req *proto.ValidateStoreBridgeRequest) (*proto.ValidateStoreBridgeResponse, error) {
	return h.service.ValidateStoreBridge(ctx, req)
}
