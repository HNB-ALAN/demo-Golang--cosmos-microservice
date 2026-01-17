package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	montypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
)

// Keeper manages the monitoring module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new monitoring keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetMetric returns a metric by its ID
func (k Keeper) GetMetric(ctx sdk.Context, id string) (montypes.Metric, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.MetricKey(id))
	if bz == nil {
		return montypes.Metric{}, fmt.Errorf("metric with ID %s not found", id)
	}

	var metric montypes.Metric
	if err := json.Unmarshal(bz, &metric); err != nil {
		return montypes.Metric{}, fmt.Errorf("failed to unmarshal metric: %w", err)
	}

	return metric, nil
}

// SetMetric sets a metric
func (k Keeper) SetMetric(ctx sdk.Context, metric montypes.Metric) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}
	store.Set(montypes.MetricKey(metric.ID), bz)
	return nil
}

// GetAllMetrics returns all metrics
func (k Keeper) GetAllMetrics(ctx sdk.Context) []montypes.Metric {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, montypes.MetricKeyPrefix)
	defer iterator.Close()

	var metrics []montypes.Metric
	for ; iterator.Valid(); iterator.Next() {
		var metric montypes.Metric
		if err := json.Unmarshal(iterator.Value(), &metric); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

// GetAlert returns an alert by its ID
func (k Keeper) GetAlert(ctx sdk.Context, id string) (montypes.Alert, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.AlertKey(id))
	if bz == nil {
		return montypes.Alert{}, fmt.Errorf("alert with ID %s not found", id)
	}

	var alert montypes.Alert
	if err := json.Unmarshal(bz, &alert); err != nil {
		return montypes.Alert{}, fmt.Errorf("failed to unmarshal alert: %w", err)
	}

	return alert, nil
}

// SetAlert sets an alert
func (k Keeper) SetAlert(ctx sdk.Context, alert montypes.Alert) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}
	store.Set(montypes.AlertKey(alert.ID), bz)
	return nil
}

// GetAllAlerts returns all alerts
func (k Keeper) GetAllAlerts(ctx sdk.Context) []montypes.Alert {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, montypes.AlertKeyPrefix)
	defer iterator.Close()

	var alerts []montypes.Alert
	for ; iterator.Valid(); iterator.Next() {
		var alert montypes.Alert
		if err := json.Unmarshal(iterator.Value(), &alert); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetPerformanceData returns performance data by its ID
func (k Keeper) GetPerformanceData(ctx sdk.Context, id string) (montypes.PerformanceData, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.PerformanceDataKey(id))
	if bz == nil {
		return montypes.PerformanceData{}, fmt.Errorf("performance data with ID %s not found", id)
	}

	var perfData montypes.PerformanceData
	if err := json.Unmarshal(bz, &perfData); err != nil {
		return montypes.PerformanceData{}, fmt.Errorf("failed to unmarshal performance data: %w", err)
	}

	return perfData, nil
}

// SetPerformanceData sets performance data
func (k Keeper) SetPerformanceData(ctx sdk.Context, perfData montypes.PerformanceData) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(perfData)
	if err != nil {
		return fmt.Errorf("failed to marshal performance data: %w", err)
	}
	store.Set(montypes.PerformanceDataKey(perfData.ID), bz)
	return nil
}

// GetAllPerformanceData returns all performance data
func (k Keeper) GetAllPerformanceData(ctx sdk.Context) []montypes.PerformanceData {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, montypes.PerformanceDataKeyPrefix)
	defer iterator.Close()

	var perfDataList []montypes.PerformanceData
	for ; iterator.Valid(); iterator.Next() {
		var perfData montypes.PerformanceData
		if err := json.Unmarshal(iterator.Value(), &perfData); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		perfDataList = append(perfDataList, perfData)
	}
	return perfDataList
}

// GetSystemHealth returns system health by its ID
func (k Keeper) GetSystemHealth(ctx sdk.Context, id string) (montypes.SystemHealth, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.SystemHealthKey(id))
	if bz == nil {
		return montypes.SystemHealth{}, fmt.Errorf("system health with ID %s not found", id)
	}

	var health montypes.SystemHealth
	if err := json.Unmarshal(bz, &health); err != nil {
		return montypes.SystemHealth{}, fmt.Errorf("failed to unmarshal system health: %w", err)
	}

	return health, nil
}

// SetSystemHealth sets system health
func (k Keeper) SetSystemHealth(ctx sdk.Context, health montypes.SystemHealth) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(health)
	if err != nil {
		return fmt.Errorf("failed to marshal system health: %w", err)
	}
	store.Set(montypes.SystemHealthKey(health.ID), bz)
	return nil
}

// GetAllSystemHealth returns all system health records
func (k Keeper) GetAllSystemHealth(ctx sdk.Context) []montypes.SystemHealth {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, montypes.SystemHealthKeyPrefix)
	defer iterator.Close()

	var healthList []montypes.SystemHealth
	for ; iterator.Valid(); iterator.Next() {
		var health montypes.SystemHealth
		if err := json.Unmarshal(iterator.Value(), &health); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		healthList = append(healthList, health)
	}
	return healthList
}

// GetMonitoringConfig returns monitoring config by its ID
func (k Keeper) GetMonitoringConfig(ctx sdk.Context, id string) (montypes.MonitoringConfig, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.MonitoringConfigKey(id))
	if bz == nil {
		return montypes.MonitoringConfig{}, fmt.Errorf("monitoring config with ID %s not found", id)
	}

	var config montypes.MonitoringConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		return montypes.MonitoringConfig{}, fmt.Errorf("failed to unmarshal monitoring config: %w", err)
	}

	return config, nil
}

// SetMonitoringConfig sets monitoring config
func (k Keeper) SetMonitoringConfig(ctx sdk.Context, config montypes.MonitoringConfig) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal monitoring config: %w", err)
	}
	store.Set(montypes.MonitoringConfigKey(config.ID), bz)
	return nil
}

// GetAllMonitoringConfigs returns all monitoring configs
func (k Keeper) GetAllMonitoringConfigs(ctx sdk.Context) []montypes.MonitoringConfig {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, montypes.MonitoringConfigKeyPrefix)
	defer iterator.Close()

	var configs []montypes.MonitoringConfig
	for ; iterator.Valid(); iterator.Next() {
		var config montypes.MonitoringConfig
		if err := json.Unmarshal(iterator.Value(), &config); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		configs = append(configs, config)
	}
	return configs
}

// GetParams returns the monitoring module's parameters
func (k Keeper) GetParams(ctx sdk.Context) montypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(montypes.ParamsKey)
	if bz == nil {
		return montypes.DefaultParams()
	}

	var params montypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return montypes.DefaultParams()
	}

	return params
}

// SetParams sets the monitoring module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params montypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(montypes.ParamsKey, bz)
}
