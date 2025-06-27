package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// LoadEndpoints scans contentDir for all directories (keys) and their files (values)
func (a *API) LoadEndpoints(contentDir string) {
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		log.Printf("Warning: Could not read content directory: %v", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			key := entry.Name()
			dirPath := filepath.Join(contentDir, key)
			files, err := os.ReadDir(dirPath)
			if err != nil {
				log.Printf("Warning: Could not read %s directory: %v", key, err)
				continue
			}
			var fileNames []string
			for _, file := range files {
				if !file.IsDir() {
					fileNames = append(fileNames, file.Name())
				}
			}
			a.Endpoints[key] = fileNames
		}
	}
}

// handleRequest serves file listings or files for any key (directory) in contentDir
func (a *API) handleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	files, exists := a.Endpoints[key]
	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Find file by name without extension
	var long string
	for _, file := range files {
		short := strings.TrimSuffix(file, filepath.Ext(file))
		if short == value {
			long = file
			break
		}
	}

	if long == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath.Join(key, long))
}
