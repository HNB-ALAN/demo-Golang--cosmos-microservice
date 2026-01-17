package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	performancetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

// IsCosmosAppAvailable checks if cosmosApp is available and initialized
// Returns true if cosmosApp and BaseApp are both not nil
// This helper reduces duplicate code in business service methods
func IsCosmosAppAvailable(cosmosApp *app.USCApp) bool {
	return cosmosApp != nil && cosmosApp.BaseApp != nil
}

// PerformanceMetricConfig holds configuration for recording performance metrics
type PerformanceMetricConfig struct {
	ServiceName string            // Service name (e.g., "usc_coin_operations")
	Operation   string            // Operation name (e.g., "get_balance", "transfer_usc")
	MetricName  string            // Metric name (e.g., "balance_query_time")
	IDPrefix    string            // ID prefix for metric ID (e.g., "balance_query", "transfer")
	Tags        map[string]string // Additional tags for the metric
	Description string            // Metric description
}

// RecordPerformanceMetric records a performance metric using the PerformanceKeeper
// This helper reduces duplicate code in business service methods
// Returns error if cosmosApp is not available or metric recording fails
func RecordPerformanceMetric(
	ctx context.Context,
	cosmosApp *app.USCApp,
	logger *logging.Logger,
	start time.Time,
	config PerformanceMetricConfig,
	identifier string, // Address, contract address, or other identifier
	success bool, // Operation success status
) error {
	// Check if cosmosApp is available
	if !IsCosmosAppAvailable(cosmosApp) {
		if logger != nil {
			logger.Debug("CosmosApp not available, skipping performance metric recording")
		}
		return nil // Not an error, just skip recording
	}

	// Get SDK context
	sdkCtx, err := GetSDKContext(ctx, cosmosApp, logger)
	if err != nil {
		if logger != nil {
			logger.Debug("Failed to get SDK context for performance metric", logging.Error(err))
		}
		return err
	}

	// Calculate duration
	duration := time.Since(start).Milliseconds()

	// Build status string
	statusStr := "success"
	if !success {
		statusStr = "failed"
	}

	// Build metric ID
	// Truncate identifier to 8 chars if longer
	idIdentifier := identifier
	if len(identifier) > 8 {
		idIdentifier = identifier[:8]
	}
	metricID := fmt.Sprintf("%s_%s_%d", config.IDPrefix, idIdentifier, time.Now().Unix())

	// Build tags
	tags := make(map[string]string)
	if config.Tags != nil {
		for k, v := range config.Tags {
			tags[k] = v
		}
	}
	tags["service"] = config.ServiceName
	tags["operation"] = config.Operation
	tags["status"] = statusStr

	// Add identifier to tags if provided
	if identifier != "" {
		// Use appropriate tag name based on operation type
		if config.Operation == "get_balance" || config.Operation == "transfer_usc" {
			tags["address"] = identifier
		} else if config.Operation == "execute_contract" || config.Operation == "deploy_contract" {
			tags["contract"] = identifier
		} else {
			tags["identifier"] = identifier
		}
	}

	// Build description
	description := config.Description
	if description == "" {
		description = fmt.Sprintf("%s for %s (%s)", config.Operation, identifier, statusStr)
	}

	// Create metric
	metric := performancetypes.PerformanceMetric{
		ID:          metricID,
		Name:        config.MetricName,
		Value:       duration,
		Unit:        "milliseconds",
		Timestamp:   time.Now(),
		Category:    "performance",
		Tags:        tags,
		Description: description,
	}

	// Record metric
	if err := cosmosApp.PerformanceKeeper.SetPerformanceMetric(sdkCtx, metric); err != nil {
		if logger != nil {
			logger.Debug("Failed to record performance metric", logging.Error(err))
		}
		return err
	}

	return nil
}

// RecordPerformanceMetricWithCustomTags records a performance metric with custom tags
// This is a convenience wrapper that allows passing custom tags directly
// NOTE: Currently unused, but kept for future use when custom tags are needed
// without creating PerformanceMetricConfig struct
func RecordPerformanceMetricWithCustomTags(
	ctx context.Context,
	cosmosApp *app.USCApp,
	logger *logging.Logger,
	start time.Time,
	serviceName string,
	operation string,
	metricName string,
	idPrefix string,
	identifier string,
	success bool,
	customTags map[string]string,
	description string,
) error {
	config := PerformanceMetricConfig{
		ServiceName: serviceName,
		Operation:   operation,
		MetricName:  metricName,
		IDPrefix:    idPrefix,
		Tags:        customTags,
		Description: description,
	}
	return RecordPerformanceMetric(ctx, cosmosApp, logger, start, config, identifier, success)
}

// PaginationConfig holds configuration for pagination normalization
type PaginationConfig struct {
	DefaultLimit  int32 // Default limit if limit <= 0
	MaxLimit      int32 // Maximum limit allowed
	DefaultOffset int32 // Default offset if offset < 0
}

// NormalizePagination normalizes pagination parameters (limit and offset)
// This helper reduces duplicate code in business service methods
// Returns normalized limit and offset
func NormalizePagination(limit, offset int32, config PaginationConfig) (int32, int32) {
	// Normalize limit
	if limit <= 0 {
		limit = config.DefaultLimit
	}
	if config.MaxLimit > 0 && limit > config.MaxLimit {
		limit = config.MaxLimit
	}

	// Normalize offset
	if offset < 0 {
		offset = config.DefaultOffset
	}

	return limit, offset
}

// NormalizePaginationWithDefaults normalizes pagination with default values
// Default: limit default=100, max=1000, offset default=0
// This is a convenience wrapper for the most common pagination pattern
func NormalizePaginationWithDefaults(limit, offset int32) (int32, int32) {
	return NormalizePagination(limit, offset, PaginationConfig{
		DefaultLimit:  100,
		MaxLimit:      1000,
		DefaultOffset: 0,
	})
}
