package sendpigeon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// EmailsService handles email operations.
type EmailsService struct {
	http *httpClient
}

// Get retrieves an email by ID.
func (s *EmailsService) Get(ctx context.Context, id string) (*EmailDetail, *Error) {
	body, err := s.http.Get(ctx, "/v1/emails/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp EmailDetail
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// ListResponse represents a paginated list response.
type ListResponse[T any] struct {
	Data   []T `json:"data"`
	Cursor struct {
		Next string `json:"next,omitempty"`
	} `json:"cursor,omitempty"`
}

// List lists emails with optional filters.
func (s *EmailsService) List(ctx context.Context, opts *ListEmailsOptions) (*ListResponse[EmailDetail], *Error) {
	path := "/v1/emails"
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
			params.Set("status", string(opts.Status))
		}
		if opts.Tag != "" {
			params.Set("tag", opts.Tag)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[EmailDetail]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Cancel cancels a scheduled email.
func (s *EmailsService) Cancel(ctx context.Context, id string) (*EmailDetail, *Error) {
	body, err := s.http.Post(ctx, fmt.Sprintf("/v1/emails/%s/cancel", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var resp EmailDetail
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
