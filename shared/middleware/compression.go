// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

// CompressionConfig represents compression configuration
type CompressionConfig struct {
	Level            int      `mapstructure:"level"`
	MinSize          int      `mapstructure:"min_size"`
	ContentTypes     []string `mapstructure:"content_types"`
	ExcludePaths     []string `mapstructure:"exclude_paths"`
	VaryHeader       bool     `mapstructure:"vary_header"`
	ForceCompression bool     `mapstructure:"force_compression"`
}

// DefaultCompressionConfig returns the default compression configuration
func DefaultCompressionConfig() CompressionConfig {
	return CompressionConfig{
		Level:            gzip.DefaultCompression,
		MinSize:          1024, // 1KB
		ContentTypes:     []string{"text/plain", "text/html", "text/css", "text/javascript", "application/json", "application/xml", "application/javascript", "application/x-javascript"},
		ExcludePaths:     []string{},
		VaryHeader:       true,
		ForceCompression: false,
	}
}

// CompressionMiddleware provides HTTP compression middleware
type CompressionMiddleware struct {
	config CompressionConfig
	pool   sync.Pool
}

// NewCompressionMiddleware creates a new compression middleware
func NewCompressionMiddleware(config CompressionConfig) *CompressionMiddleware {
	cm := &CompressionMiddleware{
		config: config,
		pool: sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(nil, config.Level)
				return w
			},
		},
	}

	return cm
}

// Middleware returns the compression middleware
func (m *CompressionMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if compression should be applied
			if !m.shouldCompress(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Create compressed response writer
			compressedWriter := &CompressedResponseWriter{
				ResponseWriter: w,
				config:         m.config,
				pool:           &m.pool,
			}

			// Call next handler
			next.ServeHTTP(compressedWriter, r)

			// Close the compressed writer
			compressedWriter.Close()
		})
	}
}

// shouldCompress checks if compression should be applied
func (m *CompressionMiddleware) shouldCompress(r *http.Request) bool {
	// Check if client accepts gzip
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	// Check if path should be excluded
	for _, path := range m.config.ExcludePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return false
		}
	}

	// Check if already compressed
	if r.Header.Get("Content-Encoding") != "" {
		return false
	}

	return true
}

// CompressedResponseWriter provides compressed response writing
type CompressedResponseWriter struct {
	http.ResponseWriter
	config     CompressionConfig
	pool       *sync.Pool
	writer     *gzip.Writer
	written    bool
	statusCode int
}

// WriteHeader captures the status code
func (crw *CompressedResponseWriter) WriteHeader(code int) {
	if !crw.written {
		crw.statusCode = code
		crw.written = true
	}
	crw.ResponseWriter.WriteHeader(code)
}

// Write writes compressed data
func (crw *CompressedResponseWriter) Write(data []byte) (int, error) {
	if !crw.written {
		crw.WriteHeader(http.StatusOK)
	}

	// Check if we should compress
	if !crw.shouldCompress(data) {
		return crw.ResponseWriter.Write(data)
	}

	// Initialize gzip writer if needed
	if crw.writer == nil {
		crw.writer = crw.pool.Get().(*gzip.Writer)
		crw.writer.Reset(crw.ResponseWriter)

		// Set headers
		crw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		if crw.config.VaryHeader {
			crw.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
		}
	}

	return crw.writer.Write(data)
}

// shouldCompress checks if data should be compressed
func (crw *CompressedResponseWriter) shouldCompress(data []byte) bool {
	// Check minimum size
	if len(data) < crw.config.MinSize {
		return false
	}

	// Check content type
	contentType := crw.ResponseWriter.Header().Get("Content-Type")
	return crw.isCompressibleContentType(contentType)
}

// isCompressibleContentType checks if content type is compressible
func (crw *CompressedResponseWriter) isCompressibleContentType(contentType string) bool {
	// If no content types specified, compress all
	if len(crw.config.ContentTypes) == 0 {
		return true
	}

	// Check if content type is in the list
	for _, ct := range crw.config.ContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}

	return false
}

// Close closes the compressed writer
func (crw *CompressedResponseWriter) Close() error {
	if crw.writer != nil {
		err := crw.writer.Close()
		crw.pool.Put(crw.writer)
		crw.writer = nil
		return err
	}
	return nil
}

// CompressionHandler provides a compression handler
type CompressionHandler struct {
	config CompressionConfig
	pool   sync.Pool
}

// NewCompressionHandler creates a new compression handler
func NewCompressionHandler(config CompressionConfig) *CompressionHandler {
	ch := &CompressionHandler{
		config: config,
		pool: sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(nil, config.Level)
				return w
			},
		},
	}

	return ch
}

// HandleCompression handles compression for a request
func (h *CompressionHandler) HandleCompression(w http.ResponseWriter, r *http.Request, data []byte) (int, error) {
	// Check if compression should be applied
	if !h.shouldCompress(r, data) {
		return w.Write(data)
	}

	// Get gzip writer from pool
	writer := h.pool.Get().(*gzip.Writer)
	defer h.pool.Put(writer)

	// Reset writer
	writer.Reset(w)

	// Set headers
	w.Header().Set("Content-Encoding", "gzip")
	if h.config.VaryHeader {
		w.Header().Set("Vary", "Accept-Encoding")
	}

	// Write compressed data
	n, err := writer.Write(data)
	if err != nil {
		return 0, err
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return 0, err
	}

	return n, nil
}

