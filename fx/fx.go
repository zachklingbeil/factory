package fx

import (
	"database/sql"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/zachklingbeil/factory/zero"
	"goauthentik.io/api/v3"
)

type Fx struct {
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	redis    *redis.Client
	Auth     *api.APIClient // Authentik API client for management
	Http     *http.Client
	*zero.Zero
}

func InitFx() *Fx {
	return &Fx{
		Zero: zero.NewZero(),
	}
}
