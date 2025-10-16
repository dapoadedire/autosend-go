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

	// Prepare email request
	req := &autosend.SendEmailRequest{
		To: autosend.EmailAddress{
			Email: "customer@example.com",
			Name:  "Jane Smith",
		},
		From: autosend.EmailAddress{
			Email: "hello@mail.yourdomain.com",
			Name:  "Your Company",
		},
		Subject: "Welcome to Our Platform!",
		HTML:    "<h1>Welcome, {{name}}!</h1><p>Thanks for signing up.</p>",
		DynamicData: map[string]any{
			"name": "Jane",
		},
		Categories: []string{"welcome", "onboarding"},
	}

	// Send email
	ctx := context.Background()
	resp, err := client.SendEmail(ctx, req)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Printf("Email sent successfully!\n")
	fmt.Printf("Email ID: %s\n", resp.Data.EmailID)
	fmt.Printf("Status: %s\n", resp.Data.Status)
	fmt.Printf("Queued At: %s\n", resp.Data.QueuedAt)
}
