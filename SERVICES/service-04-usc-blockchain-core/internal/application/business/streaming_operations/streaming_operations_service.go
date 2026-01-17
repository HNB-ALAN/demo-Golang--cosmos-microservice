package streaming_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/streaming_operations"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service handles streaming operations business logic
type Service struct {
	repo              *streaming_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new streaming operations service
func NewService(
	repo *streaming_operations.Repository,
	cosmosApp *app.USCApp,
	blockchainStorage *storage.StateManager,
	logger *logging.Logger,
	validator *validation.Validator,
	metricsService *metrics.MetricsService,
) *Service {
	return &Service{
		repo:              repo,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		logger:            logger,
		validator:         validator,
		metrics:           metricsService,
	}
}

// StreamBlocks streams blockchain blocks continuously
func (s *Service) StreamBlocks(req *proto.StreamBlocksRequest, stream proto.StreamingOperationsService_StreamBlocksServer) error {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("stream_blocks", time.Since(start))
	}()

	ctx := stream.Context()
	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Starting block streaming in business service",
		logging.String("correlation_id", correlationID),
		logging.String("clientId", req.ClientId),
		logging.String("filterType", req.FilterType))

	// Business logic validation
	if req.ClientId == "" {
		s.logger.Error("Client ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("stream_blocks", "validation_error", map[string]string{
			"client_id": req.ClientId,
		})
		return status.Errorf(codes.InvalidArgument, "client_id is required")
	}

	// Create channel for block events
	blockChan := make(chan *proto.StreamBlocksResponse, 100)
	defer close(blockChan)

	// Start block producer in background
	go s.produceBlockEvents(ctx, req, blockChan)

	// Stream blocks to client
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Block stream ended by client",
				logging.String("correlation_id", correlationID),
				logging.String("clientId", req.ClientId))
			s.metrics.RecordSuccess("stream_blocks", map[string]string{
				"client_id": req.ClientId,
				"status":    "ended_by_client",
			})
			return ctx.Err()

		case block, ok := <-blockChan:
			if !ok {
				// Channel closed
				s.logger.Info("Block stream channel closed",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId))
				s.metrics.RecordSuccess("stream_blocks", map[string]string{
					"client_id": req.ClientId,
					"status":    "channel_closed",
				})
				return nil
			}

			// Send block to client
			if err := stream.Send(block); err != nil {
				s.logger.Error("Failed to send block",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId),
					logging.Error(err))
				s.metrics.RecordFailure("stream_blocks", "send_error", map[string]string{
					"client_id": req.ClientId,
				})
				return err
			}

			s.logger.Debug("Sent block",
				logging.String("correlation_id", correlationID),
				logging.String("blockHash", block.BlockHash),
				logging.Int64("blockNumber", block.BlockNumber))
		}
	}
}

// StreamTransactions streams blockchain transactions continuously
func (s *Service) StreamTransactions(req *proto.StreamTransactionsRequest, stream proto.StreamingOperationsService_StreamTransactionsServer) error {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("stream_transactions", time.Since(start))
	}()

	ctx := stream.Context()
	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Starting transaction streaming in business service",
		logging.String("correlation_id", correlationID),
		logging.String("clientId", req.ClientId),
		logging.String("transactionType", req.TransactionType))

	// Business logic validation
	if req.ClientId == "" {
		s.logger.Error("Client ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("stream_transactions", "validation_error", map[string]string{
			"client_id": req.ClientId,
		})
		return status.Errorf(codes.InvalidArgument, "client_id is required")
	}

	// Create channel for transaction events
	txChan := make(chan *proto.StreamTransactionsResponse, 100)
	defer close(txChan)

	// Start transaction producer in background
	go s.produceTransactionEvents(ctx, req, txChan)

	// Stream transactions to client
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Transaction stream ended by client",
				logging.String("correlation_id", correlationID),
				logging.String("clientId", req.ClientId))
			s.metrics.RecordSuccess("stream_transactions", map[string]string{
				"client_id": req.ClientId,
				"status":    "ended_by_client",
			})
			return ctx.Err()

		case tx, ok := <-txChan:
			if !ok {
				// Channel closed
				s.logger.Info("Transaction stream channel closed",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId))
				s.metrics.RecordSuccess("stream_transactions", map[string]string{
					"client_id": req.ClientId,
					"status":    "channel_closed",
				})
				return nil
			}

			// Send transaction to client
			if err := stream.Send(tx); err != nil {
				s.logger.Error("Failed to send transaction",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId),
					logging.Error(err))
				s.metrics.RecordFailure("stream_transactions", "send_error", map[string]string{
					"client_id": req.ClientId,
				})
				return err
			}

			s.logger.Debug("Sent transaction",
				logging.String("correlation_id", correlationID),
				logging.String("txHash", tx.TransactionHash),
				logging.String("eventType", tx.EventType))
		}
	}
}

