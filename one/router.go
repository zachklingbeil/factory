package one

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (o *One) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(o.corsMiddleware())
	go func() {
		log.Println(":1001")
		http.ListenAndServe(":1001", router)
	}()
	return router
}

func (o *One) Circuit() {
	o.Path("/").HandlerFunc(o.servePathless)
	o.Path("/frame/{index}").HandlerFunc(o.serveFrame)
}

func (o *One) servePathless(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-FRAMES", o.FrameCount())
	fmt.Fprint(w, *o.GetPathless())
}

func (o *One) serveFrame(w http.ResponseWriter, r *http.Request) {
	indexStr := mux.Vars(r)["index"]
	index, err := strconv.Atoi(indexStr)

	if err != nil {
		http.Error(w, "Invalid frame index", http.StatusBadRequest)
		return
	}

	frame, exists := o.GetFrame(index)
	if !exists {
		http.Error(w, "Frame not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, *frame)
}

func (o *One) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "X-FRAMES", "Content-Type", "Peer", "Cache-Control", "Connection", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
	)
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
