package factory

import (
	"context"
	"sync"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/universe"
)

type Factory struct {
	Ctx    context.Context
	Mu     *sync.Mutex
	Rw     *sync.RWMutex
	When   *sync.Cond
	Router *mux.Router
	*io.IO
	*universe.Universe
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:    ctx,
		Mu:     mu,
		Rw:     rw,
		When:   when,
		IO:     io.NewIO(ctx),
		Router: mux.NewRouter().StrictSlash(true),
	}

	return factory
}

func (f *Factory) HelloUniverse(favicon, title, url string) *universe.Universe {
	if f.Universe == nil {
		f.Universe = universe.HelloUniverse(favicon, title, url)
	}
	return f.Universe
}
