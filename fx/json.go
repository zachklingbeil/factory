// Json i/o  execute & decode http requests, responses, streams
package fx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type JSON struct {
	HTTP *http.Client
	CTX  context.Context
}

// Json initializes and returns a new JSON utility instance, using http.Client and context.Context created in NewFactory.
func Json(http http.Client, ctx context.Context) *JSON {
	return &JSON{
		HTTP: &http,
		CTX:  context.Background(),
	}
}

// Print value as indented JSON to the standard output or logs error when value cannot be marshaled.
func (j *JSON) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}

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

// Out writes single response for http requests, using a function to source data and a locker to synchronize access or an HTTP 500 error when the input function fails or JSON encoding fails.
func (j *JSON) Out(w http.ResponseWriter, input func() (any, error), locker sync.Locker) {
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

// OutSSE is Out at a defined interval, streams responses until the client disconnects or the context is canceled.
func (j *JSON) OutSSE(w http.ResponseWriter, r *http.Request, input func() (any, error), interval time.Duration) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx := r.Context()
	var buf bytes.Buffer

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			data, err := input()
			if err != nil {
				return
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				return
			}

			buf.Reset()
			buf.WriteString("data: ")
			buf.Write(jsonData)
			buf.WriteString("\n\n")

			w.Write(buf.Bytes())
			flusher.Flush()
		}
	}
}
