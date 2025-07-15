package factory

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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
	*pathless.Pathless
	*path.Path
	*mux.Router
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	fx := fx.InitFx(ctx)
	router := fx.NewRouter()
	factory := &Factory{
		Ctx:      ctx,
		Mutex:    mu,
		Cond:     when,
		RWMutex:  &sync.RWMutex{},
		Fx:       fx,
		Router:   router,
		IO:       io.NewIO(ctx),
		Pathless: pathless.NewPathless("blue"),
		Path:     path.NewPath(),
	}
	return factory
}

func (f *Factory) InitPathless(body template.HTML) {
	f.Zero(body)
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
