package fx

import (
	"context"
	"html/template"

	"github.com/gorilla/mux"
)

type Universe struct {
	Pathless *Pathless
	Path     map[string]*Value
	Frame    map[string]*template.HTML
	Router   *mux.Router
	Ctx      context.Context
}

func NewUniverse(ctx context.Context, favicon, title string) *Universe {
	return &Universe{
		Frame:  make(map[string]*template.HTML),
		Path:   make(map[string]*Value),
		Router: mux.NewRouter().StrictSlash(true),
		Ctx:    ctx,
	}
}
