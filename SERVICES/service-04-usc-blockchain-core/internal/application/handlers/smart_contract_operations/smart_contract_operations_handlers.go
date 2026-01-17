package smart_contract_operations

import (
	"context"

	"service-04/internal/application/business/smart_contract_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles smart contract operations gRPC requests
type Handlers struct {
	proto.UnimplementedSmartContractOperationsServiceServer
	service *smart_contract_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new smart contract operations handlers
func NewHandlers(service *smart_contract_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// DeployContract deploys a new smart contract
func (h *Handlers) DeployContract(ctx context.Context, req *proto.DeployContractRequest) (*proto.DeployContractResponse, error) {
	return h.service.DeployContract(ctx, req)
}

// ExecuteContract executes a smart contract function
func (h *Handlers) ExecuteContract(ctx context.Context, req *proto.ExecuteContractRequest) (*proto.ExecuteContractResponse, error) {
	return h.service.ExecuteContract(ctx, req)
}

// QueryContract queries a smart contract function
func (h *Handlers) QueryContract(ctx context.Context, req *proto.QueryContractRequest) (*proto.QueryContractResponse, error) {
	return h.service.QueryContract(ctx, req)
}

// GetContractCode retrieves contract source code
func (h *Handlers) GetContractCode(ctx context.Context, req *proto.GetContractCodeRequest) (*proto.GetContractCodeResponse, error) {
	return h.service.GetContractCode(ctx, req)
}

// GetContractStorage retrieves contract storage
func (h *Handlers) GetContractStorage(ctx context.Context, req *proto.GetContractStorageRequest) (*proto.GetContractStorageResponse, error) {
	return h.service.GetContractStorage(ctx, req)
}
