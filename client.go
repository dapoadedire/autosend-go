package autosend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Autosend API.
	DefaultBaseURL = "https://api.autosend.com/v1"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is the Autosend API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Config holds the configuration for creating a new Client.
type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// NewClient creates a new Autosend API client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithConfig creates a new Autosend API client with custom configuration.
func NewClientWithConfig(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}

	if config.HTTPClient == nil {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = DefaultTimeout
		}
		config.HTTPClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}
}

// doRequest performs an HTTP request and handles the response.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, idempotencyKey string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "autosend-go/1.0.0")

	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseRateLimitHeaders extracts rate limit information from response headers.
func parseRateLimitHeaders(headers http.Header) *RateLimitInfo {
	info := &RateLimitInfo{}

	if limit := headers.Get("X-RateLimit-Limit"); limit != "" {
		info.Limit, _ = strconv.Atoi(limit)
	}

	if remaining := headers.Get("X-RateLimit-Remaining"); remaining != "" {
		info.Remaining, _ = strconv.Atoi(remaining)
	}

	if reset := headers.Get("X-RateLimit-Reset"); reset != "" {
		info.Reset, _ = strconv.ParseInt(reset, 10, 64)
	}

	return info
}

// handleErrorResponse parses and returns an appropriate error for non-2xx responses.
func handleErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %w", err)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// If we can't parse the error response, return a generic error
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	rateLimitInfo := parseRateLimitHeaders(resp.Header)

	return &APIError{
		StatusCode:    resp.StatusCode,
		Message:       errResp.Message,
		Errors:        errResp.Errors,
		RetryAfter:    errResp.RetryAfter,
		RateLimitInfo: rateLimitInfo,
	}
}
