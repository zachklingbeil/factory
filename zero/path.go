package zero

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type Value struct {
	Data []byte
	Type string
}

// Walk single directory and load files into memory
// Returns a map of file name (without extension) to *Value
func (z *Zero) AddPath(dir string) (map[string]*Value, error) {
	files := make(map[string]*Value)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		base := filepath.Base(relPath)
		name := base[:len(base)-len(filepath.Ext(base))] // file name without extension

		contentType := mime.TypeByExtension(filepath.Ext(base))
		if contentType == "" {
			contentType = http.DetectContentType(fileData)
		}

		files[name] = &Value{
			Data: fileData,
			Type: contentType,
		}
		return nil
	})
	return files, err
}

// func (z *Zero) Read(routePath string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if file, exists := z.Map[routePath]; exists {
// 			// w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
// 			w.Header().Set("Content-Type", file.Type)
// 			w.Write(file.Data)
// 			return
// 		}
// 		http.NotFound(w, r)
// 	}
// }
