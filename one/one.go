package one

import (
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	zero.Zero
	*fx.Fx
	Api map[string]zero.Coordinate
}

func NewOne() *One {
	o := &One{
		Zero: zero.NewZero(),
		Fx:   fx.InitFx(),
		Api:  make(map[string]zero.Coordinate),
	}
	o.Circuit()
	return o
}
