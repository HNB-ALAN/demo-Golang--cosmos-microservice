package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/product_certificate/v1/usc/product_certificate/v1"
)

// QueryServer defines the interface for the product_certificate module's query server
type QueryServer interface {
	QueryCertificate(context.Context, *blockchainproto.QueryCertificateRequest) (*blockchainproto.QueryCertificateResponse, error)
	QueryCertificates(context.Context, *blockchainproto.QueryCertificatesRequest) (*blockchainproto.QueryCertificatesResponse, error)
	QueryCertificateStats(context.Context, *blockchainproto.QueryCertificateStatsRequest) (*blockchainproto.QueryCertificateStatsResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// convertStatus converts string status to CertificateStatus enum
func convertStatus(status string) blockchainproto.CertificateStatus {
	switch status {
	case "active":
		return blockchainproto.CertificateStatus_CERTIFICATE_STATUS_ACTIVE
	case "revoked":
		return blockchainproto.CertificateStatus_CERTIFICATE_STATUS_REVOKED
	case "expired":
		return blockchainproto.CertificateStatus_CERTIFICATE_STATUS_EXPIRED
	default:
		return blockchainproto.CertificateStatus_CERTIFICATE_STATUS_UNSPECIFIED
	}
}

// convertTimestamp converts int64 timestamp to timestamppb.Timestamp
func convertTimestamp(timestamp int64) *timestamppb.Timestamp {
	return timestamppb.New(time.Unix(timestamp, 0))
}

// QueryCertificate handles certificate queries
func (k queryServer) QueryCertificate(ctx context.Context, req *blockchainproto.QueryCertificateRequest) (*blockchainproto.QueryCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	cert, found := k.Keeper.GetCertificate(sdkCtx, req.CertificateId)
	if !found {
		return nil, fmt.Errorf("certificate with ID %s not found", req.CertificateId)
	}

	// Convert internal certificate to blockchain-proto format
	blockchainCert := &blockchainproto.Certificate{
		Id:          cert.ID,
		ProductId:   cert.ProductID,
		Owner:       cert.Owner,
		Status:      convertStatus(cert.Status),
		ProductName: cert.Metadata,
		CreatedAt:   convertTimestamp(cert.CreatedAt),
		UpdatedAt:   convertTimestamp(cert.UpdatedAt),
	}

	return &blockchainproto.QueryCertificateResponse{Certificate: blockchainCert}, nil
}

// QueryCertificates handles all certificates queries
func (k queryServer) QueryCertificates(ctx context.Context, req *blockchainproto.QueryCertificatesRequest) (*blockchainproto.QueryCertificatesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	certificates, err := k.Keeper.GetAllCertificates(sdkCtx)
	if err != nil {
		return nil, err
	}

	// Convert internal certificates to blockchain-proto format
	var blockchainCerts []*blockchainproto.Certificate
	for _, cert := range certificates {
		blockchainCert := &blockchainproto.Certificate{
			Id:          cert.ID,
			ProductId:   cert.ProductID,
			Owner:       cert.Owner,
			Status:      convertStatus(cert.Status),
			ProductName: cert.Metadata,
			CreatedAt:   convertTimestamp(cert.CreatedAt),
			UpdatedAt:   convertTimestamp(cert.UpdatedAt),
		}
		blockchainCerts = append(blockchainCerts, blockchainCert)
	}

	return &blockchainproto.QueryCertificatesResponse{Certificates: blockchainCerts}, nil
}

// QueryCertificateStats handles certificate statistics queries
func (k queryServer) QueryCertificateStats(ctx context.Context, req *blockchainproto.QueryCertificateStatsRequest) (*blockchainproto.QueryCertificateStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all certificates for statistics
	certificates, err := k.Keeper.GetAllCertificates(sdkCtx)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalCertificates := int64(len(certificates))
	activeCertificates := int64(0)
	revokedCertificates := int64(0)

	for _, cert := range certificates {
		if cert.Status == "active" {
			activeCertificates++
		} else if cert.Status == "revoked" {
			revokedCertificates++
		}
	}

	stats := &blockchainproto.CertificateStats{
		TotalCertificates:   totalCertificates,
		ActiveCertificates:  activeCertificates,
		RevokedCertificates: revokedCertificates,
	}

	return &blockchainproto.QueryCertificateStatsResponse{Stats: stats}, nil
}
