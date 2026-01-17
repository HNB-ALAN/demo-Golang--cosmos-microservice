package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================
// BLOCK HANDLERS
// ============================================================================

// beginBlocker handles begin block logic
// Signature: func(Context) (BeginBlock, error)
func (app *USCApp) beginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	if app.mm != nil {
		app.mm.BeginBlock(ctx)
	}
	return sdk.BeginBlock{}, nil
}

// endBlocker handles end block logic
// Signature: func(Context) (EndBlock, error)
func (app *USCApp) endBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	if app.mm != nil {
		endBlock, err := app.mm.EndBlock(ctx)
		if err != nil {
			return sdk.EndBlock{}, err
		}
		return endBlock, nil
	}
	return sdk.EndBlock{}, nil
}
