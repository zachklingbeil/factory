// Json i/o
// Json in: execute & decode http requests
// Json out: handle & encode http responses, server-side streaming
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

func (j *JSON) RateLimitedIn(url, apiKey string, initialRPS int) ([]byte, error) {
	// Store limiters, limits, and backoff state per URL
	type rateLimiter struct {
		tokens     float64       // Fractional tokens for sub-second precision
		ratePerSec float64       // Current determined rate per second
		lastRefill time.Time     // Last time tokens were added
		backoff    time.Duration // Increasing backoff for repeated failures
		mu         sync.Mutex
	}

	// Static store across calls
	static := struct {
		limiters map[string]*rateLimiter
		mu       sync.Mutex
	}{
		limiters: make(map[string]*rateLimiter),
	}

	// Get or create a rate limiter for this URL
	static.mu.Lock()
	rl, exists := static.limiters[url]
	if !exists {
		rl = &rateLimiter{
			tokens:     float64(initialRPS),
			ratePerSec: float64(initialRPS),
			lastRefill: time.Now(),
			backoff:    time.Second, // Initial backoff of 1 second
		}
		static.limiters[url] = rl
	}
	static.mu.Unlock()

	// Acquire a token with backoff if needed
	rl.mu.Lock()

	// Calculate token refill based on time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	if elapsed > 0 {
		// Add tokens based on rate per second and time elapsed
		newTokens := elapsed * rl.ratePerSec
		rl.tokens = min(rl.ratePerSec, rl.tokens+newTokens) // Cap at max rate
		rl.lastRefill = now
	}

	// If no tokens available, wait until we have enough
	if rl.tokens < 1.0 {
		// Calculate wait time needed for 1 token
		waitTime := time.Duration((1.0 - rl.tokens) / rl.ratePerSec * float64(time.Second))
		rl.mu.Unlock()

		select {
		case <-time.After(waitTime):
			return j.RateLimitedIn(url, apiKey, initialRPS) // Retry after waiting
		case <-j.CTX.Done():
			return nil, j.CTX.Err() // Respect context cancellation
		}
	}

	// Consume a token
	rl.tokens -= 1.0
	rl.mu.Unlock()

	// Execute the request directly instead of using j.In to check response headers
	req, err := http.NewRequestWithContext(j.CTX, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	// Add API key header if provided
	if apiKey != "" {
		req.Header.Set("X-API-KEY", apiKey)
	}

	resp, err := j.HTTP.Do(req)
	if err != nil {
		// Network error - apply backoff and reduce rate
		rl.mu.Lock()
		rl.ratePerSec = max(0.1, rl.ratePerSec/2) // Reduce rate by half, minimum 0.1 req/sec
		rl.mu.Unlock()
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check for rate limit related status codes and headers
	if resp.StatusCode == http.StatusTooManyRequests {
		// We hit a rate limit - automatically adjust
		rl.mu.Lock()

		// Reduce the rate if we hit a rate limit
		rl.ratePerSec = max(0.1, rl.ratePerSec*0.75) // Reduce to 75% of current rate
		rl.tokens = 0                                // Force waiting

		// Check if server provides Retry-After header
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
				rl.backoff = seconds
			}
		} else {
			// No Retry-After, increase backoff exponentially
			rl.backoff *= 2
			if rl.backoff > 30*time.Second {
				rl.backoff = 30 * time.Second // Cap at 30s
			}
		}

		backoff := rl.backoff
		rl.mu.Unlock()

		// Wait for the specified backoff duration
		select {
		case <-time.After(backoff):
			return j.RateLimitedIn(url, apiKey, initialRPS)
		case <-j.CTX.Done():
			return nil, j.CTX.Err()
		}
	}

	// Check for rate limit headers to adjust our rate automatically
	// Try to find rate limit info in headers
	rateInfo := struct {
		limit     int // Max requests allowed
		remaining int // Requests remaining
		reset     int // Seconds until reset
	}{
		limit:     -1,
		remaining: -1,
		reset:     -1,
	}

	// Check common rate limit headers
	for _, limitHeader := range []string{"X-RateLimit-Limit", "RateLimit-Limit"} {
		if val := resp.Header.Get(limitHeader); val != "" {
			fmt.Sscanf(val, "%d", &rateInfo.limit)
			break
		}
	}

	for _, remainingHeader := range []string{"X-RateLimit-Remaining", "RateLimit-Remaining"} {
		if val := resp.Header.Get(remainingHeader); val != "" {
			fmt.Sscanf(val, "%d", &rateInfo.remaining)
			break
		}
	}

	for _, resetHeader := range []string{"X-RateLimit-Reset", "RateLimit-Reset"} {
		if val := resp.Header.Get(resetHeader); val != "" {
			fmt.Sscanf(val, "%d", &rateInfo.reset)
			break
		}
	}

	// If we have meaningful rate limit info, adjust our rate
	if rateInfo.limit > 0 && rateInfo.reset > 0 {
		rl.mu.Lock()
		// Calculate ideal rate: limit/reset gives requests per second
		idealRate := float64(rateInfo.limit) / float64(rateInfo.reset)
		// Use 80% of the ideal rate to leave some margin
		safeRate := idealRate * 0.8

		// Gradually adjust our rate toward the safe rate
		if safeRate > rl.ratePerSec {
			// Increase rate slowly
			rl.ratePerSec = min(safeRate, rl.ratePerSec*1.1)
		} else if safeRate < rl.ratePerSec {
			// Decrease rate more quickly
			rl.ratePerSec = max(0.1, safeRate)
		}
		rl.mu.Unlock()
	} else if rateInfo.remaining > 0 {
		// If we only have remaining info, use it as a hint
		rl.mu.Lock()
		if rateInfo.remaining > int(rl.ratePerSec*10) {
			// If we have lots of requests remaining, we might be too conservative
			rl.ratePerSec = min(float64(rateInfo.remaining)/10, rl.ratePerSec*1.2)
		}
		rl.mu.Unlock()
	}

	// Check for non-OK status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// Return nil if the response body is empty
	if resp.Body == nil || resp.ContentLength == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	// Success case - gradually increase rate if we're successful
	rl.mu.Lock()
	rl.backoff = time.Second // Reset backoff
	// Very gradually increase rate on success
	if rl.tokens < 0.1 { // If we're using tokens quickly
		rl.ratePerSec *= 1.01 // Increase by 1%
	}
	rl.mu.Unlock()

	// Read and return the response body as bytes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// BatchRateLimitedIn executes multiple HTTP GET requests with rate limiting and returns results as they complete.
// It's useful for processing a batch of requests while respecting API rate limits.
// - urls: slice of URLs to request
// - apiKey: API key to use for all requests
// - limit: maximum number of requests per minute (RPM)
func (j *JSON) BatchRateLimitedIn(urls []string, apiKey string, limit int) map[string][]byte {
	results := make(map[string][]byte)
	resultsMu := sync.Mutex{}

	// Create a wait group to track completion of all requests
	var wg sync.WaitGroup
	wg.Add(len(urls))

	// Process each URL
	for _, url := range urls {
		// Use a goroutine for each request
		go func(u string) {
			defer wg.Done()

			// Execute rate-limited request
			data, err := j.RateLimitedIn(u, apiKey, limit)

			// Store result if successful
			if err == nil {
				resultsMu.Lock()
				results[u] = data
				resultsMu.Unlock()
			}
		}(url)
	}

	// Wait for all requests to complete
	wg.Wait()
	return results
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
