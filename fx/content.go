package fx

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

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

				u.Path[key] = content
			}
		}
	}
}

// handlePath serves files for any key
func (u *Universe) handlePath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	fileData, exists := u.Path[key]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Write(fileData)
}

func corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
