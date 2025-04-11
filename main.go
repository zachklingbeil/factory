package factory

import (
	"context"
	"fmt"
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
	Mu   *sync.Mutex   // Mutex for exclusive access
	Rw   *sync.RWMutex // RWMutex for read-heavy operations
	When *sync.Cond    // Signal
}

func NewFactory(dbName string) (*Factory, error) {
	ctx := context.Background()
	http := &http.Client{}
	json := fx.Json(*http, ctx)
	rpc, eth, err := fx.Node(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum node: %w", err)
	}
	db, _ := fx.NewDatabase(dbName)
	fmt.Printf("factory [ %s ]\n", dbName)
	// peer := peer.NewPeers(json, eth, db)

	mu := &sync.Mutex{}
	factory := &Factory{
		Rpc:  rpc,
		Eth:  eth,
		Http: http,
		Json: json,
		Ctx:  ctx,
		Db:   db,
		Mu:   mu,
		Rw:   &sync.RWMutex{},
		When: sync.NewCond(mu),
	}
	return factory, nil
}
