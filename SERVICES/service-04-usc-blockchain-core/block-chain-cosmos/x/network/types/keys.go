package types

import (
	"bytes"
)

// Store key prefixes
var (
	NetworkKeyPrefix    = []byte{0x01}
	NodeKeyPrefix       = []byte{0x02}
	ConnectionKeyPrefix = []byte{0x03}
	SyncKeyPrefix       = []byte{0x04}
	HealthKeyPrefix     = []byte{0x05}
	ParamsKey           = []byte{0x06}
)

// NetworkKey returns the store key for a Network object
func NetworkKey(id string) []byte {
	return bytes.Join([][]byte{
		NetworkKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// NodeKey returns the store key for a Node object
func NodeKey(id string) []byte {
	return bytes.Join([][]byte{
		NodeKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// ConnectionKey returns the store key for a Connection object
func ConnectionKey(id string) []byte {
	return bytes.Join([][]byte{
		ConnectionKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// SyncKey returns the store key for a NetworkSync object
func SyncKey(id string) []byte {
	return bytes.Join([][]byte{
		SyncKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// HealthKey returns the store key for a NetworkHealth object
func HealthKey(id string) []byte {
	return bytes.Join([][]byte{
		HealthKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}
