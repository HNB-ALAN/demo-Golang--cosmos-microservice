package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/types"
)

// Keeper manages the custom_token module's state
type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, paramSpace paramtypes.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
	}
}

// GetParams returns the current parameters for the custom_token module
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params, nil
}

// SetParams sets the parameters for the custom_token module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}

// GetToken returns a token by ID
func (k Keeper) GetToken(ctx sdk.Context, id string) (types.CustomToken, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.CustomToken{}, false
	}

	var token types.CustomToken
	if err := k.cdc.Unmarshal(bz, &token); err != nil {
		ctx.Logger().Error("Failed to unmarshal token",
			"error", err,
			"key", string(types.GetTokenKey(id)))
		return types.CustomToken{}, false
	}
	return token, true
}

// SetToken sets a token in the store
func (k Keeper) SetToken(ctx sdk.Context, token types.CustomToken) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(token.ID)
	bz := k.cdc.MustMarshal(&token)
	store.Set(key, bz)
	return nil
}

// GetAllTokens returns all tokens
func (k Keeper) GetAllTokens(ctx sdk.Context) ([]types.CustomToken, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TokenKeyPrefix))
	defer iterator.Close()

	var tokens []types.CustomToken
	for ; iterator.Valid(); iterator.Next() {
		var token types.CustomToken
		if err := k.cdc.Unmarshal(iterator.Value(), &token); err != nil {
			ctx.Logger().Error("Failed to unmarshal token, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// GetBalance returns a balance by token ID and owner
func (k Keeper) GetBalance(ctx sdk.Context, tokenID, owner string) (types.TokenBalance, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBalanceKey(tokenID, owner)
	bz := store.Get(key)
	if bz == nil {
		return types.TokenBalance{}, false
	}

	var balance types.TokenBalance
	if err := k.cdc.Unmarshal(bz, &balance); err != nil {
		ctx.Logger().Error("Failed to unmarshal balance",
			"error", err,
			"key", string(types.GetBalanceKey(tokenID, owner)))
		return types.TokenBalance{}, false
	}
	return balance, true
}

// SetBalance sets a balance in the store
func (k Keeper) SetBalance(ctx sdk.Context, balance types.TokenBalance) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBalanceKey(balance.TokenID, balance.Owner)
	bz := k.cdc.MustMarshal(&balance)
	store.Set(key, bz)
	return nil
}

// GetAllBalances returns all balances
func (k Keeper) GetAllBalances(ctx sdk.Context) ([]types.TokenBalance, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.BalanceKeyPrefix))
	defer iterator.Close()

	var balances []types.TokenBalance
	for ; iterator.Valid(); iterator.Next() {
		var balance types.TokenBalance
		if err := k.cdc.Unmarshal(iterator.Value(), &balance); err != nil {
			ctx.Logger().Error("Failed to unmarshal balance, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

// GetTransfer returns a transfer by ID
func (k Keeper) GetTransfer(ctx sdk.Context, id string) (types.TokenTransfer, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransferKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.TokenTransfer{}, false
	}

	var transfer types.TokenTransfer
	if err := k.cdc.Unmarshal(bz, &transfer); err != nil {
		ctx.Logger().Error("Failed to unmarshal transfer",
			"error", err,
			"key", string(types.GetTransferKey(id)))
		return types.TokenTransfer{}, false
	}
	return transfer, true
}

// SetTransfer sets a transfer in the store
func (k Keeper) SetTransfer(ctx sdk.Context, transfer types.TokenTransfer) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransferKey(transfer.ID)
	bz := k.cdc.MustMarshal(&transfer)
	store.Set(key, bz)
	return nil
}

// GetAllTransfers returns all transfers
func (k Keeper) GetAllTransfers(ctx sdk.Context) ([]types.TokenTransfer, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransferKeyPrefix))
	defer iterator.Close()

	var transfers []types.TokenTransfer
	for ; iterator.Valid(); iterator.Next() {
		var transfer types.TokenTransfer
		if err := k.cdc.Unmarshal(iterator.Value(), &transfer); err != nil {
			ctx.Logger().Error("Failed to unmarshal transfer, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// CreateToken creates a new token
func (k Keeper) CreateToken(ctx sdk.Context, token types.CustomToken) error {
	// Check if token already exists
	if _, exists := k.GetToken(ctx, token.ID); exists {
		return fmt.Errorf("token with ID %s already exists", token.ID)
	}

	// Validate token
	if token.ID == "" {
		return fmt.Errorf("token ID cannot be empty")
	}
	if token.Name == "" {
		return fmt.Errorf("token name cannot be empty")
	}
	if token.Symbol == "" {
		return fmt.Errorf("token symbol cannot be empty")
	}
	if token.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if token.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}
	if token.Decimals > 18 {
		return fmt.Errorf("decimals cannot exceed 18")
	}

	// Set token
	return k.SetToken(ctx, token)
}

// UpdateToken updates an existing token
func (k Keeper) UpdateToken(ctx sdk.Context, token types.CustomToken) error {
	// Check if token exists
	if _, exists := k.GetToken(ctx, token.ID); !exists {
		return fmt.Errorf("token with ID %s does not exist", token.ID)
	}

	// Validate token
	if token.ID == "" {
		return fmt.Errorf("token ID cannot be empty")
	}
	if token.Name == "" {
		return fmt.Errorf("token name cannot be empty")
	}
	if token.Symbol == "" {
		return fmt.Errorf("token symbol cannot be empty")
	}
	if token.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if token.Status == "" {
		return fmt.Errorf("status cannot be empty")
	}
	if token.Decimals > 18 {
		return fmt.Errorf("decimals cannot exceed 18")
	}

	// Set token
	return k.SetToken(ctx, token)
}

// DeleteToken deletes a token
func (k Keeper) DeleteToken(ctx sdk.Context, id string) error {
	// Check if token exists
	if _, exists := k.GetToken(ctx, id); !exists {
		return fmt.Errorf("token with ID %s does not exist", id)
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenKey(id)
	store.Delete(key)
	return nil
}

// MintToken mints tokens to an account
func (k Keeper) MintToken(ctx sdk.Context, tokenID, to, amount string) error {
	// Check if token exists
	_, exists := k.GetToken(ctx, tokenID)
	if !exists {
		return fmt.Errorf("token with ID %s does not exist", tokenID)
	}

	// Get current balance
	balance, exists := k.GetBalance(ctx, tokenID, to)
	if !exists {
		balance = types.NewTokenBalance(tokenID, to, "0")
	}

	// Implement proper amount addition logic using big.Int for precision
	currentAmount := new(big.Int)
	if balance.Amount != "" {
		if _, ok := currentAmount.SetString(balance.Amount, 10); !ok {
			return fmt.Errorf("invalid current balance amount: %s", balance.Amount)
		}
	}

	addAmount := new(big.Int)
	if _, ok := addAmount.SetString(amount, 10); !ok {
		return fmt.Errorf("invalid amount to add: %s", amount)
	}

	// Add amounts
	newAmount := new(big.Int).Add(currentAmount, addAmount)
	balance.Amount = newAmount.String()
	balance.UpdatedAt = ctx.BlockTime().Unix()

	// Set balance
	return k.SetBalance(ctx, balance)
}

// BurnToken burns tokens from an account
func (k Keeper) BurnToken(ctx sdk.Context, tokenID, from, amount string) error {
	// Check if token exists
	if _, exists := k.GetToken(ctx, tokenID); !exists {
		return fmt.Errorf("token with ID %s does not exist", tokenID)
	}

	// Get current balance
	balance, exists := k.GetBalance(ctx, tokenID, from)
	if !exists {
		return fmt.Errorf("no balance found for token %s and account %s", tokenID, from)
	}

	// Implement proper amount subtraction logic using big.Int for precision
	currentAmount := new(big.Int)
	if balance.Amount != "" {
		if _, ok := currentAmount.SetString(balance.Amount, 10); !ok {
			return fmt.Errorf("invalid current balance amount: %s", balance.Amount)
		}
	}

	burnAmount := new(big.Int)
	if _, ok := burnAmount.SetString(amount, 10); !ok {
		return fmt.Errorf("invalid amount to burn: %s", amount)
	}

	// Check sufficient balance
	if currentAmount.Cmp(burnAmount) < 0 {
		return fmt.Errorf("insufficient balance: %s < %s", balance.Amount, amount)
	}

	// Subtract amounts
	newAmount := new(big.Int).Sub(currentAmount, burnAmount)
	balance.Amount = newAmount.String()
	balance.UpdatedAt = ctx.BlockTime().Unix()

	// Set balance
	return k.SetBalance(ctx, balance)
}

// TransferToken transfers tokens between accounts
func (k Keeper) TransferToken(ctx sdk.Context, tokenID, from, to, amount string) error {
	// Check if token exists
	if _, exists := k.GetToken(ctx, tokenID); !exists {
		return fmt.Errorf("token with ID %s does not exist", tokenID)
	}

	// Get sender balance
	senderBalance, exists := k.GetBalance(ctx, tokenID, from)
	if !exists {
		return fmt.Errorf("no balance found for token %s and account %s", tokenID, from)
	}

	// Get receiver balance
	receiverBalance, exists := k.GetBalance(ctx, tokenID, to)
	if !exists {
		receiverBalance = types.NewTokenBalance(tokenID, to, "0")
	}

	// Implement proper amount transfer logic using big.Int for precision
	senderAmount := new(big.Int)
	if senderBalance.Amount != "" {
		if _, ok := senderAmount.SetString(senderBalance.Amount, 10); !ok {
			return fmt.Errorf("invalid sender balance amount: %s", senderBalance.Amount)
		}
	}

	transferAmount := new(big.Int)
	if _, ok := transferAmount.SetString(amount, 10); !ok {
		return fmt.Errorf("invalid transfer amount: %s", amount)
	}

	// Check sufficient balance
	if senderAmount.Cmp(transferAmount) < 0 {
		return fmt.Errorf("insufficient balance: %s < %s", senderBalance.Amount, amount)
	}

	// Calculate new sender balance
	newSenderAmount := new(big.Int).Sub(senderAmount, transferAmount)
	senderBalance.Amount = newSenderAmount.String()
	senderBalance.UpdatedAt = ctx.BlockTime().Unix()

	// Calculate new receiver balance
	receiverAmount := new(big.Int)
	if receiverBalance.Amount != "" {
		if _, ok := receiverAmount.SetString(receiverBalance.Amount, 10); !ok {
			return fmt.Errorf("invalid receiver balance amount: %s", receiverBalance.Amount)
		}
	}
	newReceiverAmount := new(big.Int).Add(receiverAmount, transferAmount)
	receiverBalance.Amount = newReceiverAmount.String()
	receiverBalance.UpdatedAt = ctx.BlockTime().Unix()

	// Set balances
	if err := k.SetBalance(ctx, senderBalance); err != nil {
		return err
	}
	if err := k.SetBalance(ctx, receiverBalance); err != nil {
		return err
	}

	// Create transfer record
	transfer := types.NewTokenTransfer(fmt.Sprintf("%s_%s_%d", tokenID, from, ctx.BlockTime().Unix()), tokenID, from, to, amount)
	return k.SetTransfer(ctx, transfer)
}

// InitGenesis initializes the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set parameters
	k.SetParams(ctx, genState.Params)

	// Set tokens
	for _, token := range genState.Tokens {
		k.SetToken(ctx, token)
	}

	// Set balances
	for _, balance := range genState.Balances {
		k.SetBalance(ctx, balance)
	}

	// Set transfers
	for _, transfer := range genState.Transfers {
		k.SetTransfer(ctx, transfer)
	}
}

// ExportGenesis exports the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// Get parameters
	params, _ := k.GetParams(ctx)

	// Get all tokens
	tokens, _ := k.GetAllTokens(ctx)

	// Get all balances
	balances, _ := k.GetAllBalances(ctx)

	// Get all transfers
	transfers, _ := k.GetAllTransfers(ctx)

	return &types.GenesisState{
		Tokens:    tokens,
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}
}

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("CustomToken BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Token expiration checks
	// - Balance updates
	// - Emitting events

	// Example: Emit a block start event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(types.AttributeKeyCreatedAt, strconv.FormatInt(ctx.BlockTime().Unix(), 10)),
		),
	)
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) []abci.ValidatorUpdate {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("CustomToken EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Token cleanup
	// - Balance processing
	// - Emitting events

	// Example: Emit a block end event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenUpdated,
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, strconv.FormatInt(ctx.BlockTime().Unix(), 10)),
		),
	)

	// Return validator updates (if any)
	return []abci.ValidatorUpdate{}
}
