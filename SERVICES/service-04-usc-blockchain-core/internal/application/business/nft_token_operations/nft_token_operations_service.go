package nft_token_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/nft_token_operations"
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

// Service handles NFT token operations business logic
type Service struct {
	repo              *nft_token_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new NFT token operations service
func NewService(
	repo *nft_token_operations.Repository,
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

// MintNFT creates a new NFT token
func (s *Service) MintNFT(ctx context.Context, req *proto.MintNFTRequest) (*proto.MintNFTResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("mint_nft", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Minting NFT in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("to", req.ToAddress))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_nft", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_nft", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	if req.TokenUri == "" {
		s.logger.Error("Token URI is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("mint_nft", "validation_error", map[string]string{
			"token_uri": req.TokenUri,
		})
		return nil, status.Errorf(codes.InvalidArgument, "token_uri is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.MintNFT(ctx, req)
	if err != nil {
		s.logger.Error("Failed to mint NFT in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_nft", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"to_address":       req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to mint NFT: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT minted successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("to_address", req.ToAddress))
	s.metrics.RecordSuccess("mint_nft", map[string]string{
		"contract_address": req.ContractAddress,
		"to_address":       req.ToAddress,
	})

	// Record blockchain-specific metric if NFT was minted
	if response != nil && response.TokenId != "" {
		s.metrics.RecordNFTMinted(response.TokenId, req.ContractAddress)
	}

	return response, nil
}

// TransferNFT transfers an NFT token
func (s *Service) TransferNFT(ctx context.Context, req *proto.TransferNFTRequest) (*proto.TransferNFTResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("transfer_nft", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Transferring NFT in business service",
		logging.String("correlation_id", correlationID),
		logging.String("tokenId", req.TokenId),
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress))

	// Input validation using validator service
	if err := s.validator.ValidateTokenId(req.TokenId); err != nil {
		s.logger.Error("Token ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_nft", "validation_error", map[string]string{
			"token_id": req.TokenId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid token_id: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_nft", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_nft", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	// Call repository
	response, err := s.repo.TransferNFT(ctx, req)
	if err != nil {
		s.logger.Error("Failed to transfer NFT in repository",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.String("from_address", req.FromAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_nft", "repository_error", map[string]string{
			"token_id":     req.TokenId,
			"from_address": req.FromAddress,
			"to_address":   req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to transfer NFT: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT transferred successfully",
		logging.String("correlation_id", correlationID),
		logging.String("token_id", req.TokenId),
		logging.String("from_address", req.FromAddress),
		logging.String("to_address", req.ToAddress))
	s.metrics.RecordSuccess("transfer_nft", map[string]string{
		"token_id":     req.TokenId,
		"from_address": req.FromAddress,
		"to_address":   req.ToAddress,
	})

	return response, nil
}

// GetNFTInfo retrieves NFT information
func (s *Service) GetNFTInfo(ctx context.Context, req *proto.GetNFTInfoRequest) (*proto.GetNFTInfoResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_nft_info", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting NFT info in business service",
		logging.String("correlation_id", correlationID),
		logging.String("tokenId", req.TokenId),
		logging.String("contract", req.ContractAddress))

	// Input validation using validator service
	if err := s.validator.ValidateTokenId(req.TokenId); err != nil {
		s.logger.Error("Token ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.Error(err))
		s.metrics.RecordFailure("get_nft_info", "validation_error", map[string]string{
			"token_id": req.TokenId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid token_id: %v", err)
	}

	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_nft_info", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetNFTInfo(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get NFT info in repository",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_nft_info", "repository_error", map[string]string{
			"token_id":         req.TokenId,
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get NFT info: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT info retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("token_id", req.TokenId),
		logging.String("contract_address", req.ContractAddress))
	s.metrics.RecordSuccess("get_nft_info", map[string]string{
		"token_id":         req.TokenId,
		"contract_address": req.ContractAddress,
	})

	return response, nil
}

// GetNFTsByOwner retrieves NFTs owned by an address
func (s *Service) GetNFTsByOwner(ctx context.Context, req *proto.GetNFTsByOwnerRequest) (*proto.GetNFTsByOwnerResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_nfts_by_owner", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting NFTs by owner in business service",
		logging.String("correlation_id", correlationID),
		logging.String("owner", req.OwnerAddress),
		logging.Int32("limit", req.Limit),
		logging.Int32("offset", req.Offset))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.OwnerAddress); err != nil {
		s.logger.Error("Owner address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("owner_address", req.OwnerAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_nfts_by_owner", "validation_error", map[string]string{
			"owner_address": req.OwnerAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_address: %v", err)
	}

	// Normalize pagination
	// Use helper function to normalize pagination (reduces duplicate code)
	req.Limit, req.Offset = utils.NormalizePaginationWithDefaults(req.Limit, req.Offset)

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetNFTsByOwner(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get NFTs by owner in repository",
			logging.String("correlation_id", correlationID),
			logging.String("owner_address", req.OwnerAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_nfts_by_owner", "repository_error", map[string]string{
			"owner_address": req.OwnerAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get NFTs by owner: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFTs by owner retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("owner_address", req.OwnerAddress))
	s.metrics.RecordSuccess("get_nfts_by_owner", map[string]string{
		"owner_address": req.OwnerAddress,
	})

	return response, nil
}

// DeployNFTContract deploys a new NFT contract
func (s *Service) DeployNFTContract(ctx context.Context, req *proto.DeployNFTContractRequest) (*proto.DeployNFTContractResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("deploy_nft_contract", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Deploying NFT contract in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("contractName", req.ContractName))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_nft_contract", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if req.ContractName == "" {
		s.logger.Error("Contract name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("deploy_nft_contract", "validation_error", map[string]string{
			"contract_name": req.ContractName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "contract_name is required")
	}

	if req.ContractSymbol == "" {
		s.logger.Error("Contract symbol is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("deploy_nft_contract", "validation_error", map[string]string{
			"contract_symbol": req.ContractSymbol,
		})
		return nil, status.Errorf(codes.InvalidArgument, "contract_symbol is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.DeployNFTContract(ctx, req)
	if err != nil {
		s.logger.Error("Failed to deploy NFT contract in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("contract_name", req.ContractName),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_nft_contract", "repository_error", map[string]string{
			"from_address":  req.FromAddress,
			"contract_name": req.ContractName,
		})
		return nil, status.Errorf(codes.Internal, "failed to deploy NFT contract: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT contract deployed successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("contract_name", req.ContractName))
	s.metrics.RecordSuccess("deploy_nft_contract", map[string]string{
		"from_address":  req.FromAddress,
		"contract_name": req.ContractName,
	})

	// Record blockchain-specific metric if contract was deployed
	if response != nil && response.ContractAddress != "" {
		s.metrics.RecordContractDeployed(response.ContractAddress)
	}

	return response, nil
}

// CreateNFTCollection creates a new NFT collection
func (s *Service) CreateNFTCollection(ctx context.Context, req *proto.CreateNFTCollectionRequest) (*proto.CreateNFTCollectionResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("create_nft_collection", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Creating NFT collection in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("collectionName", req.CollectionName))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("create_nft_collection", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if req.CollectionName == "" {
		s.logger.Error("Collection name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("create_nft_collection", "validation_error", map[string]string{
			"collection_name": req.CollectionName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "collection_name is required")
	}

	if err := s.validator.ValidateWalletAddress(req.CreatorAddress); err != nil {
		s.logger.Error("Creator address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("creator_address", req.CreatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("create_nft_collection", "validation_error", map[string]string{
			"creator_address": req.CreatorAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid creator_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.CreateNFTCollection(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create NFT collection in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("collection_name", req.CollectionName),
			logging.Error(err))
		s.metrics.RecordFailure("create_nft_collection", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"collection_name":  req.CollectionName,
		})
		return nil, status.Errorf(codes.Internal, "failed to create NFT collection: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT collection created successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("collection_name", req.CollectionName))
	s.metrics.RecordSuccess("create_nft_collection", map[string]string{
		"contract_address": req.ContractAddress,
		"collection_name":  req.CollectionName,
	})

	return response, nil
}

// BurnNFT burns an NFT token
func (s *Service) BurnNFT(ctx context.Context, req *proto.BurnNFTRequest) (*proto.BurnNFTResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("burn_nft", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Burning NFT in business service",
		logging.String("correlation_id", correlationID),
		logging.String("tokenId", req.TokenId),
		logging.String("owner", req.OwnerAddress))

	// Input validation using validator service
	if err := s.validator.ValidateTokenId(req.TokenId); err != nil {
		s.logger.Error("Token ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.Error(err))
		s.metrics.RecordFailure("burn_nft", "validation_error", map[string]string{
			"token_id": req.TokenId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid token_id: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.OwnerAddress); err != nil {
		s.logger.Error("Owner address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("owner_address", req.OwnerAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_nft", "validation_error", map[string]string{
			"owner_address": req.OwnerAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_address: %v", err)
	}

	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_nft", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.BurnNFT(ctx, req)
	if err != nil {
		s.logger.Error("Failed to burn NFT in repository",
			logging.String("correlation_id", correlationID),
			logging.String("token_id", req.TokenId),
			logging.String("owner_address", req.OwnerAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_nft", "repository_error", map[string]string{
			"token_id":         req.TokenId,
			"owner_address":    req.OwnerAddress,
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to burn NFT: %v", err)
	}

	// Record success metrics
	s.logger.Info("NFT burned successfully",
		logging.String("correlation_id", correlationID),
		logging.String("token_id", req.TokenId),
		logging.String("owner_address", req.OwnerAddress))
	s.metrics.RecordSuccess("burn_nft", map[string]string{
		"token_id":         req.TokenId,
		"owner_address":    req.OwnerAddress,
		"contract_address": req.ContractAddress,
	})

	return response, nil
}
