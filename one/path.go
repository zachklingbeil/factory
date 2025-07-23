package one

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Value struct {
	Data []byte
	Type string
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
		name := base[:len(base)-len(filepath.Ext(base))] // file name without extension

		contentType := mime.TypeByExtension(filepath.Ext(base))
		if contentType == "" {
			contentType = http.DetectContentType(fileData)
		}

		routePath := "/" + strings.Trim(prefix, "/") + "/" + name
		fileVal := &Value{
			Data: fileData,
			Type: contentType,
		}
		o.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", fileVal.Type)
			w.Write(fileVal.Data)
		})

		return nil
	})
}
