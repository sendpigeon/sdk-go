package sendpigeon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type SequencesService struct {
	http *httpClient
}

func (s *SequencesService) List(ctx context.Context, opts *ListSequencesOptions) (*ListResponse[Sequence], *Error) {
	path := "/v1/sequences"
	if opts != nil {
		params := url.Values{}
		if opts.Status != "" {
			params.Set("status", opts.Status)
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
	var resp ListResponse[Sequence]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) Create(ctx context.Context, req CreateSequenceRequest) (*SequenceWithSteps, *Error) {
	return s.postSequence(ctx, "/v1/sequences", req)
}

func (s *SequencesService) Get(ctx context.Context, id string) (*SequenceWithSteps, *Error) {
	body, err := s.http.Get(ctx, "/v1/sequences/"+id, nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceWithSteps
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) Update(ctx context.Context, id string, req UpdateSequenceRequest) (*SequenceWithSteps, *Error) {
	body, err := s.http.Patch(ctx, "/v1/sequences/"+id, req, nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceWithSteps
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) Delete(ctx context.Context, id string) *Error {
	_, err := s.http.Delete(ctx, "/v1/sequences/"+id, nil)
	return err
}

func (s *SequencesService) Activate(ctx context.Context, id string) (*SequenceWithSteps, *Error) {
	return s.postSequence(ctx, fmt.Sprintf("/v1/sequences/%s/activate", id), nil)
}

func (s *SequencesService) Pause(ctx context.Context, id string) (*SequenceWithSteps, *Error) {
	return s.postSequence(ctx, fmt.Sprintf("/v1/sequences/%s/pause", id), nil)
}

func (s *SequencesService) Archive(ctx context.Context, id string) (*SequenceWithSteps, *Error) {
	return s.postSequence(ctx, fmt.Sprintf("/v1/sequences/%s/archive", id), nil)
}

func (s *SequencesService) Duplicate(ctx context.Context, id string) (*SequenceWithSteps, *Error) {
	return s.postSequence(ctx, fmt.Sprintf("/v1/sequences/%s/duplicate", id), nil)
}

func (s *SequencesService) Analytics(ctx context.Context, id string) (*SequenceAnalytics, *Error) {
	body, err := s.http.Get(ctx, fmt.Sprintf("/v1/sequences/%s/analytics", id), nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceAnalytics
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) AddStep(ctx context.Context, seqID string, req AddStepRequest) (*SequenceStep, *Error) {
	body, err := s.http.Post(ctx, fmt.Sprintf("/v1/sequences/%s/steps", seqID), req, nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceStep
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) ListSteps(ctx context.Context, seqID string) ([]SequenceStep, *Error) {
	body, err := s.http.Get(ctx, fmt.Sprintf("/v1/sequences/%s/steps", seqID), nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []SequenceStep `json:"data"`
	}
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return resp.Data, nil
}

func (s *SequencesService) DeleteStep(ctx context.Context, seqID, stepID string) *Error {
	_, err := s.http.Delete(ctx, fmt.Sprintf("/v1/sequences/%s/steps/%s", seqID, stepID), nil)
	return err
}

func (s *SequencesService) Enroll(ctx context.Context, seqID string, req EnrollRequest) (*EnrollResult, *Error) {
	body, err := s.http.Post(ctx, fmt.Sprintf("/v1/sequences/%s/enroll", seqID), req, nil)
	if err != nil {
		return nil, err
	}
	var resp EnrollResult
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) ListEnrollments(ctx context.Context, seqID string, opts *ListEnrollmentsOptions) (*ListResponse[SequenceEnrollment], *Error) {
	path := fmt.Sprintf("/v1/sequences/%s/enrollments", seqID)
	if opts != nil {
		params := url.Values{}
		if opts.Status != "" {
			params.Set("status", opts.Status)
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
	var resp ListResponse[SequenceEnrollment]
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) GetEnrollment(ctx context.Context, seqID, enrollmentID string) (*SequenceEnrollment, *Error) {
	body, err := s.http.Get(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments/%s", seqID, enrollmentID), nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceEnrollment
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}

func (s *SequencesService) PauseEnrollment(ctx context.Context, seqID, enrollmentID string) *Error {
	_, err := s.http.Post(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments/%s/pause", seqID, enrollmentID), nil, nil)
	return err
}

func (s *SequencesService) ResumeEnrollment(ctx context.Context, seqID, enrollmentID string) *Error {
	_, err := s.http.Post(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments/%s/resume", seqID, enrollmentID), nil, nil)
	return err
}

func (s *SequencesService) ExitEnrollment(ctx context.Context, seqID, enrollmentID string) *Error {
	_, err := s.http.Post(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments/%s/exit", seqID, enrollmentID), nil, nil)
	return err
}

func (s *SequencesService) postSequence(ctx context.Context, path string, req interface{}) (*SequenceWithSteps, *Error) {
	body, err := s.http.Post(ctx, path, req, nil)
	if err != nil {
		return nil, err
	}
	var resp SequenceWithSteps
	if jsonErr := json.Unmarshal(body, &resp); jsonErr != nil {
		return nil, NewError(ErrorCodeNetwork, "failed to parse response")
	}
	return &resp, nil
}
