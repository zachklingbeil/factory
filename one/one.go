package one

import (
	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	*zero.Zero
	*fx.Fx
	*mux.Router
}

func NewOne() *One {
	o := &One{
		Zero: zero.NewZero(),
		Fx:   fx.InitFx(),
	}
	o.Router = o.NewRouter()
	o.Circuit()
	return o
}
