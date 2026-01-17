package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/monitoring/v1/usc/monitoring/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
)

// QueryServer defines the gRPC querier service using blockchain-proto types
type QueryServer interface {
	QueryMonitoring(context.Context, *blockchainproto.QueryMonitoringRequest) (*blockchainproto.QueryMonitoringResponse, error)
	QueryMonitoringSessions(context.Context, *blockchainproto.QueryMonitoringSessionsRequest) (*blockchainproto.QueryMonitoringSessionsResponse, error)
	QueryMetrics(context.Context, *blockchainproto.QueryMetricsRequest) (*blockchainproto.QueryMetricsResponse, error)
	QueryMonitoringStats(context.Context, *blockchainproto.QueryMonitoringStatsRequest) (*blockchainproto.QueryMonitoringStatsResponse, error)
}

type queryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryMonitoring returns a single monitoring session
func (k queryServer) QueryMonitoring(ctx context.Context, req *blockchainproto.QueryMonitoringRequest) (*blockchainproto.QueryMonitoringResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config, err := k.Keeper.GetMonitoringConfig(sdkCtx, req.MonitoringId)
	if err != nil {
		return nil, fmt.Errorf("monitoring config not found: %w", err)
	}

	session := convertConfigToSession(config, sdkCtx.BlockTime())

	return &blockchainproto.QueryMonitoringResponse{
		Session: session,
	}, nil
}

// QueryMonitoringSessions returns filtered monitoring sessions
func (k queryServer) QueryMonitoringSessions(ctx context.Context, req *blockchainproto.QueryMonitoringSessionsRequest) (*blockchainproto.QueryMonitoringSessionsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allConfigs := k.Keeper.GetAllMonitoringConfigs(sdkCtx)
	sessions := make([]*blockchainproto.MonitoringSession, 0, len(allConfigs))

	for _, config := range allConfigs {
		if req.Monitor != "" && config.ServiceName != req.Monitor {
			continue
		}
		if req.Status != blockchainproto.MonitoringStatus_MONITORING_STATUS_UNSPECIFIED {
			status := blockchainproto.MonitoringStatus_MONITORING_STATUS_ACTIVE
			if !config.Enabled {
				status = blockchainproto.MonitoringStatus_MONITORING_STATUS_STOPPED
			}
			if status != req.Status {
				continue
			}
		}
		sessions = append(sessions, convertConfigToSession(config, sdkCtx.BlockTime()))
	}

	return &blockchainproto.QueryMonitoringSessionsResponse{
		Sessions:   sessions,
		Pagination: nil,
	}, nil
}

// QueryMetrics returns performance metrics for a monitoring session
func (k queryServer) QueryMetrics(ctx context.Context, req *blockchainproto.QueryMetricsRequest) (*blockchainproto.QueryMetricsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allPerfData := k.Keeper.GetAllPerformanceData(sdkCtx)
	metrics := make([]*blockchainproto.MetricData, 0)

	for _, perf := range allPerfData {
		if perf.ServiceName != req.MonitoringId {
			continue
		}
		if !perf.Timestamp.Before(req.StartTime.AsTime()) && !perf.Timestamp.After(req.EndTime.AsTime()) {
			continue
		}
		if req.MetricType != blockchainproto.MetricType_METRIC_TYPE_UNSPECIFIED {
			if convertMetricNameToType(perf.MetricName) != req.MetricType {
				continue
			}
		}
		metrics = append(metrics, convertPerfDataToMetricData(perf))
	}

	return &blockchainproto.QueryMetricsResponse{
		Metrics: metrics,
	}, nil
}

// QueryMonitoringStats returns monitoring statistics
func (k queryServer) QueryMonitoringStats(ctx context.Context, req *blockchainproto.QueryMonitoringStatsRequest) (*blockchainproto.QueryMonitoringStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allConfigs := k.Keeper.GetAllMonitoringConfigs(sdkCtx)
	allPerfData := k.Keeper.GetAllPerformanceData(sdkCtx)
	allAlerts := k.Keeper.GetAllAlerts(sdkCtx)

	totalSessions := int64(len(allConfigs))
	activeSessions := int64(0)
	stoppedSessions := int64(0)
	totalMetrics := int64(len(allPerfData))
	totalAlerts := int64(len(allAlerts))

	var mostMonitoredType string
	typeCounts := make(map[string]int64)

	for _, config := range allConfigs {
		if config.Enabled {
			activeSessions++
		} else {
			stoppedSessions++
		}
		typeCounts[config.ServiceName]++
		if typeCounts[config.ServiceName] > typeCounts[mostMonitoredType] {
			mostMonitoredType = config.ServiceName
		}
	}

	stats := &blockchainproto.MonitoringStats{
		TotalSessions:                 totalSessions,
		ActiveSessions:                activeSessions,
		StoppedSessions:               stoppedSessions,
		TotalMetrics:                  totalMetrics,
		TotalAlerts:                   totalAlerts,
		AverageSessionDurationSeconds: 0, // Calculate if needed
		MostMonitoredType:             mostMonitoredType,
		LastActivity:                  timestamppb.New(sdkCtx.BlockTime()),
	}

	return &blockchainproto.QueryMonitoringStatsResponse{
		Stats: stats,
	}, nil
}

