package factory

import (
	"context"
	"sync"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/element"
	"github.com/zachklingbeil/factory/frame"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/factory/path"
	"github.com/zachklingbeil/factory/pathless"
)

type Factory struct {
	Ctx context.Context
	*mux.Router
	Map map[string]*any
	*Motion
	*Lines
}

type Motion struct {
	*sync.Mutex
	*sync.RWMutex
	*sync.Cond
}

type Lines struct {
	*io.IO
	*fx.Fx
	*frame.Frame
	*pathless.Pathless
	*path.Path
	*element.Element
}

func InitFactory() *Factory {
	ctx := context.Background()
	mu := &sync.Mutex{}
	when := sync.NewCond(mu)
	fx := fx.InitFx(ctx)
	driver := fx.NewRouter()

	factory := &Factory{
		Ctx: ctx,
		Lines: &Lines{
			IO:       io.NewIO(ctx),
			Fx:       fx,
			Frame:    frame.NewFrame(driver),
			Pathless: pathless.NewPathless(),
			Path:     path.NewPath(),
			Element:  &element.Element{},
		},
		Motion: &Motion{Mutex: mu, RWMutex: &sync.RWMutex{}, Cond: when},
		Router: fx.NewRouter(),
		Map:    make(map[string]*any),
	}

	return factory
}
