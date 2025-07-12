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
	Frame    map[string]*template.HTML
	Handlers map[string]http.Handler
	Path     map[string][]byte
	Router   *mux.Router
	Ctx      context.Context
}

func NewUniverse(ctx context.Context) *Universe {
	return &Universe{
		Frame:    make(map[string]*template.HTML),
		Router:   mux.NewRouter().StrictSlash(true),
		Handlers: make(map[string]http.Handler),
		Ctx:      ctx,
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

func (u *Universe) Start(port string) error {
	return http.ListenAndServe(port, u.Router)
}

func (u *Universe) ServeHTML(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}

func (u *Universe) AddGET(pattern string, handler http.HandlerFunc) {
	u.Handlers[pattern] = handler
	u.Router.Handle(pattern, handler).Methods("GET")
}

func (u *Universe) AddPOST(pattern string, handler http.HandlerFunc) {
	u.Handlers[pattern] = handler
	u.Router.Handle(pattern, handler).Methods("POST")
}
func (u *Universe) AddPathless(handler http.HandlerFunc) {
	u.Handlers["/"] = handler
	u.Router.Handle("/", handler).Methods("GET")
}

func (u *Universe) AddPath(path string, dir http.FileSystem) {
	u.Router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(dir)))
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
