package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// APIKeysService handles API key operations.
type APIKeysService struct {
	http *httpClient
}

// Create creates a new API key.
func (s *APIKeysService) Create(ctx context.Context, req CreateAPIKeyRequest) (*APIKeyWithSecret, *Error) {
	body, err := s.http.Post(ctx, "/v1/api-keys", req, nil)
	if err != nil {
		return nil, err
	}

	var resp APIKeyWithSecret
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Get retrieves an API key by ID.
func (s *APIKeysService) Get(ctx context.Context, id string) (*APIKey, *Error) {
	body, err := s.http.Get(ctx, "/v1/api-keys/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp APIKey
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// List lists all API keys.
func (s *APIKeysService) List(ctx context.Context, opts *ListOptions) (*ListResponse[APIKey], *Error) {
	path := "/v1/api-keys"
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
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[APIKey]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete revokes an API key.
func (s *APIKeysService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/api-keys/"+id, nil)
	return err
}
