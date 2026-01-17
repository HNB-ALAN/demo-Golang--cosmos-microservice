package product_certificate_operations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	repoerrors "service-04/internal/application/repository"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/database"
	proto "service-04/proto"

	// Cosmos SDK imports
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	pctypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/types"

	"github.com/usc-platform/shared/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Repository handles product certificate operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new product certificate operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// CreateProductCertificate creates a new product certificate
func (r *Repository) CreateProductCertificate(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		r.logger.Info("Attempting to create certificate on keeper",
			logging.String("product_id", req.ProductId))
		result, err := r.createCertificateOnKeeper(ctx, req)
		if err == nil {
			r.logger.Info("Certificate created on keeper, attempting to save to database",
				logging.String("certificate_id", result.CertificateId),
				logging.Bool("db_available", r.db != nil))
			// Save to PostgreSQL for analytics (sync to ensure it's saved)
			if r.db != nil {
				r.logger.Info("Saving certificate to database",
					logging.String("certificate_id", result.CertificateId),
					logging.String("product_id", req.ProductId))
				if err := r.saveCertificateToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save certificate to database",
						logging.String("certificate_id", result.CertificateId),
						logging.String("product_id", req.ProductId),
						logging.Error(err))
				} else {
					r.logger.Info("Certificate saved to database successfully",
						logging.String("certificate_id", result.CertificateId),
						logging.String("product_id", req.ProductId))
				}
			} else {
				r.logger.Warn("Database connection not available, skipping certificate save",
					logging.String("certificate_id", result.CertificateId))
			}
			return result, nil
		} else {
			r.logger.Warn("Failed to create certificate on keeper, falling back to database",
				logging.String("product_id", req.ProductId),
				logging.Error(err))
		}
	}

	// Priority 2: PostgreSQL (fallback)
	r.logger.Info("Using database fallback for certificate creation",
		logging.String("product_id", req.ProductId))
	return r.createCertificateInDatabase(ctx, req)
}

// VerifyBlockchainProductCertificate verifies a blockchain product certificate
func (r *Repository) VerifyBlockchainProductCertificate(ctx context.Context, req *proto.VerifyBlockchainProductCertificateRequest) (*proto.VerifyBlockchainProductCertificateResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		// Use recover to catch any panic from protobuf unmarshal
		var result *proto.VerifyBlockchainProductCertificateResponse
		var err error
		func() {
			defer func() {
				if p := recover(); p != nil {
					r.logger.Debug("Panic recovered in verifyCertificateFromKeeper, will fallback to database",
						logging.String("certificate_id", req.CertificateId),
						logging.String("panic", fmt.Sprintf("%v", p)))
					err = fmt.Errorf("panic in keeper verification: %v", p)
				}
			}()
			result, err = r.verifyCertificateFromKeeper(ctx, req)
		}()

		if err == nil && result != nil {
			return result, nil
		}

		r.logger.Debug("Certificate not found in keeper or panic occurred, using database fallback",
			logging.String("certificate_id", req.CertificateId))
	}

	// Priority 2: PostgreSQL (fallback)
	return r.verifyCertificateInDatabase(ctx, req)
}

// TransferProductOwnership transfers product ownership
func (r *Repository) TransferProductOwnership(ctx context.Context, req *proto.TransferProductOwnershipRequest) (*proto.TransferProductOwnershipResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.transferOwnershipOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveOwnershipTransferToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save ownership transfer to database",
						logging.String("certificate_id", req.CertificateId),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Ownership transfer saved to database successfully",
						logging.String("certificate_id", req.CertificateId))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.transferOwnershipInDatabase(ctx, req)
}

// Helper methods for ProductCertificateKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// createCertificateOnKeeper creates a certificate on the keeper
// ROOT FIX: Use GetSDKContextForWrite to ensure certificate is committed to RocksDB
func (r *Repository) createCertificateOnKeeper(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
	// Use writable context for write operations
	sdkCtx, err := utils.GetSDKContextForWrite(ctx, r.cosmosApp, r.logger)
	if err != nil {
		return nil, err
	}

	// Generate certificate ID
	certificateID := fmt.Sprintf("cert_%s_%d", req.ProductId, time.Now().Unix())

	// Create ProductCertificate
	cert := pctypes.ProductCertificate{
		ID:         certificateID,
		ProductID:  req.ProductId,
		Owner:      req.FromAddress,
		Status:     "active",
		Metadata:   req.ProductDescription,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		VerifiedAt: 0,
		ExpiresAt:  0, // No expiration by default
	}

	// Set certificate in keeper
	if err := r.cosmosApp.ProductCertificateKeeper.SetCertificate(sdkCtx, cert); err != nil {
		return nil, repoerrors.NewDatabaseError("set_certificate", err)
	}

	// Generate transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, "", "", "create_product_certificate", certificateID, "")

	// Generate certificate hash
	dataStr := fmt.Sprintf("%s:%s:%s", certificateID, req.ProductId, req.FromAddress)
	hashBytes := sha256.Sum256([]byte(dataStr))
	certHash := hex.EncodeToString(hashBytes[:16])

	return &proto.CreateProductCertificateResponse{
		CertificateId:   certificateID,
		TransactionHash: txHash,
		Status:          1, // Confirmed
		ErrorMessage:    "",
		CertificateHash: certHash,
	}, nil
}

