package sendpigeon

import (
	"context"
	"encoding/json"
	"fmt"
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
