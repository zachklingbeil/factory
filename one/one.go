package one

import (
	"context"

	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	*fx.Fx
	*zero.Zero
}

func NewOne(ctx context.Context) *One {
	fx := fx.InitFx(ctx)
	return &One{
		Zero: zero.NewZero(),
		Fx:   fx,
	}
}