// verifyCertificateFromKeeper verifies a certificate from the keeper
func (r *Repository) verifyCertificateFromKeeper(ctx context.Context, req *proto.VerifyBlockchainProductCertificateRequest) (*proto.VerifyBlockchainProductCertificateResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get certificate from keeper
	// Use recover to catch any panic from protobuf unmarshal
	var cert pctypes.ProductCertificate
	var found bool
	func() {
		defer func() {
			if p := recover(); p != nil {
				r.logger.Debug("Panic recovered in GetCertificate, certificate may not exist in keeper",
					logging.String("certificate_id", req.CertificateId),
					logging.String("panic", fmt.Sprintf("%v", p)))
				found = false
			}
		}()
		cert, found = r.cosmosApp.ProductCertificateKeeper.GetCertificate(sdkCtx, req.CertificateId)
	}()

	if !found {
		// Certificate not found in keeper, return error to trigger database fallback
		return nil, repoerrors.NewNotFoundError("certificate", req.CertificateId)
	}

	// Check if certificate matches product ID
	isValid := cert.ProductID == req.ProductId
	isAuthentic := cert.Status == "active"

	// Get verification record if exists
	verifierAddress := req.CurrentOwnerAddress
	if verifierAddress != "" {
		verification, found := r.cosmosApp.ProductCertificateKeeper.GetVerification(sdkCtx, req.CertificateId, verifierAddress)
		if found && verification.Status == "verified" {
			isValid = true
			isAuthentic = true
		}
	}

	return &proto.VerifyBlockchainProductCertificateResponse{
		IsValid:              isValid,
		IsAuthentic:          isAuthentic,
		VerificationResult:   "Valid",
		CertificateStatus:    cert.Status,
		OriginalManufacturer: cert.Owner,
		CurrentOwner:         cert.Owner,
		ManufacturingDate:    time.Unix(cert.CreatedAt, 0).Format("2006-01-02"),
		ExpirationDate:       time.Unix(cert.ExpiresAt, 0).Format("2006-01-02"),
		OwnershipHistory:     []string{cert.Owner},
		ProductMetadata:      cert.Metadata,
	}, nil
}

// transferOwnershipOnKeeper transfers ownership on the keeper
// ROOT FIX: Use GetSDKContextForWrite to ensure ownership transfer is committed to RocksDB
func (r *Repository) transferOwnershipOnKeeper(ctx context.Context, req *proto.TransferProductOwnershipRequest) (*proto.TransferProductOwnershipResponse, error) {
	// Use writable context for write operations
	sdkCtx, err := utils.GetSDKContextForWrite(ctx, r.cosmosApp, r.logger)
	if err != nil {
		return nil, err
	}

	// Get certificate from keeper
	cert, found := r.cosmosApp.ProductCertificateKeeper.GetCertificate(sdkCtx, req.CertificateId)
	if !found {
		return nil, repoerrors.NewNotFoundError("certificate", req.CertificateId)
	}

	// Verify current owner
	if cert.Owner != req.FromAddress {
		return nil, repoerrors.NewValidationError("from_address", fmt.Sprintf("current owner mismatch: expected %s, got %s", req.FromAddress, cert.Owner))
	}

	// Update certificate owner
	cert.Owner = req.ToAddress
	cert.UpdatedAt = time.Now().Unix()

	// Set certificate in keeper
	if err := r.cosmosApp.ProductCertificateKeeper.SetCertificate(sdkCtx, cert); err != nil {
		return nil, repoerrors.NewDatabaseError("update_certificate", err)
	}

	// Generate real transaction hash using blocktypes helper
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, req.ToAddress, "", "transfer_product_ownership", req.CertificateId, "")

	// Create timestamp for transferred_at
	transferredAt := time.Now()

	return &proto.TransferProductOwnershipResponse{
		TransactionHash: txHash,
		Status:          1, // Confirmed
		ErrorMessage:    "",
		TransferredAt:   timestamppb.New(transferredAt),
		GasUsed:         "0",
		NewOwner:        req.ToAddress,
		CertificateId:   req.CertificateId,
	}, nil
}

