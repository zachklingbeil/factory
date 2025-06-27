package factory

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/fx/api"
	"github.com/zachklingbeil/factory/fx/json"
)

type Factory struct {
	Ctx  context.Context
	Eth  *ethclient.Client
	Rpc  *rpc.Client
	Json *json.JSON
	Api  *api.API
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble() *Factory {
	ctx := context.Background()
	json := json.NewJSON(ctx)
	api := api.NewAPI(json)
	rpc, eth := fx.Node(ctx)

	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)

	factory := &Factory{
		Ctx:  ctx,
		Eth:  eth,
		Rpc:  rpc,
		Api:  api,
		Mu:   mu,
		Rw:   rw,
		When: when,
	}
	return factory
}
