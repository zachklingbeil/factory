package zero

import (
	_ "embed"
	"os"
)

//go:embed embed/coordinate.css
var coordinateCSS string

//go:embed embed/coordinate.js
var coordinateJS string

//go:embed embed/one.js
var oneJS string

//go:embed embed/one.css
var oneCSS string

type Embed interface {
	CoordinateCSS() string
	CoordinateJS() string
	OneJS() string
	OneCSS(path string) string
}

type embed struct{}

func NewEmbed() Embed {
	return &embed{}
}

func (a *embed) CoordinateCSS() string { return coordinateCSS }
func (a *embed) CoordinateJS() string  { return coordinateJS }
func (a *embed) OneJS() string         { return oneJS }

// OneCSS returns the embedded oneCSS plus the contents of the file at path (if provided)
func (a *embed) OneCSS(path string) string {
	if path == "" {
		return oneCSS
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return oneCSS
	}
	return oneCSS + "\n" + string(data)
}
