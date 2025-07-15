package factory

import (
	"context"
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

func InitFactory(color string) *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:      ctx,
		Mu:       mu,
		When:     when,
		Rw:       &sync.RWMutex{},
		IO:       io.NewIO(ctx),
		Pathless: pathless.InitPathless(color),
	}
	return factory
}
