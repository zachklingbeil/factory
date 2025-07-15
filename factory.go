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
	*mux.Router
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:      ctx,
		Mutex:    mu,
		Cond:     when,
		RWMutex:  &sync.RWMutex{},
		IO:       io.NewIO(ctx),
		Fx:       fx.InitFx(ctx),
		Pathless: &pathless.Pathless{},
		Router:   mux.NewRouter().StrictSlash(true),
	}
	return factory
}

func (f *Factory) InitPathless(color string, body template.HTML) {
	f.Pathless = &pathless.Pathless{
		Font:  "'Roboto', sans-serif",
		Color: color,
		Md:    pathless.InitGoldmark(),
	}
	f.Zero(body)
	f.HandleFunc("/", f.One).Methods("GET")
	go func() {
		log.Println("Starting pathless on :10101")
		http.ListenAndServe(":10101", f.Router)
	}()
}
