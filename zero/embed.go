package zero

import (
	_ "embed"
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
	OneCSS() string
}

type embed struct{}

func NewEmbed() Embed {
	return &embed{}
}

func (a *embed) CoordinateCSS() string { return coordinateCSS }
func (a *embed) CoordinateJS() string  { return coordinateJS }
func (a *embed) OneJS() string         { return oneJS }
func (a *embed) OneCSS() string        { return oneCSS }
