package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/streaming/v1/usc/streaming/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/types"
)

// MsgServer defines the gRPC message server for the streaming module
type MsgServer interface {
	StartStream(context.Context, *blockchainproto.MsgStartStream) (*blockchainproto.MsgStartStreamResponse, error)
	StopStream(context.Context, *blockchainproto.MsgStopStream) (*blockchainproto.MsgStopStreamResponse, error)
	UpdateStream(context.Context, *blockchainproto.MsgUpdateStream) (*blockchainproto.MsgUpdateStreamResponse, error)
	StreamData(context.Context, *blockchainproto.MsgStreamData) (*blockchainproto.MsgStreamDataResponse, error)
	SubscribeStream(context.Context, *blockchainproto.MsgSubscribeStream) (*blockchainproto.MsgSubscribeStreamResponse, error)
}

// msgServer implements MsgServer
type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates a new Streaming message server
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// StartStream handles stream start messages
func (k msgServer) StartStream(ctx context.Context, msg *blockchainproto.MsgStartStream) (*blockchainproto.MsgStartStreamResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create stream from proto message
	streamID := fmt.Sprintf("%s-%s", msg.Streamer, msg.StreamName)
	stream := types.Stream{
		ID:         streamID,
		Title:      msg.StreamName,
		StreamerID: msg.Streamer,
		Status:     types.StreamStatusActive,
		StartedAt:  sdkCtx.BlockTime(),
		UpdatedAt:  sdkCtx.BlockTime(),
		Metadata:   map[string]string{"stream_type": msg.StreamType.String(), "target_audience": msg.TargetAudience},
	}

	// Validate the stream
	if err := stream.Validate(); err != nil {
		return nil, fmt.Errorf("invalid stream: %w", err)
	}

	// Set the stream
	if err := k.SetStream(sdkCtx, stream); err != nil {
		return nil, fmt.Errorf("failed to start stream: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamStarted,
			sdk.NewAttribute(types.AttributeKeyStreamID, streamID),
			sdk.NewAttribute(types.AttributeKeyStreamerID, msg.Streamer),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Streamer, "", "", "start_stream", streamID, "")
	startHash := blocktypes.CalculateHashFromString(fmt.Sprintf("start-%s", streamID))

	return &blockchainproto.MsgStartStreamResponse{
		Success:         true,
		StreamId:        streamID,
		StartHash:       startHash,
		TransactionHash: txHash,
	}, nil
}

// StopStream handles stream stop messages
func (k msgServer) StopStream(ctx context.Context, msg *blockchainproto.MsgStopStream) (*blockchainproto.MsgStopStreamResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing stream
	stream, err := k.GetStream(sdkCtx, msg.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Update stream status
	stream.Status = types.StreamStatusStopped
	endedAt := sdkCtx.BlockTime()
	stream.EndedAt = &endedAt
	stream.UpdatedAt = sdkCtx.BlockTime()

	// Set the updated stream
	if err := k.SetStream(sdkCtx, stream); err != nil {
		return nil, fmt.Errorf("failed to stop stream: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamStopped,
			sdk.NewAttribute(types.AttributeKeyStreamID, msg.StreamId),
			sdk.NewAttribute(types.AttributeKeyStreamerID, stream.StreamerID),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, stream.StreamerID, "", "", "stop_stream", msg.StreamId, "")
	stopHash := blocktypes.CalculateHashFromString(fmt.Sprintf("stop-%s", msg.StreamId))

	return &blockchainproto.MsgStopStreamResponse{
		Success:         true,
		StopHash:        stopHash,
		TransactionHash: txHash,
	}, nil
}

// UpdateStream handles stream update messages
func (k msgServer) UpdateStream(ctx context.Context, msg *blockchainproto.MsgUpdateStream) (*blockchainproto.MsgUpdateStreamResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing stream
	stream, err := k.GetStream(sdkCtx, msg.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Update stream fields from new_config
	if msg.NewConfig != nil && msg.NewConfig.StreamQuality != nil {
		stream.Quality = fmt.Sprintf("%dx%d", msg.NewConfig.StreamQuality.ResolutionWidth, msg.NewConfig.StreamQuality.ResolutionHeight)
	}
	if msg.NewAudience != "" {
		if stream.Metadata == nil {
			stream.Metadata = make(map[string]string)
		}
		stream.Metadata["target_audience"] = msg.NewAudience
	}
	stream.UpdatedAt = sdkCtx.BlockTime()

	// Validate updated stream
	if err := stream.Validate(); err != nil {
		return nil, fmt.Errorf("invalid stream update: %w", err)
	}

	// Set the updated stream
	if err := k.SetStream(sdkCtx, stream); err != nil {
		return nil, fmt.Errorf("failed to update stream: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamUpdated,
			sdk.NewAttribute(types.AttributeKeyStreamID, msg.StreamId),
			sdk.NewAttribute(types.AttributeKeyStreamerID, stream.StreamerID),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, stream.StreamerID, "", "", "update_stream", msg.StreamId, "")
	updateHash := blocktypes.CalculateHashFromString(fmt.Sprintf("update-%s", msg.StreamId))

	return &blockchainproto.MsgUpdateStreamResponse{
		Success:         true,
		UpdateHash:      updateHash,
		TransactionHash: txHash,
	}, nil
}

// StreamData handles stream data messages
func (k msgServer) StreamData(ctx context.Context, msg *blockchainproto.MsgStreamData) (*blockchainproto.MsgStreamDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing stream
	stream, err := k.GetStream(sdkCtx, msg.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Generate data ID
	dataID := fmt.Sprintf("data-%s-%d", msg.StreamId, sdkCtx.BlockHeight())

	// Store stream data in keeper (if needed, you can create a SetStreamData method)
	// For now, we just validate the stream exists and emit an event

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamQuality,
			sdk.NewAttribute(types.AttributeKeyStreamID, msg.StreamId),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Update stream last activity
	stream.UpdatedAt = sdkCtx.BlockTime()
	k.SetStream(sdkCtx, stream)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, stream.StreamerID, "", "", "stream_data", dataID, "")
	streamingHash := blocktypes.CalculateHashFromString(fmt.Sprintf("data-%s", dataID))

	return &blockchainproto.MsgStreamDataResponse{
		Success:         true,
		DataId:          dataID,
		StreamingHash:   streamingHash,
		TransactionHash: txHash,
	}, nil
}

// SubscribeStream handles stream subscription messages
func (k msgServer) SubscribeStream(ctx context.Context, msg *blockchainproto.MsgSubscribeStream) (*blockchainproto.MsgSubscribeStreamResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing stream
	stream, err := k.GetStream(sdkCtx, msg.StreamId)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Create subscription ID
	subscriptionID := fmt.Sprintf("sub-%s-%s", msg.StreamId, msg.Subscriber)

	// Create viewer (subscription) from proto message
	viewer := types.StreamViewer{
		ID:       subscriptionID,
		StreamID: msg.StreamId,
		ViewerID: msg.Subscriber,
		UserID:   msg.Subscriber,
		JoinedAt: sdkCtx.BlockTime(),
		Metadata: map[string]string{
			"subscription_type": msg.SubscriptionType.String(),
		},
	}

	// Validate the viewer
	if err := viewer.Validate(); err != nil {
		return nil, fmt.Errorf("invalid subscription: %w", err)
	}

	// Set the viewer
	if err := k.SetViewer(sdkCtx, viewer); err != nil {
		return nil, fmt.Errorf("failed to subscribe to stream: %w", err)
	}

	// Update stream viewer count
	stream.ViewerCount++
	stream.UpdatedAt = sdkCtx.BlockTime()
	k.SetStream(sdkCtx, stream)

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStreamViewer,
			sdk.NewAttribute(types.AttributeKeyStreamID, msg.StreamId),
			sdk.NewAttribute(types.AttributeKeyViewerID, msg.Subscriber),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Subscriber, stream.StreamerID, "", "subscribe_stream", subscriptionID, "")
	subscriptionHash := blocktypes.CalculateHashFromString(fmt.Sprintf("sub-%s", subscriptionID))

	return &blockchainproto.MsgSubscribeStreamResponse{
		Success:          true,
		SubscriptionId:   subscriptionID,
		SubscriptionHash: subscriptionHash,
		TransactionHash:  txHash,
	}, nil
}
