// Package constants provides application constants for USC platform services.
//
// IMPORTANT: Error codes and error messages have been moved to the errors package.
// Use github.com/usc-platform/shared/errors for all error handling.
//
// This package now focuses on:
// - HTTP status codes
// - gRPC status codes
// - Service names and ports
// - USC contract addresses
// - Configuration constants
package constants

// HTTP status codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

// gRPC status codes
const (
	GRPCStatusOK                 = 0
	GRPCStatusCancelled          = 1
	GRPCStatusUnknown            = 2
	GRPCStatusInvalidArgument    = 3
	GRPCStatusDeadlineExceeded   = 4
	GRPCStatusNotFound           = 5
	GRPCStatusAlreadyExists      = 6
	GRPCStatusPermissionDenied   = 7
	GRPCStatusResourceExhausted  = 8
	GRPCStatusFailedPrecondition = 9
	GRPCStatusAborted            = 10
	GRPCStatusOutOfRange         = 11
	GRPCStatusUnimplemented      = 12
	GRPCStatusInternal           = 13
	GRPCStatusUnavailable        = 14
	GRPCStatusDataLoss           = 15
	GRPCStatusUnauthenticated    = 16
)
