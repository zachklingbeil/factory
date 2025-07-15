package factory

import (
	"context"
	"html/template"
	"net/http"
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

func (f *Factory) AddFrame(name string, elements ...template.HTML) *template.HTML {
	frame := f.Universe.CreateFrame(elements...)
	f.Universe.Map[name] = frame
	f.Universe.HandleFunc("/0/"+name, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(*frame))
	})
	return frame
}
