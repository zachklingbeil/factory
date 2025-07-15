package frame

import (
	"html/template"
	"strings"
)

type Frame struct{}

func (f *Frame) CreateFrame(elements ...template.HTML) *template.HTML {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	return &html
}

func (f *Frame) AddCSS(frame *template.HTML, styles map[string]string) *template.HTML {
	var builder strings.Builder

	// Start with existing HTML
	builder.WriteString(string(*frame))
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

func (f *Frame) AddJS(frame *template.HTML, js string) *template.HTML {
	var builder strings.Builder
	// Start with existing HTML
	builder.WriteString(string(*frame))
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")

	html := template.HTML(builder.String())
	return &html
}
