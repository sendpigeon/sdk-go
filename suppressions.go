package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// SuppressionsService handles suppression operations.
type SuppressionsService struct {
	http *httpClient
}

// List lists suppressed email addresses.
func (s *SuppressionsService) List(ctx context.Context, opts *ListOptions) (*SuppressionListResponse, *Error) {
	path := "/v1/suppressions"
	if opts != nil {
		params := url.Values{}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp SuppressionListResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete removes an email address from the suppression list.
func (s *SuppressionsService) Delete(ctx context.Context, email string) *Error {
	_, err := s.http.Delete(ctx, "/v1/suppressions/"+url.PathEscape(email), nil)
	return err
}
