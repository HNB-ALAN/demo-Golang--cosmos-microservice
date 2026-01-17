package nft_token_operations

import (
	"context"

	"service-04/internal/application/business/nft_token_operations"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
)

// Handlers handles NFT token operations gRPC requests
type Handlers struct {
	proto.UnimplementedNFTTokenOperationsServiceServer
	service *nft_token_operations.Service
	logger  *logging.Logger
}

// NewHandlers creates a new NFT token operations handlers
func NewHandlers(service *nft_token_operations.Service, logger *logging.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// MintNFT creates a new NFT token
func (h *Handlers) MintNFT(ctx context.Context, req *proto.MintNFTRequest) (*proto.MintNFTResponse, error) {
	return h.service.MintNFT(ctx, req)
}

// TransferNFT transfers an NFT token
func (h *Handlers) TransferNFT(ctx context.Context, req *proto.TransferNFTRequest) (*proto.TransferNFTResponse, error) {
	return h.service.TransferNFT(ctx, req)
}

// GetNFTInfo retrieves NFT information
func (h *Handlers) GetNFTInfo(ctx context.Context, req *proto.GetNFTInfoRequest) (*proto.GetNFTInfoResponse, error) {
	return h.service.GetNFTInfo(ctx, req)
}

// GetNFTsByOwner retrieves NFTs owned by an address
func (h *Handlers) GetNFTsByOwner(ctx context.Context, req *proto.GetNFTsByOwnerRequest) (*proto.GetNFTsByOwnerResponse, error) {
	return h.service.GetNFTsByOwner(ctx, req)
}

// DeployNFTContract deploys a new NFT contract
func (h *Handlers) DeployNFTContract(ctx context.Context, req *proto.DeployNFTContractRequest) (*proto.DeployNFTContractResponse, error) {
	return h.service.DeployNFTContract(ctx, req)
}

// CreateNFTCollection creates a new NFT collection
func (h *Handlers) CreateNFTCollection(ctx context.Context, req *proto.CreateNFTCollectionRequest) (*proto.CreateNFTCollectionResponse, error) {
	return h.service.CreateNFTCollection(ctx, req)
}

// BurnNFT burns an NFT token
func (h *Handlers) BurnNFT(ctx context.Context, req *proto.BurnNFTRequest) (*proto.BurnNFTResponse, error) {
	return h.service.BurnNFT(ctx, req)
}
