package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/monitoring/v1/usc/monitoring/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
)

// MsgServer defines the gRPC message server using blockchain-proto types
type MsgServer interface {
	StartMonitoring(context.Context, *blockchainproto.MsgStartMonitoring) (*blockchainproto.MsgStartMonitoringResponse, error)
	StopMonitoring(context.Context, *blockchainproto.MsgStopMonitoring) (*blockchainproto.MsgStopMonitoringResponse, error)
	UpdateMonitoring(context.Context, *blockchainproto.MsgUpdateMonitoring) (*blockchainproto.MsgUpdateMonitoringResponse, error)
	RecordMetrics(context.Context, *blockchainproto.MsgRecordMetrics) (*blockchainproto.MsgRecordMetricsResponse, error)
}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// StartMonitoring handles starting monitoring session
func (k msgServer) StartMonitoring(ctx context.Context, msg *blockchainproto.MsgStartMonitoring) (*blockchainproto.MsgStartMonitoringResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create monitoring config
	monitoringID := fmt.Sprintf("monitoring-%s-%d", msg.TargetId, sdkCtx.BlockHeight())
	config := types.MonitoringConfig{
		ID:              monitoringID,
		ServiceName:     msg.TargetId,
		Enabled:         true,
		CheckInterval:   time.Duration(msg.MonitoringConfig.CollectionIntervalSeconds) * time.Second,
		RetentionPeriod: sdkCtx.BlockTime().AddDate(0, 0, int(msg.MonitoringConfig.RetentionDays)).Sub(sdkCtx.BlockTime()),
		CreatedAt:       sdkCtx.BlockTime(),
		UpdatedAt:       sdkCtx.BlockTime(),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid monitoring config: %w", err)
	}

	if err := k.SetMonitoringConfig(sdkCtx, config); err != nil {
		return nil, fmt.Errorf("failed to start monitoring: %w", err)
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMetricCreated,
			sdk.NewAttribute(types.AttributeKeyMetricID, monitoringID),
			sdk.NewAttribute(types.AttributeKeyMetricName, "monitoring_session"),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.TargetId, "", "", "start_monitoring", monitoringID, "")
	startHash := blocktypes.CalculateHashFromString(fmt.Sprintf("start-%s", monitoringID))

	return &blockchainproto.MsgStartMonitoringResponse{
		Success:         true,
		MonitoringId:    monitoringID,
		StartHash:       startHash,
		TransactionHash: txHash,
	}, nil
}

// StopMonitoring handles stopping monitoring session
func (k msgServer) StopMonitoring(ctx context.Context, msg *blockchainproto.MsgStopMonitoring) (*blockchainproto.MsgStopMonitoringResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config, err := k.GetMonitoringConfig(sdkCtx, msg.MonitoringId)
	if err != nil {
		return nil, fmt.Errorf("monitoring config not found: %w", err)
	}

	config.Enabled = false
	config.UpdatedAt = sdkCtx.BlockTime()

	if err := k.SetMonitoringConfig(sdkCtx, config); err != nil {
		return nil, fmt.Errorf("failed to stop monitoring: %w", err)
	}

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, config.ServiceName, "", "", "stop_monitoring", msg.MonitoringId, "")
	stopHash := blocktypes.CalculateHashFromString(fmt.Sprintf("stop-%s", msg.MonitoringId))

	return &blockchainproto.MsgStopMonitoringResponse{
		Success:         true,
		StopHash:        stopHash,
		TransactionHash: txHash,
	}, nil
}

// UpdateMonitoring handles updating monitoring configuration
func (k msgServer) UpdateMonitoring(ctx context.Context, msg *blockchainproto.MsgUpdateMonitoring) (*blockchainproto.MsgUpdateMonitoringResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config, err := k.GetMonitoringConfig(sdkCtx, msg.MonitoringId)
	if err != nil {
		return nil, fmt.Errorf("monitoring config not found: %w", err)
	}

	config.CheckInterval = time.Duration(msg.NewConfig.CollectionIntervalSeconds) * time.Second
	config.RetentionPeriod = sdkCtx.BlockTime().AddDate(0, 0, int(msg.NewConfig.RetentionDays)).Sub(sdkCtx.BlockTime())
	config.UpdatedAt = sdkCtx.BlockTime()

	if err := k.SetMonitoringConfig(sdkCtx, config); err != nil {
		return nil, fmt.Errorf("failed to update monitoring: %w", err)
	}

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, config.ServiceName, "", "", "update_monitoring", msg.MonitoringId, "")
	updateHash := blocktypes.CalculateHashFromString(fmt.Sprintf("update-%s", msg.MonitoringId))

	return &blockchainproto.MsgUpdateMonitoringResponse{
		Success:         true,
		UpdateHash:      updateHash,
		TransactionHash: txHash,
	}, nil
}

// RecordMetrics handles recording performance metrics
func (k msgServer) RecordMetrics(ctx context.Context, msg *blockchainproto.MsgRecordMetrics) (*blockchainproto.MsgRecordMetricsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for _, metricData := range msg.MetricsData {
		perfData := types.PerformanceData{
			ID:          metricData.Id,
			ServiceName: msg.MonitoringId,
			MetricName:  metricData.MetricType.String(),
			Value:       int64(metricData.MetricValue),
			Unit:        metricData.Unit.String(),
			Timestamp:   metricData.Timestamp.AsTime(),
			Metadata:    metricData.Metadata,
		}

		if err := k.SetPerformanceData(sdkCtx, perfData); err != nil {
			return nil, fmt.Errorf("failed to record metric %s: %w", metricData.Id, err)
		}
	}

	// Calculate real transaction hash
	recordingId := fmt.Sprintf("recording-%s-%d", msg.MonitoringId, sdkCtx.BlockHeight())
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.MonitoringId, "", "", "record_metrics", recordingId, "")
	recordingHash := blocktypes.CalculateHashFromString(fmt.Sprintf("record-%s", msg.MonitoringId))

	return &blockchainproto.MsgRecordMetricsResponse{
		Success:         true,
		RecordingHash:   recordingHash,
		TransactionHash: txHash,
	}, nil
}
