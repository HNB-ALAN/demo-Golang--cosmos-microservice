package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
)

// TestIsCosmosAppAvailable tests the IsCosmosAppAvailable helper function
func TestIsCosmosAppAvailable(t *testing.T) {
	tests := []struct {
		name      string
		cosmosApp *app.USCApp
		want      bool
	}{
		{
			name:      "nil cosmosApp",
			cosmosApp: nil,
			want:      false,
		},
		{
			name:      "nil BaseApp",
			cosmosApp: &app.USCApp{BaseApp: nil},
			want:      false,
		},
		// Note: Testing with actual BaseApp requires full Cosmos SDK setup
		// This test case is skipped - we only test nil cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCosmosAppAvailable(tt.cosmosApp)
			assert.Equal(t, tt.want, got, "IsCosmosAppAvailable() = %v, want %v", got, tt.want)
		})
	}
}

// TestNormalizePagination tests the NormalizePagination helper function
func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      int32
		offset     int32
		config     PaginationConfig
		wantLimit  int32
		wantOffset int32
	}{
		{
			name:       "valid limit and offset",
			limit:      50,
			offset:     10,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 1000, DefaultOffset: 0},
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "limit <= 0 uses default",
			limit:      0,
			offset:     10,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 1000, DefaultOffset: 0},
			wantLimit:  100,
			wantOffset: 10,
		},
		{
			name:       "limit > max uses max",
			limit:      2000,
			offset:     10,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 1000, DefaultOffset: 0},
			wantLimit:  1000,
			wantOffset: 10,
		},
		{
			name:       "offset < 0 uses default",
			limit:      50,
			offset:     -5,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 1000, DefaultOffset: 0},
			wantLimit:  50,
			wantOffset: 0,
		},
		{
			name:       "both limit and offset invalid",
			limit:      -10,
			offset:     -5,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 1000, DefaultOffset: 0},
			wantLimit:  100,
			wantOffset: 0,
		},
		{
			name:       "maxLimit = 0 means no max limit",
			limit:      5000,
			offset:     10,
			config:     PaginationConfig{DefaultLimit: 100, MaxLimit: 0, DefaultOffset: 0},
			wantLimit:  5000,
			wantOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := NormalizePagination(tt.limit, tt.offset, tt.config)
			assert.Equal(t, tt.wantLimit, gotLimit, "NormalizePagination() limit = %v, want %v", gotLimit, tt.wantLimit)
			assert.Equal(t, tt.wantOffset, gotOffset, "NormalizePagination() offset = %v, want %v", gotOffset, tt.wantOffset)
		})
	}
}

// TestNormalizePaginationWithDefaults tests the NormalizePaginationWithDefaults helper function
func TestNormalizePaginationWithDefaults(t *testing.T) {
	tests := []struct {
		name       string
		limit      int32
		offset     int32
		wantLimit  int32
		wantOffset int32
	}{
		{
			name:       "valid limit and offset",
			limit:      50,
			offset:     10,
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "limit <= 0 uses default 100",
			limit:      0,
			offset:     10,
			wantLimit:  100,
			wantOffset: 10,
		},
		{
			name:       "limit > 1000 uses max 1000",
			limit:      2000,
			offset:     10,
			wantLimit:  1000,
			wantOffset: 10,
		},
		{
			name:       "offset < 0 uses default 0",
			limit:      50,
			offset:     -5,
			wantLimit:  50,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := NormalizePaginationWithDefaults(tt.limit, tt.offset)
			assert.Equal(t, tt.wantLimit, gotLimit, "NormalizePaginationWithDefaults() limit = %v, want %v", gotLimit, tt.wantLimit)
			assert.Equal(t, tt.wantOffset, gotOffset, "NormalizePaginationWithDefaults() offset = %v, want %v", gotOffset, tt.wantOffset)
		})
	}
}

// TestRecordPerformanceMetric tests the RecordPerformanceMetric helper function
func TestRecordPerformanceMetric(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewLogger(constants.ServiceBlockchainCore, config.LogConfig{Level: "debug"})
	start := time.Now().Add(-100 * time.Millisecond) // 100ms ago

	tests := []struct {
		name       string
		cosmosApp  *app.USCApp
		logger     *logging.Logger
		config     PerformanceMetricConfig
		identifier string
		success    bool
		wantError  bool
		skipTest   bool // Skip if requires full Cosmos SDK setup
	}{
		{
			name:       "nil cosmosApp skips recording",
			cosmosApp:  nil,
			logger:     logger,
			config:     PerformanceMetricConfig{ServiceName: "test", Operation: "test_op", MetricName: "test_metric", IDPrefix: "test"},
			identifier: "test_id",
			success:    true,
			wantError:  false, // Should return nil (not an error)
		},
		{
			name:       "nil BaseApp skips recording",
			cosmosApp:  &app.USCApp{BaseApp: nil},
			logger:     logger,
			config:     PerformanceMetricConfig{ServiceName: "test", Operation: "test_op", MetricName: "test_metric", IDPrefix: "test"},
			identifier: "test_id",
			success:    true,
			wantError:  false, // Should return nil (not an error)
		},
		// Note: Testing with actual BaseApp requires full Cosmos SDK setup
		// This test case is skipped - we only test nil cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skip("Skipping test that requires full Cosmos SDK setup")
			}

			err := RecordPerformanceMetric(ctx, tt.cosmosApp, tt.logger, start, tt.config, tt.identifier, tt.success)
			if tt.wantError {
				require.Error(t, err, "RecordPerformanceMetric() should return error")
			} else {
				require.NoError(t, err, "RecordPerformanceMetric() should not return error")
			}
		})
	}
}

// TestRecordPerformanceMetric_IdentifierTruncation tests that long identifiers are truncated
// This test is skipped because it requires full Cosmos SDK setup
func TestRecordPerformanceMetric_IdentifierTruncation(t *testing.T) {
	t.Skip("Skipping test that requires full Cosmos SDK setup")
	// In a real test with mock PerformanceKeeper, we would verify that the metric ID
	// contains only the first 8 characters of the identifier when identifier length > 8
}

// TestRecordPerformanceMetric_Tags tests that tags are correctly set
// This test is skipped because it requires full Cosmos SDK setup
func TestRecordPerformanceMetric_Tags(t *testing.T) {
	t.Skip("Skipping test that requires full Cosmos SDK setup")

	// In a real test with mock PerformanceKeeper, we would verify the tags
}

// TestRecordPerformanceMetricWithCustomTags tests the RecordPerformanceMetricWithCustomTags wrapper
// This test is skipped because it requires full Cosmos SDK setup
func TestRecordPerformanceMetricWithCustomTags(t *testing.T) {
	t.Skip("Skipping test that requires full Cosmos SDK setup")
}
