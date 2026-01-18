package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// ContactsService handles contact operations.
type ContactsService struct {
	http *httpClient
}

// List lists contacts with optional filtering.
func (s *ContactsService) List(ctx context.Context, opts *ListContactsOptions) (*ListResponse[Contact], *Error) {
	path := "/v1/contacts"
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
		if len(opts.Tags) > 0 {
			params.Set("tags", strings.Join(opts.Tags, ","))
		}
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.Search != "" {
			params.Set("search", opts.Search)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := s.http.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Contact]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Stats returns audience statistics.
func (s *ContactsService) Stats(ctx context.Context) (*AudienceStats, *Error) {
	body, err := s.http.Get(ctx, "/v1/contacts/stats", nil)
	if err != nil {
		return nil, err
	}

	var resp AudienceStats
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Tags returns all unique tags.
func (s *ContactsService) Tags(ctx context.Context) ([]string, *Error) {
	body, err := s.http.Get(ctx, "/v1/contacts/tags", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []string `json:"data"`
	}
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return resp.Data, nil
}

// Create creates a new contact.
func (s *ContactsService) Create(ctx context.Context, req CreateContactRequest) (*Contact, *Error) {
	body, err := s.http.Post(ctx, "/v1/contacts", req, nil)
	if err != nil {
		return nil, err
	}

	var resp Contact
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Batch creates or updates multiple contacts.
func (s *ContactsService) Batch(ctx context.Context, contacts []BatchContactInput) (*BatchContactResponse, *Error) {
	body, err := s.http.Post(ctx, "/v1/contacts/batch", map[string]interface{}{"contacts": contacts}, nil)
	if err != nil {
		return nil, err
	}

	var resp BatchContactResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Get retrieves a contact by ID.
func (s *ContactsService) Get(ctx context.Context, id string) (*Contact, *Error) {
	body, err := s.http.Get(ctx, "/v1/contacts/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Contact
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Update updates a contact.
func (s *ContactsService) Update(ctx context.Context, id string, req UpdateContactRequest) (*Contact, *Error) {
	body, err := s.http.Patch(ctx, "/v1/contacts/"+id, req, nil)
	if err != nil {
		return nil, err
	}

	var resp Contact
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete removes a contact.
func (s *ContactsService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/contacts/"+id, nil)
	return err
}

// Unsubscribe unsubscribes a contact.
func (s *ContactsService) Unsubscribe(ctx context.Context, id string) (*Contact, *Error) {
	body, err := s.http.Post(ctx, "/v1/contacts/"+id+"/unsubscribe", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Contact
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Resubscribe resubscribes a contact.
func (s *ContactsService) Resubscribe(ctx context.Context, id string) (*Contact, *Error) {
	body, err := s.http.Post(ctx, "/v1/contacts/"+id+"/resubscribe", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Contact
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
