package product_certificate_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/product_certificate_operations"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service handles product certificate operations business logic
type Service struct {
	repo              *product_certificate_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new product certificate operations service
func NewService(
	repo *product_certificate_operations.Repository,
	cosmosApp *app.USCApp,
	blockchainStorage *storage.StateManager,
	logger *logging.Logger,
	validator *validation.Validator,
	metricsService *metrics.MetricsService,
) *Service {
	return &Service{
		repo:              repo,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		logger:            logger,
		validator:         validator,
		metrics:           metricsService,
	}
}

// CreateProductCertificate creates a new product certificate
func (s *Service) CreateProductCertificate(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("create_product_certificate", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Creating product certificate in business service",
		logging.String("correlation_id", correlationID),
		logging.String("productId", req.ProductId),
		logging.String("from", req.FromAddress))

	// Input validation using validator service
	if err := s.validator.ValidateProductId(req.ProductId); err != nil {
		s.logger.Error("Product ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("product_id", req.ProductId),
			logging.Error(err))
		s.metrics.RecordFailure("create_product_certificate", "validation_error", map[string]string{
			"product_id": req.ProductId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("create_product_certificate", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if req.ProductName == "" {
		s.logger.Error("Product name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("create_product_certificate", "validation_error", map[string]string{
			"product_name": req.ProductName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "product_name is required")
	}

	// ROOT FIX: Repository is single source of truth for certificate creation
	// Business service should NOT create certificate directly, only delegate to repository
	// This avoids duplicate certificate creation with different IDs
	response, err := s.repo.CreateProductCertificate(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create product certificate in repository",
			logging.String("correlation_id", correlationID),
			logging.String("product_id", req.ProductId),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("create_product_certificate", "repository_error", map[string]string{
			"product_id":   req.ProductId,
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to create product certificate: %v", err)
	}

	// Record success metrics
	s.logger.Info("Product certificate created successfully",
		logging.String("correlation_id", correlationID),
		logging.String("product_id", req.ProductId),
		logging.String("from_address", req.FromAddress))
	s.metrics.RecordSuccess("create_product_certificate", map[string]string{
		"product_id":   req.ProductId,
		"from_address": req.FromAddress,
	})

	// Record blockchain-specific metric if certificate was created
	if response != nil && response.CertificateId != "" {
		s.metrics.RecordCertificateCreated(response.CertificateId, req.ProductId)
	}

	return response, nil
}

// VerifyBlockchainProductCertificate verifies a blockchain product certificate
func (s *Service) VerifyBlockchainProductCertificate(ctx context.Context, req *proto.VerifyBlockchainProductCertificateRequest) (*proto.VerifyBlockchainProductCertificateResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("verify_blockchain_product_certificate", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Verifying blockchain product certificate in business service",
		logging.String("correlation_id", correlationID),
		logging.String("certificateId", req.CertificateId),
		logging.String("productId", req.ProductId))

	// Input validation using validator service
	if err := s.validator.ValidateCertificateId(req.CertificateId); err != nil {
		s.logger.Error("Certificate ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("certificate_id", req.CertificateId),
			logging.Error(err))
		s.metrics.RecordFailure("verify_blockchain_product_certificate", "validation_error", map[string]string{
			"certificate_id": req.CertificateId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid certificate_id: %v", err)
	}

	if err := s.validator.ValidateProductId(req.ProductId); err != nil {
		s.logger.Error("Product ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("product_id", req.ProductId),
			logging.Error(err))
		s.metrics.RecordFailure("verify_blockchain_product_certificate", "validation_error", map[string]string{
			"product_id": req.ProductId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.VerifyBlockchainProductCertificate(ctx, req)
	if err != nil {
		s.logger.Error("Failed to verify blockchain product certificate in repository",
			logging.String("correlation_id", correlationID),
			logging.String("certificate_id", req.CertificateId),
			logging.String("product_id", req.ProductId),
			logging.Error(err))
		s.metrics.RecordFailure("verify_blockchain_product_certificate", "repository_error", map[string]string{
			"certificate_id": req.CertificateId,
			"product_id":     req.ProductId,
		})
		return nil, status.Errorf(codes.Internal, "failed to verify blockchain product certificate: %v", err)
	}

	// Record success metrics
	isValidStr := "false"
	if response.IsValid {
		isValidStr = "true"
	}
	s.logger.Info("Blockchain product certificate verified successfully",
		logging.String("correlation_id", correlationID),
		logging.String("certificate_id", req.CertificateId),
		logging.String("product_id", req.ProductId),
		logging.Bool("is_valid", response.IsValid))
	s.metrics.RecordSuccess("verify_blockchain_product_certificate", map[string]string{
		"certificate_id": req.CertificateId,
		"product_id":     req.ProductId,
		"is_valid":       isValidStr,
	})

	return response, nil
}

// TransferProductOwnership transfers product ownership
func (s *Service) TransferProductOwnership(ctx context.Context, req *proto.TransferProductOwnershipRequest) (*proto.TransferProductOwnershipResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("transfer_product_ownership", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Transferring product ownership in business service",
		logging.String("correlation_id", correlationID),
		logging.String("certificateId", req.CertificateId),
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress))

	// Input validation using validator service
	if err := s.validator.ValidateCertificateId(req.CertificateId); err != nil {
		s.logger.Error("Certificate ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("certificate_id", req.CertificateId),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_product_ownership", "validation_error", map[string]string{
			"certificate_id": req.CertificateId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid certificate_id: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_product_ownership", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_product_ownership", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.TransferProductOwnership(ctx, req)
	if err != nil {
		s.logger.Error("Failed to transfer product ownership in repository",
			logging.String("correlation_id", correlationID),
			logging.String("certificate_id", req.CertificateId),
			logging.String("from_address", req.FromAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_product_ownership", "repository_error", map[string]string{
			"certificate_id": req.CertificateId,
			"from_address":  req.FromAddress,
			"to_address":    req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to transfer product ownership: %v", err)
	}

	// Record success metrics
	s.logger.Info("Product ownership transferred successfully",
		logging.String("correlation_id", correlationID),
		logging.String("certificate_id", req.CertificateId),
		logging.String("from_address", req.FromAddress),
		logging.String("to_address", req.ToAddress))
	s.metrics.RecordSuccess("transfer_product_ownership", map[string]string{
		"certificate_id": req.CertificateId,
		"from_address":  req.FromAddress,
		"to_address":    req.ToAddress,
	})

	return response, nil
}
