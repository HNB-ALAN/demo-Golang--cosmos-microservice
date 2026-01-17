package types

import (
	"bytes"
)

// Store key prefixes
var (
	StreamKeyPrefix    = []byte{0x01}
	ViewerKeyPrefix    = []byte{0x02}
	QualityKeyPrefix   = []byte{0x03}
	AnalyticsKeyPrefix = []byte{0x04}
	EventKeyPrefix     = []byte{0x05}
	ParamsKey          = []byte{0x06}
)

// StreamKey returns the store key for a Stream object
func StreamKey(id string) []byte {
	return bytes.Join([][]byte{
		StreamKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// ViewerKey returns the store key for a StreamViewer object
func ViewerKey(id string) []byte {
	return bytes.Join([][]byte{
		ViewerKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// QualityKey returns the store key for a StreamQualityMetrics object
func QualityKey(id string) []byte {
	return bytes.Join([][]byte{
		QualityKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// AnalyticsKey returns the store key for a StreamAnalytics object
func AnalyticsKey(id string) []byte {
	return bytes.Join([][]byte{
		AnalyticsKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}

// EventKey returns the store key for a StreamEvent object
func EventKey(id string) []byte {
	return bytes.Join([][]byte{
		EventKeyPrefix,
		[]byte(id),
	}, []byte("/"))
}
