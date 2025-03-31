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

func Json(http http.Client, ctx context.Context) *JSON {
	return &JSON{
		HTTP: &http,
		CTX:  context.Background(),
	}
}

func (j *JSON) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}

func (j *JSON) In(url string, result any, useAPIKey bool, apiKey string) error {
	req, err := http.NewRequestWithContext(j.CTX, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	if useAPIKey && apiKey != "" {
		req.Header.Set("X-API-KEY", apiKey)
	}

	resp, err := j.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get JSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get JSON: status code %d", resp.StatusCode)
	}

	if resp.Body == nil || resp.ContentLength == 0 {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

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
