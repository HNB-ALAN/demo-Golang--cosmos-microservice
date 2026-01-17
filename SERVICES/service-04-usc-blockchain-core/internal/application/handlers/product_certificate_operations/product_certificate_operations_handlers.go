package product_certificate_operations

import (
	"context"

	"service-04/internal/application/business/product_certificate_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles product certificate operations gRPC requests
type Handlers struct {
	proto.UnimplementedProductCertificateOperationsServiceServer
	service *product_certificate_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new product certificate operations handlers
func NewHandlers(service *product_certificate_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// CreateProductCertificate creates a new product certificate
func (h *Handlers) CreateProductCertificate(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
	return h.service.CreateProductCertificate(ctx, req)
}

// VerifyBlockchainProductCertificate verifies a blockchain product certificate
func (h *Handlers) VerifyBlockchainProductCertificate(ctx context.Context, req *proto.VerifyBlockchainProductCertificateRequest) (*proto.VerifyBlockchainProductCertificateResponse, error) {
	return h.service.VerifyBlockchainProductCertificate(ctx, req)
}

// TransferProductOwnership transfers product ownership
func (h *Handlers) TransferProductOwnership(ctx context.Context, req *proto.TransferProductOwnershipRequest) (*proto.TransferProductOwnershipResponse, error) {
	return h.service.TransferProductOwnership(ctx, req)
}
