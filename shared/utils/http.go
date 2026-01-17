// Package utils provides utility functions for USC platform services.
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPUtils provides HTTP utility functions
type HTTPUtils struct {
	client *http.Client
}

// NewHTTPUtils creates a new HTTP utils instance
func NewHTTPUtils() *HTTPUtils {
	return &HTTPUtils{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewHTTPUtilsWithTimeout creates a new HTTP utils instance with custom timeout
func NewHTTPUtilsWithTimeout(timeout time.Duration) *HTTPUtils {
	return &HTTPUtils{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	Query   map[string]string
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Error      error
}

// Get makes a GET request
func (hu *HTTPUtils) Get(url string, headers map[string]string) (*HTTPResponse, error) {
	return hu.makeRequest("GET", url, headers, nil, nil)
}

// Post makes a POST request
func (hu *HTTPUtils) Post(url string, headers map[string]string, body interface{}) (*HTTPResponse, error) {
	return hu.makeRequest("POST", url, headers, body, nil)
}

// Put makes a PUT request
func (hu *HTTPUtils) Put(url string, headers map[string]string, body interface{}) (*HTTPResponse, error) {
	return hu.makeRequest("PUT", url, headers, body, nil)
}

// Delete makes a DELETE request
func (hu *HTTPUtils) Delete(url string, headers map[string]string) (*HTTPResponse, error) {
	return hu.makeRequest("DELETE", url, headers, nil, nil)
}

// Patch makes a PATCH request
func (hu *HTTPUtils) Patch(url string, headers map[string]string, body interface{}) (*HTTPResponse, error) {
	return hu.makeRequest("PATCH", url, headers, body, nil)
}

// makeRequest makes an HTTP request
func (hu *HTTPUtils) makeRequest(method, url string, headers map[string]string, body interface{}, query map[string]string) (*HTTPResponse, error) {
	// Add query parameters
	if len(query) > 0 {
		url = hu.addQueryParams(url, query)
	}

	// Create request body
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set default content type for POST/PUT/PATCH
	if body != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Make request
	resp, err := hu.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert headers to map
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       respBody,
	}, nil
}

// addQueryParams adds query parameters to URL
func (hu *HTTPUtils) addQueryParams(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// ParseJSONResponse parses JSON response
func (hu *HTTPUtils) ParseJSONResponse(response *HTTPResponse, target interface{}) error {
	if response.Error != nil {
		return response.Error
	}

	return json.Unmarshal(response.Body, target)
}

// IsSuccessStatusCode checks if status code is success
func (hu *HTTPUtils) IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// IsClientErrorStatusCode checks if status code is client error
func (hu *HTTPUtils) IsClientErrorStatusCode(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// IsServerErrorStatusCode checks if status code is server error
func (hu *HTTPUtils) IsServerErrorStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// GetStatusText returns status text for status code
func (hu *HTTPUtils) GetStatusText(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		return "Unknown"
	}
}

// BuildURL builds URL with base and path
func (hu *HTTPUtils) BuildURL(base, path string) string {
	base = strings.TrimSuffix(base, "/")
	path = strings.TrimPrefix(path, "/")
	return base + "/" + path
}

// BuildURLWithParams builds URL with base, path and parameters
func (hu *HTTPUtils) BuildURLWithParams(base, path string, params map[string]string) string {
	url := hu.BuildURL(base, path)
	return hu.addQueryParams(url, params)
}

// ValidateURL validates URL format
func (hu *HTTPUtils) ValidateURL(urlStr string) bool {
	_, err := url.Parse(urlStr)
	return err == nil
}

// ExtractDomain extracts domain from URL
func (hu *HTTPUtils) ExtractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// ExtractPath extracts path from URL
func (hu *HTTPUtils) ExtractPath(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}

// ExtractQuery extracts query parameters from URL
func (hu *HTTPUtils) ExtractQuery(urlStr string) (map[string]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for key, values := range u.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params, nil
}

// SetHeader sets header in request
func (hu *HTTPUtils) SetHeader(req *http.Request, key, value string) {
	req.Header.Set(key, value)
}

// SetHeaders sets multiple headers in request
func (hu *HTTPUtils) SetHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// SetBasicAuth sets basic authentication
func (hu *HTTPUtils) SetBasicAuth(req *http.Request, username, password string) {
	req.SetBasicAuth(username, password)
}

// SetBearerToken sets bearer token
func (hu *HTTPUtils) SetBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// SetAPIKey sets API key
func (hu *HTTPUtils) SetAPIKey(req *http.Request, key, value string) {
	req.Header.Set(key, value)
}

// SetUserAgent sets user agent
func (hu *HTTPUtils) SetUserAgent(req *http.Request, userAgent string) {
	req.Header.Set("User-Agent", userAgent)
}

// SetContentType sets content type
func (hu *HTTPUtils) SetContentType(req *http.Request, contentType string) {
	req.Header.Set("Content-Type", contentType)
}

// SetAccept sets accept header
func (hu *HTTPUtils) SetAccept(req *http.Request, accept string) {
	req.Header.Set("Accept", accept)
}

// SetCustomHeader sets custom header
func (hu *HTTPUtils) SetCustomHeader(req *http.Request, key, value string) {
	req.Header.Set(key, value)
}

// GetHeader gets header value
func (hu *HTTPUtils) GetHeader(resp *http.Response, key string) string {
	return resp.Header.Get(key)
}

// GetHeaders gets all headers
func (hu *HTTPUtils) GetHeaders(resp *http.Response) map[string]string {
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return headers
}

// DownloadFile downloads a file from URL
func (hu *HTTPUtils) DownloadFile(url string, headers map[string]string) ([]byte, error) {
	resp, err := hu.Get(url, headers)
	if err != nil {
		return nil, err
	}

	if !hu.IsSuccessStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("download failed with status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// UploadFile uploads a file to URL
func (hu *HTTPUtils) UploadFile(url string, headers map[string]string, fileData []byte) (*HTTPResponse, error) {
	// Set content type if not provided
	if headers == nil {
		headers = make(map[string]string)
	}
	if headers["Content-Type"] == "" {
		headers["Content-Type"] = "application/octet-stream"
	}

	return hu.makeRequest("POST", url, headers, fileData, nil)
}

// MakeRequestWithRetry makes request with retry logic
func (hu *HTTPUtils) MakeRequestWithRetry(method, url string, headers map[string]string, body interface{}, maxRetries int) (*HTTPResponse, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		resp, err := hu.makeRequest(method, url, headers, body, nil)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Wait before retry
		if i < maxRetries {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return nil, lastErr
}

// MakeRequestWithTimeout makes request with custom timeout
func (hu *HTTPUtils) MakeRequestWithTimeout(method, url string, headers map[string]string, body interface{}, timeout time.Duration) (*HTTPResponse, error) {
	// Create temporary client with custom timeout
	tempClient := &http.Client{Timeout: timeout}

	// Add query parameters
	if query := headers["query"]; query != "" {
		url = hu.addQueryParams(url, map[string]string{"query": query})
		delete(headers, "query")
	}

	// Create request body
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make request
	resp, err := tempClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert headers to map
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       respBody,
	}, nil
}
