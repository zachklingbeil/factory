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
	frame := f.Build(class, elements...)
	f.frames = append(f.frames, frame)
	f.count++
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
	html.WriteString("<!DOCTYPE html>\n")
	html.WriteString("<html lang=\"en\">\n")
	html.WriteString("<head>\n")
	html.WriteString("  <meta charset=\"UTF-8\" />\n")
	html.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n")
	html.WriteString("  <title>hello universe</title>\n")
	if c != "" {
		html.WriteString("  <style>\n")
		html.WriteString(c)
		html.WriteString("\n  </style>\n")
	}
	if j != "" {
		html.WriteString("  <script>\n")
		html.WriteString(j)
		html.WriteString("\n  </script>\n")
	}
	html.WriteString("</head>\n")
	html.WriteString("<body>\n")
	html.WriteString("</body>\n")
	html.WriteString("</html>")
	result := One(template.HTML(html.String()))
	f.pathless = &result
}

func (f *frame) FrameCount() string {
	return strconv.Itoa(int(f.count) - 1)
}

// Retrieve a frame by index
func (f *frame) GetFrame(index int) *One {
	return f.frames[index]
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
