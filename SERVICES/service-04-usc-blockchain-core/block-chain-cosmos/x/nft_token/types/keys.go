package types

// Store key prefixes
var (
	NFTKeyPrefix        = []byte{0x01}
	CollectionKeyPrefix = []byte{0x02}
	ParamsKey           = []byte{0x03}
)

// NFTKey returns the key for an NFT
func NFTKey(id string) []byte {
	return append(NFTKeyPrefix, []byte(id)...)
}

// CollectionKey returns the key for a collection
func CollectionKey(id string) []byte {
	return append(CollectionKeyPrefix, []byte(id)...)
}
