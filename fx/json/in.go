package json

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Json struct {
	HTTP *http.Client
	Ctx  context.Context
}

// NewJson creates a new Json instance with a default HTTP client and context
func NewJson(ctx context.Context) *Json {
	return &Json{
		HTTP: &http.Client{},
		Ctx:  ctx,
	}
}

// GetOption represents a function that modifies a GET request
type GetOption func(*http.Request)

// Get executes HTTP GET requests with the given URL and options
func (j *Json) Get(endpoint string, options ...GetOption) ([]byte, error) {
	req, err := j.createRequest(endpoint)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(req)
	}
	return j.executeRequest(req)
}

// createRequest creates and configures the base HTTP GET request
func (j *Json) createRequest(endpoint string) (*http.Request, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %s: %w", endpoint, err)
	}

	req, err := http.NewRequestWithContext(j.Ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request for URL %s: %w", parsedURL.String(), err)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// executeRequest executes the HTTP request and returns the response body
func (j *Json) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := j.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful status codes (2xx range)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Handle empty response body
	if resp.ContentLength == 0 {
		return []byte{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return body, nil
}
