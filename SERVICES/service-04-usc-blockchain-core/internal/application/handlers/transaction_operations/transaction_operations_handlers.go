package transaction_operations

import (
	"context"

	"service-04/internal/application/business/transaction_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles transaction operations gRPC requests
type Handlers struct {
	proto.UnimplementedTransactionOperationsServiceServer
	service *transaction_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new transaction operations handlers
func NewHandlers(service *transaction_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// SubmitTransaction submits a new transaction to the blockchain
func (h *Handlers) SubmitTransaction(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	return h.service.SubmitTransaction(ctx, req)
}

// GetTransaction retrieves a transaction by hash
func (h *Handlers) GetTransaction(ctx context.Context, req *proto.GetTransactionRequest) (*proto.GetTransactionResponse, error) {
	return h.service.GetTransaction(ctx, req)
}

// GetTransactionStatus retrieves transaction status
func (h *Handlers) GetTransactionStatus(ctx context.Context, req *proto.GetTransactionStatusRequest) (*proto.GetTransactionStatusResponse, error) {
	return h.service.GetTransactionStatus(ctx, req)
}

// GetPendingTransactions retrieves pending transactions
func (h *Handlers) GetPendingTransactions(ctx context.Context, req *proto.GetPendingTransactionsRequest) (*proto.GetPendingTransactionsResponse, error) {
	return h.service.GetPendingTransactions(ctx, req)
}

// EstimateTransactionFee estimates transaction fee
func (h *Handlers) EstimateTransactionFee(ctx context.Context, req *proto.EstimateTransactionFeeRequest) (*proto.EstimateTransactionFeeResponse, error) {
	return h.service.EstimateTransactionFee(ctx, req)
}
