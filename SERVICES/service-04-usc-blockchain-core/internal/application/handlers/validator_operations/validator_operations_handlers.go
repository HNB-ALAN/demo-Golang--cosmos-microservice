package validator_operations

import (
	"context"

	"service-04/internal/application/business/validator_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles validator operations gRPC requests
type Handlers struct {
	proto.UnimplementedValidatorOperationsServiceServer
	service *validator_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new validator operations handlers
func NewHandlers(service *validator_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// RegisterValidator handles validator registration
func (h *Handlers) RegisterValidator(ctx context.Context, req *proto.RegisterValidatorRequest) (*proto.RegisterValidatorResponse, error) {
	return h.service.RegisterValidator(ctx, req)
}

// GetValidators handles getting validators list
func (h *Handlers) GetValidators(ctx context.Context, req *proto.GetValidatorsRequest) (*proto.GetValidatorsResponse, error) {
	return h.service.GetValidators(ctx, req)
}

// GetValidatorStatus handles getting validator status
func (h *Handlers) GetValidatorStatus(ctx context.Context, req *proto.GetValidatorStatusRequest) (*proto.GetValidatorStatusResponse, error) {
	return h.service.GetValidatorStatus(ctx, req)
}

// StakeUSC handles USC staking
func (h *Handlers) StakeUSC(ctx context.Context, req *proto.StakeUSCRequest) (*proto.StakeUSCResponse, error) {
	return h.service.StakeUSC(ctx, req)
}

// UnstakeUSC handles USC unstaking
func (h *Handlers) UnstakeUSC(ctx context.Context, req *proto.UnstakeUSCRequest) (*proto.UnstakeUSCResponse, error) {
	return h.service.UnstakeUSC(ctx, req)
}
