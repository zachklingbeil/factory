package universe

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/universe/frame"
	"github.com/zachklingbeil/factory/universe/pathless"
)

type Universe struct {
	Map map[string]*template.HTML
	*frame.Frame
	*pathless.Pathless
	*mux.Router
}

func NewUniverse() *Universe {
	return &Universe{
		Frame:  &frame.Frame{},
		Map:    make(map[string]*template.HTML),
		Router: mux.NewRouter().StrictSlash(true),
	}
}
