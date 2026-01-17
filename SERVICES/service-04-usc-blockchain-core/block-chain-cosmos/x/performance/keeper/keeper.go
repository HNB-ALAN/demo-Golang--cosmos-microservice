package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	perftypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

// Keeper manages the performance module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new performance keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetPerformanceMetric returns a performance metric by its ID
func (k Keeper) GetPerformanceMetric(ctx sdk.Context, id string) (perftypes.PerformanceMetric, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.MetricKey(id))
	if bz == nil {
		return perftypes.PerformanceMetric{}, fmt.Errorf("performance metric with ID %s not found", id)
	}

	var metric perftypes.PerformanceMetric
	if err := json.Unmarshal(bz, &metric); err != nil {
		return perftypes.PerformanceMetric{}, fmt.Errorf("failed to unmarshal performance metric: %w", err)
	}

	return metric, nil
}

// SetPerformanceMetric sets a performance metric
func (k Keeper) SetPerformanceMetric(ctx sdk.Context, metric perftypes.PerformanceMetric) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal performance metric: %w", err)
	}
	store.Set(perftypes.MetricKey(metric.ID), bz)
	return nil
}

// GetAllPerformanceMetrics returns all performance metrics
func (k Keeper) GetAllPerformanceMetrics(ctx sdk.Context) []perftypes.PerformanceMetric {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.MetricKeyPrefix)
	defer iterator.Close()

	var metrics []perftypes.PerformanceMetric
	for ; iterator.Valid(); iterator.Next() {
		var metric perftypes.PerformanceMetric
		if err := json.Unmarshal(iterator.Value(), &metric); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

// GetBenchmark returns a benchmark by its ID
func (k Keeper) GetBenchmark(ctx sdk.Context, id string) (perftypes.Benchmark, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.BenchmarkKey(id))
	if bz == nil {
		return perftypes.Benchmark{}, fmt.Errorf("benchmark with ID %s not found", id)
	}

	var benchmark perftypes.Benchmark
	if err := json.Unmarshal(bz, &benchmark); err != nil {
		return perftypes.Benchmark{}, fmt.Errorf("failed to unmarshal benchmark: %w", err)
	}

	return benchmark, nil
}

// SetBenchmark sets a benchmark
func (k Keeper) SetBenchmark(ctx sdk.Context, benchmark perftypes.Benchmark) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(benchmark)
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark: %w", err)
	}
	store.Set(perftypes.BenchmarkKey(benchmark.ID), bz)
	return nil
}

// GetAllBenchmarks returns all benchmarks
func (k Keeper) GetAllBenchmarks(ctx sdk.Context) []perftypes.Benchmark {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.BenchmarkKeyPrefix)
	defer iterator.Close()

	var benchmarks []perftypes.Benchmark
	for ; iterator.Valid(); iterator.Next() {
		var benchmark perftypes.Benchmark
		if err := json.Unmarshal(iterator.Value(), &benchmark); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		benchmarks = append(benchmarks, benchmark)
	}
	return benchmarks
}

// GetOptimization returns an optimization by its ID
func (k Keeper) GetOptimization(ctx sdk.Context, id string) (perftypes.Optimization, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.OptimizationKey(id))
	if bz == nil {
		return perftypes.Optimization{}, fmt.Errorf("optimization with ID %s not found", id)
	}

	var optimization perftypes.Optimization
	if err := json.Unmarshal(bz, &optimization); err != nil {
		return perftypes.Optimization{}, fmt.Errorf("failed to unmarshal optimization: %w", err)
	}

	return optimization, nil
}

// SetOptimization sets an optimization
func (k Keeper) SetOptimization(ctx sdk.Context, optimization perftypes.Optimization) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(optimization)
	if err != nil {
		return fmt.Errorf("failed to marshal optimization: %w", err)
	}
	store.Set(perftypes.OptimizationKey(optimization.ID), bz)
	return nil
}

// GetAllOptimizations returns all optimizations
func (k Keeper) GetAllOptimizations(ctx sdk.Context) []perftypes.Optimization {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.OptimizationKeyPrefix)
	defer iterator.Close()

	var optimizations []perftypes.Optimization
	for ; iterator.Valid(); iterator.Next() {
		var optimization perftypes.Optimization
		if err := json.Unmarshal(iterator.Value(), &optimization); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		optimizations = append(optimizations, optimization)
	}
	return optimizations
}

