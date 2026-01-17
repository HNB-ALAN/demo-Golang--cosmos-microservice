package errors

import (
	"testing"
)

func TestNewInvalidInputError(t *testing.T) {
	msg := "invalid input"
	err := NewInvalidInputError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeInvalidInput) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidInput, err.Code)
	}
}

func TestNewValidationError(t *testing.T) {
	msg := "validation failed"
	err := NewValidationError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeValidationFailed) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeValidationFailed {
		t.Errorf("Expected error code %s, got %s", ErrCodeValidationFailed, err.Code)
	}
}

func TestNewInternalError(t *testing.T) {
	msg := "internal error"
	err := NewInternalError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeInternal) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeInternal {
		t.Errorf("Expected error code %s, got %s", ErrCodeInternal, err.Code)
	}
}

func TestNewNotFoundError(t *testing.T) {
	msg := "not found"
	err := NewNotFoundError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeNotFound) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeNotFound {
		t.Errorf("Expected error code %s, got %s", ErrCodeNotFound, err.Code)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	msg := "unauthorized"
	err := NewUnauthorizedError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeUnauthorized) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeUnauthorized {
		t.Errorf("Expected error code %s, got %s", ErrCodeUnauthorized, err.Code)
	}
}

func TestNewForbiddenError(t *testing.T) {
	msg := "forbidden"
	err := NewForbiddenError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeForbidden) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeForbidden {
		t.Errorf("Expected error code %s, got %s", ErrCodeForbidden, err.Code)
	}
}

func TestNewConflictError(t *testing.T) {
	msg := "conflict"
	err := NewConflictError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeConflict) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeConflict {
		t.Errorf("Expected error code %s, got %s", ErrCodeConflict, err.Code)
	}
}

func TestNewTimeoutError(t *testing.T) {
	msg := "timeout"
	err := NewTimeoutError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeTimeout) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeTimeout {
		t.Errorf("Expected error code %s, got %s", ErrCodeTimeout, err.Code)
	}
}

func TestNewDatabaseError(t *testing.T) {
	msg := "database error"
	err := NewDatabaseError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeDatabaseQuery) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeDatabaseQuery {
		t.Errorf("Expected error code %s, got %s", ErrCodeDatabaseQuery, err.Code)
	}
}

func TestNewBusinessError(t *testing.T) {
	msg := "business error"
	err := NewBusinessError(msg)

	if err == nil {
		t.Fatal("Expected error to be created")
	}

	expectedMsg := string(ErrCodeOperationNotAllowed) + ": " + msg
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Code != ErrCodeOperationNotAllowed {
		t.Errorf("Expected error code %s, got %s", ErrCodeOperationNotAllowed, err.Code)
	}
}

// Benchmark tests
func BenchmarkNewInvalidInputError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewInvalidInputError("test message")
	}
}

func BenchmarkNewValidationError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewValidationError("test message")
	}
}

func BenchmarkNewInternalError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewInternalError("test message")
	}
}