// Helper functions

func convertConfigToSession(config types.MonitoringConfig, blockTime time.Time) *blockchainproto.MonitoringSession {
	status := blockchainproto.MonitoringStatus_MONITORING_STATUS_ACTIVE
	if !config.Enabled {
		status = blockchainproto.MonitoringStatus_MONITORING_STATUS_STOPPED
	}

	return &blockchainproto.MonitoringSession{
		Id:         config.ID,
		Monitor:    config.ServiceName,
		TargetId:   config.ServiceName,
		TargetType: blockchainproto.TargetType_TARGET_TYPE_SERVICE,
		Status:     status,
		Config: &blockchainproto.MonitoringConfig{
			CollectionIntervalSeconds: int32(config.CheckInterval.Seconds()),
			RetentionDays:             int32(config.RetentionPeriod.Hours() / 24),
			EnableAlerts:              true,
			EnableNotifications:       true,
			MonitoredMetrics:          []string{},
			CustomSettings:            map[string]string{},
		},
		AlertRules:       []*blockchainproto.AlertRule{},
		MetricThresholds: []*blockchainproto.MetricThreshold{},
		StartedAt:        timestamppb.New(config.CreatedAt),
		StoppedAt:        nil,
		MetricsCount:     0,
		AlertsCount:      0,
		LastActivity:     timestamppb.New(config.UpdatedAt),
		Memo:             "",
	}
}

func convertPerfDataToMetricData(perf types.PerformanceData) *blockchainproto.MetricData {
	return &blockchainproto.MetricData{
		Id:           perf.ID,
		MonitoringId: perf.ServiceName,
		MetricType:   convertMetricNameToType(perf.MetricName),
		MetricValue:  float64(perf.Value),
		Timestamp:    timestamppb.New(perf.Timestamp),
		Metadata:     perf.Metadata,
		Status:       blockchainproto.MetricStatus_METRIC_STATUS_NORMAL,
		Unit:         convertUnitStringToEnum(perf.Unit),
	}
}

func convertMetricNameToType(name string) blockchainproto.MetricType {
	switch name {
	case "cpu", "CPU_USAGE":
		return blockchainproto.MetricType_METRIC_TYPE_CPU_USAGE
	case "memory", "MEMORY_USAGE":
		return blockchainproto.MetricType_METRIC_TYPE_MEMORY_USAGE
	case "disk", "DISK_USAGE":
		return blockchainproto.MetricType_METRIC_TYPE_DISK_USAGE
	case "latency", "NETWORK_LATENCY":
		return blockchainproto.MetricType_METRIC_TYPE_NETWORK_LATENCY
	case "throughput":
		return blockchainproto.MetricType_METRIC_TYPE_THROUGHPUT
	case "error", "ERROR_RATE":
		return blockchainproto.MetricType_METRIC_TYPE_ERROR_RATE
	case "response", "RESPONSE_TIME":
		return blockchainproto.MetricType_METRIC_TYPE_RESPONSE_TIME
	case "availability":
		return blockchainproto.MetricType_METRIC_TYPE_AVAILABILITY
	default:
		return blockchainproto.MetricType_METRIC_TYPE_CUSTOM
	}
}

func convertUnitStringToEnum(unit string) blockchainproto.MetricUnit {
	switch unit {
	case "percentage", "%":
		return blockchainproto.MetricUnit_METRIC_UNIT_PERCENTAGE
	case "bytes":
		return blockchainproto.MetricUnit_METRIC_UNIT_BYTES
	case "ms", "milliseconds":
		return blockchainproto.MetricUnit_METRIC_UNIT_MILLISECONDS
	case "s", "seconds":
		return blockchainproto.MetricUnit_METRIC_UNIT_SECONDS
	case "count":
		return blockchainproto.MetricUnit_METRIC_UNIT_COUNT
	case "rate":
		return blockchainproto.MetricUnit_METRIC_UNIT_RATE
	default:
		return blockchainproto.MetricUnit_METRIC_UNIT_CUSTOM
	}
}
