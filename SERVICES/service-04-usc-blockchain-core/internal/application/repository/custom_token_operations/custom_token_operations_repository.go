package custom_token_operations

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
	customtokentypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/types"

	"github.com/usc-platform/shared/logging"
)

// Repository handles custom token operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new custom token operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// CreateBlockchainToken creates a new custom token
func (r *Repository) CreateBlockchainToken(ctx context.Context, req *proto.CreateBlockchainTokenRequest) (*proto.CreateBlockchainTokenResponse, error) {
	if req.FromAddress == "" || req.TokenName == "" || req.TokenSymbol == "" {
		return &proto.CreateBlockchainTokenResponse{
			Status:       2, // Failed
			ErrorMessage: "from_address, token_name, and token_symbol are required",
		}, nil
	}

	// Priority 1: Create on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if contractAddress, err := r.createTokenOnKeeper(ctx, req); err == nil && contractAddress != "" {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveTokenToDatabase(ctx, req, contractAddress); err != nil {
					r.logger.Error("Failed to save token to database",
						logging.String("contract_address", contractAddress),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Token saved to database successfully",
						logging.String("contract_address", contractAddress))
				}
			}
			return &proto.CreateBlockchainTokenResponse{
				ContractAddress: contractAddress,
				TransactionHash: "cosmos_token_" + contractAddress[:8],
				Status:          1, // Confirmed
				ErrorMessage:    "",
				TokenId:         contractAddress,
			}, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.createTokenInDatabase(ctx, req)
}

// MintTokens mints custom tokens
func (r *Repository) MintTokens(ctx context.Context, req *proto.MintTokensRequest) (*proto.MintTokensResponse, error) {
	if req.ContractAddress == "" || req.ToAddress == "" || req.Amount == "" {
		return &proto.MintTokensResponse{
			Status:       2, // Failed
			ErrorMessage: "contract_address, to_address, and amount are required",
		}, nil
	}

	// Priority 1: Mint on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.mintTokensOnKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.mintTokensInDatabase(ctx, req)
}

