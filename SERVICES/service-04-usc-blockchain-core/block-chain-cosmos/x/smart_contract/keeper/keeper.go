package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sctypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/types"
)

// Keeper manages the contract module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new Contract keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetContract returns a SmartContract by its ID
func (k Keeper) GetContract(ctx sdk.Context, id string) (sctypes.SmartContract, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.ContractKey(id))
	if bz == nil {
		return sctypes.SmartContract{}, fmt.Errorf("contract with ID %s not found", id)
	}

	var contract sctypes.SmartContract
	if err := json.Unmarshal(bz, &contract); err != nil {
		return sctypes.SmartContract{}, fmt.Errorf("failed to unmarshal contract: %w", err)
	}

	return contract, nil
}

// SetContract sets a SmartContract
func (k Keeper) SetContract(ctx sdk.Context, contract sctypes.SmartContract) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}
	store.Set(sctypes.ContractKey(contract.ID), bz)
	return nil
}

// GetAllContracts returns all SmartContracts
func (k Keeper) GetAllContracts(ctx sdk.Context) []sctypes.SmartContract {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, sctypes.ContractKeyPrefix)
	defer iterator.Close()

	var contracts []sctypes.SmartContract
	for ; iterator.Valid(); iterator.Next() {
		var contract sctypes.SmartContract
		if err := json.Unmarshal(iterator.Value(), &contract); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		contracts = append(contracts, contract)
	}
	return contracts
}

// GetExecution returns a ContractExecution by its ID
func (k Keeper) GetExecution(ctx sdk.Context, id string) (sctypes.ContractExecution, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.ExecutionKey(id))
	if bz == nil {
		return sctypes.ContractExecution{}, fmt.Errorf("execution with ID %s not found", id)
	}

	var execution sctypes.ContractExecution
	if err := json.Unmarshal(bz, &execution); err != nil {
		return sctypes.ContractExecution{}, fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	return execution, nil
}

// SetExecution sets a ContractExecution
func (k Keeper) SetExecution(ctx sdk.Context, execution sctypes.ContractExecution) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(execution)
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}
	store.Set(sctypes.ExecutionKey(execution.ID), bz)
	return nil
}

// GetAllExecutions returns all ContractExecutions
func (k Keeper) GetAllExecutions(ctx sdk.Context) []sctypes.ContractExecution {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, sctypes.ExecutionKeyPrefix)
	defer iterator.Close()

	var executions []sctypes.ContractExecution
	for ; iterator.Valid(); iterator.Next() {
		var execution sctypes.ContractExecution
		if err := json.Unmarshal(iterator.Value(), &execution); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		executions = append(executions, execution)
	}
	return executions
}

// GetDeployment returns a ContractDeployment by its ID
func (k Keeper) GetDeployment(ctx sdk.Context, id string) (sctypes.ContractDeployment, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.DeploymentKey(id))
	if bz == nil {
		return sctypes.ContractDeployment{}, fmt.Errorf("deployment with ID %s not found", id)
	}

	var deployment sctypes.ContractDeployment
	if err := json.Unmarshal(bz, &deployment); err != nil {
		return sctypes.ContractDeployment{}, fmt.Errorf("failed to unmarshal deployment: %w", err)
	}

	return deployment, nil
}

// SetDeployment sets a ContractDeployment
func (k Keeper) SetDeployment(ctx sdk.Context, deployment sctypes.ContractDeployment) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment: %w", err)
	}
	store.Set(sctypes.DeploymentKey(deployment.ID), bz)
	return nil
}

// GetAllDeployments returns all ContractDeployments
func (k Keeper) GetAllDeployments(ctx sdk.Context) []sctypes.ContractDeployment {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, sctypes.DeploymentKeyPrefix)
	defer iterator.Close()

	var deployments []sctypes.ContractDeployment
	for ; iterator.Valid(); iterator.Next() {
		var deployment sctypes.ContractDeployment
		if err := json.Unmarshal(iterator.Value(), &deployment); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		deployments = append(deployments, deployment)
	}
	return deployments
}

// GetUpgrade returns a ContractUpgrade by its ID
func (k Keeper) GetUpgrade(ctx sdk.Context, id string) (sctypes.ContractUpgrade, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.UpgradeKey(id))
	if bz == nil {
		return sctypes.ContractUpgrade{}, fmt.Errorf("upgrade with ID %s not found", id)
	}

	var upgrade sctypes.ContractUpgrade
	if err := json.Unmarshal(bz, &upgrade); err != nil {
		return sctypes.ContractUpgrade{}, fmt.Errorf("failed to unmarshal upgrade: %w", err)
	}

	return upgrade, nil
}

// SetUpgrade sets a ContractUpgrade
func (k Keeper) SetUpgrade(ctx sdk.Context, upgrade sctypes.ContractUpgrade) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(upgrade)
	if err != nil {
		return fmt.Errorf("failed to marshal upgrade: %w", err)
	}
	store.Set(sctypes.UpgradeKey(upgrade.ID), bz)
	return nil
}

// GetAllUpgrades returns all ContractUpgrades
func (k Keeper) GetAllUpgrades(ctx sdk.Context) []sctypes.ContractUpgrade {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, sctypes.UpgradeKeyPrefix)
	defer iterator.Close()

	var upgrades []sctypes.ContractUpgrade
	for ; iterator.Valid(); iterator.Next() {
		var upgrade sctypes.ContractUpgrade
		if err := json.Unmarshal(iterator.Value(), &upgrade); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		upgrades = append(upgrades, upgrade)
	}
	return upgrades
}

// GetMigration returns a ContractMigration by its ID
func (k Keeper) GetMigration(ctx sdk.Context, id string) (sctypes.ContractMigration, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.MigrationKey(id))
	if bz == nil {
		return sctypes.ContractMigration{}, fmt.Errorf("migration with ID %s not found", id)
	}

	var migration sctypes.ContractMigration
	if err := json.Unmarshal(bz, &migration); err != nil {
		return sctypes.ContractMigration{}, fmt.Errorf("failed to unmarshal migration: %w", err)
	}

	return migration, nil
}

// SetMigration sets a ContractMigration
func (k Keeper) SetMigration(ctx sdk.Context, migration sctypes.ContractMigration) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(migration)
	if err != nil {
		return fmt.Errorf("failed to marshal migration: %w", err)
	}
	store.Set(sctypes.MigrationKey(migration.ID), bz)
	return nil
}

// GetAllMigrations returns all ContractMigrations
func (k Keeper) GetAllMigrations(ctx sdk.Context) []sctypes.ContractMigration {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, sctypes.MigrationKeyPrefix)
	defer iterator.Close()

	var migrations []sctypes.ContractMigration
	for ; iterator.Valid(); iterator.Next() {
		var migration sctypes.ContractMigration
		if err := json.Unmarshal(iterator.Value(), &migration); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		migrations = append(migrations, migration)
	}
	return migrations
}

// GetParams returns the contract module's parameters
func (k Keeper) GetParams(ctx sdk.Context) sctypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(sctypes.ParamsKey)
	if bz == nil {
		return sctypes.DefaultParams()
	}

	var params sctypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return sctypes.DefaultParams()
	}

	return params
}

// SetParams sets the contract module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params sctypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(sctypes.ParamsKey, bz)
}
