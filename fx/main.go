package fx

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Universe struct {
	Pathless *Pathless
	Path     map[string]*Value
	Frame    map[string]*template.HTML
	Router   *mux.Router
	Ctx      context.Context
}

func NewUniverse(ctx context.Context) *Universe {
	return &Universe{
		Frame:  make(map[string]*template.HTML),
		Path:   make(map[string]*Value),
		Router: mux.NewRouter().StrictSlash(true),
		Ctx:    ctx,
	}
}

func (u *Universe) Circuit(favicon, title, path string) error {
	u.Pathless = &Pathless{
		Favicon:   favicon,
		Title:     title,
		Font:      "'Roboto', sans-serif",
		Primary:   "blue",
		Secondary: "red",
	}

	u.Pathless.HTML = u.Pathless.baseTemplate()
	u.Pathless.Body = template.HTML("")
	u.Router.HandleFunc("/", u.Pathless.Serve)

	u.LoadEndpoints(path)
	u.Router.Use(corsMiddleware())
	u.Router.HandleFunc("/{key}", u.handlePath).Methods("GET")
	u.Router.HandleFunc("/{key}/{value}", u.handlePath).Methods("GET")
	go func() {
		log.Fatal(http.ListenAndServe(":10001", u.Router))
	}()
	return nil
}

func (p *Pathless) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(p.HTML))
}

func (p *Pathless) Update(w http.ResponseWriter, r *http.Request, content string) {
	p.Body = template.HTML(content)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(p.Body))
}
