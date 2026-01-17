package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/performance/v1/usc/performance/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

// MsgServer defines the gRPC message server for the performance module using blockchain-proto types
type MsgServer struct {
	Keeper
}

// NewMsgServer creates a new performance message server
func NewMsgServer(keeper Keeper) *MsgServer {
	return &MsgServer{Keeper: keeper}
}

// RecordMetrics handles performance metrics recording messages
func (k MsgServer) RecordMetrics(ctx context.Context, msg *blockchainproto.MsgRecordMetrics) (*blockchainproto.MsgRecordMetricsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Recorder == "" {
		return nil, fmt.Errorf("recorder cannot be empty")
	}
	if msg.TargetId == "" {
		return nil, fmt.Errorf("target id cannot be empty")
	}
	if msg.TargetType == blockchainproto.TargetType_TARGET_TYPE_UNSPECIFIED {
		return nil, fmt.Errorf("target type cannot be unspecified")
	}
	if len(msg.MetricsData) == 0 {
		return nil, fmt.Errorf("metrics data cannot be empty")
	}

	// Create performance metric from blockchain-proto data
	metric := types.PerformanceMetric{
		ID:          fmt.Sprintf("metric_%s_%s", msg.TargetId, msg.Timestamp),
		Name:        "performance_metric",
		Value:       0, // Will be parsed from MetricsData JSON
		Unit:        "count",
		Timestamp:   time.Now(), // Using current time as placeholder
		Tags:        map[string]string{"target_id": msg.TargetId, "target_type": msg.TargetType.String()},
		Description: "Performance metric from blockchain-proto",
		Category:    "general",
	}

	// Set the performance metric
	if err := k.SetPerformanceMetric(sdkCtx, metric); err != nil {
		return nil, fmt.Errorf("failed to set performance metric: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePerformanceMetric,
			sdk.NewAttribute(types.AttributeKeyMetricID, metric.ID),
			sdk.NewAttribute(types.AttributeKeyMetricName, metric.Name),
			sdk.NewAttribute(types.AttributeKeyValue, fmt.Sprintf("%d", metric.Value)),
			sdk.NewAttribute(types.AttributeKeyUnit, metric.Unit),
		),
	)

	// Calculate real transaction hash
	recordingId := fmt.Sprintf("recording_%s_%s", msg.TargetId, msg.Timestamp)
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Recorder, msg.TargetId, "", "record_performance_metrics", recordingId, "")
	recordingHash := blocktypes.CalculateHashFromString(fmt.Sprintf("hash_%s", msg.TargetId))

	return &blockchainproto.MsgRecordMetricsResponse{
		Success:         true,
		RecordingId:     recordingId,
		RecordingHash:   recordingHash,
		TransactionHash: txHash,
	}, nil
}

