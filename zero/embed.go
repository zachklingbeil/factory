package zero

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"os"
)

//go:embed embed/coordinate.html
var coordinatePlane string

//go:embed embed/coordinate.css
var coordinateCSS string

//go:embed embed/one.js
var oneJS string

//go:embed embed/one.css
var oneCSS string

//go:embed embed/test.json
var testJSON string

type Embed interface {
	NewCoordinatePlane() *template.HTML
	OneJS() string
	OneCSS(path string) string
	TestJSON() []Coordinate
	CoordinateCSS() string
}

type embed struct{}

func NewEmbed() Embed {
	return &embed{}
}

// CoordinatePlane returns the embedded coordinate plane HTML
func (a *embed) NewCoordinatePlane() *template.HTML {
	return (*template.HTML)(&coordinatePlane)
}

func (a *embed) OneJS() string { return oneJS }

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

// UnmarshalTestJSON loads the embedded testJSON into a []Coordinate
func (a *embed) TestJSON() []Coordinate {
	if testJSON == "" {
		return nil
	}
	var coords []Coordinate
	json.Unmarshal([]byte(testJSON), &coords)
	return coords
}

// CoordinateCss returns the embedded coordinate CSS
func (a *embed) CoordinateCSS() string {
	if coordinateCSS == "" {
		return ""
	}
	return coordinateCSS
}
