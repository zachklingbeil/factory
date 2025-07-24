package zero

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type Frame interface {
	Pathless(css, js string)
	GetPathless() *One
	Build(class string, elements ...One) *One
	BuildFrame(class string, elements ...One)
	JS(js string) One
	CSS(css string) One
	AddFrame(frame *One)
	GetFrame(index int) (*One, bool)
	FrameCount() string
	Text
	Element
	Keybind
}

// --- frame Implementation ---
type frame struct {
	Text
	Element
	Keybind
	frames   []*One
	count    uint
	pathless *One
}

func NewFrame() Frame {
	return &frame{
		Text:     NewText(),
		Element:  NewElement(),
		frames:   make([]*One, 0),
		Keybind:  &keybind{},
		count:    0,
		pathless: nil,
	}
}

func (f *frame) Build(class string, elements ...One) *One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}

	if class == "" {
		result := One(template.HTML(b.String()))
		return &result
	}

	consolidatedContent := template.HTML(b.String())
	htmlResult := fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), string(consolidatedContent))
	result := One(template.HTML(htmlResult))
	return &result
}

func (f *frame) BuildFrame(class string, elements ...One) {
	f.AddFrame(f.Build(class, elements...))
}

func (f *frame) GetPathless() *One {
	return f.pathless
}

func (f *frame) Pathless(css, js string) {
	var c, j string
	if css != "" {
		c = f.FileToString(css)
	}
	if js != "" {
		j = f.FileToString(js)
	}

	var html strings.Builder
	html.WriteString(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8" /><meta name="viewport" content="width=device-width, initial-scale=1.0" /><title>hello universe</title>`)

	if c != "" {
		html.WriteString(`<style>`)
		html.WriteString(c)
		html.WriteString(`</style>`)
	}
	if j != "" {
		html.WriteString(`<script>`)
		html.WriteString(j)
		html.WriteString(`</script>`)
	}
	html.WriteString(`</head></html>`)
	result := One(template.HTML(html.String()))
	f.pathless = &result
}

func (f *frame) FrameCount() string {
	return strconv.Itoa(int(f.count))
}

// Add a finalized frame to the collection
func (f *frame) AddFrame(frame *One) {
	f.frames = append(f.frames, frame)
	f.count++
}

// Retrieve a frame by index
func (f *frame) GetFrame(index int) (*One, bool) {
	if index < 0 || index >= int(f.count) {
		return nil, false
	}
	return f.frames[index], true
}

func (f *frame) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *frame) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}
