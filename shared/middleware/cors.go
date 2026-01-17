// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins         []string `mapstructure:"allowed_origins"`
	AllowedMethods         []string `mapstructure:"allowed_methods"`
	AllowedHeaders         []string `mapstructure:"allowed_headers"`
	ExposedHeaders         []string `mapstructure:"exposed_headers"`
	AllowCredentials       bool     `mapstructure:"allow_credentials"`
	MaxAge                 int      `mapstructure:"max_age"`
	AllowWildcard          bool     `mapstructure:"allow_wildcard"`
	AllowBrowserExtensions bool     `mapstructure:"allow_browser_extensions"`
	AllowWebSockets        bool     `mapstructure:"allow_web_sockets"`
	AllowFiles             bool     `mapstructure:"allow_files"`
}

// DefaultCORSConfig returns the default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:         []string{"*"},
		AllowedMethods:         []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowedHeaders:         []string{"*"},
		ExposedHeaders:         []string{},
		AllowCredentials:       false,
		MaxAge:                 86400, // 24 hours
		AllowWildcard:          true,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}
}

// CORSMiddleware provides CORS middleware
type CORSMiddleware struct {
	config CORSConfig
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
	}
}

// Middleware returns the CORS middleware
func (m *CORSMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle preflight requests
			if r.Method == http.MethodOptions {
				m.handlePreflight(w, r)
				return
			}

			// Handle actual requests
			m.handleRequest(w, r)
			next.ServeHTTP(w, r)
		})
	}
}

// handlePreflight handles CORS preflight requests
func (m *CORSMiddleware) handlePreflight(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Check if origin is allowed
	if !m.isOriginAllowed(origin) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Set CORS headers
	m.setCORSHeaders(w, origin)

	// Set allowed methods
	if len(m.config.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
	}

	// Set allowed headers
	if len(m.config.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))
	}

	// Set max age
	if m.config.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", string(rune(m.config.MaxAge)))
	}

	w.WriteHeader(http.StatusOK)
}

// handleRequest handles CORS for actual requests
func (m *CORSMiddleware) handleRequest(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Check if origin is allowed
	if !m.isOriginAllowed(origin) {
		return
	}

	// Set CORS headers
	m.setCORSHeaders(w, origin)
}

// isOriginAllowed checks if an origin is allowed
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	// Check for wildcard
	if m.config.AllowWildcard && len(m.config.AllowedOrigins) == 1 && m.config.AllowedOrigins[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedOrigin := range m.config.AllowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	// Check for browser extensions
	if m.config.AllowBrowserExtensions && m.isBrowserExtension(origin) {
		return true
	}

	// Check for web sockets
	if m.config.AllowWebSockets && m.isWebSocket(origin) {
		return true
	}

	// Check for file protocol
	if m.config.AllowFiles && m.isFileProtocol(origin) {
		return true
	}

	return false
}

// setCORSHeaders sets CORS headers
func (m *CORSMiddleware) setCORSHeaders(w http.ResponseWriter, origin string) {
	// Set origin
	if m.config.AllowWildcard && len(m.config.AllowedOrigins) == 1 && m.config.AllowedOrigins[0] == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Set credentials
	if m.config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Set exposed headers
	if len(m.config.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.config.ExposedHeaders, ", "))
	}
}

// isBrowserExtension checks if the origin is a browser extension
func (m *CORSMiddleware) isBrowserExtension(origin string) bool {
	return strings.HasPrefix(origin, "chrome-extension://") ||
		strings.HasPrefix(origin, "moz-extension://") ||
		strings.HasPrefix(origin, "safari-extension://") ||
		strings.HasPrefix(origin, "ms-browser-extension://")
}

// isWebSocket checks if the origin is a web socket
func (m *CORSMiddleware) isWebSocket(origin string) bool {
	return strings.HasPrefix(origin, "ws://") || strings.HasPrefix(origin, "wss://")
}

