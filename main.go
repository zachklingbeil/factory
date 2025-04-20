package factory

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	Ctx  context.Context
	Db   *fx.Database
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *fx.JSON
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble(dbName string, distance time.Duration) *Factory {
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

	db, err := fx.Connect(dbName, distance, ctx, mu, rw)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	factory := &Factory{
		Ctx:  ctx,
		Db:   db,
		Eth:  eth,
		Http: http,
		Rpc:  rpc,
		Json: json,
		Mu:   mu,
		Rw:   rw,
		When: when,
	}
	return factory
}
