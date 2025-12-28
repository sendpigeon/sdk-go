package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// TemplatesService handles template operations.
type TemplatesService struct {
	http *httpClient
}

// Create creates a new template.
func (s *TemplatesService) Create(ctx context.Context, req CreateTemplateRequest) (*Template, *Error) {
	body, err := s.http.Post(ctx, "/v1/templates", req, nil)
	if err != nil {
		return nil, err
	}

	var resp Template
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Get retrieves a template by ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*Template, *Error) {
	body, err := s.http.Get(ctx, "/v1/templates/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Template
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// List lists all templates.
func (s *TemplatesService) List(ctx context.Context, opts *ListOptions) (*ListResponse[Template], *Error) {
	path := "/v1/templates"
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

	var resp ListResponse[Template]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Update updates a template.
func (s *TemplatesService) Update(ctx context.Context, id string, req UpdateTemplateRequest) (*Template, *Error) {
	body, err := s.http.Patch(ctx, "/v1/templates/"+id, req, nil)
	if err != nil {
		return nil, err
	}

	var resp Template
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Delete deletes a template.
func (s *TemplatesService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/templates/"+id, nil)
	return err
}

// Publish publishes a template.
func (s *TemplatesService) Publish(ctx context.Context, id string) (*Template, *Error) {
	body, err := s.http.Post(ctx, "/v1/templates/"+id+"/publish", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Template
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Unpublish unpublishes a template.
func (s *TemplatesService) Unpublish(ctx context.Context, id string) (*Template, *Error) {
	body, err := s.http.Post(ctx, "/v1/templates/"+id+"/unpublish", nil, nil)
	if err != nil {
		return nil, err
	}

	var resp Template
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// Test sends a test email using the template.
func (s *TemplatesService) Test(ctx context.Context, id string, req TestTemplateRequest) (*TestTemplateResponse, *Error) {
	body, err := s.http.Post(ctx, "/v1/templates/"+id+"/test", req, nil)
	if err != nil {
		return nil, err
	}

	var resp TestTemplateResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
