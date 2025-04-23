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
	Ctx  context.Context
	Db   *fx.Database
	Json *fx.JSON
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble(dbName string, rdb int) *Factory {
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

	db, err := fx.Connect(dbName, rdb, ctx, mu, rw)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	factory := &Factory{
		Ctx:  ctx,
		Db:   db,
		Json: json,
		Eth:  eth,
		Http: http,
		Rpc:  rpc,
		Mu:   mu,
		Rw:   rw,
		When: when,
	}
	return factory
}
