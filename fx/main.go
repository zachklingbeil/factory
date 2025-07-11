package fx

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/io"
	"github.com/zachklingbeil/factory/universe"
)

type Fx struct {
	*universe.Universe
	ctx  context.Context
	Json *io.Json
}

func NewFx(ctx context.Context) *Fx {
	return &Fx{
		Universe: universe.NewUniverse(ctx),
		Json:     io.NewJson(ctx),
		ctx:      ctx,
	}
}

func (f *Fx) NewPathless() *universe.Pathless {
	return universe.NewPathless("", "")
}
