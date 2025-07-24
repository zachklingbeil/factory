package fx

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory/fx/io"
)

type Fx struct {
	Ctx      context.Context
	rpc      *rpc.Client
	eth      *ethclient.Client
	postgres *sql.DB
	redis    *redis.Client
	oath     *http.Client
	*io.IO
}

func InitFx() *Fx {
	ctx := context.Background()
	fx := &Fx{
		Ctx: ctx,
		IO:  io.NewIO(ctx),
	}
	return fx
}
