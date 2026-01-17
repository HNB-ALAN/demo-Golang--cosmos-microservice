package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/metrics"
)

// MetricsService handles critical performance metrics collection for USC Blockchain Core Service
type MetricsService struct {
	performanceMetrics *metrics.PerformanceMetrics
	logger             logging.Logger
}

// NewMetricsService creates a new critical performance metrics service
func NewMetricsService(logger logging.Logger) (*MetricsService, error) {
	performanceMetrics := metrics.NewPerformanceMetrics()

	return &MetricsService{
		performanceMetrics: performanceMetrics,
		logger:             logger,
	}, nil
}

// RecordRequest records a request metric with critical performance focus
func (m *MetricsService) RecordRequest(latency time.Duration, isError bool) {
	m.logger.Debug("Recording critical request",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Duration("latency", latency),
		logging.Bool("is_error", isError))

	m.performanceMetrics.RecordRequest(latency, isError)
}

// RecordCacheOperation records a cache operation for critical services
func (m *MetricsService) RecordCacheOperation(isHit bool) {
	m.logger.Debug("Recording critical cache operation",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Bool("is_hit", isHit))

	m.performanceMetrics.RecordCacheOperation(isHit)
}

// RecordDatabaseQuery records a database query for critical performance monitoring
func (m *MetricsService) RecordDatabaseQuery(latency time.Duration) {
	m.logger.Debug("Recording critical database query",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Duration("latency", latency))

	m.performanceMetrics.RecordDatabaseQuery(latency)
}

// RecordFederationPerformance records GraphQL federation performance metrics
func (m *MetricsService) RecordFederationPerformance(latency time.Duration, queryComplexity int) {
	m.logger.Debug("Recording federation performance",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Duration("latency", latency),
		logging.Int("query_complexity", queryComplexity))

	// Add federation-specific metrics recording
	m.performanceMetrics.RecordRequest(latency, false)
}

// RecordAuthenticationMetrics records authentication success/failure rates
func (m *MetricsService) RecordAuthenticationMetrics(isSuccess bool, latency time.Duration) {
	m.logger.Debug("Recording authentication metrics",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Bool("is_success", isSuccess),
		logging.Duration("latency", latency))

	m.performanceMetrics.RecordRequest(latency, !isSuccess)
}

// RecordBlockchainMetrics records blockchain transaction throughput and consensus metrics
func (m *MetricsService) RecordBlockchainMetrics(transactionCount int, blockTime time.Duration) {
	m.logger.Debug("Recording blockchain metrics",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Int("transaction_count", transactionCount),
		logging.Duration("block_time", blockTime))

	// Add blockchain-specific metrics recording
	m.performanceMetrics.RecordRequest(blockTime, false)
}

// RecordStreamingMetrics records video streaming quality and buffering metrics
func (m *MetricsService) RecordStreamingMetrics(bufferingRate float64, qualityScore int) {
	m.logger.Debug("Recording streaming metrics",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Float64("buffering_rate", bufferingRate),
		logging.Int("quality_score", qualityScore))

	// Add streaming-specific metrics recording
}

// RecordMLMetrics records ML inference time and model accuracy
func (m *MetricsService) RecordMLMetrics(inferenceTime time.Duration, accuracy float64) {
	m.logger.Debug("Recording ML metrics",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Duration("inference_time", inferenceTime),
		logging.Float64("accuracy", accuracy))

	m.performanceMetrics.RecordRequest(inferenceTime, false)
}

// RecordKafkaMetrics records Kafka message throughput and consumer lag
func (m *MetricsService) RecordKafkaMetrics(throughput int, consumerLag time.Duration) {
	m.logger.Debug("Recording Kafka metrics",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.Int("throughput", throughput),
		logging.Duration("consumer_lag", consumerLag))

	// Add Kafka-specific metrics recording
}

// GetMetrics returns current critical performance metrics
func (m *MetricsService) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	m.logger.Debug("Getting critical performance metrics",
		logging.String("service", constants.ServiceBlockchainCore))

	metrics := m.performanceMetrics.GetAllMetrics()

	m.logger.Debug("Critical performance metrics retrieved",
		logging.String("service", constants.ServiceBlockchainCore))

	return metrics, nil
}

// Close closes the critical performance metrics service
func (m *MetricsService) Close() error {
	m.logger.Info("Closing critical performance metrics service",
		logging.String("service", constants.ServiceBlockchainCore))

	// Performance metrics doesn't need explicit closing, but you can add cleanup logic here if needed
	return nil
}

// RecordDuration records the duration of an operation
func (m *MetricsService) RecordDuration(operation string, duration time.Duration) {
	m.logger.Debug("Operation duration recorded",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("operation", operation),
		logging.Duration("duration", duration))

	m.performanceMetrics.RecordRequest(duration, false)
}

// RecordSuccess records a successful operation
func (m *MetricsService) RecordSuccess(operation string, labels map[string]string) {
	m.logger.Debug("Operation success recorded",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("operation", operation),
		logging.String("labels", formatLabels(labels)))
}

// RecordFailure records a failed operation
func (m *MetricsService) RecordFailure(operation string, errorType string, labels map[string]string) {
	m.logger.Error("Operation failure recorded",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("operation", operation),
		logging.String("error_type", errorType),
		logging.String("labels", formatLabels(labels)))

	// Record failure metric with error flag
	m.performanceMetrics.RecordRequest(0, true)
}

// formatLabels formats labels map to string
func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "{}"
	}

	result := "{"
	first := true
	for key, value := range labels {
		if !first {
			result += ","
		}
		result += key + "=" + value
		first = false
	}
	result += "}"
	return result
}

// Blockchain-specific metrics
func (m *MetricsService) RecordBlockCreated(blockNumber int64, blockHash string) {
	m.RecordSuccess("block_created", map[string]string{
		"block_number": fmt.Sprintf("%d", blockNumber),
		"block_hash":   blockHash,
	})
}

func (m *MetricsService) RecordTransactionSubmitted(txHash string, from, to string) {
	m.RecordSuccess("transaction_submitted", map[string]string{
		"tx_hash": txHash,
		"from":    from,
		"to":      to,
	})
}

func (m *MetricsService) RecordUSCTransfer(amount string, from, to string) {
	m.RecordSuccess("usc_transfer", map[string]string{
		"amount": amount,
		"from":   from,
		"to":     to,
	})
}

func (m *MetricsService) RecordContractDeployed(contractAddress string) {
	m.RecordSuccess("contract_deployed", map[string]string{
		"contract_address": contractAddress,
	})
}

func (m *MetricsService) RecordNFTMinted(tokenId, contractAddress string) {
	m.RecordSuccess("nft_minted", map[string]string{
		"token_id":         tokenId,
		"contract_address": contractAddress,
	})
}

func (m *MetricsService) RecordTokenCreated(tokenAddress, tokenName string) {
	m.RecordSuccess("token_created", map[string]string{
		"token_address": tokenAddress,
		"token_name":    tokenName,
	})
}

func (m *MetricsService) RecordCertificateCreated(certificateId, productId string) {
	m.RecordSuccess("certificate_created", map[string]string{
		"certificate_id": certificateId,
		"product_id":     productId,
	})
}

func (m *MetricsService) RecordValidatorRegistered(validatorAddress string) {
	m.RecordSuccess("validator_registered", map[string]string{
		"validator_address": validatorAddress,
	})
}
