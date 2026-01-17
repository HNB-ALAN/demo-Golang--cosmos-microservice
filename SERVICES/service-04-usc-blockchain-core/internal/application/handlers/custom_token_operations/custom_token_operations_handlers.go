package custom_token_operations

import (
	"context"

	"service-04/internal/application/business/custom_token_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles custom token operations gRPC requests
type Handlers struct {
	proto.UnimplementedCustomTokenOperationsServiceServer
	service *custom_token_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new custom token operations handlers
func NewHandlers(service *custom_token_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// CreateBlockchainToken creates a new custom token
func (h *Handlers) CreateBlockchainToken(ctx context.Context, req *proto.CreateBlockchainTokenRequest) (*proto.CreateBlockchainTokenResponse, error) {
	return h.service.CreateBlockchainToken(ctx, req)
}

// MintTokens mints custom tokens
func (h *Handlers) MintTokens(ctx context.Context, req *proto.MintTokensRequest) (*proto.MintTokensResponse, error) {
	return h.service.MintTokens(ctx, req)
}

// GetTokenBalance retrieves token balance for an address
func (h *Handlers) GetTokenBalance(ctx context.Context, req *proto.GetTokenBalanceRequest) (*proto.GetTokenBalanceResponse, error) {
	return h.service.GetTokenBalance(ctx, req)
}

// GetTokenInfo retrieves token information
func (h *Handlers) GetTokenInfo(ctx context.Context, req *proto.GetTokenInfoRequest) (*proto.GetTokenInfoResponse, error) {
	return h.service.GetTokenInfo(ctx, req)
}

// BurnTokens burns custom tokens
func (h *Handlers) BurnTokens(ctx context.Context, req *proto.BurnTokensRequest) (*proto.BurnTokensResponse, error) {
	return h.service.BurnTokens(ctx, req)
}
