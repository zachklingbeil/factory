package fx

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Universe struct {
	Path   map[string][]byte
	Router *mux.Router
	Ctx    context.Context
}

func NewUniverse(ctx context.Context, mux *mux.Router) *Universe {
	return &Universe{
		Router: mux,
		Ctx:    ctx,
	}
}

func (u *Universe) Start(port string) error {
	return http.ListenAndServe(port, u.Router)
}

func (u *Universe) ServeHTML(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}

func (u *Universe) AddPath(path string, dir http.FileSystem) {
	u.Router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(dir)))
}