// Database fallback methods

// createCertificateInDatabase creates a certificate in database
func (r *Repository) createCertificateInDatabase(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate certificate ID
	certificateID := fmt.Sprintf("cert_%s_%d", req.ProductId, time.Now().Unix())

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:cert", req.FromAddress, req.ProductId, req.SerialNumber, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])
	certHash := hex.EncodeToString(hashBytes[:16])

	// Use product_name from request, fallback to product_id if empty
	productName := req.ProductName
	if productName == "" {
		productName = req.ProductId
	}

	// Use manufacturer_address from request, fallback to from_address if empty
	manufacturerAddr := req.ManufacturerAddress
	if manufacturerAddr == "" {
		manufacturerAddr = req.FromAddress
	}

	query := `
		INSERT INTO product_certificates (
			certificate_id, product_id, product_name, manufacturer_address,
			current_owner_address, deployment_transaction_hash, status, created_at, product_metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, TO_TIMESTAMP($8), $9)
		ON CONFLICT (certificate_id) DO NOTHING
	`

	_, err := postgres.ExecContext(ctx, query,
		certificateID, req.ProductId, productName, manufacturerAddr,
		req.FromAddress, txHash, "active", time.Now().Unix(), req.ProductDescription,
	)
	if err != nil {
		return nil, err
	}

	return &proto.CreateProductCertificateResponse{
		CertificateId:   certificateID,
		TransactionHash: txHash,
		Status:          0, // Pending (database fallback)
		ErrorMessage:    "",
		CertificateHash: certHash,
	}, nil
}

// saveCertificateToDatabase saves certificate to database for analytics
func (r *Repository) saveCertificateToDatabase(ctx context.Context, req *proto.CreateProductCertificateRequest, result *proto.CreateProductCertificateResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		r.logger.Warn("PostgreSQL connection not available for saving certificate",
			logging.String("certificate_id", result.CertificateId))
		return fmt.Errorf("postgres connection not available")
	}

	// Use product_name from request, fallback to product_id if empty
	productName := req.ProductName
	if productName == "" {
		productName = req.ProductId
	}

	// Use manufacturer_address from request, fallback to from_address if empty
	manufacturerAddr := req.ManufacturerAddress
	if manufacturerAddr == "" {
		manufacturerAddr = req.FromAddress
	}

	// Extract transaction hash from result
	txHash := result.TransactionHash
	if txHash == "" {
		// Generate fallback hash
		dataStr := fmt.Sprintf("%s:%s:%s", result.CertificateId, req.ProductId, time.Now().Format(time.RFC3339))
		hashBytes := sha256.Sum256([]byte(dataStr))
		txHash = "0x" + hex.EncodeToString(hashBytes[:16])
	}

	query := `
		INSERT INTO product_certificates (
			certificate_id, product_id, product_name, manufacturer_address,
			current_owner_address, deployment_transaction_hash, status, created_at, product_metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, TO_TIMESTAMP($8), $9)
		ON CONFLICT (certificate_id) DO UPDATE SET
			status = EXCLUDED.status,
			current_owner_address = EXCLUDED.current_owner_address,
			updated_at = NOW()
	`

	_, err := postgres.ExecContext(ctx, query,
		result.CertificateId, req.ProductId, productName, manufacturerAddr,
		req.FromAddress, txHash, "active", time.Now().Unix(), req.ProductDescription,
	)
	return err
}

