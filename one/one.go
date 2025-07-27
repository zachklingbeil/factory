package one

import (
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	zero.Zero
	*fx.Fx
}

func NewOne() *One {
	o := &One{
		Zero: zero.NewZero(),
		Fx:   fx.InitFx(),
	}
	o.Circuit()
	return o
}
