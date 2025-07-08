package path

import (
	"context"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	Endpoints map[string]any
	Router    *mux.Router
	Ctx       context.Context
}

func NewAPI(ctx context.Context) *API {
	api := &API{
		Router:    mux.NewRouter().StrictSlash(true),
		Endpoints: make(map[string]any),
		Ctx:       ctx,
	}

	api.Router.Use(api.corsMiddleware())
	api.Router.HandleFunc("/{key}", api.handleRequest).Methods("GET")
	api.Router.HandleFunc("/{key}/{value}", api.handleRequest).Methods("GET")
	go func() {
		log.Fatal(http.ListenAndServe(":10002", api.Router))
	}()
	return api
}

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

	rawFiles, exists := a.Endpoints[key]
	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	files, ok := rawFiles.([]string)
	if !ok {
		http.Error(w, "Invalid endpoint data", http.StatusInternalServerError)
		return
	}

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

	filePath := filepath.Join(key, long)
	mimeType := mime.TypeByExtension(filepath.Ext(long))
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	http.ServeFile(w, r, filePath)
}
