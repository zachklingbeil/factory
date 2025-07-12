package fx

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Value struct {
	Content  []byte
	MimeType string
}

// LoadEndpoints scans path for all directories (keys) and their files (values)
func (u *Universe) LoadEndpoints(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Warning: Could not read content directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		files, err := os.ReadDir(filepath.Join(path, entry.Name()))
		if err != nil {
			log.Printf("Warning: Could not read %s directory: %v", entry.Name(), err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(path, entry.Name(), file.Name())
				content, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("Warning: Could not read file %s: %v", filePath, err)
					continue
				}

				baseName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				key := entry.Name() + "/" + baseName

				mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				u.Path[key] = &Value{
					Content:  content,
					MimeType: mimeType,
				}
			}
		}
	}
}

// handlePath serves files for any key (directory) in path
func (u *Universe) handlePath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	lookupKey := key + "/" + value
	fileData, exists := u.Path[lookupKey]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", fileData.MimeType)
	w.Write(fileData.Content)
}

func corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
