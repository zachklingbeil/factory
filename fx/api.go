package fx

import (
	"context"
	"embed"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed index.html
var content embed.FS

type API struct {
	Router    *mux.Router
	Endpoints map[string][]string
	Pathless  *template.Template // Store parsed template
	Ctx       context.Context
	HTTP      *http.Client
}

func NewAPI(ctx context.Context) *API {
	tmpl, err := template.ParseFS(content, "index.html")
	if err != nil {
		log.Fatalf("Failed to parse embedded index.html: %v", err)
	}
	api := &API{
		Router:    mux.NewRouter().StrictSlash(true),
		Endpoints: make(map[string][]string),
		Ctx:       ctx,
		HTTP:      &http.Client{},
		Pathless:  tmpl, // Store parsed template
	}
	api.Router.Use(api.corsMiddleware())
	return api
}

// ServePathless serves the embedded index.html template at "/"
func (a *API) ServePathless(dir string) error {
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.Pathless.Execute(w, map[string]any{
			"Title": "My App",
			"Main":  "<h1>Hello, World!</h1>",
		})
	})
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
