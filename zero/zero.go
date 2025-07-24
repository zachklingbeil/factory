package zero

import (
	"html/template"
)

type Zero struct {
	Frame
}

func NewZero() *Zero {
	z := &Zero{
		Frame: NewFrame(),
	}
	return z
}

type One template.HTML
