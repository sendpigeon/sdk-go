package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// DomainsService handles domain operations.
type DomainsService struct {
	http *httpClient
}

// Create adds a new domain.
func (s *DomainsService) Create(ctx context.Context, name string) (*DomainWithDNSRecords, *Error) {
	body, err := s.http.Post(ctx, "/v1/domains", map[string]string{"name": name}, nil)
	if err != nil {
		return nil, err
	}

	var resp DomainWithDNSRecords
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Get retrieves a domain by ID.
func (s *DomainsService) Get(ctx context.Context, id string) (*DomainWithDNSRecords, *Error) {
	body, err := s.http.Get(ctx, "/v1/domains/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp DomainWithDNSRecords
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// List lists all domains.
func (s *DomainsService) List(ctx context.Context, opts *ListOptions) (*ListResponse[Domain], *Error) {
	path := "/v1/domains"
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

	var resp ListResponse[Domain]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Verify triggers domain verification.
func (s *DomainsService) Verify(ctx context.Context, id string) (*DomainVerificationResult, *Error) {
	body, err := s.http.Post(ctx, "/v1/domains/"+id+"/verify", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp DomainVerificationResult
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete removes a domain.
func (s *DomainsService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/domains/"+id, nil)
	return err
}
