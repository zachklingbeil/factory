package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type API struct {
	Path     *Path
	Pathless *Pathless
	HTTP     *http.Client
	Router   *mux.Router
	CTX      context.Context
}

func NewAPI(dir, contentDir string) *API {
	api := &API{
		HTTP:   &http.Client{},
		CTX:    context.Background(),
		Router: mux.NewRouter().StrictSlash(true),
	}

	api.Router.Use(api.corsMiddleware())
	go func() {
		log.Fatal(http.ListenAndServe(":10001", api.Router))
	}()
	return api
}

func (a *API) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
