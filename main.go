package factory

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	Ctx   context.Context
	Eth   *ethclient.Client
	Http  *http.Client
	Rpc   *rpc.Client
	Redis *redis.Client
	Db    *fx.Database
	Pg    *sql.DB
	State *fx.State
	Json  *fx.JSON
	Math  *fx.Math
	Mu    *sync.Mutex
	Rw    *sync.RWMutex
	When  *sync.Cond
}

func Assemble(dbName string, dbNum int) *Factory {
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

	db := &fx.Database{
		Mu:  mu,
		Rw:  rw,
		Ctx: ctx,
	}

	pg, err := db.ConnectPostgres(dbName)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}

	redis, err := db.ConnectRedis(dbNum, ctx)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	state := fx.NewState(json, mu, rw, redis, ctx)
	factory := &Factory{
		Ctx:   ctx,
		Pg:    pg,
		Redis: redis,
		State: state,
		Json:  json,
		Eth:   eth,
		Http:  http,
		Math:  &fx.Math{},
		Rpc:   rpc,
		Mu:    mu,
		Rw:    rw,
		When:  when,
	}
	return factory
}
