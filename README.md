# SendPigeon Go SDK

Official Go SDK for [SendPigeon](https://sendpigeon.dev) - Transactional Email API.

## Installation

```bash
go get github.com/sendpigeon/sdks/go/sendpigeon
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sendpigeon/sdks/go/sendpigeon"
)

func main() {
    client := sendpigeon.New("sk_live_xxx", nil)

    resp, err := client.Send(context.Background(), sendpigeon.SendEmailRequest{
        To:      []string{"user@example.com"},
        Subject: "Hello",
        HTML:    "<p>Welcome to SendPigeon!</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Email sent:", resp.ID)
}
```

## Configuration

```go
client := sendpigeon.New("sk_live_xxx", &sendpigeon.ClientOptions{
    BaseURL:    "https://api.sendpigeon.dev",  // Custom base URL
    Timeout:    30 * time.Second,               // Request timeout
    MaxRetries: 2,                              // Retry attempts (max 5)
    Debug:      false,                          // Enable debug logging
})
```

## Local Development

Use the SendPigeon CLI to catch emails locally:

```bash
# Terminal 1: Start local server
npx @sendpigeon-sdk/cli dev

# Terminal 2: Run your app with dev mode
SENDPIGEON_DEV=true go run main.go
```

When `SENDPIGEON_DEV=true`, the SDK routes requests to `localhost:4100` instead of production.

## Sending Emails

### Basic Email

```go
resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:      []string{"user@example.com"},
    From:    "hello@yourdomain.com",
    Subject: "Welcome!",
    HTML:    "<h1>Welcome</h1><p>Thanks for signing up!</p>",
    Text:    "Welcome! Thanks for signing up.",
})
```

### With Template

```go
resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:         []string{"user@example.com"},
    TemplateID: "tmpl_xxx",
    Variables:  map[string]string{
        "name":    "John",
        "company": "Acme Inc",
    },
})
```

### With Attachments

```go
resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:      []string{"user@example.com"},
    Subject: "Your invoice",
    HTML:    "<p>Please find your invoice attached.</p>",
    Attachments: []sendpigeon.Attachment{
        {Filename: "invoice.pdf", Content: base64Content},
    },
})
```

### Scheduled Email

```go
resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:          []string{"user@example.com"},
    Subject:     "Reminder",
    HTML:        "<p>Don't forget about tomorrow's meeting!</p>",
    ScheduledAt: "2024-01-15T10:00:00Z",
})
```

### Batch Send (up to 100)

```go
resp, err := client.SendBatch(ctx, []sendpigeon.SendEmailRequest{
    {To: []string{"user1@example.com"}, Subject: "Hello", HTML: "<p>Hi User 1!</p>"},
    {To: []string{"user2@example.com"}, Subject: "Hello", HTML: "<p>Hi User 2!</p>"},
})

for _, result := range resp.Data {
    if result.Status == "sent" {
        fmt.Printf("Email %d sent: %s\n", result.Index, result.ID)
    }
}
```

### Tracking

Enable open/click tracking per email (opt-in):

```go
opens := true
clicks := true

resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    HTML:    `<p>Check out our <a href="https://example.com">site</a>!</p>`,
    Tracking: &sendpigeon.TrackingOptions{
        Opens:  &opens,
        Clicks: &clicks,
    },
})

// Response may include warnings if tracking is disabled at org level
if len(resp.Warnings) > 0 {
    fmt.Println("Warnings:", resp.Warnings)
}
```

Configure organization defaults in Settings → Tracking.

## Email Management

```go
// Get email by ID
email, err := client.Emails.Get(ctx, "email_xxx")

// Cancel scheduled email
cancelled, err := client.Emails.Cancel(ctx, "email_xxx")
```

## Templates

```go
// Create template
tmpl, err := client.Templates.Create(ctx, sendpigeon.CreateTemplateRequest{
    Name:    "welcome",
    Subject: "Welcome, {{name}}!",
    HTML:    "<h1>Hello {{name}}</h1><p>Welcome to {{company}}!</p>",
})

// Get template
tmpl, err := client.Templates.Get(ctx, "tmpl_xxx")

// List templates
templates, err := client.Templates.List(ctx, nil)

// Update template
tmpl, err := client.Templates.Update(ctx, "tmpl_xxx", sendpigeon.UpdateTemplateRequest{
    Subject: "Updated subject",
})

// Delete template
err := client.Templates.Delete(ctx, "tmpl_xxx")
```

## Domains

```go
// Add domain
domain, err := client.Domains.Create(ctx, "mail.yourdomain.com")

// DNS records are returned for setup
for _, record := range domain.DNSRecords {
    fmt.Printf("%s %s -> %s\n", record.Type, record.Name, record.Value)
}

// Verify domain
result, err := client.Domains.Verify(ctx, "dom_xxx")
if result.Verified {
    fmt.Println("Domain verified!")
}

// List domains
domains, err := client.Domains.List(ctx, nil)

// Delete domain
err := client.Domains.Delete(ctx, "dom_xxx")
```

## API Keys

```go
// Create API key
key, err := client.APIKeys.Create(ctx, sendpigeon.CreateAPIKeyRequest{
    Name:       "Production",
    Mode:       sendpigeon.APIKeyModeLive,
    Permission: sendpigeon.APIKeyPermissionFullAccess,
})

// Save key.Key - only returned once!
fmt.Println("API Key:", key.Key)

// List API keys
keys, err := client.APIKeys.List(ctx, nil)

// Delete API key
err := client.APIKeys.Delete(ctx, "key_xxx")
```

## Webhook Verification

```go
import "github.com/sendpigeon/sdks/go/sendpigeon"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)

    result := sendpigeon.VerifyWebhook(
        body,
        r.Header.Get("X-Webhook-Signature"),
        r.Header.Get("X-Webhook-Timestamp"),
        "whsec_xxx",  // Your webhook secret
        300,          // Max age in seconds
    )

    if !result.Valid {
        http.Error(w, result.Error, http.StatusUnauthorized)
        return
    }

    // Handle webhook event
    eventType := result.Payload["type"].(string)
    switch eventType {
    case "email.delivered":
        // Handle delivery
    case "email.bounced":
        // Handle bounce
    }

    w.WriteHeader(http.StatusOK)
}
```

### Inbound Email Webhooks

```go
result := sendpigeon.VerifyInboundWebhook(
    body,
    r.Header.Get("X-Webhook-Signature"),
    r.Header.Get("X-Webhook-Timestamp"),
    "whsec_inbound_xxx",
    300,
)
```

## Error Handling

All methods return `(*Response, *Error)`. Check the error:

```go
resp, err := client.Send(ctx, request)
if err != nil {
    switch err.Code {
    case sendpigeon.ErrorCodeAPI:
        // API error (validation, auth, etc.)
        fmt.Printf("API Error [%s]: %s (status: %d)\n", err.APICode, err.Message, err.Status)
    case sendpigeon.ErrorCodeNetwork:
        // Network/connection error
        fmt.Printf("Network Error: %s\n", err.Message)
    case sendpigeon.ErrorCodeTimeout:
        // Request timed out
        fmt.Printf("Timeout: %s\n", err.Message)
    }
    return
}

// Success
fmt.Println("Email ID:", resp.ID)
```

## Idempotency

Prevent duplicate sends with idempotency keys:

```go
resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
    To:             []string{"user@example.com"},
    Subject:        "Order confirmation",
    HTML:           "<p>Your order has been confirmed.</p>",
    IdempotencyKey: "order-12345-confirmation",
})
```

## Requirements

- Go 1.21+

## License

MIT
