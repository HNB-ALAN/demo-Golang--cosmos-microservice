package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "smart_contract"

	// RouterKey defines the message route for the contract module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the contract module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeContractCreated   = "contract_created"
	EventTypeContractUpdated   = "contract_updated"
	EventTypeContractDeleted   = "contract_deleted"
	EventTypeContractExecuted  = "contract_executed"
	EventTypeContractDeployed  = "contract_deployed"
	EventTypeContractUpgraded  = "contract_upgraded"
	EventTypeContractPaused    = "contract_paused"
	EventTypeContractResumed   = "contract_resumed"
	EventTypeContractMigrated  = "contract_migrated"
	EventTypeContractDestroyed = "contract_destroyed"

	AttributeKeyContractID   = "contract_id"
	AttributeKeyContractType = "contract_type"
	AttributeKeyOwner        = "owner"
	AttributeKeyExecutor     = "executor"
	AttributeKeyModule       = ModuleName
)

// ContractType represents the type of smart contract
type ContractType string

const (
	ContractTypeWASM   ContractType = "wasm"
	ContractTypeEVM    ContractType = "evm"
	ContractTypeNative ContractType = "native"
	ContractTypeCustom ContractType = "custom"
)

// ContractStatus represents the status of a smart contract
type ContractStatus string

const (
	ContractStatusActive  ContractStatus = "active"
	ContractStatusPaused  ContractStatus = "paused"
	ContractStatusUpgrade ContractStatus = "upgrade"
	ContractStatusDeleted ContractStatus = "deleted"
)

// SmartContract represents a smart contract
type SmartContract struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        ContractType      `json:"type"`
	Status      ContractStatus    `json:"status"`
	Owner       string            `json:"owner"`
	CodeHash    string            `json:"code_hash"`
	Code        []byte            `json:"code"`
	ABI         string            `json:"abi"`
	Bytecode    []byte            `json:"bytecode"`
	Address     string            `json:"address"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// ContractExecution represents a contract execution
type ContractExecution struct {
	ID         string            `json:"id"`
	ContractID string            `json:"contract_id"`
	Executor   string            `json:"executor"`
	Method     string            `json:"method"`
	Input      []byte            `json:"input"`
	Output     []byte            `json:"output"`
	GasUsed    uint64            `json:"gas_used"`
	GasLimit   uint64            `json:"gas_limit"`
	Status     string            `json:"status"` // success, failure, error
	Error      string            `json:"error,omitempty"`
	ExecutedAt time.Time         `json:"executed_at"`
	Metadata   map[string]string `json:"metadata"`
}

// ContractDeployment represents a contract deployment
type ContractDeployment struct {
	ID          string            `json:"id"`
	ContractID  string            `json:"contract_id"`
	Deployer    string            `json:"deployer"`
	Network     string            `json:"network"`
	Address     string            `json:"address"`
	TxHash      string            `json:"tx_hash"`
	BlockNumber uint64            `json:"block_number"`
	GasUsed     uint64            `json:"gas_used"`
	DeployedAt  time.Time         `json:"deployed_at"`
	Metadata    map[string]string `json:"metadata"`
}

// ContractUpgrade represents a contract upgrade
type ContractUpgrade struct {
	ID         string            `json:"id"`
	ContractID string            `json:"contract_id"`
	Upgrader   string            `json:"upgrader"`
	OldVersion string            `json:"old_version"`
	NewVersion string            `json:"new_version"`
	CodeHash   string            `json:"code_hash"`
	TxHash     string            `json:"tx_hash"`
	UpgradedAt time.Time         `json:"upgraded_at"`
	Metadata   map[string]string `json:"metadata"`
}

// ContractMigration represents a contract migration
type ContractMigration struct {
	ID          string            `json:"id"`
	ContractID  string            `json:"contract_id"`
	Migrator    string            `json:"migrator"`
	FromNetwork string            `json:"from_network"`
	ToNetwork   string            `json:"to_network"`
	FromAddress string            `json:"from_address"`
	ToAddress   string            `json:"to_address"`
	TxHash      string            `json:"tx_hash"`
	MigratedAt  time.Time         `json:"migrated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// GenesisState defines the contract module's genesis state
