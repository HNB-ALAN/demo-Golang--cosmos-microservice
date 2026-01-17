package utils

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
)

// GetSDKContext creates a sdk.Context from context.Context
// COSMOS SDK 0.53.4 STANDARD: Use sdk.UnwrapSDKContext(ctx) to get commit context
// This is the same pattern used in all Cosmos SDK query servers:
// - x/usc_coin/keeper/query_server.go: sdkCtx := sdk.UnwrapSDKContext(ctx)
// - x/transaction/keeper/query_server.go: sdkCtx := sdk.UnwrapSDKContext(ctx)
// - x/store_network/keeper/query_server.go: sdkCtx := sdk.UnwrapSDKContext(ctx)
// sdk.UnwrapSDKContext(ctx) automatically unwraps the context from gRPC and returns commit context
// Use recover to handle panic if context is not wrapped
func GetSDKContext(ctx context.Context, cosmosApp *app.USCApp, logger *logging.Logger) (sdk.Context, error) {
	if cosmosApp == nil || cosmosApp.BaseApp == nil {
		return sdk.Context{}, errors.New("cosmosApp not initialized")
	}

	// Try to unwrap SDK context from gRPC context
	var sdkCtx sdk.Context
	func() {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// Context is not wrapped, will use fallback
				if logger != nil {
					logger.Debug("Context not wrapped (panic caught), using BaseApp.NewContext(true) as fallback",
						logging.Any("panic", panicVal))
				}
			}
		}()
		sdkCtx = sdk.UnwrapSDKContext(ctx)
	}()

	// If context is not wrapped (e.g., called from non-gRPC handler), fallback to BaseApp.NewContext(true)
	if sdkCtx.IsZero() {
		if logger != nil {
			logger.Debug("Context not wrapped, using BaseApp.NewContext(true) as fallback")
		}
		// Use NewContext(true) to read from committed state
		sdkCtx = cosmosApp.BaseApp.NewContext(true)
	}

	return sdkCtx, nil
}

// GetSDKContextForCheckTx creates a sdk.Context for CheckTx operations
// Use this when you need to simulate transactions without committing
func GetSDKContextForCheckTx(ctx context.Context, cosmosApp *app.USCApp, logger *logging.Logger) (sdk.Context, error) {
	if cosmosApp == nil || cosmosApp.BaseApp == nil {
		return sdk.Context{}, errors.New("cosmosApp not initialized")
	}

	// For CheckTx, always use NewContext(false)
	return cosmosApp.BaseApp.NewContext(false), nil
}

// GetSDKContextForWrite creates a writable sdk.Context for write operations
// ROOT FIX: Use NewContext(false) to allow writes to keeper (will be committed on next block)
// This is critical for SetCertificate, UpdateCertificate, and other write operations
func GetSDKContextForWrite(ctx context.Context, cosmosApp *app.USCApp, logger *logging.Logger) (sdk.Context, error) {
	if cosmosApp == nil || cosmosApp.BaseApp == nil {
		return sdk.Context{}, errors.New("cosmosApp not initialized")
	}

	// Try to unwrap SDK context from gRPC context first
	var sdkCtx sdk.Context
	func() {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// Context is not wrapped, will use fallback
				if logger != nil {
					logger.Debug("Context not wrapped (panic caught), using BaseApp.NewContext(false) as fallback for write",
						logging.Any("panic", panicVal))
				}
			}
		}()
		sdkCtx = sdk.UnwrapSDKContext(ctx)
	}()

	// If context is not wrapped (e.g., called from non-gRPC handler), use NewContext(false) for write operations
	if sdkCtx.IsZero() {
		if logger != nil {
			logger.Debug("Context not wrapped, using BaseApp.NewContext(false) for write operation")
		}
		// Use NewContext(false) to allow writes (will be committed on next block)
		sdkCtx = cosmosApp.BaseApp.NewContext(false)
	}

	return sdkCtx, nil
}
