package factory

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	Ctx   context.Context
	Eth   *ethclient.Client
	Http  *http.Client
	Rpc   *rpc.Client
	Data  *fx.Data
	State *fx.State
	Json  *fx.JSON
	Math  *fx.Math
	Mu    *sync.Mutex
	Rw    *sync.RWMutex
	When  *sync.Cond
}

func Assemble() *Factory {
	ctx := context.Background()
	http := &http.Client{}
	json := fx.Json(*http, ctx)

	rpc, eth, err := fx.Node(ctx)
	if err != nil {
		log.Fatalf("Error creating Ethereum node: %v", err)
	}

	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)

	data, _ := fx.Source("timefactory", ctx)
	state := fx.NewState(data, ctx)
	factory := &Factory{
		Ctx:   ctx,
		Data:  data,
		State: state,
		Json:  json,
		Eth:   eth,
		Http:  http,
		Math:  &fx.Math{},
		Rpc:   rpc,
		Mu:    mu,
		Rw:    rw,
		When:  when,
	}
	return factory
}
