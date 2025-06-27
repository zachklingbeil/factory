package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type Path struct {
	Endpoints map[string][]string
	Router    *mux.Router
	CTX       context.Context
}

func NewPath(port string, router *mux.Router, ctx context.Context) *Path {
	p := &Path{
		Endpoints: make(map[string][]string),
		Router:    router,
		CTX:       ctx,
	}

	p.Router.HandleFunc("/{type}", p.handleRequest).Methods("GET")
	p.Router.HandleFunc("/{type}/{filename}", p.handleRequest).Methods("GET")

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, p.Router))
	}()

	return p
}

func (p *Path) LoadEndpoints(contentDir string) {
	for _, contentType := range []string{"text", "image", "video"} {
		dirPath := filepath.Join(contentDir, contentType)
		if files, err := os.ReadDir(dirPath); err == nil {
			var fileNames []string
			for _, file := range files {
				if !file.IsDir() {
					fileNames = append(fileNames, file.Name())
				}
			}
			p.Endpoints[contentType] = fileNames
		} else {
			log.Printf("Warning: Could not read %s directory: %v", contentType, err)
		}
	}
}

func (p *Path) handleRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contentType := vars["type"]
	filename := vars["filename"]

	files, exists := p.Endpoints[contentType]
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

	http.ServeFile(w, r, filepath.Join(p.Content, contentType, long))
}
