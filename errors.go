package autosend

import (
	"fmt"
	"strings"
)

// APIError represents an error response from the Autosend API.
type APIError struct {
	StatusCode    int
	Message       string
	Errors        []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
	RetryAfter    int // Seconds to wait before retrying (for 429 errors)
	RateLimitInfo *RateLimitInfo
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		var fields []string
		for _, err := range e.Errors {
			fields = append(fields, fmt.Sprintf("%s: %s", err.Field, err.Message))
		}
		return fmt.Sprintf("autosend API error (status %d): %s - %s",
			e.StatusCode, e.Message, strings.Join(fields, ", "))
	}

	if e.RetryAfter > 0 {
		return fmt.Sprintf("autosend API error (status %d): %s (retry after %d seconds)",
			e.StatusCode, e.Message, e.RetryAfter)
	}

	return fmt.Sprintf("autosend API error (status %d): %s", e.StatusCode, e.Message)
}

// IsRateLimitError returns true if the error is a rate limit error (429).
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == 429
}

// IsValidationError returns true if the error is a validation error (400).
func (e *APIError) IsValidationError() bool {
	return e.StatusCode == 400
}

// IsAuthenticationError returns true if the error is an authentication error (401).
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == 401
}

// IsForbiddenError returns true if the error is a forbidden error (403).
func (e *APIError) IsForbiddenError() bool {
	return e.StatusCode == 403
}

// IsServerError returns true if the error is a server error (5xx).
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// GetRetryAfter returns the number of seconds to wait before retrying.
// Returns 0 if not a rate limit error or if RetryAfter is not set.
func (e *APIError) GetRetryAfter() int {
	if e.IsRateLimitError() {
		return e.RetryAfter
	}
	return 0
}
