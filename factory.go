package factory

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/factory/universe"
	"github.com/zachklingbeil/factory/universe/pathless"
)

type Factory struct {
	Ctx  context.Context
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
	*io.IO
	*universe.Universe
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:      ctx,
		Mu:       mu,
		Rw:       rw,
		When:     when,
		IO:       io.NewIO(ctx),
		Universe: universe.NewUniverse(),
	}
	return factory
}

func (f *Factory) HelloUniverse(favicon, title, url string) error {
	f.Universe.Pathless = pathless.NewPathless(favicon, title, url)
	f.Universe.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(*f.Universe.HTML))
	})
	go func() {
		http.ListenAndServe(":10101", f.Router)
	}()
	return nil
}

func (f *Factory) AddFrame(name string, elements ...template.HTML) {
	frame := f.Universe.CreateFrame(elements...)
	f.Universe.Map[name] = frame
	f.Universe.HandleFunc("/0/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(*frame))
	})
}

// LoadMarkdownFile loads a single markdown file and creates a frame for it
func (f *Factory) LoadText(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".md" && ext != ".markdown" {
		return fmt.Errorf("file %s is not a markdown file", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Use filename without extension as frame name
	fileName := filepath.Base(filePath)
	frameName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	htmlContent := f.MarkdownToHTML(string(content))
	f.AddFrame(frameName, htmlContent)
	return nil
}
