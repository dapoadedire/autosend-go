package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dapoadedire/autosend-go"
)

func main() {
	apiKey := os.Getenv("AUTOSEND_API_KEY")
	if apiKey == "" {
		log.Fatal("AUTOSEND_API_KEY environment variable is required")
	}

	client := autosend.NewClient(apiKey)

	req := &autosend.SendEmailRequest{
		To: autosend.EmailAddress{
			Email: "customer@example.com",
			Name:  "Jane Smith",
		},
		From: autosend.EmailAddress{
			Email: "hello@mail.yourdomain.com",
			Name:  "Your Company",
		},
		Subject: "Test Email",
		HTML:    "<p>This is a test email</p>",
	}

	// Send email with error handling
	ctx := context.Background()
	resp, err := sendEmailWithRetry(ctx, client, req, 3)
	if err != nil {
		log.Fatalf("Failed to send email after retries: %v", err)
	}

	fmt.Printf("Email sent successfully! ID: %s\n", resp.Data.EmailID)
}

// sendEmailWithRetry implements retry logic with exponential backoff
func sendEmailWithRetry(ctx context.Context, client *autosend.Client, req *autosend.SendEmailRequest, maxRetries int) (*autosend.SendEmailResponse, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := client.SendEmail(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Check if it's an API error
		var apiErr *autosend.APIError
		if errors.As(err, &apiErr) {
			// Handle specific error types
			switch {
			case apiErr.IsRateLimitError():
				// Wait for the specified time before retrying
				waitTime := time.Duration(apiErr.GetRetryAfter()) * time.Second
				if waitTime == 0 {
					// If no retry-after header, use exponential backoff
					waitTime = time.Duration(1<<attempt) * time.Second
				}
				fmt.Printf("Rate limited. Waiting %v before retry %d/%d...\n", waitTime, attempt+1, maxRetries)
				time.Sleep(waitTime)
				continue

			case apiErr.IsValidationError():
				// Validation errors won't be fixed by retrying
				fmt.Printf("Validation error: %v\n", apiErr)
				return nil, apiErr

			case apiErr.IsAuthenticationError():
				// Authentication errors won't be fixed by retrying
				fmt.Printf("Authentication error: %v\n", apiErr)
				return nil, apiErr

			case apiErr.IsServerError():
				// Server errors might be temporary, retry with exponential backoff
				waitTime := time.Duration(1<<attempt) * time.Second
				fmt.Printf("Server error. Waiting %v before retry %d/%d...\n", waitTime, attempt+1, maxRetries)
				time.Sleep(waitTime)
				continue

			default:
				// Other API errors
				fmt.Printf("API error: %v\n", apiErr)
				return nil, apiErr
			}
		}

		// For non-API errors (network issues, etc.), retry with exponential backoff
		if attempt < maxRetries-1 {
			waitTime := time.Duration(1<<attempt) * time.Second
			fmt.Printf("Error: %v. Waiting %v before retry %d/%d...\n", err, waitTime, attempt+1, maxRetries)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
