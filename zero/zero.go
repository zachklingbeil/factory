package zero

import (
	_ "embed"
	"html/template"
)

type One template.HTML

//go:embed pathless.html
var pathless string

// Pathless returns the embedded pathless.html as *One
func (z *zero) Pathless() *One {
	return z.pathless
}

type Zero interface {
	Pathless() *One
	Build
}

// --- zero Implementation ---
type zero struct {
	pathless *One
	*build
}

func NewZero() Zero {
	result := One(template.HTML(pathless))
	return &zero{
		pathless: &result,
		build:    NewBuild().(*build),
	}
}
