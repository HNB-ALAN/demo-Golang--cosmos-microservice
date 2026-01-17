package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	nettypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/types"
)

// Keeper manages the network module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new Network keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetNetwork returns a Network by its ID
func (k Keeper) GetNetwork(ctx sdk.Context, id string) (nettypes.Network, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.NetworkKey(id))
	if bz == nil {
		return nettypes.Network{}, fmt.Errorf("network with ID %s not found", id)
	}

	var network nettypes.Network
	if err := json.Unmarshal(bz, &network); err != nil {
		return nettypes.Network{}, fmt.Errorf("failed to unmarshal network: %w", err)
	}

	return network, nil
}

// SetNetwork sets a Network
func (k Keeper) SetNetwork(ctx sdk.Context, network nettypes.Network) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(network)
	if err != nil {
		return fmt.Errorf("failed to marshal network: %w", err)
	}
	store.Set(nettypes.NetworkKey(network.ID), bz)
	return nil
}

// GetAllNetworks returns all Networks
func (k Keeper) GetAllNetworks(ctx sdk.Context) []nettypes.Network {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, nettypes.NetworkKeyPrefix)
	defer iterator.Close()

	var networks []nettypes.Network
	for ; iterator.Valid(); iterator.Next() {
		var network nettypes.Network
		if err := json.Unmarshal(iterator.Value(), &network); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		networks = append(networks, network)
	}
	return networks
}

// GetNode returns a Node by its ID
func (k Keeper) GetNode(ctx sdk.Context, id string) (nettypes.Node, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.NodeKey(id))
	if bz == nil {
		return nettypes.Node{}, fmt.Errorf("node with ID %s not found", id)
	}

	var node nettypes.Node
	if err := json.Unmarshal(bz, &node); err != nil {
		return nettypes.Node{}, fmt.Errorf("failed to unmarshal node: %w", err)
	}

	return node, nil
}

// SetNode sets a Node
func (k Keeper) SetNode(ctx sdk.Context, node nettypes.Node) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %w", err)
	}
	store.Set(nettypes.NodeKey(node.ID), bz)
	return nil
}

// GetAllNodes returns all Nodes
func (k Keeper) GetAllNodes(ctx sdk.Context) []nettypes.Node {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, nettypes.NodeKeyPrefix)
	defer iterator.Close()

	var nodes []nettypes.Node
	for ; iterator.Valid(); iterator.Next() {
		var node nettypes.Node
		if err := json.Unmarshal(iterator.Value(), &node); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// GetConnection returns a Connection by its ID
func (k Keeper) GetConnection(ctx sdk.Context, id string) (nettypes.Connection, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.ConnectionKey(id))
	if bz == nil {
		return nettypes.Connection{}, fmt.Errorf("connection with ID %s not found", id)
	}

	var connection nettypes.Connection
	if err := json.Unmarshal(bz, &connection); err != nil {
		return nettypes.Connection{}, fmt.Errorf("failed to unmarshal connection: %w", err)
	}

	return connection, nil
}

// SetConnection sets a Connection
func (k Keeper) SetConnection(ctx sdk.Context, connection nettypes.Connection) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(connection)
	if err != nil {
		return fmt.Errorf("failed to marshal connection: %w", err)
	}
	store.Set(nettypes.ConnectionKey(connection.ID), bz)
	return nil
}

// GetAllConnections returns all Connections
func (k Keeper) GetAllConnections(ctx sdk.Context) []nettypes.Connection {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, nettypes.ConnectionKeyPrefix)
	defer iterator.Close()

	var connections []nettypes.Connection
	for ; iterator.Valid(); iterator.Next() {
		var connection nettypes.Connection
		if err := json.Unmarshal(iterator.Value(), &connection); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		connections = append(connections, connection)
	}
	return connections
}

// GetSync returns a NetworkSync by its ID
func (k Keeper) GetSync(ctx sdk.Context, id string) (nettypes.NetworkSync, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.SyncKey(id))
	if bz == nil {
		return nettypes.NetworkSync{}, fmt.Errorf("sync with ID %s not found", id)
	}

	var sync nettypes.NetworkSync
	if err := json.Unmarshal(bz, &sync); err != nil {
		return nettypes.NetworkSync{}, fmt.Errorf("failed to unmarshal sync: %w", err)
	}

	return sync, nil
}

// SetSync sets a NetworkSync
func (k Keeper) SetSync(ctx sdk.Context, sync nettypes.NetworkSync) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(sync)
	if err != nil {
		return fmt.Errorf("failed to marshal sync: %w", err)
	}
	store.Set(nettypes.SyncKey(sync.ID), bz)
	return nil
}

// GetAllSyncs returns all NetworkSyncs
func (k Keeper) GetAllSyncs(ctx sdk.Context) []nettypes.NetworkSync {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, nettypes.SyncKeyPrefix)
	defer iterator.Close()

	var syncs []nettypes.NetworkSync
	for ; iterator.Valid(); iterator.Next() {
		var sync nettypes.NetworkSync
		if err := json.Unmarshal(iterator.Value(), &sync); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		syncs = append(syncs, sync)
	}
	return syncs
}

// GetHealth returns a NetworkHealth by its ID
func (k Keeper) GetHealth(ctx sdk.Context, id string) (nettypes.NetworkHealth, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.HealthKey(id))
	if bz == nil {
		return nettypes.NetworkHealth{}, fmt.Errorf("health with ID %s not found", id)
	}

	var health nettypes.NetworkHealth
	if err := json.Unmarshal(bz, &health); err != nil {
		return nettypes.NetworkHealth{}, fmt.Errorf("failed to unmarshal health: %w", err)
	}

	return health, nil
}

// SetHealth sets a NetworkHealth
func (k Keeper) SetHealth(ctx sdk.Context, health nettypes.NetworkHealth) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(health)
	if err != nil {
		return fmt.Errorf("failed to marshal health: %w", err)
	}
	store.Set(nettypes.HealthKey(health.ID), bz)
	return nil
}

// GetAllHealths returns all NetworkHealths
func (k Keeper) GetAllHealths(ctx sdk.Context) []nettypes.NetworkHealth {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, nettypes.HealthKeyPrefix)
	defer iterator.Close()

	var healths []nettypes.NetworkHealth
	for ; iterator.Valid(); iterator.Next() {
		var health nettypes.NetworkHealth
		if err := json.Unmarshal(iterator.Value(), &health); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		healths = append(healths, health)
	}
	return healths
}

// GetParams returns the network module's parameters
func (k Keeper) GetParams(ctx sdk.Context) nettypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nettypes.ParamsKey)
	if bz == nil {
		return nettypes.DefaultParams()
	}

	var params nettypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return nettypes.DefaultParams()
	}

	return params
}

// SetParams sets the network module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params nettypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(nettypes.ParamsKey, bz)
}