// GetMetrics handles metrics retrieval messages
func (k MsgServer) GetMetrics(ctx context.Context, msg *blockchainproto.MsgGetMetrics) (*blockchainproto.MsgGetMetricsResponse, error) {
	// Basic validation
	if msg.Requester == "" {
		return nil, fmt.Errorf("requester cannot be empty")
	}
	if msg.TargetId == "" {
		return nil, fmt.Errorf("target id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get performance metrics from keeper
	allMetrics := k.GetAllPerformanceMetrics(sdkCtx)
	var blockchainMetrics []*blockchainproto.PerformanceMetric

	// Filter metrics by target_id (extracted from Tags)
	for _, metric := range allMetrics {
		// Extract target_id from Tags if available
		targetIDFromTags := metric.Tags["target_id"]

		// Filter by target_id if provided
		if msg.TargetId != "" && targetIDFromTags != msg.TargetId {
			continue
		}

		// Extract target_type from Tags
		targetTypeStr := metric.Tags["target_type"]
		targetType := blockchainproto.TargetType_TARGET_TYPE_UNSPECIFIED
		if targetTypeStr != "" {
			// Map string to enum (simplified mapping)
			switch targetTypeStr {
			case "validator", "TARGET_TYPE_VALIDATOR":
				targetType = blockchainproto.TargetType_TARGET_TYPE_VALIDATOR
			case "node", "TARGET_TYPE_NODE":
				targetType = blockchainproto.TargetType_TARGET_TYPE_NODE
			case "service", "TARGET_TYPE_SERVICE":
				targetType = blockchainproto.TargetType_TARGET_TYPE_SERVICE
			}
		}

		// Map category to metric type
		metricType := blockchainproto.MetricType_METRIC_TYPE_UNSPECIFIED
		switch metric.Category {
		case "cpu":
			metricType = blockchainproto.MetricType_METRIC_TYPE_CPU_USAGE
		case "memory":
			metricType = blockchainproto.MetricType_METRIC_TYPE_MEMORY_USAGE
		case "network":
			metricType = blockchainproto.MetricType_METRIC_TYPE_THROUGHPUT
		case "disk":
			metricType = blockchainproto.MetricType_METRIC_TYPE_DISK_USAGE
		}

		// Convert internal metric to blockchain-proto metric
		blockchainMetric := &blockchainproto.PerformanceMetric{
			Id:          metric.ID,
			TargetId:    targetIDFromTags,
			TargetType:  targetType,
			MetricType:  metricType,
			MetricName:  metric.Name,
			MetricValue: float64(metric.Value),
			Unit:        metric.Unit,
			Timestamp:   timestamppb.New(metric.Timestamp),
			Recorder:    metric.Tags["recorder"],
			Metadata:    metric.Tags,
			Status:      blockchainproto.MetricStatus_METRIC_STATUS_NORMAL,
		}

		blockchainMetrics = append(blockchainMetrics, blockchainMetric)
	}

	// If no metrics found, return empty list (not an error)
	if len(blockchainMetrics) == 0 {
		sdkCtx.Logger().Debug("No metrics found for target", "target_id", msg.TargetId)
	}

	// Calculate real transaction hash (query operation, but still need hash)
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Requester, msg.TargetId, "", "get_metrics", "", "")

	return &blockchainproto.MsgGetMetricsResponse{
		Success:         true,
		Metrics:         blockchainMetrics,
		TransactionHash: txHash,
	}, nil
}

// AnalyzeMetrics handles metrics analysis messages
func (k MsgServer) AnalyzeMetrics(ctx context.Context, msg *blockchainproto.MsgAnalyzeMetrics) (*blockchainproto.MsgAnalyzeMetricsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Analyzer == "" {
		return nil, fmt.Errorf("analyzer cannot be empty")
	}
	if msg.TargetId == "" {
		return nil, fmt.Errorf("target id cannot be empty")
	}
	if msg.AnalysisType == blockchainproto.AnalysisType_ANALYSIS_TYPE_UNSPECIFIED {
		return nil, fmt.Errorf("analysis type cannot be unspecified")
	}

	// Get metrics for target
	allMetrics := k.GetAllPerformanceMetrics(sdkCtx)
	var targetMetrics []types.PerformanceMetric

	// Filter metrics by target_id
	for _, metric := range allMetrics {
		if targetID := metric.Tags["target_id"]; targetID == msg.TargetId {
			targetMetrics = append(targetMetrics, metric)
		}
	}

	if len(targetMetrics) == 0 {
		return nil, fmt.Errorf("no metrics found for target: %s", msg.TargetId)
	}

	// Perform analysis based on analysis type
	analysisId := fmt.Sprintf("analysis_%s_%s_%d", msg.TargetId, msg.AnalysisType.String(), sdkCtx.BlockTime().Unix())
	var benchmarkResults []*blockchainproto.BenchmarkResult
	var optimization *types.Optimization

	switch msg.AnalysisType {
	case blockchainproto.AnalysisType_ANALYSIS_TYPE_TREND:
		// Trend analysis: analyze metric values over time
		trendResult := analyzeTrends(sdkCtx, targetMetrics)
		benchmarkResults = append(benchmarkResults, trendResult)

	case blockchainproto.AnalysisType_ANALYSIS_TYPE_ANOMALY:
		// Anomaly detection: identify unusual metric values
		anomalyResult := detectAnomalies(sdkCtx, targetMetrics)
		benchmarkResults = append(benchmarkResults, anomalyResult)

	case blockchainproto.AnalysisType_ANALYSIS_TYPE_COMPARISON:
		// Comparison analysis: compare metrics across different time periods
		comparisonResult := compareMetrics(sdkCtx, targetMetrics)
		benchmarkResults = append(benchmarkResults, comparisonResult)

	case blockchainproto.AnalysisType_ANALYSIS_TYPE_OPTIMIZATION:
		// Optimization analysis: identify optimization opportunities
		optResult, opt := analyzeOptimization(sdkCtx, targetMetrics, msg.TargetId)
		benchmarkResults = append(benchmarkResults, optResult)
		optimization = opt

	case blockchainproto.AnalysisType_ANALYSIS_TYPE_FORECAST:
		// Forecast analysis: predict future metric values
		forecastResult := forecastMetrics(sdkCtx, targetMetrics)
		benchmarkResults = append(benchmarkResults, forecastResult)

	case blockchainproto.AnalysisType_ANALYSIS_TYPE_CORRELATION:
		// Correlation analysis: find correlations between metrics
		correlationResult := analyzeCorrelations(sdkCtx, targetMetrics)
		benchmarkResults = append(benchmarkResults, correlationResult)

	default:
		return nil, fmt.Errorf("unsupported analysis type: %s", msg.AnalysisType.String())
	}

	// Store optimization if created
	if optimization != nil {
		if err := k.SetOptimization(sdkCtx, *optimization); err != nil {
			sdkCtx.Logger().Error("Failed to store optimization", "error", err, "optimization_id", optimization.ID)
		}
	}

	// Emit event with analysis details
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOptimization,
			sdk.NewAttribute("analyzer", msg.Analyzer),
			sdk.NewAttribute("target_id", msg.TargetId),
			sdk.NewAttribute("analysis_type", msg.AnalysisType.String()),
			sdk.NewAttribute("analysis_id", analysisId),
			sdk.NewAttribute("results_count", fmt.Sprintf("%d", len(benchmarkResults))),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Analyzer, msg.TargetId, "", "analyze_metrics", analysisId, "")

	return &blockchainproto.MsgAnalyzeMetricsResponse{
		Success:          true,
		AnalysisId:       analysisId,
		BenchmarkResults: benchmarkResults,
		TransactionHash:  txHash,
	}, nil
}