// shouldCompress checks if compression should be applied
func (h *CompressionHandler) shouldCompress(r *http.Request, data []byte) bool {
	// Check if client accepts gzip
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	// Check if path should be excluded
	for _, path := range h.config.ExcludePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return false
		}
	}

	// Check minimum size
	if len(data) < h.config.MinSize {
		return false
	}

	// Check content type
	contentType := r.Header.Get("Content-Type")
	return h.isCompressibleContentType(contentType)
}

// isCompressibleContentType checks if content type is compressible
func (h *CompressionHandler) isCompressibleContentType(contentType string) bool {
	// If no content types specified, compress all
	if len(h.config.ContentTypes) == 0 {
		return true
	}

	// Check if content type is in the list
	for _, ct := range h.config.ContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}

	return false
}

// CompressionValidator provides compression validation
type CompressionValidator struct {
	config CompressionConfig
}

// NewCompressionValidator creates a new compression validator
func NewCompressionValidator(config CompressionConfig) *CompressionValidator {
	return &CompressionValidator{
		config: config,
	}
}

// ValidateCompression validates if compression should be applied
func (v *CompressionValidator) ValidateCompression(r *http.Request, data []byte) bool {
	// Check if client accepts gzip
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	// Check if path should be excluded
	for _, path := range v.config.ExcludePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return false
		}
	}

	// Check minimum size
	if len(data) < v.config.MinSize {
		return false
	}

	// Check content type
	contentType := r.Header.Get("Content-Type")
	return v.isCompressibleContentType(contentType)
}

// isCompressibleContentType checks if content type is compressible
func (v *CompressionValidator) isCompressibleContentType(contentType string) bool {
	// If no content types specified, compress all
	if len(v.config.ContentTypes) == 0 {
		return true
	}

	// Check if content type is in the list
	for _, ct := range v.config.ContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}

	return false
}

// CompressionConfigBuilder provides a builder for compression configuration
type CompressionConfigBuilder struct {
	config CompressionConfig
}

// NewCompressionConfigBuilder creates a new compression config builder
func NewCompressionConfigBuilder() *CompressionConfigBuilder {
	return &CompressionConfigBuilder{
		config: DefaultCompressionConfig(),
	}
}

// WithLevel sets the compression level
func (b *CompressionConfigBuilder) WithLevel(level int) *CompressionConfigBuilder {
	b.config.Level = level
	return b
}

// WithMinSize sets the minimum size for compression
func (b *CompressionConfigBuilder) WithMinSize(size int) *CompressionConfigBuilder {
	b.config.MinSize = size
	return b
}

// WithContentTypes sets the content types to compress
func (b *CompressionConfigBuilder) WithContentTypes(types ...string) *CompressionConfigBuilder {
	b.config.ContentTypes = types
	return b
}

// WithExcludePaths sets the paths to exclude from compression
func (b *CompressionConfigBuilder) WithExcludePaths(paths ...string) *CompressionConfigBuilder {
	b.config.ExcludePaths = paths
	return b
}

// WithVaryHeader sets whether to include Vary header
func (b *CompressionConfigBuilder) WithVaryHeader(vary bool) *CompressionConfigBuilder {
	b.config.VaryHeader = vary
	return b
}

// WithForceCompression sets whether to force compression
func (b *CompressionConfigBuilder) WithForceCompression(force bool) *CompressionConfigBuilder {
	b.config.ForceCompression = force
	return b
}

// Build builds the compression configuration
func (b *CompressionConfigBuilder) Build() CompressionConfig {
	return b.config
}

// CompressionMetrics provides compression metrics
type CompressionMetrics struct {
	TotalRequests      int64   `json:"total_requests"`
	CompressedRequests int64   `json:"compressed_requests"`
	TotalBytes         int64   `json:"total_bytes"`
	CompressedBytes    int64   `json:"compressed_bytes"`
	CompressionRatio   float64 `json:"compression_ratio"`
}

// CompressionMetricsCollector collects compression metrics
type CompressionMetricsCollector struct {
	metrics *CompressionMetrics
}

// NewCompressionMetricsCollector creates a new compression metrics collector
func NewCompressionMetricsCollector() *CompressionMetricsCollector {
	return &CompressionMetricsCollector{
		metrics: &CompressionMetrics{},
	}
}

// RecordRequest records a request
func (cmc *CompressionMetricsCollector) RecordRequest(compressed bool, originalSize, compressedSize int64) {
	cmc.metrics.TotalRequests++
	cmc.metrics.TotalBytes += originalSize

	if compressed {
		cmc.metrics.CompressedRequests++
		cmc.metrics.CompressedBytes += compressedSize
	}

	// Calculate compression ratio
	if cmc.metrics.TotalBytes > 0 {
		cmc.metrics.CompressionRatio = float64(cmc.metrics.CompressedBytes) / float64(cmc.metrics.TotalBytes)
	}
}

// GetMetrics returns current compression metrics
func (cmc *CompressionMetricsCollector) GetMetrics() *CompressionMetrics {
	return cmc.metrics
}

// Reset resets compression metrics
func (cmc *CompressionMetricsCollector) Reset() {
	cmc.metrics = &CompressionMetrics{}
}
