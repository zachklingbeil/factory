package fx

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type API struct {
	Router    *mux.Router
	Endpoints map[string][]string
	Pathless  string
	Ctx       context.Context
	HTTP      *http.Client
}

func NewAPI(ctx context.Context) *API {
	api := &API{
		Router:    mux.NewRouter().StrictSlash(true),
		Endpoints: make(map[string][]string),
		Ctx:       ctx,
		HTTP:      &http.Client{},
	}
	api.Router.Use(api.corsMiddleware())
	return api
}

func (a *API) ServePathless(dir string) error {
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))
	go func() {
		log.Fatal(http.ListenAndServe(":10002", a.Router))
	}()
	return nil
}

func (a *API) NewPath(dir string) error {
	a.LoadEndpoints(dir)
	a.Router.HandleFunc("/{key}", a.handleRequest).Methods("GET")
	a.Router.HandleFunc("/{key}/{value}", a.handleRequest).Methods("GET")

	go func() {
		log.Fatal(http.ListenAndServe(":10003", a.Router))
	}()
	return nil
}

func (a *API) corsMiddleware() mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
