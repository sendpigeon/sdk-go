package sendpigeon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL    = "https://api.sendpigeon.dev"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 2
	maxRetries        = 5
)

// ClientOptions configures the HTTP client.
type ClientOptions struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
	Debug      bool
	HTTPClient *http.Client
}

// httpClient handles HTTP requests with retry logic.
type httpClient struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	maxRetries int
	debug      bool
	client     *http.Client
}

func newHTTPClient(apiKey string, opts *ClientOptions) *httpClient {
	baseURL := defaultBaseURL
	timeout := defaultTimeout
	maxRetries := defaultMaxRetries
	debug := false
	var client *http.Client

	if opts != nil {
		if opts.BaseURL != "" {
			baseURL = opts.BaseURL
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.MaxRetries >= 0 {
			maxRetries = opts.MaxRetries
			if maxRetries > 5 {
				maxRetries = 5
			}
		}
		debug = opts.Debug
		client = opts.HTTPClient
	}

	if client == nil {
		client = &http.Client{Timeout: timeout}
	}

	return &httpClient{
		apiKey:     apiKey,
		baseURL:    baseURL,
		timeout:    timeout,
		maxRetries: maxRetries,
		debug:      debug,
		client:     client,
	}
}

// request makes an HTTP request with retry logic.
func (c *httpClient) request(ctx context.Context, method, path string, body interface{}, headers map[string]string) ([]byte, *Error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, NewError(ErrorCodeNetwork, fmt.Sprintf("failed to marshal request body: %v", err))
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	var lastErr *Error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Reset body reader for retry
		if body != nil {
			jsonBody, _ := json.Marshal(body)
			bodyReader = bytes.NewReader(jsonBody)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, NewError(ErrorCodeNetwork, fmt.Sprintf("failed to create request: %v", err))
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "sendpigeon-go/1.0.0")

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return nil, NewError(ErrorCodeTimeout, "request timed out")
			}
			lastErr = NewError(ErrorCodeNetwork, fmt.Sprintf("request failed: %v", err))
			if attempt < c.maxRetries {
				c.sleep(attempt, 0)
				continue
			}
			return nil, lastErr
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = NewError(ErrorCodeNetwork, fmt.Sprintf("failed to read response: %v", err))
			continue
		}

		// Success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		// Parse error response
		var apiErr struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		json.Unmarshal(respBody, &apiErr)

		message := apiErr.Error.Message
		if message == "" {
			message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}

		lastErr = NewAPIError(resp.StatusCode, apiErr.Error.Code, message)

		// Should retry?
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < c.maxRetries {
				retryAfter := c.parseRetryAfter(resp.Header.Get("Retry-After"))
				c.sleep(attempt, retryAfter)
				continue
			}
		}

		return nil, lastErr
	}

	return nil, lastErr
}

// sleep implements exponential backoff.
func (c *httpClient) sleep(attempt int, retryAfter time.Duration) {
	if retryAfter > 0 {
		time.Sleep(retryAfter)
		return
	}
	backoff := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
	if backoff > 10*time.Second {
		backoff = 10 * time.Second
	}
	time.Sleep(backoff)
}

// parseRetryAfter parses the Retry-After header.
func (c *httpClient) parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return 0
}

// Get makes a GET request.
func (c *httpClient) Get(ctx context.Context, path string, headers map[string]string) ([]byte, *Error) {
	return c.request(ctx, http.MethodGet, path, nil, headers)
}

// Post makes a POST request.
func (c *httpClient) Post(ctx context.Context, path string, body interface{}, headers map[string]string) ([]byte, *Error) {
	return c.request(ctx, http.MethodPost, path, body, headers)
}

// Put makes a PUT request.
func (c *httpClient) Put(ctx context.Context, path string, body interface{}, headers map[string]string) ([]byte, *Error) {
	return c.request(ctx, http.MethodPut, path, body, headers)
}

// Patch makes a PATCH request.
func (c *httpClient) Patch(ctx context.Context, path string, body interface{}, headers map[string]string) ([]byte, *Error) {
	return c.request(ctx, http.MethodPatch, path, body, headers)
}

// Delete makes a DELETE request.
func (c *httpClient) Delete(ctx context.Context, path string, headers map[string]string) ([]byte, *Error) {
	return c.request(ctx, http.MethodDelete, path, nil, headers)
}
