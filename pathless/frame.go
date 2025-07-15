package pathless

import (
	"html/template"
	"strings"
)

func (p *Pathless) CreateFrame(elements ...template.HTML) *template.HTML {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	return &html
}

func (p *Pathless) AddCSS(frame *template.HTML, styles map[string]string) *template.HTML {
	var builder strings.Builder
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

func (p *Pathless) AddJS(frame *template.HTML, js string) *template.HTML {
	var builder strings.Builder
	builder.WriteString(string(*frame))
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	html := template.HTML(builder.String())
	return &html
}
