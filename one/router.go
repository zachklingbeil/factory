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
	"github.com/zachklingbeil/factory/zero"
)

func (o *One) Circuit() {
	o.Path("/").HandlerFunc(o.servePathless)
	o.Path("/frame/{index}").HandlerFunc(o.serveFrame)
}

func (o *One) RegisterAPI(prefix string, coords []zero.Coordinate) {
	o.Path(prefix).HandlerFunc(o.makeAPIHandler(prefix))
	o.Path(prefix + "/{x}").HandlerFunc(o.makeAPIHandler(prefix))
	o.Path(prefix + "/{x}/{y}").HandlerFunc(o.makeAPIHandler(prefix))
}

func (o *One) makeAPIHandler(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		x, xOK := vars["x"]
		y, yOK := vars["y"]

		_, ok := o.Api[prefix]
		if !ok {
			http.Error(w, "No data for prefix", http.StatusNotFound)
			return
		}

		var result any
		switch {
		case xOK && yOK:
			result = o.GetZ(x, y)
		case xOK:
			result = o.GetY(x)
		default:
			result = o.GetX()
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
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
