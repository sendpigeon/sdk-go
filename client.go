// Package sendpigeon provides a Go client for the SendPigeon email API.
package sendpigeon

import (
	"context"
	"encoding/json"
)

// Client is the SendPigeon API client.
type Client struct {
	http         *httpClient
	Emails       *EmailsService
	Templates    *TemplatesService
	Domains      *DomainsService
	APIKeys      *APIKeysService
	Suppressions *SuppressionsService
}

// New creates a new SendPigeon client.
//
// Example:
//
//	client := sendpigeon.New("sk_live_xxx", nil)
//	resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello",
//	    HTML:    "<p>Hi there!</p>",
//	})
func New(apiKey string, opts *ClientOptions) *Client {
	http := newHTTPClient(apiKey, opts)
	return &Client{
		http:         http,
		Emails:       &EmailsService{http: http},
		Templates:    &TemplatesService{http: http},
		Domains:      &DomainsService{http: http},
		APIKeys:      &APIKeysService{http: http},
		Suppressions: &SuppressionsService{http: http},
	}
}

// Send sends a single email.
//
// Example:
//
//	resp, err := client.Send(ctx, sendpigeon.SendEmailRequest{
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello",
//	    HTML:    "<p>Hi there!</p>",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Email ID:", resp.ID)
func (c *Client) Send(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, *Error) {
	headers := make(map[string]string)
	if req.IdempotencyKey != "" {
		headers["Idempotency-Key"] = req.IdempotencyKey
	}

	body, err := c.http.Post(ctx, "/v1/emails", req, headers)
	if err != nil {
		return nil, err
	}

	var resp SendEmailResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// SendBatch sends multiple emails in a single request (max 100).
//
// Example:
//
//	resp, err := client.SendBatch(ctx, []sendpigeon.SendEmailRequest{
//	    {To: []string{"user1@example.com"}, Subject: "Hello", HTML: "<p>Hi User 1!</p>"},
//	    {To: []string{"user2@example.com"}, Subject: "Hello", HTML: "<p>Hi User 2!</p>"},
//	})
func (c *Client) SendBatch(ctx context.Context, emails []SendEmailRequest) (*SendBatchResponse, *Error) {
	body, err := c.http.Post(ctx, "/v1/emails/batch", map[string]interface{}{"emails": emails}, nil)
	if err != nil {
		return nil, err
	}

	var resp SendBatchResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
