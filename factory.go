package factory

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/frame"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/factory/path"
	"github.com/zachklingbeil/factory/pathless"
)

type Factory struct {
	Ctx context.Context
	*sync.Mutex
	*sync.RWMutex
	*sync.Cond
	*io.IO
	*fx.Fx
	*frame.Frame
	*pathless.Pathless
	*path.Path
	*mux.Router
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	fx := fx.InitFx(ctx)
	factory := &Factory{
		Ctx:      ctx,
		Mutex:    mu,
		Cond:     when,
		RWMutex:  &sync.RWMutex{},
		Fx:       fx,
		Router:   fx.NewRouter(),
		Pathless: pathless.NewPathless(),
		Path:     path.NewPath(),
		Frame:    frame.NewFrame(),
		IO:       io.NewIO(ctx),
	}
	return factory
}

func (f *Factory) InitPathless(body template.HTML, cssPath string) {
	f.Zero(body, cssPath)
	f.HandleFunc("/", f.One).Methods("GET")
	go func() {
		log.Println("Starting pathless on :1001")
		http.ListenAndServe(":1001", f.Router)
	}()
}

func (f *Factory) Swap(newBody template.HTML) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	f.Pathless.HTML = &newBody
}

func (f *Factory) AddConstant(dir string) {
	f.AddConstants(dir, f.Router)
}

func (f *Factory) AddText(file string, elements ...template.HTML) template.HTML {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	markdown := f.FromMarkdown(file, elements...)
	return markdown
}
