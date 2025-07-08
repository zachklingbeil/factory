package fx

import (
	"context"

	"github.com/zachklingbeil/factory/fx/json"
	"github.com/zachklingbeil/factory/fx/path"
	"github.com/zachklingbeil/factory/fx/pathless"
	"github.com/zachklingbeil/factory/fx/universe"
)

type Fx struct {
	json     *json.Json
	api      *path.API
	pathless *pathless.Pathless
	universe *universe.Universe
	Ctx      context.Context
}

func NewFx(ctx context.Context) *Fx {
	return &Fx{
		json:     json.NewJson(ctx),
		api:      path.NewAPI(ctx),
		pathless: pathless.NewPathless(),
	}
}
