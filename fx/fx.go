package fx

import (
	"database/sql"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/zachklingbeil/factory/zero"
)

type Fx struct {
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	redis    *redis.Client
	Http     *http.Client
	*zero.Zero
}

func Init() *Fx {
	return &Fx{
		Zero: zero.NewZero(),
	}
}