type GenesisState struct {
	Contracts   []SmartContract      `json:"contracts"`
	Executions  []ContractExecution  `json:"executions"`
	Deployments []ContractDeployment `json:"deployments"`
	Upgrades    []ContractUpgrade    `json:"upgrades"`
	Migrations  []ContractMigration  `json:"migrations"`
	Params      Params               `json:"params"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Contracts:   []SmartContract{},
		Executions:  []ContractExecution{},
		Deployments: []ContractDeployment{},
		Upgrades:    []ContractUpgrade{},
		Migrations:  []ContractMigration{},
		Params:      DefaultParams(),
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState()
}

// Params defines the parameters for the contract module
type Params struct {
	MaxContracts      uint32        `json:"max_contracts"`
	MaxExecutions     uint32        `json:"max_executions"`
	MaxDeployments    uint32        `json:"max_deployments"`
	MaxUpgrades       uint32        `json:"max_upgrades"`
	MaxMigrations     uint32        `json:"max_migrations"`
	GasLimit          uint64        `json:"gas_limit"`
	ExecutionTimeout  time.Duration `json:"execution_timeout"`
	DeploymentTimeout time.Duration `json:"deployment_timeout"`
	UpgradeTimeout    time.Duration `json:"upgrade_timeout"`
	MigrationTimeout  time.Duration `json:"migration_timeout"`
}

// NewParams creates a new Params object
func NewParams(
	maxContracts uint32,
	maxExecutions uint32,
	maxDeployments uint32,
	maxUpgrades uint32,
	maxMigrations uint32,
	gasLimit uint64,
	executionTimeout time.Duration,
	deploymentTimeout time.Duration,
	upgradeTimeout time.Duration,
	migrationTimeout time.Duration,
) Params {
	return Params{
		MaxContracts:      maxContracts,
		MaxExecutions:     maxExecutions,
		MaxDeployments:    maxDeployments,
		MaxUpgrades:       maxUpgrades,
		MaxMigrations:     maxMigrations,
		GasLimit:          gasLimit,
		ExecutionTimeout:  executionTimeout,
		DeploymentTimeout: deploymentTimeout,
		UpgradeTimeout:    upgradeTimeout,
		MigrationTimeout:  migrationTimeout,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		1000,           // MaxContracts
		10000,          // MaxExecutions
		1000,           // MaxDeployments
		500,            // MaxUpgrades
		200,            // MaxMigrations
		10000000,       // GasLimit (10M gas)
		30*time.Second, // ExecutionTimeout
		5*time.Minute,  // DeploymentTimeout
		2*time.Minute,  // UpgradeTimeout
		10*time.Minute, // MigrationTimeout
	)
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxContracts == 0 {
		return fmt.Errorf("max contracts cannot be zero")
	}
	if p.MaxExecutions == 0 {
		return fmt.Errorf("max executions cannot be zero")
	}
	if p.MaxDeployments == 0 {
		return fmt.Errorf("max deployments cannot be zero")
	}
	if p.MaxUpgrades == 0 {
		return fmt.Errorf("max upgrades cannot be zero")
	}
	if p.MaxMigrations == 0 {
		return fmt.Errorf("max migrations cannot be zero")
	}
	if p.GasLimit == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}
	if p.ExecutionTimeout <= 0 {
		return fmt.Errorf("execution timeout must be positive")
	}
	if p.DeploymentTimeout <= 0 {
		return fmt.Errorf("deployment timeout must be positive")
	}
	if p.UpgradeTimeout <= 0 {
		return fmt.Errorf("upgrade timeout must be positive")
	}
	if p.MigrationTimeout <= 0 {
		return fmt.Errorf("migration timeout must be positive")
	}
	return nil
}

// Validate validates a SmartContract
func (c SmartContract) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if c.Name == "" {
		return fmt.Errorf("contract name cannot be empty")
	}
	if c.Owner == "" {
		return fmt.Errorf("contract owner cannot be empty")
	}
	if c.CodeHash == "" {
		return fmt.Errorf("contract code hash cannot be empty")
	}
	if len(c.Code) == 0 {
		return fmt.Errorf("contract code cannot be empty")
	}
	return nil
}

// Validate validates a ContractExecution
func (e ContractExecution) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("execution ID cannot be empty")
	}
	if e.ContractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if e.Executor == "" {
		return fmt.Errorf("executor cannot be empty")
	}
	if e.Method == "" {
		return fmt.Errorf("method cannot be empty")
	}
	return nil
}

// Validate validates a ContractDeployment
func (d ContractDeployment) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("deployment ID cannot be empty")
	}
	if d.ContractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if d.Deployer == "" {
		return fmt.Errorf("deployer cannot be empty")
	}
	if d.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	return nil
}

// Validate validates a ContractUpgrade
func (u ContractUpgrade) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("upgrade ID cannot be empty")
	}
	if u.ContractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if u.Upgrader == "" {
		return fmt.Errorf("upgrader cannot be empty")
	}
	if u.OldVersion == "" {
		return fmt.Errorf("old version cannot be empty")
	}
	if u.NewVersion == "" {
		return fmt.Errorf("new version cannot be empty")
	}
	return nil
}

// Validate validates a ContractMigration
func (m ContractMigration) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("migration ID cannot be empty")
	}
	if m.ContractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if m.Migrator == "" {
		return fmt.Errorf("migrator cannot be empty")
	}
	if m.FromNetwork == "" {
		return fmt.Errorf("from network cannot be empty")
	}
	if m.ToNetwork == "" {
		return fmt.Errorf("to network cannot be empty")
	}
	return nil
}
