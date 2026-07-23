package sendpigeon

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type EventsService struct {
	http *httpClient
}

func (s *EventsService) Send(ctx context.Context, req FireEventRequest) (*FireEventResponse, *Error) {
	body, err := s.http.Post(ctx, "/v1/events", req, nil)
	if err != nil {
		return nil, err
	}
	var resp FireEventResponse
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *EventsService) List(ctx context.Context, opts *ListEventsOptions) (*ListResponse[Event], *Error) {
	path := "/v1/events"
	if opts != nil {
		params := url.Values{}
		if opts.Name != "" {
			params.Set("name", opts.Name)
		}
		if opts.ContactID != "" {
			params.Set("contactId", opts.ContactID)
		}
		if opts.Since != "" {
			params.Set("since", opts.Since)
		}
		if opts.Until != "" {
			params.Set("until", opts.Until)
		}
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
	var resp ListResponse[Event]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *EventsService) ContactEvents(ctx context.Context, contactID string, opts *ListEventsOptions) (*ListResponse[Event], *Error) {
	path := "/v1/events/contacts/" + contactID + "/events"
	if opts != nil {
		params := url.Values{}
		if opts.Name != "" {
			params.Set("name", opts.Name)
		}
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
	var resp ListResponse[Event]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}