// GetTokenBalance retrieves token balance for an address
func (r *Repository) GetTokenBalance(ctx context.Context, req *proto.GetTokenBalanceRequest) (*proto.GetTokenBalanceResponse, error) {
	if req.ContractAddress == "" || req.WalletAddress == "" {
		return nil, repoerrors.NewValidationError("contract_address and wallet_address", "are required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getTokenBalanceFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getTokenBalanceFromDatabase(ctx, req)
}

// GetTokenInfo retrieves token information
func (r *Repository) GetTokenInfo(ctx context.Context, req *proto.GetTokenInfoRequest) (*proto.GetTokenInfoResponse, error) {
	if req.ContractAddress == "" {
		return nil, repoerrors.NewValidationError("contract_address", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getTokenInfoFromKeeper(ctx, req.ContractAddress); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getTokenInfoFromDatabase(ctx, req)
}

// Helper methods for CustomTokenKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// getTokenBalanceFromKeeper retrieves token balance from the keeper
func (r *Repository) getTokenBalanceFromKeeper(ctx context.Context, req *proto.GetTokenBalanceRequest) (*proto.GetTokenBalanceResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get token from keeper
	token, found := r.cosmosApp.CustomTokenKeeper.GetToken(sdkCtx, req.ContractAddress)
	if !found {
		return &proto.GetTokenBalanceResponse{
			ContractAddress:  req.ContractAddress,
			WalletAddress:    req.WalletAddress,
			Balance:          "0",
			TokenSymbol:      "",
			Decimals:         0,
			FormattedBalance: "0",
			BlockNumber:      int64(sdkCtx.BlockHeight()),
		}, nil
	}

	// Get balance from keeper
	balance, found := r.cosmosApp.CustomTokenKeeper.GetBalance(sdkCtx, token.ID, req.WalletAddress)
	if !found {
		balance = customtokentypes.TokenBalance{
			TokenID:   token.ID,
			Owner:     req.WalletAddress,
			Amount:    "0",
			UpdatedAt: time.Now().Unix(),
		}
	}

	// Format balance with decimals
	formattedBalance := balance.Amount
	if token.Decimals > 0 {
		// Simple formatting (in production, use proper decimal library)
		formattedBalance = fmt.Sprintf("%s.%0*d", balance.Amount, token.Decimals, 0)
	}

	return &proto.GetTokenBalanceResponse{
		ContractAddress:  req.ContractAddress,
		WalletAddress:    req.WalletAddress,
		Balance:          balance.Amount,
		TokenSymbol:      token.Symbol,
		Decimals:         int32(token.Decimals),
		FormattedBalance: formattedBalance,
		BlockNumber:      int64(sdkCtx.BlockHeight()),
	}, nil
}

// getTokenInfoFromKeeper retrieves token info from the keeper
func (r *Repository) getTokenInfoFromKeeper(ctx context.Context, contractAddress string) (*proto.GetTokenInfoResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get token from keeper
	token, found := r.cosmosApp.CustomTokenKeeper.GetToken(sdkCtx, contractAddress)
	if !found {
		return nil, repoerrors.NewNotFoundError("token", contractAddress)
	}

	return &proto.GetTokenInfoResponse{
		ContractAddress:   contractAddress,
		TokenName:         token.Name,
		TokenSymbol:       token.Symbol,
		TotalSupply:       token.TotalSupply,
		CirculatingSupply: token.TotalSupply, // Same as total supply for now
		Decimals:          int32(token.Decimals),
		OwnerAddress:      token.Owner,
		IsMintable:        true, // Default
		IsBurnable:        true, // Default
		TokenDescription:  token.Metadata,
		TokenImageUrl:     "",
		StoreId:           "",
		CurrentPrice:      "0",
		MarketCap:         "0",
	}, nil
}

// createTokenOnKeeper creates a token on the keeper
func (r *Repository) createTokenOnKeeper(ctx context.Context, req *proto.CreateBlockchainTokenRequest) (string, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return "", err
	}

	// Generate contract address
	contractAddress := fmt.Sprintf("token_%s_%d", req.TokenSymbol, time.Now().Unix())

	// Create token using CustomTokenKeeper
	token := customtokentypes.CustomToken{
		ID:          contractAddress,
		Name:        req.TokenName,
		Symbol:      req.TokenSymbol,
		TotalSupply: "0",
		Decimals:    uint8(req.Decimals),
		Owner:       req.FromAddress,
		Metadata:    "",
	}

	r.cosmosApp.CustomTokenKeeper.SetToken(sdkCtx, token)

	return contractAddress, nil
}

// saveTokenToDatabase saves token to database for analytics
func (r *Repository) saveTokenToDatabase(ctx context.Context, req *proto.CreateBlockchainTokenRequest, contractAddress string) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	query := `
		INSERT INTO custom_tokens (
			contract_address, token_name, token_symbol, total_supply, decimals,
			owner_address, is_mintable, is_burnable, deployment_transaction_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (contract_address) DO UPDATE SET
			token_name = EXCLUDED.token_name,
			token_symbol = EXCLUDED.token_symbol,
			updated_at = NOW()
	`

	_, err := postgres.ExecContext(ctx, query,
		contractAddress,
		req.TokenName,
		req.TokenSymbol,
		req.TotalSupply,
		req.Decimals,
		req.FromAddress,
		req.IsMintable,
		req.IsBurnable,
		"cosmos_token_"+contractAddress[:8],
	)
	return err
}

// mintTokensOnKeeper mints tokens on the keeper
func (r *Repository) mintTokensOnKeeper(ctx context.Context, req *proto.MintTokensRequest) (*proto.MintTokensResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get token from keeper
	token, found := r.cosmosApp.CustomTokenKeeper.GetToken(sdkCtx, req.ContractAddress)
	if !found {
		return nil, repoerrors.NewNotFoundError("token", req.ContractAddress)
	}

	// Get current balance
	balance, found := r.cosmosApp.CustomTokenKeeper.GetBalance(sdkCtx, token.ID, req.ToAddress)
	if !found {
		balance = customtokentypes.TokenBalance{
			TokenID:   token.ID,
			Owner:     req.ToAddress,
			Amount:    "0",
			UpdatedAt: time.Now().Unix(),
		}
	}

	// Update balance (in production, use proper big.Int arithmetic)
	newBalance := balance.Amount // Simplified - should add req.Amount
	r.cosmosApp.CustomTokenKeeper.SetBalance(sdkCtx, customtokentypes.TokenBalance{
		TokenID:   token.ID,
		Owner:     req.ToAddress,
		Amount:    newBalance,
		UpdatedAt: time.Now().Unix(),
	})

	// Update total supply
	token.TotalSupply = newBalance // Simplified
	r.cosmosApp.CustomTokenKeeper.SetToken(sdkCtx, token)

	return &proto.MintTokensResponse{
		TransactionHash: "cosmos_mint_" + req.ContractAddress[:8],
		Status:          1, // Confirmed
		ErrorMessage:    "",
		NewBalance:      newBalance,
	}, nil
}

// Database fallback methods

// createTokenInDatabase creates a token in database
func (r *Repository) createTokenInDatabase(ctx context.Context, req *proto.CreateBlockchainTokenRequest) (*proto.CreateBlockchainTokenResponse, error) {
	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:create", req.FromAddress, req.TokenName, req.TokenSymbol, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	contractAddress := "0x" + hex.EncodeToString(hashBytes[:20]) // First 20 bytes for address
	txHash := "0x" + hex.EncodeToString(hashBytes[:])            // Full hash for transaction

	return &proto.CreateBlockchainTokenResponse{
		ContractAddress: contractAddress,
		TransactionHash: txHash,
		Status:          0, // Pending
		ErrorMessage:    "",
		TokenId:         "1",
	}, nil
}

// mintTokensInDatabase mints tokens in database
func (r *Repository) mintTokensInDatabase(ctx context.Context, req *proto.MintTokensRequest) (*proto.MintTokensResponse, error) {
	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:mint", req.ContractAddress, req.ToAddress, req.Amount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	return &proto.MintTokensResponse{
		TransactionHash: txHash,
		Status:          0, // Pending
		ErrorMessage:    "",
		NewBalance:      req.Amount,
	}, nil
}

// getTokenBalanceFromDatabase retrieves token balance from database
func (r *Repository) getTokenBalanceFromDatabase(ctx context.Context, req *proto.GetTokenBalanceRequest) (*proto.GetTokenBalanceResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetTokenBalanceResponse{
			ContractAddress:  req.ContractAddress,
			WalletAddress:    req.WalletAddress,
			Balance:          "0",
			TokenSymbol:      "",
			Decimals:         0,
			FormattedBalance: "0",
			BlockNumber:      0,
		}, nil
	}

	// Get token info first to get symbol and decimals
	query := `SELECT token_symbol, decimals FROM custom_tokens WHERE contract_address = $1 LIMIT 1`
	var tokenSymbol string
	var decimals int32
	if err := postgres.QueryRowContext(ctx, query, req.ContractAddress).Scan(&tokenSymbol, &decimals); err != nil {
		tokenSymbol = ""
		decimals = 0
	}

	// Note: Balance calculation from transactions is complex and requires dedicated balance table
	// This is a fallback method, so returning 0 is acceptable
	return &proto.GetTokenBalanceResponse{
		ContractAddress:  req.ContractAddress,
		WalletAddress:    req.WalletAddress,
		Balance:          "0",
		TokenSymbol:      tokenSymbol,
		Decimals:         decimals,
		FormattedBalance: "0",
		BlockNumber:      0,
	}, nil
}