// StreamValidatorEvents streams validator events continuously
func (s *Service) StreamValidatorEvents(req *proto.StreamValidatorEventsRequest, stream proto.StreamingOperationsService_StreamValidatorEventsServer) error {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("stream_validator_events", time.Since(start))
	}()

	ctx := stream.Context()
	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Starting validator events streaming in business service",
		logging.String("correlation_id", correlationID),
		logging.String("clientId", req.ClientId),
		logging.String("eventType", req.EventType))

	// Business logic validation
	if req.ClientId == "" {
		s.logger.Error("Client ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("stream_validator_events", "validation_error", map[string]string{
			"client_id": req.ClientId,
		})
		return status.Errorf(codes.InvalidArgument, "client_id is required")
	}

	// Create channel for validator events
	eventChan := make(chan *proto.StreamValidatorEventsResponse, 100)
	defer close(eventChan)

	// Start validator event producer in background
	go s.produceValidatorEvents(ctx, req, eventChan)

	// Stream validator events to client
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Validator events stream ended by client",
				logging.String("correlation_id", correlationID),
				logging.String("clientId", req.ClientId))
			s.metrics.RecordSuccess("stream_validator_events", map[string]string{
				"client_id": req.ClientId,
				"status":    "ended_by_client",
			})
			return ctx.Err()

		case event, ok := <-eventChan:
			if !ok {
				// Channel closed
				s.logger.Info("Validator events stream channel closed",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId))
				s.metrics.RecordSuccess("stream_validator_events", map[string]string{
					"client_id": req.ClientId,
					"status":    "channel_closed",
				})
				return nil
			}

			// Send validator event to client
			if err := stream.Send(event); err != nil {
				s.logger.Error("Failed to send validator event",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId),
					logging.Error(err))
				s.metrics.RecordFailure("stream_validator_events", "send_error", map[string]string{
					"client_id": req.ClientId,
				})
				return err
			}

			s.logger.Debug("Sent validator event",
				logging.String("correlation_id", correlationID),
				logging.String("validatorAddress", event.ValidatorAddress),
				logging.String("eventType", event.EventType))
		}
	}
}

// StreamNetworkEvents streams network events continuously
func (s *Service) StreamNetworkEvents(req *proto.StreamNetworkEventsRequest, stream proto.StreamingOperationsService_StreamNetworkEventsServer) error {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("stream_network_events", time.Since(start))
	}()

	ctx := stream.Context()
	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Starting network events streaming in business service",
		logging.String("correlation_id", correlationID),
		logging.String("clientId", req.ClientId),
		logging.String("eventType", req.EventType))

	// Business logic validation
	if req.ClientId == "" {
		s.logger.Error("Client ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("stream_network_events", "validation_error", map[string]string{
			"client_id": req.ClientId,
		})
		return status.Errorf(codes.InvalidArgument, "client_id is required")
	}

	// Create channel for network events
	eventChan := make(chan *proto.StreamNetworkEventsResponse, 100)
	defer close(eventChan)

	// Start network event producer in background
	go s.produceNetworkEvents(ctx, req, eventChan)

	// Stream network events to client
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Network events stream ended by client",
				logging.String("correlation_id", correlationID),
				logging.String("clientId", req.ClientId))
			s.metrics.RecordSuccess("stream_network_events", map[string]string{
				"client_id": req.ClientId,
				"status":    "ended_by_client",
			})
			return ctx.Err()

		case event, ok := <-eventChan:
			if !ok {
				// Channel closed
				s.logger.Info("Network events stream channel closed",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId))
				s.metrics.RecordSuccess("stream_network_events", map[string]string{
					"client_id": req.ClientId,
					"status":    "channel_closed",
				})
				return nil
			}

			// Send network event to client
			if err := stream.Send(event); err != nil {
				s.logger.Error("Failed to send network event",
					logging.String("correlation_id", correlationID),
					logging.String("clientId", req.ClientId),
					logging.Error(err))
				s.metrics.RecordFailure("stream_network_events", "send_error", map[string]string{
					"client_id": req.ClientId,
				})
				return err
			}

			s.logger.Debug("Sent network event",
				logging.String("correlation_id", correlationID),
				logging.String("eventId", event.EventId),
				logging.String("eventType", event.EventType))
		}
	}
}