// OptimizePerformance handles performance optimization messages
func (k MsgServer) OptimizePerformance(ctx context.Context, msg *blockchainproto.MsgOptimizePerformance) (*blockchainproto.MsgOptimizePerformanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Optimizer == "" {
		return nil, fmt.Errorf("optimizer cannot be empty")
	}
	if msg.TargetId == "" {
		return nil, fmt.Errorf("target id cannot be empty")
	}
	if msg.OptimizationType == blockchainproto.OptimizationType_OPTIMIZATION_TYPE_UNSPECIFIED {
		return nil, fmt.Errorf("optimization type cannot be unspecified")
	}

	// Perform optimization (simplified implementation)
	// TODO: Implement actual optimization logic

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOptimization,
			sdk.NewAttribute("optimizer", msg.Optimizer),
			sdk.NewAttribute("target_id", msg.TargetId),
			sdk.NewAttribute("optimization_type", msg.OptimizationType.String()),
		),
	)

	// Calculate real transaction hash
	optimizationId := fmt.Sprintf("optimization_%s_%s", msg.TargetId, msg.OptimizationType.String())
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Optimizer, msg.TargetId, "", "optimize_performance", optimizationId, "")

	return &blockchainproto.MsgOptimizePerformanceResponse{
		Success:             true,
		OptimizationId:      optimizationId,
		OptimizationResults: []*blockchainproto.BenchmarkResult{}, // Empty for now
		TransactionHash:     txHash,
	}, nil
}

