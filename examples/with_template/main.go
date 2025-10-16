package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dapoadedire/autosend-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("AUTOSEND_API_KEY")
	if apiKey == "" {
		log.Fatal("AUTOSEND_API_KEY environment variable is required")
	}

	// Create a new client
	client := autosend.NewClient(apiKey)

	// Prepare email request with template
	req := &autosend.SendEmailRequest{
		To: autosend.EmailAddress{
			Email: "customer@example.com",
			Name:  "Jane Smith",
		},
		From: autosend.EmailAddress{
			Email: "hello@mail.yourdomain.com",
			Name:  "Your Company",
		},
		TemplateID: "tmpl_abc123",
		DynamicData: map[string]any{
			"firstName":   "Jane",
			"orderNumber": "ORD-12345",
			"orderTotal":  "$99.99",
		},
		Categories: []string{"order", "confirmation"},
	}

	// Send email with idempotency key for safe retries
	ctx := context.Background()
	idempotencyKey := "order-confirmation-ORD-12345"
	resp, err := client.SendEmailWithIdempotency(ctx, req, idempotencyKey)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Printf("Email sent successfully!\n")
	fmt.Printf("Email ID: %s\n", resp.Data.EmailID)
	fmt.Printf("Status: %s\n", resp.Data.Status)
}
