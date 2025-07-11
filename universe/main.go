package universe

import (
	"context"
	"html/template"

	"github.com/gorilla/mux"
)

type Universe struct {
	Frame    map[string]*template.HTML
	Router   *mux.Router
	Pathless *Pathless
	Path     map[string]*Value
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