// isFileProtocol checks if the origin is a file protocol
func (m *CORSMiddleware) isFileProtocol(origin string) bool {
	return strings.HasPrefix(origin, "file://")
}

// CORSHandler provides a CORS handler
type CORSHandler struct {
	config CORSConfig
}

// NewCORSHandler creates a new CORS handler
func NewCORSHandler(config CORSConfig) *CORSHandler {
	return &CORSHandler{
		config: config,
	}
}

// HandleCORS handles CORS for a request
func (h *CORSHandler) HandleCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Check if origin is allowed
	if !h.isOriginAllowed(origin) {
		return false
	}

	// Set CORS headers
	h.setCORSHeaders(w, origin)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		h.handlePreflight(w, r)
		return true
	}

	return true
}

// isOriginAllowed checks if an origin is allowed
func (h *CORSHandler) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	// Check for wildcard
	if h.config.AllowWildcard && len(h.config.AllowedOrigins) == 1 && h.config.AllowedOrigins[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedOrigin := range h.config.AllowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	// Check for browser extensions
	if h.config.AllowBrowserExtensions && h.isBrowserExtension(origin) {
		return true
	}

	// Check for web sockets
	if h.config.AllowWebSockets && h.isWebSocket(origin) {
		return true
	}

	// Check for file protocol
	if h.config.AllowFiles && h.isFileProtocol(origin) {
		return true
	}

	return false
}

// setCORSHeaders sets CORS headers
func (h *CORSHandler) setCORSHeaders(w http.ResponseWriter, origin string) {
	// Set origin
	if h.config.AllowWildcard && len(h.config.AllowedOrigins) == 1 && h.config.AllowedOrigins[0] == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Set credentials
	if h.config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Set exposed headers
	if len(h.config.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(h.config.ExposedHeaders, ", "))
	}
}

// handlePreflight handles CORS preflight requests
func (h *CORSHandler) handlePreflight(w http.ResponseWriter, r *http.Request) {
	// Set allowed methods
	if len(h.config.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(h.config.AllowedMethods, ", "))
	}

	// Set allowed headers
	if len(h.config.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(h.config.AllowedHeaders, ", "))
	}

	// Set max age
	if h.config.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", string(rune(h.config.MaxAge)))
	}

	w.WriteHeader(http.StatusOK)
}

// isBrowserExtension checks if the origin is a browser extension
func (h *CORSHandler) isBrowserExtension(origin string) bool {
	return strings.HasPrefix(origin, "chrome-extension://") ||
		strings.HasPrefix(origin, "moz-extension://") ||
		strings.HasPrefix(origin, "safari-extension://") ||
		strings.HasPrefix(origin, "ms-browser-extension://")
}

// isWebSocket checks if the origin is a web socket
func (h *CORSHandler) isWebSocket(origin string) bool {
	return strings.HasPrefix(origin, "ws://") || strings.HasPrefix(origin, "wss://")
}

// isFileProtocol checks if the origin is a file protocol
func (h *CORSHandler) isFileProtocol(origin string) bool {
	return strings.HasPrefix(origin, "file://")
}

// CORSValidator provides CORS validation
type CORSValidator struct {
	config CORSConfig
}

// NewCORSValidator creates a new CORS validator
func NewCORSValidator(config CORSConfig) *CORSValidator {
	return &CORSValidator{
		config: config,
	}
}

// ValidateOrigin validates an origin
func (v *CORSValidator) ValidateOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	// Check for wildcard
	if v.config.AllowWildcard && len(v.config.AllowedOrigins) == 1 && v.config.AllowedOrigins[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedOrigin := range v.config.AllowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	// Check for browser extensions
	if v.config.AllowBrowserExtensions && v.isBrowserExtension(origin) {
		return true
	}

	// Check for web sockets
	if v.config.AllowWebSockets && v.isWebSocket(origin) {
		return true
	}

	// Check for file protocol
	if v.config.AllowFiles && v.isFileProtocol(origin) {
		return true
	}

	return false
}

