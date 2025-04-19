package factory

import (
	"context"
	"fmt"
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
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *fx.JSON
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	When *sync.Cond
}

func Assemble(dbName string) *Factory {
	ctx := context.Background()
	http := &http.Client{}
	json := fx.Json(*http, ctx)

	rpc, eth, err := fx.Node(ctx)
	if err != nil {
		log.Fatalf("Error creating Ethereum node: %v", err)
	}

	db, err := fx.Connect(ctx, dbName)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	fmt.Printf("[ %s ]\n", dbName)

	mu := &sync.Mutex{}
	rw := &sync.RWMutex{}
	when := sync.NewCond(mu)

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