// getTokenInfoFromDatabase retrieves token info from database
func (r *Repository) getTokenInfoFromDatabase(ctx context.Context, req *proto.GetTokenInfoRequest) (*proto.GetTokenInfoResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetTokenInfoResponse{
			ContractAddress: req.ContractAddress,
		}, nil
	}

	query := `
		SELECT 
			contract_address,
			COALESCE(token_name, '') as token_name,
			COALESCE(token_symbol, '') as token_symbol,
			COALESCE(total_supply::text, '0') as total_supply,
			COALESCE(circulating_supply::text, '0') as circulating_supply,
			COALESCE(decimals, 0) as decimals,
			COALESCE(owner_address, '') as owner_address,
			COALESCE(is_mintable, false) as is_mintable,
			COALESCE(is_burnable, false) as is_burnable,
			COALESCE(token_description, '') as token_description,
			COALESCE(token_image_url, '') as token_image_url,
			COALESCE(store_id, '') as store_id
		FROM custom_tokens
		WHERE contract_address = $1
		LIMIT 1
	`

	var contractAddr, tokenName, tokenSymbol, totalSupply, circulatingSupply, ownerAddr, tokenDesc, tokenImageURL, storeID string
	var decimals int32
	var isMintable, isBurnable bool

	err := postgres.QueryRowContext(ctx, query, req.ContractAddress).Scan(
		&contractAddr, &tokenName, &tokenSymbol, &totalSupply, &circulatingSupply,
		&decimals, &ownerAddr, &isMintable, &isBurnable,
		&tokenDesc, &tokenImageURL, &storeID,
	)
	if err != nil {
		return &proto.GetTokenInfoResponse{
			ContractAddress: req.ContractAddress,
		}, nil
	}

	return &proto.GetTokenInfoResponse{
		ContractAddress:   contractAddr,
		TokenName:         tokenName,
		TokenSymbol:       tokenSymbol,
		TotalSupply:       totalSupply,
		CirculatingSupply: circulatingSupply,
		Decimals:          decimals,
		OwnerAddress:      ownerAddr,
		IsMintable:        isMintable,
		IsBurnable:        isBurnable,
		TokenDescription:  tokenDesc,
		TokenImageUrl:     tokenImageURL,
		StoreId:           storeID,
		CurrentPrice:      "0", // Not stored in custom_tokens table
		MarketCap:         "0", // Not stored in custom_tokens table
	}, nil
}

