package json

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetOption represents a function that modifies a GET request
type GetOption func(*http.Request)

// Get executes HTTP GET requests with the given URL and options
func (j *JSON) Get(endpoint string, options ...GetOption) ([]byte, error) {
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
func (j *JSON) createRequest(endpoint string) (*http.Request, error) {
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
func (j *JSON) executeRequest(req *http.Request) ([]byte, error) {
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

	// Read response body with size limit for safety
	const maxResponseSize = 50 * 1024 * 1024 // 50MB limit
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return body, nil
}

// Header option builders
func WithHeader(key, value string) GetOption {
	return func(req *http.Request) {
		if value != "" {
			req.Header.Set(key, value)
		}
	}
}

func WithAPIKey(apiKey string) GetOption {
	return WithHeader("X-API-Key", apiKey)
}

func WithBearerToken(token string) GetOption {
	return func(req *http.Request) {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func WithBasicAuth(username, password string) GetOption {
	return func(req *http.Request) {
		if username != "" || password != "" {
			req.SetBasicAuth(username, password)
		}
	}
}

// Query parameter option builders
func WithQuery(key string, values ...string) GetOption {
	return func(req *http.Request) {
		if len(values) == 0 {
			return
		}

		query := req.URL.Query()
		if len(values) == 1 {
			query.Set(key, values[0])
		} else {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		req.URL.RawQuery = query.Encode()
	}
}

func WithQueryMap(params map[string]string) GetOption {
	return func(req *http.Request) {
		if len(params) == 0 {
			return
		}

		query := req.URL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}
}
