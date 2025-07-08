package factory

import (
	"context"
	"html/template"
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

	factory.NewComponent(
		[]template.HTML{u.H1("Hello"), u.Paragraph("World")},
		factory.WithCSS(map[string]string{".home": "display:flex;"}),
		universe.WithJS(template.JS("console.log('hi')")),
	)

	return factory
}
