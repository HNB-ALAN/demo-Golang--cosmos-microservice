package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	streamtypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/types"
)

// Keeper manages the streaming module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new Streaming keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetStream returns a Stream by its ID
func (k Keeper) GetStream(ctx sdk.Context, id string) (streamtypes.Stream, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.StreamKey(id))
	if bz == nil {
		return streamtypes.Stream{}, fmt.Errorf("stream with ID %s not found", id)
	}

	var stream streamtypes.Stream
	if err := json.Unmarshal(bz, &stream); err != nil {
		return streamtypes.Stream{}, fmt.Errorf("failed to unmarshal stream: %w", err)
	}

	return stream, nil
}

// SetStream sets a Stream
func (k Keeper) SetStream(ctx sdk.Context, stream streamtypes.Stream) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(stream)
	if err != nil {
		return fmt.Errorf("failed to marshal stream: %w", err)
	}
	store.Set(streamtypes.StreamKey(stream.ID), bz)
	return nil
}

// GetAllStreams returns all Streams
func (k Keeper) GetAllStreams(ctx sdk.Context) []streamtypes.Stream {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, streamtypes.StreamKeyPrefix)
	defer iterator.Close()

	var streams []streamtypes.Stream
	for ; iterator.Valid(); iterator.Next() {
		var stream streamtypes.Stream
		if err := json.Unmarshal(iterator.Value(), &stream); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		streams = append(streams, stream)
	}
	return streams
}

// GetViewer returns a StreamViewer by its ID
func (k Keeper) GetViewer(ctx sdk.Context, id string) (streamtypes.StreamViewer, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.ViewerKey(id))
	if bz == nil {
		return streamtypes.StreamViewer{}, fmt.Errorf("viewer with ID %s not found", id)
	}

	var viewer streamtypes.StreamViewer
	if err := json.Unmarshal(bz, &viewer); err != nil {
		return streamtypes.StreamViewer{}, fmt.Errorf("failed to unmarshal viewer: %w", err)
	}

	return viewer, nil
}

// SetViewer sets a StreamViewer
func (k Keeper) SetViewer(ctx sdk.Context, viewer streamtypes.StreamViewer) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(viewer)
	if err != nil {
		return fmt.Errorf("failed to marshal viewer: %w", err)
	}
	store.Set(streamtypes.ViewerKey(viewer.ID), bz)
	return nil
}

// GetAllViewers returns all StreamViewers
func (k Keeper) GetAllViewers(ctx sdk.Context) []streamtypes.StreamViewer {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, streamtypes.ViewerKeyPrefix)
	defer iterator.Close()

	var viewers []streamtypes.StreamViewer
	for ; iterator.Valid(); iterator.Next() {
		var viewer streamtypes.StreamViewer
		if err := json.Unmarshal(iterator.Value(), &viewer); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		viewers = append(viewers, viewer)
	}
	return viewers
}

// GetQuality returns a StreamQualityMetrics by its ID
func (k Keeper) GetQuality(ctx sdk.Context, id string) (streamtypes.StreamQualityMetrics, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.QualityKey(id))
	if bz == nil {
		return streamtypes.StreamQualityMetrics{}, fmt.Errorf("quality metrics with ID %s not found", id)
	}

	var quality streamtypes.StreamQualityMetrics
	if err := json.Unmarshal(bz, &quality); err != nil {
		return streamtypes.StreamQualityMetrics{}, fmt.Errorf("failed to unmarshal quality metrics: %w", err)
	}

	return quality, nil
}

// SetQuality sets a StreamQualityMetrics
func (k Keeper) SetQuality(ctx sdk.Context, quality streamtypes.StreamQualityMetrics) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(quality)
	if err != nil {
		return fmt.Errorf("failed to marshal quality metrics: %w", err)
	}
	store.Set(streamtypes.QualityKey(quality.ID), bz)
	return nil
}

// GetAllQualities returns all StreamQualityMetrics
func (k Keeper) GetAllQualities(ctx sdk.Context) []streamtypes.StreamQualityMetrics {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, streamtypes.QualityKeyPrefix)
	defer iterator.Close()

	var qualities []streamtypes.StreamQualityMetrics
	for ; iterator.Valid(); iterator.Next() {
		var quality streamtypes.StreamQualityMetrics
		if err := json.Unmarshal(iterator.Value(), &quality); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		qualities = append(qualities, quality)
	}
	return qualities
}

// GetAnalytics returns a StreamAnalytics by its ID
func (k Keeper) GetAnalytics(ctx sdk.Context, id string) (streamtypes.StreamAnalytics, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.AnalyticsKey(id))
	if bz == nil {
		return streamtypes.StreamAnalytics{}, fmt.Errorf("analytics with ID %s not found", id)
	}

	var analytics streamtypes.StreamAnalytics
	if err := json.Unmarshal(bz, &analytics); err != nil {
		return streamtypes.StreamAnalytics{}, fmt.Errorf("failed to unmarshal analytics: %w", err)
	}

	return analytics, nil
}

// SetAnalytics sets a StreamAnalytics
func (k Keeper) SetAnalytics(ctx sdk.Context, analytics streamtypes.StreamAnalytics) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(analytics)
	if err != nil {
		return fmt.Errorf("failed to marshal analytics: %w", err)
	}
	store.Set(streamtypes.AnalyticsKey(analytics.ID), bz)
	return nil
}

// GetAllAnalytics returns all StreamAnalytics
func (k Keeper) GetAllAnalytics(ctx sdk.Context) []streamtypes.StreamAnalytics {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, streamtypes.AnalyticsKeyPrefix)
	defer iterator.Close()

	var analytics []streamtypes.StreamAnalytics
	for ; iterator.Valid(); iterator.Next() {
		var analyticsItem streamtypes.StreamAnalytics
		if err := json.Unmarshal(iterator.Value(), &analyticsItem); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		analytics = append(analytics, analyticsItem)
	}
	return analytics
}

// GetEvent returns a StreamEvent by its ID
func (k Keeper) GetEvent(ctx sdk.Context, id string) (streamtypes.StreamEvent, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.EventKey(id))
	if bz == nil {
		return streamtypes.StreamEvent{}, fmt.Errorf("event with ID %s not found", id)
	}

	var event streamtypes.StreamEvent
	if err := json.Unmarshal(bz, &event); err != nil {
		return streamtypes.StreamEvent{}, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return event, nil
}

// SetEvent sets a StreamEvent
func (k Keeper) SetEvent(ctx sdk.Context, event streamtypes.StreamEvent) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	store.Set(streamtypes.EventKey(event.ID), bz)
	return nil
}

// GetAllEvents returns all StreamEvents
func (k Keeper) GetAllEvents(ctx sdk.Context) []streamtypes.StreamEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, streamtypes.EventKeyPrefix)
	defer iterator.Close()

	var events []streamtypes.StreamEvent
	for ; iterator.Valid(); iterator.Next() {
		var event streamtypes.StreamEvent
		if err := json.Unmarshal(iterator.Value(), &event); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		events = append(events, event)
	}
	return events
}

// GetParams returns the streaming module's parameters
func (k Keeper) GetParams(ctx sdk.Context) streamtypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(streamtypes.ParamsKey)
	if bz == nil {
		return streamtypes.DefaultParams()
	}

	var params streamtypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return streamtypes.DefaultParams()
	}

	return params
}

// SetParams sets the streaming module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params streamtypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(streamtypes.ParamsKey, bz)
}
