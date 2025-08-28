package zero

import (
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type One template.HTML

//go:embed embed/index.html
var pathless string

// Pathless returns the embedded index.html as *One
func (z *zero) Pathless() *One {
	result := One(template.HTML(pathless))
	return &result
}

type Zero interface {
	Pathless() *One
	Build(class string, elements ...One) *One
	BuildFrame(class string, elements ...One)
	JS(js string) One
	CSS(css string) One
	GetPathless() *One
	GetFrame(index int) *One
	FrameCount() string
	Text
	Element
	Keybind
	Universe
}

// --- zero Implementation ---
type zero struct {
	*text
	*element
	*keybind
	*coordinates
	frames   []*One
	count    uint
	pathless *One
}

func NewZero() Zero {
	return &zero{
		text:        NewText().(*text),
		element:     NewElement().(*element),
		keybind:     &keybind{},
		frames:      make([]*One, 0),
		coordinates: NewUniverse().(*coordinates),
		count:       0,
		pathless:    nil,
	}
}

func (f *zero) Build(class string, elements ...One) *One {
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

func (f *zero) BuildFrame(class string, elements ...One) {
	zero := f.Build(class, elements...)
	f.frames = append(f.frames, zero)
	f.count++
}

func (f *zero) GetPathless() *One {
	return f.pathless
}

func (f *zero) FrameCount() string {
	return strconv.Itoa(int(f.count) - 1)
}

// Retrieve a zero by index
func (f *zero) GetFrame(index int) *One {
	return f.frames[index]
}

func (f *zero) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *zero) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}
