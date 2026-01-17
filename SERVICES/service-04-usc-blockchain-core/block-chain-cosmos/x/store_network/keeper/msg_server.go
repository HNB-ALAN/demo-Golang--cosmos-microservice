package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/store_network/v1/usc/store_network/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/types"
)

// MsgServer defines the gRPC message server using blockchain-proto types
type MsgServer interface {
	StoreData(context.Context, *blockchainproto.MsgStoreData) (*blockchainproto.MsgStoreDataResponse, error)
	RetrieveData(context.Context, *blockchainproto.MsgRetrieveData) (*blockchainproto.MsgRetrieveDataResponse, error)
	DeleteData(context.Context, *blockchainproto.MsgDeleteData) (*blockchainproto.MsgDeleteDataResponse, error)
	UpdateData(context.Context, *blockchainproto.MsgUpdateData) (*blockchainproto.MsgUpdateDataResponse, error)
}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) MsgServer { return &msgServer{Keeper: keeper} }

// StoreData handles storing data per proto
func (k msgServer) StoreData(ctx context.Context, msg *blockchainproto.MsgStoreData) (*blockchainproto.MsgStoreDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	valueBytes := []byte(msg.Value)
	stored := types.StoredData{
		Key:         msg.Key,
		Value:       valueBytes,
		Size:        int64(len(valueBytes)),
		ContentType: mapDataTypeToContentType(msg.DataType),
		Tags:        map[string]string{},
		Metadata:    map[string]string{"storer": msg.Storer},
		CreatedAt:   sdkCtx.BlockTime(),
		UpdatedAt:   sdkCtx.BlockTime(),
		ExpiresAt:   time.Time{},
		Version:     1,
	}
	if err := stored.Validate(); err != nil {
		return nil, fmt.Errorf("invalid stored data: %w", err)
	}
	if err := k.SetStoredData(sdkCtx, stored); err != nil {
		return nil, fmt.Errorf("failed to store data: %w", err)
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDataStored,
			sdk.NewAttribute(types.AttributeKeyDataKey, stored.Key),
			sdk.NewAttribute(types.AttributeKeyDataSize, fmt.Sprintf("%d", stored.Size)),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Storer, "", "", "store_data", stored.Key, "")
	storageHash := blocktypes.CalculateHashFromString(fmt.Sprintf("store-%s", stored.Key))

	return &blockchainproto.MsgStoreDataResponse{Success: true, StorageHash: storageHash, TransactionHash: txHash}, nil
}

// RetrieveData returns the value for a key
func (k msgServer) RetrieveData(ctx context.Context, msg *blockchainproto.MsgRetrieveData) (*blockchainproto.MsgRetrieveDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	data, err := k.GetStoredData(sdkCtx, msg.Key)
	if err != nil {
		return nil, fmt.Errorf("stored data not found: %w", err)
	}
	// Calculate real transaction hash (query operation)
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Retriever, "", "", "retrieve_data", msg.Key, "")

	return &blockchainproto.MsgRetrieveDataResponse{
		Success:         true,
		Value:           string(data.Value),
		DataType:        mapContentTypeToDataType(data.ContentType),
		StoredAt:        timestamppb.New(data.CreatedAt),
		TransactionHash: txHash,
	}, nil
}

// DeleteData marks a key as deleted
func (k msgServer) DeleteData(ctx context.Context, msg *blockchainproto.MsgDeleteData) (*blockchainproto.MsgDeleteDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	data, err := k.GetStoredData(sdkCtx, msg.Key)
	if err != nil {
		return nil, fmt.Errorf("stored data not found: %w", err)
	}
	data.Value = nil
	data.Size = 0
	data.UpdatedAt = sdkCtx.BlockTime()
	if err := k.SetStoredData(sdkCtx, data); err != nil {
		return nil, fmt.Errorf("failed to delete data: %w", err)
	}
	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Deleter, "", "", "delete_data", msg.Key, "")
	deletionHash := blocktypes.CalculateHashFromString(fmt.Sprintf("del-%s", msg.Key))

	return &blockchainproto.MsgDeleteDataResponse{Success: true, DeletionHash: deletionHash, TransactionHash: txHash}, nil
}

// UpdateData updates a stored key
func (k msgServer) UpdateData(ctx context.Context, msg *blockchainproto.MsgUpdateData) (*blockchainproto.MsgUpdateDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	data, err := k.GetStoredData(sdkCtx, msg.Key)
	if err != nil {
		return nil, fmt.Errorf("stored data not found: %w", err)
	}
	newBytes := []byte(msg.NewValue)
	data.Value = newBytes
	data.Size = int64(len(newBytes))
	data.ContentType = mapDataTypeToContentType(msg.DataType)
	data.UpdatedAt = sdkCtx.BlockTime()
	if err := k.SetStoredData(sdkCtx, data); err != nil {
		return nil, fmt.Errorf("failed to update data: %w", err)
	}
	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Updater, "", "", "update_data", msg.Key, "")
	updateHash := blocktypes.CalculateHashFromString(fmt.Sprintf("upd-%s", msg.Key))

	return &blockchainproto.MsgUpdateDataResponse{Success: true, UpdateHash: updateHash, TransactionHash: txHash}, nil
}

func mapDataTypeToContentType(dt blockchainproto.DataType) string {
	switch dt {
	case blockchainproto.DataType_DATA_TYPE_JSON:
		return "application/json"
	case blockchainproto.DataType_DATA_TYPE_BINARY:
		return "application/octet-stream"
	case blockchainproto.DataType_DATA_TYPE_TEXT:
		return "text/plain"
	case blockchainproto.DataType_DATA_TYPE_ENCRYPTED:
		return "application/encrypted"
	case blockchainproto.DataType_DATA_TYPE_COMPRESSED:
		return "application/compressed"
	default:
		return "application/octet-stream"
	}
}

func mapContentTypeToDataType(ct string) blockchainproto.DataType {
	switch ct {
	case "application/json":
		return blockchainproto.DataType_DATA_TYPE_JSON
	case "text/plain":
		return blockchainproto.DataType_DATA_TYPE_TEXT
	case "application/encrypted":
		return blockchainproto.DataType_DATA_TYPE_ENCRYPTED
	case "application/compressed":
		return blockchainproto.DataType_DATA_TYPE_COMPRESSED
	default:
		return blockchainproto.DataType_DATA_TYPE_BINARY
	}
}
