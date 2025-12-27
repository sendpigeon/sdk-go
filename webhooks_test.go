package sendpigeon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func sign(payload, secret string, timestamp int64) string {
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookValid(t *testing.T) {
	secret := "whsec_test123"
	payload := `{"type":"email.delivered","data":{"id":"email_123"}}`
	timestamp := time.Now().Unix()
	signature := sign(payload, secret, timestamp)

	result := VerifyWebhook(
		[]byte(payload),
		signature,
		strconv.FormatInt(timestamp, 10),
		secret,
		300,
	)

	if !result.Valid {
		t.Errorf("expected valid, got error: %s", result.Error)
	}
	if result.Payload["type"] != "email.delivered" {
		t.Errorf("unexpected payload type: %v", result.Payload["type"])
	}
}

func TestVerifyWebhookInvalidSignature(t *testing.T) {
	secret := "whsec_test123"
	payload := `{"type":"email.delivered"}`
	timestamp := time.Now().Unix()

	result := VerifyWebhook(
		[]byte(payload),
		"invalid_signature",
		strconv.FormatInt(timestamp, 10),
		secret,
		300,
	)

	if result.Valid {
		t.Error("expected invalid")
	}
	if result.Error != "Invalid signature" {
		t.Errorf("expected Invalid signature, got %s", result.Error)
	}
}

func TestVerifyWebhookExpiredTimestamp(t *testing.T) {
	secret := "whsec_test123"
	payload := `{"type":"email.delivered"}`
	timestamp := time.Now().Add(-10 * time.Minute).Unix() // 10 minutes ago
	signature := sign(payload, secret, timestamp)

	result := VerifyWebhook(
		[]byte(payload),
		signature,
		strconv.FormatInt(timestamp, 10),
		secret,
		300, // 5 minutes max age
	)

	if result.Valid {
		t.Error("expected invalid due to old timestamp")
	}
	if result.Error != "Timestamp too old" {
		t.Errorf("expected Timestamp too old, got %s", result.Error)
	}
}

func TestVerifyWebhookInvalidTimestamp(t *testing.T) {
	secret := "whsec_test123"
	payload := `{"type":"email.delivered"}`

	result := VerifyWebhook(
		[]byte(payload),
		"signature",
		"not-a-number",
		secret,
		300,
	)

	if result.Valid {
		t.Error("expected invalid")
	}
	if result.Error != "Invalid timestamp" {
		t.Errorf("expected Invalid timestamp, got %s", result.Error)
	}
}

func TestVerifyWebhookInvalidJSON(t *testing.T) {
	secret := "whsec_test123"
	payload := `not json`
	timestamp := time.Now().Unix()
	signature := sign(payload, secret, timestamp)

	result := VerifyWebhook(
		[]byte(payload),
		signature,
		strconv.FormatInt(timestamp, 10),
		secret,
		300,
	)

	if result.Valid {
		t.Error("expected invalid")
	}
	if result.Error != "Invalid JSON payload" {
		t.Errorf("expected Invalid JSON payload, got %s", result.Error)
	}
}

func TestVerifyInboundWebhook(t *testing.T) {
	secret := "whsec_inbound123"
	payload := `{"type":"email.received","data":{"from":"sender@example.com"}}`
	timestamp := time.Now().Unix()
	signature := sign(payload, secret, timestamp)

	result := VerifyInboundWebhook(
		[]byte(payload),
		signature,
		strconv.FormatInt(timestamp, 10),
		secret,
		300,
	)

	if !result.Valid {
		t.Errorf("expected valid, got error: %s", result.Error)
	}
}
