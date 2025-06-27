package factory

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/fx/api"
)

type Factory struct {
	Ctx  context.Context
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Api  *api.API
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble() *Factory {
	ctx := context.Background()
	http := &http.Client{}
	api := api.NewAPI("/factory/interface", "/factory/content")

	rpc, eth, err := fx.Node(ctx)
	if err != nil {
		log.Fatalf("Error creating Ethereum node: %v", err)
	}

	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)

	factory := &Factory{
		Ctx:  ctx,
		Eth:  eth,
		Http: http,
		Rpc:  rpc,
		Api:  api,
		Mu:   mu,
		Rw:   rw,
		When: when,
	}
	return factory
}
