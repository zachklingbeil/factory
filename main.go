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
	"github.com/zachklingbeil/factory/fx/peer"
)

type Factory struct {
	Ctx  context.Context
	Db   *fx.Database
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *fx.JSON
	Peer *peer.Peers
	Mu   sync.Mutex
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
	peer := peer.NewPeers(json, eth, db)
	go peer.HelloUniverse()

	factory := &Factory{
		Rpc:  rpc,
		Eth:  eth,
		Http: http,
		Json: json,
		Ctx:  ctx,
		Db:   db,
		Peer: peer,
	}

	fmt.Printf("factory [ %s ]\n", dbName)
	return factory, nil
}
