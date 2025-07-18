package frame

import (
	"html/template"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Frame struct {
	Md *goldmark.Markdown
}

func NewFrame(mux *mux.Router) *Frame {
	frame := &Frame{
		Md: initGoldmark(),
	}
	return frame
}

func (f *Frame) AddFrame(reference string, elements ...template.HTML) *template.HTML {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	result := builder.String()
	if reference != "" {
		result = `<div class="` + reference + `">` + result + `</div>`
	}
	frame := template.HTML(result)
	return &frame
}

func (f *Frame) AddCSS(styles map[string]string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	html := template.HTML(builder.String())
	return &html
}

func (f *Frame) AddJS(js string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	html := template.HTML(builder.String())
	return &html
}
