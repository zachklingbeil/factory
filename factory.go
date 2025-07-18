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
	*mux.Router
	Map map[string]*any
	*Motion
	*Lines
}

type Motion struct {
	*sync.Mutex
	*sync.RWMutex
	*sync.Cond
}

type Lines struct {
	*io.IO
	*fx.Fx
	*frame.Frame
	*pathless.Pathless
	*path.Path
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	fx := fx.InitFx(ctx)
	driver := fx.NewRouter()

	factory := &Factory{
		Ctx: ctx,
		Lines: &Lines{
			IO:       io.NewIO(ctx),
			Fx:       fx,
			Frame:    frame.NewFrame(driver),
			Pathless: pathless.NewPathless(),
			Path:     path.NewPath(),
		},
		Motion: &Motion{Mutex: mu, RWMutex: &sync.RWMutex{}, Cond: when},
		Router: fx.NewRouter(),
		Map:    make(map[string]*any),
	}
	return factory
}

// SetMapValue sets a value in the Factory's Map.
func (f *Factory) Input(key string, value any) {
	f.RWMutex.Lock()
	defer f.RWMutex.Unlock()
	f.Map[key] = &value
}

// GetMapValue retrieves a value from the Factory's Map.
func (f *Factory) Output(key string) (any, bool) {
	f.RWMutex.RLock()
	defer f.RWMutex.RUnlock()
	valPtr, ok := f.Map[key]
	if !ok || valPtr == nil {
		return nil, false
	}
	return *valPtr, true
}

// Register endpoint template and route
func (f *Factory) RegisterFrameRoutes() {
	for path, tmpl := range f.Map {
		f.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if htmlTmpl, ok := (*tmpl).(*template.HTML); ok {
				f.WriteResponse(w, htmlTmpl)
			} else {
				f.WriteResponse(w, nil)
			}
		})
	}
}

// Write HTML response
func (f *Factory) WriteResponse(w http.ResponseWriter, tmpl *template.HTML) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if tmpl != nil {
		w.Write([]byte(string(*tmpl)))
	} else {
		w.Write([]byte("<div>404 Not Found</div>"))
	}
}

func (f *Factory) InitPathless(body template.HTML, cssPath string) {
	f.Zero(body, cssPath)
	f.HandleFunc("/", f.One).Methods("GET")
	go func() {
		log.Println("Starting pathless on :1001")
		http.ListenAndServe(":1001", f.Router)
	}()
}

func (f *Factory) AddText(file string, elements ...template.HTML) template.HTML {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	markdown := f.FromMarkdown(file, elements...)
	return markdown
}