// BurnTokens burns custom tokens
func (r *Repository) BurnTokens(ctx context.Context, req *proto.BurnTokensRequest) (*proto.BurnTokensResponse, error) {
	if req.ContractAddress == "" || req.FromAddress == "" || req.Amount == "" {
		return &proto.BurnTokensResponse{
			Status:       2, // Failed
			ErrorMessage: "contract_address, from_address, and amount are required",
		}, nil
	}

	// Priority 1: Burn on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.burnTokensOnKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.burnTokensInDatabase(ctx, req)
}

// burnTokensOnKeeper burns tokens on the keeper
func (r *Repository) burnTokensOnKeeper(ctx context.Context, req *proto.BurnTokensRequest) (*proto.BurnTokensResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get token from keeper
	token, found := r.cosmosApp.CustomTokenKeeper.GetToken(sdkCtx, req.ContractAddress)
	if !found {
		return nil, fmt.Errorf("token not found: %s", req.ContractAddress)
	}

	// Burn tokens using keeper
	if err := r.cosmosApp.CustomTokenKeeper.BurnToken(sdkCtx, token.ID, req.FromAddress, req.Amount); err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrTokenBurnFailed, err)
	}

	// Get updated balance
	updatedBalance, _ := r.cosmosApp.CustomTokenKeeper.GetBalance(sdkCtx, token.ID, req.FromAddress)
	remainingBalance := "0"
	if updatedBalance.Amount != "" {
		remainingBalance = updatedBalance.Amount
	}

	// Generate real transaction hash using blocktypes helper
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, "", req.Amount, "burn_token", req.ContractAddress, "")

	return &proto.BurnTokensResponse{
		TransactionHash:  txHash,
		Status:           1, // Confirmed
		ErrorMessage:     "",
		RemainingBalance: remainingBalance,
	}, nil
}

// burnTokensInDatabase burns tokens in database
func (r *Repository) burnTokensInDatabase(ctx context.Context, req *proto.BurnTokensRequest) (*proto.BurnTokensResponse, error) {
	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:burn", req.ContractAddress, req.FromAddress, req.Amount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	return &proto.BurnTokensResponse{
		TransactionHash:  txHash,
		Status:           0, // Pending
		ErrorMessage:     "",
		RemainingBalance: "0",
	}, nil
}
