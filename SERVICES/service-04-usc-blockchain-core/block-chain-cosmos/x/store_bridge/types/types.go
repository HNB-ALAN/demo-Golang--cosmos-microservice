package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "store_bridge"

	// RouterKey defines the message route for the bridge module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the bridge module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeBridgeCreated     = "bridge_created"
	EventTypeBridgeUpdated     = "bridge_updated"
	EventTypeBridgeDeleted     = "bridge_deleted"
	EventTypeTransferInitiated = "transfer_initiated"
	EventTypeTransferCompleted = "transfer_completed"
	EventTypeTransferFailed    = "transfer_failed"
	EventTypeValidatorAdded    = "validator_added"
	EventTypeValidatorRemoved  = "validator_removed"
)

// Event attribute keys
const (
	AttributeKeyBridgeID      = "bridge_id"
	AttributeKeyBridgeName    = "bridge_name"
	AttributeKeyTransferID    = "transfer_id"
	AttributeKeyFromChain     = "from_chain"
	AttributeKeyToChain       = "to_chain"
	AttributeKeyAmount        = "amount"
	AttributeKeyToken         = "token"
	AttributeKeyValidatorID   = "validator_id"
	AttributeKeyValidatorAddr = "validator_address"
	AttributeKeyStatus        = "status"
	AttributeKeyTimestamp     = "timestamp"
)

// Bridge represents a cross-chain bridge
type Bridge struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	FromChain   string            `json:"from_chain"`
	ToChain     string            `json:"to_chain"`
	Type        string            `json:"type"`   // token, nft, data
	Status      string            `json:"status"` // active, inactive, maintenance
	Config      map[string]string `json:"config"`
	Validators  []string          `json:"validators"`
	Threshold   int64             `json:"threshold"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Tags        map[string]string `json:"tags"`
}

// Transfer represents a cross-chain transfer
type Transfer struct {
	ID            string            `json:"id"`
	BridgeID      string            `json:"bridge_id"`
	FromChain     string            `json:"from_chain"`
	ToChain       string            `json:"to_chain"`
	FromAddress   string            `json:"from_address"`
	ToAddress     string            `json:"to_address"`
	Amount        string            `json:"amount"`
	Token         string            `json:"token"`
	Status        string            `json:"status"` // pending, confirmed, completed, failed
	TxHash        string            `json:"tx_hash"`
	BlockHeight   int64             `json:"block_height"`
	CreatedAt     time.Time         `json:"created_at"`
	ConfirmedAt   time.Time         `json:"confirmed_at"`
	CompletedAt   time.Time         `json:"completed_at"`
	FailedAt      time.Time         `json:"failed_at"`
	FailureReason string            `json:"failure_reason"`
	Metadata      map[string]string `json:"metadata"`
	Tags          map[string]string `json:"tags"`
}

// Validator represents a bridge validator
type Validator struct {
	ID          string            `json:"id"`
	Address     string            `json:"address"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"` // active, inactive, suspended
	Stake       string            `json:"stake"`
	Commission  string            `json:"commission"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Tags        map[string]string `json:"tags"`
}

// BridgeConfig represents bridge configuration
type BridgeConfig struct {
	ID          string            `json:"id"`
	BridgeID    string            `json:"bridge_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Config      map[string]string `json:"config"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Tags        map[string]string `json:"tags"`
}

// BridgeFee represents bridge fees
type BridgeFee struct {
	ID        string            `json:"id"`
	BridgeID  string            `json:"bridge_id"`
	Token     string            `json:"token"`
	Fee       string            `json:"fee"`
	MinFee    string            `json:"min_fee"`
	MaxFee    string            `json:"max_fee"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Tags      map[string]string `json:"tags"`
}

// BridgeLimit represents bridge limits
type BridgeLimit struct {
	ID         string            `json:"id"`
	BridgeID   string            `json:"bridge_id"`
	Token      string            `json:"token"`
	DailyLimit string            `json:"daily_limit"`
	TxLimit    string            `json:"tx_limit"`
	UserLimit  string            `json:"user_limit"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Tags       map[string]string `json:"tags"`
}

