// Json i/o
// Json in: execute & decode http requests
// Json out: handle & encode http responses, server-side streaming
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// Execute HTTP GET requests, with X-API-KEY headers as needed, decode the response or return an error when the request fails or the response cannot be decoded.
func (j *JSON) In(url string, useAPIKey bool, apiKey string) (any, error) {
	// Create a new HTTP GET request with the provided context
	req, err := http.NewRequestWithContext(j.CTX, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	// Add API key header if required
	if useAPIKey && apiKey != "" {
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

	// Decode the JSON response into a generic `any` type
	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return result, nil
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
