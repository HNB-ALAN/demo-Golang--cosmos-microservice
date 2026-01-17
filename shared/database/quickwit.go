package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/usc-platform/shared/config"
)

// QuickwitHealthChecker implements health checking for Quickwit
type QuickwitHealthChecker struct {
	client QuickwitClient
}

// NewQuickwitHealthChecker creates a new Quickwit health checker
func NewQuickwitHealthChecker(client QuickwitClient) *QuickwitHealthChecker {
	return &QuickwitHealthChecker{client: client}
}

// Check performs a health check on Quickwit
func (h *QuickwitHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx); err != nil {
		return fmt.Errorf("Quickwit ping failed: %w", err)
	}

	return nil
}

// Name returns the name of the health checker
func (h *QuickwitHealthChecker) Name() string {
	return "quickwit"
}

// Description returns the description of the health checker
func (h *QuickwitHealthChecker) Description() string {
	return "Quickwit search engine health check"
}

// QuickwitConnection represents a Quickwit connection
type QuickwitConnection struct {
	client  *http.Client
	baseURL string
}

// NewQuickwitConnection creates a new Quickwit connection with retry logic
func NewQuickwitConnection(cfg *config.Config) (*QuickwitConnection, error) {
	url := cfg.GetQuickwitURL()
	maxRetries := 5
	baseDelay := 2 * time.Second

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := baseDelay * time.Duration(1<<attempt)
			time.Sleep(delay)
		}

		// Test connection with health check (use /api/v1/version as Quickwit doesn't have /health endpoint)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/version", url), nil)
		if err != nil {
			cancel()
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("failed to create Quickwit health check request: %w", err)
			}
			continue
		}

		resp, err := client.Do(req)
		cancel()

		if err != nil {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("failed to connect to Quickwit after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success
			return &QuickwitConnection{
				client:  client,
				baseURL: url,
			}, nil
		}

		if attempt == maxRetries-1 {
			return nil, fmt.Errorf("Quickwit health check failed with status %d after %d attempts", resp.StatusCode, maxRetries)
		}
	}

	return nil, fmt.Errorf("failed to create Quickwit connection after %d attempts", maxRetries)
}

// Client returns the underlying HTTP client
func (q *QuickwitConnection) Client() *http.Client {
	return q.client
}

// Ping tests the connection (use /api/v1/version as Quickwit doesn't have /health endpoint)
func (q *QuickwitConnection) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/version", q.baseURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create ping request: %w", err)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return fmt.Errorf("ping request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ping failed with status %d", resp.StatusCode)
	}

	return nil
}

// Index indexes a document
func (q *QuickwitConnection) Index(ctx context.Context, index string, document interface{}) error {
	// Quickwit uses /api/v1/{index}/ingest endpoint
	url := fmt.Sprintf("%s/api/v1/%s/ingest", q.baseURL, index)

	var body io.Reader
	if docBytes, ok := document.([]byte); ok {
		body = bytes.NewReader(docBytes)
	} else {
		docJSON, err := json.Marshal(document)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		body = bytes.NewReader(docJSON)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create index request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)
	if err != nil {
		return fmt.Errorf("index request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("index failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Search performs a search query
func (q *QuickwitConnection) Search(ctx context.Context, index string, query interface{}) (*SearchResult, error) {
	// Quickwit uses /api/v1/{index}/search endpoint
	url := fmt.Sprintf("%s/api/v1/%s/search", q.baseURL, index)

	var body io.Reader
	if queryBytes, ok := query.([]byte); ok {
		body = bytes.NewReader(queryBytes)
	} else {
		queryJSON, err := json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query: %w", err)
		}
		body = bytes.NewReader(queryJSON)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return &SearchResult{
		response: bodyBytes,
	}, nil
}

// SearchResult represents the result of a Quickwit search
type SearchResult struct {
	response []byte
}

// IsError returns true if the search resulted in an error
func (s *SearchResult) IsError() bool {
	// Simple error check - in production, implement proper error checking
	return false
}

// String returns the response body as a string
func (s *SearchResult) String() string {
	return string(s.response)
}

// Bytes returns the response body as bytes
func (s *SearchResult) Bytes() []byte {
	return s.response
}

// Close closes the search result
func (s *SearchResult) Close() error {
	// Simple close - in production, implement proper close
	return nil
}

// Close closes the connection
func (q *QuickwitConnection) Close() error {
	// HTTP client doesn't need explicit close
	return nil
}

// HealthCheck performs a health check
func (q *QuickwitConnection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := q.Ping(ctx); err != nil {
		return fmt.Errorf("Quickwit ping failed: %w", err)
	}

	return nil
}

// initializeQuickwit initializes Quickwit connection with retry logic
func (m *DatabaseManager) initializeQuickwit() error {
	url := m.config.GetQuickwitURL()
	maxRetries := 5
	baseDelay := 2 * time.Second

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := baseDelay * time.Duration(1<<attempt)
			time.Sleep(delay)
		}

		// Test connection with health check (use /api/v1/version as Quickwit doesn't have /health endpoint)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/version", url), nil)
		if err != nil {
			cancel()
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to create Quickwit health check request: %w", err)
			}
			continue
		}

		resp, err := client.Do(req)
		cancel()

		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to connect to Quickwit after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success - store connection
			m.quickwit = &QuickwitConnection{
				client:  client,
				baseURL: url,
			}
			return nil
		}

		if attempt == maxRetries-1 {
			return fmt.Errorf("Quickwit health check failed with status %d after %d attempts", resp.StatusCode, maxRetries)
		}
	}

	return fmt.Errorf("failed to initialize Quickwit after %d attempts", maxRetries)
}

// IsQuickwitError checks if an error is Quickwit-specific
func IsQuickwitError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common Quickwit error patterns
	errStr := err.Error()
	quickwitErrors := []string{
		"Quickwit",
		"quickwit",
		"index",
		"query",
		"search",
		"document",
		"ingest",
	}

	for _, pattern := range quickwitErrors {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}
