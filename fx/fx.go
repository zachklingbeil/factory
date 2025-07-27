package fx

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	*mux.Router
	*io.IO
}

func InitFx() *Fx {
	ctx := context.Background()
	fx := &Fx{
		Ctx: ctx,
		IO:  io.NewIO(ctx),
	}
	fx.Router = fx.NewRouter()
	return fx
}

func (f *Fx) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(f.corsMiddleware())
	go func() {
		http.ListenAndServe(":1001", router)
	}()
	return router
}

func (f *Fx) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "X-FRAMES", "Content-Type", "Peer", "Cache-Control", "Connection", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
	)
}
