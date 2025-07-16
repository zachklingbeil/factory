package frame

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

type Frame struct {
	Md *goldmark.Markdown
}

func NewFrame() *Frame {
	return &Frame{
		Md: initGoldmark(),
	}
}

func (f *Frame) CreateFrame(elements ...template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}

	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}

	return template.HTML(builder.String())
}

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return template.HTML("")
	}
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, template.HTML(buf.String()))
	allElements = append(allElements, elements...)
	return f.CreateFrame(allElements...)
}

func (f *Frame) AddCSS(styles map[string]string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	return template.HTML(builder.String())
}

func (f *Frame) AddJS(js string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	return template.HTML(builder.String())
}