// GetPerformanceAlert returns a performance alert by its ID
func (k Keeper) GetPerformanceAlert(ctx sdk.Context, id string) (perftypes.PerformanceAlert, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.AlertKey(id))
	if bz == nil {
		return perftypes.PerformanceAlert{}, fmt.Errorf("performance alert with ID %s not found", id)
	}

	var alert perftypes.PerformanceAlert
	if err := json.Unmarshal(bz, &alert); err != nil {
		return perftypes.PerformanceAlert{}, fmt.Errorf("failed to unmarshal performance alert: %w", err)
	}

	return alert, nil
}

// SetPerformanceAlert sets a performance alert
func (k Keeper) SetPerformanceAlert(ctx sdk.Context, alert perftypes.PerformanceAlert) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal performance alert: %w", err)
	}
	store.Set(perftypes.AlertKey(alert.ID), bz)
	return nil
}

// GetAllPerformanceAlerts returns all performance alerts
func (k Keeper) GetAllPerformanceAlerts(ctx sdk.Context) []perftypes.PerformanceAlert {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.AlertKeyPrefix)
	defer iterator.Close()

	var alerts []perftypes.PerformanceAlert
	for ; iterator.Valid(); iterator.Next() {
		var alert perftypes.PerformanceAlert
		if err := json.Unmarshal(iterator.Value(), &alert); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

// GetPerformanceProfile returns a performance profile by its ID
func (k Keeper) GetPerformanceProfile(ctx sdk.Context, id string) (perftypes.PerformanceProfile, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.ProfileKey(id))
	if bz == nil {
		return perftypes.PerformanceProfile{}, fmt.Errorf("performance profile with ID %s not found", id)
	}

	var profile perftypes.PerformanceProfile
	if err := json.Unmarshal(bz, &profile); err != nil {
		return perftypes.PerformanceProfile{}, fmt.Errorf("failed to unmarshal performance profile: %w", err)
	}

	return profile, nil
}

// SetPerformanceProfile sets a performance profile
func (k Keeper) SetPerformanceProfile(ctx sdk.Context, profile perftypes.PerformanceProfile) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to marshal performance profile: %w", err)
	}
	store.Set(perftypes.ProfileKey(profile.ID), bz)
	return nil
}

// GetAllPerformanceProfiles returns all performance profiles
func (k Keeper) GetAllPerformanceProfiles(ctx sdk.Context) []perftypes.PerformanceProfile {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.ProfileKeyPrefix)
	defer iterator.Close()

	var profiles []perftypes.PerformanceProfile
	for ; iterator.Valid(); iterator.Next() {
		var profile perftypes.PerformanceProfile
		if err := json.Unmarshal(iterator.Value(), &profile); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		profiles = append(profiles, profile)
	}
	return profiles
}

// GetPerformanceReport returns a performance report by its ID
func (k Keeper) GetPerformanceReport(ctx sdk.Context, id string) (perftypes.PerformanceReport, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.ReportKey(id))
	if bz == nil {
		return perftypes.PerformanceReport{}, fmt.Errorf("performance report with ID %s not found", id)
	}

	var report perftypes.PerformanceReport
	if err := json.Unmarshal(bz, &report); err != nil {
		return perftypes.PerformanceReport{}, fmt.Errorf("failed to unmarshal performance report: %w", err)
	}

	return report, nil
}

// SetPerformanceReport sets a performance report
func (k Keeper) SetPerformanceReport(ctx sdk.Context, report perftypes.PerformanceReport) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal performance report: %w", err)
	}
	store.Set(perftypes.ReportKey(report.ID), bz)
	return nil
}

// GetAllPerformanceReports returns all performance reports
func (k Keeper) GetAllPerformanceReports(ctx sdk.Context) []perftypes.PerformanceReport {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, perftypes.ReportKeyPrefix)
	defer iterator.Close()

	var reports []perftypes.PerformanceReport
	for ; iterator.Valid(); iterator.Next() {
		var report perftypes.PerformanceReport
		if err := json.Unmarshal(iterator.Value(), &report); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		reports = append(reports, report)
	}
	return reports
}

// GetParams returns the performance module's parameters
func (k Keeper) GetParams(ctx sdk.Context) perftypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(perftypes.ParamsKey)
	if bz == nil {
		return perftypes.DefaultParams()
	}

	var params perftypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return perftypes.DefaultParams()
	}

	return params
}

// SetParams sets the performance module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params perftypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(perftypes.ParamsKey, bz)
}
