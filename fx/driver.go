package fx

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Driver struct {
	Router   *mux.Router
	Handlers map[string]http.Handler
}

func NewDriver() *Driver {
	return &Driver{
		Router:   mux.NewRouter().StrictSlash(true),
		Handlers: make(map[string]http.Handler),
	}
}

func (d *Driver) Start(port string) error {
	return http.ListenAndServe(port, d.Router)
}

func (d *Driver) ServeHTML(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}

func (d *Driver) AddGET(pattern string, handler http.HandlerFunc) {
	d.Handlers[pattern] = handler
	d.Router.Handle(pattern, handler).Methods("GET")
}

func (d *Driver) AddPOST(pattern string, handler http.HandlerFunc) {
	d.Handlers[pattern] = handler
	d.Router.Handle(pattern, handler).Methods("POST")
}
func (d *Driver) AddPathless(handler http.HandlerFunc) {
	d.Handlers["/"] = handler
	d.Router.Handle("/", handler).Methods("GET")
}

func (d *Driver) AddPath(path string, dir http.FileSystem) {
	d.Router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(dir)))
}
