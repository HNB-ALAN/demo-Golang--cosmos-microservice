package logging

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field represents a log field
type Field = zap.Field

// String creates a string field
func String(key, value string) Field {
	return zap.String(key, value)
}

// Int creates an integer field
func Int(key string, value int) Field {
	return zap.Int(key, value)
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

// Float64 creates a float64 field
func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

// Bool creates a boolean field
func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

// Duration creates a duration field
func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

// Time creates a time field
func Time(key string, value time.Time) Field {
	return zap.Time(key, value)
}

// Error creates an error field
func Error(err error) Field {
	return zap.Error(err)
}

// Any creates an any field
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

// Binary creates a binary field
func Binary(key string, value []byte) Field {
	return zap.Binary(key, value)
}

// ByteString creates a byte string field
func ByteString(key string, value []byte) Field {
	return zap.ByteString(key, value)
}

// Complex128 creates a complex128 field
func Complex128(key string, value complex128) Field {
	return zap.Complex128(key, value)
}

// Complex64 creates a complex64 field
func Complex64(key string, value complex64) Field {
	return zap.Complex64(key, value)
}

// Float32 creates a float32 field
func Float32(key string, value float32) Field {
	return zap.Float32(key, value)
}

// Int32 creates an int32 field
func Int32(key string, value int32) Field {
	return zap.Int32(key, value)
}

// Int16 creates an int16 field
func Int16(key string, value int16) Field {
	return zap.Int16(key, value)
}

// Int8 creates an int8 field
func Int8(key string, value int8) Field {
	return zap.Int8(key, value)
}

// Uint creates a uint field
func Uint(key string, value uint) Field {
	return zap.Uint(key, value)
}

// Uint64 creates a uint64 field
func Uint64(key string, value uint64) Field {
	return zap.Uint64(key, value)
}

// Uint32 creates a uint32 field
func Uint32(key string, value uint32) Field {
	return zap.Uint32(key, value)
}

// Uint16 creates a uint16 field
func Uint16(key string, value uint16) Field {
	return zap.Uint16(key, value)
}

// Uint8 creates a uint8 field
func Uint8(key string, value uint8) Field {
	return zap.Uint8(key, value)
}

// Uintptr creates a uintptr field
func Uintptr(key string, value uintptr) Field {
	return zap.Uintptr(key, value)
}

// Reflect creates a reflect field
func Reflect(key string, value interface{}) Field {
	return zap.Reflect(key, value)
}

// Namespace creates a namespace field
func Namespace(key string) Field {
	return zap.Namespace(key)
}

// Stringer creates a stringer field
func Stringer(key string, value fmt.Stringer) Field {
	return zap.Stringer(key, value)
}

// Strings creates a strings field
func Strings(key string, value []string) Field {
	return zap.Strings(key, value)
}

// Ints creates an ints field
func Ints(key string, value []int) Field {
	return zap.Ints(key, value)
}

// Int64s creates an int64s field
func Int64s(key string, value []int64) Field {
	return zap.Int64s(key, value)
}

// Float64s creates a float64s field
func Float64s(key string, value []float64) Field {
	return zap.Float64s(key, value)
}

// Bools creates a bools field
func Bools(key string, value []bool) Field {
	return zap.Bools(key, value)
}

// Durations creates a durations field
func Durations(key string, value []time.Duration) Field {
	return zap.Durations(key, value)
}

// Times creates a times field
func Times(key string, value []time.Time) Field {
	return zap.Times(key, value)
}

// Errors creates an errors field
func Errors(key string, value []error) Field {
	return zap.Errors(key, value)
}

// Object creates an object field
func Object(key string, value zapcore.ObjectMarshaler) Field {
	return zap.Object(key, value)
}

// Array creates an array field
func Array(key string, value zapcore.ArrayMarshaler) Field {
	return zap.Array(key, value)
}

// Skip creates a skip field
func Skip() Field {
	return zap.Skip()
}

// Stack creates a stack field
func Stack(key string) Field {
	return zap.Stack(key)
}

// StackSkip creates a stack skip field
func StackSkip(key string, skip int) Field {
	return zap.StackSkip(key, skip)
}

// Common fields for consistent logging across services
var (
	// Service fields
	ServiceName    = func(name string) Field { return String("service", name) }
	ServiceVersion = func(version string) Field { return String("version", version) }
	ServiceID      = func(id string) Field { return String("service_id", id) }

	// Request fields
	RequestID     = func(id string) Field { return String("request_id", id) }
	RequestMethod = func(method string) Field { return String("method", method) }
	RequestPath   = func(path string) Field { return String("path", path) }
	RequestSize   = func(size int64) Field { return Int64("request_size", size) }
	ResponseSize  = func(size int64) Field { return Int64("response_size", size) }

	// Timing fields
	DurationField = func(d time.Duration) Field { return Duration("duration", d) }
	StartTime     = func(t time.Time) Field { return Time("start_time", t) }
	EndTime       = func(t time.Time) Field { return Time("end_time", t) }

	// Database fields
	DatabaseName  = func(name string) Field { return String("database", name) }
	TableName     = func(name string) Field { return String("table", name) }
	QueryType     = func(queryType string) Field { return String("query_type", queryType) }
	QueryDuration = func(d time.Duration) Field { return Duration("query_duration", d) }

	// User fields
	UserID    = func(id string) Field { return String("user_id", id) }
	UserEmail = func(email string) Field { return String("user_email", email) }
	UserRole  = func(role string) Field { return String("user_role", role) }

	// Error fields
	ErrorCode    = func(code string) Field { return String("error_code", code) }
	ErrorMessage = func(msg string) Field { return String("error_message", msg) }
	ErrorType    = func(errorType string) Field { return String("error_type", errorType) }

	// Environment fields
	Environment = func(env string) Field { return String("environment", env) }
	Region      = func(region string) Field { return String("region", region) }
	Instance    = func(instance string) Field { return String("instance", instance) }

	// Performance fields
	MemoryUsage = func(usage int64) Field { return Int64("memory_usage", usage) }
	CPUUsage    = func(usage float64) Field { return Float64("cpu_usage", usage) }
	Goroutines  = func(count int) Field { return Int("goroutines", count) }
)
