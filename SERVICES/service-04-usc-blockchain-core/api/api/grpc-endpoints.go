package api

import (
	"net"

	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
	"google.golang.org/grpc"
)

// GRPCEndpoints handles gRPC API endpoints for USC Blockchain Core Service
type GRPCEndpoints struct {
	server *grpc.Server
	logger *logging.Logger
}

// NewGRPCEndpoints creates new gRPC endpoints
func NewGRPCEndpoints(logger *logging.Logger) *GRPCEndpoints {
	// Create gRPC server with options
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
	)

	return &GRPCEndpoints{
		server: server,
		logger: logger,
	}
}

// RegisterServices registers all gRPC services
func (g *GRPCEndpoints) RegisterServices() {
	// Register Blockchain service
	// Note: This is a template - actual service registration would depend on your proto definitions
	// Example: pb.RegisterBlockchainServiceServer(g.server, BlockchainHandlers)

	g.logger.Info("Registered gRPC services",
		logging.String("service", constants.ServiceBlockchainCore))
}

// Start starts the gRPC server
func (g *GRPCEndpoints) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		g.logger.Error("Failed to listen on gRPC address",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("address", address),
			logging.Error(err))
		return err
	}

	g.logger.Info("Starting gRPC server",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("address", address))

	return g.server.Serve(lis)
}

// Stop stops the gRPC server gracefully
func (g *GRPCEndpoints) Stop() {
	g.logger.Info("Stopping gRPC server",
		logging.String("service", constants.ServiceBlockchainCore))

	g.server.GracefulStop()
}

// GetServer returns the gRPC server instance
func (g *GRPCEndpoints) GetServer() *grpc.Server {
	return g.server
}
