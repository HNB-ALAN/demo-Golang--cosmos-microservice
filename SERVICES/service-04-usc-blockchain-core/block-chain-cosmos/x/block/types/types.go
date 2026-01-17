package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "block"

	// StoreKey defines the primary store key for this module
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_block"

	// RouterKey defines the routing key for this module
	RouterKey = ModuleName
)

// Event types
const (
	EventTypeBlockCreated   = "block_created"
	EventTypeBlockUpdated   = "block_updated"
	EventTypeBlockDeleted   = "block_deleted"
	EventTypeBlockValidated = "block_validated"
	EventTypeBlockFinalized = "block_finalized"
)

// Event attribute keys
const (
	AttributeKeyBlockID        = "block_id"
	AttributeKeyBlockHeight    = "block_height"
	AttributeKeyBlockHash      = "block_hash"
	AttributeKeyBlockTime      = "block_time"
	AttributeKeyBlockValidator = "block_validator"
	AttributeKeyBlockSize      = "block_size"
	AttributeKeyBlockTxCount   = "block_tx_count"
	AttributeKeyBlockFinalizer = "block_finalizer"
)

// Block represents a blockchain block
type Block struct {
	ID           string    `json:"id"`
	Height       int64     `json:"height"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previous_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Validator    string    `json:"validator"`
	Size         int64     `json:"size"`
	TxCount      int64     `json:"tx_count"`
	GasUsed      int64     `json:"gas_used"`
	GasLimit     int64     `json:"gas_limit"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BlockHeader represents block header information
type BlockHeader struct {
	Height       int64     `json:"height"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previous_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Validator    string    `json:"validator"`
	Size         int64     `json:"size"`
	TxCount      int64     `json:"tx_count"`
	GasUsed      int64     `json:"gas_used"`
	GasLimit     int64     `json:"gas_limit"`
}

// BlockData represents block data
type BlockData struct {
	BlockID    string    `json:"block_id"`
	Height     int64     `json:"height"`
	Hash       string    `json:"hash"`
	Data       []byte    `json:"data"`
	Size       int64     `json:"size"`
	Compressed bool      `json:"compressed"`
	CreatedAt  time.Time `json:"created_at"`
}

// BlockValidation represents block validation result
type BlockValidation struct {
	BlockID        string    `json:"block_id"`
	Height         int64     `json:"height"`
	Hash           string    `json:"hash"`
	IsValid        bool      `json:"is_valid"`
	ValidationTime time.Time `json:"validation_time"`
	Validator      string    `json:"validator"`
	Errors         []string  `json:"errors"`
	Warnings       []string  `json:"warnings"`
}

// GenesisState defines the block module's genesis state
type GenesisState struct {
	Blocks      []Block           `json:"blocks"`
	BlockData   []BlockData       `json:"block_data"`
	Validations []BlockValidation `json:"validations"`
	Params      Params            `json:"params"`
}

// Params defines the parameters for the block module
type Params struct {
	MaxBlockSize       int64  `json:"max_block_size"`
	MaxTxCount         int64  `json:"max_tx_count"`
	MaxGasLimit        int64  `json:"max_gas_limit"`
	BlockTime          string `json:"block_time"`
	ValidationEnabled  bool   `json:"validation_enabled"`
	CompressionEnabled bool   `json:"compression_enabled"`
}

// DefaultParams returns the default parameters for the block module
func DefaultParams() Params {
	return Params{
		MaxBlockSize:       1048576, // 1MB
		MaxTxCount:         10000,
		MaxGasLimit:        10000000,
		BlockTime:          "3s", // Production-optimized: 3 seconds per block
		ValidationEnabled:  true,
		CompressionEnabled: true,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxBlockSize <= 0 {
		return fmt.Errorf("max block size must be positive")
	}
	if p.MaxTxCount <= 0 {
		return fmt.Errorf("max tx count must be positive")
	}
	if p.MaxGasLimit <= 0 {
		return fmt.Errorf("max gas limit must be positive")
	}
	if p.BlockTime == "" {
		return fmt.Errorf("block time cannot be empty")
	}
	return nil
}

// Validate validates a block
func (b Block) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("block ID cannot be empty")
	}
	if b.Height < 0 {
		return fmt.Errorf("block height must be non-negative")
	}
	if b.Hash == "" {
		return fmt.Errorf("block hash cannot be empty")
	}
	if b.Validator == "" {
		return fmt.Errorf("block validator cannot be empty")
	}
	if b.Size < 0 {
		return fmt.Errorf("block size must be non-negative")
	}
	if b.TxCount < 0 {
		return fmt.Errorf("block tx count must be non-negative")
	}
	if b.GasUsed < 0 {
		return fmt.Errorf("block gas used must be non-negative")
	}
	if b.GasLimit < 0 {
		return fmt.Errorf("block gas limit must be non-negative")
	}
	if b.Status == "" {
		return fmt.Errorf("block status cannot be empty")
	}
	return nil
}

// Validate validates block data
func (bd BlockData) Validate() error {
	if bd.BlockID == "" {
		return fmt.Errorf("block ID cannot be empty")
	}
	if bd.Height < 0 {
		return fmt.Errorf("block height must be non-negative")
	}
	if bd.Hash == "" {
		return fmt.Errorf("block hash cannot be empty")
	}
	if bd.Size < 0 {
		return fmt.Errorf("block data size must be non-negative")
	}
	return nil
}

// Validate validates block validation
func (bv BlockValidation) Validate() error {
	if bv.BlockID == "" {
		return fmt.Errorf("block ID cannot be empty")
	}
	if bv.Height < 0 {
		return fmt.Errorf("block height must be non-negative")
	}
	if bv.Hash == "" {
		return fmt.Errorf("block hash cannot be empty")
	}
	if bv.Validator == "" {
		return fmt.Errorf("validator cannot be empty")
	}
	return nil
}

// ============================================================================
// HASH AND ADDRESS UTILITIES
// ============================================================================

// CalculateTransactionHash calculates a deterministic transaction hash from transaction data
// This is used when we don't have access to the raw transaction bytes
// Format: chain_id:height:timestamp:from:to:amount:type:data:memo
func CalculateTransactionHash(ctx sdk.Context, from, to, amount, txType, data, memo string) string {
	header := ctx.BlockHeader()
	timestamp := ctx.BlockTime().Unix()

	// Combine transaction fields to create unique hash
	dataStr := fmt.Sprintf("%s:%d:%d:%s:%s:%s:%s:%s:%s",
		header.ChainID,
		header.Height,
		timestamp,
		from,
		to,
		amount,
		txType,
		data,
		memo,
	)
	hash := sha256.Sum256([]byte(dataStr))
	return hex.EncodeToString(hash[:])
}

// CalculateContractAddress generates a deterministic contract address
// Format: deployer_address:chain_id:height:nonce:code_hash
func CalculateContractAddress(ctx sdk.Context, deployer, codeHash string, nonce uint64) string {
	header := ctx.BlockHeader()

	// Combine fields to create unique address
	dataStr := fmt.Sprintf("%s:%s:%d:%d:%s",
		deployer,
		header.ChainID,
		header.Height,
		nonce,
		codeHash,
	)
	hash := sha256.Sum256([]byte(dataStr))
	// Take first 20 bytes (40 hex chars) for Ethereum-style address
	addressHex := hex.EncodeToString(hash[:20])
	return fmt.Sprintf("0x%s", addressHex)
}

// CalculateHashFromData calculates a hash from arbitrary data
func CalculateHashFromData(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CalculateHashFromString calculates a hash from a string
func CalculateHashFromString(data string) string {
	return CalculateHashFromData([]byte(data))
}