// verifyCertificateInDatabase verifies a certificate in database
func (r *Repository) verifyCertificateInDatabase(ctx context.Context, req *proto.VerifyBlockchainProductCertificateRequest) (*proto.VerifyBlockchainProductCertificateResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.VerifyBlockchainProductCertificateResponse{
			IsValid:            false,
			IsAuthentic:        false,
			VerificationResult: "Database not available",
			CertificateStatus:  "unknown",
		}, nil
	}

	query := `
		SELECT certificate_id, product_id, current_owner_address, status, 
		       EXTRACT(EPOCH FROM created_at)::BIGINT as created_at,
		       COALESCE(expires_at, 0) as expires_at,
		       COALESCE(product_metadata, '') as product_metadata
		FROM product_certificates
		WHERE certificate_id = $1
		LIMIT 1
	`

	var certID, productID, ownerAddr, status, metadata string
	var createdAt, expiresAt int64

	err := postgres.QueryRowContext(ctx, query, req.CertificateId).Scan(
		&certID, &productID, &ownerAddr, &status, &createdAt, &expiresAt, &metadata,
	)
	if err != nil {
		return &proto.VerifyBlockchainProductCertificateResponse{
			IsValid:            false,
			IsAuthentic:        false,
			VerificationResult: "Certificate not found",
			CertificateStatus:  "not_found",
		}, nil
	}

	isValid := productID == req.ProductId
	isAuthentic := status == "active"

	manufacturingDate := "1970-01-01"
	expirationDate := "1970-01-01"
	if createdAt > 0 {
		manufacturingDate = time.Unix(createdAt, 0).Format("2006-01-02")
	}
	if expiresAt > 0 {
		expirationDate = time.Unix(expiresAt, 0).Format("2006-01-02")
	}

	return &proto.VerifyBlockchainProductCertificateResponse{
		IsValid:              isValid,
		IsAuthentic:          isAuthentic,
		VerificationResult:   "Valid",
		CertificateStatus:    status,
		OriginalManufacturer: ownerAddr,
		CurrentOwner:         ownerAddr,
		ManufacturingDate:    manufacturingDate,
		ExpirationDate:       expirationDate,
		OwnershipHistory:     []string{ownerAddr},
		ProductMetadata:      metadata,
	}, nil
}

// transferOwnershipInDatabase transfers ownership in database
func (r *Repository) transferOwnershipInDatabase(ctx context.Context, req *proto.TransferProductOwnershipRequest) (*proto.TransferProductOwnershipResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		r.logger.Warn("PostgreSQL connection not available for transfer ownership",
			logging.String("certificate_id", req.CertificateId))
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	r.logger.Info("Transferring ownership in database",
		logging.String("certificate_id", req.CertificateId),
		logging.String("from_address", req.FromAddress),
		logging.String("to_address", req.ToAddress))

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:transfer", req.CertificateId, req.FromAddress, req.ToAddress, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	query := `
		UPDATE product_certificates
		SET current_owner_address = $1, updated_at = TO_TIMESTAMP($2)
		WHERE certificate_id = $3 AND current_owner_address = $4
	`

	result, err := postgres.ExecContext(ctx, query, req.ToAddress, time.Now().Unix(), req.CertificateId, req.FromAddress)
	if err != nil {
		r.logger.Error("Failed to execute ownership transfer query",
			logging.String("certificate_id", req.CertificateId),
			logging.Error(err))
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn("No rows affected in ownership transfer",
			logging.String("certificate_id", req.CertificateId),
			logging.String("from_address", req.FromAddress))
		return nil, fmt.Errorf("certificate not found or ownership mismatch")
	}

	r.logger.Info("Ownership transferred successfully in database",
		logging.String("certificate_id", req.CertificateId),
		logging.String("new_owner", req.ToAddress))

	return &proto.TransferProductOwnershipResponse{
		TransactionHash: txHash,
		Status:          0, // Pending (database fallback)
		ErrorMessage:    "",
		NewOwner:        req.ToAddress,
		CertificateId:   req.CertificateId,
	}, nil
}

// saveOwnershipTransferToDatabase saves ownership transfer to database for analytics
func (r *Repository) saveOwnershipTransferToDatabase(ctx context.Context, req *proto.TransferProductOwnershipRequest, result *proto.TransferProductOwnershipResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	// Update certificate owner
	updateQuery := `
		UPDATE product_certificates
		SET current_owner_address = $1, updated_at = TO_TIMESTAMP($2)
		WHERE certificate_id = $3
	`

	if _, err := postgres.ExecContext(ctx, updateQuery, req.ToAddress, time.Now().Unix(), req.CertificateId); err != nil {
		return fmt.Errorf("failed to update certificate owner: %w", err)
	}

	// Save ownership history
	historyQuery := `
		INSERT INTO product_certificate_ownership_history (
			certificate_id, from_address, to_address, transaction_hash, transferred_at
		) VALUES ($1, $2, $3, $4, TO_TIMESTAMP($5))
		ON CONFLICT (certificate_id, transaction_hash) DO NOTHING
	`

	if _, err := postgres.ExecContext(ctx, historyQuery,
		req.CertificateId,
		req.FromAddress,
		req.ToAddress,
		result.TransactionHash,
		time.Now().Unix(),
	); err != nil {
		return fmt.Errorf("failed to save ownership history: %w", err)
	}

	return nil
}
