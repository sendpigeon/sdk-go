package sendpigeon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("sk_test_xxx", nil)
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
	if client.Emails == nil {
		t.Error("expected Emails service to be non-nil")
	}
	if client.Templates == nil {
		t.Error("expected Templates service to be non-nil")
	}
	if client.Domains == nil {
		t.Error("expected Domains service to be non-nil")
	}
	if client.APIKeys == nil {
		t.Error("expected APIKeys service to be non-nil")
	}
}

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/emails" {
			t.Errorf("expected /v1/emails, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer sk_test_xxx" {
			t.Errorf("expected Bearer sk_test_xxx, got %s", auth)
		}

		var body SendEmailRequest
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.To) != 1 || body.To[0] != "user@example.com" {
			t.Errorf("unexpected To: %v", body.To)
		}
		if body.Subject != "Hello" {
			t.Errorf("expected subject Hello, got %s", body.Subject)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "email_123",
			"status": "pending",
		})
	}))
	defer server.Close()

	client := New("sk_test_xxx", &ClientOptions{BaseURL: server.URL})
	resp, err := client.Send(context.Background(), SendEmailRequest{
		To:      []string{"user@example.com"},
		Subject: "Hello",
		HTML:    "<p>Hi</p>",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "email_123" {
		t.Errorf("expected id email_123, got %s", resp.ID)
	}
	if resp.Status != "pending" {
		t.Errorf("expected status pending, got %s", resp.Status)
	}
}

func TestSendWithIdempotencyKey(t *testing.T) {
	var receivedKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.Header.Get("Idempotency-Key")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     "email_123",
			"status": "pending",
		})
	}))
	defer server.Close()

	client := New("sk_test_xxx", &ClientOptions{BaseURL: server.URL})
	_, err := client.Send(context.Background(), SendEmailRequest{
		To:             []string{"user@example.com"},
		Subject:        "Hello",
		IdempotencyKey: "unique-key-123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedKey != "unique-key-123" {
		t.Errorf("expected idempotency key unique-key-123, got %s", receivedKey)
	}
}

func TestSendBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/emails/batch" {
			t.Errorf("expected /v1/emails/batch, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"index": 0, "status": "sent", "id": "email_1"},
				{"index": 1, "status": "sent", "id": "email_2"},
			},
			"summary": map[string]interface{}{"sent": 2, "failed": 0},
		})
	}))
	defer server.Close()

	client := New("sk_test_xxx", &ClientOptions{BaseURL: server.URL})
	resp, err := client.SendBatch(context.Background(), []SendEmailRequest{
		{To: []string{"user1@example.com"}, Subject: "Hello 1"},
		{To: []string{"user2@example.com"}, Subject: "Hello 2"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 results, got %d", len(resp.Data))
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "validation_error",
				"message": "Subject is required",
			},
		})
	}))
	defer server.Close()

	client := New("sk_test_xxx", &ClientOptions{BaseURL: server.URL})
	_, err := client.Send(context.Background(), SendEmailRequest{
		To: []string{"user@example.com"},
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if err.Code != ErrorCodeAPI {
		t.Errorf("expected api_error, got %s", err.Code)
	}
	if err.APICode != "validation_error" {
		t.Errorf("expected validation_error, got %s", err.APICode)
	}
	if err.Status != 400 {
		t.Errorf("expected status 400, got %d", err.Status)
	}
}
