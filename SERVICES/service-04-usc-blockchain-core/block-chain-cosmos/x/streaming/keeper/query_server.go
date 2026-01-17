package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/streaming/v1/usc/streaming/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/types"
)

// QueryServer defines the gRPC querier service for the streaming module
type QueryServer interface {
	QueryStream(context.Context, *blockchainproto.QueryStreamRequest) (*blockchainproto.QueryStreamResponse, error)
	QueryStreams(context.Context, *blockchainproto.QueryStreamsRequest) (*blockchainproto.QueryStreamsResponse, error)
	QueryStreamData(context.Context, *blockchainproto.QueryStreamDataRequest) (*blockchainproto.QueryStreamDataResponse, error)
	QueryStreamSubscriptions(context.Context, *blockchainproto.QueryStreamSubscriptionsRequest) (*blockchainproto.QueryStreamSubscriptionsResponse, error)
	QueryStreamStats(context.Context, *blockchainproto.QueryStreamStatsRequest) (*blockchainproto.QueryStreamStatsResponse, error)
}

// queryServer implements QueryServer
type queryServer struct {
	Keeper
}

// NewQueryServerImpl creates a new Streaming query server
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryStream handles stream queries by ID
func (k queryServer) QueryStream(ctx context.Context, req *blockchainproto.QueryStreamRequest) (*blockchainproto.QueryStreamResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stream, err := k.Keeper.GetStream(sdkCtx, req.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Convert internal stream to blockchainproto.Stream
	blockchainStream := convertStreamToProto(stream)

	return &blockchainproto.QueryStreamResponse{
		Stream: blockchainStream,
	}, nil
}

// QueryStreams handles queries for multiple streams
func (k queryServer) QueryStreams(ctx context.Context, req *blockchainproto.QueryStreamsRequest) (*blockchainproto.QueryStreamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	streams := k.Keeper.GetAllStreams(sdkCtx)

	// Apply filters
	filteredStreams := []types.Stream{}
	for _, stream := range streams {
		// Filter by streamer
		if req.Streamer != "" && stream.StreamerID != req.Streamer {
			continue
		}
		// Filter by stream_type (can't map directly, skip for now)
		// Filter by status
		if req.Status != blockchainproto.StreamStatus_STREAM_STATUS_UNSPECIFIED {
			protoStatus := convertStatusToProto(stream.Status)
			if protoStatus != req.Status {
				continue
			}
		}
		filteredStreams = append(filteredStreams, stream)
	}

	// Convert to proto
	protoStreams := make([]*blockchainproto.Stream, len(filteredStreams))
	for i, stream := range filteredStreams {
		protoStreams[i] = convertStreamToProto(stream)
	}

	return &blockchainproto.QueryStreamsResponse{
		Streams: protoStreams,
	}, nil
}

// QueryStreamData handles queries for stream data
func (k queryServer) QueryStreamData(ctx context.Context, req *blockchainproto.QueryStreamDataRequest) (*blockchainproto.QueryStreamDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get stream to verify it exists
	_, err := k.Keeper.GetStream(sdkCtx, req.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// For now, return empty data array
	// In a real implementation, you would query stream data from keeper
	// based on stream_id, start_time, end_time, and data_type filters

	return &blockchainproto.QueryStreamDataResponse{
		Data: []*blockchainproto.StreamData{},
	}, nil
}

// QueryStreamSubscriptions handles queries for stream subscriptions
func (k queryServer) QueryStreamSubscriptions(ctx context.Context, req *blockchainproto.QueryStreamSubscriptionsRequest) (*blockchainproto.QueryStreamSubscriptionsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all viewers for the stream
	viewers := k.Keeper.GetAllViewers(sdkCtx)

	// Filter by stream_id and subscriber
	filteredViewers := []types.StreamViewer{}
	for _, viewer := range viewers {
		if viewer.StreamID != req.StreamId {
			continue
		}
		if req.Subscriber != "" && viewer.UserID != req.Subscriber {
			continue
		}
		// Filter by subscription_type (can't map directly, skip for now)
		filteredViewers = append(filteredViewers, viewer)
	}

	// Convert to proto StreamSubscription
	protoSubscriptions := make([]*blockchainproto.StreamSubscription, len(filteredViewers))
	for i, viewer := range filteredViewers {
		protoSubscriptions[i] = convertViewerToSubscription(viewer)
	}

	return &blockchainproto.QueryStreamSubscriptionsResponse{
		Subscriptions: protoSubscriptions,
	}, nil
}

// QueryStreamStats handles queries for stream statistics
func (k queryServer) QueryStreamStats(ctx context.Context, req *blockchainproto.QueryStreamStatsRequest) (*blockchainproto.QueryStreamStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allStreams := k.Keeper.GetAllStreams(sdkCtx)
	allViewers := k.Keeper.GetAllViewers(sdkCtx)

	// Calculate statistics
	totalStreams := int64(len(allStreams))
	activeStreams := int64(0)
	stoppedStreams := int64(0)
	totalSubscribers := int64(len(allViewers))
	totalDataPoints := int64(0)
	var mostPopularStreamType blockchainproto.StreamType = blockchainproto.StreamType_STREAM_TYPE_UNSPECIFIED

	for _, stream := range allStreams {
		if stream.Status == types.StreamStatusActive {
			activeStreams++
		}
		if stream.Status == types.StreamStatusStopped {
			stoppedStreams++
		}
	}

	stats := &blockchainproto.StreamStats{
		TotalStreams:                 totalStreams,
		ActiveStreams:                activeStreams,
		StoppedStreams:               stoppedStreams,
		TotalSubscribers:             totalSubscribers,
		TotalDataPoints:              totalDataPoints,
		AverageStreamDurationSeconds: 0, // Calculate if needed
		MostPopularStreamType:        mostPopularStreamType,
		TotalBandwidthUsedBytes:      0,
		AverageBandwidthUtilization:  0,
		PeakConcurrentStreams:        0,
		StreamSuccessRate:            0,
		LastActivity:                 timestamppb.New(sdkCtx.BlockTime()),
	}

	return &blockchainproto.QueryStreamStatsResponse{
		Stats: stats,
	}, nil
}

// Helper functions

func convertStreamToProto(stream types.Stream) *blockchainproto.Stream {
	protoStream := &blockchainproto.Stream{
		Id:               stream.ID,
		Name:             stream.Title,
		Streamer:         stream.StreamerID,
		StreamType:       blockchainproto.StreamType_STREAM_TYPE_VIDEO, // Default, can be improved
		Status:           convertStatusToProto(stream.Status),
		TargetAudience:   "",
		StartedAt:        timestamppb.New(stream.StartedAt),
		StoppedAt:        nil,
		TotalSubscribers: int64(stream.ViewerCount),
		TotalDataPoints:  0,
		LastActivity:     timestamppb.New(stream.UpdatedAt),
		Memo:             "",
	}

	if stream.EndedAt != nil {
		protoStream.StoppedAt = timestamppb.New(*stream.EndedAt)
	}

	// Parse quality from stream.Quality if needed
	// Set StreamConfig from metadata if needed

	return protoStream
}

func convertStatusToProto(status types.StreamStatus) blockchainproto.StreamStatus {
	switch status {
	case types.StreamStatusActive:
		return blockchainproto.StreamStatus_STREAM_STATUS_ACTIVE
	case types.StreamStatusStopped:
		return blockchainproto.StreamStatus_STREAM_STATUS_STOPPED
	case types.StreamStatusPaused:
		return blockchainproto.StreamStatus_STREAM_STATUS_PAUSED
	default:
		return blockchainproto.StreamStatus_STREAM_STATUS_UNSPECIFIED
	}
}

func convertViewerToSubscription(viewer types.StreamViewer) *blockchainproto.StreamSubscription {
	// If viewer has left, subscription is cancelled
	if viewer.LeftAt != nil {
		return &blockchainproto.StreamSubscription{
			Id:               viewer.ID,
			StreamId:         viewer.StreamID,
			Subscriber:       viewer.UserID,
			SubscriptionType: blockchainproto.SubscriptionType_SUBSCRIPTION_TYPE_FULL,
			Status:           blockchainproto.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELLED,
			SubscribedAt:     timestamppb.New(viewer.JoinedAt),
			LastActivity:     timestamppb.New(*viewer.LeftAt),
			DataReceived:     0,
			Memo:             "",
		}
	}

	// Active subscription
	return &blockchainproto.StreamSubscription{
		Id:               viewer.ID,
		StreamId:         viewer.StreamID,
		Subscriber:       viewer.UserID,
		SubscriptionType: blockchainproto.SubscriptionType_SUBSCRIPTION_TYPE_FULL,
		Status:           blockchainproto.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE,
		SubscribedAt:     timestamppb.New(viewer.JoinedAt),
		LastActivity:     nil,
		DataReceived:     0,
		Memo:             "",
	}
}