// BridgeEvent represents a bridge event
type BridgeEvent struct {
	ID        string            `json:"id"`
	BridgeID  string            `json:"bridge_id"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	Tags      map[string]string `json:"tags"`
}

// GenesisState represents the genesis state of the bridge module
type GenesisState struct {
	Bridges    []Bridge       `json:"bridges"`
	Transfers  []Transfer     `json:"transfers"`
	Validators []Validator    `json:"validators"`
	Configs    []BridgeConfig `json:"configs"`
	Fees       []BridgeFee    `json:"fees"`
	Limits     []BridgeLimit  `json:"limits"`
	Events     []BridgeEvent  `json:"events"`
	Params     Params         `json:"params"`
}

// Params represents the parameters for the bridge module
type Params struct {
	MaxBridges       int64         `json:"max_bridges"`
	MaxValidators    int64         `json:"max_validators"`
	MinStake         string        `json:"min_stake"`
	MaxStake         string        `json:"max_stake"`
	DefaultThreshold int64         `json:"default_threshold"`
	TransferTimeout  time.Duration `json:"transfer_timeout"`
	FeePercentage    string        `json:"fee_percentage"`
	MaxDailyVolume   string        `json:"max_daily_volume"`
}

// DefaultParams returns the default parameters for the bridge module
func DefaultParams() Params {
	return Params{
		MaxBridges:       100,
		MaxValidators:    50,
		MinStake:         "1000000",   // 1M tokens
		MaxStake:         "100000000", // 100M tokens
		DefaultThreshold: 2,
		TransferTimeout:  24 * time.Hour, // 24 hours
		FeePercentage:    "0.1",          // 0.1%
		MaxDailyVolume:   "1000000000",   // 1B tokens
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxBridges <= 0 {
		return fmt.Errorf("max bridges must be positive")
	}
	if p.MaxValidators <= 0 {
		return fmt.Errorf("max validators must be positive")
	}
	if p.DefaultThreshold <= 0 {
		return fmt.Errorf("default threshold must be positive")
	}
	if p.TransferTimeout <= 0 {
		return fmt.Errorf("transfer timeout must be positive")
	}
	return nil
}

// Validate validates a bridge
func (b Bridge) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if b.Name == "" {
		return fmt.Errorf("bridge name cannot be empty")
	}
	if b.FromChain == "" {
		return fmt.Errorf("from chain cannot be empty")
	}
	if b.ToChain == "" {
		return fmt.Errorf("to chain cannot be empty")
	}
	if b.Type != "token" && b.Type != "nft" && b.Type != "data" {
		return fmt.Errorf("invalid bridge type: %s", b.Type)
	}
	if b.Status != "active" && b.Status != "inactive" && b.Status != "maintenance" {
		return fmt.Errorf("invalid bridge status: %s", b.Status)
	}
	if b.Threshold <= 0 {
		return fmt.Errorf("threshold must be positive")
	}
	return nil
}

// Validate validates a transfer
func (t Transfer) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("transfer ID cannot be empty")
	}
	if t.BridgeID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if t.FromChain == "" {
		return fmt.Errorf("from chain cannot be empty")
	}
	if t.ToChain == "" {
		return fmt.Errorf("to chain cannot be empty")
	}
	if t.FromAddress == "" {
		return fmt.Errorf("from address cannot be empty")
	}
	if t.ToAddress == "" {
		return fmt.Errorf("to address cannot be empty")
	}
	if t.Amount == "" {
		return fmt.Errorf("amount cannot be empty")
	}
	if t.Status != "pending" && t.Status != "confirmed" && t.Status != "completed" && t.Status != "failed" {
		return fmt.Errorf("invalid transfer status: %s", t.Status)
	}
	return nil
}

// Validate validates a validator
func (v Validator) Validate() error {
	if v.ID == "" {
		return fmt.Errorf("validator ID cannot be empty")
	}
	if v.Address == "" {
		return fmt.Errorf("validator address cannot be empty")
	}
	if v.Name == "" {
		return fmt.Errorf("validator name cannot be empty")
	}
	if v.Status != "active" && v.Status != "inactive" && v.Status != "suspended" {
		return fmt.Errorf("invalid validator status: %s", v.Status)
	}
	return nil
}

// Validate validates a bridge config
func (c BridgeConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("config ID cannot be empty")
	}
	if c.BridgeID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if c.Name == "" {
		return fmt.Errorf("config name cannot be empty")
	}
	return nil
}

// Validate validates a bridge fee
func (f BridgeFee) Validate() error {
	if f.ID == "" {
		return fmt.Errorf("fee ID cannot be empty")
	}
	if f.BridgeID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if f.Token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	if f.Fee == "" {
		return fmt.Errorf("fee cannot be empty")
	}
	return nil
}

// Validate validates a bridge limit
func (l BridgeLimit) Validate() error {
	if l.ID == "" {
		return fmt.Errorf("limit ID cannot be empty")
	}
	if l.BridgeID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if l.Token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	return nil
}

// Validate validates a bridge event
func (e BridgeEvent) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if e.BridgeID == "" {
		return fmt.Errorf("bridge ID cannot be empty")
	}
	if e.Type == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	return nil
}
