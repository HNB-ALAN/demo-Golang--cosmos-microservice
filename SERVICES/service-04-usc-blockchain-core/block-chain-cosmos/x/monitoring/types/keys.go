package types

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Store key prefixes
var (
	MetricKeyPrefix           = []byte{0x01}
	AlertKeyPrefix            = []byte{0x02}
	PerformanceDataKeyPrefix  = []byte{0x03}
	SystemHealthKeyPrefix     = []byte{0x04}
	MonitoringConfigKeyPrefix = []byte{0x05}
	ParamsKey                 = []byte{0x06}
)

// MetricKey returns the key for a metric
func MetricKey(id string) []byte {
	return append(MetricKeyPrefix, []byte(id)...)
}

// AlertKey returns the key for an alert
func AlertKey(id string) []byte {
	return append(AlertKeyPrefix, []byte(id)...)
}

// PerformanceDataKey returns the key for performance data
func PerformanceDataKey(id string) []byte {
	return append(PerformanceDataKeyPrefix, []byte(id)...)
}

// SystemHealthKey returns the key for system health
func SystemHealthKey(id string) []byte {
	return append(SystemHealthKeyPrefix, []byte(id)...)
}

// MonitoringConfigKey returns the key for monitoring config
func MonitoringConfigKey(id string) []byte {
	return append(MonitoringConfigKeyPrefix, []byte(id)...)
}

// MetricByServiceKey returns the key for metrics by service
func MetricByServiceKey(serviceName, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(serviceName)...), []byte(metricID)...)
}

// AlertBySeverityKey returns the key for alerts by severity
func AlertBySeverityKey(severity, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(severity)...), []byte(alertID)...)
}

// PerformanceDataByServiceKey returns the key for performance data by service
func PerformanceDataByServiceKey(serviceName, dataID string) []byte {
	return append(append(PerformanceDataKeyPrefix, []byte(serviceName)...), []byte(dataID)...)
}

// SystemHealthByTimestampKey returns the key for system health by timestamp
func SystemHealthByTimestampKey(timestamp time.Time) []byte {
	return append(SystemHealthKeyPrefix, []byte(timestamp.Format(time.RFC3339))...)
}

// MonitoringConfigByServiceKey returns the key for monitoring config by service
func MonitoringConfigByServiceKey(serviceName string) []byte {
	return append(MonitoringConfigKeyPrefix, []byte(serviceName)...)
}

// MetricByTimestampKey returns the key for metrics by timestamp
func MetricByTimestampKey(timestamp time.Time, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(timestamp.Format(time.RFC3339))...), []byte(metricID)...)
}

// AlertByStatusKey returns the key for alerts by status
func AlertByStatusKey(status, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(status)...), []byte(alertID)...)
}

// PerformanceDataByTimestampKey returns the key for performance data by timestamp
func PerformanceDataByTimestampKey(timestamp time.Time, dataID string) []byte {
	return append(append(PerformanceDataKeyPrefix, []byte(timestamp.Format(time.RFC3339))...), []byte(dataID)...)
}

// SystemHealthByStatusKey returns the key for system health by status
func SystemHealthByStatusKey(status, healthID string) []byte {
	return append(append(SystemHealthKeyPrefix, []byte(status)...), []byte(healthID)...)
}

// MetricByNameKey returns the key for metrics by name
func MetricByNameKey(metricName, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(metricName)...), []byte(metricID)...)
}

// AlertByMetricKey returns the key for alerts by metric
func AlertByMetricKey(metricID, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(metricID)...), []byte(alertID)...)
}

// PerformanceDataByMetricKey returns the key for performance data by metric
func PerformanceDataByMetricKey(metricName, dataID string) []byte {
	return append(append(PerformanceDataKeyPrefix, []byte(metricName)...), []byte(dataID)...)
}

// SystemHealthByComponentKey returns the key for system health by component
func SystemHealthByComponentKey(componentName, healthID string) []byte {
	return append(append(SystemHealthKeyPrefix, []byte(componentName)...), []byte(healthID)...)
}

// MonitoringConfigByEnabledKey returns the key for monitoring config by enabled status
func MonitoringConfigByEnabledKey(enabled bool, configID string) []byte {
	enabledStr := "false"
	if enabled {
		enabledStr = "true"
	}
	return append(append(MonitoringConfigKeyPrefix, []byte(enabledStr)...), []byte(configID)...)
}

// MetricByValueRangeKey returns the key for metrics by value range
func MetricByValueRangeKey(minValue, maxValue int64, metricID string) []byte {
	minStr := strconv.FormatInt(minValue, 10)
	maxStr := strconv.FormatInt(maxValue, 10)
	rangeStr := fmt.Sprintf("%s-%s", minStr, maxStr)
	return append(append(MetricKeyPrefix, []byte(rangeStr)...), []byte(metricID)...)
}

// AlertByCreatedTimeKey returns the key for alerts by creation time
func AlertByCreatedTimeKey(createdAt time.Time, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(createdAt.Format(time.RFC3339))...), []byte(alertID)...)
}

// PerformanceDataByValueRangeKey returns the key for performance data by value range
func PerformanceDataByValueRangeKey(minValue, maxValue int64, dataID string) []byte {
	minStr := strconv.FormatInt(minValue, 10)
	maxStr := strconv.FormatInt(maxValue, 10)
	rangeStr := fmt.Sprintf("%s-%s", minStr, maxStr)
	return append(append(PerformanceDataKeyPrefix, []byte(rangeStr)...), []byte(dataID)...)
}

// SystemHealthByScoreRangeKey returns the key for system health by score range
func SystemHealthByScoreRangeKey(minScore, maxScore int64, healthID string) []byte {
	minStr := strconv.FormatInt(minScore, 10)
	maxStr := strconv.FormatInt(maxScore, 10)
	rangeStr := fmt.Sprintf("%s-%s", minStr, maxStr)
	return append(append(SystemHealthKeyPrefix, []byte(rangeStr)...), []byte(healthID)...)
}

// MonitoringConfigByIntervalKey returns the key for monitoring config by check interval
func MonitoringConfigByIntervalKey(interval time.Duration, configID string) []byte {
	intervalStr := interval.String()
	return append(append(MonitoringConfigKeyPrefix, []byte(intervalStr)...), []byte(configID)...)
}

// GetMetricIDFromKey extracts metric ID from key
func GetMetricIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, MetricKeyPrefix) {
		return ""
	}
	return string(key[len(MetricKeyPrefix):])
}

// GetAlertIDFromKey extracts alert ID from key
func GetAlertIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, AlertKeyPrefix) {
		return ""
	}
	return string(key[len(AlertKeyPrefix):])
}

// GetPerformanceDataIDFromKey extracts performance data ID from key
func GetPerformanceDataIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, PerformanceDataKeyPrefix) {
		return ""
	}
	return string(key[len(PerformanceDataKeyPrefix):])
}

// GetSystemHealthIDFromKey extracts system health ID from key
func GetSystemHealthIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, SystemHealthKeyPrefix) {
		return ""
	}
	return string(key[len(SystemHealthKeyPrefix):])
}

// GetMonitoringConfigIDFromKey extracts monitoring config ID from key
func GetMonitoringConfigIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, MonitoringConfigKeyPrefix) {
		return ""
	}
	return string(key[len(MonitoringConfigKeyPrefix):])
}
