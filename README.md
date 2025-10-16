# Autosend Go SDK

Go SDK for the [Autosend](https://autosend.com) email API. Send transactional and marketing emails with ease.

## Features

- Simple and intuitive API
- Full support for the Autosend email API
- Idempotency support for safe retries
- Comprehensive error handling
- Rate limit information in error responses
- Template support with dynamic data
- Context-aware requests
- Configurable HTTP client

## Installation

```bash
go get github.com/dapoadedire/autosend-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/dapoadedire/autosend-go"
)

func main() {
    // Create a new client
    client := autosend.NewClient(os.Getenv("AUTOSEND_API_KEY"))

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
        Subject: "Welcome!",
        HTML:    "<h1>Welcome to our platform!</h1>",
    }

    // Send email
    resp, err := client.SendEmail(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Email sent! ID: %s\n", resp.Data.EmailID)
}
```

## Usage

### Creating a Client

**Basic client:**
```go
client := autosend.NewClient("your-api-key")
```

**With custom configuration:**
```go
client := autosend.NewClientWithConfig(autosend.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://api.autosend.com/v1",
    Timeout: 30 * time.Second,
})
```

**Using functional options:**
```go
client := autosend.NewClientWithOptions(
    "your-api-key",
    autosend.WithTimeout(60 * time.Second),
    autosend.WithBaseURL("https://api.autosend.com/v1"),
)
```

### Sending Emails

**Basic email:**
```go
req := &autosend.SendEmailRequest{
    To: autosend.EmailAddress{
        Email: "recipient@example.com",
        Name:  "John Doe",
    },
    From: autosend.EmailAddress{
        Email: "sender@yourdomain.com",
        Name:  "Your Company",
    },
    Subject: "Hello!",
    HTML:    "<p>This is a test email</p>",
    Text:    "This is a test email",
}

resp, err := client.SendEmail(context.Background(), req)
```

**With template:**
```go
req := &autosend.SendEmailRequest{
    To: autosend.EmailAddress{
        Email: "customer@example.com",
    },
    From: autosend.EmailAddress{
        Email: "hello@mail.yourdomain.com",
    },
    TemplateID: "tmpl_abc123",
    DynamicData: map[string]any{
        "firstName": "Jane",
        "orderNumber": "ORD-12345",
    },
}

resp, err := client.SendEmail(context.Background(), req)
```

**With idempotency key:**
```go
idempotencyKey := "unique-request-id-12345"
resp, err := client.SendEmailWithIdempotency(context.Background(), req, idempotencyKey)
```

**With all options:**
```go
req := &autosend.SendEmailRequest{
    To: autosend.EmailAddress{
        Email: "customer@example.com",
        Name:  "Jane Smith",
    },
    From: autosend.EmailAddress{
        Email: "hello@mail.yourdomain.com",
        Name:  "Your Company",
    },
    ReplyTo: &autosend.EmailAddress{
        Email: "support@yourdomain.com",
        Name:  "Support Team",
    },
    Subject: "Welcome!",
    HTML:    "<h1>Welcome, {{name}}!</h1>",
    DynamicData: map[string]any{
        "name": "Jane",
    },
    Categories: []string{"welcome", "onboarding"},
    CampaignName: "Welcome Campaign",
    ScheduledAt: "2025-10-20T10:00:00Z",
    Test: false,
}

resp, err := client.SendEmail(context.Background(), req)
```

### Error Handling

The SDK provides comprehensive error handling with typed errors:

```go
resp, err := client.SendEmail(ctx, req)
if err != nil {
    var apiErr *autosend.APIError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsRateLimitError():
            fmt.Printf("Rate limited. Retry after %d seconds\n", apiErr.GetRetryAfter())

        case apiErr.IsValidationError():
            fmt.Printf("Validation error: %v\n", apiErr)

        case apiErr.IsAuthenticationError():
            fmt.Println("Invalid API key")

        case apiErr.IsServerError():
            fmt.Println("Server error, try again later")
        }
    }
    return err
}
```

### Retry Logic with Exponential Backoff

```go
func sendWithRetry(client *autosend.Client, req *autosend.SendEmailRequest, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := client.SendEmail(context.Background(), req)
        if err == nil {
            return nil
        }

        var apiErr *autosend.APIError
        if errors.As(err, &apiErr) {
            if apiErr.IsRateLimitError() {
                waitTime := time.Duration(apiErr.GetRetryAfter()) * time.Second
                if waitTime == 0 {
                    waitTime = time.Duration(1<<attempt) * time.Second
                }
                time.Sleep(waitTime)
                continue
            }
            // Don't retry validation or auth errors
            if apiErr.IsValidationError() || apiErr.IsAuthenticationError() {
                return err
            }
        }

        // Exponential backoff for other errors
        if attempt < maxRetries-1 {
            time.Sleep(time.Duration(1<<attempt) * time.Second)
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

## API Reference

### Types

#### `EmailAddress`
```go
type EmailAddress struct {
    Email string `json:"email"`
    Name  string `json:"name,omitempty"`
}
```

#### `SendEmailRequest`
```go
type SendEmailRequest struct {
    To                 EmailAddress       `json:"to"`                   // Required
    From               EmailAddress       `json:"from"`                 // Required
    Subject            string             `json:"subject,omitempty"`    // Required (unless using template)
    HTML               string             `json:"html,omitempty"`       // Required (unless using template)
    Text               string             `json:"text,omitempty"`
    TemplateID         string             `json:"templateId,omitempty"`
    ReplyTo            *EmailAddress      `json:"replyTo,omitempty"`
    UnsubscribeGroupID string             `json:"unsubscribeGroupId,omitempty"`
    Categories         []string           `json:"categories,omitempty"`
    DynamicData        map[string]any     `json:"dynamicData,omitempty"`
    ScheduledAt        string             `json:"scheduledAt,omitempty"`
    CampaignName       string             `json:"campaignName,omitempty"`
    Test               bool               `json:"test,omitempty"`
}
```

#### `SendEmailResponse`
```go
type SendEmailResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        EmailID  string    `json:"emailId"`
        Status   string    `json:"status"`
        QueuedAt time.Time `json:"queuedAt"`
    } `json:"data"`
}
```

### Methods

#### `NewClient(apiKey string) *Client`
Creates a new client with default configuration.

#### `NewClientWithConfig(config Config) *Client`
Creates a new client with custom configuration.

#### `NewClientWithOptions(apiKey string, opts ...ClientOption) *Client`
Creates a new client with functional options.

#### `SendEmail(ctx context.Context, req *SendEmailRequest) (*SendEmailResponse, error)`
Sends an email.

#### `SendEmailWithIdempotency(ctx context.Context, req *SendEmailRequest, idempotencyKey string) (*SendEmailResponse, error)`
Sends an email with an idempotency key for safe retries.

### Error Methods

#### `IsRateLimitError() bool`
Returns true if the error is a rate limit error (429).

#### `IsValidationError() bool`
Returns true if the error is a validation error (400).

#### `IsAuthenticationError() bool`
Returns true if the error is an authentication error (401).

#### `IsForbiddenError() bool`
Returns true if the error is a forbidden error (403).

#### `IsServerError() bool`
Returns true if the error is a server error (5xx).

#### `GetRetryAfter() int`
Returns the number of seconds to wait before retrying for rate limit errors.

## Examples

See the [examples](./examples) directory for complete working examples:

- [Basic usage](./examples/basic/main.go)
- [Using templates](./examples/with_template/main.go)
- [Error handling and retries](./examples/error_handling/main.go)

## Best Practices

1. **Use environment variables for API keys:**
   ```go
   client := autosend.NewClient(os.Getenv("AUTOSEND_API_KEY"))
   ```

2. **Always use context for cancellation and timeouts:**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   resp, err := client.SendEmail(ctx, req)
   ```

3. **Implement retry logic for transient failures:**
   ```go
   // See examples/error_handling/main.go for a complete implementation
   ```

4. **Use idempotency keys for critical emails:**
   ```go
   idempotencyKey := fmt.Sprintf("order-%s", orderID)
   resp, err := client.SendEmailWithIdempotency(ctx, req, idempotencyKey)
   ```

5. **Validate email addresses before sending:**
   ```go
   if !isValidEmail(email) {
       return errors.New("invalid email address")
   }
   ```

6. **Use verified domains for the From address**

## Rate Limits

The Autosend API has the following rate limits:
- 2 requests per second per API key
- 50 requests per minute per API key

Rate limit information is included in the `APIError` type when you receive a 429 response.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Support

For issues with this SDK, please [open an issue](https://github.com/dapoadedire/autosend-go/issues).

For API-related questions, contact [Autosend support](https://autosend.com/support).
