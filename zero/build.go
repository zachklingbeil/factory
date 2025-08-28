package zero

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type Build interface {
	Lego(class string, elements ...One) *One
	Build(class string, elements ...One)
	JS(js string) One
	CSS(css string) One
	GetFrame(index int) *One
	FrameCount() string
	AddKeybind(containerId string, keyHandlers map[string]string) *One
	Text
	Element
	Universe
}

type build struct {
	*text
	*element
	*coordinates
	frames []*One
	count  uint
}

func NewBuild() Build {
	return &build{
		text:        NewText().(*text),
		element:     NewElement().(*element),
		coordinates: NewUniverse().(*coordinates),
		frames:      make([]*One, 0),
		count:       0,
	}
}

func (f *build) Lego(class string, elements ...One) *One {
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

func (f *build) Build(class string, elements ...One) {
	zero := f.Lego(class, elements...)
	f.frames = append(f.frames, zero)
	f.count++
}

func (f *build) FrameCount() string {
	return strconv.Itoa(int(f.count) - 1)
}

// Retrieve a zero by index
func (f *build) GetFrame(index int) *One {
	return f.frames[index]
}

func (f *build) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *build) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}

func (f *build) AddKeybind(containerId string, keyHandlers map[string]string) *One {
	var handlers strings.Builder
	for key, handlerCode := range keyHandlers {
		handlers.WriteString(fmt.Sprintf(`
         if (event.key === %q) {
            %s
         }
        `, key, handlerCode))
	}
	js := fmt.Sprintf(`
document.addEventListener('DOMContentLoaded', () => {
   const container = document.getElementById(%q);
   if (!container) return;
   container.tabIndex = 0;
   container.addEventListener('keydown', (event) => {
      %s
   });
});
`, containerId, handlers.String())
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}
