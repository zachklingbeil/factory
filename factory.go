package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/api/json"
	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	Ctx  context.Context
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *json.JSON
	Sync *State
}

type State struct {
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble() *Factory {
	ctx := context.Background()
	http := &http.Client{}
	json := json.Json(*http, ctx)

	rpc, eth, err := fx.Node(ctx)
	if err != nil {
		log.Fatalf("Error creating Ethereum node: %v", err)
	}

	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)
	sync := &State{
		Mu:   mu,
		Rw:   rw,
		When: when,
	}

	factory := &Factory{
		Ctx:  ctx,
		Json: json,
		Eth:  eth,
		Http: http,
		Rpc:  rpc,
		Sync: sync,
	}
	return factory
}
