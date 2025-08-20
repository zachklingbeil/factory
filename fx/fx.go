package fx

import (
	"context"
	"database/sql"
	"net/http"
	"os"

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
	Rpc      *rpc.Client
	Eth      *ethclient.Client
	postgres *sql.DB
	redis    *redis.Client
	*mux.Router
	*io.IO
	apiKey string
}

func InitFx() *Fx {
	ctx := context.Background()
	fx := &Fx{
		Ctx:    ctx,
		IO:     io.NewIO(ctx),
		apiKey: os.Getenv("API_KEY"),
	}
	fx.Node()
	fx.Router = fx.NewRouter()
	fx.HandleFunc("/geth", fx.withAPIKey(fx.GethHandler)).Methods(http.MethodPost)
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

// withAPIKey wraps a handler with API key validation
func (f *Fx) withAPIKey(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != f.apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	})
}

func (f *Fx) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "X-FRAMES", "Content-Type", "Peer", "Cache-Control", "Connection", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
	)
}
