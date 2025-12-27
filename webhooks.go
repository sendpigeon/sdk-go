package sendpigeon

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Webhook event types
const (
	WebhookEventDelivered  = "email.delivered"
	WebhookEventBounced    = "email.bounced"
	WebhookEventComplained = "email.complained"
	WebhookEventOpened     = "email.opened"
	WebhookEventClicked    = "email.clicked"
	WebhookEventTest       = "webhook.test"
)

// WebhookPayloadData represents the typed webhook payload data.
type WebhookPayloadData struct {
	EmailID       string `json:"emailId,omitempty"`
	ToAddress     string `json:"toAddress,omitempty"`
	FromAddress   string `json:"fromAddress,omitempty"`
	Subject       string `json:"subject,omitempty"`
	BounceType    string `json:"bounceType,omitempty"`
	ComplaintType string `json:"complaintType,omitempty"`
	// Present for email.opened events
	OpenedAt string `json:"openedAt,omitempty"`
	// Present for email.clicked events
	ClickedAt string `json:"clickedAt,omitempty"`
	LinkURL   string `json:"linkUrl,omitempty"`
	LinkIndex *int   `json:"linkIndex,omitempty"`
}

// WebhookPayload represents a typed webhook event.
type WebhookPayload struct {
	Event     string             `json:"event"`
	Timestamp string             `json:"timestamp"`
	Data      WebhookPayloadData `json:"data"`
}

// WebhookVerifyResult represents the result of webhook verification.
type WebhookVerifyResult struct {
	Valid   bool
	Payload map[string]interface{}
	Error   string
}

// ParseWebhookPayload parses a raw payload into a typed WebhookPayload.
func ParseWebhookPayload(payload []byte) (*WebhookPayload, error) {
	var p WebhookPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// VerifyWebhook verifies a webhook signature from SendPigeon.
//
// Parameters:
//   - payload: Raw request body (string or []byte)
//   - signature: Value of X-Webhook-Signature header
//   - timestamp: Value of X-Webhook-Timestamp header
//   - secret: Your webhook secret from dashboard
//   - maxAge: Maximum age of webhook in seconds (default: 300 = 5 minutes)
//
// Example:
//
//	result := sendpigeon.VerifyWebhook(
//	    body,
//	    req.Header.Get("X-Webhook-Signature"),
//	    req.Header.Get("X-Webhook-Timestamp"),
//	    "whsec_xxx",
//	    300,
//	)
//	if result.Valid {
//	    handleEvent(result.Payload)
//	}
func VerifyWebhook(payload []byte, signature, timestamp, secret string, maxAge int) WebhookVerifyResult {
	if maxAge <= 0 {
		maxAge = 300
	}

	// Validate timestamp
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return WebhookVerifyResult{Valid: false, Error: "Invalid timestamp"}
	}

	now := time.Now().Unix()
	if abs(now-ts) > int64(maxAge) {
		return WebhookVerifyResult{Valid: false, Error: "Timestamp too old"}
	}

	// Compute expected signature
	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	// Timing-safe comparison
	if subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) != 1 {
		return WebhookVerifyResult{Valid: false, Error: "Invalid signature"}
	}

	// Parse payload
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return WebhookVerifyResult{Valid: false, Error: "Invalid JSON payload"}
	}

	return WebhookVerifyResult{Valid: true, Payload: data}
}

// VerifyInboundWebhook verifies an inbound email webhook signature.
// Same verification logic as regular webhooks.
func VerifyInboundWebhook(payload []byte, signature, timestamp, secret string, maxAge int) WebhookVerifyResult {
	return VerifyWebhook(payload, signature, timestamp, secret, maxAge)
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
