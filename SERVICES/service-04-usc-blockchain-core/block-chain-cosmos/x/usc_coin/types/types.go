package types

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
)

// USC module constants
const (
	ModuleName   = "usc_coin"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeTransfer   = "transfer"
	EventTypeMint       = "mint"
	EventTypeBurn       = "burn"
	EventTypeBlockStart = "block_start"
	EventTypeBlockEnd   = "block_end"
)

// Event attribute keys
const (
	AttributeKeyFrom        = "from"
	AttributeKeyTo          = "to"
	AttributeKeyAmount      = "amount"
	AttributeKeyBlockHeight = "block_height"
	AttributeKeyBlockTime   = "block_time"
)

// USC token types
type USC struct {
	Denom    string `json:"denom"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

// Balance represents a user's USC balance
type Balance struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
	Denom   string `json:"denom"`
}

// ToBlockchainProto converts Balance to blockchain-proto USCHolder
func (b Balance) ToBlockchainProto() *blockchainproto.USCHolder {
	return &blockchainproto.USCHolder{
		Address: b.Address,
		Balance: &sdk.Coin{Denom: b.Denom, Amount: math.NewInt(0)}, // TODO: parse b.Amount
	}
}

// Transfer represents a USC transfer
type Transfer struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
	Denom       string `json:"denom"`
	Timestamp   int64  `json:"timestamp"`
}

// GenesisState represents the genesis state of the USC module
type GenesisState struct {
	Balances  []Balance  `json:"balances"`
	Transfers []Transfer `json:"transfers"`
	Params    Params     `json:"params"`
}

// Params defines the parameters for the USC module
type Params struct {
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals uint8  `json:"token_decimals"`
	MaxSupply     string `json:"max_supply"`
	MintEnabled   bool   `json:"mint_enabled"`
	BurnEnabled   bool   `json:"burn_enabled"`
}

// DefaultParams returns default parameters for the USC module
func DefaultParams() Params {
	return Params{
		TokenName:     "Universal Social Coin",
		TokenSymbol:   "USC",
		TokenDecimals: 18,
		MaxSupply:     "10000000000000000000000000000", // 10B USC
		MintEnabled:   true,
		BurnEnabled:   true,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.TokenName == "" {
		return fmt.Errorf("token name cannot be empty")
	}
	if p.TokenSymbol == "" {
		return fmt.Errorf("token symbol cannot be empty")
	}
	if p.TokenDecimals > 18 {
		return fmt.Errorf("token decimals cannot exceed 18")
	}
	return nil
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("TokenName"), &p.TokenName, validateString),
		paramtypes.NewParamSetPair([]byte("TokenSymbol"), &p.TokenSymbol, validateString),
		paramtypes.NewParamSetPair([]byte("TokenDecimals"), &p.TokenDecimals, validateUint8),
		paramtypes.NewParamSetPair([]byte("MaxSupply"), &p.MaxSupply, validateString),
		paramtypes.NewParamSetPair([]byte("MintEnabled"), &p.MintEnabled, validateBool),
		paramtypes.NewParamSetPair([]byte("BurnEnabled"), &p.BurnEnabled, validateBool),
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

func validateUint8(i interface{}) error {
	_, ok := i.(uint8)
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

// NewTransfer creates a new transfer
func NewTransfer(fromAddress, toAddress, amount, denom string) Transfer {
	return Transfer{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
		Denom:       denom,
		Timestamp:   time.Now().Unix(),
	}
}

// BankKeeper defines the expected interface for the bank keeper
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// RegisterMsgServer registers the msg server
func RegisterMsgServer(server interface{}, impl interface{}) {
	// Implementation will be added when needed
}

// RegisterQueryServer registers the query server
func RegisterQueryServer(server interface{}, impl interface{}) {
	// Implementation will be added when needed
}

// NewUSCHolderFromBalance creates a USCHolder from Balance
func NewUSCHolderFromBalance(balance Balance) *blockchainproto.USCHolder {
	return &blockchainproto.USCHolder{
		Address: balance.Address,
		Balance: &sdk.Coin{Denom: balance.Denom, Amount: math.NewInt(0)}, // TODO: parse balance.Amount
	}
}

// NewCoinFromString creates a sdk.Coin from string amount
func NewCoinFromString(amount, denom string) sdk.Coin {
	// TODO: implement proper string to sdk.Int parsing
	return sdk.Coin{Denom: denom, Amount: math.NewInt(0)}
}
