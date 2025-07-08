package fx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// GetOption represents a function that modifies a GET request
type GetOption func(*http.Request)

// Get executes HTTP GET requests with the given URL and options
func (a *API) Get(endpoint string, options ...GetOption) ([]byte, error) {
	req, err := a.createRequest(endpoint)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(req)
	}
	return a.executeRequest(req)
}

// createRequest creates and configures the base HTTP GET request
func (a *API) createRequest(endpoint string) (*http.Request, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %s: %w", endpoint, err)
	}

	req, err := http.NewRequestWithContext(a.Ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request for URL %s: %w", parsedURL.String(), err)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// executeRequest executes the HTTP request and returns the response body
func (a *API) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := a.HTTP.Do(req)
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

// Out writes single response for http requests, using a function to source data and a locker to synchronize access or an HTTP 500 error when the input function fails or JSON encoding fails.
func (a *API) Out(w http.ResponseWriter, input func() (any, error), locker sync.Locker) {
	locker.Lock()
	data, err := input()
	locker.Unlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Print value as indented JSON to the standard output or logs error when value cannot be marshaled.
func (a *API) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}

// Simplify processes any input, flattens it, and removes empty fields.
func (a *API) Simplify(input any) any {
	result := make(map[string]any)

	// Use a stack to avoid recursion
	type stackItem struct {
		value  any
		prefix string
	}

	stack := []stackItem{{value: input, prefix: ""}}

	for len(stack) > 0 {
		// Pop from stack
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch v := item.value.(type) {
		case map[string]any:
			for key, value := range v {
				newKey := a.buildKey(item.prefix, key)
				stack = append(stack, stackItem{value: value, prefix: newKey})
			}
		case []any:
			for i, value := range v {
				arrayKey := a.buildArrayKey(item.prefix, i)
				stack = append(stack, stackItem{value: value, prefix: arrayKey})
			}
		default:
			if !a.isEmpty(v) {
				result[item.prefix] = v
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return any(result)
}

// buildKey constructs the flattened key with proper prefix handling (empty string only)
func (a *API) buildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// buildArrayKey constructs array keys with proper prefix handling (empty string only)
func (a *API) buildArrayKey(prefix string, index int) string {
	arrayKey := fmt.Sprintf("[%d]", index)
	if prefix == "" {
		return arrayKey
	}
	return prefix + arrayKey
}

// isEmpty checks if a value is considered empty
func (a *API) isEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return false
	}
}
