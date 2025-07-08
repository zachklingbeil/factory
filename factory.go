package factory

import (
	"context"
	"sync"

	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	Ctx  context.Context
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
	*fx.Fx
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)
	factory := &Factory{
		Ctx:  ctx,
		Mu:   mu,
		Rw:   rw,
		When: when,
		Fx:   fx.NewFx(ctx),
	}
	return factory
}
