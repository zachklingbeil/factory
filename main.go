// Factory provides a common context for sourcing and distrubting data.
// Includes an Ethereum, HTTP, RPC client, a database connection, and json i/o logic.
package factory

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zachklingbeil/factory/cmd"
)

type Factory struct {
	Ctx  context.Context
	Db   *sql.DB
	Eth  *ethclient.Client
	Http *http.Client
	Rpc  *rpc.Client
	Json *cmd.JSON
	Mu   sync.Mutex
}

func NewFactory() (*Factory, error) {
	ctx := context.Background()
	http := &http.Client{}

	rpc, eth, err := cmd.Node(ctx)
	if err != nil {
		return nil, err
	}

	db := cmd.Database()
	if db == nil {
		return nil, err
	}

	json := cmd.Json(*http, ctx)
	return &Factory{
		Rpc:  rpc,
		Eth:  eth,
		Http: http,
		Json: json,
		Db:   db,
		Ctx:  ctx,
	}, nil
}
