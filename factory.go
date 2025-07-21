package factory

import (
	"context"
	"sync"

	"github.com/zachklingbeil/factory/one"
)

type Factory struct {
	Ctx context.Context
	*one.One
}

type Motion struct {
	*sync.Mutex
	*sync.RWMutex
	*sync.Cond
}

func InitFactory() *Factory {
	ctx := context.Background()
	factory := &Factory{
		Ctx: ctx,
		One: one.NewOne(ctx),
	}
	return factory
}
