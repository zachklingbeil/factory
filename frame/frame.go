package frame

import (
	"html/template"
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

func (f *Frame) AddCSS(styles map[string]string) template.CSS {
	var builder strings.Builder
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	return template.CSS(builder.String())
}

func (f *Frame) AddJS(js string) template.JS {
	return template.JS(js)
}
