package json

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

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
