// Json i/o  execute & decode http requests, responses, streams
package json

import (
	"fmt"
	"io"
	"net/http"
)

// Execute HTTP GET requests, with X-API-KEY headers as needed, and return the response body as bytes.
func (j *JSON) In(url, apiKey string) ([]byte, error) {
	// Create a new HTTP GET request with the provided context
	req, err := http.NewRequestWithContext(j.CTX, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	// Add API key header if apiKey is provided
	if apiKey != "" {
		req.Header.Set("X-API-KEY", apiKey)
	}

	// Execute the HTTP request
	resp, err := j.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// Return nil if the response body is empty
	if resp.Body == nil || resp.ContentLength == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	// Read and return the response body as bytes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

// Simplify processes a slice of objects ([]any), flattens each object, and removes empty fields.
func (j *JSON) Simplify(input []any, prefix string) []any {
	var result []any

	for _, item := range input {
		obj, ok := item.(map[string]any)
		if !ok {
			continue // Skip non-map items
		}

		flatMap := make(map[string]any)
		var flatten func(map[string]any, string)
		flatten = func(m map[string]any, pfx string) {
			for key, value := range m {
				newKey := key
				if pfx != "" {
					newKey = pfx + "." + key
				}
				switch v := value.(type) {
				case map[string]any:
					flatten(v, newKey) // Recursively flatten nested maps
				case []any:
					for i, item := range v {
						arrayKey := fmt.Sprintf("%s[%d]", newKey, i)
						if nestedMap, ok := item.(map[string]any); ok {
							flatten(nestedMap, arrayKey) // Flatten nested maps in arrays
						} else if !isEmpty(item) {
							flatMap[arrayKey] = item // Add non-empty array items
						}
					}
				default:
					if !isEmpty(v) {
						flatMap[newKey] = v // Add non-empty values
					}
				}
			}
		}
		flatten(obj, prefix)
		result = append(result, flatMap)
	}
	return result
}

// isEmpty checks if a value is an empty string, empty slice, or empty map.
func isEmpty(v any) bool {
	switch val := v.(type) {
	case string:
		return val == ""
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return v == nil
	}
}
