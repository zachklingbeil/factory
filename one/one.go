package one

import (
	"github.com/zachklingbeil/factory/fx"
)

type One struct {
	*fx.Fx
}

func NewOne() *One {
	return &One{
		Fx: fx.InitFx(),
	}
}
