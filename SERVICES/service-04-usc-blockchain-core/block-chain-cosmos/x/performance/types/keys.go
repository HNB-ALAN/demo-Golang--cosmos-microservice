package types

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Store key prefixes
var (
	MetricKeyPrefix       = []byte{0x01}
	BenchmarkKeyPrefix    = []byte{0x02}
	OptimizationKeyPrefix = []byte{0x03}
	AlertKeyPrefix        = []byte{0x04}
	ProfileKeyPrefix      = []byte{0x05}
	ReportKeyPrefix       = []byte{0x06}
	ParamsKey             = []byte{0x07}
)

// MetricKey returns the key for a performance metric
func MetricKey(id string) []byte {
	return append(MetricKeyPrefix, []byte(id)...)
}

// BenchmarkKey returns the key for a benchmark
func BenchmarkKey(id string) []byte {
	return append(BenchmarkKeyPrefix, []byte(id)...)
}

// OptimizationKey returns the key for an optimization
func OptimizationKey(id string) []byte {
	return append(OptimizationKeyPrefix, []byte(id)...)
}

// AlertKey returns the key for a performance alert
func AlertKey(id string) []byte {
	return append(AlertKeyPrefix, []byte(id)...)
}

// ProfileKey returns the key for a performance profile
func ProfileKey(id string) []byte {
	return append(ProfileKeyPrefix, []byte(id)...)
}

// ReportKey returns the key for a performance report
func ReportKey(id string) []byte {
	return append(ReportKeyPrefix, []byte(id)...)
}

// MetricByServiceKey returns the key for metrics by service
func MetricByServiceKey(serviceName, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(serviceName)...), []byte(metricID)...)
}

// MetricByCategoryKey returns the key for metrics by category
func MetricByCategoryKey(category, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(category)...), []byte(metricID)...)
}

// BenchmarkByStatusKey returns the key for benchmarks by status
func BenchmarkByStatusKey(status, benchmarkID string) []byte {
	return append(append(BenchmarkKeyPrefix, []byte(status)...), []byte(benchmarkID)...)
}

// OptimizationByTypeKey returns the key for optimizations by type
func OptimizationByTypeKey(optType, optimizationID string) []byte {
	return append(append(OptimizationKeyPrefix, []byte(optType)...), []byte(optimizationID)...)
}

// OptimizationByStatusKey returns the key for optimizations by status
func OptimizationByStatusKey(status, optimizationID string) []byte {
	return append(append(OptimizationKeyPrefix, []byte(status)...), []byte(optimizationID)...)
}

// AlertBySeverityKey returns the key for alerts by severity
func AlertBySeverityKey(severity, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(severity)...), []byte(alertID)...)
}

// AlertByStatusKey returns the key for alerts by status
func AlertByStatusKey(status, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(status)...), []byte(alertID)...)
}

// ProfileByServiceKey returns the key for profiles by service
func ProfileByServiceKey(serviceName, profileID string) []byte {
	return append(append(ProfileKeyPrefix, []byte(serviceName)...), []byte(profileID)...)
}

// ReportByTimeRangeKey returns the key for reports by time range
func ReportByTimeRangeKey(startTime, endTime time.Time, reportID string) []byte {
	timeRange := fmt.Sprintf("%s-%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	return append(append(ReportKeyPrefix, []byte(timeRange)...), []byte(reportID)...)
}

// MetricByTimestampKey returns the key for metrics by timestamp
func MetricByTimestampKey(timestamp time.Time, metricID string) []byte {
	return append(append(MetricKeyPrefix, []byte(timestamp.Format(time.RFC3339))...), []byte(metricID)...)
}

// BenchmarkByDurationKey returns the key for benchmarks by duration
func BenchmarkByDurationKey(duration time.Duration, benchmarkID string) []byte {
	durationStr := duration.String()
	return append(append(BenchmarkKeyPrefix, []byte(durationStr)...), []byte(benchmarkID)...)
}

// OptimizationByImpactKey returns the key for optimizations by impact
func OptimizationByImpactKey(impact, optimizationID string) []byte {
	return append(append(OptimizationKeyPrefix, []byte(impact)...), []byte(optimizationID)...)
}

// AlertByCreatedTimeKey returns the key for alerts by creation time
func AlertByCreatedTimeKey(createdAt time.Time, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(createdAt.Format(time.RFC3339))...), []byte(alertID)...)
}

