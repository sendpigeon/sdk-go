package sendpigeon

import (
	"context"
	"encoding/json"
)

// TrackingService handles tracking operations.
type TrackingService struct {
	http *httpClient
}

// GetDefaults retrieves organization tracking defaults.
func (s *TrackingService) GetDefaults(ctx context.Context) (*TrackingDefaults, *Error) {
	body, err := s.http.Get(ctx, "/v1/tracking/defaults", nil)
	if err != nil {
		return nil, err
	}

	var resp TrackingDefaults
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}

// UpdateDefaults updates organization tracking defaults.
func (s *TrackingService) UpdateDefaults(ctx context.Context, req UpdateTrackingDefaultsRequest) (*TrackingDefaults, *Error) {
	body, err := s.http.Patch(ctx, "/v1/tracking/defaults", req, nil)
	if err != nil {
		return nil, err
	}

	var resp TrackingDefaults
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}

	return &resp, nil
}
