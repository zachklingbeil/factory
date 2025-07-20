package zero

import (
	"html/template"

	"github.com/zachklingbeil/factory/element"
)

type Zero struct {
	Pathless *template.HTML
	X        []element.Frame
}

func NewZero() *Zero {
	return &Zero{
		X: make([]element.Frame, 0),
	}
}
