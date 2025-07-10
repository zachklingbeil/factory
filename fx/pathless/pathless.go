package pathless

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed index.html
var index embed.FS

type Pathless struct {
	Favicon string
	Title   string
	router  *mux.Router
	zero    *template.Template
	Body    template.HTML
}

func NewPathless() *Pathless {
	one, err := template.ParseFS(index, "index.html")
	if err != nil {
		log.Fatalf("failed to parse embedded index.html: %v", err)
	}

	p := &Pathless{
		router: mux.NewRouter().StrictSlash(true),
		zero:   one,
	}

	p.router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	))
	p.router.HandleFunc("/", p.serve)
	go func() {
		log.Fatal(http.ListenAndServe(":10001", p.router))
	}()
	return p
}

func (p *Pathless) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := p.zero.Execute(w, nil); err != nil {
		http.Error(w, "failed to render index.html", http.StatusInternalServerError)
	}
}
