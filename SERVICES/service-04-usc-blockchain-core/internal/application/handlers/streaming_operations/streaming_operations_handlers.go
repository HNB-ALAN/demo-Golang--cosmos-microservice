package streaming_operations

import (
	"service-04/internal/application/business/streaming_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles streaming operations gRPC requests
type Handlers struct {
	proto.UnimplementedStreamingOperationsServiceServer
	service *streaming_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new streaming operations handlers
func NewHandlers(service *streaming_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// StreamBlocks streams blockchain blocks
func (h *Handlers) StreamBlocks(req *proto.StreamBlocksRequest, stream proto.StreamingOperationsService_StreamBlocksServer) error {
	return h.service.StreamBlocks(req, stream)
}

// StreamTransactions streams blockchain transactions
func (h *Handlers) StreamTransactions(req *proto.StreamTransactionsRequest, stream proto.StreamingOperationsService_StreamTransactionsServer) error {
	return h.service.StreamTransactions(req, stream)
}

// StreamValidatorEvents streams validator events
func (h *Handlers) StreamValidatorEvents(req *proto.StreamValidatorEventsRequest, stream proto.StreamingOperationsService_StreamValidatorEventsServer) error {
	return h.service.StreamValidatorEvents(req, stream)
}

// StreamNetworkEvents streams network events
func (h *Handlers) StreamNetworkEvents(req *proto.StreamNetworkEventsRequest, stream proto.StreamingOperationsService_StreamNetworkEventsServer) error {
	return h.service.StreamNetworkEvents(req, stream)
}