// ValidateMethod validates a method
func (v *CORSValidator) ValidateMethod(method string) bool {
	if method == "" {
		return false
	}

	// Check for wildcard
	if len(v.config.AllowedMethods) == 1 && v.config.AllowedMethods[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedMethod := range v.config.AllowedMethods {
		if method == allowedMethod {
			return true
		}
	}

	return false
}

// ValidateHeader validates a header
func (v *CORSValidator) ValidateHeader(header string) bool {
	if header == "" {
		return false
	}

	// Check for wildcard
	if len(v.config.AllowedHeaders) == 1 && v.config.AllowedHeaders[0] == "*" {
		return true
	}

	// Check exact matches
	for _, allowedHeader := range v.config.AllowedHeaders {
		if header == allowedHeader {
			return true
		}
	}

	return false
}

// isBrowserExtension checks if the origin is a browser extension
func (v *CORSValidator) isBrowserExtension(origin string) bool {
	return strings.HasPrefix(origin, "chrome-extension://") ||
		strings.HasPrefix(origin, "moz-extension://") ||
		strings.HasPrefix(origin, "safari-extension://") ||
		strings.HasPrefix(origin, "ms-browser-extension://")
}

// isWebSocket checks if the origin is a web socket
func (v *CORSValidator) isWebSocket(origin string) bool {
	return strings.HasPrefix(origin, "ws://") || strings.HasPrefix(origin, "wss://")
}

// isFileProtocol checks if the origin is a file protocol
func (v *CORSValidator) isFileProtocol(origin string) bool {
	return strings.HasPrefix(origin, "file://")
}

// CORSConfigBuilder provides a builder for CORS configuration
type CORSConfigBuilder struct {
	config CORSConfig
}

// NewCORSConfigBuilder creates a new CORS config builder
func NewCORSConfigBuilder() *CORSConfigBuilder {
	return &CORSConfigBuilder{
		config: DefaultCORSConfig(),
	}
}

// WithOrigins sets allowed origins
func (b *CORSConfigBuilder) WithOrigins(origins ...string) *CORSConfigBuilder {
	b.config.AllowedOrigins = origins
	return b
}

// WithMethods sets allowed methods
func (b *CORSConfigBuilder) WithMethods(methods ...string) *CORSConfigBuilder {
	b.config.AllowedMethods = methods
	return b
}

// WithHeaders sets allowed headers
func (b *CORSConfigBuilder) WithHeaders(headers ...string) *CORSConfigBuilder {
	b.config.AllowedHeaders = headers
	return b
}

// WithExposedHeaders sets exposed headers
func (b *CORSConfigBuilder) WithExposedHeaders(headers ...string) *CORSConfigBuilder {
	b.config.ExposedHeaders = headers
	return b
}

// WithCredentials sets allow credentials
func (b *CORSConfigBuilder) WithCredentials(allow bool) *CORSConfigBuilder {
	b.config.AllowCredentials = allow
	return b
}

// WithMaxAge sets max age
func (b *CORSConfigBuilder) WithMaxAge(maxAge int) *CORSConfigBuilder {
	b.config.MaxAge = maxAge
	return b
}

// WithWildcard sets allow wildcard
func (b *CORSConfigBuilder) WithWildcard(allow bool) *CORSConfigBuilder {
	b.config.AllowWildcard = allow
	return b
}

// WithBrowserExtensions sets allow browser extensions
func (b *CORSConfigBuilder) WithBrowserExtensions(allow bool) *CORSConfigBuilder {
	b.config.AllowBrowserExtensions = allow
	return b
}

// WithWebSockets sets allow web sockets
func (b *CORSConfigBuilder) WithWebSockets(allow bool) *CORSConfigBuilder {
	b.config.AllowWebSockets = allow
	return b
}

// WithFiles sets allow files
func (b *CORSConfigBuilder) WithFiles(allow bool) *CORSConfigBuilder {
	b.config.AllowFiles = allow
	return b
}

// Build builds the CORS configuration
func (b *CORSConfigBuilder) Build() CORSConfig {
	return b.config
}
