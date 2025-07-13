package factory

import (
	"context"
	"sync"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/io"
)

type Factory struct {
	Ctx    context.Context
	Mu     *sync.Mutex
	Rw     *sync.RWMutex
	When   *sync.Cond
	Json   *io.Json
	Router *mux.Router
	*fx.Universe
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
		Router: mux.NewRouter().StrictSlash(true),
	}
	return factory
}
