package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/types"
)

// Keeper manages the USC module state
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace
	bk         keeper.Keeper
}

// NewKeeper creates a new USC keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	bk keeper.Keeper,
) Keeper {
	// Set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
		bk:         bk,
	}
}

// GetParams returns the current parameters
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params, nil
}

// SetParams sets the parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("params validation failed: %w", err)
	}

	// Ensure paramSpace has KeyTable before setting params
	if !k.paramSpace.HasKeyTable() {
		k.paramSpace = k.paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	// Set params using paramSpace
	// NOTE: During InitGenesis, store service may not be available in context.
	// If this fails, caller should handle gracefully (e.g., use default params).
	k.paramSpace.SetParamSet(ctx, &params)

	return nil
}

// GetBalance returns the balance for an address
func (k Keeper) GetBalance(ctx sdk.Context, address string) (types.Balance, error) {
	if k.storeKey == nil {
		return types.Balance{}, fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetBalanceKey(address)

	if !store.Has(key) {
		return types.Balance{}, fmt.Errorf("balance not found for address: %s", address)
	}

	var balance types.Balance
	if err := json.Unmarshal(store.Get(key), &balance); err != nil {
		return types.Balance{}, err
	}

	return balance, nil
}

// SetBalance sets the balance for an address
func (k Keeper) SetBalance(ctx sdk.Context, address string, balance types.Balance) error {
	if k.storeKey == nil {
		return fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetBalanceKey(address)

	bz, err := json.Marshal(&balance)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// GetAllBalances returns all balances
func (k Keeper) GetAllBalances(ctx sdk.Context) ([]types.Balance, error) {
	if k.storeKey == nil {
		return nil, fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.BalanceKeyPrefix))
	defer iterator.Close()

	var balances []types.Balance
	for ; iterator.Valid(); iterator.Next() {
		var balance types.Balance
		if err := json.Unmarshal(iterator.Value(), &balance); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

// GetTransfer returns a transfer by transaction hash
func (k Keeper) GetTransfer(ctx sdk.Context, txHash string) (types.Transfer, error) {
	if k.storeKey == nil {
		return types.Transfer{}, fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransferKey(txHash)

	if !store.Has(key) {
		return types.Transfer{}, fmt.Errorf("transfer not found: %s", txHash)
	}

	var transfer types.Transfer
	if err := json.Unmarshal(store.Get(key), &transfer); err != nil {
		return types.Transfer{}, err
	}

	return transfer, nil
}

// SetTransfer sets a transfer
func (k Keeper) SetTransfer(ctx sdk.Context, transfer types.Transfer) error {
	if k.storeKey == nil {
		return fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransferKey(fmt.Sprintf("%s_%d", transfer.FromAddress, transfer.Timestamp))

	bz, err := json.Marshal(&transfer)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// GetAllTransfers returns all transfers
func (k Keeper) GetAllTransfers(ctx sdk.Context) ([]types.Transfer, error) {
	if k.storeKey == nil {
		return nil, fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransferKeyPrefix))
	defer iterator.Close()

	var transfers []types.Transfer
	for ; iterator.Valid(); iterator.Next() {
		var transfer types.Transfer
		if err := json.Unmarshal(iterator.Value(), &transfer); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// GetTotalSupply returns the total supply
func (k Keeper) GetTotalSupply(ctx sdk.Context) (string, error) {
	if k.storeKey == nil {
		return "0", fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetSupplyKey()

	if !store.Has(key) {
		return "0", nil
	}

	return string(store.Get(key)), nil
}

// SetTotalSupply sets the total supply
func (k Keeper) SetTotalSupply(ctx sdk.Context, supply string) error {
	if k.storeKey == nil {
		return fmt.Errorf("storeKey is nil - keeper not properly initialized")
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetSupplyKey()

	store.Set(key, []byte(supply))
	return nil
}

// TransferUSC transfers USC tokens between addresses using blockchain-proto message
func (k Keeper) TransferUSC(ctx sdk.Context, msg *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error) {
	// Get sender balance
	senderBalance, err := k.GetBalance(ctx, msg.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("sender balance not found: %w", err)
	}

	// Check if sender has sufficient balance
	senderAmount, err := strconv.ParseInt(senderBalance.Amount, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sender amount: %w", err)
	}

	transferAmount := msg.Amount.Amount.Int64()

	if senderAmount < transferAmount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Update sender balance
	senderBalance.Amount = strconv.FormatInt(senderAmount-transferAmount, 10)
	if err := k.SetBalance(ctx, msg.FromAddress, senderBalance); err != nil {
		return nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	// Get or create receiver balance
	receiverBalance, err := k.GetBalance(ctx, msg.ToAddress)
	if err != nil {
		// Create new balance for receiver
		receiverBalance = types.Balance{
			Address: msg.ToAddress,
			Amount:  "0",
			Denom:   senderBalance.Denom,
		}
	}

	// Update receiver balance
	receiverAmount, err := strconv.ParseInt(receiverBalance.Amount, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver amount: %w", err)
	}

	receiverBalance.Amount = strconv.FormatInt(receiverAmount+transferAmount, 10)
	if err := k.SetBalance(ctx, msg.ToAddress, receiverBalance); err != nil {
		return nil, fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// Create transfer record
	transfer := types.NewTransfer(msg.FromAddress, msg.ToAddress, msg.Amount.Amount.String(), senderBalance.Denom)
	if err := k.SetTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("failed to record transfer: %w", err)
	}

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(ctx, msg.FromAddress, msg.ToAddress, msg.Amount.Amount.String(), "transfer_usc", "", "")

	return &blockchainproto.MsgTransferUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// MintUSC mints new USC tokens using blockchain-proto message
func (k Keeper) MintUSC(ctx sdk.Context, msg *blockchainproto.MsgMintUSC) (*blockchainproto.MsgMintUSCResponse, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	if !params.MintEnabled {
		return nil, fmt.Errorf("minting is disabled")
	}

	// Get current total supply
	currentSupply, err := k.GetTotalSupply(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %w", err)
	}

	currentSupplyInt, err := strconv.ParseInt(currentSupply, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid current supply: %w", err)
	}

	mintAmount := msg.Amount.Amount.Int64()

	maxSupplyInt, err := strconv.ParseInt(params.MaxSupply, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid max supply: %w", err)
	}

	if currentSupplyInt+mintAmount > maxSupplyInt {
		return nil, fmt.Errorf("minting would exceed max supply")
	}

	// Get or create balance for recipient
	balance, err := k.GetBalance(ctx, msg.Minter)
	if err != nil {
		// Create new balance
		balance = types.Balance{
			Address: msg.Minter,
			Amount:  "0",
			Denom:   params.TokenSymbol,
		}
	}

	// Update balance
	currentAmount, err := strconv.ParseInt(balance.Amount, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid current amount: %w", err)
	}

	balance.Amount = strconv.FormatInt(currentAmount+mintAmount, 10)
	if err := k.SetBalance(ctx, msg.Minter, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Update total supply
	newSupply := strconv.FormatInt(currentSupplyInt+mintAmount, 10)
	if err := k.SetTotalSupply(ctx, newSupply); err != nil {
		return nil, fmt.Errorf("failed to update total supply: %w", err)
	}

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(ctx, msg.Minter, "", msg.Amount.Amount.String(), "mint_usc", "", "")

	return &blockchainproto.MsgMintUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// BurnUSC burns USC tokens using blockchain-proto message
func (k Keeper) BurnUSC(ctx sdk.Context, msg *blockchainproto.MsgBurnUSC) (*blockchainproto.MsgBurnUSCResponse, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	if !params.BurnEnabled {
		return nil, fmt.Errorf("burning is disabled")
	}

	// Get sender balance
	balance, err := k.GetBalance(ctx, msg.Burner)
	if err != nil {
		return nil, fmt.Errorf("balance not found: %w", err)
	}

	// Check if sender has sufficient balance
	currentAmount, err := strconv.ParseInt(balance.Amount, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid current amount: %w", err)
	}

	burnAmount := msg.Amount.Amount.Int64()

	if currentAmount < burnAmount {
		return nil, fmt.Errorf("insufficient balance to burn")
	}

	// Update balance
	balance.Amount = strconv.FormatInt(currentAmount-burnAmount, 10)
	if err := k.SetBalance(ctx, msg.Burner, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Update total supply
	currentSupply, err := k.GetTotalSupply(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total supply: %w", err)
	}

	currentSupplyInt, err := strconv.ParseInt(currentSupply, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid current supply: %w", err)
	}

	newSupply := strconv.FormatInt(currentSupplyInt-burnAmount, 10)
	if err := k.SetTotalSupply(ctx, newSupply); err != nil {
		return nil, fmt.Errorf("failed to update total supply: %w", err)
	}

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(ctx, msg.Burner, "", msg.Amount.Amount.String(), "burn_usc", "", "")

	return &blockchainproto.MsgBurnUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}
