package fx

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/fx/api"
	"github.com/zachklingbeil/factory/fx/json"
	"github.com/zachklingbeil/factory/fx/universe"
)

type Fx struct {
	*api.API
	*universe.Universe
	ctx  context.Context
	Json *json.Json
}

func NewFx(ctx context.Context) *Fx {
	return &Fx{
		API:      api.NewAPI(ctx),
		Universe: universe.NewUniverse(),
		Json:     json.NewJson(ctx),
		ctx:      ctx,
	}
}
