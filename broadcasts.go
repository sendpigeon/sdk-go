package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// BroadcastsService handles broadcast operations.
type BroadcastsService struct {
	http *httpClient
}

// List lists broadcasts with optional filtering.
func (s *BroadcastsService) List(ctx context.Context, opts *ListBroadcastsOptions) (*ListResponse[Broadcast], *Error) {
	path := "/v1/broadcasts"
	if opts != nil {
		params := url.Values{}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Broadcast]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Create creates a new broadcast.
func (s *BroadcastsService) Create(ctx context.Context, req CreateBroadcastRequest) (*Broadcast, *Error) {
	body, err := s.http.Post(ctx, "/v1/broadcasts", req, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Get retrieves a broadcast by ID.
func (s *BroadcastsService) Get(ctx context.Context, id string) (*Broadcast, *Error) {
	body, err := s.http.Get(ctx, "/v1/broadcasts/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Update updates a broadcast.
func (s *BroadcastsService) Update(ctx context.Context, id string, req UpdateBroadcastRequest) (*Broadcast, *Error) {
	body, err := s.http.Patch(ctx, "/v1/broadcasts/"+id, req, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete removes a broadcast.
func (s *BroadcastsService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/broadcasts/"+id, nil)
	return err
}

// Duplicate creates a copy of a broadcast.
func (s *BroadcastsService) Duplicate(ctx context.Context, id string) (*Broadcast, *Error) {
	body, err := s.http.Post(ctx, "/v1/broadcasts/"+id+"/duplicate", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Recipients lists recipients of a broadcast.
func (s *BroadcastsService) Recipients(ctx context.Context, id string, opts *ListRecipientsOptions) (*ListResponse[BroadcastRecipient], *Error) {
	path := "/v1/broadcasts/" + id + "/recipients"
	if opts != nil {
		params := url.Values{}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[BroadcastRecipient]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Send sends a broadcast immediately with optional targeting.
func (s *BroadcastsService) Send(ctx context.Context, id string, req *SendBroadcastRequest) (*Broadcast, *Error) {
	var reqBody interface{}
	if req != nil {
		reqBody = req
	}
	body, err := s.http.Post(ctx, "/v1/broadcasts/"+id+"/send", reqBody, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Schedule schedules a broadcast for later.
func (s *BroadcastsService) Schedule(ctx context.Context, id string, req ScheduleBroadcastRequest) (*Broadcast, *Error) {
	body, err := s.http.Post(ctx, "/v1/broadcasts/"+id+"/schedule", req, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Cancel cancels a scheduled broadcast.
func (s *BroadcastsService) Cancel(ctx context.Context, id string) (*Broadcast, *Error) {
	body, err := s.http.Post(ctx, "/v1/broadcasts/"+id+"/cancel", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Broadcast
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Test sends a test email for a broadcast.
func (s *BroadcastsService) Test(ctx context.Context, id string, req TestBroadcastRequest) (*TestBroadcastResponse, *Error) {
	body, err := s.http.Post(ctx, "/v1/broadcasts/"+id+"/test", req, nil)
	if err != nil {
		return nil, err
	}

	var resp TestBroadcastResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Analytics retrieves detailed analytics for a broadcast.
func (s *BroadcastsService) Analytics(ctx context.Context, id string) (*BroadcastAnalytics, *Error) {
	body, err := s.http.Get(ctx, "/v1/broadcasts/"+id+"/analytics", nil)
	if err != nil {
		return nil, err
	}

	var resp BroadcastAnalytics
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
