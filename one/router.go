package one

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (o *One) Circuit() {
	o.Path("/").HandlerFunc(o.servePathless)
	o.Path("/frame/{index}").HandlerFunc(o.serveFrame)
	o.Path("/api").HandlerFunc(o.serveAPI)
	o.Path("/api/{file}").HandlerFunc(o.serveAPIFile) // <-- Add this line
}

func (o *One) servePathless(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-FRAMES", o.FrameCount())
	pathless := o.GetPathless()
	if pathless == nil {
		http.Error(w, "No pathless content", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, *pathless)
}

// Handler for /api/{file} that serves a JSON file from the ./api directory
func (o *One) serveAPIFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]
	if !strings.HasSuffix(file, ".json") {
		file += ".json"
	}
	filePath := filepath.Join("api", file) // Adjust directory as needed

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Handler for /api that returns JSON
func (o *One) serveAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"status":  "ok",
		"message": "API endpoint",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (o *One) serveFrame(w http.ResponseWriter, r *http.Request) {
	indexStr := mux.Vars(r)["index"]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid frame index", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-FRAMES", o.FrameCount())
	frame := o.GetFrame(index)
	if frame == nil {
		http.Error(w, "Frame not found", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, *frame)
}

// Walk directory and load files into memory, determine Content-Type based on file extension. Register route/<prefix/<file without extension>.
func (o *One) AddPath(dir string, prefix string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		fileData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]
		contentType := o.getType(base, fileData)
		routePath := "/" + strings.Trim(prefix, "/") + "/" + name

		o.addRoute(routePath, fileData, contentType)
		return nil
	})
}

func (o *One) getType(filename string, data []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func (o *One) addRoute(path string, data []byte, contentType string) {
	o.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})
}
