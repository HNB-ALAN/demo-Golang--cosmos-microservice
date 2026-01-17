package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// CustomToken module constants
const (
	ModuleName   = "custom_token"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeTokenCreated     = "token_created"
	EventTypeTokenUpdated     = "token_updated"
	EventTypeTokenDeleted     = "token_deleted"
	EventTypeTokenMinted      = "token_minted"
	EventTypeTokenBurned      = "token_burned"
	EventTypeTokenTransferred = "token_transferred"
)

// Event attribute keys
const (
	AttributeKeyTokenID     = "token_id"
	AttributeKeyTokenName   = "token_name"
	AttributeKeyTokenSymbol = "token_symbol"
	AttributeKeyOwner       = "owner"
	AttributeKeyAmount      = "amount"
	AttributeKeyFrom        = "from"
	AttributeKeyTo          = "to"
	AttributeKeyCreatedAt   = "created_at"
	AttributeKeyUpdatedAt   = "updated_at"
)

// CustomToken represents a custom token
type CustomToken struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply string `json:"total_supply"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	Metadata    string `json:"metadata"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// ProtoMessage implements proto.Message interface
func (ct *CustomToken) ProtoMessage() {}

// Reset implements proto.Message interface
func (ct *CustomToken) Reset() {
	*ct = CustomToken{}
}

// String implements proto.Message interface
func (ct *CustomToken) String() string {
	return fmt.Sprintf("CustomToken{ID: %s, Name: %s, Symbol: %s, Decimals: %d, TotalSupply: %s, Owner: %s, Status: %s, Metadata: %s, CreatedAt: %d, UpdatedAt: %d}",
		ct.ID, ct.Name, ct.Symbol, ct.Decimals, ct.TotalSupply, ct.Owner, ct.Status, ct.Metadata, ct.CreatedAt, ct.UpdatedAt)
}

// TokenBalance represents a token balance for an account
type TokenBalance struct {
	TokenID   string `json:"token_id"`
	Owner     string `json:"owner"`
	Amount    string `json:"amount"`
	UpdatedAt int64  `json:"updated_at"`
}

// ProtoMessage implements proto.Message interface
func (tb *TokenBalance) ProtoMessage() {}

// Reset implements proto.Message interface
func (tb *TokenBalance) Reset() {
	*tb = TokenBalance{}
}

// String implements proto.Message interface
func (tb *TokenBalance) String() string {
	return fmt.Sprintf("TokenBalance{TokenID: %s, Owner: %s, Amount: %s, UpdatedAt: %d}",
		tb.TokenID, tb.Owner, tb.Amount, tb.UpdatedAt)
}

