package fx

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/pathless"
)

type Universe struct {
	Pathless *pathless.Pathless
	Frame    map[string]*template.HTML
	Handlers map[string]http.Handler
	Path     map[string][]byte
	Router   *mux.Router
	Ctx      context.Context
}

func NewUniverse(ctx context.Context, mux *mux.Router) *Universe {
	return &Universe{
		Router:   mux,
		Ctx:      ctx,
		Frame:    make(map[string]*template.HTML),
		Handlers: make(map[string]http.Handler),
	}
}

func (u *Universe) Circuit(favicon, title, path string) error {
	u.Pathless = &pathless.Pathless{
		Favicon:   favicon,
		Title:     title,
		Font:      "'Roboto', sans-serif",
		Primary:   "blue",
		Secondary: "red",
	}

	u.Pathless.HTML = u.Pathless.Zero()
	u.Pathless.Body = template.HTML("")
	u.Router.HandleFunc("/", u.Serve)

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

func (u *Universe) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(u.Pathless.HTML))
}

func (u *Universe) Update(w http.ResponseWriter, r *http.Request, content string) {
	u.Pathless.Body = template.HTML(content)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(u.Pathless.Body))
}