// produceBlockEvents produces block events for streaming
func (s *Service) produceBlockEvents(ctx context.Context, req *proto.StreamBlocksRequest, blockChan chan<- *proto.StreamBlocksResponse) {
	// Determine ticker interval based on max_blocks_per_second
	tickerInterval := 2 * time.Second // Default: 0.5 blocks per second
	if req.MaxBlocksPerSecond > 0 {
		tickerInterval = time.Duration(1000/req.MaxBlocksPerSecond) * time.Millisecond
		if tickerInterval < 100*time.Millisecond {
			tickerInterval = 100 * time.Millisecond // Minimum 100ms
		}
	}

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	lastBlockNumber := int64(0)
	lastBlockHash := ""

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// Get block from repository (repository handles Keeper → Database fallback)
			block, err := s.repo.StreamBlocks(ctx, req)

			if err != nil {
				s.logger.Warn("Failed to get block",
					logging.Error(err),
					logging.String("clientId", req.ClientId))
				continue
			}

			// Only send if this is a new block (different hash or higher block number)
			if block.BlockHash != lastBlockHash && block.BlockNumber > lastBlockNumber {
				lastBlockNumber = block.BlockNumber
				lastBlockHash = block.BlockHash

				// Send block to channel
				select {
				case blockChan <- block:
					s.logger.Debug("Produced block event",
						logging.String("blockHash", block.BlockHash),
						logging.Int64("blockNumber", block.BlockNumber))
				default:
					s.logger.Warn("Block channel full, dropping block")
				}
			}
		}
	}
}

// produceTransactionEvents produces transaction events for streaming
func (s *Service) produceTransactionEvents(ctx context.Context, req *proto.StreamTransactionsRequest, txChan chan<- *proto.StreamTransactionsResponse) {
	// Determine ticker interval based on max_transactions_per_second
	tickerInterval := 1 * time.Second // Default: 1 transaction per second
	if req.MaxTransactionsPerSecond > 0 {
		tickerInterval = time.Duration(1000/req.MaxTransactionsPerSecond) * time.Millisecond
		if tickerInterval < 100*time.Millisecond {
			tickerInterval = 100 * time.Millisecond // Minimum 100ms
		}
	}

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// Get transaction from repository
			tx, err := s.repo.StreamTransactions(ctx, req)
			if err != nil {
				s.logger.Warn("Failed to get transaction",
					logging.Error(err),
					logging.String("clientId", req.ClientId))
				continue
			}

			// Send transaction to channel
			select {
			case txChan <- tx:
				s.logger.Debug("Produced transaction event",
					logging.String("txHash", tx.TransactionHash),
					logging.String("eventType", tx.EventType))
			default:
				s.logger.Warn("Transaction channel full, dropping transaction")
			}
		}
	}
}

// produceValidatorEvents produces validator events for streaming
func (s *Service) produceValidatorEvents(ctx context.Context, req *proto.StreamValidatorEventsRequest, eventChan chan<- *proto.StreamValidatorEventsResponse) {
	ticker := time.NewTicker(5 * time.Second) // Default: 1 event per 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// Get validator event from repository
			event, err := s.repo.StreamValidatorEvents(ctx, req)
			if err != nil {
				s.logger.Warn("Failed to get validator event",
					logging.Error(err),
					logging.String("clientId", req.ClientId))
				continue
			}

			// Send validator event to channel
			select {
			case eventChan <- event:
				s.logger.Debug("Produced validator event",
					logging.String("validatorAddress", event.ValidatorAddress),
					logging.String("eventType", event.EventType))
			default:
				s.logger.Warn("Validator event channel full, dropping event")
			}
		}
	}
}

// produceNetworkEvents produces network events for streaming
func (s *Service) produceNetworkEvents(ctx context.Context, req *proto.StreamNetworkEventsRequest, eventChan chan<- *proto.StreamNetworkEventsResponse) {
	ticker := time.NewTicker(10 * time.Second) // Default: 1 event per 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			// Get network event from repository
			event, err := s.repo.StreamNetworkEvents(ctx, req)
			if err != nil {
				s.logger.Warn("Failed to get network event",
					logging.Error(err),
					logging.String("clientId", req.ClientId))
				continue
			}

			// Send network event to channel
			select {
			case eventChan <- event:
				s.logger.Debug("Produced network event",
					logging.String("eventId", event.EventId),
					logging.String("eventType", event.EventType))
			default:
				s.logger.Warn("Network event channel full, dropping event")
			}
		}
	}
}

