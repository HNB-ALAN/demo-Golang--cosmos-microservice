package store_network_operations

import (
	"context"

	"service-04/internal/application/business/store_network_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles store network operations gRPC requests
type Handlers struct {
	proto.UnimplementedStoreNetworkOperationsServiceServer
	service *store_network_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new store network operations handlers
func NewHandlers(service *store_network_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// SyncStoreNetworkState syncs external network state
func (h *Handlers) SyncStoreNetworkState(ctx context.Context, req *proto.SyncStoreNetworkStateRequest) (*proto.SyncStoreNetworkStateResponse, error) {
	return h.service.SyncStoreNetworkState(ctx, req)
}

// GetStoreNetworkInfo retrieves store network information
func (h *Handlers) GetStoreNetworkInfo(ctx context.Context, req *proto.GetStoreNetworkInfoRequest) (*proto.GetStoreNetworkInfoResponse, error) {
	return h.service.GetStoreNetworkInfo(ctx, req)
}

// UpdateStoreBridgeConfig updates bridge configuration
func (h *Handlers) UpdateStoreBridgeConfig(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest) (*proto.UpdateStoreBridgeConfigResponse, error) {
	return h.service.UpdateStoreBridgeConfig(ctx, req)
}
