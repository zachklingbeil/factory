package zero

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type Frame interface {
	Pathless()
	Build(class string, elements ...One) *One
	BuildFrame(class string, elements ...One)
	JS(js string) One
	CSS(css string) One
	GetPathless() *One
	GetFrame(index int) *One
	FrameCount() string
	CoordinatePlane() One
	Text
	Element
	Keybind
	Embed
}

// --- frame Implementation ---
type frame struct {
	*text
	*element
	*keybind
	*embed
	frames   []*One
	count    uint
	pathless *One
}

func NewFrame() Frame {
	return &frame{
		text:     NewText().(*text),
		element:  NewElement().(*element),
		keybind:  &keybind{},
		embed:    NewEmbed().(*embed),
		frames:   make([]*One, 0),
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

func (f *frame) Pathless() {
	var html strings.Builder
	html.WriteString("<!DOCTYPE html>\n")
	html.WriteString("<html lang=\"en\">\n")
	html.WriteString("<head>\n")
	html.WriteString("  <meta charset=\"UTF-8\" />\n")
	html.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n")
	html.WriteString("  <title>hello universe</title>\n")

	html.WriteString("  <style>\n")
	html.WriteString(f.OneCSS())
	html.WriteString("\n  </style>\n")

	html.WriteString("  <script>\n")
	html.WriteString(f.OneJS())
	html.WriteString("\n  </script>\n")

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

func (f *frame) CoordinatePlane() One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(f.CoordinateCSS())
	b.WriteString(`</style>`)
	b.WriteString(`<div class="coordinate-plane" id="coordinate-plane"></div>`)
	b.WriteString(`<script>`)
	b.WriteString(f.CoordinateJS())
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}
