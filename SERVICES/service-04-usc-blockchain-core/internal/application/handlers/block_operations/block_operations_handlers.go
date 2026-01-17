package block_operations

import (
	"context"

	"service-04/internal/application/business/block_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handlers handles block operations gRPC requests
type Handlers struct {
	proto.UnimplementedBlockOperationsServiceServer
	service *block_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new block operations handlers
func NewHandlers(service *block_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// ProduceBlock creates a new blockchain block
func (h *Handlers) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
	return h.service.ProduceBlock(ctx, req)
}

// ValidateBlock validates block integrity
func (h *Handlers) ValidateBlock(ctx context.Context, req *proto.ValidateBlockRequest) (*proto.ValidateBlockResponse, error) {
	return h.service.ValidateBlock(ctx, req)
}

// GetBlock retrieves specific block
func (h *Handlers) GetBlock(ctx context.Context, req *proto.GetBlockRequest) (*proto.GetBlockResponse, error) {
	return h.service.GetBlock(ctx, req)
}

// GetBlockByHash gets block by hash
func (h *Handlers) GetBlockByHash(ctx context.Context, req *proto.GetBlockByHashRequest) (*proto.GetBlockResponse, error) {
	return h.service.GetBlockByHash(ctx, req)
}

// GetLatestBlock gets current blockchain head
func (h *Handlers) GetLatestBlock(ctx context.Context, req *emptypb.Empty) (*proto.GetBlockResponse, error) {
	return h.service.GetLatestBlock(ctx, req)
}

// GetBlockRange gets blockchain block range data
func (h *Handlers) GetBlockRange(ctx context.Context, req *proto.GetBlockRangeRequest) (*proto.GetBlockRangeResponse, error) {
	return h.service.GetBlockRange(ctx, req)
}
