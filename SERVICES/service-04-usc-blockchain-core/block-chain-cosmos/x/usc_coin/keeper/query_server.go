package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
)

// QueryServer defines the query server interface using blockchain-proto types
type QueryServer interface {
	QueryUSCBalance(context.Context, *blockchainproto.QueryUSCBalanceRequest) (*blockchainproto.QueryUSCBalanceResponse, error)
	QueryUSCSupply(context.Context, *blockchainproto.QueryUSCSupplyRequest) (*blockchainproto.QueryUSCSupplyResponse, error)
	QueryUSCHolders(context.Context, *blockchainproto.QueryUSCHoldersRequest) (*blockchainproto.QueryUSCHoldersResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryUSCBalance returns the USC balance for a specific address using blockchain-proto types
func (k queryServer) QueryUSCBalance(ctx context.Context, req *blockchainproto.QueryUSCBalanceRequest) (*blockchainproto.QueryUSCBalanceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	if req.Address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	balance, err := k.GetBalance(sdkCtx, req.Address)
	if err != nil {
		return nil, fmt.Errorf("balance not found: %w", err)
	}

	// Convert to blockchain-proto Coin type
	amount, ok := math.NewIntFromString(balance.Amount)
	if !ok {
		return nil, fmt.Errorf("invalid amount format: %s", balance.Amount)
	}
	blockchainBalance := &sdk.Coin{
		Denom:  balance.Denom,
		Amount: amount,
	}

	return &blockchainproto.QueryUSCBalanceResponse{
		Balance: blockchainBalance,
	}, nil
}

// QueryUSCSupply returns the total supply of USC tokens using blockchain-proto types
func (k queryServer) QueryUSCSupply(ctx context.Context, req *blockchainproto.QueryUSCSupplyRequest) (*blockchainproto.QueryUSCSupplyResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	supply, err := k.GetTotalSupply(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %w", err)
	}

	// Convert to blockchain-proto Coin type
	amount, ok := math.NewIntFromString(supply)
	if !ok {
		return nil, fmt.Errorf("invalid supply format: %s", supply)
	}
	blockchainSupply := &sdk.Coin{
		Denom:  "usc",
		Amount: amount,
	}

	return &blockchainproto.QueryUSCSupplyResponse{
		Supply: blockchainSupply,
	}, nil
}

// QueryUSCHolders returns all USC holders with pagination using blockchain-proto types
func (k queryServer) QueryUSCHolders(ctx context.Context, req *blockchainproto.QueryUSCHoldersRequest) (*blockchainproto.QueryUSCHoldersResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all balances
	balances, err := k.GetAllBalances(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	// Convert to blockchain-proto USCHolder types
	var blockchainHolders []*blockchainproto.USCHolder
	for _, balance := range balances {
		amount, ok := math.NewIntFromString(balance.Amount)
		if !ok {
			continue // Skip invalid amounts
		}
		if amount.IsPositive() {
			blockchainHolder := &blockchainproto.USCHolder{
				Address: balance.Address,
				Balance: &sdk.Coin{
					Denom:  balance.Denom,
					Amount: amount,
				},
			}
			blockchainHolders = append(blockchainHolders, blockchainHolder)
		}
	}

	// Create pagination response
	var pagination *query.PageResponse
	if req.Pagination != nil {
		pagination = &query.PageResponse{
			NextKey: nil, // TODO: Implement proper pagination
			Total:   uint64(len(blockchainHolders)),
		}
	}

	return &blockchainproto.QueryUSCHoldersResponse{
		Holders:    blockchainHolders,
		Pagination: pagination,
	}, nil
}

// Note: Custom query types removed as they are replaced by blockchain-proto query types
// The blockchain-proto interface provides QueryBalance, QueryBalances, QueryTransfer, QueryTransfers, QueryTotalSupply, and QueryParams
