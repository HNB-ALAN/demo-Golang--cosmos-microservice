package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/store_network/v1/usc/store_network/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/types"
)

// QueryServer defines the gRPC querier service using blockchain-proto types
type QueryServer interface {
	QueryData(context.Context, *blockchainproto.QueryDataRequest) (*blockchainproto.QueryDataResponse, error)
	QueryDataList(context.Context, *blockchainproto.QueryDataListRequest) (*blockchainproto.QueryDataListResponse, error)
	QueryStorageStats(context.Context, *blockchainproto.QueryStorageStatsRequest) (*blockchainproto.QueryStorageStatsResponse, error)
}

type queryServer struct{ Keeper }

func NewQueryServerImpl(keeper Keeper) QueryServer { return &queryServer{Keeper: keeper} }

// QueryData returns a single stored data item
func (k queryServer) QueryData(ctx context.Context, req *blockchainproto.QueryDataRequest) (*blockchainproto.QueryDataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	data, err := k.Keeper.GetStoredData(sdkCtx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("stored data not found: %w", err)
	}
	return &blockchainproto.QueryDataResponse{Data: convertStoredToProto(data)}, nil
}

// QueryDataList returns a filtered list
func (k queryServer) QueryDataList(ctx context.Context, req *blockchainproto.QueryDataListRequest) (*blockchainproto.QueryDataListResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	all := k.Keeper.GetAllStoredData(sdkCtx)
	out := make([]*blockchainproto.StoredData, 0, len(all))
	for _, d := range all {
		if req.DataType != blockchainproto.DataType_DATA_TYPE_UNSPECIFIED {
			if convertTypeToProtoFromContentType(d.ContentType) != req.DataType {
				continue
			}
		}
		out = append(out, convertStoredToProto(d))
	}
	return &blockchainproto.QueryDataListResponse{DataList: out, Pagination: nil}, nil
}

// QueryStorageStats aggregates stats
func (k queryServer) QueryStorageStats(ctx context.Context, req *blockchainproto.QueryStorageStatsRequest) (*blockchainproto.QueryStorageStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	all := k.Keeper.GetAllStoredData(sdkCtx)
	var total, active, deleted, size int64
	var mostCommon blockchainproto.DataType = blockchainproto.DataType_DATA_TYPE_UNSPECIFIED
	counts := map[blockchainproto.DataType]int64{}
	for _, d := range all {
		total++
		size += d.Size
		active++
		dt := convertTypeToProtoFromContentType(d.ContentType)
		counts[dt]++
		if counts[dt] > counts[mostCommon] {
			mostCommon = dt
		}
	}
	avg := int64(0)
	if total > 0 {
		avg = size / total
	}
	stats := &blockchainproto.StorageStats{
		TotalEntries:          total,
		ActiveEntries:         active,
		DeletedEntries:        deleted,
		TotalSizeBytes:        size,
		AverageSizeBytes:      avg,
		TotalStorers:          0,
		MostCommonType:        mostCommon,
		TotalReplicas:         0,
		ShardedEntries:        0,
		ReplicationEfficiency: 0,
		LastActivity:          timestamppb.New(sdkCtx.BlockTime()),
	}
	return &blockchainproto.QueryStorageStatsResponse{Stats: stats}, nil
}

func convertStoredToProto(d types.StoredData) *blockchainproto.StoredData {
	var deletedAt *timestamppb.Timestamp
	return &blockchainproto.StoredData{
		Key:                 d.Key,
		Value:               string(d.Value),
		Storer:              "",
		DataType:            convertTypeToProtoFromContentType(d.ContentType),
		Status:              blockchainproto.DataStatus_DATA_STATUS_ACTIVE,
		StoredAt:            timestamppb.New(d.CreatedAt),
		UpdatedAt:           timestamppb.New(d.UpdatedAt),
		DeletedAt:           deletedAt,
		StorageHash:         "",
		EncryptionKeyHash:   "",
		SizeBytes:           d.Size,
		ReplicationFactor:   1,
		DataShardingEnabled: false,
		Metadata:            d.Metadata,
		Memo:                "",
	}
}

func convertTypeToProtoFromContentType(ct string) blockchainproto.DataType {
	switch ct {
	case "application/json":
		return blockchainproto.DataType_DATA_TYPE_JSON
	case "text/plain":
		return blockchainproto.DataType_DATA_TYPE_TEXT
	case "application/encrypted":
		return blockchainproto.DataType_DATA_TYPE_ENCRYPTED
	case "application/compressed":
		return blockchainproto.DataType_DATA_TYPE_COMPRESSED
	case "application/octet-stream":
		return blockchainproto.DataType_DATA_TYPE_BINARY
	default:
		return blockchainproto.DataType_DATA_TYPE_UNSPECIFIED
	}
}
