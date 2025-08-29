package zero

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"os"
	"strings"

	math "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	h "github.com/yuin/goldmark/renderer/html"
)

//go:embed pathless.html
var pathless string

type Build interface {
	Pathless() *One
	Lego(class string, elements ...One) *One
	JS(js string) One
	CSS(css string) One
	AddKeybind(containerId string, keyHandlers map[string]string) *One
	AddMarkdown(file string) *One
}

func NewBuild() Build {
	b := &build{
		element: NewElement().(*element),
		Md:      initGoldmark(),
	}
	one := One(template.HTML(pathless))
	b.pathless = &one
	return b
}

func (f *build) AddMarkdown(file string) *One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		empty := One("")
		return &empty
	}

	result := One(template.HTML(buf.String()))
	return &result
}

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, math.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
		),
		goldmark.WithRendererOptions(
			h.WithHardWraps(),
			h.WithXHTML(),
		),
	)
	return &md
}

// Add this field to the build struct:
type build struct {
	*element
	pathless *One
	Md       *goldmark.Markdown
}

func (f *build) Pathless() *One {
	return f.pathless

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
