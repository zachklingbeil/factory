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
	Peer *fx.Peers
	Mu   sync.Mutex
}

// NewFactory initializes the Factory with all required components, including the database connection.
func NewFactory(dbName string) (*Factory, error) {
	ctx := context.Background()
	http := &http.Client{}
	json := fx.Json(*http, ctx)
	rpc, eth, _ := fx.Node(ctx)
	db, _ := fx.NewDatabase(dbName)
	peer := fx.NewPeers(json, eth, db)

	factory := &Factory{
		Rpc:  rpc,
		Eth:  eth,
		Http: http,
		Json: json,
		Ctx:  ctx,
		Db:   db,
		Peer: peer,
	}
	fmt.Printf("- Initialized Ethereum RPC (ipc), http, json, and peers\n")
	fmt.Printf("- Database: %s\n", dbName)
	return factory, nil
}