// Helper functions for different analysis types

// analyzeTrends performs trend analysis on metrics
func analyzeTrends(ctx sdk.Context, metrics []types.PerformanceMetric) *blockchainproto.BenchmarkResult {
	if len(metrics) < 2 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "trend_analysis",
			TestType:           "trend",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}
	}

	// Calculate average trend (simplified: compare first and last values)
	firstValue := float64(metrics[0].Value)
	lastValue := float64(metrics[len(metrics)-1].Value)
	trendValue := ((lastValue - firstValue) / firstValue) * 100.0

	return &blockchainproto.BenchmarkResult{
		TestName:           "trend_analysis",
		TestType:           "trend",
		BaselineValue:      firstValue,
		CurrentValue:       lastValue,
		ImprovementPercent: trendValue,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"metric_count": fmt.Sprintf("%d", len(metrics))},
	}
}

// detectAnomalies detects anomalies in metrics
func detectAnomalies(ctx sdk.Context, metrics []types.PerformanceMetric) *blockchainproto.BenchmarkResult {
	if len(metrics) == 0 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "anomaly_detection",
			TestType:           "anomaly",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}
	}

	// Calculate mean
	var sum int64
	for _, m := range metrics {
		sum += m.Value
	}
	mean := float64(sum) / float64(len(metrics))

	// Count anomalies (values > 1.5x mean)
	threshold := mean * 1.5
	anomalyCount := 0
	for _, m := range metrics {
		if float64(m.Value) > threshold {
			anomalyCount++
		}
	}

	anomalyPercent := (float64(anomalyCount) / float64(len(metrics))) * 100.0

	return &blockchainproto.BenchmarkResult{
		TestName:           "anomaly_detection",
		TestType:           "anomaly",
		BaselineValue:      mean,
		CurrentValue:       float64(anomalyCount),
		ImprovementPercent: anomalyPercent,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"threshold": fmt.Sprintf("%.2f", threshold), "total_metrics": fmt.Sprintf("%d", len(metrics))},
	}
}

// compareMetrics compares metrics across time periods
func compareMetrics(ctx sdk.Context, metrics []types.PerformanceMetric) *blockchainproto.BenchmarkResult {
	if len(metrics) < 2 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "comparison_analysis",
			TestType:           "comparison",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}
	}

	// Compare first half vs second half
	mid := len(metrics) / 2
	var firstHalfSum, secondHalfSum int64
	for i := 0; i < mid; i++ {
		firstHalfSum += metrics[i].Value
	}
	for i := mid; i < len(metrics); i++ {
		secondHalfSum += metrics[i].Value
	}

	firstHalfAvg := float64(firstHalfSum) / float64(mid)
	secondHalfAvg := float64(secondHalfSum) / float64(len(metrics)-mid)
	comparisonValue := ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100.0

	return &blockchainproto.BenchmarkResult{
		TestName:           "comparison_analysis",
		TestType:           "period_comparison",
		BaselineValue:      firstHalfAvg,
		CurrentValue:       secondHalfAvg,
		ImprovementPercent: comparisonValue,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"first_half_count": fmt.Sprintf("%d", mid), "second_half_count": fmt.Sprintf("%d", len(metrics)-mid)},
	}
}

