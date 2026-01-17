package types

import (
	"bytes"
)

// Store key prefixes
var (
	ContractKeyPrefix   = []byte{0x01}
	ExecutionKeyPrefix  = []byte{0x02}
	DeploymentKeyPrefix = []byte{0x03}
	UpgradeKeyPrefix    = []byte{0x04}
	MigrationKeyPrefix  = []byte{0x05}
	ParamsKey           = []byte{0x06}
)

// ContractKey returns the store key for a SmartContract object
func ContractKey(id string) []byte {
	return bytes.Join([][]byte{
		ContractKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// ExecutionKey returns the store key for a ContractExecution object
func ExecutionKey(id string) []byte {
	return bytes.Join([][]byte{
		ExecutionKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// DeploymentKey returns the store key for a ContractDeployment object
func DeploymentKey(id string) []byte {
	return bytes.Join([][]byte{
		DeploymentKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// UpgradeKey returns the store key for a ContractUpgrade object
func UpgradeKey(id string) []byte {
	return bytes.Join([][]byte{
		UpgradeKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// MigrationKey returns the store key for a ContractMigration object
func MigrationKey(id string) []byte {
	return bytes.Join([][]byte{
		MigrationKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}
