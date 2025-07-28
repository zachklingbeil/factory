package zero

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type One template.HTML

type Zero interface {
	Pathless(cssPath string)
	Build(class string, elements ...One) *One
	BuildFrame(class string, elements ...One)
	JS(js string) One
	CSS(css string) One
	GetPathless() *One
	GetFrame(index int) *One
	FrameCount() string
	CoordinatePlane()
	InjectFrame(index int) string

	Text
	Element
	Keybind
	Embed
	Universe
}

// --- zero Implementation ---
type zero struct {
	*text
	*element
	*keybind
	*embed
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
		embed:       NewEmbed().(*embed),
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

func (f *zero) InjectFrame(index int) string {
	if index < 0 || index >= len(f.frames) {
		return ""
	}
	return string(*f.frames[index])
}
func (f *zero) BuildFrame(class string, elements ...One) {
	zero := f.Build(class, elements...)
	f.frames = append(f.frames, zero)
	f.count++
}

func (f *zero) GetPathless() *One {
	return f.pathless
}

func (f *zero) Pathless(cssPath string, frameIndex int) {
	var html strings.Builder
	html.WriteString("<!DOCTYPE html>\n")
	html.WriteString("<html lang=\"en\">\n")
	html.WriteString("<head>\n")
	html.WriteString("  <meta charset=\"UTF-8\" />\n")
	html.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n")
	html.WriteString("  <title>hello universe</title>\n")

	html.WriteString("  <style>\n")
	html.WriteString(f.OneCSS(cssPath))
	html.WriteString("\n  </style>\n")

	html.WriteString("  <script>\n")
	html.WriteString(f.OneJS())
	html.WriteString("\n  </script>\n")

	html.WriteString("</head>\n")
	html.WriteString("<body>\n")
	html.WriteString(f.InjectFrame(frameIndex)) // inject the frame here
	html.WriteString("</body>\n")
	html.WriteString("</html>")
	result := One(template.HTML(html.String()))
	f.pathless = &result
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

func (f *zero) CoordinatePlane() {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(f.CoordinateCSS())
	b.WriteString(`</style>`)
	b.WriteString(`<div class="coordinate-plane" id="coordinate-plane"></div>`)
	b.WriteString(`<script>
document.addEventListener('DOMContentLoaded', async () => {
    try {
        const response = await fetch('/api/test');
        const data = await response.json();
        const plane = new CoordinatePlane(
            document.getElementById('coordinate-plane')
        );
        plane.initFromJson(data);
    } catch (err) {
        console.error('Failed to load test.json:', err);
    }
});
</script>`)
	one := One(template.HTML(b.String()))
	f.frames = append(f.frames, &one)
	f.count++
}
