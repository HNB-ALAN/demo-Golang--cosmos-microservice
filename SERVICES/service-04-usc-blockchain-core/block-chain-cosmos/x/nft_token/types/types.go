package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "nft_token"

	// RouterKey defines the message route for the nft module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the nft module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeNFTCreated        = "nft_created"
	EventTypeNFTTransferred    = "nft_transferred"
	EventTypeNFTBurned         = "nft_burned"
	EventTypeNFTUpdated        = "nft_updated"
	EventTypeCollectionCreated = "collection_created"
	EventTypeCollectionUpdated = "collection_updated"
)

// Event attribute keys
const (
	AttributeKeyNFTID        = "nft_id"
	AttributeKeyCollectionID = "collection_id"
	AttributeKeyOwner        = "owner"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyTokenURI     = "token_uri"
	AttributeKeyModule       = ModuleName
)

// NFT represents a non-fungible token
type NFT struct {
	ID           string            `json:"id"`
	CollectionID string            `json:"collection_id"`
	Owner        string            `json:"owner"`
	TokenURI     string            `json:"token_uri"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Image        string            `json:"image"`
	Attributes   map[string]string `json:"attributes"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Metadata     map[string]string `json:"metadata"`
}

// Collection represents an NFT collection
type Collection struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Symbol        string            `json:"symbol"`
	Image         string            `json:"image"`
	Owner         string            `json:"owner"`
	MaxSupply     int64             `json:"max_supply"`
	CurrentSupply int64             `json:"current_supply"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Metadata      map[string]string `json:"metadata"`
}

// GenesisState represents the genesis state of the nft module
type GenesisState struct {
	NFTs        []NFT        `json:"nfts"`
	Collections []Collection `json:"collections"`
	Params      Params       `json:"params"`
}

// Params represents the parameters for the nft module
type Params struct {
	MaxNFTsPerCollection int64         `json:"max_nfts_per_collection"`
	MaxCollections       int64         `json:"max_collections"`
	MintingFee           string        `json:"minting_fee"`
	TransferFee          string        `json:"transfer_fee"`
	BurnFee              string        `json:"burn_fee"`
	MetadataRetention    time.Duration `json:"metadata_retention"`
}

// DefaultParams returns the default parameters for the nft module
func DefaultParams() Params {
	return Params{
		MaxNFTsPerCollection: 10000,
		MaxCollections:       1000,
		MintingFee:           "1000000usc",         // 1 USC
		TransferFee:          "100000usc",          // 0.1 USC
		BurnFee:              "50000usc",           // 0.05 USC
		MetadataRetention:    365 * 24 * time.Hour, // 1 year
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxNFTsPerCollection <= 0 {
		return fmt.Errorf("max NFTs per collection must be positive")
	}
	if p.MaxCollections <= 0 {
		return fmt.Errorf("max collections must be positive")
	}
	if p.MintingFee == "" {
		return fmt.Errorf("minting fee cannot be empty")
	}
	if p.TransferFee == "" {
		return fmt.Errorf("transfer fee cannot be empty")
	}
	if p.BurnFee == "" {
		return fmt.Errorf("burn fee cannot be empty")
	}
	if p.MetadataRetention <= 0 {
		return fmt.Errorf("metadata retention must be positive")
	}
	return nil
}

// Validate validates an NFT
func (n NFT) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("NFT ID cannot be empty")
	}
	if n.CollectionID == "" {
		return fmt.Errorf("collection ID cannot be empty")
	}
	if n.Owner == "" {
		return fmt.Errorf("NFT owner cannot be empty")
	}
	if n.TokenURI == "" {
		return fmt.Errorf("token URI cannot be empty")
	}
	return nil
}

// Validate validates a collection
func (c Collection) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("collection ID cannot be empty")
	}
	if c.Name == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	if c.Symbol == "" {
		return fmt.Errorf("collection symbol cannot be empty")
	}
	if c.Owner == "" {
		return fmt.Errorf("collection owner cannot be empty")
	}
	if c.MaxSupply <= 0 {
		return fmt.Errorf("max supply must be positive")
	}
	if c.CurrentSupply < 0 {
		return fmt.Errorf("current supply cannot be negative")
	}
	if c.CurrentSupply > c.MaxSupply {
		return fmt.Errorf("current supply cannot exceed max supply")
	}
	return nil
}
