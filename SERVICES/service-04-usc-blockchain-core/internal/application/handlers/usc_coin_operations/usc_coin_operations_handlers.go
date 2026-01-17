package usc_coin_operations

import (
	"context"

	"service-04/internal/application/business/usc_coin_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handlers handles USC coin operations gRPC requests
type Handlers struct {
	proto.UnimplementedUSCCoinOperationsServiceServer
	service *usc_coin_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new USC coin operations handlers
func NewHandlers(service *usc_coin_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// GetWalletBalance retrieves USC balance for an address (matches proto)
func (h *Handlers) GetWalletBalance(ctx context.Context, req *proto.GetWalletBalanceRequest) (*proto.GetWalletBalanceResponse, error) {
	return h.service.GetUSCBalance(ctx, req)
}

// TransferUSCBlockchain transfers USC between addresses (matches proto)
func (h *Handlers) TransferUSCBlockchain(ctx context.Context, req *proto.TransferUSCBlockchainRequest) (*proto.TransferUSCBlockchainResponse, error) {
	return h.service.TransferUSC(ctx, req)
}

// GetUSCSupply retrieves total USC supply
func (h *Handlers) GetUSCSupply(ctx context.Context, req *emptypb.Empty) (*proto.GetUSCSupplyResponse, error) {
	return h.service.GetUSCSupply(ctx)
}

// GetTransactionHistory retrieves transaction history for an address
func (h *Handlers) GetTransactionHistory(ctx context.Context, req *proto.GetTransactionHistoryRequest) (*proto.GetTransactionHistoryResponse, error) {
	return h.service.GetTransactionHistory(ctx, req)
}

// GetUSCTransactions retrieves USC-specific transactions
func (h *Handlers) GetUSCTransactions(ctx context.Context, req *proto.GetUSCTransactionsRequest) (*proto.GetUSCTransactionsResponse, error) {
	return h.service.GetUSCTransactions(ctx, req)
}
