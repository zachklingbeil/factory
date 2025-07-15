package factory

import (
	"context"
	"html/template"
	"sync"

	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/factory/pathless"
)

type Factory struct {
	Ctx  context.Context
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
	*io.IO
	*pathless.Pathless
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:  ctx,
		Mu:   mu,
		When: when,
		Rw:   &sync.RWMutex{},
		IO:   io.NewIO(ctx),
	}
	return factory
}

func (f *Factory) InitPathless(color string, body template.HTML) {
	f.Pathless = pathless.InitPathless(color, body)
}
