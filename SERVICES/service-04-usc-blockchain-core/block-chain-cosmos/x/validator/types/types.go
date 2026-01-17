package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Validator module constants
const (
	ModuleName   = "validator"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeValidatorCreate  = "validator_create"
	EventTypeValidatorUpdate  = "validator_update"
	EventTypeValidatorRemove  = "validator_remove"
	EventTypeDelegationCreate = "delegation_create"
	EventTypeDelegationUpdate = "delegation_update"
	EventTypeDelegationRemove = "delegation_remove"
	EventTypeBlockStart       = "block_start"
	EventTypeBlockEnd         = "block_end"
)

// Event attribute keys
const (
	AttributeKeyValidatorAddress = "validator_address"
	AttributeKeyDelegatorAddress = "delegator_address"
	AttributeKeyAmount           = "amount"
	AttributeKeyBlockHeight      = "block_height"
	AttributeKeyBlockTime        = "block_time"
)

// Validator represents a validator
type Validator struct {
	Address     string `json:"address"`
	PubKey      string `json:"pub_key"`
	Power       int64  `json:"power"`
	Description string `json:"description"`
	Commission  string `json:"commission"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

// Delegation represents a delegation
type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           string `json:"amount"`
	CreatedAt        int64  `json:"created_at"`
}

// GenesisState represents the genesis state of the validator module
type GenesisState struct {
	Validators  []Validator  `json:"validators"`
	Delegations []Delegation `json:"delegations"`
	Params      Params       `json:"params"`
}

// ProtoMessage implements proto.Message interface
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements proto.Message interface
func (gs *GenesisState) String() string {
	return fmt.Sprintf("GenesisState{Validators: %v, Delegations: %v, Params: %v}",
		gs.Validators, gs.Delegations, gs.Params)
}

// Params defines the parameters for the validator module
type Params struct {
	MaxValidators   uint32 `json:"max_validators"`
	MinDelegation   string `json:"min_delegation"`
	MaxCommission   string `json:"max_commission"`
	UnbondingTime   int64  `json:"unbonding_time"`
	SlashingEnabled bool   `json:"slashing_enabled"`
	JailingEnabled  bool   `json:"jailing_enabled"`
}

// DefaultParams returns default parameters for the validator module
func DefaultParams() Params {
	return Params{
		MaxValidators:   100,
		MinDelegation:   "1000000",  // 1 USC
		MaxCommission:   "0.20",     // 20%
		UnbondingTime:   86400 * 21, // 21 days
		SlashingEnabled: true,
		JailingEnabled:  true,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxValidators == 0 {
		return fmt.Errorf("max validators must be positive")
	}
	if p.MinDelegation == "" {
		return fmt.Errorf("min delegation cannot be empty")
	}
	if p.MaxCommission == "" {
		return fmt.Errorf("max commission cannot be empty")
	}
	if p.UnbondingTime <= 0 {
		return fmt.Errorf("unbonding time must be positive")
	}
	return nil
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("MaxValidators"), &p.MaxValidators, validateUint32),
		paramtypes.NewParamSetPair([]byte("MinDelegation"), &p.MinDelegation, validateString),
		paramtypes.NewParamSetPair([]byte("MaxCommission"), &p.MaxCommission, validateString),
		paramtypes.NewParamSetPair([]byte("UnbondingTime"), &p.UnbondingTime, validateInt64),
		paramtypes.NewParamSetPair([]byte("SlashingEnabled"), &p.SlashingEnabled, validateBool),
		paramtypes.NewParamSetPair([]byte("JailingEnabled"), &p.JailingEnabled, validateBool),
	}
}

// ParamKeyTable returns the parameter key table
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Validation functions
func validateString(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUint32(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateInt64(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// NewValidator creates a new validator
func NewValidator(address, pubKey, description, commission string) Validator {
	return Validator{
		Address:     address,
		PubKey:      pubKey,
		Power:       0,
		Description: description,
		Commission:  commission,
		Status:      "active",
		CreatedAt:   time.Now().Unix(),
	}
}

// NewDelegation creates a new delegation
func NewDelegation(delegatorAddress, validatorAddress, amount string) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		CreatedAt:        time.Now().Unix(),
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Validators:  []Validator{},
		Delegations: []Delegation{},
		Params:      DefaultParams(),
	}
}

// Validate validates the genesis state
func (gs GenesisState) Validate() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate validators
	seenAddresses := make(map[string]bool)
	for _, validator := range gs.Validators {
		if validator.Address == "" {
			return fmt.Errorf("validator address cannot be empty")
		}
		if validator.PubKey == "" {
			return fmt.Errorf("validator pub key cannot be empty")
		}
		if validator.Description == "" {
			return fmt.Errorf("validator description cannot be empty")
		}
		if validator.Commission == "" {
			return fmt.Errorf("validator commission cannot be empty")
		}

		// Check for duplicate addresses
		if seenAddresses[validator.Address] {
			return fmt.Errorf("duplicate validator address: %s", validator.Address)
		}
		seenAddresses[validator.Address] = true
	}

	// Validate delegations
	for _, delegation := range gs.Delegations {
		if delegation.DelegatorAddress == "" {
			return fmt.Errorf("delegation delegator address cannot be empty")
		}
		if delegation.ValidatorAddress == "" {
			return fmt.Errorf("delegation validator address cannot be empty")
		}
		if delegation.Amount == "" {
			return fmt.Errorf("delegation amount cannot be empty")
		}
		if delegation.CreatedAt <= 0 {
			return fmt.Errorf("delegation timestamp must be positive")
		}
	}

	return nil
}
