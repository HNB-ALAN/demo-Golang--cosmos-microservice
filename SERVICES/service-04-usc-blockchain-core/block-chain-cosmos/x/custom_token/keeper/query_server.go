package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/custom_token/v1/usc/custom_token/v1"
)

// QueryServer defines the interface for the custom_token module's query server
type QueryServer interface {
	QueryToken(context.Context, *blockchainproto.QueryTokenRequest) (*blockchainproto.QueryTokenResponse, error)
	QueryTokens(context.Context, *blockchainproto.QueryTokensRequest) (*blockchainproto.QueryTokensResponse, error)
	QueryTokenBalance(context.Context, *blockchainproto.QueryTokenBalanceRequest) (*blockchainproto.QueryTokenBalanceResponse, error)
	QueryTokenHolders(context.Context, *blockchainproto.QueryTokenHoldersRequest) (*blockchainproto.QueryTokenHoldersResponse, error)
	QueryTokenStats(context.Context, *blockchainproto.QueryTokenStatsRequest) (*blockchainproto.QueryTokenStatsResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// Token handles token queries
func (k queryServer) QueryToken(ctx context.Context, req *blockchainproto.QueryTokenRequest) (*blockchainproto.QueryTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	token, found := k.Keeper.GetToken(sdkCtx, req.TokenId)
	if !found {
		return nil, fmt.Errorf("token with ID %s not found", req.TokenId)
	}

	protoToken := &blockchainproto.CustomToken{
		Id:             token.ID,
		Address:        "",
		Creator:        token.Owner,
		Name:           token.Name,
		Symbol:         token.Symbol,
		Description:    "",
		Metadata:       nil,                                             // TODO: convert string to TokenMetadata
		CurrentSupply:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse token.TotalSupply
		MaxSupply:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse token.TotalSupply
		Decimals:       int32(token.Decimals),
		TokenType:      blockchainproto.TokenType_TOKEN_TYPE_FUNGIBLE,
		Status:         blockchainproto.TokenStatus_TOKEN_STATUS_ACTIVE,
		CreatedAt:      timestamppb.New(time.Unix(token.CreatedAt, 0)),
		UpdatedAt:      timestamppb.New(time.Unix(token.UpdatedAt, 0)),
		TotalMinted:    0,
		TotalBurned:    0,
		TotalTransfers: 0,
		HolderCount:    0,
		Memo:           "",
	}

	return &blockchainproto.QueryTokenResponse{Token: protoToken}, nil
}

// AllTokens handles all tokens queries
func (k queryServer) QueryTokens(ctx context.Context, req *blockchainproto.QueryTokensRequest) (*blockchainproto.QueryTokensResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	tokens, err := k.Keeper.GetAllTokens(sdkCtx)
	if err != nil {
		return nil, err
	}

	var protoTokens []*blockchainproto.CustomToken
	for _, token := range tokens {
		protoTokens = append(protoTokens, &blockchainproto.CustomToken{
			Id:             token.ID,
			Address:        "",
			Creator:        token.Owner,
			Name:           token.Name,
			Symbol:         token.Symbol,
			Description:    "",
			Metadata:       nil,                                             // TODO: convert string to TokenMetadata
			CurrentSupply:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse token.TotalSupply
			MaxSupply:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse token.TotalSupply
			Decimals:       int32(token.Decimals),
			TokenType:      blockchainproto.TokenType_TOKEN_TYPE_FUNGIBLE,
			Status:         blockchainproto.TokenStatus_TOKEN_STATUS_ACTIVE,
			CreatedAt:      timestamppb.New(time.Unix(token.CreatedAt, 0)),
			UpdatedAt:      timestamppb.New(time.Unix(token.UpdatedAt, 0)),
			TotalMinted:    0,
			TotalBurned:    0,
			TotalTransfers: 0,
			HolderCount:    0,
			Memo:           "",
		})
	}

	pageRes := &query.PageResponse{NextKey: nil, Total: uint64(len(protoTokens))}
	return &blockchainproto.QueryTokensResponse{Tokens: protoTokens, Pagination: pageRes}, nil
}

// Balance handles balance queries
func (k queryServer) QueryTokenBalance(ctx context.Context, req *blockchainproto.QueryTokenBalanceRequest) (*blockchainproto.QueryTokenBalanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	balance, found := k.Keeper.GetBalance(sdkCtx, req.TokenId, req.Holder)
	if !found {
		return nil, fmt.Errorf("balance not found for token %s and holder %s", req.TokenId, req.Holder)
	}

	protoBal := &blockchainproto.TokenBalance{
		TokenId:          req.TokenId,
		Holder:           req.Holder,
		Balance:          &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse balance.Amount
		LockedBalance:    &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		AvailableBalance: &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		LastUpdated:      timestamppb.New(time.Unix(balance.UpdatedAt, 0)),
	}

	return &blockchainproto.QueryTokenBalanceResponse{Balance: protoBal}, nil
}

// AllBalances handles all balances queries
// QueryTokenHolders returns holders for a token with pagination
func (k queryServer) QueryTokenHolders(ctx context.Context, req *blockchainproto.QueryTokenHoldersRequest) (*blockchainproto.QueryTokenHoldersResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	balances, err := k.Keeper.GetAllBalances(sdkCtx)
	if err != nil {
		return nil, err
	}

	var protoHolders []*blockchainproto.TokenHolder
	for _, b := range balances {
		if b.TokenID != req.TokenId {
			continue
		}
		protoHolders = append(protoHolders, &blockchainproto.TokenHolder{
			Holder:           b.Owner,
			TokenId:          b.TokenID,
			Balance:          &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse b.Amount
			Percentage:       0.0,
			FirstAcquired:    timestamppb.New(time.Unix(b.UpdatedAt, 0)),
			LastActivity:     timestamppb.New(time.Unix(b.UpdatedAt, 0)),
			TransactionCount: 0,
		})
	}
	pageRes := &query.PageResponse{NextKey: nil, Total: uint64(len(protoHolders))}

	return &blockchainproto.QueryTokenHoldersResponse{Holders: protoHolders, Pagination: pageRes}, nil
}

// QueryTokenStats returns aggregate stats for a token
func (k queryServer) QueryTokenStats(ctx context.Context, req *blockchainproto.QueryTokenStatsRequest) (*blockchainproto.QueryTokenStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	tokens, _ := k.Keeper.GetAllTokens(sdkCtx)
	balances, _ := k.Keeper.GetAllBalances(sdkCtx)

	var totalHolders int64
	for _, b := range balances {
		if b.TokenID == req.TokenId {
			totalHolders++
		}
	}

	protoStats := &blockchainproto.TokenStats{
		TotalTokens:     int64(len(tokens)),
		ActiveTokens:    int64(len(tokens)),
		PausedTokens:    0,
		BurnedTokens:    0,
		TotalSupply:     &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		TotalMinted:     &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		TotalBurned:     &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		TotalTransfers:  0,
		TotalHolders:    totalHolders,
		MostActiveToken: "",
		LastActivity:    timestamppb.New(time.Now()),
	}

	return &blockchainproto.QueryTokenStatsResponse{Stats: protoStats}, nil
}
