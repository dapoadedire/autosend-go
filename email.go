package autosend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// SendEmail sends an email using the Autosend API.
// It returns the response data or an error if the request fails.
func (c *Client) SendEmail(ctx context.Context, req *SendEmailRequest) (*SendEmailResponse, error) {
	return c.SendEmailWithIdempotency(ctx, req, "")
}

// SendEmailWithIdempotency sends an email with an idempotency key.
// The idempotency key allows you to safely retry requests without sending duplicate emails.
// If you retry a request with the same idempotency key within 24 hours,
// you'll receive the same response without sending a duplicate email.
func (c *Client) SendEmailWithIdempotency(ctx context.Context, req *SendEmailRequest, idempotencyKey string) (*SendEmailResponse, error) {
	if err := validateSendEmailRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", "/mails/send", req, idempotencyKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var emailResp SendEmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &emailResp, nil
}

// validateSendEmailRequest validates the required fields in the send email request.
func validateSendEmailRequest(req *SendEmailRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.To.Email == "" {
		return fmt.Errorf("to.email is required")
	}

	if req.From.Email == "" {
		return fmt.Errorf("from.email is required")
	}

	// Either HTML/Text or TemplateID must be provided
	if req.TemplateID == "" && req.HTML == "" && req.Text == "" {
		return fmt.Errorf("either templateId or html/text content must be provided")
	}

	// If not using a template, subject is required
	if req.TemplateID == "" && req.Subject == "" {
		return fmt.Errorf("subject is required when not using a template")
	}

	return nil
}
