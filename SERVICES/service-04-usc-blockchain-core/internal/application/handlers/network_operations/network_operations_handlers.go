package network_operations

import (
	"context"

	"service-04/internal/application/business/network_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handlers handles network operations gRPC requests
type Handlers struct {
	proto.UnimplementedNetworkOperationsServiceServer
	service *network_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new network operations handlers
func NewHandlers(service *network_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// GetNetworkInfo retrieves network information
func (h *Handlers) GetNetworkInfo(ctx context.Context, req *emptypb.Empty) (*proto.GetNetworkInfoResponse, error) {
	return h.service.GetNetworkInfo(ctx)
}

// GetChainInfo retrieves chain information
func (h *Handlers) GetChainInfo(ctx context.Context, req *emptypb.Empty) (*proto.GetChainInfoResponse, error) {
	return h.service.GetChainInfo(ctx)
}

// GetPeers retrieves list of peers
func (h *Handlers) GetPeers(ctx context.Context, req *proto.GetPeersRequest) (*proto.GetPeersResponse, error) {
	return h.service.GetPeers(ctx, req)
}

// GetNetworkStats retrieves network statistics
func (h *Handlers) GetNetworkStats(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
	return h.service.GetNetworkStats(ctx, req)
}
