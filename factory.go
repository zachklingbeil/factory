package factory

import (
	"context"
	"sync"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero/frame"
	"github.com/zachklingbeil/factory/zero/path"
	"github.com/zachklingbeil/factory/zero/pathless"
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
	*fx.Fx
	*frame.Frame
	*pathless.Pathless
	*path.Path
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
			Fx:       fx,
			Frame:    frame.NewFrame(driver),
			Pathless: pathless.NewPathless(),
			Path:     path.NewPath(),
		},
		Motion: &Motion{Mutex: mu, RWMutex: &sync.RWMutex{}, Cond: when},
		Router: fx.NewRouter(),
		Map:    make(map[string]*any),
	}
	return factory
}
