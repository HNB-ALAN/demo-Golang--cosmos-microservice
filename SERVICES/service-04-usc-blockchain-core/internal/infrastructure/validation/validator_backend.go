package validation

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/validation"
)

// Validator handles input validation for USC Blockchain Core Service (Backend Services)
type Validator struct {
	validator *validation.Validator
	logger    logging.Logger
}

// NewValidator creates a new validator instance for backend services
func NewValidator(logger logging.Logger) *Validator {
	return &Validator{
		validator: validation.NewValidator(),
		logger:    logger,
	}
}

// ValidateTransaction validates blockchain transaction data
func (v *Validator) ValidateTransaction(txData map[string]interface{}) error {
	v.logger.Debug("Validating transaction data",
		logging.String("service", constants.ServiceBlockchainCore))

	// Transaction validation for backend services
	if txData["hash"] == nil || txData["from"] == nil || txData["to"] == nil {
		v.logger.Error("Transaction validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("txData", fmt.Sprintf("%+v", txData)))
		return &validation.ValidationError{
			Field:    "transaction",
			Tag:      "transaction",
			Value:    txData,
			Message:  "missing required transaction fields",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("Transaction validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// ValidateModelInput validates ML model input data
func (v *Validator) ValidateModelInput(inputData map[string]interface{}) error {
	v.logger.Debug("Validating model input data",
		logging.String("service", constants.ServiceBlockchainCore))

	// Model input validation for AI services
	if inputData["features"] == nil || inputData["model_id"] == nil {
		v.logger.Error("Model input validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("inputData", fmt.Sprintf("%+v", inputData)))
		return &validation.ValidationError{
			Field:    "model_input",
			Tag:      "model_input",
			Value:    inputData,
			Message:  "missing required model input fields",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("Model input validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// ValidateQuery validates search query parameters
func (v *Validator) ValidateQuery(queryParams map[string]interface{}) error {
	v.logger.Debug("Validating query parameters",
		logging.String("service", constants.ServiceBlockchainCore))

	// Query validation for search services
	if queryParams["query"] == nil || queryParams["filters"] == nil {
		v.logger.Error("Query validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("queryParams", fmt.Sprintf("%+v", queryParams)))
		return &validation.ValidationError{
			Field:    "query",
			Tag:      "query",
			Value:    queryParams,
			Message:  "missing required query parameters",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("Query validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// ValidateData validates analytics data structure
func (v *Validator) ValidateData(data map[string]interface{}) error {
	v.logger.Debug("Validating analytics data",
		logging.String("service", constants.ServiceBlockchainCore))

	// Data validation for analytics services
	if data["metrics"] == nil || data["timestamp"] == nil {
		v.logger.Error("Data validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("data", fmt.Sprintf("%+v", data)))
		return &validation.ValidationError{
			Field:    "data",
			Tag:      "data",
			Value:    data,
			Message:  "missing required data fields",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("Data validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// ValidateBlock validates blockchain block structure
func (v *Validator) ValidateBlock(blockData map[string]interface{}) error {
	v.logger.Debug("Validating block data",
		logging.String("service", constants.ServiceBlockchainCore))

	// Block validation for blockchain services
	if blockData["hash"] == nil || blockData["previous_hash"] == nil || blockData["transactions"] == nil {
		v.logger.Error("Block validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("blockData", fmt.Sprintf("%+v", blockData)))
		return &validation.ValidationError{
			Field:    "block",
			Tag:      "block",
			Value:    blockData,
			Message:  "missing required block fields",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("Block validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// ValidateStruct validates a struct using tags
func (v *Validator) ValidateStruct(s interface{}) error {
	v.logger.Debug("Validating struct",
		logging.String("service", constants.ServiceBlockchainCore))

	errors := v.validator.Validate(s)
	if errors.HasErrors() {
		v.logger.Error("Struct validation failed",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("errors", errors.Error()))
		return errors
	}

	v.logger.Debug("Struct validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// SanitizeInput sanitizes user input
func (v *Validator) SanitizeInput(input string) string {
	v.logger.Debug("Sanitizing input",
		logging.String("service", constants.ServiceBlockchainCore))

	// Simple sanitization - implement proper sanitization based on your requirements
	sanitized := strings.TrimSpace(input)

	v.logger.Debug("Input sanitized",
		logging.String("service", constants.ServiceBlockchainCore))

	return sanitized
}

// ValidateBlockNumber validates blockchain block number
func (v *Validator) ValidateBlockNumber(blockNumber int64) error {
	if blockNumber <= 0 {
		return &validation.ValidationError{
			Field:    "block_number",
			Tag:      "block_number",
			Value:    blockNumber,
			Message:  "block number must be greater than 0",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateBlockHash validates blockchain block hash
func (v *Validator) ValidateBlockHash(blockHash string) error {
	if blockHash == "" {
		return &validation.ValidationError{
			Field:    "block_hash",
			Tag:      "block_hash",
			Value:    blockHash,
			Message:  "block hash cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be hex string with reasonable length
	if len(blockHash) < 32 || len(blockHash) > 128 {
		return &validation.ValidationError{
			Field:    "block_hash",
			Tag:      "block_hash",
			Value:    blockHash,
			Message:  "block hash must be between 32 and 128 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateValidatorId validates validator ID
func (v *Validator) ValidateValidatorId(validatorId string) error {
	if validatorId == "" {
		return &validation.ValidationError{
			Field:    "validator_id",
			Tag:      "validator_id",
			Value:    validatorId,
			Message:  "validator_id cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateWalletAddress validates blockchain wallet address
func (v *Validator) ValidateWalletAddress(address string) error {
	if address == "" {
		return &validation.ValidationError{
			Field:    "wallet_address",
			Tag:      "wallet_address",
			Value:    address,
			Message:  "wallet address cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be hex string with reasonable length
	if len(address) < 20 || len(address) > 128 {
		return &validation.ValidationError{
			Field:    "wallet_address",
			Tag:      "wallet_address",
			Value:    address,
			Message:  "wallet address must be between 20 and 128 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateAmount validates transaction amount
func (v *Validator) ValidateAmount(amount string) error {
	if amount == "" {
		return &validation.ValidationError{
			Field:    "amount",
			Tag:      "amount",
			Value:    amount,
			Message:  "amount cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be a valid number string
	// More comprehensive validation can be added (decimal places, precision, etc.)
	if len(amount) > 50 {
		return &validation.ValidationError{
			Field:    "amount",
			Tag:      "amount",
			Value:    amount,
			Message:  "amount must be less than 50 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateTransactionHash validates transaction hash
func (v *Validator) ValidateTransactionHash(hash string) error {
	if hash == "" {
		return &validation.ValidationError{
			Field:    "transaction_hash",
			Tag:      "transaction_hash",
			Value:    hash,
			Message:  "transaction hash cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be hex string with reasonable length
	if len(hash) < 32 || len(hash) > 128 {
		return &validation.ValidationError{
			Field:    "transaction_hash",
			Tag:      "transaction_hash",
			Value:    hash,
			Message:  "transaction hash must be between 32 and 128 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateGasPrice validates gas price
func (v *Validator) ValidateGasPrice(gasPrice string) error {
	if gasPrice == "" {
		return &validation.ValidationError{
			Field:    "gas_price",
			Tag:      "gas_price",
			Value:    gasPrice,
			Message:  "gas price cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be a valid number string
	if len(gasPrice) > 50 {
		return &validation.ValidationError{
			Field:    "gas_price",
			Tag:      "gas_price",
			Value:    gasPrice,
			Message:  "gas price must be less than 50 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateGasLimit validates gas limit
func (v *Validator) ValidateGasLimit(gasLimit int64) error {
	if gasLimit <= 0 {
		return &validation.ValidationError{
			Field:    "gas_limit",
			Tag:      "gas_limit",
			Value:    gasLimit,
			Message:  "gas limit must be greater than 0",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Maximum gas limit validation (reasonable upper bound)
	if gasLimit > 100000000 {
		return &validation.ValidationError{
			Field:    "gas_limit",
			Tag:      "gas_limit",
			Value:    gasLimit,
			Message:  "gas limit must be less than 100,000,000",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateContractAddress validates smart contract address
func (v *Validator) ValidateContractAddress(address string) error {
	if address == "" {
		return &validation.ValidationError{
			Field:    "contract_address",
			Tag:      "contract_address",
			Value:    address,
			Message:  "contract address cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be hex string with reasonable length
	if len(address) < 20 || len(address) > 128 {
		return &validation.ValidationError{
			Field:    "contract_address",
			Tag:      "contract_address",
			Value:    address,
			Message:  "contract address must be between 20 and 128 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateTokenId validates NFT token ID
func (v *Validator) ValidateTokenId(tokenId string) error {
	if tokenId == "" {
		return &validation.ValidationError{
			Field:    "token_id",
			Tag:      "token_id",
			Value:    tokenId,
			Message:  "token_id cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be reasonable length
	if len(tokenId) > 256 {
		return &validation.ValidationError{
			Field:    "token_id",
			Tag:      "token_id",
			Value:    tokenId,
			Message:  "token_id must be less than 256 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateCertificateId validates product certificate ID
func (v *Validator) ValidateCertificateId(certificateId string) error {
	if certificateId == "" {
		return &validation.ValidationError{
			Field:    "certificate_id",
			Tag:      "certificate_id",
			Value:    certificateId,
			Message:  "certificate_id cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be reasonable length
	if len(certificateId) > 256 {
		return &validation.ValidationError{
			Field:    "certificate_id",
			Tag:      "certificate_id",
			Value:    certificateId,
			Message:  "certificate_id must be less than 256 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateProductId validates product ID
func (v *Validator) ValidateProductId(productId string) error {
	if productId == "" {
		return &validation.ValidationError{
			Field:    "product_id",
			Tag:      "product_id",
			Value:    productId,
			Message:  "product_id cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be reasonable length
	if len(productId) > 256 {
		return &validation.ValidationError{
			Field:    "product_id",
			Tag:      "product_id",
			Value:    productId,
			Message:  "product_id must be less than 256 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateValidatorAddress validates validator address
func (v *Validator) ValidateValidatorAddress(address string) error {
	if address == "" {
		return &validation.ValidationError{
			Field:    "validator_address",
			Tag:      "validator_address",
			Value:    address,
			Message:  "validator_address cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}
	// Basic format validation: should be hex string with reasonable length
	if len(address) < 20 || len(address) > 128 {
		return &validation.ValidationError{
			Field:    "validator_address",
			Tag:      "validator_address",
			Value:    address,
			Message:  "validator_address must be between 20 and 128 characters",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}
	return nil
}

// ValidateUUID validates UUID format (Tableau style - strict validation)
func (v *Validator) ValidateUUID(uuidStr string) error {
	v.logger.Debug("Validating UUID",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("uuid", uuidStr))

	if uuidStr == "" {
		v.logger.Error("UUID validation failed - empty",
			logging.String("service", constants.ServiceBlockchainCore))
		return &validation.ValidationError{
			Field:    "uuid",
			Tag:      "uuid",
			Value:    uuidStr,
			Message:  "UUID cannot be empty",
			Type:     validation.ErrorTypeRequired,
			Severity: validation.SeverityError,
		}
	}

	// Try to parse as UUID using github.com/google/uuid
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		v.logger.Error("UUID validation failed - invalid format",
			logging.String("service", constants.ServiceBlockchainCore),
			logging.String("uuid", uuidStr),
			logging.Error(err))
		return &validation.ValidationError{
			Field:    "uuid",
			Tag:      "uuid",
			Value:    uuidStr,
			Message:  "invalid UUID format",
			Type:     validation.ErrorTypeFormat,
			Severity: validation.SeverityError,
		}
	}

	v.logger.Debug("UUID validation successful",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}
