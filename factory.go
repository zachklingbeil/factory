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

func InitFactory(css, js string) *Factory {
	return &Factory{
		One: one.NewOne(css, js),
	}
}