// TokenTransfer represents a token transfer
type TokenTransfer struct {
	ID        string `json:"id"`
	TokenID   string `json:"token_id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
	CreatedAt int64  `json:"created_at"`
}

// ProtoMessage implements proto.Message interface
func (tt *TokenTransfer) ProtoMessage() {}

// Reset implements proto.Message interface
func (tt *TokenTransfer) Reset() {
	*tt = TokenTransfer{}
}

// String implements proto.Message interface
func (tt *TokenTransfer) String() string {
	return fmt.Sprintf("TokenTransfer{ID: %s, TokenID: %s, From: %s, To: %s, Amount: %s, CreatedAt: %d}",
		tt.ID, tt.TokenID, tt.From, tt.To, tt.Amount, tt.CreatedAt)
}

// GenesisState represents the genesis state of the custom_token module
type GenesisState struct {
	Tokens    []CustomToken   `json:"tokens"`
	Balances  []TokenBalance  `json:"balances"`
	Transfers []TokenTransfer `json:"transfers"`
	Params    Params          `json:"params"`
}

// ProtoMessage implements proto.Message interface
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements proto.Message interface
func (gs *GenesisState) String() string {
	return fmt.Sprintf("GenesisState{Tokens: %v, Balances: %v, Transfers: %v, Params: %v}",
		gs.Tokens, gs.Balances, gs.Transfers, gs.Params)
}

// Params defines the parameters for the custom_token module
type Params struct {
	MaxTokens       uint32 `json:"max_tokens"`
	MinNameLength   uint32 `json:"min_name_length"`
	MaxNameLength   uint32 `json:"max_name_length"`
	MinSymbolLength uint32 `json:"min_symbol_length"`
	MaxSymbolLength uint32 `json:"max_symbol_length"`
	MaxDecimals     uint8  `json:"max_decimals"`
	MintingFee      string `json:"minting_fee"`
	TransferFee     string `json:"transfer_fee"`
	BurnFee         string `json:"burn_fee"`
}

// DefaultParams returns the default parameters for the custom_token module
func DefaultParams() Params {
	return Params{
		MaxTokens:       10000,
		MinNameLength:   3,
		MaxNameLength:   50,
		MinSymbolLength: 3,
		MaxSymbolLength: 10,
		MaxDecimals:     18,
		MintingFee:      "1000000", // 1 USC
		TransferFee:     "100000",  // 0.1 USC
		BurnFee:         "50000",   // 0.05 USC
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxTokens == 0 {
		return fmt.Errorf("max tokens must be positive")
	}
	if p.MinNameLength == 0 {
		return fmt.Errorf("min name length must be positive")
	}
	if p.MaxNameLength < p.MinNameLength {
		return fmt.Errorf("max name length must be greater than min name length")
	}
	if p.MinSymbolLength == 0 {
		return fmt.Errorf("min symbol length must be positive")
	}
	if p.MaxSymbolLength < p.MinSymbolLength {
		return fmt.Errorf("max symbol length must be greater than min symbol length")
	}
	if p.MaxDecimals > 18 {
		return fmt.Errorf("max decimals cannot exceed 18")
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
	return nil
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("MaxTokens"), &p.MaxTokens, validateUint32),
		paramtypes.NewParamSetPair([]byte("MinNameLength"), &p.MinNameLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("MaxNameLength"), &p.MaxNameLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("MinSymbolLength"), &p.MinSymbolLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("MaxSymbolLength"), &p.MaxSymbolLength, validateUint32),
		paramtypes.NewParamSetPair([]byte("MaxDecimals"), &p.MaxDecimals, validateUint8),
		paramtypes.NewParamSetPair([]byte("MintingFee"), &p.MintingFee, validateString),
		paramtypes.NewParamSetPair([]byte("TransferFee"), &p.TransferFee, validateString),
		paramtypes.NewParamSetPair([]byte("BurnFee"), &p.BurnFee, validateString),
	}
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

func validateUint8(i interface{}) error {
	_, ok := i.(uint8)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// NewCustomToken creates a new CustomToken
func NewCustomToken(id, name, symbol string, decimals uint8, totalSupply, owner, status, metadata string) CustomToken {
	now := time.Now().Unix()
	return CustomToken{
		ID:          id,
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Owner:       owner,
		Status:      status,
		Metadata:    metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewTokenBalance creates a new TokenBalance
func NewTokenBalance(tokenID, owner, amount string) TokenBalance {
	return TokenBalance{
		TokenID:   tokenID,
		Owner:     owner,
		Amount:    amount,
		UpdatedAt: time.Now().Unix(),
	}
}

// NewTokenTransfer creates a new TokenTransfer
func NewTokenTransfer(id, tokenID, from, to, amount string) TokenTransfer {
	return TokenTransfer{
		ID:        id,
		TokenID:   tokenID,
		From:      from,
		To:        to,
		Amount:    amount,
		CreatedAt: time.Now().Unix(),
	}
}

// DefaultGenesis returns the default genesis state for the custom_token module
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Tokens:    []CustomToken{},
		Balances:  []TokenBalance{},
		Transfers: []TokenTransfer{},
		Params:    DefaultParams(),
	}
}

// Validate validates the genesis state
func (gs GenesisState) Validate() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate tokens
	seenIDs := make(map[string]bool)
	for _, token := range gs.Tokens {
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

		// Check for duplicate IDs
		if seenIDs[token.ID] {
			return fmt.Errorf("duplicate token ID: %s", token.ID)
		}
		seenIDs[token.ID] = true
	}

	// Validate balances
	for _, balance := range gs.Balances {
		if balance.TokenID == "" {
			return fmt.Errorf("balance token ID cannot be empty")
		}
		if balance.Owner == "" {
			return fmt.Errorf("balance owner cannot be empty")
		}
		if balance.Amount == "" {
			return fmt.Errorf("balance amount cannot be empty")
		}
		if balance.UpdatedAt <= 0 {
			return fmt.Errorf("balance timestamp must be positive")
		}
	}

	// Validate transfers
	for _, transfer := range gs.Transfers {
		if transfer.ID == "" {
			return fmt.Errorf("transfer ID cannot be empty")
		}
		if transfer.TokenID == "" {
			return fmt.Errorf("transfer token ID cannot be empty")
		}
		if transfer.From == "" {
			return fmt.Errorf("transfer from cannot be empty")
		}
		if transfer.To == "" {
			return fmt.Errorf("transfer to cannot be empty")
		}
		if transfer.Amount == "" {
			return fmt.Errorf("transfer amount cannot be empty")
		}
		if transfer.CreatedAt <= 0 {
			return fmt.Errorf("transfer timestamp must be positive")
		}
	}

	return nil
}
