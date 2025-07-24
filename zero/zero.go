package zero

import (
	"html/template"
)

type Zero struct {
	Frame
}

func NewZero() *Zero {
	return &Zero{
		Frame: NewFrame(),
	}
}

type One template.HTML
