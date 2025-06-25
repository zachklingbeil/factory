package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type API struct {
	Interface string
	Content   string
	Endpoints map[string][]string
}

func NewAPI(dir string, contentDir string) *API {
	api := &API{
		Interface: dir,
		Content:   contentDir,
		Endpoints: make(map[string][]string),
	}

	api.loadEndpoints()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(api.corsMiddleware())
	router.HandleFunc("/{type}", api.handleRequest).Methods("GET")
	router.HandleFunc("/{type}/{filename}", api.handleRequest).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(api.Interface)))

	go func() {
		log.Fatal(http.ListenAndServe(":10001", router))
	}()
	return api
}

func (a *API) loadEndpoints() {
	for _, contentType := range []string{"text", "image", "video"} {
		dirPath := filepath.Join(a.Content, contentType)
		if files, err := os.ReadDir(dirPath); err == nil {
			var fileNames []string
			for _, file := range files {
				if !file.IsDir() {
					fileNames = append(fileNames, file.Name())
				}
			}
			a.Endpoints[contentType] = fileNames
		} else {
			log.Printf("Warning: Could not read %s directory: %v", contentType, err)
		}
	}
}

func (a *API) handleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentType := vars["type"]
	filename := vars["filename"]

	files, exists := a.Endpoints[contentType]
	if !exists {
		http.Error(w, "Content type not found", http.StatusNotFound)
		return
	}

	// If no filename, return file listing with type info
	if filename == "" {
		result := map[string]any{
			"type":  contentType,
			"files": files,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	// Find file by name without extension
	var long string
	for _, file := range files {
		short := strings.TrimSuffix(file, filepath.Ext(file))
		if short == filename {
			long = file
			break
		}
	}

	if long == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath.Join(a.Content, contentType, long))
}

func (a *API) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
