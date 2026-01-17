package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/transaction/v1/usc/transaction/v1"
)

const (
	// ModuleName defines the module name
	ModuleName = "transaction"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Key prefixes for store
const (
	TransactionKeyPrefix = "transaction:"
	StatsKeyPrefix       = "stats:"
)

// Event types
const (
	EventTypeCreateTransaction   = "create_transaction"
	EventTypeValidateTransaction = "validate_transaction"
	EventTypeExecuteTransaction  = "execute_transaction"
	EventTypeCancelTransaction   = "cancel_transaction"
)

// Event attribute keys
const (
	AttributeKeyTransactionHash = "transaction_hash"
	AttributeKeyTransactionID   = "transaction_id"
	AttributeKeyFromAddress     = "from_address"
	AttributeKeyToAddress       = "to_address"
	AttributeKeyAmount          = "amount"
	AttributeKeyTransactionType = "transaction_type"
	AttributeKeyStatus          = "status"
	AttributeKeyValidator       = "validator"
	AttributeKeyExecutor        = "executor"
	AttributeKeyCanceller       = "canceller"
)

// Transaction represents a blockchain transaction
// COSMOS SDK 0.53.4: Must have protobuf tags if implementing ProtoMessage()
type Transaction struct {
	Hash            string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	ID              string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	FromAddress     string `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToAddress       string `protobuf:"bytes,4,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`
	Amount          string `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
	TransactionType string `protobuf:"bytes,6,opt,name=transaction_type,json=transactionType,proto3" json:"transaction_type,omitempty"`
	Status          string `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"` // "pending", "validated", "executed", "failed", "cancelled"
	Data            string `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`     // JSON string
	Memo            string `protobuf:"bytes,9,opt,name=memo,proto3" json:"memo,omitempty"`
	CreatedAt       int64  `protobuf:"varint,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ValidatedAt     int64  `protobuf:"varint,11,opt,name=validated_at,json=validatedAt,proto3" json:"validated_at,omitempty"`
	ExecutedAt      int64  `protobuf:"varint,12,opt,name=executed_at,json=executedAt,proto3" json:"executed_at,omitempty"`
	FailedAt        int64  `protobuf:"varint,13,opt,name=failed_at,json=failedAt,proto3" json:"failed_at,omitempty"`
	ValidationProof string `protobuf:"bytes,14,opt,name=validation_proof,json=validationProof,proto3" json:"validation_proof,omitempty"`
	ExecutionProof  string `protobuf:"bytes,15,opt,name=execution_proof,json=executionProof,proto3" json:"execution_proof,omitempty"`
	ErrorMessage    string `protobuf:"bytes,16,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	GasUsed         int64  `protobuf:"varint,17,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	GasLimit        int64  `protobuf:"varint,18,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	Fee             string `protobuf:"bytes,19,opt,name=fee,proto3" json:"fee,omitempty"`
}

// ProtoMessage implements proto.Message interface
func (t *Transaction) ProtoMessage() {}

// Reset implements proto.Message interface
func (t *Transaction) Reset() {
	*t = Transaction{}
}

// String implements proto.Message interface
func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction{Hash: %s, ID: %s, From: %s, To: %s, Amount: %s, Type: %s, Status: %s, CreatedAt: %d}",
		t.Hash, t.ID, t.FromAddress, t.ToAddress, t.Amount, t.TransactionType, t.Status, t.CreatedAt)
}

// TransactionStats represents transaction statistics
type TransactionStats struct {
	TotalTransactions       int64  `json:"total_transactions"`
	PendingTransactions     int64  `json:"pending_transactions"`
	ValidatedTransactions   int64  `json:"validated_transactions"`
	ExecutedTransactions    int64  `json:"executed_transactions"`
	FailedTransactions      int64  `json:"failed_transactions"`
	CancelledTransactions   int64  `json:"cancelled_transactions"`
	TotalVolume             string `json:"total_volume"`
	AverageTransactionValue string `json:"average_transaction_value"`
	SuccessRate             string `json:"success_rate"`
	AverageExecutionTime    int64  `json:"average_execution_time"`
	CurrentHeight           int64  `json:"current_height"`
	LastTransactionTime     string `json:"last_transaction_time"`
}

// ProtoMessage implements proto.Message interface
func (ts *TransactionStats) ProtoMessage() {}

// Reset implements proto.Message interface
func (ts *TransactionStats) Reset() {
	*ts = TransactionStats{}
}

// String implements proto.Message interface
func (ts *TransactionStats) String() string {
	return fmt.Sprintf("TransactionStats{Total: %d, Pending: %d, Executed: %d, SuccessRate: %s}",
		ts.TotalTransactions, ts.PendingTransactions, ts.ExecutedTransactions, ts.SuccessRate)
}

// GenesisState defines the transaction module's genesis state
type GenesisState struct {
	Transactions     []Transaction    `json:"transactions"`
	TransactionStats TransactionStats `json:"transaction_stats"`
	Params           Params           `json:"params"`
}

// ProtoMessage implements proto.Message interface
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements proto.Message interface
func (gs *GenesisState) String() string {
	return fmt.Sprintf("GenesisState{Transactions: %d, Stats: %s, Params: %s}",
		len(gs.Transactions), gs.TransactionStats.String(), gs.Params.String())
}

// Params defines the parameters for the transaction module
type Params struct {
	MaxTransactionSize     uint32 `json:"max_transaction_size"`
	MinTransactionFee      string `json:"min_transaction_fee"`
	MaxTransactionFee      string `json:"max_transaction_fee"`
	TransactionTimeout     int64  `json:"transaction_timeout"`
	EnableTransactionCache bool   `json:"enable_transaction_cache"`
	MaxCacheSize           uint32 `json:"max_cache_size"`
}

// DefaultParams returns default parameters for the transaction module
var DefaultParams = Params{
	MaxTransactionSize:     1024 * 1024, // 1MB
	MinTransactionFee:      "1000",
	MaxTransactionFee:      "1000000",
	TransactionTimeout:     300, // 5 minutes
	EnableTransactionCache: true,
	MaxCacheSize:           10000,
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("MaxTransactionSize"), &p.MaxTransactionSize, validateUint32),
		paramtypes.NewParamSetPair([]byte("MinTransactionFee"), &p.MinTransactionFee, validateString),
		paramtypes.NewParamSetPair([]byte("MaxTransactionFee"), &p.MaxTransactionFee, validateString),
		paramtypes.NewParamSetPair([]byte("TransactionTimeout"), &p.TransactionTimeout, validateInt64),
		paramtypes.NewParamSetPair([]byte("EnableTransactionCache"), &p.EnableTransactionCache, validateBool),
		paramtypes.NewParamSetPair([]byte("MaxCacheSize"), &p.MaxCacheSize, validateUint32),
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxTransactionSize == 0 {
		return fmt.Errorf("max transaction size must be positive")
	}
	if p.TransactionTimeout <= 0 {
		return fmt.Errorf("transaction timeout must be positive")
	}
	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	return fmt.Sprintf("Params{MaxSize: %d, MinFee: %s, MaxFee: %s, Timeout: %d, Cache: %t, CacheSize: %d}",
		p.MaxTransactionSize, p.MinTransactionFee, p.MaxTransactionFee, p.TransactionTimeout, p.EnableTransactionCache, p.MaxCacheSize)
}

// NewTransaction creates a new transaction
func NewTransaction(hash, id, from, to, amount, txType, data, memo string) Transaction {
	now := time.Now().Unix()
	return Transaction{
		Hash:            hash,
		ID:              id,
		FromAddress:     from,
		ToAddress:       to,
		Amount:          amount,
		TransactionType: txType,
		Status:          "pending",
		Data:            data,
		Memo:            memo,
		CreatedAt:       now,
		GasUsed:         0,
		GasLimit:        100000,
		Fee:             "1000",
	}
}

// NewTransactionStats creates new transaction statistics
func NewTransactionStats() TransactionStats {
	return TransactionStats{
		TotalTransactions:       0,
		PendingTransactions:     0,
		ValidatedTransactions:   0,
		ExecutedTransactions:    0,
		FailedTransactions:      0,
		CancelledTransactions:   0,
		TotalVolume:             "0",
		AverageTransactionValue: "0",
		SuccessRate:             "0.00",
		AverageExecutionTime:    0,
		CurrentHeight:           0,
		LastTransactionTime:     time.Now().Format(time.RFC3339),
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Transactions:     []Transaction{},
		TransactionStats: NewTransactionStats(),
		Params:           DefaultParams,
	}
}

// Validate validates the genesis state
func (gs GenesisState) Validate() error {
	for _, tx := range gs.Transactions {
		if tx.Hash == "" {
			return fmt.Errorf("transaction hash cannot be empty")
		}
		if tx.FromAddress == "" {
			return fmt.Errorf("transaction from address cannot be empty")
		}
	}
	return gs.Params.Validate()
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

// ParamKeyTable returns the parameter key table for the transaction module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx interface{}, addr interface{}) interface{}
}

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	SendCoins(ctx interface{}, fromAddr, toAddr interface{}, amt interface{}) error
	GetBalance(ctx interface{}, addr interface{}, denom string) interface{}
}

// ToBlockchainProto converts Transaction to blockchain-proto Transaction
func (t Transaction) ToBlockchainProto() *blockchainproto.Transaction {
	// Convert string status to enum
	var status blockchainproto.TransactionStatus
	switch t.Status {
	case "pending":
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_PENDING
	case "validated":
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_VALIDATED
	case "executed":
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_EXECUTED
	case "failed":
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_FAILED
	case "cancelled":
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_CANCELLED
	default:
		status = blockchainproto.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
	}

	// Convert string type to enum
	var txType blockchainproto.TransactionType
	switch t.TransactionType {
	case "transfer":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_TRANSFER
	case "mint":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_MINT
	case "burn":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_BURN
	default:
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}

	return &blockchainproto.Transaction{
		Hash:            t.Hash,
		Id:              t.ID,
		FromAddress:     t.FromAddress,
		ToAddress:       t.ToAddress,
		Amount:          &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse t.Amount
		TransactionType: txType,
		Status:          status,
		Data:            &blockchainproto.TransactionData{}, // TODO: parse t.Data
		Memo:            t.Memo,
		CreatedAt:       timestamppb.New(time.Unix(t.CreatedAt, 0)),
		ValidatedAt:     timestamppb.New(time.Unix(t.ValidatedAt, 0)),
		ExecutedAt:      timestamppb.New(time.Unix(t.ExecutedAt, 0)),
		FailedAt:        timestamppb.New(time.Unix(t.FailedAt, 0)),
		ValidationProof: t.ValidationProof,
		ExecutionProof:  t.ExecutionProof,
		ErrorMessage:    t.ErrorMessage,
		GasUsed:         t.GasUsed,
		GasLimit:        uint64(t.GasLimit),
		Fee:             &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // TODO: parse t.Fee
	}
}

// NewTransactionFromBlockchainProto creates Transaction from blockchain-proto
func NewTransactionFromBlockchainProto(tx *blockchainproto.Transaction) Transaction {
	return Transaction{
		Hash:            tx.Hash,
		ID:              tx.Id,
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		Amount:          tx.Amount.String(), // TODO: proper conversion
		TransactionType: tx.TransactionType.String(),
		Status:          tx.Status.String(),
		Data:            "", // TODO: convert tx.Data
		Memo:            tx.Memo,
		CreatedAt:       tx.CreatedAt.Seconds,
		ValidatedAt:     tx.ValidatedAt.Seconds,
		ExecutedAt:      tx.ExecutedAt.Seconds,
		FailedAt:        tx.FailedAt.Seconds,
		ValidationProof: tx.ValidationProof,
		ExecutionProof:  tx.ExecutionProof,
		ErrorMessage:    tx.ErrorMessage,
		GasUsed:         int64(tx.GasUsed),
		GasLimit:        int64(tx.GasLimit),
		Fee:             tx.Fee.String(), // TODO: proper conversion
	}
}
