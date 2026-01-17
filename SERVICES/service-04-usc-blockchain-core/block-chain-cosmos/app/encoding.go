package app

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// ============================================================================
// ENCODING CONFIGURATION
// ============================================================================

// EncodingConfig represents the encoding configuration
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates a new encoding configuration
func MakeEncodingConfig() EncodingConfig {
	// Legacy Amino is disabled to avoid go-amino float panics; use protobuf only
	var amino *codec.LegacyAmino = nil

	// Create interface registry
	interfaceRegistry := types.NewInterfaceRegistry()

	// Create proto codec
	protoCodec := codec.NewProtoCodec(interfaceRegistry)

	// Create transaction config
	txConfig := tx.NewTxConfig(protoCodec, tx.DefaultSignModes)

	// Register standard Cosmos SDK interfaces (protobuf)
	std.RegisterInterfaces(interfaceRegistry)

	// Register module basics (interfaces only)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             protoCodec,
		TxConfig:          txConfig,
		Amino:             amino,
	}
}

// readGenesisChainID reads chain-id from genesis.json file
// Searches multiple common paths for genesis.json
func readGenesisChainID() string {
	genesisPaths := []string{
		"./genesis.json",
		"./config/genesis.json",
		"../config/genesis.json",
		"../../config/genesis.json",
		"/app/block-chain-cosmos/data/config/genesis.json",
	}
	for _, path := range genesisPaths {
		if data, err := os.ReadFile(path); err == nil {
			var genesis map[string]interface{}
			if err := json.Unmarshal(data, &genesis); err == nil {
				if chainID, ok := genesis["chain_id"].(string); ok && chainID != "" {
					return chainID
				}
			}
		}
	}
	return ""
}
