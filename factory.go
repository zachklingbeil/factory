package factory

import (
	"context"
	"sync"

	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/fx/json"
)

type Factory struct {
	Ctx  context.Context
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
	Json *json.Json
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
		Json: json.NewJson(ctx),
		Fx:   fx.NewFx(ctx),
	}
	return factory
}
