package autosend

import "time"

// EmailAddress represents an email address with an optional name.
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// SendEmailRequest represents the request body for sending an email.
type SendEmailRequest struct {
	// Required fields
	To   EmailAddress `json:"to"`
	From EmailAddress `json:"from"`

	// Content fields (either HTML/Text or TemplateID required)
	Subject    string `json:"subject,omitempty"`
	HTML       string `json:"html,omitempty"`
	Text       string `json:"text,omitempty"`
	TemplateID string `json:"templateId,omitempty"`

	// Optional fields
	ReplyTo             *EmailAddress      `json:"replyTo,omitempty"`
	UnsubscribeGroupID  string             `json:"unsubscribeGroupId,omitempty"`
	Categories          []string           `json:"categories,omitempty"`
	DynamicData         map[string]any     `json:"dynamicData,omitempty"`
	ScheduledAt         string             `json:"scheduledAt,omitempty"`
	CampaignName        string             `json:"campaignName,omitempty"`
	Test                bool               `json:"test,omitempty"`
}

// SendEmailResponse represents the successful response from the send email API.
type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		EmailID  string    `json:"emailId"`
		Status   string    `json:"status"`
		QueuedAt time.Time `json:"queuedAt"`
	} `json:"data"`
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
	RetryAfter int `json:"retryAfter,omitempty"` // For 429 responses
}

// RateLimitInfo contains rate limit information from response headers.
type RateLimitInfo struct {
	Limit     int   // X-RateLimit-Limit
	Remaining int   // X-RateLimit-Remaining
	Reset     int64 // X-RateLimit-Reset (Unix timestamp)
}
