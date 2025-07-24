package factory

import (
	"sync"

	"github.com/zachklingbeil/factory/one"
)

type Factory struct {
	*one.One
}

type Motion struct {
	*sync.Mutex
	*sync.RWMutex
	*sync.Cond
}

func InitFactory() *Factory {
	return &Factory{
		One: one.NewOne(),
	}
}
