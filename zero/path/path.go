package path

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Path struct {
	Map map[string]*Value
}

func NewPath() *Path {
	return &Path{
		Map: make(map[string]*Value),
	}
}

type Value struct {
	Data []byte
	Type string
}

// Walk single directory and load files into memory
// prefix from directory name
func (p *Path) AddConstants(dir string, mux *mux.Router) {
	prefix := "/" + filepath.Base(dir) + "/"

	files := make(map[string]*Value)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		fileData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		urlPath := filepath.ToSlash(relPath)
		ext := filepath.Ext(urlPath)
		urlPathWithoutExt := urlPath[:len(urlPath)-len(ext)]
		routePath := prefix + urlPathWithoutExt

		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = http.DetectContentType(fileData)
		}

		files[routePath] = &Value{
			Data: fileData,
			Type: contentType,
		}
		return nil
	})

	for routePath := range files {
		p.Map[routePath] = files[routePath]
		mux.HandleFunc(routePath, p.Read(routePath)).Methods("GET")
	}
}

func (p *Path) Read(routePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, exists := p.Map[routePath]
		if !exists {
			http.NotFound(w, r)
			return
		}
		// w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		w.Header().Set("Content-Type", file.Type)
		w.Write(file.Data)
	}
}
