package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the monitoring module's genesis state
func InitGenesis(ctx sdk.Context, k interface{}, genState GenesisState) {
	// TODO: Implement genesis initialization
	// This would typically involve:
	// - Setting parameters
	// - Initializing metrics
	// - Initializing alerts
	// - Initializing performance data
	// - Initializing system health
	// - Initializing monitoring configs
}

// ExportGenesis returns the monitoring module's exported genesis state
func ExportGenesis(ctx sdk.Context, k interface{}) GenesisState {
	// TODO: Implement genesis export
	// This would typically involve:
	// - Getting all metrics
	// - Getting all alerts
	// - Getting all performance data
	// - Getting all system health
	// - Getting all monitoring configs
	// - Getting parameters

	return GenesisState{
		Metrics:          []Metric{},
		Alerts:           []Alert{},
		PerformanceData:  []PerformanceData{},
		SystemHealth:     []SystemHealth{},
		MonitoringConfig: []MonitoringConfig{},
		Params:           DefaultParams(),
	}
}

// ValidateGenesis validates the monitoring module's genesis state
func ValidateGenesis(genState GenesisState) error {
	// Validate parameters
	if err := genState.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate metrics
	for _, metric := range genState.Metrics {
		if err := metric.Validate(); err != nil {
			return fmt.Errorf("invalid metric: %w", err)
		}
	}

	// Validate alerts
	for _, alert := range genState.Alerts {
		if err := alert.Validate(); err != nil {
			return fmt.Errorf("invalid alert: %w", err)
		}
	}

	// Validate performance data
	for _, perfData := range genState.PerformanceData {
		if err := perfData.Validate(); err != nil {
			return fmt.Errorf("invalid performance data: %w", err)
		}
	}

	// Validate system health
	for _, health := range genState.SystemHealth {
		if err := health.Validate(); err != nil {
			return fmt.Errorf("invalid system health: %w", err)
		}
	}

	// Validate monitoring configs
	for _, config := range genState.MonitoringConfig {
		if err := config.Validate(); err != nil {
			return fmt.Errorf("invalid monitoring config: %w", err)
		}
	}

	return nil
}