// ProfileByCreatedTimeKey returns the key for profiles by creation time
func ProfileByCreatedTimeKey(createdAt time.Time, profileID string) []byte {
	return append(append(ProfileKeyPrefix, []byte(createdAt.Format(time.RFC3339))...), []byte(profileID)...)
}

// ReportByCreatedTimeKey returns the key for reports by creation time
func ReportByCreatedTimeKey(createdAt time.Time, reportID string) []byte {
	return append(append(ReportKeyPrefix, []byte(createdAt.Format(time.RFC3339))...), []byte(reportID)...)
}

// MetricByValueRangeKey returns the key for metrics by value range
func MetricByValueRangeKey(minValue, maxValue int64, metricID string) []byte {
	minStr := strconv.FormatInt(minValue, 10)
	maxStr := strconv.FormatInt(maxValue, 10)
	rangeStr := fmt.Sprintf("%s-%s", minStr, maxStr)
	return append(append(MetricKeyPrefix, []byte(rangeStr)...), []byte(metricID)...)
}

// BenchmarkByResultRangeKey returns the key for benchmarks by result range
func BenchmarkByResultRangeKey(minResult, maxResult int64, benchmarkID string) []byte {
	minStr := strconv.FormatInt(minResult, 10)
	maxStr := strconv.FormatInt(maxResult, 10)
	rangeStr := fmt.Sprintf("%s-%s", minStr, maxStr)
	return append(append(BenchmarkKeyPrefix, []byte(rangeStr)...), []byte(benchmarkID)...)
}

// OptimizationByMetricsKey returns the key for optimizations by metrics
func OptimizationByMetricsKey(metricName string, optimizationID string) []byte {
	return append(append(OptimizationKeyPrefix, []byte(metricName)...), []byte(optimizationID)...)
}

// AlertByMetricKey returns the key for alerts by metric
func AlertByMetricKey(metricID, alertID string) []byte {
	return append(append(AlertKeyPrefix, []byte(metricID)...), []byte(alertID)...)
}

// ProfileByMetricsKey returns the key for profiles by metrics
func ProfileByMetricsKey(metricName, profileID string) []byte {
	return append(append(ProfileKeyPrefix, []byte(metricName)...), []byte(profileID)...)
}

// ReportBySummaryKey returns the key for reports by summary
func ReportBySummaryKey(summaryKey string, reportID string) []byte {
	return append(append(ReportKeyPrefix, []byte(summaryKey)...), []byte(reportID)...)
}

// GetMetricIDFromKey extracts metric ID from key
func GetMetricIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, MetricKeyPrefix) {
		return ""
	}
	return string(key[len(MetricKeyPrefix):])
}

// GetBenchmarkIDFromKey extracts benchmark ID from key
func GetBenchmarkIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, BenchmarkKeyPrefix) {
		return ""
	}
	return string(key[len(BenchmarkKeyPrefix):])
}

// GetOptimizationIDFromKey extracts optimization ID from key
func GetOptimizationIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, OptimizationKeyPrefix) {
		return ""
	}
	return string(key[len(OptimizationKeyPrefix):])
}

// GetAlertIDFromKey extracts alert ID from key
func GetAlertIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, AlertKeyPrefix) {
		return ""
	}
	return string(key[len(AlertKeyPrefix):])
}

// GetProfileIDFromKey extracts profile ID from key
func GetProfileIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, ProfileKeyPrefix) {
		return ""
	}
	return string(key[len(ProfileKeyPrefix):])
}

// GetReportIDFromKey extracts report ID from key
func GetReportIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, ReportKeyPrefix) {
		return ""
	}
	return string(key[len(ReportKeyPrefix):])
}