// analyzeOptimization identifies optimization opportunities
func analyzeOptimization(ctx sdk.Context, metrics []types.PerformanceMetric, targetID string) (*blockchainproto.BenchmarkResult, *types.Optimization) {
	if len(metrics) == 0 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "optimization_analysis",
			TestType:           "optimization",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}, nil
	}

	// Calculate average performance
	var sum int64
	for _, m := range metrics {
		sum += m.Value
	}
	avgValue := float64(sum) / float64(len(metrics))

	// Identify optimization opportunity (simplified: if average is high, suggest optimization)
	optimizationScore := 0.0
	optimizationType := "configuration"
	if avgValue > 80 {
		optimizationScore = 75.0 // High optimization potential
		optimizationType = "resource"
	} else if avgValue > 60 {
		optimizationScore = 50.0 // Medium optimization potential
		optimizationType = "algorithm"
	}

	// Create optimization record
	var optimization *types.Optimization
	if optimizationScore > 0 {
		optimization = &types.Optimization{
			ID:          fmt.Sprintf("opt_%s_%d", targetID, ctx.BlockTime().Unix()),
			Name:        fmt.Sprintf("Optimization for %s", targetID),
			Description: fmt.Sprintf("Performance optimization opportunity identified with score %.2f", optimizationScore),
			Type:        optimizationType,
			Impact:      "medium",
			Status:      "pending",
			CreatedAt:   ctx.BlockTime(),
			Metrics:     map[string]int64{"avg_value": int64(avgValue)},
			Tags:        map[string]string{"target_id": targetID},
		}
	}

	return &blockchainproto.BenchmarkResult{
		TestName:           "optimization_analysis",
		TestType:           "optimization",
		BaselineValue:      avgValue,
		CurrentValue:       optimizationScore,
		ImprovementPercent: optimizationScore,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"optimization_type": optimizationType, "target_id": targetID},
	}, optimization
}

// forecastMetrics predicts future metric values
func forecastMetrics(ctx sdk.Context, metrics []types.PerformanceMetric) *blockchainproto.BenchmarkResult {
	if len(metrics) < 3 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "forecast_analysis",
			TestType:           "forecast",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}
	}

	// Simple linear forecast: use last 3 values to predict next
	last3 := metrics[len(metrics)-3:]
	var sum int64
	for _, m := range last3 {
		sum += m.Value
	}
	forecastValue := float64(sum) / 3.0
	lastValue := float64(metrics[len(metrics)-1].Value)
	improvement := ((forecastValue - lastValue) / lastValue) * 100.0

	return &blockchainproto.BenchmarkResult{
		TestName:           "forecast_analysis",
		TestType:           "forecast",
		BaselineValue:      lastValue,
		CurrentValue:       forecastValue,
		ImprovementPercent: improvement,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"forecast_method": "moving_average", "window_size": "3"},
	}
}

// analyzeCorrelations finds correlations between metrics
func analyzeCorrelations(ctx sdk.Context, metrics []types.PerformanceMetric) *blockchainproto.BenchmarkResult {
	if len(metrics) < 2 {
		return &blockchainproto.BenchmarkResult{
			TestName:           "correlation_analysis",
			TestType:           "correlation",
			BaselineValue:      0.0,
			CurrentValue:       0.0,
			ImprovementPercent: 0.0,
			TestDate:           timestamppb.New(ctx.BlockTime()),
			TestParameters:     make(map[string]string),
		}
	}

	// Simplified correlation: calculate coefficient of variation
	var sum int64
	for _, m := range metrics {
		sum += m.Value
	}
	mean := float64(sum) / float64(len(metrics))

	var varianceSum float64
	for _, m := range metrics {
		diff := float64(m.Value) - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(metrics))
	coefficientOfVariation := variance / (mean * mean) // Coefficient of variation

	return &blockchainproto.BenchmarkResult{
		TestName:           "correlation_analysis",
		TestType:           "correlation",
		BaselineValue:      mean,
		CurrentValue:       coefficientOfVariation,
		ImprovementPercent: coefficientOfVariation * 100.0,
		TestDate:           timestamppb.New(ctx.BlockTime()),
		TestParameters:     map[string]string{"metric_count": fmt.Sprintf("%d", len(metrics)), "variance": fmt.Sprintf("%.2f", variance)},
	}
}
